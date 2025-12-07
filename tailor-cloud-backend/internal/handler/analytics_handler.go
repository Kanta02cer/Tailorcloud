package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"tailor-cloud/backend/internal/middleware"
	"tailor-cloud/backend/internal/service"
)

// AnalyticsHandler ダッシュボードAPI
type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

// NewAnalyticsHandler コンストラクタ
func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetSummary GET /api/analytics/summary
func (h *AnalyticsHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.analyticsService == nil {
		http.Error(w, "Analytics service unavailable", http.StatusServiceUnavailable)
		return
	}

	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		tenantID := r.URL.Query().Get("tenant_id")
		if tenantID == "" {
			http.Error(w, "Authentication required or tenant_id must be provided", http.StatusUnauthorized)
			return
		}
		authUser = &middleware.AuthUser{TenantID: tenantID}
	}

	rangeDays := 30
	if v := r.URL.Query().Get("range_days"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			rangeDays = parsed
		}
	}

	summary, err := h.analyticsService.GetSummary(r.Context(), &service.AnalyticsSummaryRequest{
		TenantID:  authUser.TenantID,
		RangeDays: rangeDays,
	})
	if err != nil {
		http.Error(w, "Failed to calculate analytics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
