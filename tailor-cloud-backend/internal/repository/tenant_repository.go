package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tailor-cloud/backend/internal/config/domain"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// TenantRepository テナントリポジトリインターフェース
type TenantRepository interface {
	GetByID(ctx context.Context, tenantID string) (*domain.Tenant, error)
	GetByDomain(ctx context.Context, domain string) (*domain.Tenant, error) // 将来の拡張用
	Create(ctx context.Context, tenant *domain.Tenant) error
	Update(ctx context.Context, tenant *domain.Tenant) error
}

// PostgreSQLTenantRepository PostgreSQLを使ったテナントリポジトリ実装
type PostgreSQLTenantRepository struct {
	db *sql.DB
}

// NewPostgreSQLTenantRepository PostgreSQLTenantRepositoryのコンストラクタ
func NewPostgreSQLTenantRepository(db *sql.DB) TenantRepository {
	return &PostgreSQLTenantRepository{
		db: db,
	}
}

// GetByID テナントIDで取得
func (r *PostgreSQLTenantRepository) GetByID(ctx context.Context, tenantID string) (*domain.Tenant, error) {
	query := `
		SELECT 
			id, type, legal_name, address,
			invoice_registration_no, tax_rounding_method,
			created_at, updated_at
		FROM tenants
		WHERE id = $1
	`

	var tenant domain.Tenant
	var legalName, address, invoiceRegNo, taxRoundingMethod sql.NullString
	var typeStr string

	err := r.db.QueryRowContext(ctx, query, tenantID).Scan(
		&tenant.ID,
		&typeStr,
		&legalName,
		&address,
		&invoiceRegNo,
		&taxRoundingMethod,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tenant not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	tenant.Type = domain.TenantType(typeStr)
	if legalName.Valid {
		tenant.LegalName = legalName.String
	}
	if address.Valid {
		tenant.Address = address.String
	}
	if invoiceRegNo.Valid {
		tenant.InvoiceRegistrationNo = invoiceRegNo.String
	}
	if taxRoundingMethod.Valid {
		tenant.TaxRoundingMethod = domain.TaxRoundingMethod(taxRoundingMethod.String)
	} else {
		tenant.TaxRoundingMethod = domain.TaxRoundingMethodHalfUp // デフォルト
	}

	return &tenant, nil
}

// GetByDomain ドメインでテナントを取得（将来の拡張用）
func (r *PostgreSQLTenantRepository) GetByDomain(ctx context.Context, domain string) (*domain.Tenant, error) {
	// TODO: テナントテーブルにdomainカラムを追加して実装
	return nil, fmt.Errorf("not implemented")
}

// Create テナントを作成
func (r *PostgreSQLTenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		INSERT INTO tenants (id, name, type, legal_name, address, invoice_registration_no, tax_rounding_method, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	tenantID := tenant.ID
	if tenantID == "" {
		tenantID = uuid.New().String()
		tenant.ID = tenantID
	}

	now := time.Now()
	if tenant.CreatedAt.IsZero() {
		tenant.CreatedAt = now
	}
	if tenant.UpdatedAt.IsZero() {
		tenant.UpdatedAt = now
	}

	_, err := r.db.ExecContext(ctx, query,
		tenantID,
		tenant.LegalName, // nameカラムにlegal_nameを使用
		string(tenant.Type),
		tenant.LegalName,
		tenant.Address,
		tenant.InvoiceRegistrationNo,
		string(tenant.TaxRoundingMethod),
		tenant.CreatedAt,
		tenant.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}

	return nil
}

// Update テナントを更新
func (r *PostgreSQLTenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	tenant.UpdatedAt = time.Now()

	query := `
		UPDATE tenants
		SET legal_name = $2,
		    address = $3,
		    invoice_registration_no = $4,
		    tax_rounding_method = $5,
		    updated_at = $6
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		tenant.ID,
		tenant.LegalName,
		tenant.Address,
		tenant.InvoiceRegistrationNo,
		tenant.TaxRoundingMethod,
		tenant.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update tenant: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tenant not found")
	}

	return nil
}
