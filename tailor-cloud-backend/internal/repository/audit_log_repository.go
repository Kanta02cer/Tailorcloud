package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"tailor-cloud/backend/internal/config/domain"
)

// AuditLogRepository 監査ログリポジトリインターフェース
type AuditLogRepository interface {
	Create(ctx context.Context, log *domain.AuditLog) error
	GetByResourceID(ctx context.Context, resourceType, resourceID string) ([]*domain.AuditLog, error)
	GetByTenantID(ctx context.Context, tenantID string, limit, offset int) ([]*domain.AuditLog, error)
	GetLogsByDateRange(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]*domain.AuditLog, error)
	MarkAsArchived(ctx context.Context, tenantID string, startDate, endDate time.Time, archiveLocation string) error
	UpdateLogHash(ctx context.Context, logID string, hash string) error
}

// ComplianceDocumentViewLogRepository 契約書閲覧ログリポジトリインターフェース
type ComplianceDocumentViewLogRepository interface {
	Create(ctx context.Context, log *domain.ComplianceDocumentViewLog) error
	GetByOrderID(ctx context.Context, orderID string) ([]*domain.ComplianceDocumentViewLog, error)
}

// PostgreSQLAuditLogRepository PostgreSQLを使った監査ログリポジトリ実装
type PostgreSQLAuditLogRepository struct {
	db *sql.DB
}

// NewPostgreSQLAuditLogRepository PostgreSQLAuditLogRepositoryのコンストラクタ
func NewPostgreSQLAuditLogRepository(db *sql.DB) AuditLogRepository {
	return &PostgreSQLAuditLogRepository{
		db: db,
	}
}

// Create 監査ログを作成
func (r *PostgreSQLAuditLogRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	query := `
		INSERT INTO audit_logs (
			id, tenant_id, user_id, action, resource_type, resource_id,
			old_value, new_value, changed_fields, ip_address, user_agent, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	
	// ChangedFieldsをJSON配列に変換
	changedFieldsJSON, err := json.Marshal(log.ChangedFields)
	if err != nil {
		return fmt.Errorf("failed to marshal changed_fields: %w", err)
	}
	
	_, err = r.db.ExecContext(ctx, query,
		log.ID,
		log.TenantID,
		log.UserID,
		string(log.Action),
		log.ResourceType,
		log.ResourceID,
		log.OldValue,
		log.NewValue,
		string(changedFieldsJSON),
		log.IPAddress,
		log.UserAgent,
		log.CreatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}
	
	return nil
}

// GetByResourceID リソースIDで監査ログを取得
func (r *PostgreSQLAuditLogRepository) GetByResourceID(ctx context.Context, resourceType, resourceID string) ([]*domain.AuditLog, error) {
	query := `
		SELECT 
			id, tenant_id, user_id, action, resource_type, resource_id,
			old_value, new_value, changed_fields, ip_address, user_agent, created_at
		FROM audit_logs
		WHERE resource_type = $1 AND resource_id = $2
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, resourceType, resourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()
	
	logs := make([]*domain.AuditLog, 0)
	
	for rows.Next() {
		var log domain.AuditLog
		var actionStr, changedFieldsJSON string
		
		err := rows.Scan(
			&log.ID,
			&log.TenantID,
			&log.UserID,
			&actionStr,
			&log.ResourceType,
			&log.ResourceID,
			&log.OldValue,
			&log.NewValue,
			&changedFieldsJSON,
			&log.IPAddress,
			&log.UserAgent,
			&log.CreatedAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		
		log.Action = domain.AuditAction(actionStr)
		
		// ChangedFieldsをJSON配列から復元
		if changedFieldsJSON != "" {
			err = json.Unmarshal([]byte(changedFieldsJSON), &log.ChangedFields)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal changed_fields: %w", err)
			}
		}
		
		logs = append(logs, &log)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating audit logs: %w", err)
	}
	
	return logs, nil
}

// GetByTenantID テナントIDで監査ログを取得（ページネーション対応）
func (r *PostgreSQLAuditLogRepository) GetByTenantID(ctx context.Context, tenantID string, limit, offset int) ([]*domain.AuditLog, error) {
	query := `
		SELECT 
			id, tenant_id, user_id, action, resource_type, resource_id,
			old_value, new_value, changed_fields, ip_address, user_agent, created_at
		FROM audit_logs
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.QueryContext(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()
	
	logs := make([]*domain.AuditLog, 0)
	
	for rows.Next() {
		var log domain.AuditLog
		var actionStr, changedFieldsJSON string
		
		err := rows.Scan(
			&log.ID,
			&log.TenantID,
			&log.UserID,
			&actionStr,
			&log.ResourceType,
			&log.ResourceID,
			&log.OldValue,
			&log.NewValue,
			&changedFieldsJSON,
			&log.IPAddress,
			&log.UserAgent,
			&log.CreatedAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		
		log.Action = domain.AuditAction(actionStr)
		
		// ChangedFieldsをJSON配列から復元
		if changedFieldsJSON != "" {
			err = json.Unmarshal([]byte(changedFieldsJSON), &log.ChangedFields)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal changed_fields: %w", err)
			}
		}
		
		logs = append(logs, &log)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating audit logs: %w", err)
	}
	
	return logs, nil
}

// GetLogsByDateRange 日付範囲で監査ログを取得
func (r *PostgreSQLAuditLogRepository) GetLogsByDateRange(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]*domain.AuditLog, error) {
	query := `
		SELECT 
			id, tenant_id, user_id, action, resource_type, resource_id,
			old_value, new_value, changed_fields, ip_address, user_agent, created_at
		FROM audit_logs
		WHERE tenant_id = $1 AND created_at >= $2 AND created_at < $3
		ORDER BY created_at ASC
	`
	
	rows, err := r.db.QueryContext(ctx, query, tenantID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs by date range: %w", err)
	}
	defer rows.Close()
	
	logs := make([]*domain.AuditLog, 0)
	
	for rows.Next() {
		var log domain.AuditLog
		var actionStr, changedFieldsJSON string
		
		err := rows.Scan(
			&log.ID,
			&log.TenantID,
			&log.UserID,
			&actionStr,
			&log.ResourceType,
			&log.ResourceID,
			&log.OldValue,
			&log.NewValue,
			&changedFieldsJSON,
			&log.IPAddress,
			&log.UserAgent,
			&log.CreatedAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		
		log.Action = domain.AuditAction(actionStr)
		
		// ChangedFieldsをJSON配列から復元
		if changedFieldsJSON != "" {
			err = json.Unmarshal([]byte(changedFieldsJSON), &log.ChangedFields)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal changed_fields: %w", err)
			}
		}
		
		logs = append(logs, &log)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating audit logs: %w", err)
	}
	
	return logs, nil
}

// MarkAsArchived ログをアーカイブ済みとしてマーク
func (r *PostgreSQLAuditLogRepository) MarkAsArchived(ctx context.Context, tenantID string, startDate, endDate time.Time, archiveLocation string) error {
	query := `
		UPDATE audit_logs
		SET archived_at = NOW(), archive_location = $4
		WHERE tenant_id = $1 AND created_at >= $2 AND created_at < $3 AND archived_at IS NULL
	`
	
	result, err := r.db.ExecContext(ctx, query, tenantID, startDate, endDate, archiveLocation)
	if err != nil {
		return fmt.Errorf("failed to mark logs as archived: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("no logs found to archive")
	}
	
	return nil
}

// UpdateLogHash ログのハッシュ値を更新
func (r *PostgreSQLAuditLogRepository) UpdateLogHash(ctx context.Context, logID string, hash string) error {
	query := `
		UPDATE audit_logs
		SET log_hash = $2
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query, logID, hash)
	if err != nil {
		return fmt.Errorf("failed to update log hash: %w", err)
	}
	
	return nil
}

// PostgreSQLComplianceDocumentViewLogRepository PostgreSQLを使った契約書閲覧ログリポジトリ実装
type PostgreSQLComplianceDocumentViewLogRepository struct {
	db *sql.DB
}

// NewPostgreSQLComplianceDocumentViewLogRepository PostgreSQLComplianceDocumentViewLogRepositoryのコンストラクタ
func NewPostgreSQLComplianceDocumentViewLogRepository(db *sql.DB) ComplianceDocumentViewLogRepository {
	return &PostgreSQLComplianceDocumentViewLogRepository{
		db: db,
	}
}

// Create 契約書閲覧ログを作成
func (r *PostgreSQLComplianceDocumentViewLogRepository) Create(ctx context.Context, log *domain.ComplianceDocumentViewLog) error {
	query := `
		INSERT INTO compliance_document_view_logs (
			id, order_id, tenant_id, user_id, document_url, document_hash,
			viewed_at, ip_address, user_agent
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		log.ID,
		log.OrderID,
		log.TenantID,
		log.UserID,
		log.DocumentURL,
		log.DocumentHash,
		log.ViewedAt,
		log.IPAddress,
		log.UserAgent,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create compliance document view log: %w", err)
	}
	
	return nil
}

// GetByOrderID 注文IDで契約書閲覧ログを取得
func (r *PostgreSQLComplianceDocumentViewLogRepository) GetByOrderID(ctx context.Context, orderID string) ([]*domain.ComplianceDocumentViewLog, error) {
	query := `
		SELECT 
			id, order_id, tenant_id, user_id, document_url, document_hash,
			viewed_at, ip_address, user_agent
		FROM compliance_document_view_logs
		WHERE order_id = $1
		ORDER BY viewed_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query compliance document view logs: %w", err)
	}
	defer rows.Close()
	
	logs := make([]*domain.ComplianceDocumentViewLog, 0)
	
	for rows.Next() {
		var log domain.ComplianceDocumentViewLog
		
		err := rows.Scan(
			&log.ID,
			&log.OrderID,
			&log.TenantID,
			&log.UserID,
			&log.DocumentURL,
			&log.DocumentHash,
			&log.ViewedAt,
			&log.IPAddress,
			&log.UserAgent,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan compliance document view log: %w", err)
		}
		
		logs = append(logs, &log)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating compliance document view logs: %w", err)
	}
	
	return logs, nil
}

