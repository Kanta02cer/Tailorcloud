package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"tailor-cloud/backend/internal/middleware"
	"tailor-cloud/backend/internal/service"
)

// AmbassadorHandler アンバサダーハンドラー
type AmbassadorHandler struct {
	ambassadorService *service.AmbassadorService
}

// NewAmbassadorHandler AmbassadorHandlerのコンストラクタ
func NewAmbassadorHandler(ambassadorService *service.AmbassadorService) *AmbassadorHandler {
	return &AmbassadorHandler{
		ambassadorService: ambassadorService,
	}
}

// CreateAmbassadorRequest アンバサダー作成リクエスト
type CreateAmbassadorRequest struct {
	TenantID       string  `json:"tenant_id"`
	UserID         string  `json:"user_id"` // Firebase AuthのUserID
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	Phone          string  `json:"phone"`
	CommissionRate float64 `json:"commission_rate"` // 成果報酬率（オプション）
}

// CreateAmbassador POST /api/ambassadors - アンバサダーを作成
func (h *AmbassadorHandler) CreateAmbassador(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req CreateAmbassadorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		if req.TenantID == "" {
			http.Error(w, "tenant_id is required", http.StatusBadRequest)
			return
		}
		authUser = &middleware.AuthUser{TenantID: req.TenantID}
	}
	
	// テナントID: 認証ユーザーから取得、またはリクエストから
	tenantID := authUser.TenantID
	if tenantID == "" {
		tenantID = req.TenantID
		if tenantID == "" {
			http.Error(w, "tenant_id is required", http.StatusBadRequest)
			return
		}
	}
	
	// サービス層で作成
	serviceReq := &service.CreateAmbassadorRequest{
		TenantID:       tenantID,
		UserID:         req.UserID,
		Name:           req.Name,
		Email:          req.Email,
		Phone:          req.Phone,
		CommissionRate: req.CommissionRate,
	}
	
	ambassador, err := h.ambassadorService.CreateAmbassador(r.Context(), serviceReq)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "ambassador already exists" {
			statusCode = http.StatusConflict
		}
		http.Error(w, "Failed to create ambassador: "+err.Error(), statusCode)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ambassador)
}

// GetAmbassadorByUserID GET /api/ambassadors/me - 自分のアンバサダー情報を取得
func (h *AmbassadorHandler) GetAmbassadorByUserID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Authentication required: "+err.Error(), http.StatusUnauthorized)
		return
	}
	
	// サービス層で取得
	ambassador, err := h.ambassadorService.GetAmbassadorByUserID(r.Context(), authUser.ID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "ambassador not found" {
			statusCode = http.StatusNotFound
		}
		http.Error(w, "Failed to get ambassador: "+err.Error(), statusCode)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ambassador)
}

// ListAmbassadors GET /api/ambassadors - アンバサダー一覧を取得
func (h *AmbassadorHandler) ListAmbassadors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		tenantID := r.URL.Query().Get("tenant_id")
		if tenantID == "" {
			http.Error(w, "tenant_id is required", http.StatusBadRequest)
			return
		}
		authUser = &middleware.AuthUser{TenantID: tenantID}
	}
	
	tenantID := authUser.TenantID
	if tenantID == "" {
		tenantID = r.URL.Query().Get("tenant_id")
		if tenantID == "" {
			http.Error(w, "tenant_id is required", http.StatusBadRequest)
			return
		}
	}
	
	// サービス層で一覧取得
	req := &service.ListAmbassadorsRequest{
		TenantID: tenantID,
	}
	
	ambassadors, err := h.ambassadorService.ListAmbassadors(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to list ambassadors: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ambassadors)
}

// GetCommissions GET /api/ambassadors/{ambassador_id}/commissions - 成果報酬一覧を取得
func (h *AmbassadorHandler) GetCommissions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Authentication required: "+err.Error(), http.StatusUnauthorized)
		return
	}
	
	// URLからambassador_idを取得
	ambassadorID := r.URL.Query().Get("ambassador_id")
	if ambassadorID == "" {
		// 自分の成果報酬を取得する場合
		ambassador, err := h.ambassadorService.GetAmbassadorByUserID(r.Context(), authUser.ID)
		if err != nil {
			http.Error(w, "Failed to get ambassador: "+err.Error(), http.StatusInternalServerError)
			return
		}
		ambassadorID = ambassador.ID
	}
	
	// ページネーション
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	
	limit := 20
	offset := 0
	
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}
	
	// サービス層で取得
	req := &service.GetCommissionsByAmbassadorRequest{
		AmbassadorID: ambassadorID,
		Limit:        limit,
		Offset:       offset,
	}
	
	commissions, err := h.ambassadorService.GetCommissionsByAmbassador(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to get commissions: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"commissions": commissions,
		"total":       len(commissions),
		"limit":       limit,
		"offset":      offset,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

