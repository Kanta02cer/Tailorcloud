package middleware

import (
	"net/http"
	"time"

	"tailor-cloud/backend/internal/metrics"
)

// MetricsMiddleware メトリクス収集ミドルウェア
type MetricsMiddleware struct {
	collector *metrics.MetricsCollector
}

// NewMetricsMiddleware MetricsMiddlewareのコンストラクタ
func NewMetricsMiddleware(collector *metrics.MetricsCollector) *MetricsMiddleware {
	return &MetricsMiddleware{
		collector: collector,
	}
}

// Collect メトリクスを収集するミドルウェア
func (m *MetricsMiddleware) Collect(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// レスポンスライターをラップしてステータスコードを記録
		rr := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// 次のハンドラーを実行
		next.ServeHTTP(rr, r)

		// レイテンシーを計算
		latency := time.Since(startTime)

		// メトリクスを記録
		m.collector.RecordRequest(latency, rr.statusCode)
	}
}

