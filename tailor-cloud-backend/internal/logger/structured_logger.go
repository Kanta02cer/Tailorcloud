package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

// LogLevel ログレベル
type LogLevel string

const (
	LogLevelDebug   LogLevel = "DEBUG"
	LogLevelInfo    LogLevel = "INFO"
	LogLevelWarning LogLevel = "WARNING"
	LogLevelError   LogLevel = "ERROR"
	LogLevelFatal   LogLevel = "FATAL"
)

// StructuredLog 構造化ログエントリ
type StructuredLog struct {
	Timestamp   string                 `json:"timestamp"`
	Level       LogLevel               `json:"level"`
	Message     string                 `json:"message"`
	TraceID     string                 `json:"trace_id,omitempty"`
	Service     string                 `json:"service"`
	Fields      map[string]interface{} `json:"fields,omitempty"`
	Error       *ErrorInfo             `json:"error,omitempty"`
	HTTPRequest *HTTPRequestInfo       `json:"http_request,omitempty"`
}

// ErrorInfo エラー情報
type ErrorInfo struct {
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
	Stack   string `json:"stack,omitempty"`
}

// HTTPRequestInfo HTTPリクエスト情報
type HTTPRequestInfo struct {
	Method     string `json:"method"`
	Path       string `json:"path"`
	StatusCode int    `json:"status_code,omitempty"`
	UserAgent  string `json:"user_agent,omitempty"`
	IPAddress  string `json:"ip_address,omitempty"`
	Latency    string `json:"latency,omitempty"`
}

// StructuredLogger 構造化ロガー
type StructuredLogger struct {
	service string
	level   LogLevel
	output  *os.File
}

// LoggerOption ロガーオプション
type LoggerOption func(*StructuredLogger)

// WithService サービス名を設定
func WithService(service string) LoggerOption {
	return func(l *StructuredLogger) {
		l.service = service
	}
}

// WithLevel ログレベルを設定
func WithLevel(level LogLevel) LoggerOption {
	return func(l *StructuredLogger) {
		l.level = level
	}
}

// WithOutput 出力先を設定
func WithOutput(output *os.File) LoggerOption {
	return func(l *StructuredLogger) {
		l.output = output
	}
}

// NewStructuredLogger 構造化ロガーを作成
func NewStructuredLogger(opts ...LoggerOption) *StructuredLogger {
	logger := &StructuredLogger{
		service: "tailorcloud-backend",
		level:   LogLevelInfo,
		output:  os.Stdout,
	}

	for _, opt := range opts {
		opt(logger)
	}

	return logger
}

// log 内部ログ関数
func (l *StructuredLogger) log(level LogLevel, message string, fields map[string]interface{}) {
	// ログレベルのフィルタリング
	if !l.shouldLog(level) {
		return
	}

	logEntry := StructuredLog{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level,
		Message:   message,
		Service:   l.service,
		Fields:    fields,
	}

	// JSON形式で出力
	jsonBytes, err := json.Marshal(logEntry)
	if err != nil {
		// JSONエンコードに失敗した場合は、フォールバック
		fmt.Fprintf(l.output, "{\"timestamp\":\"%s\",\"level\":\"ERROR\",\"message\":\"Failed to encode log: %v\"}\n",
			time.Now().UTC().Format(time.RFC3339Nano), err)
		return
	}

	fmt.Fprintf(l.output, "%s\n", string(jsonBytes))
}

// shouldLog ログを出力すべきか判定
func (l *StructuredLogger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		LogLevelDebug:   0,
		LogLevelInfo:    1,
		LogLevelWarning: 2,
		LogLevelError:   3,
		LogLevelFatal:   4,
	}

	return levels[level] >= levels[l.level]
}

// Debug デバッグログ
func (l *StructuredLogger) Debug(ctx context.Context, message string, fields map[string]interface{}) {
	l.logWithContext(ctx, LogLevelDebug, message, fields, nil)
}

// Info 情報ログ
func (l *StructuredLogger) Info(ctx context.Context, message string, fields map[string]interface{}) {
	l.logWithContext(ctx, LogLevelInfo, message, fields, nil)
}

// Warning 警告ログ
func (l *StructuredLogger) Warning(ctx context.Context, message string, fields map[string]interface{}) {
	l.logWithContext(ctx, LogLevelWarning, message, fields, nil)
}

// Error エラーログ
func (l *StructuredLogger) Error(ctx context.Context, message string, err error, fields map[string]interface{}) {
	var errorInfo *ErrorInfo
	if err != nil {
		errorInfo = &ErrorInfo{
			Message: err.Error(),
			Type:    fmt.Sprintf("%T", err),
		}
	}
	l.logWithContext(ctx, LogLevelError, message, fields, errorInfo)
}

// Fatal 致命的エラーログ
func (l *StructuredLogger) Fatal(ctx context.Context, message string, err error, fields map[string]interface{}) {
	var errorInfo *ErrorInfo
	if err != nil {
		errorInfo = &ErrorInfo{
			Message: err.Error(),
			Type:    fmt.Sprintf("%T", err),
		}
	}
	l.logWithContext(ctx, LogLevelFatal, message, fields, errorInfo)
	os.Exit(1)
}

// logWithContext コンテキスト付きログ
func (l *StructuredLogger) logWithContext(ctx context.Context, level LogLevel, message string, fields map[string]interface{}, errorInfo *ErrorInfo) {
	// ログレベルのフィルタリング
	if !l.shouldLog(level) {
		return
	}

	logEntry := StructuredLog{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level,
		Message:   message,
		Service:   l.service,
		Fields:    fields,
		Error:     errorInfo,
	}

	// コンテキストからトレースIDを取得
	if traceID := ctx.Value("trace_id"); traceID != nil {
		if id, ok := traceID.(string); ok {
			logEntry.TraceID = id
		}
	}

	// HTTPリクエスト情報をコンテキストから取得
	if httpInfo := ctx.Value("http_request"); httpInfo != nil {
		if info, ok := httpInfo.(*HTTPRequestInfo); ok {
			logEntry.HTTPRequest = info
		}
	}

	// JSON形式で出力
	jsonBytes, err := json.Marshal(logEntry)
	if err != nil {
		// JSONエンコードに失敗した場合は、フォールバック
		fmt.Fprintf(l.output, "{\"timestamp\":\"%s\",\"level\":\"ERROR\",\"message\":\"Failed to encode log: %v\"}\n",
			time.Now().UTC().Format(time.RFC3339Nano), err)
		return
	}

	fmt.Fprintf(l.output, "%s\n", string(jsonBytes))
}

// WithTraceID トレースIDをコンテキストに追加
func WithTraceID(ctx context.Context, traceID string) context.Context {
	if traceID == "" {
		traceID = uuid.New().String()
	}
	return context.WithValue(ctx, "trace_id", traceID)
}

// GetTraceID コンテキストからトレースIDを取得
func GetTraceID(ctx context.Context) string {
	if traceID := ctx.Value("trace_id"); traceID != nil {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	return ""
}

// WithHTTPRequest HTTPリクエスト情報をコンテキストに追加
func WithHTTPRequest(ctx context.Context, httpInfo *HTTPRequestInfo) context.Context {
	return context.WithValue(ctx, "http_request", httpInfo)
}

