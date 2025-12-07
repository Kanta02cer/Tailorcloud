package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"tailor-cloud/backend/internal/config"
	"tailor-cloud/backend/internal/handler"
	"tailor-cloud/backend/internal/logger"
	"tailor-cloud/backend/internal/metrics"
	"tailor-cloud/backend/internal/middleware"
	"tailor-cloud/backend/internal/repository"
	"tailor-cloud/backend/internal/service"
)

func main() {
	ctx := context.Background()

	// 1. Initialize PostgreSQL (Primary DB)
	// 仕様書要件: Primary DB (RDBMS): PostgreSQL - 注文データ、顧客台帳、会計データ、決済トランザクション
	dbConfig := config.LoadDatabaseConfig()

	// PostgreSQL接続（オプショナル: 接続失敗時は警告のみ）
	db, err := config.ConnectPostgreSQL(dbConfig)
	if err != nil {
		log.Printf("WARNING: Failed to connect to PostgreSQL: %v", err)
		log.Printf("WARNING: Continuing with Firestore only mode (not recommended for production)")
		db = nil
	} else {
		log.Println("PostgreSQL connection established successfully")

		// 接続プール設定を適用（パフォーマンス最適化）
		if err := config.ConfigurePool(db, config.DefaultPoolConfig()); err != nil {
			log.Printf("WARNING: Failed to configure connection pool: %v", err)
		} else {
			log.Println("Database connection pool configured")
		}

		defer db.Close()
	}

	// 2. Initialize Firebase/Firestore (Secondary DB)
	// 仕様書要件: Secondary DB (NoSQL): Firestore - 案件チャットログ、一時的なUIステータス、通知バッジ
	var app *firebase.App
	gcpProjectID := os.Getenv("GCP_PROJECT_ID")
	serviceAccountJSON := os.Getenv("FIREBASE_SERVICE_ACCOUNT_JSON")

	if gcpProjectID != "" {
		var opts []option.ClientOption

		// 認証情報の読み込み（Render環境とローカル環境の両方に対応）
		if serviceAccountJSON != "" {
			// Render環境: 環境変数からJSON認証情報を読み込む
			log.Println("Loading Firebase credentials from FIREBASE_SERVICE_ACCOUNT_JSON environment variable")

			// 改行文字の修正: Renderの環境変数でエスケープされた \n を実際の改行コードに置換
			// これを行わないと "private key should be a PEM or plain PKCS1 or PKCS8" エラーになる場合がある
			if strings.Contains(serviceAccountJSON, "\\n") {
				log.Println("Detected escaped newlines in Firebase credentials, replacing them...")
				serviceAccountJSON = strings.ReplaceAll(serviceAccountJSON, "\\n", "\n")
			}

			opts = append(opts, option.WithCredentialsJSON([]byte(serviceAccountJSON)))
		} else if credPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); credPath != "" {
			// ローカル環境: ファイルパスから読み込む
			log.Printf("Loading Firebase credentials from file: %s", credPath)
			opts = append(opts, option.WithCredentialsFile(credPath))
		} else {
			// Application Default Credentials (ADC) を使用
			log.Println("Using Application Default Credentials for Firebase")
		}

		conf := &firebase.Config{ProjectID: gcpProjectID}
		firebaseApp, err := firebase.NewApp(ctx, conf, opts...)
		if err != nil {
			log.Printf("WARNING: Failed to initialize Firebase app: %v", err)
			log.Printf("WARNING: Continuing without Firebase (authentication features will not work)")
			app = nil
		} else {
			app = firebaseApp
			log.Println("Firebase app initialized successfully")
		}
	} else {
		log.Println("WARNING: GCP_PROJECT_ID not set. Firebase features will not be available.")
		app = nil
	}

	var firestoreClient *firestore.Client
	if app != nil {
		client, err := app.Firestore(ctx)
		if err != nil {
			log.Printf("WARNING: Failed to initialize Firestore client: %v", err)
		} else {
			firestoreClient = client
			defer client.Close()
		}
	}

	// 3. Dependency Injection (Wiring)
	// Repository -> Service -> Handler の順で依存性を注入

	// 注文リポジトリ: PostgreSQLを使用（Primary DB）
	// 仕様書準拠: 注文データはPostgreSQLに保存
	var orderRepo repository.OrderRepository
	if db != nil {
		orderRepo = repository.NewPostgreSQLOrderRepository(db)
		log.Println("Using PostgreSQL for orders (Primary DB)")
	} else {
		// フォールバック: PostgreSQLが利用できない場合はFirestoreを使用（開発環境用）
		if firestoreClient != nil {
			orderRepo = repository.NewFirestoreOrderRepository(firestoreClient)
			log.Println("WARNING: Using Firestore for orders (fallback mode - not recommended for production)")
		} else {
			// 開発環境: データベースがなくても起動可能（認証エンドポイントなどは動作）
			log.Println("WARNING: Neither PostgreSQL nor Firestore is available.")
			log.Println("WARNING: Order-related endpoints will not work, but authentication endpoints will work.")
			log.Println("WARNING: This is acceptable for development/testing purposes.")
			orderRepo = nil
		}
	}

	// 監査ログリポジトリ: PostgreSQLを使用
	var auditLogRepo repository.AuditLogRepository
	if db != nil {
		auditLogRepo = repository.NewPostgreSQLAuditLogRepository(db)
		// complianceViewLogRepoは将来の契約書閲覧ログ機能で使用予定
		// _ = repository.NewPostgreSQLComplianceDocumentViewLogRepository(db)
		log.Println("Audit log repository initialized")
	} else {
		log.Println("WARNING: Audit logging disabled (PostgreSQL not available)")
	}

	// 3.1 Initialize Firebase Auth Middleware
	var authMiddleware *middleware.FirebaseAuthMiddleware
	if app != nil {
		authMw, err := middleware.NewFirebaseAuthMiddleware(app)
		if err != nil {
			log.Printf("WARNING: Failed to initialize Firebase Auth middleware: %v", err)
			log.Printf("WARNING: Continuing without authentication (not recommended for production)")
		} else {
			authMiddleware = authMw
			log.Println("Firebase Auth middleware initialized")
		}
	}

	// 権限リポジトリ
	var permissionRepo repository.PermissionRepository
	if db != nil {
		permissionRepo = repository.NewPostgreSQLPermissionRepository(db)
		log.Println("Permission repository initialized")
	}

	// RBACサービス
	var rbacService *service.RBACService
	if permissionRepo != nil {
		rbacService = service.NewRBACService(permissionRepo)
		log.Println("RBAC service initialized")
	}

	// RBACミドルウェア
	rbacMiddleware := middleware.NewRBACMiddleware(rbacService)

	// 構造化ロガーを初期化
	structuredLogger := logger.NewStructuredLogger(
		logger.WithService("tailorcloud-backend"),
		logger.WithLevel(logger.LogLevelInfo),
	)
	log.Println("Structured logger initialized")

	// メトリクスコレクターを初期化
	metricsCollector := metrics.NewMetricsCollector()
	log.Println("Metrics collector initialized")

	// トレースミドルウェアを初期化
	traceMiddleware := middleware.NewTraceMiddleware(structuredLogger)
	log.Println("Trace middleware initialized")

	// ロギングミドルウェアを初期化
	loggingMiddleware := middleware.NewLoggingMiddleware(structuredLogger)
	log.Println("Logging middleware initialized")

	// メトリクスミドルウェアを初期化
	metricsMiddleware := middleware.NewMetricsMiddleware(metricsCollector)
	log.Println("Metrics middleware initialized")

	// 生地リポジトリ: PostgreSQLを使用
	var fabricRepo repository.FabricRepository
	if db != nil {
		fabricRepo = repository.NewPostgreSQLFabricRepository(db)
		log.Println("Fabric repository initialized")
	}

	// アンバサダーリポジトリ: PostgreSQLを使用
	var ambassadorRepo repository.AmbassadorRepository
	var commissionRepo repository.CommissionRepository
	if db != nil {
		ambassadorRepo = repository.NewPostgreSQLAmbassadorRepository(db)
		commissionRepo = repository.NewPostgreSQLCommissionRepository(db)
		log.Println("Ambassador repositories initialized")
	}

	// 顧客リポジトリ: PostgreSQLを使用
	var customerRepo repository.CustomerRepository
	if db != nil {
		customerRepo = repository.NewPostgreSQLCustomerRepository(db)
		log.Println("Customer repository initialized")
	}

	// テナントリポジトリ: PostgreSQLを使用
	var tenantRepo repository.TenantRepository
	if db != nil {
		tenantRepo = repository.NewPostgreSQLTenantRepository(db)
		log.Println("Tenant repository initialized")
	}

	// ユーザーリポジトリ: PostgreSQLを使用
	var userRepo repository.UserRepository
	if db != nil {
		userRepo = repository.NewPostgreSQLUserRepository(db)
		log.Println("User repository initialized")
	}

	// 反物（Roll）リポジトリ: PostgreSQLを使用
	var fabricRollRepo repository.FabricRollRepository
	if db != nil {
		fabricRollRepo = repository.NewPostgreSQLFabricRollRepository(db)
		log.Println("Fabric roll repository initialized")
	}

	// 反物引当リポジトリ: PostgreSQLを使用
	var fabricAllocationRepo repository.FabricAllocationRepository
	if db != nil {
		fabricAllocationRepo = repository.NewPostgreSQLFabricAllocationRepository(db)
		log.Println("Fabric allocation repository initialized")
	}

	// 診断リポジトリ: PostgreSQLを使用（Suit-MBTI統合）
	var diagnosisRepo repository.DiagnosisRepository
	if db != nil {
		diagnosisRepo = repository.NewPostgreSQLDiagnosisRepository(db)
		log.Println("Diagnosis repository initialized")
	}

	// 予約リポジトリ: PostgreSQLを使用（Suit-MBTI統合）
	var appointmentRepo repository.AppointmentRepository
	if db != nil {
		appointmentRepo = repository.NewPostgreSQLAppointmentRepository(db)
		log.Println("Appointment repository initialized")
	}

	// サービス層の依存性注入
	// アンバサダーサービス（成果報酬管理用）
	var ambassadorService *service.AmbassadorService
	if ambassadorRepo != nil && commissionRepo != nil {
		ambassadorService = service.NewAmbassadorService(ambassadorRepo, commissionRepo)
		log.Println("Ambassador service initialized")
	}

	// 注文サービス: 監査ログリポジトリとアンバサダーサービスを注入
	// orderRepoがnilの場合はnil
	var orderService *service.OrderService
	if orderRepo != nil {
		orderService = service.NewOrderService(orderRepo, auditLogRepo, ambassadorService)
		log.Println("Order service initialized")
	} else {
		log.Println("WARNING: Order service not initialized (no database available)")
	}

	// 生地サービス
	var fabricService *service.FabricService
	if fabricRepo != nil {
		fabricService = service.NewFabricService(fabricRepo)
	}

	// 顧客サービス
	var customerService *service.CustomerService
	if customerRepo != nil && orderRepo != nil {
		customerService = service.NewCustomerService(customerRepo, orderRepo)
		log.Println("Customer service initialized")
	}

	// アナリティクスサービス
	var analyticsService *service.AnalyticsService
	if orderRepo != nil && customerRepo != nil {
		analyticsService = service.NewAnalyticsService(orderRepo, customerRepo)
		log.Println("Analytics service initialized")
	}

	// 在庫引当サービス（エンタープライズ実装の核心）
	var inventoryAllocationService *service.InventoryAllocationService
	if fabricRollRepo != nil && fabricAllocationRepo != nil && fabricRepo != nil && db != nil {
		inventoryAllocationService = service.NewInventoryAllocationService(
			fabricRollRepo,
			fabricAllocationRepo,
			fabricRepo,
			db,
		)
		log.Println("Inventory allocation service initialized")
	}

	// Cloud Storageサービス（PDF保存用）
	var storageService service.StorageService
	bucketName := os.Getenv("GCS_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "tailorcloud-compliance-docs" // デフォルトバケット名
		log.Println("WARNING: GCS_BUCKET_NAME not set, using default:", bucketName)
	}

	// StorageServiceの初期化（認証情報ファイルは環境変数から取得）
	credentialsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	gcsStorage, err := service.NewGCSStorageService(ctx, credentialsPath)
	if err != nil {
		log.Printf("WARNING: Failed to initialize Cloud Storage service: %v", err)
		log.Printf("WARNING: PDF upload will fail, but PDF generation will work")
	} else {
		storageService = gcsStorage
		defer gcsStorage.Close()
		log.Println("Cloud Storage service initialized")
	}

	// コンプライアンス文書リポジトリ
	var complianceDocRepo repository.ComplianceDocumentRepository
	if db != nil {
		complianceDocRepo = repository.NewPostgreSQLComplianceDocumentRepository(db)
		log.Println("Compliance document repository initialized")
	}

	// コンプライアンスサービス（PDF生成用）
	var complianceService *service.ComplianceService
	if complianceDocRepo != nil {
		complianceService = service.NewComplianceService(storageService, bucketName, complianceDocRepo)
		log.Println("Compliance service initialized")
	} else {
		// リポジトリがない場合はnilで作成（履歴管理なし）
		complianceService = service.NewComplianceService(storageService, bucketName, nil)
		log.Println("Compliance service initialized (without history management)")
	}

	// 税率計算サービス（インボイス制度対応）
	var taxService *service.TaxCalculationService
	if tenantRepo != nil {
		taxService = service.NewTaxCalculationService(tenantRepo)
		log.Println("Tax calculation service initialized")
	}

	// 請求書サービス（インボイスPDF生成用）
	var invoiceService *service.InvoiceService
	if orderRepo != nil && tenantRepo != nil && customerRepo != nil && storageService != nil && taxService != nil {
		invoiceService = service.NewInvoiceService(
			orderRepo,
			tenantRepo,
			customerRepo,
			storageService,
			bucketName,
			taxService,
		)
		log.Println("Invoice service initialized")
	}

	// 診断サービス（Suit-MBTI統合）
	var diagnosisService *service.DiagnosisService
	if diagnosisRepo != nil {
		diagnosisService = service.NewDiagnosisService(diagnosisRepo)
		log.Println("Diagnosis service initialized")
	}

	// 予約サービス（Suit-MBTI統合）
	var appointmentService *service.AppointmentService
	if appointmentRepo != nil {
		appointmentService = service.NewAppointmentService(appointmentRepo)
		log.Println("Appointment service initialized")
	}

	// 自動補正エンジンサービス（The "Auto Patterner"）
	var measurementCorrectionService *service.MeasurementCorrectionService
	if diagnosisService != nil && fabricRepo != nil {
		measurementCorrectionService = service.NewMeasurementCorrectionService(
			diagnosisService,
			fabricRepo,
		)
		log.Println("Measurement correction service (Auto Patterner) initialized")
	}

	// 採寸データバリデーションサービス
	var measurementValidationService *service.MeasurementValidationService
	if orderRepo != nil {
		measurementValidationService = service.NewMeasurementValidationService(orderRepo)
		log.Println("Measurement validation service initialized")
	}

	// ユーザーサービス
	var userService *service.UserService
	if userRepo != nil && tenantRepo != nil {
		userService = service.NewUserService(userRepo, tenantRepo)
		log.Println("User service initialized")
	}

	// ハンドラー
	// orderServiceがnilの場合はnil
	var orderHandler *handler.OrderHandler
	if orderService != nil {
		orderHandler = handler.NewOrderHandler(orderService)
		log.Println("Order handler initialized")
	} else {
		log.Println("WARNING: Order handler not initialized (order service not available)")
	}

	// 生地ハンドラー
	var fabricHandler *handler.FabricHandler
	if fabricService != nil {
		fabricHandler = handler.NewFabricHandler(fabricService)
	}

	// アンバサダーハンドラー
	var ambassadorHandler *handler.AmbassadorHandler
	if ambassadorService != nil {
		ambassadorHandler = handler.NewAmbassadorHandler(ambassadorService)
	}

	// コンプライアンスハンドラー（発注書生成用）
	// complianceServiceまたはorderServiceがnilの場合はnil
	var complianceHandler *handler.ComplianceHandler
	if complianceService != nil && orderService != nil {
		complianceHandler = handler.NewComplianceHandler(complianceService, orderService)
		log.Println("Compliance handler initialized")
	} else {
		log.Println("WARNING: Compliance handler not initialized (compliance or order service not available)")
	}

	// 顧客ハンドラー
	var customerHandler *handler.CustomerHandler
	if customerService != nil {
		customerHandler = handler.NewCustomerHandler(customerService)
		log.Println("Customer handler initialized")
	} else {
		log.Println("WARNING: Customer handler not initialized (customer service not available)")
	}

	// アナリティクスハンドラー
	var analyticsHandler *handler.AnalyticsHandler
	if analyticsService != nil {
		analyticsHandler = handler.NewAnalyticsHandler(analyticsService)
		log.Println("Analytics handler initialized")
	} else {
		log.Println("WARNING: Analytics handler not initialized (analytics service not available)")
	}

	// 反物（Roll）ハンドラー
	var fabricRollHandler *handler.FabricRollHandler
	if fabricRollRepo != nil {
		fabricRollHandler = handler.NewFabricRollHandler(fabricRollRepo)
		log.Println("Fabric roll handler initialized")
	}

	// 在庫引当ハンドラー
	var inventoryAllocationHandler *handler.InventoryAllocationHandler
	if inventoryAllocationService != nil {
		inventoryAllocationHandler = handler.NewInventoryAllocationHandler(inventoryAllocationService)
		log.Println("Inventory allocation handler initialized")
	}

	// 請求書ハンドラー（インボイスPDF生成用）
	var invoiceHandler *handler.InvoiceHandler
	if invoiceService != nil {
		invoiceHandler = handler.NewInvoiceHandler(invoiceService)
		log.Println("Invoice handler initialized")
	}

	// 権限ハンドラー
	var permissionHandler *handler.PermissionHandler
	if rbacService != nil {
		permissionHandler = handler.NewPermissionHandler(rbacService)
		log.Println("Permission handler initialized")
	}

	// 診断ハンドラー（Suit-MBTI統合）
	var diagnosisHandler *handler.DiagnosisHandler
	if diagnosisService != nil {
		diagnosisHandler = handler.NewDiagnosisHandler(diagnosisService)
		log.Println("Diagnosis handler initialized")
	}

	// 予約ハンドラー（Suit-MBTI統合）
	var appointmentHandler *handler.AppointmentHandler
	if appointmentService != nil {
		appointmentHandler = handler.NewAppointmentHandler(appointmentService)
		log.Println("Appointment handler initialized")
	}

	// 自動補正エンジンハンドラー（The "Auto Patterner"）
	var measurementCorrectionHandler *handler.MeasurementCorrectionHandler
	if measurementCorrectionService != nil {
		measurementCorrectionHandler = handler.NewMeasurementCorrectionHandler(measurementCorrectionService)
		log.Println("Measurement correction handler (Auto Patterner) initialized")
	}

	// 採寸データバリデーションハンドラー
	var measurementValidationHandler *handler.MeasurementValidationHandler
	if measurementValidationService != nil {
		measurementValidationHandler = handler.NewMeasurementValidationHandler(measurementValidationService)
		log.Println("Measurement validation handler initialized")
	}

	// 認証ハンドラー（Firebase Auth用）
	var authEndpointHandler *handler.AuthHandler
	if app != nil {
		authHdlr, err := handler.NewAuthHandler(app, userService)
		if err != nil {
			log.Printf("WARNING: Failed to initialize Auth handler: %v", err)
		} else {
			authEndpointHandler = authHdlr
			log.Println("Auth handler initialized")
		}
	}

	// 4. Routing
	mux := http.NewServeMux()

	// Order endpoints with authentication
	// 環境変数ENVに基づいて認証モードを切り替え
	var authHandler func(http.HandlerFunc) http.HandlerFunc
	if authMiddleware != nil {
		env := os.Getenv("ENV")
		if env == "production" {
			// 本番環境: 認証必須（Authenticate）
			authHandler = authMiddleware.Authenticate
			log.Println("Using Authenticate middleware (production mode - authentication required)")
		} else {
			// 開発環境: 認証オプショナル（OptionalAuth）
			authHandler = authMiddleware.OptionalAuth
			log.Println("Using OptionalAuth middleware (development mode - authentication optional)")
		}
	} else {
		// 認証ミドルウェアがない場合はパススルー
		authHandler = func(next http.HandlerFunc) http.HandlerFunc {
			return next
		}
		log.Println("WARNING: Authentication middleware not available - all requests will pass through")
	}

	// CORSミドルウェアをインポート
	// 開発環境ではすべてのオリジンを許可
	corsMiddleware := middleware.CORS

	// ミドルウェアチェーンを作成
	// 順序: CORS -> Trace -> Logging -> Metrics -> Auth -> RBAC -> Handler
	chainMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		// CORS -> トレースIDを付与 -> ロギング -> メトリクス収集
		return corsMiddleware(
			traceMiddleware.Trace(
				loggingMiddleware.Log(
					metricsMiddleware.Collect(next),
				),
			),
		)
	}

	// 認証付きミドルウェアチェーン
	authChainMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return chainMiddleware(authHandler(next))
	}

	// Order endpoints（orderHandlerがnilの場合は登録しない）
	if orderHandler != nil {
		mux.HandleFunc("POST /api/orders", authChainMiddleware(orderHandler.CreateOrder))
		mux.HandleFunc("POST /api/orders/confirm", authChainMiddleware(rbacMiddleware.RequireOwnerOrStaff()(orderHandler.ConfirmOrder)))
		mux.HandleFunc("GET /api/orders", authChainMiddleware(func(w http.ResponseWriter, r *http.Request) {
			// order_idがあれば単一取得、なければ一覧取得
			if r.URL.Query().Get("order_id") != "" {
				orderHandler.GetOrder(w, r)
			} else {
				orderHandler.ListOrders(w, r)
			}
		}))
		log.Println("Order endpoints registered")
	} else {
		log.Println("WARNING: Order endpoints not registered (order service not available)")
	}

	// Compliance endpoints (PDF生成)
	// 注意: パスパターンは /api/orders/{id}/generate-document の形式
	// Go 1.22+ の新しいルーティング機能を使用
	if complianceHandler != nil {
		mux.HandleFunc("POST /api/orders/{id}/generate-document", authChainMiddleware(rbacMiddleware.RequireOwnerOrStaff()(complianceHandler.GenerateDocument)))
		mux.HandleFunc("POST /api/orders/{id}/generate-amendment", authChainMiddleware(rbacMiddleware.RequireOwnerOrStaff()(complianceHandler.GenerateAmendmentDocument)))
	}

	// Fabric (Inventory) endpoints
	if fabricHandler != nil {
		mux.HandleFunc("GET /api/fabrics", authChainMiddleware(fabricHandler.ListFabrics))
		mux.HandleFunc("GET /api/fabrics/detail", authChainMiddleware(fabricHandler.GetFabric))
		mux.HandleFunc("POST /api/fabrics/reserve", authChainMiddleware(fabricHandler.ReserveFabric))
	}

	// Ambassador endpoints
	if ambassadorHandler != nil {
		mux.HandleFunc("POST /api/ambassadors", authChainMiddleware(rbacMiddleware.RequireOwnerOnly()(ambassadorHandler.CreateAmbassador)))
		mux.HandleFunc("GET /api/ambassadors/me", authChainMiddleware(ambassadorHandler.GetAmbassadorByUserID))
		mux.HandleFunc("GET /api/ambassadors", authChainMiddleware(ambassadorHandler.ListAmbassadors))
		mux.HandleFunc("GET /api/ambassadors/commissions", authChainMiddleware(ambassadorHandler.GetCommissions))
	}

	// Customer endpoints (CRM)
	if customerHandler != nil {
		mux.HandleFunc("POST /api/customers", authChainMiddleware(rbacMiddleware.RequireOwnerOrStaff()(customerHandler.CreateCustomer)))
		mux.HandleFunc("GET /api/customers/{id}", authChainMiddleware(customerHandler.GetCustomer))
		mux.HandleFunc("GET /api/customers", authChainMiddleware(customerHandler.ListCustomers))
		mux.HandleFunc("PUT /api/customers/{id}", authChainMiddleware(rbacMiddleware.RequireOwnerOrStaff()(customerHandler.UpdateCustomer)))
		mux.HandleFunc("DELETE /api/customers/{id}", authChainMiddleware(rbacMiddleware.RequireOwnerOnly()(customerHandler.DeleteCustomer)))
		mux.HandleFunc("GET /api/customers/{id}/orders", authChainMiddleware(customerHandler.GetCustomerOrders))
		log.Println("Customer endpoints registered")
	} else {
		// 開発環境用: データベースなしでもエンドポイントを登録（モックレスポンス）
		mux.HandleFunc("GET /api/customers", authChainMiddleware(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"customers": []interface{}{},
				"total":     0,
				"message":   "Database not connected. This is a mock response for development.",
			})
		}))
		mux.HandleFunc("POST /api/customers", authChainMiddleware(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Database not connected. Customer service unavailable.",
			})
		}))
		log.Println("Customer endpoints registered (mock mode - database not connected)")
	}

	// Analytics endpoints
	if analyticsHandler != nil {
		mux.HandleFunc("GET /api/analytics/summary", authChainMiddleware(analyticsHandler.GetSummary))
		log.Println("Analytics endpoints registered")
	} else {
		// 開発環境用: データベースなしでもエンドポイントを登録（モックレスポンス）
		mux.HandleFunc("GET /api/analytics/summary", authChainMiddleware(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"total_customers":      0,
				"total_orders":         0,
				"total_revenue":        0.0,
				"revenue_trend":        []interface{}{},
				"top_customers":        []interface{}{},
				"order_status_summary": map[string]int{},
				"message":              "Database not connected. This is a mock response for development.",
			})
		}))
		log.Println("Analytics endpoints registered (mock mode - database not connected)")
	}

	// Fabric Roll (反物管理) endpoints
	if fabricRollHandler != nil {
		mux.HandleFunc("POST /api/fabric-rolls", authChainMiddleware(rbacMiddleware.RequireOwnerOrStaff()(fabricRollHandler.CreateFabricRoll)))
		mux.HandleFunc("GET /api/fabric-rolls/{id}", authChainMiddleware(fabricRollHandler.GetFabricRoll))
		mux.HandleFunc("GET /api/fabric-rolls", authChainMiddleware(fabricRollHandler.ListFabricRolls))
		mux.HandleFunc("PUT /api/fabric-rolls/{id}", authChainMiddleware(rbacMiddleware.RequireOwnerOrStaff()(fabricRollHandler.UpdateFabricRoll)))
	}

	// Inventory Allocation (在庫引当) endpoints
	if inventoryAllocationHandler != nil {
		mux.HandleFunc("POST /api/inventory/allocate", authChainMiddleware(rbacMiddleware.RequireOwnerOrStaff()(inventoryAllocationHandler.AllocateInventory)))
		mux.HandleFunc("POST /api/inventory/release", authChainMiddleware(rbacMiddleware.RequireOwnerOrStaff()(inventoryAllocationHandler.ReleaseAllocation)))
	}

	// Invoice (請求書・インボイス) endpoints
	if invoiceHandler != nil {
		mux.HandleFunc("POST /api/orders/{id}/generate-invoice", authChainMiddleware(rbacMiddleware.RequireOwnerOrStaff()(invoiceHandler.GenerateInvoice)))
	}

	// Permission (権限管理) endpoints
	if permissionHandler != nil {
		mux.HandleFunc("POST /api/permissions", authChainMiddleware(rbacMiddleware.RequireOwnerOnly()(permissionHandler.CreatePermission)))
		mux.HandleFunc("GET /api/permissions", authChainMiddleware(permissionHandler.GetPermissions))
		mux.HandleFunc("POST /api/permissions/check", authChainMiddleware(permissionHandler.CheckPermission))
	}

	// Diagnosis (診断) endpoints (Suit-MBTI統合)
	if diagnosisHandler != nil {
		mux.HandleFunc("POST /api/diagnoses", authChainMiddleware(diagnosisHandler.CreateDiagnosis))
		mux.HandleFunc("GET /api/diagnoses/{id}", authChainMiddleware(diagnosisHandler.GetDiagnosis))
		mux.HandleFunc("GET /api/diagnoses", authChainMiddleware(diagnosisHandler.ListDiagnoses))
	}

	// Measurement Correction (自動補正エンジン) endpoints
	if measurementCorrectionHandler != nil {
		mux.HandleFunc("POST /api/measurements/convert", authChainMiddleware(rbacMiddleware.RequireOwnerOrStaff()(measurementCorrectionHandler.ConvertToFinalMeasurements)))
	}

	// Measurement Validation (採寸データバリデーション) endpoints
	if measurementValidationHandler != nil {
		mux.HandleFunc("POST /api/measurements/validate", authChainMiddleware(rbacMiddleware.RequireOwnerOrStaff()(measurementValidationHandler.ValidateMeasurements)))
		mux.HandleFunc("POST /api/measurements/validate-range", authChainMiddleware(rbacMiddleware.RequireOwnerOrStaff()(measurementValidationHandler.ValidateMeasurementRange)))
	}

	// Appointment (予約) endpoints (Suit-MBTI統合)
	if appointmentHandler != nil {
		mux.HandleFunc("POST /api/appointments", authChainMiddleware(rbacMiddleware.RequireOwnerOrStaff()(appointmentHandler.CreateAppointment)))
		mux.HandleFunc("GET /api/appointments/{id}", authChainMiddleware(appointmentHandler.GetAppointment))
		mux.HandleFunc("GET /api/appointments", authChainMiddleware(appointmentHandler.ListAppointments))
		mux.HandleFunc("PUT /api/appointments/{id}", authChainMiddleware(rbacMiddleware.RequireOwnerOrStaff()(appointmentHandler.UpdateAppointment)))
		mux.HandleFunc("DELETE /api/appointments/{id}", authChainMiddleware(rbacMiddleware.RequireOwnerOrStaff()(appointmentHandler.CancelAppointment)))
	}

	// Metrics (メトリクス) endpoint
	metricsHandler := handler.NewMetricsHandler(metricsCollector)
	mux.HandleFunc("GET /api/metrics", chainMiddleware(metricsHandler.GetMetrics))

	// Auth (認証) endpoints
	if authEndpointHandler != nil {
		mux.HandleFunc("POST /api/auth/verify", chainMiddleware(authEndpointHandler.VerifyToken))
		log.Println("Auth endpoints registered")
	}

	// Health Check (監視ミドルウェアは適用しない)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		// データベース接続状態も確認
		status := "OK"
		if db != nil {
			if err := db.Ping(); err != nil {
				status = "WARNING: PostgreSQL connection failed"
			}
		} else {
			status = "WARNING: PostgreSQL not connected"
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(status))
	})

	// 5. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// CORS対応のHTTPハンドラーを作成
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORSヘッダーを設定
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// OPTIONSリクエスト（プリフライト）の処理
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// ルーターにリクエストを渡す
		mux.ServeHTTP(w, r)
	})

	log.Printf("TailorCloud Backend running on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
