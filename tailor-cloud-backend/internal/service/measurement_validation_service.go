package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"tailor-cloud/backend/internal/repository"
)

// MeasurementValidationService 採寸データバリデーションサービス
// 前回採寸データとの比較、異常値検出を担当
type MeasurementValidationService struct {
	orderRepo repository.OrderRepository
}

// NewMeasurementValidationService MeasurementValidationServiceのコンストラクタ
func NewMeasurementValidationService(orderRepo repository.OrderRepository) *MeasurementValidationService {
	return &MeasurementValidationService{
		orderRepo: orderRepo,
	}
}

// MeasurementData 採寸データ（JSON形式からパース）
type MeasurementData struct {
	Height      *float64 `json:"height"`      // 身長
	Bust        *float64 `json:"bust"`        // バスト
	Waist       *float64 `json:"waist"`        // ウエスト
	Hip         *float64 `json:"hip"`         // ヒップ
	Thigh       *float64 `json:"thigh"`       // 太もも
	Knee        *float64 `json:"knee"`         // 膝
	Calf        *float64 `json:"calf"`         // ふくらはぎ
	OB          *float64 `json:"ob"`           // OB
	JacketLength *float64 `json:"jacket_length"` // ジャケット長
	Sleeve      *float64 `json:"sleeve"`       // 袖長
	Chest       *float64 `json:"chest"`        // 胸囲
}

// ValidationAlert バリデーションアラート
type ValidationAlert struct {
	Field       string  `json:"field"`        // フィールド名
	Current     float64 `json:"current"`       // 現在の値
	Previous    float64 `json:"previous"`     // 前回の値
	Difference  float64 `json:"difference"`   // 差分（絶対値）
	Threshold   float64 `json:"threshold"`    // 閾値（5.0cm）
	Severity    string  `json:"severity"`     // "warning" or "error"
	Message     string  `json:"message"`      // アラートメッセージ
}

// ValidateMeasurementsRequest バリデーションリクエスト
type ValidateMeasurementsRequest struct {
	CustomerID      string                   `json:"customer_id"`
	TenantID        string                   `json:"tenant_id"`
	CurrentMeasurements json.RawMessage     `json:"current_measurements"`
}

// ValidateMeasurementsResponse バリデーションレスポンス
type ValidateMeasurementsResponse struct {
	IsValid      bool              `json:"is_valid"`      // バリデーション成功
	Alerts       []ValidationAlert `json:"alerts"`        // アラート一覧
	HasWarnings  bool              `json:"has_warnings"`  // 警告があるか
	HasErrors    bool              `json:"has_errors"`    // エラーがあるか
	PreviousData *MeasurementData  `json:"previous_data"` // 前回の採寸データ（存在する場合）
}

// ValidateMeasurements 採寸データをバリデーション
// 前回採寸データとの比較、異常値検出を実行
func (s *MeasurementValidationService) ValidateMeasurements(
	ctx context.Context,
	req *ValidateMeasurementsRequest,
) (*ValidateMeasurementsResponse, error) {
	// 1. 現在の採寸データをパース
	var currentData MeasurementData
	if err := json.Unmarshal(req.CurrentMeasurements, &currentData); err != nil {
		return nil, fmt.Errorf("failed to parse current measurements: %w", err)
	}

	// 2. 顧客の注文履歴から前回の採寸データを取得
	previousData, err := s.getPreviousMeasurements(ctx, req.CustomerID, req.TenantID)
	if err != nil {
		// 前回データがない場合は警告のみ（エラーではない）
		return &ValidateMeasurementsResponse{
			IsValid:     true,
			Alerts:      []ValidationAlert{},
			HasWarnings: false,
			HasErrors:   false,
			PreviousData: nil,
		}, nil
	}

	// 3. 前回データと比較してアラートを生成
	alerts := s.compareMeasurements(currentData, *previousData)

	// 4. レスポンスを構築
	hasErrors := false
	hasWarnings := false
	for _, alert := range alerts {
		if alert.Severity == "error" {
			hasErrors = true
		} else if alert.Severity == "warning" {
			hasWarnings = true
		}
	}

	return &ValidateMeasurementsResponse{
		IsValid:      !hasErrors, // エラーがない場合は有効
		Alerts:       alerts,
		HasWarnings:  hasWarnings,
		HasErrors:    hasErrors,
		PreviousData: previousData,
	}, nil
}

// getPreviousMeasurements 顧客の前回採寸データを取得
func (s *MeasurementValidationService) getPreviousMeasurements(
	ctx context.Context,
	customerID string,
	tenantID string,
) (*MeasurementData, error) {
	// 顧客の注文履歴を取得（最新の注文から順に）
	orders, err := s.orderRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	// 顧客の注文をフィルターして、採寸データがある最新の注文を探す
	for _, order := range orders {
		if order.CustomerID == customerID && order.Details != nil && order.Details.MeasurementData != nil {
			var data MeasurementData
			if err := json.Unmarshal(order.Details.MeasurementData, &data); err == nil {
				// 有効な採寸データが見つかった
				return &data, nil
			}
		}
	}

	return nil, fmt.Errorf("no previous measurements found")
}

// compareMeasurements 採寸データを比較してアラートを生成
// 閾値: ±5cm以上の差分で警告
func (s *MeasurementValidationService) compareMeasurements(
	current MeasurementData,
	previous MeasurementData,
) []ValidationAlert {
	var alerts []ValidationAlert
	threshold := 5.0 // 閾値: 5cm

	// 比較するフィールドのリスト
	fields := []struct {
		name     string
		current  *float64
		previous *float64
	}{
		{"height", current.Height, previous.Height},
		{"bust", current.Bust, previous.Bust},
		{"waist", current.Waist, previous.Waist},
		{"hip", current.Hip, previous.Hip},
		{"thigh", current.Thigh, previous.Thigh},
		{"knee", current.Knee, previous.Knee},
		{"calf", current.Calf, previous.Calf},
		{"ob", current.OB, previous.OB},
		{"jacket_length", current.JacketLength, previous.JacketLength},
		{"sleeve", current.Sleeve, previous.Sleeve},
		{"chest", current.Chest, previous.Chest},
	}

	for _, field := range fields {
		// 両方の値が存在する場合のみ比較
		if field.current == nil || field.previous == nil {
			continue
		}

		currentValue := *field.current
		previousValue := *field.previous
		difference := math.Abs(currentValue - previousValue)

		// 閾値（5cm）を超えている場合
		if difference >= threshold {
			severity := "warning"
			// 10cm以上の差分はエラー
			if difference >= 10.0 {
				severity = "error"
			}

			alert := ValidationAlert{
				Field:      field.name,
				Current:    currentValue,
				Previous:   previousValue,
				Difference: difference,
				Threshold:  threshold,
				Severity:   severity,
				Message:    fmt.Sprintf(
					"%sの値が前回より%.1fcm異なります（現在: %.1fcm, 前回: %.1fcm）",
					field.name, difference, currentValue, previousValue,
				),
			}
			alerts = append(alerts, alert)
		}
	}

	return alerts
}

// ValidateMeasurementRange 採寸データの範囲をバリデーション
// 異常値（例: 身長が50cm以下、300cm以上など）を検出
func (s *MeasurementValidationService) ValidateMeasurementRange(
	measurements json.RawMessage,
) []ValidationAlert {
	var data MeasurementData
	if err := json.Unmarshal(measurements, &data); err != nil {
		return []ValidationAlert{{
			Field:    "parse_error",
			Severity: "error",
			Message:  fmt.Sprintf("Failed to parse measurements: %v", err),
		}}
	}

	var alerts []ValidationAlert

	// 身長の範囲チェック（50cm - 250cm）
	if data.Height != nil {
		height := *data.Height
		if height < 50.0 || height > 250.0 {
			alerts = append(alerts, ValidationAlert{
				Field:    "height",
				Current:  height,
				Severity: "error",
				Message:  fmt.Sprintf("身長の値が異常です: %.1fcm（正常範囲: 50-250cm）", height),
			})
		}
	}

	// バストの範囲チェック（50cm - 200cm）
	if data.Bust != nil {
		bust := *data.Bust
		if bust < 50.0 || bust > 200.0 {
			alerts = append(alerts, ValidationAlert{
				Field:    "bust",
				Current:  bust,
				Severity: "error",
				Message:  fmt.Sprintf("バストの値が異常です: %.1fcm（正常範囲: 50-200cm）", bust),
			})
		}
	}

	// ウエストの範囲チェック（40cm - 150cm）
	if data.Waist != nil {
		waist := *data.Waist
		if waist < 40.0 || waist > 150.0 {
			alerts = append(alerts, ValidationAlert{
				Field:    "waist",
				Current:  waist,
				Severity: "error",
				Message:  fmt.Sprintf("ウエストの値が異常です: %.1fcm（正常範囲: 40-150cm）", waist),
			})
		}
	}

	// その他のフィールドも同様にチェック（必要に応じて追加）

	return alerts
}

