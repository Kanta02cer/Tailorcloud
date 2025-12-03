package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Diagnosis 診断モデル
// Suit-MBTI診断結果を記録
type Diagnosis struct {
	ID              string          `json:"id" db:"id"`
	UserID          string          `json:"user_id" db:"user_id"`
	TenantID        string          `json:"tenant_id" db:"tenant_id"`
	Archetype       Archetype       `json:"archetype" db:"archetype"`
	PlanType        PlanType        `json:"plan_type" db:"plan_type"`
	DiagnosisResult json.RawMessage `json:"diagnosis_result" db:"diagnosis_result"` // 診断結果詳細（JSON形式）
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

// Archetype アーキタイプ（RATタイプ）
type Archetype string

const (
	ArchetypeClassic Archetype = "Classic" // クラシック
	ArchetypeModern  Archetype = "Modern"  // モダン
	ArchetypeElegant Archetype = "Elegant" // エレガント
	ArchetypeSporty  Archetype = "Sporty"  // スポーティ
	ArchetypeCasual  Archetype = "Casual"  // カジュアル
)

// PlanType プランタイプ
type PlanType string

const (
	PlanTypeBestValue PlanType = "Best Value" // ベストバリュー
	PlanTypeAuthentic PlanType = "Authentic"  // オーセンティック
)

// IsValid アーキタイプが有効かチェック
func (a Archetype) IsValid() bool {
	switch a {
	case ArchetypeClassic, ArchetypeModern, ArchetypeElegant, ArchetypeSporty, ArchetypeCasual:
		return true
	default:
		return false
	}
}

// IsValid プランタイプが有効かチェック
func (p PlanType) IsValid() bool {
	switch p {
	case PlanTypeBestValue, PlanTypeAuthentic:
		return true
	default:
		return false
	}
}

// NewDiagnosis 新しい診断を作成
func NewDiagnosis(userID, tenantID string, archetype Archetype, planType PlanType, diagnosisResult json.RawMessage) *Diagnosis {
	now := time.Now()
	return &Diagnosis{
		ID:              uuid.New().String(),
		UserID:          userID,
		TenantID:        tenantID,
		Archetype:       archetype,
		PlanType:        planType,
		DiagnosisResult: diagnosisResult,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

