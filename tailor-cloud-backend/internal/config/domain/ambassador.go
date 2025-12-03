package domain

import (
	"time"

	"github.com/google/uuid"
)

// Ambassador アンバサダーモデル
// Phase 1: 学生アンバサダーが「誰が売ったか」を記録し、成果報酬を自動計算する基盤
type Ambassador struct {
	ID          string    `json:"id" db:"id"`
	TenantID    string    `json:"tenant_id" db:"tenant_id"`
	UserID      string    `json:"user_id" db:"user_id"` // Firebase AuthのUserIDとリンク
	Name        string    `json:"name" db:"name"`
	Email       string    `json:"email" db:"email"`
	Phone       string    `json:"phone" db:"phone"`
	Status      AmbassadorStatus `json:"status" db:"status"`
	CommissionRate float64 `json:"commission_rate" db:"commission_rate"` // 成果報酬率（例: 0.10 = 10%）
	TotalSales  int64     `json:"total_sales" db:"total_sales"` // 累計売上（円）
	TotalCommission int64 `json:"total_commission" db:"total_commission"` // 累計報酬（円）
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// AmbassadorStatus アンバサダーステータス
type AmbassadorStatus string

const (
	AmbassadorStatusActive   AmbassadorStatus = "Active"   // アクティブ
	AmbassadorStatusInactive AmbassadorStatus = "Inactive" // 非アクティブ
	AmbassadorStatusSuspended AmbassadorStatus = "Suspended" // 一時停止
)

// Commission 成果報酬モデル
// 注文ごとの成果報酬を記録
type Commission struct {
	ID            string    `json:"id" db:"id"`
	OrderID       string    `json:"order_id" db:"order_id"`
	AmbassadorID  string    `json:"ambassador_id" db:"ambassador_id"`
	TenantID      string    `json:"tenant_id" db:"tenant_id"`
	OrderAmount   int64     `json:"order_amount" db:"order_amount"` // 注文金額
	CommissionRate float64  `json:"commission_rate" db:"commission_rate"` // 報酬率
	CommissionAmount int64  `json:"commission_amount" db:"commission_amount"` // 報酬額（円）
	Status        CommissionStatus `json:"status" db:"status"`
	PaidAt        *time.Time `json:"paid_at" db:"paid_at"` // 支払日（nil = 未払い）
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// CommissionStatus 成果報酬ステータス
type CommissionStatus string

const (
	CommissionStatusPending   CommissionStatus = "Pending"   // 未確定（注文が確定していない）
	CommissionStatusApproved  CommissionStatus = "Approved"  // 確定（支払い待ち）
	CommissionStatusPaid      CommissionStatus = "Paid"      // 支払済み
	CommissionStatusCancelled CommissionStatus = "Cancelled" // キャンセル（注文がキャンセルされた場合）
)

// CalculateCommission 成果報酬を計算
func CalculateCommission(orderAmount int64, commissionRate float64) int64 {
	return int64(float64(orderAmount) * commissionRate)
}

// NewAmbassador 新しいアンバサダーを作成
func NewAmbassador(tenantID, userID, name, email string) *Ambassador {
	now := time.Now()
	
	// デフォルトの成果報酬率: 10%（Phase 1では固定値）
	defaultCommissionRate := 0.10
	
	return &Ambassador{
		ID:            uuid.New().String(),
		TenantID:      tenantID,
		UserID:        userID,
		Name:          name,
		Email:         email,
		Status:        AmbassadorStatusActive,
		CommissionRate: defaultCommissionRate,
		TotalSales:    0,
		TotalCommission: 0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// NewCommission 新しい成果報酬を作成
func NewCommission(orderID, ambassadorID, tenantID string, orderAmount int64, commissionRate float64) *Commission {
	commissionAmount := CalculateCommission(orderAmount, commissionRate)
	
	return &Commission{
		ID:              uuid.New().String(),
		OrderID:         orderID,
		AmbassadorID:    ambassadorID,
		TenantID:        tenantID,
		OrderAmount:     orderAmount,
		CommissionRate:  commissionRate,
		CommissionAmount: commissionAmount,
		Status:          CommissionStatusPending, // 注文確定までPending
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

