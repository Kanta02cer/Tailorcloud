package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// MetricsCollector メトリクス収集器
type MetricsCollector struct {
	// HTTPメトリクス
	totalRequests      int64
	totalErrors        int64
	totalLatency       int64 // マイクロ秒単位
	requestCount       int64

	// データベースメトリクス
	dbConnections      int64
	dbConnectionsInUse int64

	// ロック
	mu sync.RWMutex
}

// NewMetricsCollector MetricsCollectorを作成
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{}
}

// RecordRequest リクエストを記録
func (m *MetricsCollector) RecordRequest(latency time.Duration, statusCode int) {
	atomic.AddInt64(&m.totalRequests, 1)
	atomic.AddInt64(&m.requestCount, 1)

	// レイテンシーをマイクロ秒で記録
	latencyMicroseconds := latency.Microseconds()
	atomic.AddInt64(&m.totalLatency, latencyMicroseconds)

	// エラーを記録（4xx, 5xx）
	if statusCode >= 400 {
		atomic.AddInt64(&m.totalErrors, 1)
	}
}

// RecordDBConnection DB接続を記録
func (m *MetricsCollector) RecordDBConnection(connections, inUse int) {
	atomic.StoreInt64(&m.dbConnections, int64(connections))
	atomic.StoreInt64(&m.dbConnectionsInUse, int64(inUse))
}

// GetMetrics メトリクスを取得
func (m *MetricsCollector) GetMetrics() *Metrics {
	totalRequests := atomic.LoadInt64(&m.totalRequests)
	totalErrors := atomic.LoadInt64(&m.totalErrors)
	totalLatency := atomic.LoadInt64(&m.totalLatency)
	requestCount := atomic.LoadInt64(&m.requestCount)

	var avgLatency time.Duration
	if requestCount > 0 {
		avgLatencyMicroseconds := totalLatency / requestCount
		avgLatency = time.Duration(avgLatencyMicroseconds) * time.Microsecond
	}

	var errorRate float64
	if totalRequests > 0 {
		errorRate = float64(totalErrors) / float64(totalRequests) * 100
	}

	return &Metrics{
		TotalRequests:      totalRequests,
		TotalErrors:        totalErrors,
		ErrorRate:          errorRate,
		AverageLatency:     avgLatency,
		RequestCount:       requestCount,
		DBConnections:      atomic.LoadInt64(&m.dbConnections),
		DBConnectionsInUse: atomic.LoadInt64(&m.dbConnectionsInUse),
		Timestamp:          time.Now(),
	}
}

// ResetMetrics メトリクスをリセット
func (m *MetricsCollector) ResetMetrics() {
	atomic.StoreInt64(&m.totalRequests, 0)
	atomic.StoreInt64(&m.totalErrors, 0)
	atomic.StoreInt64(&m.totalLatency, 0)
	atomic.StoreInt64(&m.requestCount, 0)
}

// Metrics メトリクス情報
type Metrics struct {
	TotalRequests      int64         `json:"total_requests"`
	TotalErrors        int64         `json:"total_errors"`
	ErrorRate          float64       `json:"error_rate"`
	AverageLatency     time.Duration `json:"average_latency"`
	RequestCount       int64         `json:"request_count"`
	DBConnections      int64         `json:"db_connections"`
	DBConnectionsInUse int64         `json:"db_connections_in_use"`
	Timestamp          time.Time     `json:"timestamp"`
}

