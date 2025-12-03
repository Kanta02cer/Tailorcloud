package repository

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"tailor-cloud/backend/internal/config/domain"
)

// PermissionRepository 権限リポジトリインターフェース
type PermissionRepository interface {
	GetPermission(ctx context.Context, tenantID string, resourceType domain.ResourceType, action domain.Action, role domain.UserRole, resourceID *string) (*domain.Permission, error)
	GetPermissionsByRole(ctx context.Context, tenantID string, role domain.UserRole) ([]*domain.Permission, error)
	GetPermissionsByResource(ctx context.Context, tenantID string, resourceType domain.ResourceType, resourceID *string) ([]*domain.Permission, error)
	CreatePermission(ctx context.Context, permission *domain.Permission) error
	UpdatePermission(ctx context.Context, permission *domain.Permission) error
	DeletePermission(ctx context.Context, permissionID string, tenantID string) error
}

// PostgreSQLPermissionRepository PostgreSQL実装
type PostgreSQLPermissionRepository struct {
	db *sql.DB
}

// NewPostgreSQLPermissionRepository PostgreSQLPermissionRepositoryのコンストラクタ
func NewPostgreSQLPermissionRepository(db *sql.DB) PermissionRepository {
	return &PostgreSQLPermissionRepository{
		db: db,
	}
}

// GetPermission 権限を取得（最も具体的な権限を優先）
func (r *PostgreSQLPermissionRepository) GetPermission(
	ctx context.Context,
	tenantID string,
	resourceType domain.ResourceType,
	action domain.Action,
	role domain.UserRole,
	resourceID *string,
) (*domain.Permission, error) {
	// 1. 特定リソースに対する権限を探す（最優先）
	// 2. リソースタイプに対する権限を探す
	// 3. 全リソースに対する権限を探す（最後）
	
	var permission domain.Permission
	var resourceTypeStr, actionStr, roleStr string
	var resourceIDPtr sql.NullString
	var granted bool
	
	query := `
		SELECT id, tenant_id, resource_type, resource_id, action, role, granted
		FROM permissions
		WHERE tenant_id = $1 AND role = $2
		AND (
			(resource_type = $3 AND action = $4 AND resource_id = $5) OR
			(resource_type = $3 AND action = $4 AND resource_id IS NULL) OR
			(resource_type = 'ALL' AND action = 'ALL') OR
			(resource_type = $3 AND action = 'ALL')
		)
		ORDER BY 
			CASE WHEN resource_id IS NOT NULL THEN 1 ELSE 2 END,
			CASE WHEN resource_type = 'ALL' THEN 2 ELSE 1 END
		LIMIT 1
	`
	
	var resourceIDVal interface{}
	if resourceID != nil {
		resourceIDVal = *resourceID
	}
	
	err := r.db.QueryRowContext(ctx, query,
		tenantID,
		string(role),
		string(resourceType),
		string(action),
		resourceIDVal,
	).Scan(
		&permission.ID,
		&permission.TenantID,
		&resourceTypeStr,
		&resourceIDPtr,
		&actionStr,
		&roleStr,
		&granted,
	)
	
	if err == sql.ErrNoRows {
		// 権限が見つからない場合は拒否
		return nil, fmt.Errorf("permission denied: no matching permission found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}
	
	permission.ResourceType = domain.ResourceType(resourceTypeStr)
	if resourceIDPtr.Valid {
		permission.ResourceID = &resourceIDPtr.String
	}
	permission.Action = domain.Action(actionStr)
	permission.Role = domain.UserRole(roleStr)
	permission.Granted = granted
	
	return &permission, nil
}

// GetPermissionsByRole ロールごとの権限一覧を取得
func (r *PostgreSQLPermissionRepository) GetPermissionsByRole(ctx context.Context, tenantID string, role domain.UserRole) ([]*domain.Permission, error) {
	query := `
		SELECT id, tenant_id, resource_type, resource_id, action, role, granted
		FROM permissions
		WHERE tenant_id = $1 AND role = $2
		ORDER BY resource_type, action
	`
	
	rows, err := r.db.QueryContext(ctx, query, tenantID, string(role))
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions by role: %w", err)
	}
	defer rows.Close()
	
	var permissions []*domain.Permission
	for rows.Next() {
		var permission domain.Permission
		var resourceTypeStr, actionStr, roleStr string
		var resourceIDPtr sql.NullString
		var granted bool
		
		err := rows.Scan(
			&permission.ID,
			&permission.TenantID,
			&resourceTypeStr,
			&resourceIDPtr,
			&actionStr,
			&roleStr,
			&granted,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		
		permission.ResourceType = domain.ResourceType(resourceTypeStr)
		if resourceIDPtr.Valid {
			permission.ResourceID = &resourceIDPtr.String
		}
		permission.Action = domain.Action(actionStr)
		permission.Role = domain.UserRole(roleStr)
		permission.Granted = granted
		
		permissions = append(permissions, &permission)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate permissions: %w", err)
	}
	
	return permissions, nil
}

// GetPermissionsByResource リソースごとの権限一覧を取得
func (r *PostgreSQLPermissionRepository) GetPermissionsByResource(ctx context.Context, tenantID string, resourceType domain.ResourceType, resourceID *string) ([]*domain.Permission, error) {
	query := `
		SELECT id, tenant_id, resource_type, resource_id, action, role, granted
		FROM permissions
		WHERE tenant_id = $1 AND resource_type = $2
		AND (resource_id = $3 OR resource_id IS NULL)
		ORDER BY role, action
	`
	
	var resourceIDVal interface{}
	if resourceID != nil {
		resourceIDVal = *resourceID
	}
	
	rows, err := r.db.QueryContext(ctx, query, tenantID, string(resourceType), resourceIDVal)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions by resource: %w", err)
	}
	defer rows.Close()
	
	var permissions []*domain.Permission
	for rows.Next() {
		var permission domain.Permission
		var resourceTypeStr, actionStr, roleStr string
		var resourceIDPtr sql.NullString
		var granted bool
		
		err := rows.Scan(
			&permission.ID,
			&permission.TenantID,
			&resourceTypeStr,
			&resourceIDPtr,
			&actionStr,
			&roleStr,
			&granted,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		
		permission.ResourceType = domain.ResourceType(resourceTypeStr)
		if resourceIDPtr.Valid {
			permission.ResourceID = &resourceIDPtr.String
		}
		permission.Action = domain.Action(actionStr)
		permission.Role = domain.UserRole(roleStr)
		permission.Granted = granted
		
		permissions = append(permissions, &permission)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate permissions: %w", err)
	}
	
	return permissions, nil
}

// CreatePermission 権限を作成
func (r *PostgreSQLPermissionRepository) CreatePermission(ctx context.Context, permission *domain.Permission) error {
	query := `
		INSERT INTO permissions (tenant_id, resource_type, resource_id, action, role, granted)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	
	var resourceIDVal interface{}
	if permission.ResourceID != nil {
		resourceIDVal = *permission.ResourceID
	}
	
	err := r.db.QueryRowContext(ctx, query,
		permission.TenantID,
		string(permission.ResourceType),
		resourceIDVal,
		string(permission.Action),
		string(permission.Role),
		permission.Granted,
	).Scan(&permission.ID)
	
	if err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}
	
	return nil
}

// UpdatePermission 権限を更新
func (r *PostgreSQLPermissionRepository) UpdatePermission(ctx context.Context, permission *domain.Permission) error {
	query := `
		UPDATE permissions
		SET resource_type = $3, resource_id = $4, action = $5, role = $6, granted = $7, updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`
	
	var resourceIDVal interface{}
	if permission.ResourceID != nil {
		resourceIDVal = *permission.ResourceID
	}
	
	result, err := r.db.ExecContext(ctx, query,
		permission.ID,
		permission.TenantID,
		string(permission.ResourceType),
		resourceIDVal,
		string(permission.Action),
		string(permission.Role),
		permission.Granted,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update permission: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("permission not found")
	}
	
	return nil
}

// DeletePermission 権限を削除
func (r *PostgreSQLPermissionRepository) DeletePermission(ctx context.Context, permissionID string, tenantID string) error {
	query := `
		DELETE FROM permissions
		WHERE id = $1 AND tenant_id = $2
	`
	
	result, err := r.db.ExecContext(ctx, query, permissionID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("permission not found")
	}
	
	return nil
}

