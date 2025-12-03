-- ============================================================================
-- TailorCloud: Suit-MBTI統合 - 診断ログテーブル作成
-- ============================================================================
-- 目的: Suit-MBTI診断結果をTailorCloudに統合
-- ============================================================================

-- Diagnoses (診断ログ) テーブル
CREATE TABLE IF NOT EXISTS diagnoses (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(255) NOT NULL,
    archetype VARCHAR(50) NOT NULL, -- RATタイプ（例: "Classic", "Modern", etc.）
    plan_type VARCHAR(50), -- プラン（"Best Value" / "Authentic"）
    diagnosis_result JSONB, -- 診断結果詳細（JSON形式）
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
    -- 注意: tenantsテーブルが存在する場合のみ外部キー制約を有効化
    -- FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT diagnoses_archetype_check CHECK (archetype IN ('Classic', 'Modern', 'Elegant', 'Sporty', 'Casual')),
    CONSTRAINT diagnoses_plan_type_check CHECK (plan_type IN ('Best Value', 'Authentic') OR plan_type IS NULL)
);

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_diagnoses_user_id ON diagnoses(user_id);
CREATE INDEX IF NOT EXISTS idx_diagnoses_tenant_id ON diagnoses(tenant_id);
CREATE INDEX IF NOT EXISTS idx_diagnoses_archetype ON diagnoses(archetype);
CREATE INDEX IF NOT EXISTS idx_diagnoses_tenant_user ON diagnoses(tenant_id, user_id);
CREATE INDEX IF NOT EXISTS idx_diagnoses_created_at ON diagnoses(created_at DESC);

-- コメント追加
COMMENT ON TABLE diagnoses IS 'Suit-MBTI診断ログテーブル: 診断結果を記録';
COMMENT ON COLUMN diagnoses.archetype IS '選択アーキタイプ（RATタイプ）';
COMMENT ON COLUMN diagnoses.plan_type IS 'プラン（Best Value / Authentic）';
COMMENT ON COLUMN diagnoses.diagnosis_result IS '診断結果詳細（JSON形式）';

