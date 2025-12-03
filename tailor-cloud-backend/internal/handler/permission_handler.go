package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/middleware"
	"tailor-cloud/backend/internal/service"
)

// PermissionHandler 権限ハンドラー
type PermissionHandler struct {
	rbacService *service.RBACService
}

// NewPermissionHandler PermissionHandlerのコンストラクタ
func NewPermissionHandler(rbacService *service.RBACService) *PermissionHandler {
	return &PermissionHandler{
		rbacService: rbacService,
	}
}

// CreatePermissionRequest 権限作成リクエスト
type CreatePermissionRequest struct {
	ResourceType string  `json:"resource_type"`
	ResourceID   *string `json:"resource_id,omitempty"`
	Action       string  `json:"action"`
	Role         string  `json:"role"`
	Granted      bool    `json:"granted"`
}

// CreatePermission POST /api/permissions - 権限を作成
func (h *PermissionHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// 権限オブジェクトを作成
	permission := &domain.Permission{
		ID:           uuid.New().String(),
		TenantID:     authUser.TenantID,
		ResourceType: domain.ResourceType(req.ResourceType),
		ResourceID:   req.ResourceID,
		Action:       domain.Action(req.Action),
		Role:         domain.UserRole(req.Role),
		Granted:      req.Granted,
	}

	// 権限を作成
	if err := h.rbacService.GrantPermission(r.Context(), permission); err != nil {
		http.Error(w, "Failed to create permission: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(permission)
}

// GetPermissions GET /api/permissions - 権限一覧を取得
func (h *PermissionHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// ロールフィルター
	role := r.URL.Query().Get("role")
	if role == "" {
		role = authUser.Role // デフォルトは自分のロール
	}

	// 権限一覧を取得
	permissions, err := h.rbacService.GetUserPermissions(r.Context(), authUser.TenantID, domain.UserRole(role))
	if err != nil {
		http.Error(w, "Failed to get permissions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"permissions": permissions,
		"total":       len(permissions),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CheckPermissionRequest 権限チェックリクエスト
type CheckPermissionRequest struct {
	ResourceType string  `json:"resource_type"`
	ResourceID   *string `json:"resource_id,omitempty"`
	Action       string  `json:"action"`
	Role         *string `json:"role,omitempty"` // 指定しない場合は自分のロール
}

// CheckPermission POST /api/permissions/check - 権限をチェック
func (h *PermissionHandler) CheckPermission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CheckPermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// ロールの決定
	role := authUser.Role
	if req.Role != nil {
		role = *req.Role
	}

	// 権限をチェック
	check, err := h.rbacService.CheckPermissionWithReason(
		r.Context(),
		authUser.TenantID,
		domain.ResourceType(req.ResourceType),
		domain.Action(req.Action),
		domain.UserRole(role),
		req.ResourceID,
	)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "permission denied") {
			statusCode = http.StatusForbidden
		}
		http.Error(w, "Failed to check permission: "+err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if check.Allowed {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
	json.NewEncoder(w).Encode(check)
}

