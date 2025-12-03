package service

import (
	"context"
	"fmt"

	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/repository"
)

// RBACService ロールベースアクセス制御サービス
// リソースベースの細かい権限管理を実現
type RBACService struct {
	permissionRepo repository.PermissionRepository
	cacheEnabled   bool // 権限キャッシュ（将来的に実装）
}

// NewRBACService RBACServiceのコンストラクタ
func NewRBACService(permissionRepo repository.PermissionRepository) *RBACService {
	return &RBACService{
		permissionRepo: permissionRepo,
		cacheEnabled:   false, // 将来実装
	}
}

// CheckPermission 権限をチェック
// tenantID: テナントID
// resourceType: リソースタイプ
// action: 操作
// role: ユーザーロール
// resourceID: 特定リソースID（オプション）
func (s *RBACService) CheckPermission(
	ctx context.Context,
	tenantID string,
	resourceType domain.ResourceType,
	action domain.Action,
	role domain.UserRole,
	resourceID *string,
) (bool, error) {
	// Ownerロールは常に全権限を持つ
	if role == domain.RoleOwner {
		return true, nil
	}
	
	// 権限を取得
	permission, err := s.permissionRepo.GetPermission(ctx, tenantID, resourceType, action, role, resourceID)
	if err != nil {
		// 権限が見つからない場合は拒否
		return false, fmt.Errorf("permission denied: %w", err)
	}
	
	// 許可されているかチェック
	return permission.Granted, nil
}

// CheckPermissionWithReason 権限をチェック（理由付き）
func (s *RBACService) CheckPermissionWithReason(
	ctx context.Context,
	tenantID string,
	resourceType domain.ResourceType,
	action domain.Action,
	role domain.UserRole,
	resourceID *string,
) (*domain.PermissionCheck, error) {
	// Ownerロールは常に全権限を持つ
	if role == domain.RoleOwner {
		return &domain.PermissionCheck{
			Allowed:  true,
			Reason:   "Owner role has all permissions",
			Permission: nil,
		}, nil
	}
	
	// 権限を取得
	permission, err := s.permissionRepo.GetPermission(ctx, tenantID, resourceType, action, role, resourceID)
	if err != nil {
		return &domain.PermissionCheck{
			Allowed:  false,
			Reason:   fmt.Sprintf("Permission not found: %v", err),
			Permission: nil,
		}, nil
	}
	
	return &domain.PermissionCheck{
		Allowed:    permission.Granted,
		Reason:     fmt.Sprintf("Permission %s on %s: %v", action, resourceType, permission.Granted),
		Permission: permission,
	}, nil
}

// GetUserPermissions ユーザーの権限一覧を取得
func (s *RBACService) GetUserPermissions(ctx context.Context, tenantID string, role domain.UserRole) ([]*domain.Permission, error) {
	return s.permissionRepo.GetPermissionsByRole(ctx, tenantID, role)
}

// GrantPermission 権限を付与
func (s *RBACService) GrantPermission(ctx context.Context, permission *domain.Permission) error {
	return s.permissionRepo.CreatePermission(ctx, permission)
}

// RevokePermission 権限を剥奪
func (s *RBACService) RevokePermission(ctx context.Context, permissionID string, tenantID string) error {
	return s.permissionRepo.DeletePermission(ctx, permissionID, tenantID)
}

// UpdatePermission 権限を更新
func (s *RBACService) UpdatePermission(ctx context.Context, permission *domain.Permission) error {
	return s.permissionRepo.UpdatePermission(ctx, permission)
}

