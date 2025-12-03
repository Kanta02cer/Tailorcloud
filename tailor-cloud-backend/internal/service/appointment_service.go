package service

import (
	"context"
	"fmt"
	"time"

	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/repository"
)

// AppointmentService 予約サービス
type AppointmentService struct {
	appointmentRepo repository.AppointmentRepository
}

// NewAppointmentService AppointmentServiceのコンストラクタ
func NewAppointmentService(appointmentRepo repository.AppointmentRepository) *AppointmentService {
	return &AppointmentService{
		appointmentRepo: appointmentRepo,
	}
}

// CreateAppointmentRequest 予約作成リクエスト
type CreateAppointmentRequest struct {
	UserID              string
	TenantID            string
	FitterID            string
	AppointmentDateTime time.Time
	DurationMinutes     int
	Notes               string
}

// CreateAppointment 予約を作成
func (s *AppointmentService) CreateAppointment(ctx context.Context, req *CreateAppointmentRequest) (*domain.Appointment, error) {
	// バリデーション
	if req.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if req.TenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}
	if req.FitterID == "" {
		return nil, fmt.Errorf("fitter_id is required")
	}
	if req.AppointmentDateTime.IsZero() {
		return nil, fmt.Errorf("appointment_datetime is required")
	}
	if req.AppointmentDateTime.Before(time.Now()) {
		return nil, fmt.Errorf("appointment_datetime must be in the future")
	}
	if req.DurationMinutes <= 0 {
		req.DurationMinutes = 60 // デフォルト60分
	}
	if req.DurationMinutes > 480 { // 最大8時間
		return nil, fmt.Errorf("duration_minutes must be less than or equal to 480")
	}

	// 空き状況チェック
	available, err := s.appointmentRepo.CheckAvailability(ctx, req.FitterID, req.TenantID, req.AppointmentDateTime, req.DurationMinutes)
	if err != nil {
		return nil, fmt.Errorf("failed to check availability: %w", err)
	}
	if !available {
		return nil, fmt.Errorf("appointment time slot is not available")
	}

	// 予約オブジェクトを作成
	appointment := domain.NewAppointment(
		req.UserID,
		req.TenantID,
		req.FitterID,
		req.AppointmentDateTime,
		req.DurationMinutes,
	)
	appointment.Notes = req.Notes
	appointment.Status = domain.AppointmentStatusPending

	// リポジトリに保存
	if err := s.appointmentRepo.Create(ctx, appointment); err != nil {
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	return appointment, nil
}

// GetAppointment 予約を取得
func (s *AppointmentService) GetAppointment(ctx context.Context, appointmentID string, tenantID string) (*domain.Appointment, error) {
	if appointmentID == "" {
		return nil, fmt.Errorf("appointment_id is required")
	}
	if tenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}

	appointment, err := s.appointmentRepo.GetByID(ctx, appointmentID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	return appointment, nil
}

// ListAppointmentsRequest 予約一覧取得リクエスト
type ListAppointmentsRequest struct {
	TenantID  string
	UserID    string   // オプション: ユーザーIDでフィルター
	FitterID  string   // オプション: フィッターIDでフィルター
	StartDate *time.Time // オプション: 開始日
	EndDate   *time.Time // オプション: 終了日
}

// ListAppointments 予約一覧を取得
func (s *AppointmentService) ListAppointments(ctx context.Context, req *ListAppointmentsRequest) ([]*domain.Appointment, error) {
	if req.TenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}

	// デフォルトの期間設定（開始日・終了日が指定されていない場合）
	startDate := req.StartDate
	endDate := req.EndDate

	if startDate == nil {
		now := time.Now()
		startDate = &now // 今日から
	}
	if endDate == nil {
		future := startDate.AddDate(0, 1, 0) // 1ヶ月後まで
		endDate = &future
	}

	var appointments []*domain.Appointment
	var err error

	if req.UserID != "" {
		// ユーザーIDで取得
		appointments, err = s.appointmentRepo.GetByUserID(ctx, req.UserID, req.TenantID)
		if err != nil {
			return nil, fmt.Errorf("failed to get appointments by user: %w", err)
		}
	} else if req.FitterID != "" {
		// フィッターIDで取得（期間指定）
		appointments, err = s.appointmentRepo.GetByFitterID(ctx, req.FitterID, req.TenantID, *startDate, *endDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get appointments by fitter: %w", err)
		}
	} else {
		// テナントIDで取得（期間指定）
		appointments, err = s.appointmentRepo.GetByTenantID(ctx, req.TenantID, *startDate, *endDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get appointments by tenant: %w", err)
		}
	}

	// 期間フィルター（リポジトリで対応していない場合のフォールバック）
	if req.UserID != "" && (startDate != nil || endDate != nil) {
		filtered := make([]*domain.Appointment, 0)
		for _, apt := range appointments {
			if startDate != nil && apt.AppointmentDateTime.Before(*startDate) {
				continue
			}
			if endDate != nil && apt.AppointmentDateTime.After(*endDate) {
				continue
			}
			filtered = append(filtered, apt)
		}
		appointments = filtered
	}

	return appointments, nil
}

// UpdateAppointmentRequest 予約更新リクエスト
type UpdateAppointmentRequest struct {
	AppointmentID      string
	TenantID           string
	FitterID           string
	AppointmentDateTime *time.Time
	DurationMinutes    int
	Status             domain.AppointmentStatus
	Notes              string
}

// UpdateAppointment 予約を更新
func (s *AppointmentService) UpdateAppointment(ctx context.Context, req *UpdateAppointmentRequest) (*domain.Appointment, error) {
	if req.AppointmentID == "" {
		return nil, fmt.Errorf("appointment_id is required")
	}
	if req.TenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}

	// 既存の予約を取得
	appointment, err := s.appointmentRepo.GetByID(ctx, req.AppointmentID, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	// 更新可能なフィールドを更新
	if req.FitterID != "" {
		appointment.FitterID = req.FitterID
	}
	if req.AppointmentDateTime != nil {
		// 時間変更時は空き状況チェック
		available, err := s.appointmentRepo.CheckAvailability(ctx, appointment.FitterID, req.TenantID, *req.AppointmentDateTime, appointment.DurationMinutes)
		if err != nil {
			return nil, fmt.Errorf("failed to check availability: %w", err)
		}
		if !available {
			return nil, fmt.Errorf("appointment time slot is not available")
		}
		appointment.AppointmentDateTime = *req.AppointmentDateTime
	}
	if req.DurationMinutes > 0 {
		appointment.DurationMinutes = req.DurationMinutes
	}
	if req.Status != "" {
		if !req.Status.IsValid() {
			return nil, fmt.Errorf("invalid status: %s", req.Status)
		}
		appointment.Status = req.Status
	}
	if req.Notes != "" {
		appointment.Notes = req.Notes
	}

	// リポジトリに保存
	if err := s.appointmentRepo.Update(ctx, appointment); err != nil {
		return nil, fmt.Errorf("failed to update appointment: %w", err)
	}

	return appointment, nil
}

// CancelAppointment 予約をキャンセル
func (s *AppointmentService) CancelAppointment(ctx context.Context, appointmentID string, tenantID string, reason string) error {
	if appointmentID == "" {
		return fmt.Errorf("appointment_id is required")
	}
	if tenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}

	// 既存の予約を取得
	appointment, err := s.appointmentRepo.GetByID(ctx, appointmentID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get appointment: %w", err)
	}

	// キャンセル可能かチェック
	if appointment.Status == domain.AppointmentStatusCancelled {
		return fmt.Errorf("appointment is already cancelled")
	}
	if appointment.Status == domain.AppointmentStatusCompleted {
		return fmt.Errorf("cannot cancel completed appointment")
	}

	// キャンセル実行
	if err := s.appointmentRepo.Cancel(ctx, appointmentID, tenantID, reason); err != nil {
		return fmt.Errorf("failed to cancel appointment: %w", err)
	}

	return nil
}

// CheckAvailabilityRequest 空き状況チェックリクエスト
type CheckAvailabilityRequest struct {
	FitterID          string
	TenantID          string
	AppointmentDateTime time.Time
	DurationMinutes   int
}

// CheckAvailability 空き状況をチェック
func (s *AppointmentService) CheckAvailability(ctx context.Context, req *CheckAvailabilityRequest) (bool, error) {
	if req.FitterID == "" {
		return false, fmt.Errorf("fitter_id is required")
	}
	if req.TenantID == "" {
		return false, fmt.Errorf("tenant_id is required")
	}
	if req.AppointmentDateTime.IsZero() {
		return false, fmt.Errorf("appointment_datetime is required")
	}
	if req.DurationMinutes <= 0 {
		return false, fmt.Errorf("duration_minutes must be greater than 0")
	}

	available, err := s.appointmentRepo.CheckAvailability(ctx, req.FitterID, req.TenantID, req.AppointmentDateTime, req.DurationMinutes)
	if err != nil {
		return false, fmt.Errorf("failed to check availability: %w", err)
	}

	return available, nil
}

