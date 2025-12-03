package handler

import (
	"encoding/json"
	"net/http"

	"tailor-cloud/backend/internal/middleware"
	"tailor-cloud/backend/internal/service"
)

// MeasurementValidationHandler 採寸データバリデーションハンドラー
type MeasurementValidationHandler struct {
	validationService *service.MeasurementValidationService
}

// NewMeasurementValidationHandler MeasurementValidationHandlerのコンストラクタ
func NewMeasurementValidationHandler(validationService *service.MeasurementValidationService) *MeasurementValidationHandler {
	return &MeasurementValidationHandler{
		validationService: validationService,
	}
}

// ValidateMeasurementsRequest APIリクエスト
type ValidateMeasurementsRequest struct {
	CustomerID         string          `json:"customer_id"`
	CurrentMeasurements json.RawMessage `json:"current_measurements"`
}

// ValidateMeasurements POST /api/measurements/validate - 採寸データをバリデーション
func (h *MeasurementValidationHandler) ValidateMeasurements(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ValidateMeasurementsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// バリデーション
	if req.CustomerID == "" {
		http.Error(w, "customer_id is required", http.StatusBadRequest)
		return
	}
	if len(req.CurrentMeasurements) == 0 {
		http.Error(w, "current_measurements is required", http.StatusBadRequest)
		return
	}

	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	tenantID := ""
	if err == nil {
		tenantID = authUser.TenantID
	} else {
		// 認証がない場合はクエリパラメータから取得
		tenantID = r.URL.Query().Get("tenant_id")
		if tenantID == "" {
			http.Error(w, "Authentication required or tenant_id must be provided", http.StatusUnauthorized)
			return
		}
	}

	// サービス層でバリデーション
	serviceReq := &service.ValidateMeasurementsRequest{
		CustomerID:         req.CustomerID,
		TenantID:           tenantID,
		CurrentMeasurements: req.CurrentMeasurements,
	}

	response, err := h.validationService.ValidateMeasurements(r.Context(), serviceReq)
	if err != nil {
		http.Error(w, "Failed to validate measurements: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ValidateMeasurementRange POST /api/measurements/validate-range - 採寸データの範囲をバリデーション
func (h *MeasurementValidationHandler) ValidateMeasurementRange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Measurements json.RawMessage `json:"measurements"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Measurements) == 0 {
		http.Error(w, "measurements is required", http.StatusBadRequest)
		return
	}

	// サービス層で範囲バリデーション
	alerts := h.validationService.ValidateMeasurementRange(req.Measurements)

	response := map[string]interface{}{
		"is_valid": len(alerts) == 0,
		"alerts":   alerts,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

