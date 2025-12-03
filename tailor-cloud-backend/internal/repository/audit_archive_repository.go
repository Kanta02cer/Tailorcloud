package repository

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"tailor-cloud/backend/internal/config/domain"
)

// AuditLogArchiveRepository 監査ログアーカイブリポジトリインターフェース
type AuditLogArchiveRepository interface {
	Create(ctx context.Context, archive *domain.AuditLogArchive) error
	GetByTenantID(ctx context.Context, tenantID string) ([]*domain.AuditLogArchive, error)
	GetByID(ctx context.Context, archiveID string) (*domain.AuditLogArchive, error)
}

// PostgreSQLAuditLogArchiveRepository PostgreSQL実装
type PostgreSQLAuditLogArchiveRepository struct {
	db *sql.DB
}

// NewPostgreSQLAuditLogArchiveRepository PostgreSQLAuditLogArchiveRepositoryのコンストラクタ
func NewPostgreSQLAuditLogArchiveRepository(db *sql.DB) AuditLogArchiveRepository {
	return &PostgreSQLAuditLogArchiveRepository{
		db: db,
	}
}

// Create アーカイブメタデータを作成
func (r *PostgreSQLAuditLogArchiveRepository) Create(ctx context.Context, archive *domain.AuditLogArchive) error {
	query := `
		INSERT INTO audit_log_archives (
			id, tenant_id, archive_period_start, archive_period_end,
			log_count, archive_location, archive_hash, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		archive.ID,
		archive.TenantID,
		archive.ArchivePeriodStart,
		archive.ArchivePeriodEnd,
		archive.LogCount,
		archive.ArchiveLocation,
		archive.ArchiveHash,
		archive.CreatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create audit log archive: %w", err)
	}
	
	return nil
}

// GetByTenantID テナントIDでアーカイブメタデータを取得
func (r *PostgreSQLAuditLogArchiveRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*domain.AuditLogArchive, error) {
	query := `
		SELECT id, tenant_id, archive_period_start, archive_period_end,
		       log_count, archive_location, archive_hash, created_at
		FROM audit_log_archives
		WHERE tenant_id = $1
		ORDER BY archive_period_start DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit log archives: %w", err)
	}
	defer rows.Close()
	
	archives := make([]*domain.AuditLogArchive, 0)
	
	for rows.Next() {
		var archive domain.AuditLogArchive
		
		err := rows.Scan(
			&archive.ID,
			&archive.TenantID,
			&archive.ArchivePeriodStart,
			&archive.ArchivePeriodEnd,
			&archive.LogCount,
			&archive.ArchiveLocation,
			&archive.ArchiveHash,
			&archive.CreatedAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log archive: %w", err)
		}
		
		archives = append(archives, &archive)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating audit log archives: %w", err)
	}
	
	return archives, nil
}

// GetByID アーカイブIDでアーカイブメタデータを取得
func (r *PostgreSQLAuditLogArchiveRepository) GetByID(ctx context.Context, archiveID string) (*domain.AuditLogArchive, error) {
	query := `
		SELECT id, tenant_id, archive_period_start, archive_period_end,
		       log_count, archive_location, archive_hash, created_at
		FROM audit_log_archives
		WHERE id = $1
	`
	
	var archive domain.AuditLogArchive
	
	err := r.db.QueryRowContext(ctx, query, archiveID).Scan(
		&archive.ID,
		&archive.TenantID,
		&archive.ArchivePeriodStart,
		&archive.ArchivePeriodEnd,
		&archive.LogCount,
		&archive.ArchiveLocation,
		&archive.ArchiveHash,
		&archive.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("audit log archive not found: %s", archiveID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get audit log archive: %w", err)
	}
	
	return &archive, nil
}

