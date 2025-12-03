package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/middleware"
	"tailor-cloud/backend/internal/service"
)

// DiagnosisHandler 診断ハンドラー
type DiagnosisHandler struct {
	diagnosisService *service.DiagnosisService
}

// NewDiagnosisHandler DiagnosisHandlerのコンストラクタ
func NewDiagnosisHandler(diagnosisService *service.DiagnosisService) *DiagnosisHandler {
	return &DiagnosisHandler{
		diagnosisService: diagnosisService,
	}
}

// CreateDiagnosisRequest 診断作成リクエスト
type CreateDiagnosisRequest struct {
	UserID          string          `json:"user_id"`
	Archetype       string          `json:"archetype"`
	PlanType        string          `json:"plan_type,omitempty"`
	DiagnosisResult json.RawMessage `json:"diagnosis_result"`
}

// CreateDiagnosis POST /api/diagnoses - 診断を作成
func (h *DiagnosisHandler) CreateDiagnosis(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateDiagnosisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		tenantID := r.URL.Query().Get("tenant_id")
		if tenantID == "" {
			http.Error(w, "Authentication required or tenant_id must be provided", http.StatusUnauthorized)
			return
		}
		authUser = &middleware.AuthUser{TenantID: tenantID}
	}

	// サービス層で作成
	serviceReq := &service.CreateDiagnosisRequest{
		UserID:          req.UserID,
		TenantID:        authUser.TenantID,
		Archetype:       domain.Archetype(req.Archetype),
		PlanType:        domain.PlanType(req.PlanType),
		DiagnosisResult: req.DiagnosisResult,
	}

	diagnosis, err := h.diagnosisService.CreateDiagnosis(r.Context(), serviceReq)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "invalid") {
			statusCode = http.StatusBadRequest
		}
		http.Error(w, "Failed to create diagnosis: "+err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(diagnosis)
}

// GetDiagnosis GET /api/diagnoses/{id} - 診断を取得
func (h *DiagnosisHandler) GetDiagnosis(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLパスまたはクエリパラメータから診断IDを取得
	diagnosisID := r.URL.Query().Get("diagnosis_id")
	if diagnosisID == "" {
		if id := r.PathValue("id"); id != "" {
			diagnosisID = id
		}
	}
	if diagnosisID == "" {
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) >= 3 && pathParts[len(pathParts)-2] == "diagnoses" {
			diagnosisID = pathParts[len(pathParts)-1]
		}
	}

	if diagnosisID == "" {
		http.Error(w, "diagnosis_id is required", http.StatusBadRequest)
		return
	}

	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		tenantID := r.URL.Query().Get("tenant_id")
		if tenantID == "" {
			http.Error(w, "Authentication required or tenant_id must be provided", http.StatusUnauthorized)
			return
		}
		authUser = &middleware.AuthUser{TenantID: tenantID}
	}

	// サービス層で取得
	diagnosis, err := h.diagnosisService.GetDiagnosis(r.Context(), diagnosisID, authUser.TenantID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		http.Error(w, "Failed to get diagnosis: "+err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(diagnosis)
}

// ListDiagnoses GET /api/diagnoses - 診断一覧を取得
func (h *DiagnosisHandler) ListDiagnoses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		tenantID := r.URL.Query().Get("tenant_id")
		if tenantID == "" {
			http.Error(w, "Authentication required or tenant_id must be provided", http.StatusUnauthorized)
			return
		}
		authUser = &middleware.AuthUser{TenantID: tenantID}
	}

	// クエリパラメータからフィルター条件を取得
	userID := r.URL.Query().Get("user_id")
	archetype := r.URL.Query().Get("archetype")
	planType := r.URL.Query().Get("plan_type")
	
	limit := 20 // デフォルト値
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
			if limit > 100 {
				limit = 100 // 最大値
			}
		}
	}

	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	var diagnoses []*domain.Diagnosis
	var listErr error

	if userID != "" {
		// ユーザーIDで取得
		diagnoses, listErr = h.diagnosisService.GetDiagnosesByUser(r.Context(), userID, authUser.TenantID)
	} else {
		// テナントIDで取得（ページネーション対応）
		serviceReq := &service.GetDiagnosesByTenantRequest{
			TenantID: authUser.TenantID,
			Limit:    limit,
			Offset:   offset,
		}
		diagnoses, listErr = h.diagnosisService.GetDiagnosesByTenant(r.Context(), serviceReq)
	}

	if listErr != nil {
		http.Error(w, "Failed to list diagnoses: "+listErr.Error(), http.StatusInternalServerError)
		return
	}

		// フィルター適用（必要に応じて）
	if archetype != "" || planType != "" {
		filtered := make([]*domain.Diagnosis, 0)
		for _, d := range diagnoses {
			if archetype != "" && string(d.Archetype) != archetype {
				continue
			}
			if planType != "" && string(d.PlanType) != planType {
				continue
			}
			filtered = append(filtered, d)
		}
		diagnoses = filtered
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": diagnoses,
		"total": len(diagnoses),
	})
}

