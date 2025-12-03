package domain

import (
	"time"

	"github.com/google/uuid"
)

// AuditLog 監査ログモデル
// 仕様書要件: 「誰が」「いつ」「どの数値を」変更したかの完全なログ保存（法的証拠能力のため）
type AuditLog struct {
	ID              string     `json:"id" db:"id"`
	TenantID        string     `json:"tenant_id" db:"tenant_id"`
	UserID          string     `json:"user_id" db:"user_id"`
	Action          AuditAction `json:"action" db:"action"`
	ResourceType    string     `json:"resource_type" db:"resource_type"` // "order", "customer", "fabric" など
	ResourceID      string     `json:"resource_id" db:"resource_id"`
	OldValue        string     `json:"old_value" db:"old_value"` // JSON形式で変更前の値を保存
	NewValue        string     `json:"new_value" db:"new_value"` // JSON形式で変更後の値を保存
	ChangedFields   []string   `json:"changed_fields" db:"changed_fields"` // 変更されたフィールド名のリスト
	IPAddress       string     `json:"ip_address" db:"ip_address"`
	UserAgent       string     `json:"user_agent" db:"user_agent"`
	DeviceID        string     `json:"device_id" db:"device_id"`               // デバイスID（エンタープライズ要件）
	ChangeSummary   string     `json:"change_summary" db:"change_summary"`     // 変更サマリー（人間が読みやすい形式）
	LogHash         string     `json:"log_hash" db:"log_hash"`                 // ログのハッシュ値（改ざん検知用）
	ArchivedAt      *time.Time `json:"archived_at" db:"archived_at"`           // アーカイブ日時
	ArchiveLocation string     `json:"archive_location" db:"archive_location"` // アーカイブ先（Cloud Storageパス）
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}

// AuditAction 監査アクションタイプ
type AuditAction string

const (
	AuditActionCreate   AuditAction = "CREATE"   // 作成
	AuditActionUpdate   AuditAction = "UPDATE"   // 更新
	AuditActionDelete   AuditAction = "DELETE"   // 削除
	AuditActionView     AuditAction = "VIEW"     // 閲覧（契約書PDF閲覧など）
	AuditActionConfirm  AuditAction = "CONFIRM"  // 確定（注文確定など）
	AuditActionStatusChange AuditAction = "STATUS_CHANGE" // ステータス変更
)

// ComplianceDocumentViewLog 契約書閲覧ログ
// 仕様書要件: 「いつ契約書を閲覧したか」の完全なログ保存
type ComplianceDocumentViewLog struct {
	ID              string    `json:"id" db:"id"`
	OrderID         string    `json:"order_id" db:"order_id"`
	TenantID        string    `json:"tenant_id" db:"tenant_id"`
	UserID          string    `json:"user_id" db:"user_id"`
	DocumentURL     string    `json:"document_url" db:"document_url"`
	DocumentHash    string    `json:"document_hash" db:"document_hash"`
	ViewedAt        time.Time `json:"viewed_at" db:"viewed_at"`
	IPAddress       string    `json:"ip_address" db:"ip_address"`
	UserAgent       string    `json:"user_agent" db:"user_agent"`
}

// NewAuditLog 新しい監査ログを作成
func NewAuditLog(tenantID, userID string, action AuditAction, resourceType, resourceID string) *AuditLog {
	return &AuditLog{
		ID:           generateUUID(),
		TenantID:     tenantID,
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		CreatedAt:    time.Now(),
	}
}

// generateUUID UUID生成ヘルパー
func generateUUID() string {
	return uuid.New().String()
}

