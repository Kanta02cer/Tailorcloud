package domain

import (
	"time"

	"github.com/google/uuid"
)

// Appointment 予約モデル
// フィッティング予約を管理
type Appointment struct {
	ID                      string           `json:"id" db:"id"`
	UserID                  string           `json:"user_id" db:"user_id"`
	TenantID                string           `json:"tenant_id" db:"tenant_id"`
	FitterID                string           `json:"fitter_id" db:"fitter_id"` // フィッターID（スタッフ）
	AppointmentDateTime     time.Time        `json:"appointment_datetime" db:"appointment_datetime"`
	DurationMinutes         int              `json:"duration_minutes" db:"duration_minutes"` // 予約時間（分）
	Status                  AppointmentStatus `json:"status" db:"status"`
	DepositAmount           *int64           `json:"deposit_amount" db:"deposit_amount"`                    // デポジット金額（円）
	DepositPaymentIntentID  string           `json:"deposit_payment_intent_id" db:"deposit_payment_intent_id"` // Stripe Payment Intent ID
	DepositStatus           DepositStatus    `json:"deposit_status" db:"deposit_status"`
	Notes                   string           `json:"notes" db:"notes"`                   // メモ
	CancelledAt             *time.Time       `json:"cancelled_at" db:"cancelled_at"`     // キャンセル日時
	CancelledReason         string           `json:"cancelled_reason" db:"cancelled_reason"` // キャンセル理由
	CreatedAt               time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time        `json:"updated_at" db:"updated_at"`
}

// AppointmentStatus 予約ステータス
type AppointmentStatus string

const (
	AppointmentStatusPending   AppointmentStatus = "Pending"   // 保留中
	AppointmentStatusConfirmed AppointmentStatus = "Confirmed" // 確定
	AppointmentStatusCancelled AppointmentStatus = "Cancelled" // キャンセル
	AppointmentStatusCompleted AppointmentStatus = "Completed" // 完了
	AppointmentStatusNoShow    AppointmentStatus = "NoShow"    // 無断キャンセル
)

// DepositStatus デポジットステータス
type DepositStatus string

const (
	DepositStatusPending  DepositStatus = "pending"  // 保留中
	DepositStatusSucceeded DepositStatus = "succeeded" // 成功
	DepositStatusFailed   DepositStatus = "failed"   // 失敗
	DepositStatusRefunded DepositStatus = "refunded" // 返金済み
)

// IsValid 予約ステータスが有効かチェック
func (s AppointmentStatus) IsValid() bool {
	switch s {
	case AppointmentStatusPending, AppointmentStatusConfirmed, AppointmentStatusCancelled, AppointmentStatusCompleted, AppointmentStatusNoShow:
		return true
	default:
		return false
	}
}

// IsValid デポジットステータスが有効かチェック
func (s DepositStatus) IsValid() bool {
	switch s {
	case DepositStatusPending, DepositStatusSucceeded, DepositStatusFailed, DepositStatusRefunded:
		return true
	default:
		return false
	}
}

// NewAppointment 新しい予約を作成
func NewAppointment(userID, tenantID, fitterID string, appointmentDateTime time.Time, durationMinutes int) *Appointment {
	now := time.Now()
	return &Appointment{
		ID:                  uuid.New().String(),
		UserID:              userID,
		TenantID:            tenantID,
		FitterID:            fitterID,
		AppointmentDateTime: appointmentDateTime,
		DurationMinutes:     durationMinutes,
		Status:              AppointmentStatusPending,
		DepositStatus:       DepositStatusPending,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}

// IsCancellable キャンセル可能かチェック（3日前までキャンセル可能）
func (a *Appointment) IsCancellable() bool {
	if a.Status != AppointmentStatusConfirmed && a.Status != AppointmentStatusPending {
		return false
	}
	
	// 3日前かどうかチェック
	threeDaysBefore := a.AppointmentDateTime.AddDate(0, 0, -3)
	return time.Now().Before(threeDaysBefore)
}

// CanRefundDeposit デポジットを返金可能かチェック
func (a *Appointment) CanRefundDeposit() bool {
	// 3日前までなら返金可能
	return a.IsCancellable() && a.DepositStatus == DepositStatusSucceeded
}

