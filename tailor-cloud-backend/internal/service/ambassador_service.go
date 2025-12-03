package service

import (
	"context"
	"fmt"

	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/repository"
)

// AmbassadorService アンバサダーサービス
type AmbassadorService struct {
	ambassadorRepo repository.AmbassadorRepository
	commissionRepo repository.CommissionRepository
}

// NewAmbassadorService AmbassadorServiceのコンストラクタ
func NewAmbassadorService(ambassadorRepo repository.AmbassadorRepository, commissionRepo repository.CommissionRepository) *AmbassadorService {
	return &AmbassadorService{
		ambassadorRepo: ambassadorRepo,
		commissionRepo: commissionRepo,
	}
}

// CreateAmbassadorRequest アンバサダー作成リクエスト
type CreateAmbassadorRequest struct {
	TenantID        string  `json:"tenant_id"`
	UserID          string  `json:"user_id"` // Firebase AuthのUserID
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	Phone           string  `json:"phone"`
	CommissionRate  float64 `json:"commission_rate"` // 成果報酬率（オプション、デフォルト10%）
}

// CreateAmbassador アンバサダーを作成
func (s *AmbassadorService) CreateAmbassador(ctx context.Context, req *CreateAmbassadorRequest) (*domain.Ambassador, error) {
	// バリデーション
	if req.TenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}
	if req.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	
	// 既存のアンバサダーをチェック
	existing, err := s.ambassadorRepo.GetByUserID(ctx, req.UserID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("ambassador already exists for user_id: %s", req.UserID)
	}
	
	// 成果報酬率の設定（デフォルト10%）
	commissionRate := req.CommissionRate
	if commissionRate == 0 {
		commissionRate = 0.10 // 10%
	}
	if commissionRate < 0 || commissionRate > 1 {
		return nil, fmt.Errorf("commission_rate must be between 0 and 1")
	}
	
	// アンバサダーを作成
	ambassador := domain.NewAmbassador(req.TenantID, req.UserID, req.Name, req.Email)
	ambassador.Phone = req.Phone
	ambassador.CommissionRate = commissionRate
	
	if err := s.ambassadorRepo.Create(ctx, ambassador); err != nil {
		return nil, fmt.Errorf("failed to create ambassador: %w", err)
	}
	
	return ambassador, nil
}

// GetAmbassadorByUserID ユーザーIDでアンバサダーを取得
func (s *AmbassadorService) GetAmbassadorByUserID(ctx context.Context, userID string) (*domain.Ambassador, error) {
	ambassador, err := s.ambassadorRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ambassador: %w", err)
	}
	
	return ambassador, nil
}

// ListAmbassadorsRequest アンバサダー一覧取得リクエスト
type ListAmbassadorsRequest struct {
	TenantID string
}

// ListAmbassadors アンバサダー一覧を取得
func (s *AmbassadorService) ListAmbassadors(ctx context.Context, req *ListAmbassadorsRequest) ([]*domain.Ambassador, error) {
	ambassadors, err := s.ambassadorRepo.GetByTenantID(ctx, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list ambassadors: %w", err)
	}
	
	return ambassadors, nil
}

// CreateCommissionForOrder 注文に対して成果報酬を作成
// 注文作成時に自動的に呼び出される
func (s *AmbassadorService) CreateCommissionForOrder(ctx context.Context, order *domain.Order) (*domain.Commission, error) {
	// 注文のCreatedByからアンバサダーを取得
	ambassador, err := s.ambassadorRepo.GetByUserID(ctx, order.CreatedBy)
	if err != nil {
		// アンバサダーが存在しない場合は成果報酬を作成しない（プロが直接注文した場合など）
		return nil, nil
	}
	
	// 成果報酬を作成（Pendingステータス）
	commission := domain.NewCommission(
		order.ID,
		ambassador.ID,
		order.TenantID,
		order.TotalAmount,
		ambassador.CommissionRate,
	)
	
	if err := s.commissionRepo.Create(ctx, commission); err != nil {
		return nil, fmt.Errorf("failed to create commission: %w", err)
	}
	
	return commission, nil
}

// ApproveCommission 成果報酬を確定（注文確定時に呼び出される）
func (s *AmbassadorService) ApproveCommission(ctx context.Context, orderID string) error {
	commission, err := s.commissionRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get commission: %w", err)
	}
	
	// 成果報酬が存在しない場合は何もしない
	if commission == nil {
		return nil
	}
	
	// ステータスをApprovedに変更
	if err := s.commissionRepo.UpdateStatus(ctx, commission.ID, domain.CommissionStatusApproved); err != nil {
		return fmt.Errorf("failed to approve commission: %w", err)
	}
	
	// アンバサダーの売上統計を更新
	if err := s.ambassadorRepo.UpdateSalesStats(ctx, commission.AmbassadorID, commission.OrderAmount, commission.CommissionAmount); err != nil {
		// エラーはログに記録するが、メイン処理は継続
		fmt.Printf("WARNING: Failed to update sales stats: %v\n", err)
	}
	
	return nil
}

// GetCommissionsByAmbassadorRequest アンバサダーの成果報酬一覧取得リクエスト
type GetCommissionsByAmbassadorRequest struct {
	AmbassadorID string
	Limit        int
	Offset       int
}

// GetCommissionsByAmbassador アンバサダーの成果報酬一覧を取得
func (s *AmbassadorService) GetCommissionsByAmbassador(ctx context.Context, req *GetCommissionsByAmbassadorRequest) ([]*domain.Commission, error) {
	if req.Limit <= 0 {
		req.Limit = 20 // デフォルト
	}
	
	commissions, err := s.commissionRepo.GetByAmbassadorID(ctx, req.AmbassadorID, req.Limit, req.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get commissions: %w", err)
	}
	
	return commissions, nil
}

