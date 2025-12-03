package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/service"
)

// RBACMiddleware ロールベースアクセス制御ミドルウェア
type RBACMiddleware struct {
	rbacService *service.RBACService
}

// NewRBACMiddleware RBACMiddlewareのコンストラクタ
func NewRBACMiddleware(rbacService *service.RBACService) *RBACMiddleware {
	return &RBACMiddleware{
		rbacService: rbacService,
	}
}

// RequireRole 特定のロールを要求するミドルウェア
func (m *RBACMiddleware) RequireRole(allowedRoles ...domain.UserRole) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, err := GetUserFromContext(r.Context())
			if err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			// ユーザーのロールを取得
			userRole := domain.UserRole(user.Role)

			// 許可されたロールかチェック
			for _, allowedRole := range allowedRoles {
				if userRole == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			// 権限不足
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Insufficient permissions. Required role: %v", allowedRoles),
			})
		}
	}
}

// RequireAnyRole いずれかのロールを要求するミドルウェア
func (m *RBACMiddleware) RequireAnyRole(allowedRoles ...domain.UserRole) func(http.HandlerFunc) http.HandlerFunc {
	return m.RequireRole(allowedRoles...)
}

// RequireOwnerOrStaff OwnerまたはStaffロールを要求
func (m *RBACMiddleware) RequireOwnerOrStaff() func(http.HandlerFunc) http.HandlerFunc {
	return m.RequireRole(domain.RoleOwner, domain.RoleStaff)
}

// RequireOwnerOnly Ownerのみ許可
func (m *RBACMiddleware) RequireOwnerOnly() func(http.HandlerFunc) http.HandlerFunc {
	return m.RequireRole(domain.RoleOwner)
}

// CheckTenantAccess テナントアクセスチェック
// リクエストのtenant_idが、ユーザーのtenant_idと一致しているかチェック
func CheckTenantAccess(r *http.Request, requestedTenantID string) error {
	user, err := GetUserFromContext(r.Context())
	if err != nil {
		return fmt.Errorf("authentication required: %w", err)
	}

	// Ownerロールは全テナントにアクセス可能（将来的なマルチテナント管理用）
	if user.Role == string(domain.RoleOwner) {
		return nil
	}

	// ユーザーのテナントIDと一致しているかチェック
	if user.TenantID != requestedTenantID {
		return fmt.Errorf("unauthorized: tenant_id mismatch")
	}

	return nil
}

// RequirePermission リソースベースの権限チェックミドルウェア
// 特定のリソースに対する特定の操作を許可するかチェック
func (m *RBACMiddleware) RequirePermission(
	resourceType domain.ResourceType,
	action domain.Action,
	resourceIDExtractor func(*http.Request) *string, // リクエストからリソースIDを抽出する関数
) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// RBACServiceが設定されていない場合は従来のロールチェックを使用
			if m.rbacService == nil {
				// フォールバック: 従来のロールチェック
				next.ServeHTTP(w, r)
				return
			}
			
			user, err := GetUserFromContext(r.Context())
			if err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
			
			// リソースIDを抽出
			var resourceID *string
			if resourceIDExtractor != nil {
				resourceID = resourceIDExtractor(r)
			}
			
			// 権限をチェック
			allowed, err := m.rbacService.CheckPermission(
				r.Context(),
				user.TenantID,
				resourceType,
				action,
				domain.UserRole(user.Role),
				resourceID,
			)
			
			if err != nil || !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]string{
					"error": fmt.Sprintf("Permission denied: %s on %s", action, resourceType),
				})
				return
			}
			
			next.ServeHTTP(w, r)
		}
	}
}

