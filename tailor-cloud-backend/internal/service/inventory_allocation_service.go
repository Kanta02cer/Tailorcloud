package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/repository"
)

// InventoryAllocationService 在庫引当サービス
// エンタープライズ実装の核心: 反物（Roll）単位での在庫引当ロジック
type InventoryAllocationService struct {
	fabricRollRepo        repository.FabricRollRepository
	fabricAllocationRepo  repository.FabricAllocationRepository
	fabricRepo            repository.FabricRepository
	db                    *sql.DB // トランザクション管理用
}

// NewInventoryAllocationService InventoryAllocationServiceのコンストラクタ
func NewInventoryAllocationService(
	fabricRollRepo repository.FabricRollRepository,
	fabricAllocationRepo repository.FabricAllocationRepository,
	fabricRepo repository.FabricRepository,
	db *sql.DB,
) *InventoryAllocationService {
	return &InventoryAllocationService{
		fabricRollRepo:       fabricRollRepo,
		fabricAllocationRepo: fabricAllocationRepo,
		fabricRepo:           fabricRepo,
		db:                   db,
	}
}

// AllocationStrategy 引当戦略
type AllocationStrategy string

const (
	AllocationStrategyFIFO AllocationStrategy = "FIFO" // First In First Out（古い反物から）
	AllocationStrategyLIFO AllocationStrategy = "LIFO" // Last In First Out（新しい反物から）
	AllocationStrategyBestFit AllocationStrategy = "BEST_FIT" // 最適フィット（最小の無駄）
)

// AllocateInventoryRequest 在庫引当リクエスト
type AllocateInventoryRequest struct {
	TenantID        string
	OrderID         string
	FabricID        string
	RequiredLength  float64 // 必要な長さ（メートル）
	Strategy        AllocationStrategy
}

// AllocateInventoryResponse 在庫引当レスポンス
type AllocateInventoryResponse struct {
	Allocations     []*domain.FabricAllocation
	TotalAllocated  float64
	RemainingNeeded float64
}

// AllocateInventory 在庫を引当（反物単位で管理）
// エンタープライズ実装の核心: 物理的な反物から必要な長さを引当
func (s *InventoryAllocationService) AllocateInventory(ctx context.Context, req *AllocateInventoryRequest) (*AllocateInventoryResponse, error) {
	if req.RequiredLength <= 0 {
		return nil, fmt.Errorf("required_length must be greater than 0")
	}
	
	// トランザクション開始（排他制御のため）
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	// 生地が存在するか確認
	_, err = s.fabricRepo.GetByID(ctx, req.FabricID)
	if err != nil {
		return nil, fmt.Errorf("fabric not found: %w", err)
	}
	
	// 利用可能な反物を検索（ロック付きで取得）
	availableRolls, err := s.findAvailableRollsWithLock(ctx, tx, req.TenantID, req.FabricID, req.RequiredLength)
	if err != nil {
		return nil, fmt.Errorf("failed to find available rolls: %w", err)
	}
	
	if len(availableRolls) == 0 {
		return nil, fmt.Errorf("insufficient inventory: no available fabric rolls for fabric_id=%s, required_length=%.2fm", req.FabricID, req.RequiredLength)
	}
	
	// 引当戦略に基づいて反物を選択
	selectedRolls := s.selectRollsByStrategy(availableRolls, req.RequiredLength, req.Strategy)
	if len(selectedRolls) == 0 {
		return nil, fmt.Errorf("insufficient inventory: cannot allocate %.2fm from available rolls", req.RequiredLength)
	}
	
	// 反物から引当
	var allocations []*domain.FabricAllocation
	remainingNeeded := req.RequiredLength
	totalAllocated := 0.0
	
	for _, roll := range selectedRolls {
		if remainingNeeded <= 0 {
			break
		}
		
		// この反物から引当できる長さ
		allocateLength := remainingNeeded
		if roll.CurrentLength < remainingNeeded {
			allocateLength = roll.CurrentLength
		}
		
		// 反物の残り長さを更新（ロック中なので安全）
		newLength := roll.CurrentLength - allocateLength
		err := s.updateRollLengthInTx(ctx, tx, roll.ID, req.TenantID, newLength)
		if err != nil {
			return nil, fmt.Errorf("failed to update roll length: %w", err)
		}
		
		// 引当レコードを作成
		allocation := &domain.FabricAllocation{
			ID:              uuid.New().String(),
			TenantID:        req.TenantID,
			OrderID:         req.OrderID,
			FabricRollID:    roll.ID,
			AllocatedLength: allocateLength,
			Status:          domain.FabricAllocationStatusReserved,
			AllocatedAt:     time.Now(),
		}
		
		// 引当を保存（トランザクション内）
		err = s.createAllocationInTx(ctx, tx, allocation)
		if err != nil {
			return nil, fmt.Errorf("failed to create allocation: %w", err)
		}
		
		allocations = append(allocations, allocation)
		remainingNeeded -= allocateLength
		totalAllocated += allocateLength
	}
	
	// まだ必要量に満たない場合
	if remainingNeeded > 0.01 { // 0.01mの誤差は許容
		return nil, fmt.Errorf("insufficient inventory: still need %.2fm after allocating from all available rolls", remainingNeeded)
	}
	
	// トランザクションコミット
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return &AllocateInventoryResponse{
		Allocations:     allocations,
		TotalAllocated:  totalAllocated,
		RemainingNeeded: remainingNeeded,
	}, nil
}

// findAvailableRollsWithLock 利用可能な反物を検索（ロック付き）
// SELECT FOR UPDATE で行ロックを取得し、同時発注時の重複引当を防止
func (s *InventoryAllocationService) findAvailableRollsWithLock(ctx context.Context, tx *sql.Tx, tenantID string, fabricID string, requiredLength float64) ([]*domain.FabricRoll, error) {
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
		FOR UPDATE SKIP LOCKED
	`
	
	rows, err := tx.QueryContext(ctx, query, tenantID, fabricID, requiredLength)
	if err != nil {
		return nil, fmt.Errorf("failed to query available rolls: %w", err)
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
			return nil, fmt.Errorf("failed to scan roll: %w", err)
		}
		
		if width.Valid {
			roll.Width = &width.Float64
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
	
	return rolls, nil
}

// selectRollsByStrategy 引当戦略に基づいて反物を選択
func (s *InventoryAllocationService) selectRollsByStrategy(rolls []*domain.FabricRoll, requiredLength float64, strategy AllocationStrategy) []*domain.FabricRoll {
	var selected []*domain.FabricRoll
	remaining := requiredLength
	
	switch strategy {
	case AllocationStrategyFIFO:
		// 古い反物から（received_atが古い順、またはcreated_atが古い順）
		// 既にORDER BYでソートされているのでそのまま使用
		for _, roll := range rolls {
			if remaining <= 0 {
				break
			}
			selected = append(selected, roll)
			if roll.CurrentLength >= remaining {
				break
			}
			remaining -= roll.CurrentLength
		}
		
	case AllocationStrategyLIFO:
		// 新しい反物から（逆順）
		for i := len(rolls) - 1; i >= 0; i-- {
			roll := rolls[i]
			if remaining <= 0 {
				break
			}
			selected = append(selected, roll)
			if roll.CurrentLength >= remaining {
				break
			}
			remaining -= roll.CurrentLength
		}
		
	case AllocationStrategyBestFit:
		// 最適フィット: 必要な長さに最も近い反物を選択（最小の無駄）
		selected = s.selectBestFit(rolls, requiredLength)
		
	default:
		// デフォルトはFIFO
		return s.selectRollsByStrategy(rolls, requiredLength, AllocationStrategyFIFO)
	}
	
	return selected
}

// selectBestFit 最適フィットアルゴリズム
// 必要な長さに最も近い反物を選択し、無駄を最小化
func (s *InventoryAllocationService) selectBestFit(rolls []*domain.FabricRoll, requiredLength float64) []*domain.FabricRoll {
	// シンプルな実装: 1本の反物で足りる場合はそれを使用
	// 複数の反物が必要な場合は、最小の無駄になる組み合わせを選択
	var selected []*domain.FabricRoll
	remaining := requiredLength
	
	// 1本で足りる反物を探す
	for _, roll := range rolls {
		if roll.CurrentLength >= remaining {
			// この反物だけで足りる
			selected = append(selected, roll)
			return selected
		}
	}
	
	// 複数の反物が必要な場合: 残り長さが最小になるように選択
	for _, roll := range rolls {
		if remaining <= 0 {
			break
		}
		if roll.CurrentLength > remaining {
			// この反物で足りるので、これで終了
			selected = append(selected, roll)
			break
		}
		selected = append(selected, roll)
		remaining -= roll.CurrentLength
	}
	
	return selected
}

// updateRollLengthInTx トランザクション内で反物の残り長さを更新
func (s *InventoryAllocationService) updateRollLengthInTx(ctx context.Context, tx *sql.Tx, rollID string, tenantID string, newLength float64) error {
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
	
	result, err := tx.ExecContext(ctx, query, rollID, tenantID, newLength)
	if err != nil {
		return fmt.Errorf("failed to update roll length: %w", err)
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

// createAllocationInTx トランザクション内で引当レコードを作成
func (s *InventoryAllocationService) createAllocationInTx(ctx context.Context, tx *sql.Tx, allocation *domain.FabricAllocation) error {
	query := `
		INSERT INTO fabric_allocations (
			id, tenant_id, order_id, order_item_id, fabric_roll_id,
			allocated_length, actual_used_length, remnant_length,
			allocation_status, allocated_at, confirmed_at, cut_at,
			notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	
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
	
	_, err := tx.ExecContext(ctx, query,
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
		return fmt.Errorf("failed to create allocation: %w", err)
	}
	
	return nil
}

// ReleaseAllocation 引当を解除（キャンセル時など）
func (s *InventoryAllocationService) ReleaseAllocation(ctx context.Context, allocationID string, tenantID string) error {
	// 引当を取得
	allocation, err := s.fabricAllocationRepo.GetByID(ctx, allocationID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get allocation: %w", err)
	}
	
	if allocation.Status == domain.FabricAllocationStatusCancelled {
		return nil // 既にキャンセル済み
	}
	
	// トランザクション開始
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	// 反物の残り長さを戻す
	roll, err := s.fabricRollRepo.GetByID(ctx, allocation.FabricRollID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get roll: %w", err)
	}
	
	newLength := roll.CurrentLength + allocation.AllocatedLength
	err = s.updateRollLengthInTx(ctx, tx, roll.ID, tenantID, newLength)
	if err != nil {
		return fmt.Errorf("failed to restore roll length: %w", err)
	}
	
	// 引当をキャンセル
	allocation.Cancel()
	err = s.fabricAllocationRepo.Update(ctx, allocation)
	if err != nil {
		return fmt.Errorf("failed to cancel allocation: %w", err)
	}
	
	// トランザクションコミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

