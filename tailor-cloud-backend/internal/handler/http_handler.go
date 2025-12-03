package handler

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/middleware"
	"tailor-cloud/backend/internal/service"
)

// OrderHandler 注文ハンドラー
type OrderHandler struct {
	orderService *service.OrderService
}

// NewOrderHandler OrderHandlerのコンストラクタ
func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrderRequest HTTPリクエストボディ
type CreateOrderRequest struct {
	TenantID     string `json:"tenant_id"`
	CustomerID   string `json:"customer_id"`
	FabricID     string `json:"fabric_id"`
	TotalAmount  int64  `json:"total_amount"`
	DeliveryDate string `json:"delivery_date"` // ISO 8601形式 (例: "2025-12-31T00:00:00Z")
	Details      struct {
		MeasurementData json.RawMessage `json:"measurement_data"`
		Adjustments     json.RawMessage `json:"adjustments"`
		Description     string          `json:"description"` // 給付の内容（コンプライアンス用）
	} `json:"details"`
	CreatedBy string `json:"created_by"`
}

// CreateOrder POST /api/orders - 注文を作成
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// DeliveryDateをパース
	deliveryDate, err := time.Parse(time.RFC3339, req.DeliveryDate)
	if err != nil {
		// ISO 8601形式でも試す
		deliveryDate, err = time.Parse("2006-01-02T15:04:05Z", req.DeliveryDate)
		if err != nil {
			http.Error(w, "Invalid delivery_date format (expected ISO 8601): "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		// 認証情報がない場合は、リクエストのCreatedByを使用（開発環境用フォールバック）
		if req.CreatedBy == "" {
			http.Error(w, "Authentication required or created_by must be provided", http.StatusUnauthorized)
			return
		}
		authUser = &middleware.AuthUser{ID: req.CreatedBy}
	}

	// テナントID: 認証ユーザーから取得、またはリクエストから（フォールバック）
	tenantID := authUser.TenantID
	if tenantID == "" {
		tenantID = req.TenantID
		if tenantID == "" {
			http.Error(w, "tenant_id is required", http.StatusBadRequest)
			return
		}
	}

	// サービス層のリクエストに変換
	serviceReq := &service.CreateOrderRequest{
		TenantID:     tenantID,
		CustomerID:   req.CustomerID,
		FabricID:     req.FabricID,
		TotalAmount:  req.TotalAmount,
		DeliveryDate: deliveryDate,
		Details: &domain.OrderDetails{
			MeasurementData: req.Details.MeasurementData,
			Adjustments:     req.Details.Adjustments,
			Description:     req.Details.Description,
		},
		CreatedBy: authUser.ID, // 認証済みユーザーIDを使用
		IPAddress: extractIPAddress(r),
		UserAgent: r.UserAgent(),
	}

	// サービス層で注文を作成
	order, err := h.orderService.CreateOrder(r.Context(), serviceReq)
	if err != nil {
		http.Error(w, "Failed to create order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// レスポンス
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// ConfirmOrderRequest 注文確定リクエスト
type ConfirmOrderRequest struct {
	OrderID       string `json:"order_id"`
	TenantID      string `json:"tenant_id"`
	PrincipalName string `json:"principal_name"` // 委託をする者の氏名
}

// ConfirmOrder POST /api/orders/{order_id}/confirm - 注文を確定
func (h *OrderHandler) ConfirmOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ConfirmOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Authentication required: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// テナントID: 認証ユーザーから取得、またはリクエストから（フォールバック）
	tenantID := authUser.TenantID
	if tenantID == "" {
		tenantID = req.TenantID
		if tenantID == "" {
			http.Error(w, "tenant_id is required", http.StatusBadRequest)
			return
		}
	}

	// サービス層で注文を確定
	confirmReq := &service.ConfirmOrderRequest{
		OrderID:       req.OrderID,
		TenantID:      tenantID,
		PrincipalName: req.PrincipalName,
		UserID:        authUser.ID,
		IPAddress:     extractIPAddress(r),
		UserAgent:     r.UserAgent(),
	}

	order, err := h.orderService.ConfirmOrder(r.Context(), confirmReq)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "unauthorized: tenant_id mismatch" {
			statusCode = http.StatusUnauthorized
		} else if err.Error() == "order status must be Draft to confirm" {
			statusCode = http.StatusBadRequest
		}
		http.Error(w, "Failed to confirm order: "+err.Error(), statusCode)
		return
	}

	// レスポンス
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

// GetOrder GET /api/orders/{order_id} - 注文を取得
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLからorder_idを取得（簡易実装）
	orderID := r.URL.Query().Get("order_id")
	tenantID := r.URL.Query().Get("tenant_id")

	if orderID == "" {
		http.Error(w, "order_id is required", http.StatusBadRequest)
		return
	}
	if tenantID == "" {
		http.Error(w, "tenant_id is required", http.StatusBadRequest)
		return
	}

	order, err := h.orderService.GetOrder(r.Context(), orderID, tenantID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "unauthorized: tenant_id mismatch" {
			statusCode = http.StatusUnauthorized
		}
		http.Error(w, "Failed to get order: "+err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

// ListOrders GET /api/orders?tenant_id={tenant_id} - 注文一覧を取得
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		http.Error(w, "tenant_id is required", http.StatusBadRequest)
		return
	}

	orders, err := h.orderService.ListOrders(r.Context(), tenantID)
	if err != nil {
		http.Error(w, "Failed to list orders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

// extractIPAddress HTTPリクエストからIPアドレスを抽出
func extractIPAddress(r *http.Request) string {
	// X-Forwarded-For ヘッダーを確認（ロードバランサー経由の場合）
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// X-Real-IP ヘッダーを確認
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// RemoteAddrから取得
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
