-- ============================================================================
-- TailorCloud: CRM拡張 - 顧客テーブルにタグ/履歴/状態フィールドを追加
-- 目的: 顧客管理機能の高度化に必要なデータ項目を追加する
-- ============================================================================

ALTER TABLE customers
    ADD COLUMN IF NOT EXISTS customer_status VARCHAR(32) DEFAULT 'lead',
    ADD COLUMN IF NOT EXISTS tags TEXT[] DEFAULT ARRAY[]::text[],
    ADD COLUMN IF NOT EXISTS vip_rank INTEGER DEFAULT 0,
    ADD COLUMN IF NOT EXISTS lifetime_value NUMERIC(12,2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS last_interaction_at TIMESTAMP NULL,
    ADD COLUMN IF NOT EXISTS interaction_notes JSONB DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS preferred_channel VARCHAR(50),
    ADD COLUMN IF NOT EXISTS lead_source VARCHAR(50);

-- 既存のnotesカラムをCRMメモ用途に活用するためコメントを更新
COMMENT ON COLUMN customers.notes IS '社内共有用メモ（顧客の背景や注意事項）';
COMMENT ON COLUMN customers.customer_status IS '顧客ステータス（lead/prospect/active/inactive/vip等）';
COMMENT ON COLUMN customers.tags IS '顧客タグ（嗜好やキャンペーン等）';
COMMENT ON COLUMN customers.vip_rank IS 'VIPランク（数値が高いほど優先度高）';
COMMENT ON COLUMN customers.lifetime_value IS '累計LTV（円）';
COMMENT ON COLUMN customers.last_interaction_at IS '最終接触日時';
COMMENT ON COLUMN customers.interaction_notes IS '問い合わせ・フォロー履歴（JSON配列）';
COMMENT ON COLUMN customers.preferred_channel IS '顧客が好むコミュニケーションチャネル';
COMMENT ON COLUMN customers.lead_source IS '顧客獲得経路';

CREATE INDEX IF NOT EXISTS idx_customers_status ON customers(customer_status);
CREATE INDEX IF NOT EXISTS idx_customers_tags ON customers USING GIN(tags);
CREATE INDEX IF NOT EXISTS idx_customers_last_interaction ON customers(last_interaction_at DESC);

