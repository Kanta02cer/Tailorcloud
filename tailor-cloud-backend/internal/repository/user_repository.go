package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tailor-cloud/backend/internal/config/domain"

	"github.com/google/uuid"
)

// UserRepository ユーザーリポジトリインターフェース
type UserRepository interface {
	GetByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User, firebaseUID string) error
	UpdateLastLogin(ctx context.Context, userID string) error
}

// PostgreSQLUserRepository PostgreSQLを使ったユーザーリポジトリ実装
type PostgreSQLUserRepository struct {
	db *sql.DB
}

// NewPostgreSQLUserRepository PostgreSQLUserRepositoryのコンストラクタ
func NewPostgreSQLUserRepository(db *sql.DB) UserRepository {
	return &PostgreSQLUserRepository{
		db: db,
	}
}

// GetByFirebaseUID Firebase UIDでユーザーを取得
func (r *PostgreSQLUserRepository) GetByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	query := `
		SELECT 
			id, tenant_id, email, name, role, created_at, updated_at
		FROM users
		WHERE firebase_uid = $1 AND is_active = TRUE
	`

	var user domain.User
	var roleStr string
	err := r.db.QueryRowContext(ctx, query, firebaseUID).Scan(
		&user.ID,
		&user.TenantID,
		&user.Email,
		&user.Name,
		&roleStr,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by firebase_uid: %w", err)
	}

	user.Role = domain.UserRole(roleStr)
	return &user, nil
}

// GetByEmail メールアドレスでユーザーを取得
func (r *PostgreSQLUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT 
			id, tenant_id, email, name, role, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = TRUE
		ORDER BY created_at DESC
		LIMIT 1
	`

	var user domain.User
	var roleStr string
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.TenantID,
		&user.Email,
		&user.Name,
		&roleStr,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	user.Role = domain.UserRole(roleStr)
	return &user, nil
}

// Create ユーザーを作成
func (r *PostgreSQLUserRepository) Create(ctx context.Context, user *domain.User, firebaseUID string) error {
	query := `
		INSERT INTO users (id, tenant_id, email, name, role, firebase_uid, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	userID := user.ID
	if userID == "" {
		userID = uuid.New().String()
		user.ID = userID
	}

	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}

	// Firebase UIDを保存
	var firebaseUIDNull sql.NullString
	if firebaseUID != "" {
		firebaseUIDNull = sql.NullString{String: firebaseUID, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query,
		userID,
		user.TenantID,
		user.Email,
		user.Name,
		string(user.Role),
		firebaseUIDNull,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// UpdateLastLogin 最終ログイン時刻を更新
func (r *PostgreSQLUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET last_login_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

