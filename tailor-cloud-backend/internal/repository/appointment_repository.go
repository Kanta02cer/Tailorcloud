package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"tailor-cloud/backend/internal/config/domain"
)

// AppointmentRepository 予約リポジトリインターフェース
type AppointmentRepository interface {
	Create(ctx context.Context, appointment *domain.Appointment) error
	GetByID(ctx context.Context, appointmentID string, tenantID string) (*domain.Appointment, error)
	GetByUserID(ctx context.Context, userID string, tenantID string) ([]*domain.Appointment, error)
	GetByFitterID(ctx context.Context, fitterID string, tenantID string, startDate, endDate time.Time) ([]*domain.Appointment, error)
	GetByTenantID(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]*domain.Appointment, error)
	Update(ctx context.Context, appointment *domain.Appointment) error
	Cancel(ctx context.Context, appointmentID string, tenantID string, reason string) error
	CheckAvailability(ctx context.Context, fitterID string, tenantID string, appointmentDateTime time.Time, durationMinutes int) (bool, error)
}

// PostgreSQLAppointmentRepository PostgreSQLを使った予約リポジトリ実装
type PostgreSQLAppointmentRepository struct {
	db *sql.DB
}

// NewPostgreSQLAppointmentRepository PostgreSQLAppointmentRepositoryのコンストラクタ
func NewPostgreSQLAppointmentRepository(db *sql.DB) AppointmentRepository {
	return &PostgreSQLAppointmentRepository{
		db: db,
	}
}

// Create 予約を作成
func (r *PostgreSQLAppointmentRepository) Create(ctx context.Context, appointment *domain.Appointment) error {
	now := time.Now()
	if appointment.CreatedAt.IsZero() {
		appointment.CreatedAt = now
	}
	if appointment.UpdatedAt.IsZero() {
		appointment.UpdatedAt = now
	}

	var depositAmount sql.NullInt64
	if appointment.DepositAmount != nil {
		depositAmount = sql.NullInt64{Int64: *appointment.DepositAmount, Valid: true}
	}

	var depositStatus sql.NullString
	if appointment.DepositStatus != "" {
		depositStatus = sql.NullString{String: string(appointment.DepositStatus), Valid: true}
	}

	query := `
		INSERT INTO appointments (
			id, user_id, tenant_id, fitter_id, appointment_datetime, duration_minutes,
			status, deposit_amount, deposit_payment_intent_id, deposit_status, notes,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.ExecContext(ctx, query,
		appointment.ID,
		appointment.UserID,
		appointment.TenantID,
		appointment.FitterID,
		appointment.AppointmentDateTime,
		appointment.DurationMinutes,
		string(appointment.Status),
		depositAmount,
		appointment.DepositPaymentIntentID,
		depositStatus,
		appointment.Notes,
		appointment.CreatedAt,
		appointment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create appointment: %w", err)
	}

	return nil
}

// GetByID 予約IDで取得（テナントIDもチェック）
func (r *PostgreSQLAppointmentRepository) GetByID(ctx context.Context, appointmentID string, tenantID string) (*domain.Appointment, error) {
	query := `
		SELECT 
			id, user_id, tenant_id, fitter_id, appointment_datetime, duration_minutes,
			status, deposit_amount, deposit_payment_intent_id, deposit_status, notes,
			cancelled_at, cancelled_reason, created_at, updated_at
		FROM appointments
		WHERE id = $1 AND tenant_id = $2
	`

	var appointment domain.Appointment
	var statusStr, depositStatusStr, depositPaymentIntentID, notes, cancelledReason sql.NullString
	var depositAmount sql.NullInt64
	var cancelledAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, appointmentID, tenantID).Scan(
		&appointment.ID,
		&appointment.UserID,
		&appointment.TenantID,
		&appointment.FitterID,
		&appointment.AppointmentDateTime,
		&appointment.DurationMinutes,
		&statusStr,
		&depositAmount,
		&depositPaymentIntentID,
		&depositStatusStr,
		&notes,
		&cancelledAt,
		&cancelledReason,
		&appointment.CreatedAt,
		&appointment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("appointment not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	appointment.Status = domain.AppointmentStatus(statusStr.String)
	if depositAmount.Valid {
		deposit := depositAmount.Int64
		appointment.DepositAmount = &deposit
	}
	if depositPaymentIntentID.Valid {
		appointment.DepositPaymentIntentID = depositPaymentIntentID.String
	}
	if depositStatusStr.Valid {
		appointment.DepositStatus = domain.DepositStatus(depositStatusStr.String)
	}
	if notes.Valid {
		appointment.Notes = notes.String
	}
	if cancelledAt.Valid {
		appointment.CancelledAt = &cancelledAt.Time
	}
	if cancelledReason.Valid {
		appointment.CancelledReason = cancelledReason.String
	}

	return &appointment, nil
}

// GetByUserID ユーザーIDで予約一覧を取得
func (r *PostgreSQLAppointmentRepository) GetByUserID(ctx context.Context, userID string, tenantID string) ([]*domain.Appointment, error) {
	query := `
		SELECT 
			id, user_id, tenant_id, fitter_id, appointment_datetime, duration_minutes,
			status, deposit_amount, deposit_payment_intent_id, deposit_status, notes,
			cancelled_at, cancelled_reason, created_at, updated_at
		FROM appointments
		WHERE user_id = $1 AND tenant_id = $2
		ORDER BY appointment_datetime DESC
	`

	return r.scanAppointments(ctx, query, userID, tenantID)
}

// GetByFitterID フィッターIDで予約一覧を取得（期間指定）
func (r *PostgreSQLAppointmentRepository) GetByFitterID(ctx context.Context, fitterID string, tenantID string, startDate, endDate time.Time) ([]*domain.Appointment, error) {
	query := `
		SELECT 
			id, user_id, tenant_id, fitter_id, appointment_datetime, duration_minutes,
			status, deposit_amount, deposit_payment_intent_id, deposit_status, notes,
			cancelled_at, cancelled_reason, created_at, updated_at
		FROM appointments
		WHERE fitter_id = $1 AND tenant_id = $2
		  AND appointment_datetime >= $3
		  AND appointment_datetime <= $4
		  AND status != 'Cancelled'
		ORDER BY appointment_datetime ASC
	`

	return r.scanAppointments(ctx, query, fitterID, tenantID, startDate, endDate)
}

// GetByTenantID テナントIDで予約一覧を取得（期間指定）
func (r *PostgreSQLAppointmentRepository) GetByTenantID(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]*domain.Appointment, error) {
	query := `
		SELECT 
			id, user_id, tenant_id, fitter_id, appointment_datetime, duration_minutes,
			status, deposit_amount, deposit_payment_intent_id, deposit_status, notes,
			cancelled_at, cancelled_reason, created_at, updated_at
		FROM appointments
		WHERE tenant_id = $1
		  AND appointment_datetime >= $2
		  AND appointment_datetime <= $3
		ORDER BY appointment_datetime ASC
	`

	return r.scanAppointments(ctx, query, tenantID, startDate, endDate)
}

// Update 予約を更新
func (r *PostgreSQLAppointmentRepository) Update(ctx context.Context, appointment *domain.Appointment) error {
	appointment.UpdatedAt = time.Now()

	var depositAmount sql.NullInt64
	if appointment.DepositAmount != nil {
		depositAmount = sql.NullInt64{Int64: *appointment.DepositAmount, Valid: true}
	}

	var depositStatus sql.NullString
	if appointment.DepositStatus != "" {
		depositStatus = sql.NullString{String: string(appointment.DepositStatus), Valid: true}
	}

	query := `
		UPDATE appointments
		SET fitter_id = $3, appointment_datetime = $4, duration_minutes = $5,
		    status = $6, deposit_amount = $7, deposit_payment_intent_id = $8,
		    deposit_status = $9, notes = $10, updated_at = $11
		WHERE id = $1 AND tenant_id = $2
	`

	result, err := r.db.ExecContext(ctx, query,
		appointment.ID,
		appointment.TenantID,
		appointment.FitterID,
		appointment.AppointmentDateTime,
		appointment.DurationMinutes,
		string(appointment.Status),
		depositAmount,
		appointment.DepositPaymentIntentID,
		depositStatus,
		appointment.Notes,
		appointment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update appointment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("appointment not found or tenant_id mismatch")
	}

	return nil
}

// Cancel 予約をキャンセル
func (r *PostgreSQLAppointmentRepository) Cancel(ctx context.Context, appointmentID string, tenantID string, reason string) error {
	now := time.Now()

	query := `
		UPDATE appointments
		SET status = 'Cancelled', cancelled_at = $3, cancelled_reason = $4, updated_at = $3
		WHERE id = $1 AND tenant_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, appointmentID, tenantID, now, reason)
	if err != nil {
		return fmt.Errorf("failed to cancel appointment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("appointment not found or tenant_id mismatch")
	}

	return nil
}

// CheckAvailability 予約可能かチェック（時間重複チェック）
func (r *PostgreSQLAppointmentRepository) CheckAvailability(ctx context.Context, fitterID string, tenantID string, appointmentDateTime time.Time, durationMinutes int) (bool, error) {
	endDateTime := appointmentDateTime.Add(time.Duration(durationMinutes) * time.Minute)

	query := `
		SELECT COUNT(*) 
		FROM appointments
		WHERE fitter_id = $1 
		  AND tenant_id = $2
		  AND status != 'Cancelled'
		  AND (
			(appointment_datetime <= $3 AND appointment_datetime + (duration_minutes || ' minutes')::interval > $3) OR
			(appointment_datetime < $4 AND appointment_datetime + (duration_minutes || ' minutes')::interval >= $4) OR
			(appointment_datetime >= $3 AND appointment_datetime + (duration_minutes || ' minutes')::interval <= $4)
		  )
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, fitterID, tenantID, appointmentDateTime, endDateTime).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check availability: %w", err)
	}

	return count == 0, nil
}

// scanAppointments 予約一覧をスキャン（ヘルパー関数）
func (r *PostgreSQLAppointmentRepository) scanAppointments(ctx context.Context, query string, args ...interface{}) ([]*domain.Appointment, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query appointments: %w", err)
	}
	defer rows.Close()

	var appointments []*domain.Appointment
	for rows.Next() {
		var appointment domain.Appointment
		var statusStr, depositStatusStr, depositPaymentIntentID, notes, cancelledReason sql.NullString
		var depositAmount sql.NullInt64
		var cancelledAt sql.NullTime

		err := rows.Scan(
			&appointment.ID,
			&appointment.UserID,
			&appointment.TenantID,
			&appointment.FitterID,
			&appointment.AppointmentDateTime,
			&appointment.DurationMinutes,
			&statusStr,
			&depositAmount,
			&depositPaymentIntentID,
			&depositStatusStr,
			&notes,
			&cancelledAt,
			&cancelledReason,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan appointment: %w", err)
		}

		appointment.Status = domain.AppointmentStatus(statusStr.String)
		if depositAmount.Valid {
			deposit := depositAmount.Int64
			appointment.DepositAmount = &deposit
		}
		if depositPaymentIntentID.Valid {
			appointment.DepositPaymentIntentID = depositPaymentIntentID.String
		}
		if depositStatusStr.Valid {
			appointment.DepositStatus = domain.DepositStatus(depositStatusStr.String)
		}
		if notes.Valid {
			appointment.Notes = notes.String
		}
		if cancelledAt.Valid {
			appointment.CancelledAt = &cancelledAt.Time
		}
		if cancelledReason.Valid {
			appointment.CancelledReason = cancelledReason.String
		}

		appointments = append(appointments, &appointment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate appointments: %w", err)
	}

	return appointments, nil
}

