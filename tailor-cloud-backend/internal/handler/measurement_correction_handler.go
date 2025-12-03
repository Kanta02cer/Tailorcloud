package handler

import (
	"encoding/json"
	"net/http"

	"tailor-cloud/backend/internal/middleware"
	"tailor-cloud/backend/internal/service"
)

// MeasurementCorrectionHandler 自動補正エンジンハンドラー
type MeasurementCorrectionHandler struct {
	correctionService *service.MeasurementCorrectionService
}

// NewMeasurementCorrectionHandler MeasurementCorrectionHandlerのコンストラクタ
func NewMeasurementCorrectionHandler(correctionService *service.MeasurementCorrectionService) *MeasurementCorrectionHandler {
	return &MeasurementCorrectionHandler{
		correctionService: correctionService,
	}
}

// ConvertToFinalMeasurementsRequest APIリクエスト
type ConvertToFinalMeasurementsRequest struct {
	RawMeasurements *service.RawMeasurement `json:"raw_measurements"`
	UserID          string                   `json:"user_id"` // オプション（診断プロファイル取得用）
	FabricID        string                   `json:"fabric_id"`
}

// ConvertToFinalMeasurements POST /api/measurements/convert - ヌード寸を仕上がり寸法に変換
func (h *MeasurementCorrectionHandler) ConvertToFinalMeasurements(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ConvertToFinalMeasurementsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// バリデーション
	if req.RawMeasurements == nil {
		http.Error(w, "raw_measurements is required", http.StatusBadRequest)
		return
	}
	if req.FabricID == "" {
		http.Error(w, "fabric_id is required", http.StatusBadRequest)
		return
	}

	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	tenantID := ""
	if err == nil {
		tenantID = authUser.TenantID
		// UserIDが指定されていない場合は認証ユーザーのIDを使用
		if req.UserID == "" {
			req.UserID = authUser.ID
		}
	} else {
		// 認証がない場合はクエリパラメータから取得
		tenantID = r.URL.Query().Get("tenant_id")
		if tenantID == "" {
			http.Error(w, "Authentication required or tenant_id must be provided", http.StatusUnauthorized)
			return
		}
	}

	// サービス層で変換
	serviceReq := &service.ConvertToFinalMeasurementsRequest{
		RawMeasurements: req.RawMeasurements,
		UserID:          req.UserID,
		TenantID:        tenantID,
		FabricID:        req.FabricID,
	}

	response, err := h.correctionService.ConvertToFinalMeasurements(r.Context(), serviceReq)
	if err != nil {
		http.Error(w, "Failed to convert measurements: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

