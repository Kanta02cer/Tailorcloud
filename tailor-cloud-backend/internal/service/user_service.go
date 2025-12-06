package service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/repository"
)

// UserService ユーザーサービス
type UserService struct {
	userRepo   repository.UserRepository
	tenantRepo repository.TenantRepository
}

// NewUserService UserServiceのコンストラクタ
func NewUserService(userRepo repository.UserRepository, tenantRepo repository.TenantRepository) *UserService {
	return &UserService{
		userRepo:   userRepo,
		tenantRepo: tenantRepo,
	}
}

// GetOrCreateUserRequest ユーザー取得または作成リクエスト
type GetOrCreateUserRequest struct {
	FirebaseUID string
	Email       string
	Name        string
	TenantID    string // オプショナル: 指定がない場合はデフォルトテナントを使用
}

// GetOrCreateUser ユーザーを取得または作成
// Firebase UIDで検索し、存在しない場合は新規作成
func (s *UserService) GetOrCreateUser(ctx context.Context, req *GetOrCreateUserRequest) (*domain.User, error) {
	// まずFirebase UIDで検索
	user, err := s.userRepo.GetByFirebaseUID(ctx, req.FirebaseUID)
	if err == nil {
		// ユーザーが存在する場合は最終ログイン時刻を更新
		_ = s.userRepo.UpdateLastLogin(ctx, user.ID)
		return user, nil
	}

	// ユーザーが存在しない場合は新規作成
	tenantID := req.TenantID
	if tenantID == "" {
		// 環境変数からデフォルトテナントIDを取得
		tenantID = os.Getenv("DEFAULT_TENANT_ID")
		if tenantID == "" {
			// 環境変数が設定されていない場合は、メールドメインから推測（将来の拡張用）
			// TODO: テナントリポジトリにGetByDomainメソッドを実装して、ドメインからテナントを検索
			tenantID = "default-tenant-id" // フォールバック値
		}
	}

	// テナントの存在確認
	_, err = s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		// テナントが存在しない場合はエラー
		// ログに記録して、より詳細なエラーメッセージを返す
		return nil, fmt.Errorf("tenant not found (tenant_id: %s). Please ensure the tenant exists or set DEFAULT_TENANT_ID environment variable: %w", tenantID, err)
	}

	// ユーザー名が空の場合はメールアドレスから生成
	name := req.Name
	if name == "" {
		emailParts := strings.Split(req.Email, "@")
		name = emailParts[0]
	}

	// 新規ユーザーを作成
	newUser := &domain.User{
		ID:        "", // UUIDはリポジトリで生成
		TenantID:  tenantID,
		Email:     req.Email,
		Name:      name,
		Role:      domain.RoleStaff, // デフォルトロール: Staff
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, newUser, req.FirebaseUID); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Firebase UIDを更新（別途実装が必要）
	// ここでは一旦作成したユーザーを返す
	return newUser, nil
}

// GetUserByFirebaseUID Firebase UIDでユーザーを取得
func (s *UserService) GetUserByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	user, err := s.userRepo.GetByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 最終ログイン時刻を更新
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	return user, nil
}
