package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// TaxRate と TaxRoundingMethod は tax.go で定義

// User ユーザーモデル
// マルチテナントアーキテクチャにおけるユーザー情報
type User struct {
	ID        string    `json:"id" firestore:"id"`
	TenantID  string    `json:"tenant_id" firestore:"tenant_id"`
	Name      string    `json:"name" firestore:"name"`
	Email     string    `json:"email" firestore:"email"`
	Role      UserRole  `json:"role" firestore:"role"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`
}

// UserRole ユーザー権限ロール
type UserRole string

const (
	RoleOwner         UserRole = "Owner"          // 契約、決済、全データ閲覧権限
	RoleStaff         UserRole = "Staff"          // 接客、発注作成権限
	RoleFactoryManager UserRole = "Factory_Manager" // 受注承認、工程管理権限
	RoleWorker        UserRole = "Worker"         // 作業完了チェックのみ
)

// Tenant テナントモデル
// テーラー（発注者）または縫製工場（受注者）の情報
type Tenant struct {
	ID                      string             `json:"id" firestore:"id" db:"id"`
	Type                    TenantType         `json:"type" firestore:"type" db:"type"`
	LegalName               string             `json:"legal_name" firestore:"legal_name" db:"legal_name"`
	Address                 string             `json:"address" firestore:"address" db:"address"`
	InvoiceRegistrationNo   string             `json:"invoice_registration_no" firestore:"invoice_registration_no" db:"invoice_registration_no"` // インボイス登録番号（T番号）
	TaxRoundingMethod       TaxRoundingMethod  `json:"tax_rounding_method" firestore:"tax_rounding_method" db:"tax_rounding_method"`               // 端数処理方法
	CreatedAt               time.Time          `json:"created_at" firestore:"created_at" db:"created_at"`
	UpdatedAt               time.Time          `json:"updated_at" firestore:"updated_at" db:"updated_at"`
}

// TenantType テナントタイプ
type TenantType string

const (
	TenantTypeTailor  TenantType = "Tailor"  // テーラー（発注者）
	TenantTypeFactory TenantType = "Factory" // 縫製工場（受注者）
)

// Order 注文モデル
// コンプライアンスエンジンの中心となるデータモデル
type Order struct {
	ID                string             `json:"id" firestore:"id" db:"id"`
	TenantID          string             `json:"tenant_id" firestore:"tenant_id" db:"tenant_id"`
	CustomerID        string             `json:"customer_id" firestore:"customer_id" db:"customer_id"`
	FabricID          string             `json:"fabric_id" firestore:"fabric_id" db:"fabric_id"`
	Status            OrderStatus        `json:"status" firestore:"status" db:"status"`
	ComplianceDocURL  string             `json:"compliance_doc_url" firestore:"compliance_doc_url" db:"compliance_doc_url"`
	ComplianceDocHash string             `json:"compliance_doc_hash" firestore:"compliance_doc_hash" db:"compliance_doc_hash"`
	TotalAmount       int64              `json:"total_amount" firestore:"total_amount" db:"total_amount"` // 税抜金額（円）
	TaxAmount         int64              `json:"tax_amount" firestore:"tax_amount" db:"tax_amount"`       // 消費税額（円）
	TaxRate           TaxRate            `json:"tax_rate" firestore:"tax_rate" db:"tax_rate"`             // 消費税率（0.10 = 10%, 0.08 = 8%）
	TaxExcludedAmount *int64             `json:"tax_excluded_amount" firestore:"tax_excluded_amount" db:"tax_excluded_amount"` // 税抜金額（明示的な場合）
	InvoiceIssuedAt   *time.Time         `json:"invoice_issued_at" firestore:"invoice_issued_at" db:"invoice_issued_at"`       // 請求書発行日時
	PaymentDueDate    time.Time          `json:"payment_due_date" firestore:"payment_due_date" db:"payment_due_date"`          // 支払期日（下請法60日ルール準拠）
	DeliveryDate      time.Time          `json:"delivery_date" firestore:"delivery_date" db:"delivery_date"`                   // 納期
	Details           *OrderDetails      `json:"details" firestore:"details" db:"details"`
	CreatedAt         time.Time          `json:"created_at" firestore:"created_at" db:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" firestore:"updated_at" db:"updated_at"`
	CreatedBy         string             `json:"created_by" firestore:"created_by" db:"created_by"` // ユーザーID
}

// OrderStatus 注文ステータス
// ワークフロー設計書のステータス遷移図に基づく
type OrderStatus string

const (
	OrderStatusDraft          OrderStatus = "Draft"           // 接客中(保存)
	OrderStatusConfirmed      OrderStatus = "Confirmed"       // 発注確定(法的拘束力発生)
	OrderStatusMaterialSecured OrderStatus = "Material_Secured" // 生地確保済み
	OrderStatusCutting        OrderStatus = "Cutting"         // 裁断開始
	OrderStatusSewing         OrderStatus = "Sewing"          // 縫製中
	OrderStatusInspection     OrderStatus = "Inspection"      // 検品中
	OrderStatusShipped        OrderStatus = "Shipped"         // 発送済み
	OrderStatusDelivered      OrderStatus = "Delivered"       // 納品完了
	OrderStatusPaid           OrderStatus = "Paid"            // 支払い完了
	OrderStatusCancelled      OrderStatus = "Cancelled"       // キャンセル
)

// OrderDetails 注文詳細
// 採寸データや調整情報を含む
type OrderDetails struct {
	MeasurementData json.RawMessage `json:"measurement_data" firestore:"measurement_data" db:"measurement_data"` // 採寸データ（JSON形式）
	Adjustments     json.RawMessage `json:"adjustments" firestore:"adjustments" db:"adjustments"`                 // 補正情報（JSON形式）
	Description     string          `json:"description" firestore:"description" db:"description"`                 // 給付の内容（コンプライアンス用）
}

// Customer 顧客モデル
type Customer struct {
	ID        string    `json:"id" firestore:"id"`
	TenantID  string    `json:"tenant_id" firestore:"tenant_id"`
	Name      string    `json:"name" firestore:"name"`
	Email     string    `json:"email" firestore:"email"`
	Phone     string    `json:"phone" firestore:"phone"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`
}

// Fabric 生地モデル
// 在庫連携APIとの連携に使用
type Fabric struct {
	ID           string      `json:"id" firestore:"id" db:"id"`
	SupplierID   string      `json:"supplier_id" firestore:"supplier_id" db:"supplier_id"`
	Name         string      `json:"name" firestore:"name" db:"name"`
	StockAmount  float64     `json:"stock_amount" firestore:"stock_amount" db:"stock_amount"` // 在庫数量（メートル）
	Price        int64       `json:"price" firestore:"price" db:"price"`                      // 単価（円/メートル）
	StockStatus  StockStatus `json:"stock_status" firestore:"stock_status" db:"stock_status"` // 在庫ステータス（計算フィールド）
	ImageURL     string      `json:"image_url" firestore:"image_url" db:"image_url"`          // 生地画像URL（UI表示用）
	MinimumOrder float64     `json:"minimum_order" firestore:"minimum_order" db:"minimum_order"` // 最小発注数量（デフォルト3.2m = スーツ1着分）
	CreatedAt    time.Time   `json:"created_at" firestore:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" firestore:"updated_at" db:"updated_at"`
}

// StockStatus 在庫ステータス
// 在庫連携ロジックに基づく
type StockStatus string

const (
	StockStatusAvailable StockStatus = "Available" // ◎ (在庫あり) Stock > 3.2m
	StockStatusLimited   StockStatus = "Limited"   // △ (要確認) 3.2m > Stock > 0
	StockStatusSoldOut   StockStatus = "SoldOut"   // × (Sold Out) Stock = 0
)

// CalculateStockStatus 在庫数量からステータスを計算
func (f *Fabric) CalculateStockStatus() {
	if f.StockAmount > 3.2 {
		f.StockStatus = StockStatusAvailable
	} else if f.StockAmount > 0 {
		f.StockStatus = StockStatusLimited
	} else {
		f.StockStatus = StockStatusSoldOut
	}
}

// Transaction 取引モデル
// 決済トランザクション（Phase 3で使用）
type Transaction struct {
	ID            string           `json:"id" db:"id"`
	OrderID       string           `json:"order_id" db:"order_id"`
	Status        TransactionStatus `json:"status" db:"status"`
	PaymentMethod string           `json:"payment_method" db:"payment_method"`
	Amount        int64            `json:"amount" db:"amount"`
	CreatedAt     time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at" db:"updated_at"`
}

// TransactionStatus 取引ステータス
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "Pending"
	TransactionStatusCompleted TransactionStatus = "Completed"
	TransactionStatusFailed    TransactionStatus = "Failed"
)

// ComplianceDocument は compliance.go で定義されています（履歴管理対応版）

// NewOrder 新しい注文を作成
func NewOrder(tenantID, customerID, fabricID, createdBy string, totalAmount int64, deliveryDate time.Time) *Order {
	now := time.Now()
	
	// 下請法60日ルール: 納期から60日後に支払期日を設定
	paymentDueDate := deliveryDate.AddDate(0, 0, 60)
	
	order := &Order{
		ID:             uuid.New().String(),
		TenantID:       tenantID,
		CustomerID:     customerID,
		FabricID:       fabricID,
		Status:         OrderStatusDraft,
		TotalAmount:    totalAmount,
		PaymentDueDate: paymentDueDate,
		DeliveryDate:   deliveryDate,
		CreatedAt:      now,
		UpdatedAt:      now,
		CreatedBy:      createdBy,
	}
	
	return order
}

