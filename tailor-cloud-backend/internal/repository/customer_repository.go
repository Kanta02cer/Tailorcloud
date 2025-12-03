package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"tailor-cloud/backend/internal/config/domain"
)

// CustomerRepository 顧客リポジトリインターフェース
type CustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	GetByID(ctx context.Context, customerID string, tenantID string) (*domain.Customer, error)
	GetByTenantID(ctx context.Context, tenantID string) ([]*domain.Customer, error)
	Search(ctx context.Context, tenantID string, keyword string) ([]*domain.Customer, error)
	Update(ctx context.Context, customer *domain.Customer) error
	Delete(ctx context.Context, customerID string, tenantID string) error
}

// PostgreSQLCustomerRepository PostgreSQLを使った顧客リポジトリ実装
type PostgreSQLCustomerRepository struct {
	db *sql.DB
}

// NewPostgreSQLCustomerRepository PostgreSQLCustomerRepositoryのコンストラクタ
func NewPostgreSQLCustomerRepository(db *sql.DB) CustomerRepository {
	return &PostgreSQLCustomerRepository{
		db: db,
	}
}

// Create 顧客を作成
func (r *PostgreSQLCustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	if customer.ID == "" {
		customer.ID = uuid.New().String()
	}
	
	now := time.Now()
	if customer.CreatedAt.IsZero() {
		customer.CreatedAt = now
	}
	if customer.UpdatedAt.IsZero() {
		customer.UpdatedAt = now
	}
	
	query := `
		INSERT INTO customers (
			id, tenant_id, name, email, phone, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		customer.ID,
		customer.TenantID,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.CreatedAt,
		customer.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}
	
	return nil
}

// GetByID 顧客IDで取得（テナントIDもチェック）
func (r *PostgreSQLCustomerRepository) GetByID(ctx context.Context, customerID string, tenantID string) (*domain.Customer, error) {
	query := `
		SELECT 
			id, tenant_id, name, email, phone, created_at, updated_at
		FROM customers
		WHERE id = $1 AND tenant_id = $2
	`
	
	var customer domain.Customer
	var email, phone sql.NullString
	
	err := r.db.QueryRowContext(ctx, query, customerID, tenantID).Scan(
		&customer.ID,
		&customer.TenantID,
		&customer.Name,
		&email,
		&phone,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("customer not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	
	if email.Valid {
		customer.Email = email.String
	}
	if phone.Valid {
		customer.Phone = phone.String
	}
	
	return &customer, nil
}

// GetByTenantID テナントIDで顧客一覧を取得
func (r *PostgreSQLCustomerRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*domain.Customer, error) {
	query := `
		SELECT 
			id, tenant_id, name, email, phone, created_at, updated_at
		FROM customers
		WHERE tenant_id = $1
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}
	defer rows.Close()
	
	var customers []*domain.Customer
	for rows.Next() {
		var customer domain.Customer
		var email, phone sql.NullString
		
		err := rows.Scan(
			&customer.ID,
			&customer.TenantID,
			&customer.Name,
			&email,
			&phone,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}
		
		if email.Valid {
			customer.Email = email.String
		}
		if phone.Valid {
			customer.Phone = phone.String
		}
		
		customers = append(customers, &customer)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate customers: %w", err)
	}
	
	return customers, nil
}

// Search 顧客を検索（名前、メール、電話番号で検索）
func (r *PostgreSQLCustomerRepository) Search(ctx context.Context, tenantID string, keyword string) ([]*domain.Customer, error) {
	query := `
		SELECT 
			id, tenant_id, name, email, phone, created_at, updated_at
		FROM customers
		WHERE tenant_id = $1
		  AND (
			name ILIKE $2 OR
			email ILIKE $2 OR
			phone ILIKE $2
		  )
		ORDER BY created_at DESC
	`
	
	searchPattern := "%" + keyword + "%"
	rows, err := r.db.QueryContext(ctx, query, tenantID, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search customers: %w", err)
	}
	defer rows.Close()
	
	var customers []*domain.Customer
	for rows.Next() {
		var customer domain.Customer
		var email, phone sql.NullString
		
		err := rows.Scan(
			&customer.ID,
			&customer.TenantID,
			&customer.Name,
			&email,
			&phone,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}
		
		if email.Valid {
			customer.Email = email.String
		}
		if phone.Valid {
			customer.Phone = phone.String
		}
		
		customers = append(customers, &customer)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate customers: %w", err)
	}
	
	return customers, nil
}

// Update 顧客を更新
func (r *PostgreSQLCustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	customer.UpdatedAt = time.Now()
	
	query := `
		UPDATE customers
		SET name = $3, email = $4, phone = $5, updated_at = $6
		WHERE id = $1 AND tenant_id = $2
	`
	
	result, err := r.db.ExecContext(ctx, query,
		customer.ID,
		customer.TenantID,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("customer not found or tenant_id mismatch")
	}
	
	return nil
}

// Delete 顧客を削除
func (r *PostgreSQLCustomerRepository) Delete(ctx context.Context, customerID string, tenantID string) error {
	query := `
		DELETE FROM customers
		WHERE id = $1 AND tenant_id = $2
	`
	
	result, err := r.db.ExecContext(ctx, query, customerID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("customer not found or tenant_id mismatch")
	}
	
	return nil
}

