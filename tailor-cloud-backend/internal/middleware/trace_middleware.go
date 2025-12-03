package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"tailor-cloud/backend/internal/logger"
)

// TraceMiddleware トレースIDミドルウェア
// 各リクエストにトレースIDを付与
type TraceMiddleware struct {
	logger *logger.StructuredLogger
}

// NewTraceMiddleware TraceMiddlewareのコンストラクタ
func NewTraceMiddleware(logger *logger.StructuredLogger) *TraceMiddleware {
	return &TraceMiddleware{
		logger: logger,
	}
}

// Trace トレースIDを付与するミドルウェア
func (m *TraceMiddleware) Trace(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// トレースIDを生成またはヘッダーから取得
		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// コンテキストにトレースIDを追加
		ctx := logger.WithTraceID(r.Context(), traceID)

		// レスポンスヘッダーにトレースIDを追加
		w.Header().Set("X-Trace-ID", traceID)

		// 次のハンドラーにコンテキストを渡す
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

// LoggingMiddleware 構造化ログミドルウェア
// HTTPリクエスト・レスポンスをログに記録
type LoggingMiddleware struct {
	logger *logger.StructuredLogger
}

// NewLoggingMiddleware LoggingMiddlewareのコンストラクタ
func NewLoggingMiddleware(logger *logger.StructuredLogger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// Log リクエスト・レスポンスをログに記録するミドルウェア
func (m *LoggingMiddleware) Log(next http.HandlerFunc) http.HandlerFunc {
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

		// HTTPリクエスト情報を構築
		httpInfo := &logger.HTTPRequestInfo{
			Method:     r.Method,
			Path:       r.URL.Path,
			StatusCode: rr.statusCode,
			UserAgent:  r.Header.Get("User-Agent"),
			IPAddress:  getIPAddress(r),
			Latency:    latency.String(),
		}

		// コンテキストにHTTPリクエスト情報を追加
		ctx := logger.WithHTTPRequest(r.Context(), httpInfo)

		// ログレベルを決定
		logLevel := logger.LogLevelInfo
		if rr.statusCode >= 500 {
			logLevel = logger.LogLevelError
		} else if rr.statusCode >= 400 {
			logLevel = logger.LogLevelWarning
		}

		// ログに記録
		fields := map[string]interface{}{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status_code": rr.statusCode,
			"latency_ms":  latency.Milliseconds(),
		}

		if logLevel == logger.LogLevelInfo {
			m.logger.Info(ctx, "HTTP request completed", fields)
		} else if logLevel == logger.LogLevelWarning {
			m.logger.Warning(ctx, "HTTP request completed with client error", fields)
		} else {
			m.logger.Error(ctx, "HTTP request completed with server error", nil, fields)
		}
	}
}

// responseRecorder レスポンスを記録するラッパー
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader ステータスコードを記録
func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

// getIPAddress IPアドレスを取得
func getIPAddress(r *http.Request) string {
	// X-Forwarded-Forヘッダーを確認（プロキシ経由の場合）
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// X-Real-IPヘッダーを確認
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// RemoteAddrから取得
	return r.RemoteAddr
}

