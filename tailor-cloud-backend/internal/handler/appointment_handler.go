package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/middleware"
	"tailor-cloud/backend/internal/service"
)

// AppointmentHandler 予約ハンドラー
type AppointmentHandler struct {
	appointmentService *service.AppointmentService
}

// NewAppointmentHandler AppointmentHandlerのコンストラクタ
func NewAppointmentHandler(appointmentService *service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
	}
}

// CreateAppointmentRequest 予約作成リクエスト
type CreateAppointmentRequest struct {
	UserID              string `json:"user_id"`
	FitterID            string `json:"fitter_id"`
	AppointmentDateTime string `json:"appointment_datetime"` // ISO 8601形式
	DurationMinutes     int    `json:"duration_minutes"`
	Notes               string `json:"notes,omitempty"`
}

// CreateAppointment POST /api/appointments - 予約を作成
func (h *AppointmentHandler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateAppointmentRequest
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

	// 日時文字列をパース
	appointmentDateTime, err := time.Parse(time.RFC3339, req.AppointmentDateTime)
	if err != nil {
		http.Error(w, "Invalid appointment_datetime format. Use ISO 8601 (RFC3339)", http.StatusBadRequest)
		return
	}

	// サービス層で作成
	serviceReq := &service.CreateAppointmentRequest{
		UserID:              req.UserID,
		TenantID:            authUser.TenantID,
		FitterID:            req.FitterID,
		AppointmentDateTime: appointmentDateTime,
		DurationMinutes:     req.DurationMinutes,
		Notes:               req.Notes,
	}

	appointment, err := h.appointmentService.CreateAppointment(r.Context(), serviceReq)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not available") {
			statusCode = http.StatusConflict
		} else if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "required") {
			statusCode = http.StatusBadRequest
		}
		http.Error(w, "Failed to create appointment: "+err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(appointment)
}

// GetAppointment GET /api/appointments/{id} - 予約を取得
func (h *AppointmentHandler) GetAppointment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLパスまたはクエリパラメータから予約IDを取得
	appointmentID := r.URL.Query().Get("appointment_id")
	if appointmentID == "" {
		if id := r.PathValue("id"); id != "" {
			appointmentID = id
		}
	}
	if appointmentID == "" {
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) >= 3 && pathParts[len(pathParts)-2] == "appointments" {
			appointmentID = pathParts[len(pathParts)-1]
		}
	}

	if appointmentID == "" {
		http.Error(w, "appointment_id is required", http.StatusBadRequest)
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
	appointment, err := h.appointmentService.GetAppointment(r.Context(), appointmentID, authUser.TenantID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		http.Error(w, "Failed to get appointment: "+err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appointment)
}

// ListAppointments GET /api/appointments - 予約一覧を取得
func (h *AppointmentHandler) ListAppointments(w http.ResponseWriter, r *http.Request) {
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
	fitterID := r.URL.Query().Get("fitter_id")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate *time.Time
	if startDateStr != "" {
		parsed, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			http.Error(w, "Invalid start_date format. Use ISO 8601 (RFC3339)", http.StatusBadRequest)
			return
		}
		startDate = &parsed
	}
	if endDateStr != "" {
		parsed, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			http.Error(w, "Invalid end_date format. Use ISO 8601 (RFC3339)", http.StatusBadRequest)
			return
		}
		endDate = &parsed
	}

	// サービス層で取得
	serviceReq := &service.ListAppointmentsRequest{
		TenantID:  authUser.TenantID,
		UserID:    userID,
		FitterID:  fitterID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	appointments, err := h.appointmentService.ListAppointments(r.Context(), serviceReq)
	if err != nil {
		http.Error(w, "Failed to list appointments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": appointments,
		"total": len(appointments),
	})
}

// UpdateAppointmentRequest 予約更新リクエスト
type UpdateAppointmentRequest struct {
	FitterID            string `json:"fitter_id,omitempty"`
	AppointmentDateTime string `json:"appointment_datetime,omitempty"` // ISO 8601形式
	DurationMinutes     int    `json:"duration_minutes,omitempty"`
	Status              string `json:"status,omitempty"`
	Notes               string `json:"notes,omitempty"`
}

// UpdateAppointment PUT /api/appointments/{id} - 予約を更新
func (h *AppointmentHandler) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLパスまたはクエリパラメータから予約IDを取得
	appointmentID := r.URL.Query().Get("appointment_id")
	if appointmentID == "" {
		if id := r.PathValue("id"); id != "" {
			appointmentID = id
		}
	}
	if appointmentID == "" {
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) >= 3 && pathParts[len(pathParts)-2] == "appointments" {
			appointmentID = pathParts[len(pathParts)-1]
		}
	}

	if appointmentID == "" {
		http.Error(w, "appointment_id is required", http.StatusBadRequest)
		return
	}

	var req UpdateAppointmentRequest
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

	// サービス層で更新
	serviceReq := &service.UpdateAppointmentRequest{
		AppointmentID: appointmentID,
		TenantID:      authUser.TenantID,
		FitterID:      req.FitterID,
		DurationMinutes: req.DurationMinutes,
		Status:        domain.AppointmentStatus(req.Status),
		Notes:         req.Notes,
	}

	if req.AppointmentDateTime != "" {
		parsed, err := time.Parse(time.RFC3339, req.AppointmentDateTime)
		if err != nil {
			http.Error(w, "Invalid appointment_datetime format. Use ISO 8601 (RFC3339)", http.StatusBadRequest)
			return
		}
		serviceReq.AppointmentDateTime = &parsed
	}

	appointment, err := h.appointmentService.UpdateAppointment(r.Context(), serviceReq)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "not available") {
			statusCode = http.StatusConflict
		} else if strings.Contains(err.Error(), "invalid") {
			statusCode = http.StatusBadRequest
		}
		http.Error(w, "Failed to update appointment: "+err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appointment)
}

// CancelAppointmentRequest 予約キャンセルリクエスト
type CancelAppointmentRequest struct {
	Reason string `json:"reason"`
}

// CancelAppointment DELETE /api/appointments/{id} - 予約をキャンセル
func (h *AppointmentHandler) CancelAppointment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLパスまたはクエリパラメータから予約IDを取得
	appointmentID := r.URL.Query().Get("appointment_id")
	if appointmentID == "" {
		if id := r.PathValue("id"); id != "" {
			appointmentID = id
		}
	}
	if appointmentID == "" {
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) >= 3 && pathParts[len(pathParts)-2] == "appointments" {
			appointmentID = pathParts[len(pathParts)-1]
		}
	}

	if appointmentID == "" {
		http.Error(w, "appointment_id is required", http.StatusBadRequest)
		return
	}

	var req CancelAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// リクエストボディがない場合は空文字列でキャンセル
		req.Reason = ""
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

	// サービス層でキャンセル
	if err := h.appointmentService.CancelAppointment(r.Context(), appointmentID, authUser.TenantID, req.Reason); err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "already") || strings.Contains(err.Error(), "cannot") {
			statusCode = http.StatusBadRequest
		}
		http.Error(w, "Failed to cancel appointment: "+err.Error(), statusCode)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

