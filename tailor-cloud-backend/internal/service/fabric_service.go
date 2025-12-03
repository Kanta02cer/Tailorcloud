package service

import (
	"context"
	"fmt"

	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/repository"
)

// FabricService 生地サービス
type FabricService struct {
	fabricRepo repository.FabricRepository
}

// NewFabricService FabricServiceのコンストラクタ
func NewFabricService(fabricRepo repository.FabricRepository) *FabricService {
	return &FabricService{
		fabricRepo: fabricRepo,
	}
}

// GetFabricRequest 生地取得リクエスト
type GetFabricRequest struct {
	FabricID string
	TenantID string
}

// GetFabric 生地を取得
func (s *FabricService) GetFabric(ctx context.Context, req *GetFabricRequest) (*domain.Fabric, error) {
	fabric, err := s.fabricRepo.GetByID(ctx, req.FabricID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fabric: %w", err)
	}
	
	return fabric, nil
}

// ListFabricsRequest 生地一覧取得リクエスト
type ListFabricsRequest struct {
	TenantID string
	Status   []domain.StockStatus // フィルター: Available, Limited, SoldOut
	Search   string               // 検索キーワード
}

// ListFabrics 生地一覧を取得（フィルター・検索対応）
func (s *FabricService) ListFabrics(ctx context.Context, req *ListFabricsRequest) ([]*domain.Fabric, error) {
	filters := &repository.FabricFilters{
		Status: req.Status,
		Search: req.Search,
	}
	
	fabrics, err := s.fabricRepo.GetAll(ctx, req.TenantID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list fabrics: %w", err)
	}
	
	// 在庫ステータスを再計算（念のため）
	for _, fabric := range fabrics {
		fabric.CalculateStockStatus()
	}
	
	return fabrics, nil
}

// SearchFabricsRequest 生地検索リクエスト
type SearchFabricsRequest struct {
	TenantID string
	Keyword  string
}

// SearchFabrics 生地名で検索
func (s *FabricService) SearchFabrics(ctx context.Context, req *SearchFabricsRequest) ([]*domain.Fabric, error) {
	fabrics, err := s.fabricRepo.Search(ctx, req.TenantID, req.Keyword)
	if err != nil {
		return nil, fmt.Errorf("failed to search fabrics: %w", err)
	}
	
	return fabrics, nil
}

// ReserveFabricRequest 生地確保リクエスト（発注フロー開始）
type ReserveFabricRequest struct {
	FabricID string
	TenantID string
	Amount   float64 // 確保したい数量（メートル）
}

// ReserveFabric 生地を確保（在庫確認・確保）
// Phase 1では簡易実装（在庫チェックのみ）
func (s *FabricService) ReserveFabric(ctx context.Context, req *ReserveFabricRequest) error {
	// 1. 生地を取得
	fabric, err := s.fabricRepo.GetByID(ctx, req.FabricID)
	if err != nil {
		return fmt.Errorf("failed to get fabric: %w", err)
	}
	
	// 2. 在庫ステータスを再計算
	fabric.CalculateStockStatus()
	
	// 3. 在庫チェック
	if fabric.StockStatus == domain.StockStatusSoldOut {
		return fmt.Errorf("fabric is sold out")
	}
	
	// 4. 要求数量チェック
	if req.Amount > fabric.StockAmount {
		return fmt.Errorf("insufficient stock: requested %.2fm, available %.2fm", req.Amount, fabric.StockAmount)
	}
	
	// 5. 最小発注数量チェック
	if req.Amount < fabric.MinimumOrder {
		return fmt.Errorf("amount below minimum order: requested %.2fm, minimum %.2fm", req.Amount, fabric.MinimumOrder)
	}
	
	// Phase 1では、在庫確保は実際には行わない（発注フロー開始のトリガーのみ）
	// Phase 2で、在庫確保機能を実装予定
	
	return nil
}

