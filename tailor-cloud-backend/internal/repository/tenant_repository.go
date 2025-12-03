package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"tailor-cloud/backend/internal/config/domain"
)

// TenantRepository テナントリポジトリインターフェース
type TenantRepository interface {
	GetByID(ctx context.Context, tenantID string) (*domain.Tenant, error)
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

