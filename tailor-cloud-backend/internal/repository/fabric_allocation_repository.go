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

// FabricAllocationRepository 反物引当リポジトリインターフェース
type FabricAllocationRepository interface {
	Create(ctx context.Context, allocation *domain.FabricAllocation) error
	GetByID(ctx context.Context, allocationID string, tenantID string) (*domain.FabricAllocation, error)
	GetByOrderID(ctx context.Context, orderID string, tenantID string) ([]*domain.FabricAllocation, error)
	GetByFabricRollID(ctx context.Context, fabricRollID string, tenantID string) ([]*domain.FabricAllocation, error)
	Update(ctx context.Context, allocation *domain.FabricAllocation) error
	UpdateStatus(ctx context.Context, allocationID string, tenantID string, status domain.FabricAllocationStatus) error
	Delete(ctx context.Context, allocationID string, tenantID string) error
}

// PostgreSQLFabricAllocationRepository PostgreSQLを使った反物引当リポジトリ実装
type PostgreSQLFabricAllocationRepository struct {
	db *sql.DB
}

// NewPostgreSQLFabricAllocationRepository PostgreSQLFabricAllocationRepositoryのコンストラクタ
func NewPostgreSQLFabricAllocationRepository(db *sql.DB) FabricAllocationRepository {
	return &PostgreSQLFabricAllocationRepository{
		db: db,
	}
}

// Create 反物引当を作成
func (r *PostgreSQLFabricAllocationRepository) Create(ctx context.Context, allocation *domain.FabricAllocation) error {
	if allocation.ID == "" {
		allocation.ID = uuid.New().String()
	}
	
	now := time.Now()
	if allocation.CreatedAt.IsZero() {
		allocation.CreatedAt = now
	}
	if allocation.UpdatedAt.IsZero() {
		allocation.UpdatedAt = now
	}
	if allocation.AllocatedAt.IsZero() {
		allocation.AllocatedAt = now
	}
	
	query := `
		INSERT INTO fabric_allocations (
			id, tenant_id, order_id, order_item_id, fabric_roll_id,
			allocated_length, actual_used_length, remnant_length,
			allocation_status, allocated_at, confirmed_at, cut_at,
			notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		allocation.ID,
		allocation.TenantID,
		allocation.OrderID,
		allocation.OrderItemID,
		allocation.FabricRollID,
		allocation.AllocatedLength,
		allocation.ActualUsedLength,
		allocation.RemnantLength,
		allocation.Status,
		allocation.AllocatedAt,
		allocation.ConfirmedAt,
		allocation.CutAt,
		allocation.Notes,
		allocation.CreatedAt,
		allocation.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create fabric allocation: %w", err)
	}
	
	return nil
}

// GetByID 反物引当IDで取得
func (r *PostgreSQLFabricAllocationRepository) GetByID(ctx context.Context, allocationID string, tenantID string) (*domain.FabricAllocation, error) {
	query := `
		SELECT 
			id, tenant_id, order_id, order_item_id, fabric_roll_id,
			allocated_length, actual_used_length, remnant_length,
			allocation_status, allocated_at, confirmed_at, cut_at,
			notes, created_at, updated_at
		FROM fabric_allocations
		WHERE id = $1 AND tenant_id = $2
	`
	
	var allocation domain.FabricAllocation
	var orderItemID, notes sql.NullString
	var actualUsedLength, remnantLength sql.NullFloat64
	var confirmedAt, cutAt sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, allocationID, tenantID).Scan(
		&allocation.ID,
		&allocation.TenantID,
		&allocation.OrderID,
		&orderItemID,
		&allocation.FabricRollID,
		&allocation.AllocatedLength,
		&actualUsedLength,
		&remnantLength,
		&allocation.Status,
		&allocation.AllocatedAt,
		&confirmedAt,
		&cutAt,
		&notes,
		&allocation.CreatedAt,
		&allocation.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("fabric allocation not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get fabric allocation: %w", err)
	}
	
	if orderItemID.Valid {
		allocation.OrderItemID = &orderItemID.String
	}
	if actualUsedLength.Valid {
		allocation.ActualUsedLength = &actualUsedLength.Float64
	}
	if remnantLength.Valid {
		allocation.RemnantLength = &remnantLength.Float64
	}
	if confirmedAt.Valid {
		allocation.ConfirmedAt = &confirmedAt.Time
	}
	if cutAt.Valid {
		allocation.CutAt = &cutAt.Time
	}
	if notes.Valid {
		allocation.Notes = &notes.String
	}
	
	return &allocation, nil
}

// GetByOrderID 注文IDで反物引当一覧を取得
func (r *PostgreSQLFabricAllocationRepository) GetByOrderID(ctx context.Context, orderID string, tenantID string) ([]*domain.FabricAllocation, error) {
	query := `
		SELECT 
			id, tenant_id, order_id, order_item_id, fabric_roll_id,
			allocated_length, actual_used_length, remnant_length,
			allocation_status, allocated_at, confirmed_at, cut_at,
			notes, created_at, updated_at
		FROM fabric_allocations
		WHERE order_id = $1 AND tenant_id = $2
		ORDER BY allocated_at ASC
	`
	
	rows, err := r.db.QueryContext(ctx, query, orderID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fabric allocations by order: %w", err)
	}
	defer rows.Close()
	
	var allocations []*domain.FabricAllocation
	for rows.Next() {
		var allocation domain.FabricAllocation
		var orderItemID, notes sql.NullString
		var actualUsedLength, remnantLength sql.NullFloat64
		var confirmedAt, cutAt sql.NullTime
		
		err := rows.Scan(
			&allocation.ID,
			&allocation.TenantID,
			&allocation.OrderID,
			&orderItemID,
			&allocation.FabricRollID,
			&allocation.AllocatedLength,
			&actualUsedLength,
			&remnantLength,
			&allocation.Status,
			&allocation.AllocatedAt,
			&confirmedAt,
			&cutAt,
			&notes,
			&allocation.CreatedAt,
			&allocation.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fabric allocation: %w", err)
		}
		
		if orderItemID.Valid {
			allocation.OrderItemID = &orderItemID.String
		}
		if actualUsedLength.Valid {
			allocation.ActualUsedLength = &actualUsedLength.Float64
		}
		if remnantLength.Valid {
			allocation.RemnantLength = &remnantLength.Float64
		}
		if confirmedAt.Valid {
			allocation.ConfirmedAt = &confirmedAt.Time
		}
		if cutAt.Valid {
			allocation.CutAt = &cutAt.Time
		}
		if notes.Valid {
			allocation.Notes = &notes.String
		}
		
		allocations = append(allocations, &allocation)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate fabric allocations: %w", err)
	}
	
	return allocations, nil
}

// GetByFabricRollID 反物（Roll）IDで反物引当一覧を取得
func (r *PostgreSQLFabricAllocationRepository) GetByFabricRollID(ctx context.Context, fabricRollID string, tenantID string) ([]*domain.FabricAllocation, error) {
	query := `
		SELECT 
			id, tenant_id, order_id, order_item_id, fabric_roll_id,
			allocated_length, actual_used_length, remnant_length,
			allocation_status, allocated_at, confirmed_at, cut_at,
			notes, created_at, updated_at
		FROM fabric_allocations
		WHERE fabric_roll_id = $1 AND tenant_id = $2
		  AND allocation_status != 'CANCELLED'
		ORDER BY allocated_at ASC
	`
	
	rows, err := r.db.QueryContext(ctx, query, fabricRollID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fabric allocations by roll: %w", err)
	}
	defer rows.Close()
	
	var allocations []*domain.FabricAllocation
	for rows.Next() {
		var allocation domain.FabricAllocation
		var orderItemID, notes sql.NullString
		var actualUsedLength, remnantLength sql.NullFloat64
		var confirmedAt, cutAt sql.NullTime
		
		err := rows.Scan(
			&allocation.ID,
			&allocation.TenantID,
			&allocation.OrderID,
			&orderItemID,
			&allocation.FabricRollID,
			&allocation.AllocatedLength,
			&actualUsedLength,
			&remnantLength,
			&allocation.Status,
			&allocation.AllocatedAt,
			&confirmedAt,
			&cutAt,
			&notes,
			&allocation.CreatedAt,
			&allocation.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fabric allocation: %w", err)
		}
		
		if orderItemID.Valid {
			allocation.OrderItemID = &orderItemID.String
		}
		if actualUsedLength.Valid {
			allocation.ActualUsedLength = &actualUsedLength.Float64
		}
		if remnantLength.Valid {
			allocation.RemnantLength = &remnantLength.Float64
		}
		if confirmedAt.Valid {
			allocation.ConfirmedAt = &confirmedAt.Time
		}
		if cutAt.Valid {
			allocation.CutAt = &cutAt.Time
		}
		if notes.Valid {
			allocation.Notes = &notes.String
		}
		
		allocations = append(allocations, &allocation)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate fabric allocations: %w", err)
	}
	
	return allocations, nil
}

// Update 反物引当を更新
func (r *PostgreSQLFabricAllocationRepository) Update(ctx context.Context, allocation *domain.FabricAllocation) error {
	allocation.UpdatedAt = time.Now()
	
	query := `
		UPDATE fabric_allocations
		SET allocated_length = $3,
		    actual_used_length = $4,
		    remnant_length = $5,
		    allocation_status = $6,
		    confirmed_at = $7,
		    cut_at = $8,
		    notes = $9,
		    updated_at = $10
		WHERE id = $1 AND tenant_id = $2
	`
	
	result, err := r.db.ExecContext(ctx, query,
		allocation.ID,
		allocation.TenantID,
		allocation.AllocatedLength,
		allocation.ActualUsedLength,
		allocation.RemnantLength,
		allocation.Status,
		allocation.ConfirmedAt,
		allocation.CutAt,
		allocation.Notes,
		allocation.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update fabric allocation: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("fabric allocation not found or tenant_id mismatch")
	}
	
	return nil
}

// UpdateStatus 反物引当の状態を更新
func (r *PostgreSQLFabricAllocationRepository) UpdateStatus(ctx context.Context, allocationID string, tenantID string, status domain.FabricAllocationStatus) error {
	query := `
		UPDATE fabric_allocations
		SET allocation_status = $3,
		    updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`
	
	result, err := r.db.ExecContext(ctx, query, allocationID, tenantID, status)
	if err != nil {
		return fmt.Errorf("failed to update fabric allocation status: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("fabric allocation not found or tenant_id mismatch")
	}
	
	return nil
}

// Delete 反物引当を削除
func (r *PostgreSQLFabricAllocationRepository) Delete(ctx context.Context, allocationID string, tenantID string) error {
	query := `
		DELETE FROM fabric_allocations
		WHERE id = $1 AND tenant_id = $2
	`
	
	result, err := r.db.ExecContext(ctx, query, allocationID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete fabric allocation: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("fabric allocation not found or tenant_id mismatch")
	}
	
	return nil
}

