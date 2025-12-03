package service

import (
	"context"
	"fmt"

	"tailor-cloud/backend/internal/repository"
)

// MeasurementCorrectionService 自動補正エンジン（The "Auto Patterner"）
// 特許出願済みロジックに基づく採寸データの自動補正
type MeasurementCorrectionService struct {
	diagnosisService *DiagnosisService
	fabricRepo       repository.FabricRepository
}

// NewMeasurementCorrectionService MeasurementCorrectionServiceのコンストラクタ
func NewMeasurementCorrectionService(
	diagnosisService *DiagnosisService,
	fabricRepo repository.FabricRepository,
) *MeasurementCorrectionService {
	return &MeasurementCorrectionService{
		diagnosisService: diagnosisService,
		fabricRepo:       fabricRepo,
	}
}

// RawMeasurement ヌード寸（身体の実寸）
type RawMeasurement struct {
	Height float64 `json:"height"` // 身長 (cm)
	Bust   float64 `json:"bust"`   // バスト (cm)
	Waist  float64 `json:"waist"`  // ウエスト (cm)
	Hip    float64 `json:"hip"`    // ヒップ (cm)
	Thigh  float64 `json:"thigh"`  // 太もも (cm)
	Knee   float64 `json:"knee"`   // 膝 (cm)
	Calf   float64 `json:"calf"`   // ふくらはぎ (cm)
	OB     float64 `json:"ob"`     // OB (Over Bust) - バスト上
}

// DiagnosisProfile 診断プロファイル
type DiagnosisProfile struct {
	Archetype     string  `json:"archetype"`      // "Classic", "Modern", "Elegant", etc.
	FitPreference string  `json:"fit_preference"` // "tight", "relaxed"
	Silhouette    string  `json:"silhouette"`     // "tapered", "straight"
	Ease          float64 `json:"ease"`           // ゆとり量 (cm)
}

// Correction 補正情報
type Correction struct {
	Type        string  `json:"type"`        // 補正タイプ
	Value       float64 `json:"value"`       // 補正前の値
	Adjustment  float64 `json:"adjustment"`  // 補正量
	Description string  `json:"description"` // 補正説明
}

// FinalMeasurement 仕上がり寸法（製造指示値）
type FinalMeasurement struct {
	JacketLength float64      `json:"jacket_length"` // ジャケット長
	SleeveLength float64      `json:"sleeve_length"` // 袖長
	Chest        float64      `json:"chest"`         // 胸囲
	Waist        float64      `json:"waist"`         // ウエスト
	Hip          float64      `json:"hip"`           // ヒップ
	Thigh        float64      `json:"thigh"`         // 太もも
	Knee         float64      `json:"knee"`          // 膝
	Calf         float64      `json:"calf"`          // ふくらはぎ
	Hem          float64      `json:"hem"`           // 裾幅
	Corrections  []Correction `json:"corrections"`   // 適用された補正の履歴
}

// ConvertToFinalMeasurementsRequest 変換リクエスト
type ConvertToFinalMeasurementsRequest struct {
	RawMeasurements *RawMeasurement `json:"raw_measurements"`
	UserID          string           `json:"user_id"`
	TenantID        string           `json:"tenant_id"`
	FabricID        string           `json:"fabric_id"`
}

// ConvertToFinalMeasurementsResponse 変換レスポンス
type ConvertToFinalMeasurementsResponse struct {
	FinalMeasurements *FinalMeasurement `json:"final_measurements"`
}

// ConvertToFinalMeasurements ヌード寸を仕上がり寸法に変換
// 特許出願済みロジックに基づく自動補正を適用
func (s *MeasurementCorrectionService) ConvertToFinalMeasurements(
	ctx context.Context,
	req *ConvertToFinalMeasurementsRequest,
) (*ConvertToFinalMeasurementsResponse, error) {
	// 1. 診断プロファイルを取得（オプション）
	var profile *DiagnosisProfile
	var err error
	if req.UserID != "" {
		profile, err = s.getDiagnosisProfile(ctx, req.UserID, req.TenantID)
		if err != nil {
			// 診断プロファイルがない場合はデフォルト値を使用
			profile = &DiagnosisProfile{
				Archetype:     "Classic",
				FitPreference: "relaxed",
				Silhouette:    "tapered",
				Ease:          3.0, // デフォルトゆとり量
			}
		}
	} else {
		// UserIDがない場合はデフォルト値を使用
		profile = &DiagnosisProfile{
			Archetype:     "Classic",
			FitPreference: "relaxed",
			Silhouette:    "tapered",
			Ease:          3.0, // デフォルトゆとり量
		}
	}

	// 2. 生地の特性を取得
	fabric, err := s.fabricRepo.GetByID(ctx, req.FabricID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fabric: %w", err)
	}

	// 3. 初期値を設定（ヌード寸をベースに）
	final := &FinalMeasurement{
		Chest:       req.RawMeasurements.Bust,
		Waist:       req.RawMeasurements.Waist,
		Hip:         req.RawMeasurements.Hip,
		Thigh:       req.RawMeasurements.Thigh,
		Knee:        req.RawMeasurements.Knee,
		Calf:        req.RawMeasurements.Calf,
		Corrections: []Correction{},
	}

	// 4. OB差分補正を適用
	if err := s.applyOBDifferenceCorrection(req.RawMeasurements, final); err != nil {
		return nil, fmt.Errorf("OB difference correction failed: %w", err)
	}

	// 5. シルエット計算
	if err := s.calculateSilhouette(req.RawMeasurements, profile, final); err != nil {
		return nil, fmt.Errorf("silhouette calculation failed: %w", err)
	}

	// 6. 生地の伸縮性を考慮した補正（将来実装）
	// if err := s.applyFabricStretchCorrection(final, fabric); err != nil {
	// 	return nil, fmt.Errorf("fabric stretch correction failed: %w", err)
	// }

	// 7. バリデーション
	if err := s.validateMeasurements(final, req.RawMeasurements); err != nil {
		return nil, fmt.Errorf("measurement validation failed: %w", err)
	}

	// fabric を使用して未使用変数エラーを回避
	_ = fabric

	return &ConvertToFinalMeasurementsResponse{
		FinalMeasurements: final,
	}, nil
}

// applyOBDifferenceCorrection OB差分補正を適用
// ロジック: IF (OB - Bust >= 20) THEN Add Correction
func (s *MeasurementCorrectionService) applyOBDifferenceCorrection(
	raw *RawMeasurement,
	final *FinalMeasurement,
) error {
	obDiff := raw.OB - raw.Bust

	if obDiff >= 20.0 {
		// OB差分が20cm以上の場合は補正を適用
		// 補正係数: 0.5（要調整・特許ロジックに基づく）
		adjustment := obDiff * 0.5

		correction := Correction{
			Type:        "OB_DIFFERENCE",
			Value:       obDiff,
			Adjustment:  adjustment,
			Description: fmt.Sprintf("OB差分補正: %.1fcm (調整量: %.1fcm)", obDiff, adjustment),
		}

		// 胸囲に補正を適用
		final.Chest += adjustment
		final.Corrections = append(final.Corrections, correction)
	}

	return nil
}

// calculateSilhouette シルエット計算
// ロジック: Knee/2 + Ease - 5.0cm (Tapered)
func (s *MeasurementCorrectionService) calculateSilhouette(
	raw *RawMeasurement,
	profile *DiagnosisProfile,
	final *FinalMeasurement,
) error {
	if profile.Silhouette == "tapered" {
		// テーパード: Knee/2 + Ease - 5.0cm
		kneeHalf := raw.Knee / 2.0
		hem := kneeHalf + profile.Ease - 5.0

		// リミッター: IF Hem < (Calf/2 - 1.5) THEN Error
		minHem := (raw.Calf / 2.0) - 1.5
		if hem < minHem {
			return fmt.Errorf(
				"hem width (%.1fcm) is less than minimum (%.1fcm). Silhouette calculation failed",
				hem, minHem,
			)
		}

		final.Hem = hem
		final.Corrections = append(final.Corrections, Correction{
			Type:        "SILHOUETTE_TAPERED",
			Value:       hem,
			Description: fmt.Sprintf("テーパードシルエット: 裾幅 %.1fcm (膝幅/2 + ゆとり%.1fcm - 5.0cm)", hem, profile.Ease),
		})
	} else if profile.Silhouette == "straight" {
		// ストレート: 膝幅と同程度
		final.Hem = raw.Knee / 2.0
		final.Corrections = append(final.Corrections, Correction{
			Type:        "SILHOUETTE_STRAIGHT",
			Value:       final.Hem,
			Description: fmt.Sprintf("ストレートシルエット: 裾幅 %.1fcm (膝幅/2)", final.Hem),
		})
	} else {
		// デフォルト: テーパード
		kneeHalf := raw.Knee / 2.0
		hem := kneeHalf + profile.Ease - 5.0

		// リミッター: IF Hem < (Calf/2 - 1.5) THEN Error
		minHem := (raw.Calf / 2.0) - 1.5
		if hem < minHem {
			return fmt.Errorf(
				"hem width (%.1fcm) is less than minimum (%.1fcm). Silhouette calculation failed",
				hem, minHem,
			)
		}

		final.Hem = hem
		final.Corrections = append(final.Corrections, Correction{
			Type:        "SILHOUETTE_DEFAULT",
			Value:       hem,
			Description: fmt.Sprintf("デフォルトシルエット（テーパード）: 裾幅 %.1fcm (膝幅/2 + ゆとり%.1fcm - 5.0cm)", hem, profile.Ease),
		})
	}

	return nil
}

// validateMeasurements バリデーション
// リミッター: IF Hem < (Calf/2 - 1.5) THEN Error
func (s *MeasurementCorrectionService) validateMeasurements(
	final *FinalMeasurement,
	raw *RawMeasurement,
) error {
	// リミッター: IF Hem < (Calf/2 - 1.5) THEN Error
	if final.Hem > 0 {
		minHem := (raw.Calf / 2.0) - 1.5
		if final.Hem < minHem {
			return fmt.Errorf(
				"hem width validation failed: %.1fcm < %.1fcm (minimum: Calf/2 - 1.5cm)",
				final.Hem, minHem,
			)
		}
	}

	// その他のバリデーション
	if final.Chest <= 0 || final.Waist <= 0 {
		return fmt.Errorf("invalid measurement values: chest=%.1f, waist=%.1f", final.Chest, final.Waist)
	}

	return nil
}

// getDiagnosisProfile 診断プロファイルを取得
func (s *MeasurementCorrectionService) getDiagnosisProfile(
	ctx context.Context,
	userID string,
	tenantID string,
) (*DiagnosisProfile, error) {
	if s.diagnosisService == nil {
		return nil, fmt.Errorf("diagnosis service is not available")
	}

	// 最新の診断結果を取得
	diagnosis, err := s.diagnosisService.GetLatestByUserID(ctx, userID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get diagnosis: %w", err)
	}

	// アーキタイプからゆとり量を決定
	ease := s.getEaseByArchetype(string(diagnosis.Archetype))

	// シルエットを決定（診断結果から、またはデフォルト）
	silhouette := "tapered" // デフォルト
	if diagnosis.PlanType == "Best Value" {
		silhouette = "tapered"
	}

	return &DiagnosisProfile{
		Archetype:     string(diagnosis.Archetype),
		FitPreference: string(diagnosis.PlanType),
		Silhouette:    silhouette,
		Ease:          ease,
	}, nil
}

// getEaseByArchetype アーキタイプからゆとり量を決定
func (s *MeasurementCorrectionService) getEaseByArchetype(archetype string) float64 {
	// アーキタイプに応じたゆとり量（要調整・特許ロジックに基づく）
	switch archetype {
	case "Classic":
		return 3.0
	case "Modern":
		return 2.0
	case "Elegant":
		return 4.0
	case "Sporty":
		return 1.5
	case "Casual":
		return 2.5
	default:
		return 3.0 // デフォルト
	}
}

