package service

import (
	"context"
	"fmt"

	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/repository"
)

// NewTaxCalculationService TaxCalculationServiceのコンストラクタ
func NewTaxCalculationService(tenantRepo repository.TenantRepository) *TaxCalculationService {
	return &TaxCalculationService{
		tenantRepo: tenantRepo,
	}
}

// TaxCalculationService 税率計算サービス
// インボイス制度対応: 標準税率（10%）と軽減税率（8%）の混在に対応
type TaxCalculationService struct {
	tenantRepo repository.TenantRepository
}

// CalculateTaxRequest 税率計算リクエスト
type CalculateTaxRequest struct {
	TenantID          string
	TaxExcludedAmount int64           // 税抜金額
	TaxRate           domain.TaxRate  // 消費税率（0.10 = 10%, 0.08 = 8%）
}

// CalculateTaxResponse 税率計算レスポンス
type CalculateTaxResponse struct {
	TaxExcludedAmount int64                    // 税抜金額
	TaxAmount         int64                    // 消費税額
	TaxIncludedAmount int64                    // 税込金額
	TaxRate           domain.TaxRate           // 消費税率
	RoundingMethod    domain.TaxRoundingMethod
}

// CalculateTax 消費税額を計算
// テナントの端数処理方法に基づいて計算
func (s *TaxCalculationService) CalculateTax(ctx context.Context, req *CalculateTaxRequest) (*CalculateTaxResponse, error) {
	if req.TaxExcludedAmount < 0 {
		return nil, fmt.Errorf("tax_excluded_amount must be >= 0")
	}
	if req.TaxRate <= 0 || req.TaxRate > 1 {
		return nil, fmt.Errorf("tax_rate must be between 0 and 1")
	}

	// テナント情報を取得（端数処理方法を取得）
	tenant, err := s.tenantRepo.GetByID(ctx, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	// デフォルトの端数処理方法
	roundingMethod := tenant.TaxRoundingMethod
	if roundingMethod == "" {
		roundingMethod = domain.TaxRoundingMethodHalfUp // デフォルトは四捨五入
	}

	// 消費税額を計算
	taxAmount := domain.CalculateTax(req.TaxExcludedAmount, req.TaxRate, roundingMethod)
	taxIncludedAmount := req.TaxExcludedAmount + taxAmount

	return &CalculateTaxResponse{
		TaxExcludedAmount: req.TaxExcludedAmount,
		TaxAmount:         taxAmount,
		TaxIncludedAmount: taxIncludedAmount,
		TaxRate:           req.TaxRate,
		RoundingMethod:    roundingMethod,
	}, nil
}

// CalculateTaxForOrder 注文に対して消費税を計算
// 注文の税抜金額から自動的に消費税額を計算
func (s *TaxCalculationService) CalculateTaxForOrder(ctx context.Context, order *domain.Order) (*CalculateTaxResponse, error) {
	// 税率が指定されていない場合は標準税率（10%）を使用
	taxRate := order.TaxRate
	if taxRate == 0 {
		taxRate = domain.TaxRateStandard
	}

	// 税抜金額を取得（TaxExcludedAmountがあればそれを使用、なければTotalAmountを使用）
	taxExcludedAmount := order.TotalAmount
	if order.TaxExcludedAmount != nil {
		taxExcludedAmount = *order.TaxExcludedAmount
	}

	req := &CalculateTaxRequest{
		TenantID:          order.TenantID,
		TaxExcludedAmount: taxExcludedAmount,
		TaxRate:           taxRate,
	}

	return s.CalculateTax(ctx, req)
}

