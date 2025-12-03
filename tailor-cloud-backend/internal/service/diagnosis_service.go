package service

import (
	"context"
	"encoding/json"
	"fmt"

	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/repository"
)

// DiagnosisService 診断サービス
type DiagnosisService struct {
	diagnosisRepo repository.DiagnosisRepository
}

// NewDiagnosisService DiagnosisServiceのコンストラクタ
func NewDiagnosisService(diagnosisRepo repository.DiagnosisRepository) *DiagnosisService {
	return &DiagnosisService{
		diagnosisRepo: diagnosisRepo,
	}
}

// CreateDiagnosisRequest 診断作成リクエスト
type CreateDiagnosisRequest struct {
	UserID          string
	TenantID        string
	Archetype       domain.Archetype
	PlanType        domain.PlanType
	DiagnosisResult json.RawMessage
}

// CreateDiagnosis 診断を作成
func (s *DiagnosisService) CreateDiagnosis(ctx context.Context, req *CreateDiagnosisRequest) (*domain.Diagnosis, error) {
	// バリデーション
	if req.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if req.TenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}
	if !req.Archetype.IsValid() {
		return nil, fmt.Errorf("invalid archetype: %s", req.Archetype)
	}
	if req.PlanType != "" && !req.PlanType.IsValid() {
		return nil, fmt.Errorf("invalid plan_type: %s", req.PlanType)
	}

	// 診断オブジェクトを作成
	diagnosis := domain.NewDiagnosis(
		req.UserID,
		req.TenantID,
		req.Archetype,
		req.PlanType,
		req.DiagnosisResult,
	)

	// リポジトリに保存
	if err := s.diagnosisRepo.Create(ctx, diagnosis); err != nil {
		return nil, fmt.Errorf("failed to create diagnosis: %w", err)
	}

	return diagnosis, nil
}

// GetDiagnosis 診断を取得
func (s *DiagnosisService) GetDiagnosis(ctx context.Context, diagnosisID string, tenantID string) (*domain.Diagnosis, error) {
	if diagnosisID == "" {
		return nil, fmt.Errorf("diagnosis_id is required")
	}
	if tenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}

	diagnosis, err := s.diagnosisRepo.GetByID(ctx, diagnosisID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get diagnosis: %w", err)
	}

	return diagnosis, nil
}

// GetDiagnosesByUser ユーザーの診断一覧を取得
func (s *DiagnosisService) GetDiagnosesByUser(ctx context.Context, userID string, tenantID string) ([]*domain.Diagnosis, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if tenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}

	diagnoses, err := s.diagnosisRepo.GetByUserID(ctx, userID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get diagnoses by user: %w", err)
	}

	return diagnoses, nil
}

// GetLatestByUserID ユーザーの最新の診断を取得
func (s *DiagnosisService) GetLatestByUserID(ctx context.Context, userID string, tenantID string) (*domain.Diagnosis, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if tenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}

	diagnoses, err := s.diagnosisRepo.GetByUserID(ctx, userID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get diagnoses by user: %w", err)
	}

	if len(diagnoses) == 0 {
		return nil, fmt.Errorf("no diagnosis found for user")
	}

	// 最新の診断を返す（GetByUserIDはcreated_at DESCでソート済み）
	return diagnoses[0], nil
}

// GetDiagnosesByTenantRequest テナントの診断一覧取得リクエスト
type GetDiagnosesByTenantRequest struct {
	TenantID string
	Limit    int
	Offset   int
}

// GetDiagnosesByTenant テナントの診断一覧を取得（ページネーション対応）
func (s *DiagnosisService) GetDiagnosesByTenant(ctx context.Context, req *GetDiagnosesByTenantRequest) ([]*domain.Diagnosis, error) {
	if req.TenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20 // デフォルト値
	}
	if limit > 100 {
		limit = 100 // 最大値
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	diagnoses, err := s.diagnosisRepo.GetByTenantID(ctx, req.TenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get diagnoses by tenant: %w", err)
	}

	return diagnoses, nil
}

// ListDiagnosesRequest 診断一覧取得リクエスト（フィルター対応）
type ListDiagnosesRequest struct {
	TenantID  string
	Archetype *domain.Archetype
	PlanType  *domain.PlanType
	StartDate *string // ISO 8601形式
	EndDate   *string // ISO 8601形式
	Limit     int
	Offset    int
}

// ListDiagnoses 診断一覧を取得（フィルター対応）
func (s *DiagnosisService) ListDiagnoses(ctx context.Context, req *ListDiagnosesRequest) ([]*domain.Diagnosis, error) {
	if req.TenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}

	filter := repository.DiagnosisFilter{
		Archetype: req.Archetype,
		PlanType:  req.PlanType,
		Limit:     req.Limit,
		Offset:    req.Offset,
	}

	// 日付フィルターの変換（必要に応じて実装）
	// if req.StartDate != nil {
	//     startDate, err := time.Parse(time.RFC3339, *req.StartDate)
	//     if err != nil {
	//         return nil, fmt.Errorf("invalid start_date format: %w", err)
	//     }
	//     filter.StartDate = &startDate
	// }

	diagnoses, err := s.diagnosisRepo.List(ctx, req.TenantID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list diagnoses: %w", err)
	}

	return diagnoses, nil
}

