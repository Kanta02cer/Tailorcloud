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

// FabricRollRepository 反物（Roll）リポジトリインターフェース
type FabricRollRepository interface {
	Create(ctx context.Context, roll *domain.FabricRoll) error
	GetByID(ctx context.Context, rollID string, tenantID string) (*domain.FabricRoll, error)
	GetByRollNumber(ctx context.Context, tenantID string, rollNumber string) (*domain.FabricRoll, error)
	ListByFabricID(ctx context.Context, tenantID string, fabricID string, status *domain.FabricRollStatus) ([]*domain.FabricRoll, error)
	FindAvailableRolls(ctx context.Context, tenantID string, fabricID string, requiredLength float64) ([]*domain.FabricRoll, error)
	Update(ctx context.Context, roll *domain.FabricRoll) error
	UpdateLength(ctx context.Context, rollID string, tenantID string, newLength float64) error
	UpdateStatus(ctx context.Context, rollID string, tenantID string, status domain.FabricRollStatus) error
	Delete(ctx context.Context, rollID string, tenantID string) error
}

// PostgreSQLFabricRollRepository PostgreSQLを使った反物（Roll）リポジトリ実装
type PostgreSQLFabricRollRepository struct {
	db *sql.DB
}

// NewPostgreSQLFabricRollRepository PostgreSQLFabricRollRepositoryのコンストラクタ
func NewPostgreSQLFabricRollRepository(db *sql.DB) FabricRollRepository {
	return &PostgreSQLFabricRollRepository{
		db: db,
	}
}

// Create 反物（Roll）を作成
func (r *PostgreSQLFabricRollRepository) Create(ctx context.Context, roll *domain.FabricRoll) error {
	if roll.ID == "" {
		roll.ID = uuid.New().String()
	}
	
	now := time.Now()
	if roll.CreatedAt.IsZero() {
		roll.CreatedAt = now
	}
	if roll.UpdatedAt.IsZero() {
		roll.UpdatedAt = now
	}
	
	query := `
		INSERT INTO fabric_rolls (
			id, tenant_id, fabric_id, roll_number,
			initial_length, current_length, width,
			supplier_lot_no, received_at, location,
			status, notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		roll.ID,
		roll.TenantID,
		roll.FabricID,
		roll.RollNumber,
		roll.InitialLength,
		roll.CurrentLength,
		roll.Width,
		roll.SupplierLotNo,
		roll.ReceivedAt,
		roll.Location,
		roll.Status,
		roll.Notes,
		roll.CreatedAt,
		roll.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create fabric roll: %w", err)
	}
	
	return nil
}

// GetByID 反物（Roll）IDで取得
func (r *PostgreSQLFabricRollRepository) GetByID(ctx context.Context, rollID string, tenantID string) (*domain.FabricRoll, error) {
	query := `
		SELECT 
			id, tenant_id, fabric_id, roll_number,
			initial_length, current_length, width,
			supplier_lot_no, received_at, location,
			status, notes, created_at, updated_at
		FROM fabric_rolls
		WHERE id = $1 AND tenant_id = $2
	`
	
	var roll domain.FabricRoll
	var width sql.NullFloat64
	var supplierLotNo, location, notes sql.NullString
	var receivedAt sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, rollID, tenantID).Scan(
		&roll.ID,
		&roll.TenantID,
		&roll.FabricID,
		&roll.RollNumber,
		&roll.InitialLength,
		&roll.CurrentLength,
		&width,
		&supplierLotNo,
		&receivedAt,
		&location,
		&roll.Status,
		&notes,
		&roll.CreatedAt,
		&roll.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("fabric roll not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get fabric roll: %w", err)
	}
	
	if width.Valid {
		w := width.Float64
		roll.Width = &w
	}
	if supplierLotNo.Valid {
		roll.SupplierLotNo = &supplierLotNo.String
	}
	if receivedAt.Valid {
		roll.ReceivedAt = &receivedAt.Time
	}
	if location.Valid {
		roll.Location = &location.String
	}
	if notes.Valid {
		roll.Notes = &notes.String
	}
	
	return &roll, nil
}

// GetByRollNumber ロール番号で取得
func (r *PostgreSQLFabricRollRepository) GetByRollNumber(ctx context.Context, tenantID string, rollNumber string) (*domain.FabricRoll, error) {
	query := `
		SELECT 
			id, tenant_id, fabric_id, roll_number,
			initial_length, current_length, width,
			supplier_lot_no, received_at, location,
			status, notes, created_at, updated_at
		FROM fabric_rolls
		WHERE tenant_id = $1 AND roll_number = $2
	`
	
	var roll domain.FabricRoll
	var width sql.NullFloat64
	var supplierLotNo, location, notes sql.NullString
	var receivedAt sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, tenantID, rollNumber).Scan(
		&roll.ID,
		&roll.TenantID,
		&roll.FabricID,
		&roll.RollNumber,
		&roll.InitialLength,
		&roll.CurrentLength,
		&width,
		&supplierLotNo,
		&receivedAt,
		&location,
		&roll.Status,
		&notes,
		&roll.CreatedAt,
		&roll.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("fabric roll not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get fabric roll: %w", err)
	}
	
	if width.Valid {
		w := width.Float64
		roll.Width = &w
	}
	if supplierLotNo.Valid {
		roll.SupplierLotNo = &supplierLotNo.String
	}
	if receivedAt.Valid {
		roll.ReceivedAt = &receivedAt.Time
	}
	if location.Valid {
		roll.Location = &location.String
	}
	if notes.Valid {
		roll.Notes = &notes.String
	}
	
	return &roll, nil
}

// ListByFabricID 生地IDで反物（Roll）一覧を取得
func (r *PostgreSQLFabricRollRepository) ListByFabricID(ctx context.Context, tenantID string, fabricID string, status *domain.FabricRollStatus) ([]*domain.FabricRoll, error) {
	query := `
		SELECT 
			id, tenant_id, fabric_id, roll_number,
			initial_length, current_length, width,
			supplier_lot_no, received_at, location,
			status, notes, created_at, updated_at
		FROM fabric_rolls
		WHERE tenant_id = $1 AND fabric_id = $2
	`
	
	args := []interface{}{tenantID, fabricID}
	if status != nil {
		query += " AND status = $3"
		args = append(args, *status)
	}
	
	query += " ORDER BY created_at DESC"
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list fabric rolls: %w", err)
	}
	defer rows.Close()
	
	var rolls []*domain.FabricRoll
	for rows.Next() {
		var roll domain.FabricRoll
		var width sql.NullFloat64
		var supplierLotNo, location, notes sql.NullString
		var receivedAt sql.NullTime
		
		err := rows.Scan(
			&roll.ID,
			&roll.TenantID,
			&roll.FabricID,
			&roll.RollNumber,
			&roll.InitialLength,
			&roll.CurrentLength,
			&width,
			&supplierLotNo,
			&receivedAt,
			&location,
			&roll.Status,
			&notes,
			&roll.CreatedAt,
			&roll.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fabric roll: %w", err)
		}
		
		if width.Valid {
			w := width.Float64
			roll.Width = &w
		}
		if supplierLotNo.Valid {
			roll.SupplierLotNo = &supplierLotNo.String
		}
		if receivedAt.Valid {
			roll.ReceivedAt = &receivedAt.Time
		}
		if location.Valid {
			roll.Location = &location.String
		}
		if notes.Valid {
			roll.Notes = &notes.String
		}
		
		rolls = append(rolls, &roll)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate fabric rolls: %w", err)
	}
	
	return rolls, nil
}

// FindAvailableRolls 利用可能な反物（Roll）を検索
// requiredLength以上の残り長さを持つAVAILABLE状態の反物を返す
func (r *PostgreSQLFabricRollRepository) FindAvailableRolls(ctx context.Context, tenantID string, fabricID string, requiredLength float64) ([]*domain.FabricRoll, error) {
	query := `
		SELECT 
			id, tenant_id, fabric_id, roll_number,
			initial_length, current_length, width,
			supplier_lot_no, received_at, location,
			status, notes, created_at, updated_at
		FROM fabric_rolls
		WHERE tenant_id = $1 
		  AND fabric_id = $2
		  AND status = 'AVAILABLE'
		  AND current_length >= $3
		ORDER BY current_length ASC
	`
	
	rows, err := r.db.QueryContext(ctx, query, tenantID, fabricID, requiredLength)
	if err != nil {
		return nil, fmt.Errorf("failed to find available fabric rolls: %w", err)
	}
	defer rows.Close()
	
	var rolls []*domain.FabricRoll
	for rows.Next() {
		var roll domain.FabricRoll
		var width sql.NullFloat64
		var supplierLotNo, location, notes sql.NullString
		var receivedAt sql.NullTime
		
		err := rows.Scan(
			&roll.ID,
			&roll.TenantID,
			&roll.FabricID,
			&roll.RollNumber,
			&roll.InitialLength,
			&roll.CurrentLength,
			&width,
			&supplierLotNo,
			&receivedAt,
			&location,
			&roll.Status,
			&notes,
			&roll.CreatedAt,
			&roll.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fabric roll: %w", err)
		}
		
		if width.Valid {
			w := width.Float64
			roll.Width = &w
		}
		if supplierLotNo.Valid {
			roll.SupplierLotNo = &supplierLotNo.String
		}
		if receivedAt.Valid {
			roll.ReceivedAt = &receivedAt.Time
		}
		if location.Valid {
			roll.Location = &location.String
		}
		if notes.Valid {
			roll.Notes = &notes.String
		}
		
		rolls = append(rolls, &roll)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate fabric rolls: %w", err)
	}
	
	return rolls, nil
}

// Update 反物（Roll）を更新
func (r *PostgreSQLFabricRollRepository) Update(ctx context.Context, roll *domain.FabricRoll) error {
	roll.UpdatedAt = time.Now()
	
	query := `
		UPDATE fabric_rolls
		SET roll_number = $3,
		    initial_length = $4,
		    current_length = $5,
		    width = $6,
		    supplier_lot_no = $7,
		    received_at = $8,
		    location = $9,
		    status = $10,
		    notes = $11,
		    updated_at = $12
		WHERE id = $1 AND tenant_id = $2
	`
	
	result, err := r.db.ExecContext(ctx, query,
		roll.ID,
		roll.TenantID,
		roll.RollNumber,
		roll.InitialLength,
		roll.CurrentLength,
		roll.Width,
		roll.SupplierLotNo,
		roll.ReceivedAt,
		roll.Location,
		roll.Status,
		roll.Notes,
		roll.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update fabric roll: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("fabric roll not found or tenant_id mismatch")
	}
	
	return nil
}

// UpdateLength 反物（Roll）の残り長さを更新（引当時など）
func (r *PostgreSQLFabricRollRepository) UpdateLength(ctx context.Context, rollID string, tenantID string, newLength float64) error {
	query := `
		UPDATE fabric_rolls
		SET current_length = $3,
		    status = CASE 
		        WHEN $3 = 0 THEN 'CONSUMED'
		        WHEN status = 'AVAILABLE' THEN 'ALLOCATED'
		        ELSE status
		    END,
		    updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`
	
	result, err := r.db.ExecContext(ctx, query, rollID, tenantID, newLength)
	if err != nil {
		return fmt.Errorf("failed to update fabric roll length: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("fabric roll not found or tenant_id mismatch")
	}
	
	return nil
}

// UpdateStatus 反物（Roll）の状態を更新
func (r *PostgreSQLFabricRollRepository) UpdateStatus(ctx context.Context, rollID string, tenantID string, status domain.FabricRollStatus) error {
	query := `
		UPDATE fabric_rolls
		SET status = $3,
		    updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`
	
	result, err := r.db.ExecContext(ctx, query, rollID, tenantID, status)
	if err != nil {
		return fmt.Errorf("failed to update fabric roll status: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("fabric roll not found or tenant_id mismatch")
	}
	
	return nil
}

// Delete 反物（Roll）を削除
func (r *PostgreSQLFabricRollRepository) Delete(ctx context.Context, rollID string, tenantID string) error {
	query := `
		DELETE FROM fabric_rolls
		WHERE id = $1 AND tenant_id = $2
	`
	
	result, err := r.db.ExecContext(ctx, query, rollID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete fabric roll: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("fabric roll not found or tenant_id mismatch")
	}
	
	return nil
}

