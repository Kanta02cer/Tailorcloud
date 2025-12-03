package domain

import (
	"time"
)

// AuditLogArchive 監査ログアーカイブメタデータ
// WORMストレージに移行したログの情報
type AuditLogArchive struct {
	ID                string    `json:"id" db:"id"`
	TenantID          string    `json:"tenant_id" db:"tenant_id"`
	ArchivePeriodStart time.Time `json:"archive_period_start" db:"archive_period_start"`
	ArchivePeriodEnd   time.Time `json:"archive_period_end" db:"archive_period_end"`
	LogCount           int64     `json:"log_count" db:"log_count"`
	ArchiveLocation    string    `json:"archive_location" db:"archive_location"`
	ArchiveHash        string    `json:"archive_hash" db:"archive_hash"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

