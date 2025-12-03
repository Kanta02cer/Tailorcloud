# 自動補正エンジン（The "Auto Patterner"）実装計画

**作成日**: 2025-01  
**優先度**: 🔴 最優先  
**工数見積**: 2-3週間  
**目的**: システムの核心価値「感性×製造データの自動変換」を実現

---

## 🎯 概要

開発計画書に記載された特許出願済みロジックに基づく自動補正エンジンを実装します。

**核心ロジック**:
1. **OB差分補正**: `IF (OB - Bust >= 20) THEN Add Correction`
2. **シルエット計算**: `Knee/2 + Ease - 5.0cm (Tapered)`
3. **リミッター**: `IF Hem < (Calf/2 - 1.5) THEN Error`

---

## 📋 実装要件

### 入力データ

```go
type RawMeasurement struct {
    // ヌード寸（身体の実寸）
    Height      float64 // 身長 (cm)
    Bust        float64 // バスト (cm)
    Waist       float64 // ウエスト (cm)
    Hip         float64 // ヒップ (cm)
    Thigh       float64 // 太もも (cm)
    Knee        float64 // 膝 (cm)
    Calf        float64 // ふくらはぎ (cm)
    OB          float64 // OB (Over Bust) - バスト上
    // ... その他の寸法
}

type DiagnosisProfile struct {
    Archetype    string // "Classic", "Modern", "Elegant", etc.
    FitPreference string // "tight", "relaxed"
    Silhouette   string // "tapered", "straight"
    Ease         float64 // ゆとり量 (cm)
}
```

### 出力データ

```go
type FinalMeasurement struct {
    // 仕上がり寸法（製造指示値）
    JacketLength float64 // ジャケット長
    SleeveLength float64 // 袖長
    Chest        float64 // 胸囲
    Waist        float64 // ウエスト
    Hip          float64 // ヒップ
    Thigh        float64 // 太もも
    Knee         float64 // 膝
    Calf         float64 // ふくらはぎ
    Hem          float64 // 裾幅
    // ... その他の寸法
    Corrections  []Correction // 適用された補正の履歴
}
```

---

## 🔧 実装内容

### 1. 補正ロジックサービスの作成

**ファイル**: `internal/service/measurement_correction_service.go` (新規作成)

#### 1.1 OB差分補正

```go
func (s *MeasurementCorrectionService) ApplyOBDifferenceCorrection(
    raw *RawMeasurement,
    final *FinalMeasurement,
) error {
    obDiff := raw.OB - raw.Bust
    
    if obDiff >= 20.0 {
        // OB差分が20cm以上の場合は補正を適用
        correction := Correction{
            Type:        "OB_DIFFERENCE",
            Value:       obDiff,
            Adjustment:  obDiff * 0.5, // 補正係数（要調整）
            Description: fmt.Sprintf("OB差分補正: %.1fcm", obDiff),
        }
        
        // 胸囲に補正を適用
        final.Chest += correction.Adjustment
        final.Corrections = append(final.Corrections, correction)
    }
    
    return nil
}
```

#### 1.2 シルエット計算

```go
func (s *MeasurementCorrectionService) CalculateSilhouette(
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
            return fmt.Errorf("hem width (%.1fcm) is less than minimum (%.1fcm)", hem, minHem)
        }
        
        final.Hem = hem
        final.Corrections = append(final.Corrections, Correction{
            Type:        "SILHOUETTE_TAPERED",
            Value:       hem,
            Description: fmt.Sprintf("テーパードシルエット: 裾幅 %.1fcm", hem),
        })
    } else if profile.Silhouette == "straight" {
        // ストレート: 膝幅と同程度
        final.Hem = raw.Knee / 2.0
    }
    
    return nil
}
```

#### 1.3 リミッター（バリデーション）

```go
func (s *MeasurementCorrectionService) ValidateMeasurements(
    final *FinalMeasurement,
    raw *RawMeasurement,
) error {
    // リミッター: IF Hem < (Calf/2 - 1.5) THEN Error
    minHem := (raw.Calf / 2.0) - 1.5
    if final.Hem < minHem {
        return fmt.Errorf(
            "hem width validation failed: %.1fcm < %.1fcm (minimum)",
            final.Hem, minHem,
        )
    }
    
    // その他のバリデーション
    if final.Chest <= 0 || final.Waist <= 0 {
        return fmt.Errorf("invalid measurement values")
    }
    
    return nil
}
```

### 2. 診断プロファイル連携

**ファイル**: `internal/service/diagnosis_service.go` (既存を拡張)

```go
func (s *DiagnosisService) GetDiagnosisProfile(
    ctx context.Context,
    userID string,
    tenantID string,
) (*DiagnosisProfile, error) {
    // 最新の診断結果を取得
    diagnosis, err := s.diagnosisRepo.GetLatestByUserID(ctx, userID, tenantID)
    if err != nil {
        return nil, fmt.Errorf("failed to get diagnosis: %w", err)
    }
    
    // アーキタイプからゆとり量を決定
    ease := s.getEaseByArchetype(diagnosis.Archetype)
    
    return &DiagnosisProfile{
        Archetype:    diagnosis.Archetype,
        FitPreference: diagnosis.PlanType, // "Best Value" → "relaxed"
        Silhouette:   diagnosis.Silhouette, // 診断結果から取得
        Ease:         ease,
    }, nil
}
```

### 3. メイン変換関数

```go
func (s *MeasurementCorrectionService) ConvertToFinalMeasurements(
    ctx context.Context,
    raw *RawMeasurement,
    userID string,
    tenantID string,
    fabricID string,
) (*FinalMeasurement, error) {
    // 1. 診断プロファイルを取得
    profile, err := s.diagnosisService.GetDiagnosisProfile(ctx, userID, tenantID)
    if err != nil {
        return nil, fmt.Errorf("failed to get diagnosis profile: %w", err)
    }
    
    // 2. 生地の特性を取得
    fabric, err := s.fabricRepo.GetByID(ctx, fabricID, tenantID)
    if err != nil {
        return nil, fmt.Errorf("failed to get fabric: %w", err)
    }
    
    // 3. 初期値を設定（ヌード寸をベースに）
    final := &FinalMeasurement{
        Chest: raw.Bust,
        Waist: raw.Waist,
        Hip:   raw.Hip,
        Thigh: raw.Thigh,
        Knee:  raw.Knee,
        Calf:  raw.Calf,
    }
    
    // 4. OB差分補正を適用
    if err := s.ApplyOBDifferenceCorrection(raw, final); err != nil {
        return nil, fmt.Errorf("OB difference correction failed: %w", err)
    }
    
    // 5. シルエット計算
    if err := s.CalculateSilhouette(raw, profile, final); err != nil {
        return nil, fmt.Errorf("silhouette calculation failed: %w", err)
    }
    
    // 6. 生地の伸縮性を考慮した補正
    if err := s.ApplyFabricStretchCorrection(final, fabric); err != nil {
        return nil, fmt.Errorf("fabric stretch correction failed: %w", err)
    }
    
    // 7. バリデーション
    if err := s.ValidateMeasurements(final, raw); err != nil {
        return nil, fmt.Errorf("measurement validation failed: %w", err)
    }
    
    return final, nil
}
```

---

## 📡 APIエンドポイント

### POST /api/measurements/convert

**リクエスト**:
```json
{
  "raw_measurements": {
    "height": 170.0,
    "bust": 90.0,
    "waist": 75.0,
    "hip": 95.0,
    "thigh": 55.0,
    "knee": 40.0,
    "calf": 35.0,
    "ob": 110.0
  },
  "user_id": "user_123",
  "fabric_id": "fabric_456"
}
```

**レスポンス**:
```json
{
  "final_measurements": {
    "chest": 95.0,
    "waist": 75.0,
    "hip": 95.0,
    "thigh": 55.0,
    "knee": 40.0,
    "calf": 35.0,
    "hem": 20.0
  },
  "corrections": [
    {
      "type": "OB_DIFFERENCE",
      "value": 20.0,
      "adjustment": 10.0,
      "description": "OB差分補正: 20.0cm"
    },
    {
      "type": "SILHOUETTE_TAPERED",
      "value": 20.0,
      "description": "テーパードシルエット: 裾幅 20.0cm"
    }
  ]
}
```

---

## 🧪 テスト計画

### 単体テスト

1. **OB差分補正テスト**
   - OB - Bust = 20cm以上の場合、補正が適用される
   - OB - Bust < 20cmの場合、補正が適用されない

2. **シルエット計算テスト**
   - テーパード: Knee/2 + Ease - 5.0cm
   - ストレート: Knee/2

3. **リミッターテスト**
   - Hem < (Calf/2 - 1.5) の場合、エラーが返る
   - Hem >= (Calf/2 - 1.5) の場合、正常に処理される

---

## 📅 実装スケジュール

### Week 1: コアロジック実装
- [ ] 補正ロジックサービスの作成
- [ ] OB差分補正の実装
- [ ] シルエット計算の実装
- [ ] リミッターの実装

### Week 2: 統合・API実装
- [ ] 診断プロファイル連携
- [ ] 生地特性連携
- [ ] APIエンドポイント実装
- [ ] 単体テスト

### Week 3: UI実装・テスト
- [ ] ヌード寸入力フォーム拡張
- [ ] 仕上がり寸法表示UI
- [ ] 統合テスト
- [ ] ドキュメント作成

---

**最終更新日**: 2025-01

