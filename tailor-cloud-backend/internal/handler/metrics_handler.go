package handler

import (
	"encoding/json"
	"net/http"

	"tailor-cloud/backend/internal/metrics"
)

// MetricsHandler メトリクスハンドラー
type MetricsHandler struct {
	collector *metrics.MetricsCollector
}

// NewMetricsHandler MetricsHandlerのコンストラクタ
func NewMetricsHandler(collector *metrics.MetricsCollector) *MetricsHandler {
	return &MetricsHandler{
		collector: collector,
	}
}

// GetMetrics GET /api/metrics - メトリクスを取得
func (h *MetricsHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics := h.collector.GetMetrics()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}

