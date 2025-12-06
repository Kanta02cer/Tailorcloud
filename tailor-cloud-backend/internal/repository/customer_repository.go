package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"tailor-cloud/backend/internal/config/domain"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// CustomerRepository 顧客リポジトリインターフェース
type CustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	GetByID(ctx context.Context, customerID string, tenantID string) (*domain.Customer, error)
	GetByTenantID(ctx context.Context, tenantID string) ([]*domain.Customer, error)
	Search(ctx context.Context, tenantID string, keyword string) ([]*domain.Customer, error)
	Update(ctx context.Context, customer *domain.Customer) error
	Delete(ctx context.Context, customerID string, tenantID string) error
}

// PostgreSQLCustomerRepository PostgreSQLを使った顧客リポジトリ実装
type PostgreSQLCustomerRepository struct {
	db *sql.DB
}

type rowScanner interface {
	Scan(dest ...any) error
}

// NewPostgreSQLCustomerRepository PostgreSQLCustomerRepositoryのコンストラクタ
func NewPostgreSQLCustomerRepository(db *sql.DB) CustomerRepository {
	return &PostgreSQLCustomerRepository{
		db: db,
	}
}

func scanCustomerFromRow(scanner rowScanner) (*domain.Customer, error) {
	var (
		customer         domain.Customer
		email            sql.NullString
		phone            sql.NullString
		address          sql.NullString
		status           sql.NullString
		tags             pq.StringArray
		vipRank          sql.NullInt32
		ltvScore         sql.NullFloat64
		lifetimeValue    sql.NullFloat64
		lastInteraction  sql.NullTime
		interactionRaw   []byte
		preferredChannel sql.NullString
		leadSource       sql.NullString
		notes            sql.NullString
		occupation       sql.NullString
		incomeRange      sql.NullString
		archetype        sql.NullString
		diagnosisCount   sql.NullInt32
	)

	err := scanner.Scan(
		&customer.ID,
		&customer.TenantID,
		&customer.Name,
		&email,
		&phone,
		&address,
		&status,
		pq.Array(&tags),
		&vipRank,
		&ltvScore,
		&lifetimeValue,
		&lastInteraction,
		&interactionRaw,
		&preferredChannel,
		&leadSource,
		&notes,
		&occupation,
		&incomeRange,
		&archetype,
		&diagnosisCount,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if email.Valid {
		customer.Email = email.String
	}
	if phone.Valid {
		customer.Phone = phone.String
	}
	if address.Valid {
		customer.Address = address.String
	}
	if status.Valid {
		customer.CustomerStatus = status.String
	}
	if tags != nil {
		customer.Tags = []string(tags)
	}
	if vipRank.Valid {
		customer.VIPRank = int(vipRank.Int32)
	}
	if ltvScore.Valid {
		customer.LTVScore = ltvScore.Float64
	}
	if lifetimeValue.Valid {
		customer.LifetimeValue = lifetimeValue.Float64
	}
	if lastInteraction.Valid {
		ts := lastInteraction.Time
		customer.LastInteractionAt = &ts
	}
	if len(interactionRaw) > 0 {
		if err := json.Unmarshal(interactionRaw, &customer.InteractionNotes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal interaction notes: %w", err)
		}
	}
	if preferredChannel.Valid {
		customer.PreferredChannel = preferredChannel.String
	}
	if leadSource.Valid {
		customer.LeadSource = leadSource.String
	}
	if notes.Valid {
		customer.Notes = notes.String
	}
	if occupation.Valid {
		customer.Occupation = occupation.String
	}
	if incomeRange.Valid {
		customer.AnnualIncomeRange = incomeRange.String
	}
	if archetype.Valid {
		customer.PreferredArchetype = archetype.String
	}
	if diagnosisCount.Valid {
		customer.DiagnosisCount = int(diagnosisCount.Int32)
	}

	return &customer, nil
}

// Create 顧客を作成
func (r *PostgreSQLCustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	if customer.ID == "" {
		customer.ID = uuid.New().String()
	}

	if customer.CustomerStatus == "" {
		customer.CustomerStatus = "lead"
	}
	if customer.Tags == nil {
		customer.Tags = []string{}
	}

	now := time.Now()
	if customer.CreatedAt.IsZero() {
		customer.CreatedAt = now
	}
	if customer.UpdatedAt.IsZero() {
		customer.UpdatedAt = now
	}

	var (
		lastInteraction interface{}
		interactionJSON []byte
		err             error
	)
	if customer.LastInteractionAt != nil && !customer.LastInteractionAt.IsZero() {
		lastInteraction = customer.LastInteractionAt
	}
	if len(customer.InteractionNotes) > 0 {
		interactionJSON, err = json.Marshal(customer.InteractionNotes)
		if err != nil {
			return fmt.Errorf("failed to marshal interaction notes: %w", err)
		}
	} else {
		interactionJSON = []byte("[]")
	}

	query := `
		INSERT INTO customers (
			id, tenant_id, name, email, phone, address, customer_status, tags, vip_rank,
			ltv_score, lifetime_value, last_interaction_at, interaction_notes,
			preferred_channel, lead_source, notes, occupation, annual_income_range,
			preferred_archetype, diagnosis_count, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13,
			$14, $15, $16, $17, $18,
			$19, $20, $21, $22
		)
	`

	_, err = r.db.ExecContext(ctx, query,
		customer.ID,
		customer.TenantID,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.Address,
		customer.CustomerStatus,
		pq.Array(customer.Tags),
		customer.VIPRank,
		customer.LTVScore,
		customer.LifetimeValue,
		lastInteraction,
		interactionJSON,
		customer.PreferredChannel,
		customer.LeadSource,
		customer.Notes,
		customer.Occupation,
		customer.AnnualIncomeRange,
		customer.PreferredArchetype,
		customer.DiagnosisCount,
		customer.CreatedAt,
		customer.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}

	return nil
}

// GetByID 顧客IDで取得（テナントIDもチェック）
func (r *PostgreSQLCustomerRepository) GetByID(ctx context.Context, customerID string, tenantID string) (*domain.Customer, error) {
	query := `
		SELECT 
			id, tenant_id, name, email, phone, address, customer_status, tags, vip_rank,
			ltv_score, lifetime_value, last_interaction_at, interaction_notes,
			preferred_channel, lead_source, notes, occupation, annual_income_range,
			preferred_archetype, diagnosis_count, created_at, updated_at
		FROM customers
		WHERE id = $1 AND tenant_id = $2
	`

	row := r.db.QueryRowContext(ctx, query, customerID, tenantID)
	customer, err := scanCustomerFromRow(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("customer not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return customer, nil
}

// GetByTenantID テナントIDで顧客一覧を取得
func (r *PostgreSQLCustomerRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*domain.Customer, error) {
	query := `
		SELECT 
			id, tenant_id, name, email, phone, address, customer_status, tags, vip_rank,
			ltv_score, lifetime_value, last_interaction_at, interaction_notes,
			preferred_channel, lead_source, notes, occupation, annual_income_range,
			preferred_archetype, diagnosis_count, created_at, updated_at
		FROM customers
		WHERE tenant_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}
	defer rows.Close()

	var customers []*domain.Customer
	for rows.Next() {
		customer, err := scanCustomerFromRow(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}

		customers = append(customers, customer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate customers: %w", err)
	}

	return customers, nil
}

// Search 顧客を検索（名前、メール、電話番号で検索）
func (r *PostgreSQLCustomerRepository) Search(ctx context.Context, tenantID string, keyword string) ([]*domain.Customer, error) {
	query := `
		SELECT 
			id, tenant_id, name, email, phone, address, customer_status, tags, vip_rank,
			ltv_score, lifetime_value, last_interaction_at, interaction_notes,
			preferred_channel, lead_source, notes, occupation, annual_income_range,
			preferred_archetype, diagnosis_count, created_at, updated_at
		FROM customers
		WHERE tenant_id = $1
		  AND (
			name ILIKE $2 OR
			email ILIKE $2 OR
			phone ILIKE $2
			OR $2 = ANY(tags)
		  )
		ORDER BY created_at DESC
	`

	searchPattern := "%" + keyword + "%"
	rows, err := r.db.QueryContext(ctx, query, tenantID, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search customers: %w", err)
	}
	defer rows.Close()

	var customers []*domain.Customer
	for rows.Next() {
		customer, err := scanCustomerFromRow(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}

		customers = append(customers, customer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate customers: %w", err)
	}

	return customers, nil
}

// Update 顧客を更新
func (r *PostgreSQLCustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	customer.UpdatedAt = time.Now()
	if customer.CustomerStatus == "" {
		customer.CustomerStatus = "lead"
	}
	if customer.Tags == nil {
		customer.Tags = []string{}
	}

	var (
		lastInteraction interface{}
		interactionJSON []byte
		err             error
	)
	if customer.LastInteractionAt != nil && !customer.LastInteractionAt.IsZero() {
		lastInteraction = customer.LastInteractionAt
	}
	if len(customer.InteractionNotes) > 0 {
		interactionJSON, err = json.Marshal(customer.InteractionNotes)
		if err != nil {
			return fmt.Errorf("failed to marshal interaction notes: %w", err)
		}
	} else {
		interactionJSON = []byte("[]")
	}

	query := `
		UPDATE customers
		SET name = $3,
		    email = $4,
		    phone = $5,
		    address = $6,
		    customer_status = $7,
		    tags = $8,
		    vip_rank = $9,
		    ltv_score = $10,
		    lifetime_value = $11,
		    last_interaction_at = $12,
		    interaction_notes = $13,
		    preferred_channel = $14,
		    lead_source = $15,
		    notes = $16,
		    occupation = $17,
		    annual_income_range = $18,
		    preferred_archetype = $19,
		    diagnosis_count = $20,
		    updated_at = $21
		WHERE id = $1 AND tenant_id = $2
	`

	result, err := r.db.ExecContext(ctx, query,
		customer.ID,
		customer.TenantID,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.Address,
		customer.CustomerStatus,
		pq.Array(customer.Tags),
		customer.VIPRank,
		customer.LTVScore,
		customer.LifetimeValue,
		lastInteraction,
		interactionJSON,
		customer.PreferredChannel,
		customer.LeadSource,
		customer.Notes,
		customer.Occupation,
		customer.AnnualIncomeRange,
		customer.PreferredArchetype,
		customer.DiagnosisCount,
		customer.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("customer not found or tenant_id mismatch")
	}

	return nil
}

// Delete 顧客を削除
func (r *PostgreSQLCustomerRepository) Delete(ctx context.Context, customerID string, tenantID string) error {
	query := `
		DELETE FROM customers
		WHERE id = $1 AND tenant_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, customerID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("customer not found or tenant_id mismatch")
	}

	return nil
}
