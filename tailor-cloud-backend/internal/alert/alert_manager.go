package alert

import (
	"fmt"
	"time"

	"tailor-cloud/backend/internal/metrics"
)

// AlertLevel アラートレベル
type AlertLevel string

const (
	AlertLevelInfo    AlertLevel = "INFO"
	AlertLevelWarning AlertLevel = "WARNING"
	AlertLevelCritical AlertLevel = "CRITICAL"
)

// AlertRule アラートルール
type AlertRule struct {
	Name        string
	Condition   func(*metrics.Metrics) bool
	Level       AlertLevel
	Description string
}

// Alert アラート情報
type Alert struct {
	Rule        *AlertRule
	Metrics     *metrics.Metrics
	TriggeredAt time.Time
	Message     string
}

// AlertHandler アラートハンドラーインターフェース
type AlertHandler interface {
	Handle(alert *Alert) error
}

// AlertManager アラートマネージャー
type AlertManager struct {
	rules   []*AlertRule
	handlers []AlertHandler
}

// NewAlertManager AlertManagerを作成
func NewAlertManager() *AlertManager {
	return &AlertManager{
		rules:    make([]*AlertRule, 0),
		handlers: make([]AlertHandler, 0),
	}
}

// AddRule アラートルールを追加
func (m *AlertManager) AddRule(rule *AlertRule) {
	m.rules = append(m.rules, rule)
}

// AddHandler アラートハンドラーを追加
func (m *AlertManager) AddHandler(handler AlertHandler) {
	m.handlers = append(m.handlers, handler)
}

// Check メトリクスをチェックしてアラートを発火
func (m *AlertManager) Check(metrics *metrics.Metrics) []*Alert {
	alerts := make([]*Alert, 0)

	for _, rule := range m.rules {
		if rule.Condition(metrics) {
			alert := &Alert{
				Rule:        rule,
				Metrics:     metrics,
				TriggeredAt: time.Now(),
				Message:     fmt.Sprintf("[%s] %s", rule.Level, rule.Description),
			}

			alerts = append(alerts, alert)

			// アラートハンドラーを実行
			for _, handler := range m.handlers {
				if err := handler.Handle(alert); err != nil {
					// ハンドラーエラーはログに記録（非同期で実行されるため）
					fmt.Printf("WARNING: Failed to handle alert: %v\n", err)
				}
			}
		}
	}

	return alerts
}

// DefaultAlertRules デフォルトのアラートルールを作成
func DefaultAlertRules() []*AlertRule {
	return []*AlertRule{
		// エラー率が5%を超えた場合
		{
			Name:        "HighErrorRate",
			Condition:   func(m *metrics.Metrics) bool { return m.ErrorRate > 5.0 },
			Level:       AlertLevelWarning,
			Description: fmt.Sprintf("Error rate is high: %.2f%%", 5.0),
		},
		// エラー率が10%を超えた場合
		{
			Name:        "CriticalErrorRate",
			Condition:   func(m *metrics.Metrics) bool { return m.ErrorRate > 10.0 },
			Level:       AlertLevelCritical,
			Description: fmt.Sprintf("Error rate is critical: %.2f%%", 10.0),
		},
		// 平均レイテンシーが1秒を超えた場合
		{
			Name:        "HighLatency",
			Condition:   func(m *metrics.Metrics) bool { return m.AverageLatency > 1*time.Second },
			Level:       AlertLevelWarning,
			Description: "Average latency is high",
		},
		// 平均レイテンシーが5秒を超えた場合
		{
			Name:        "CriticalLatency",
			Condition:   func(m *metrics.Metrics) bool { return m.AverageLatency > 5*time.Second },
			Level:       AlertLevelCritical,
			Description: "Average latency is critical",
		},
		// DB接続数が最大の80%を超えた場合
		{
			Name: "HighDBConnections",
			Condition: func(m *metrics.Metrics) bool {
				maxConns := int64(25) // DefaultPoolConfigのMaxOpenConns
				return m.DBConnections > int64(float64(maxConns)*0.8)
			},
			Level:       AlertLevelWarning,
			Description: "Database connections are high",
		},
	}
}

// LogAlertHandler ログ出力アラートハンドラー
type LogAlertHandler struct{}

// Handle アラートをログに出力
func (h *LogAlertHandler) Handle(alert *Alert) error {
	fmt.Printf("[ALERT] %s: %s\n", alert.Rule.Level, alert.Message)
	return nil
}

