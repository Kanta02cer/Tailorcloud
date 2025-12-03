package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"tailor-cloud/backend/internal/config/domain"
)

// DiagnosisRepository 診断リポジトリインターフェース
type DiagnosisRepository interface {
	Create(ctx context.Context, diagnosis *domain.Diagnosis) error
	GetByID(ctx context.Context, diagnosisID string, tenantID string) (*domain.Diagnosis, error)
	GetByUserID(ctx context.Context, userID string, tenantID string) ([]*domain.Diagnosis, error)
	GetByTenantID(ctx context.Context, tenantID string, limit, offset int) ([]*domain.Diagnosis, error)
	List(ctx context.Context, tenantID string, filter DiagnosisFilter) ([]*domain.Diagnosis, error)
}

// DiagnosisFilter 診断フィルター
type DiagnosisFilter struct {
	Archetype *domain.Archetype
	PlanType  *domain.PlanType
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
}

// PostgreSQLDiagnosisRepository PostgreSQLを使った診断リポジトリ実装
type PostgreSQLDiagnosisRepository struct {
	db *sql.DB
}

// NewPostgreSQLDiagnosisRepository PostgreSQLDiagnosisRepositoryのコンストラクタ
func NewPostgreSQLDiagnosisRepository(db *sql.DB) DiagnosisRepository {
	return &PostgreSQLDiagnosisRepository{
		db: db,
	}
}

// Create 診断を作成
func (r *PostgreSQLDiagnosisRepository) Create(ctx context.Context, diagnosis *domain.Diagnosis) error {
	now := time.Now()
	if diagnosis.CreatedAt.IsZero() {
		diagnosis.CreatedAt = now
	}
	if diagnosis.UpdatedAt.IsZero() {
		diagnosis.UpdatedAt = now
	}

	query := `
		INSERT INTO diagnoses (
			id, user_id, tenant_id, archetype, plan_type, diagnosis_result, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		diagnosis.ID,
		diagnosis.UserID,
		diagnosis.TenantID,
		string(diagnosis.Archetype),
		string(diagnosis.PlanType),
		diagnosis.DiagnosisResult,
		diagnosis.CreatedAt,
		diagnosis.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create diagnosis: %w", err)
	}

	return nil
}

// GetByID 診断IDで取得（テナントIDもチェック）
func (r *PostgreSQLDiagnosisRepository) GetByID(ctx context.Context, diagnosisID string, tenantID string) (*domain.Diagnosis, error) {
	query := `
		SELECT 
			id, user_id, tenant_id, archetype, plan_type, diagnosis_result, created_at, updated_at
		FROM diagnoses
		WHERE id = $1 AND tenant_id = $2
	`

	var diagnosis domain.Diagnosis
	var archetypeStr, planTypeStr string
	var diagnosisResultBytes []byte

	err := r.db.QueryRowContext(ctx, query, diagnosisID, tenantID).Scan(
		&diagnosis.ID,
		&diagnosis.UserID,
		&diagnosis.TenantID,
		&archetypeStr,
		&planTypeStr,
		&diagnosisResultBytes,
		&diagnosis.CreatedAt,
		&diagnosis.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("diagnosis not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get diagnosis: %w", err)
	}

	diagnosis.Archetype = domain.Archetype(archetypeStr)
	diagnosis.PlanType = domain.PlanType(planTypeStr)
	diagnosis.DiagnosisResult = json.RawMessage(diagnosisResultBytes)

	return &diagnosis, nil
}

// GetByUserID ユーザーIDで診断一覧を取得
func (r *PostgreSQLDiagnosisRepository) GetByUserID(ctx context.Context, userID string, tenantID string) ([]*domain.Diagnosis, error) {
	query := `
		SELECT 
			id, user_id, tenant_id, archetype, plan_type, diagnosis_result, created_at, updated_at
		FROM diagnoses
		WHERE user_id = $1 AND tenant_id = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get diagnoses by user ID: %w", err)
	}
	defer rows.Close()

	var diagnoses []*domain.Diagnosis
	for rows.Next() {
		var diagnosis domain.Diagnosis
		var archetypeStr, planTypeStr string
		var diagnosisResultBytes []byte

		err := rows.Scan(
			&diagnosis.ID,
			&diagnosis.UserID,
			&diagnosis.TenantID,
			&archetypeStr,
			&planTypeStr,
			&diagnosisResultBytes,
			&diagnosis.CreatedAt,
			&diagnosis.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan diagnosis: %w", err)
		}

		diagnosis.Archetype = domain.Archetype(archetypeStr)
		diagnosis.PlanType = domain.PlanType(planTypeStr)
		diagnosis.DiagnosisResult = json.RawMessage(diagnosisResultBytes)

		diagnoses = append(diagnoses, &diagnosis)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate diagnoses: %w", err)
	}

	return diagnoses, nil
}

// GetByTenantID テナントIDで診断一覧を取得（ページネーション対応）
func (r *PostgreSQLDiagnosisRepository) GetByTenantID(ctx context.Context, tenantID string, limit, offset int) ([]*domain.Diagnosis, error) {
	query := `
		SELECT 
			id, user_id, tenant_id, archetype, plan_type, diagnosis_result, created_at, updated_at
		FROM diagnoses
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get diagnoses by tenant ID: %w", err)
	}
	defer rows.Close()

	var diagnoses []*domain.Diagnosis
	for rows.Next() {
		var diagnosis domain.Diagnosis
		var archetypeStr, planTypeStr string
		var diagnosisResultBytes []byte

		err := rows.Scan(
			&diagnosis.ID,
			&diagnosis.UserID,
			&diagnosis.TenantID,
			&archetypeStr,
			&planTypeStr,
			&diagnosisResultBytes,
			&diagnosis.CreatedAt,
			&diagnosis.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan diagnosis: %w", err)
		}

		diagnosis.Archetype = domain.Archetype(archetypeStr)
		diagnosis.PlanType = domain.PlanType(planTypeStr)
		diagnosis.DiagnosisResult = json.RawMessage(diagnosisResultBytes)

		diagnoses = append(diagnoses, &diagnosis)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate diagnoses: %w", err)
	}

	return diagnoses, nil
}

// List 診断一覧を取得（フィルター対応）
func (r *PostgreSQLDiagnosisRepository) List(ctx context.Context, tenantID string, filter DiagnosisFilter) ([]*domain.Diagnosis, error) {
	query := `
		SELECT 
			id, user_id, tenant_id, archetype, plan_type, diagnosis_result, created_at, updated_at
		FROM diagnoses
		WHERE tenant_id = $1
	`
	args := []interface{}{tenantID}
	argIndex := 2

	// フィルター条件を追加
	if filter.Archetype != nil {
		query += fmt.Sprintf(" AND archetype = $%d", argIndex)
		args = append(args, string(*filter.Archetype))
		argIndex++
	}

	if filter.PlanType != nil {
		query += fmt.Sprintf(" AND plan_type = $%d", argIndex)
		args = append(args, string(*filter.PlanType))
		argIndex++
	}

	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, *filter.EndDate)
		argIndex++
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
		argIndex++
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list diagnoses: %w", err)
	}
	defer rows.Close()

	var diagnoses []*domain.Diagnosis
	for rows.Next() {
		var diagnosis domain.Diagnosis
		var archetypeStr, planTypeStr string
		var diagnosisResultBytes []byte

		err := rows.Scan(
			&diagnosis.ID,
			&diagnosis.UserID,
			&diagnosis.TenantID,
			&archetypeStr,
			&planTypeStr,
			&diagnosisResultBytes,
			&diagnosis.CreatedAt,
			&diagnosis.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan diagnosis: %w", err)
		}

		diagnosis.Archetype = domain.Archetype(archetypeStr)
		diagnosis.PlanType = domain.PlanType(planTypeStr)
		diagnosis.DiagnosisResult = json.RawMessage(diagnosisResultBytes)

		diagnoses = append(diagnoses, &diagnosis)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate diagnoses: %w", err)
	}

	return diagnoses, nil
}

