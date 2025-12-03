package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"tailor-cloud/backend/internal/config/domain"
)

// AmbassadorRepository アンバサダーリポジトリインターフェース
type AmbassadorRepository interface {
	Create(ctx context.Context, ambassador *domain.Ambassador) error
	GetByID(ctx context.Context, ambassadorID string) (*domain.Ambassador, error)
	GetByUserID(ctx context.Context, userID string) (*domain.Ambassador, error)
	GetByTenantID(ctx context.Context, tenantID string) ([]*domain.Ambassador, error)
	Update(ctx context.Context, ambassador *domain.Ambassador) error
	UpdateSalesStats(ctx context.Context, ambassadorID string, orderAmount, commissionAmount int64) error
}

// CommissionRepository 成果報酬リポジトリインターフェース
type CommissionRepository interface {
	Create(ctx context.Context, commission *domain.Commission) error
	GetByID(ctx context.Context, commissionID string) (*domain.Commission, error)
	GetByOrderID(ctx context.Context, orderID string) (*domain.Commission, error)
	GetByAmbassadorID(ctx context.Context, ambassadorID string, limit, offset int) ([]*domain.Commission, error)
	GetByTenantID(ctx context.Context, tenantID string, limit, offset int) ([]*domain.Commission, error)
	UpdateStatus(ctx context.Context, commissionID string, status domain.CommissionStatus) error
}

// PostgreSQLAmbassadorRepository PostgreSQLを使ったアンバサダーリポジトリ実装
type PostgreSQLAmbassadorRepository struct {
	db *sql.DB
}

// NewPostgreSQLAmbassadorRepository PostgreSQLAmbassadorRepositoryのコンストラクタ
func NewPostgreSQLAmbassadorRepository(db *sql.DB) AmbassadorRepository {
	return &PostgreSQLAmbassadorRepository{
		db: db,
	}
}

// Create アンバサダーを作成
func (r *PostgreSQLAmbassadorRepository) Create(ctx context.Context, ambassador *domain.Ambassador) error {
	query := `
		INSERT INTO ambassadors (
			id, tenant_id, user_id, name, email, phone,
			status, commission_rate, total_sales, total_commission,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		ambassador.ID,
		ambassador.TenantID,
		ambassador.UserID,
		ambassador.Name,
		ambassador.Email,
		ambassador.Phone,
		string(ambassador.Status),
		ambassador.CommissionRate,
		ambassador.TotalSales,
		ambassador.TotalCommission,
		ambassador.CreatedAt,
		ambassador.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create ambassador: %w", err)
	}
	
	return nil
}

// GetByID アンバサダーIDで取得
func (r *PostgreSQLAmbassadorRepository) GetByID(ctx context.Context, ambassadorID string) (*domain.Ambassador, error) {
	query := `
		SELECT 
			id, tenant_id, user_id, name, email, phone,
			status, commission_rate, total_sales, total_commission,
			created_at, updated_at
		FROM ambassadors
		WHERE id = $1
	`
	
	var ambassador domain.Ambassador
	var statusStr string
	
	err := r.db.QueryRowContext(ctx, query, ambassadorID).Scan(
		&ambassador.ID,
		&ambassador.TenantID,
		&ambassador.UserID,
		&ambassador.Name,
		&ambassador.Email,
		&ambassador.Phone,
		&statusStr,
		&ambassador.CommissionRate,
		&ambassador.TotalSales,
		&ambassador.TotalCommission,
		&ambassador.CreatedAt,
		&ambassador.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ambassador not found: %s", ambassadorID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get ambassador: %w", err)
	}
	
	ambassador.Status = domain.AmbassadorStatus(statusStr)
	
	return &ambassador, nil
}

// GetByUserID ユーザーIDでアンバサダーを取得
func (r *PostgreSQLAmbassadorRepository) GetByUserID(ctx context.Context, userID string) (*domain.Ambassador, error) {
	query := `
		SELECT 
			id, tenant_id, user_id, name, email, phone,
			status, commission_rate, total_sales, total_commission,
			created_at, updated_at
		FROM ambassadors
		WHERE user_id = $1
	`
	
	var ambassador domain.Ambassador
	var statusStr string
	
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&ambassador.ID,
		&ambassador.TenantID,
		&ambassador.UserID,
		&ambassador.Name,
		&ambassador.Email,
		&ambassador.Phone,
		&statusStr,
		&ambassador.CommissionRate,
		&ambassador.TotalSales,
		&ambassador.TotalCommission,
		&ambassador.CreatedAt,
		&ambassador.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ambassador not found for user_id: %s", userID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get ambassador by user_id: %w", err)
	}
	
	ambassador.Status = domain.AmbassadorStatus(statusStr)
	
	return &ambassador, nil
}

// GetByTenantID テナントIDでアンバサダー一覧を取得
func (r *PostgreSQLAmbassadorRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*domain.Ambassador, error) {
	query := `
		SELECT 
			id, tenant_id, user_id, name, email, phone,
			status, commission_rate, total_sales, total_commission,
			created_at, updated_at
		FROM ambassadors
		WHERE tenant_id = $1
		ORDER BY total_sales DESC, created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query ambassadors: %w", err)
	}
	defer rows.Close()
	
	ambassadors := make([]*domain.Ambassador, 0)
	
	for rows.Next() {
		var ambassador domain.Ambassador
		var statusStr string
		
		err := rows.Scan(
			&ambassador.ID,
			&ambassador.TenantID,
			&ambassador.UserID,
			&ambassador.Name,
			&ambassador.Email,
			&ambassador.Phone,
			&statusStr,
			&ambassador.CommissionRate,
			&ambassador.TotalSales,
			&ambassador.TotalCommission,
			&ambassador.CreatedAt,
			&ambassador.UpdatedAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan ambassador: %w", err)
		}
		
		ambassador.Status = domain.AmbassadorStatus(statusStr)
		ambassadors = append(ambassadors, &ambassador)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ambassadors: %w", err)
	}
	
	return ambassadors, nil
}

// Update アンバサダーを更新
func (r *PostgreSQLAmbassadorRepository) Update(ctx context.Context, ambassador *domain.Ambassador) error {
	query := `
		UPDATE ambassadors SET
			name = $2,
			email = $3,
			phone = $4,
			status = $5,
			commission_rate = $6,
			updated_at = $7
		WHERE id = $1 AND tenant_id = $8
	`
	
	ambassador.UpdatedAt = time.Now()
	
	result, err := r.db.ExecContext(ctx, query,
		ambassador.ID,
		ambassador.Name,
		ambassador.Email,
		ambassador.Phone,
		string(ambassador.Status),
		ambassador.CommissionRate,
		ambassador.UpdatedAt,
		ambassador.TenantID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update ambassador: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("ambassador not found or tenant_id mismatch")
	}
	
	return nil
}

// UpdateSalesStats 売上統計を更新
func (r *PostgreSQLAmbassadorRepository) UpdateSalesStats(ctx context.Context, ambassadorID string, orderAmount, commissionAmount int64) error {
	query := `
		UPDATE ambassadors SET
			total_sales = total_sales + $2,
			total_commission = total_commission + $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query, ambassadorID, orderAmount, commissionAmount)
	if err != nil {
		return fmt.Errorf("failed to update sales stats: %w", err)
	}
	
	return nil
}

// PostgreSQLCommissionRepository PostgreSQLを使った成果報酬リポジトリ実装
type PostgreSQLCommissionRepository struct {
	db *sql.DB
}

// NewPostgreSQLCommissionRepository PostgreSQLCommissionRepositoryのコンストラクタ
func NewPostgreSQLCommissionRepository(db *sql.DB) CommissionRepository {
	return &PostgreSQLCommissionRepository{
		db: db,
	}
}

// Create 成果報酬を作成
func (r *PostgreSQLCommissionRepository) Create(ctx context.Context, commission *domain.Commission) error {
	query := `
		INSERT INTO commissions (
			id, order_id, ambassador_id, tenant_id,
			order_amount, commission_rate, commission_amount,
			status, paid_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	
	var paidAt *time.Time = commission.PaidAt
	
	_, err := r.db.ExecContext(ctx, query,
		commission.ID,
		commission.OrderID,
		commission.AmbassadorID,
		commission.TenantID,
		commission.OrderAmount,
		commission.CommissionRate,
		commission.CommissionAmount,
		string(commission.Status),
		paidAt,
		commission.CreatedAt,
		commission.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create commission: %w", err)
	}
	
	return nil
}

// GetByID 成果報酬IDで取得
func (r *PostgreSQLCommissionRepository) GetByID(ctx context.Context, commissionID string) (*domain.Commission, error) {
	query := `
		SELECT 
			id, order_id, ambassador_id, tenant_id,
			order_amount, commission_rate, commission_amount,
			status, paid_at, created_at, updated_at
		FROM commissions
		WHERE id = $1
	`
	
	var commission domain.Commission
	var statusStr string
	var paidAt sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, commissionID).Scan(
		&commission.ID,
		&commission.OrderID,
		&commission.AmbassadorID,
		&commission.TenantID,
		&commission.OrderAmount,
		&commission.CommissionRate,
		&commission.CommissionAmount,
		&statusStr,
		&paidAt,
		&commission.CreatedAt,
		&commission.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("commission not found: %s", commissionID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get commission: %w", err)
	}
	
	commission.Status = domain.CommissionStatus(statusStr)
	
	if paidAt.Valid {
		commission.PaidAt = &paidAt.Time
	}
	
	return &commission, nil
}

// GetByOrderID 注文IDで成果報酬を取得
func (r *PostgreSQLCommissionRepository) GetByOrderID(ctx context.Context, orderID string) (*domain.Commission, error) {
	query := `
		SELECT 
			id, order_id, ambassador_id, tenant_id,
			order_amount, commission_rate, commission_amount,
			status, paid_at, created_at, updated_at
		FROM commissions
		WHERE order_id = $1
	`
	
	var commission domain.Commission
	var statusStr string
	var paidAt sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&commission.ID,
		&commission.OrderID,
		&commission.AmbassadorID,
		&commission.TenantID,
		&commission.OrderAmount,
		&commission.CommissionRate,
		&commission.CommissionAmount,
		&statusStr,
		&paidAt,
		&commission.CreatedAt,
		&commission.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil // 成果報酬が無い場合はnilを返す
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get commission by order_id: %w", err)
	}
	
	commission.Status = domain.CommissionStatus(statusStr)
	
	if paidAt.Valid {
		commission.PaidAt = &paidAt.Time
	}
	
	return &commission, nil
}

// GetByAmbassadorID アンバサダーIDで成果報酬一覧を取得
func (r *PostgreSQLCommissionRepository) GetByAmbassadorID(ctx context.Context, ambassadorID string, limit, offset int) ([]*domain.Commission, error) {
	query := `
		SELECT 
			id, order_id, ambassador_id, tenant_id,
			order_amount, commission_rate, commission_amount,
			status, paid_at, created_at, updated_at
		FROM commissions
		WHERE ambassador_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.QueryContext(ctx, query, ambassadorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query commissions: %w", err)
	}
	defer rows.Close()
	
	commissions := make([]*domain.Commission, 0)
	
	for rows.Next() {
		var commission domain.Commission
		var statusStr string
		var paidAt sql.NullTime
		
		err := rows.Scan(
			&commission.ID,
			&commission.OrderID,
			&commission.AmbassadorID,
			&commission.TenantID,
			&commission.OrderAmount,
			&commission.CommissionRate,
			&commission.CommissionAmount,
			&statusStr,
			&paidAt,
			&commission.CreatedAt,
			&commission.UpdatedAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan commission: %w", err)
		}
		
		commission.Status = domain.CommissionStatus(statusStr)
		
		if paidAt.Valid {
			commission.PaidAt = &paidAt.Time
		}
		
		commissions = append(commissions, &commission)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating commissions: %w", err)
	}
	
	return commissions, nil
}

// GetByTenantID テナントIDで成果報酬一覧を取得
func (r *PostgreSQLCommissionRepository) GetByTenantID(ctx context.Context, tenantID string, limit, offset int) ([]*domain.Commission, error) {
	query := `
		SELECT 
			id, order_id, ambassador_id, tenant_id,
			order_amount, commission_rate, commission_amount,
			status, paid_at, created_at, updated_at
		FROM commissions
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.QueryContext(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query commissions: %w", err)
	}
	defer rows.Close()
	
	commissions := make([]*domain.Commission, 0)
	
	for rows.Next() {
		var commission domain.Commission
		var statusStr string
		var paidAt sql.NullTime
		
		err := rows.Scan(
			&commission.ID,
			&commission.OrderID,
			&commission.AmbassadorID,
			&commission.TenantID,
			&commission.OrderAmount,
			&commission.CommissionRate,
			&commission.CommissionAmount,
			&statusStr,
			&paidAt,
			&commission.CreatedAt,
			&commission.UpdatedAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan commission: %w", err)
		}
		
		commission.Status = domain.CommissionStatus(statusStr)
		
		if paidAt.Valid {
			commission.PaidAt = &paidAt.Time
		}
		
		commissions = append(commissions, &commission)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating commissions: %w", err)
	}
	
	return commissions, nil
}

// UpdateStatus 成果報酬ステータスを更新
func (r *PostgreSQLCommissionRepository) UpdateStatus(ctx context.Context, commissionID string, status domain.CommissionStatus) error {
	query := `
		UPDATE commissions SET
			status = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query, commissionID, string(status))
	if err != nil {
		return fmt.Errorf("failed to update commission status: %w", err)
	}
	
	return nil
}

