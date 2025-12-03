-- ============================================================================
-- TailorCloud: Suit-MBTI統合 - 顧客テーブル拡張
-- ============================================================================
-- 目的: 顧客テーブルに診断関連フィールドを追加
-- ============================================================================

-- Customers テーブルに新しいカラムを追加
ALTER TABLE customers 
ADD COLUMN IF NOT EXISTS occupation VARCHAR(255),
ADD COLUMN IF NOT EXISTS annual_income_range VARCHAR(50),
ADD COLUMN IF NOT EXISTS ltv_score DECIMAL(10,2) DEFAULT 0,
ADD COLUMN IF NOT EXISTS preferred_archetype VARCHAR(50),
ADD COLUMN IF NOT EXISTS diagnosis_count INTEGER DEFAULT 0;

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_customers_archetype ON customers(preferred_archetype);
CREATE INDEX IF NOT EXISTS idx_customers_ltv_score ON customers(ltv_score DESC);

-- コメント追加
COMMENT ON COLUMN customers.occupation IS '職業';
COMMENT ON COLUMN customers.annual_income_range IS '年収感（例: "500-1000万円"）';
COMMENT ON COLUMN customers.ltv_score IS 'LTVスコア（生涯価値、計算フィールド）';
COMMENT ON COLUMN customers.preferred_archetype IS '好みのアーキタイプ（RATタイプ）';
COMMENT ON COLUMN customers.diagnosis_count IS '診断回数';

