-- ============================================================================
-- TailorCloud Enterprise: インボイス制度対応フィールド追加
-- ============================================================================
-- 目的: 2023年10月施行のインボイス制度（適格請求書等保存方式）に対応
-- ============================================================================

-- テナントテーブルにインボイス制度対応フィールドを追加
ALTER TABLE tenants
ADD COLUMN IF NOT EXISTS invoice_registration_no VARCHAR(50), -- インボイス登録番号（T番号）
ADD COLUMN IF NOT EXISTS legal_name VARCHAR(255), -- 法人名（請求書記載用）
ADD COLUMN IF NOT EXISTS address TEXT, -- 住所（請求書記載用）
ADD COLUMN IF NOT EXISTS tax_rounding_method VARCHAR(20) DEFAULT 'HALF_UP'; -- 端数処理方法（HALF_UP, ROUND_DOWN, ROUND_UP）

-- 注文テーブルにインボイス制度対応フィールドを追加
ALTER TABLE orders
ADD COLUMN IF NOT EXISTS tax_amount DECIMAL(12, 2) DEFAULT 0, -- 消費税額
ADD COLUMN IF NOT EXISTS tax_rate DECIMAL(3, 2) DEFAULT 0.10, -- 消費税率（10% or 8%）
ADD COLUMN IF NOT EXISTS tax_excluded_amount DECIMAL(12, 2), -- 税抜金額
ADD COLUMN IF NOT EXISTS invoice_issued_at TIMESTAMPTZ; -- 請求書発行日時

-- インデックス追加
CREATE INDEX IF NOT EXISTS idx_tenants_invoice_registration_no ON tenants(invoice_registration_no);

-- コメント追加
COMMENT ON COLUMN tenants.invoice_registration_no IS 'インボイス登録番号（T番号）。適格請求書に記載';
COMMENT ON COLUMN tenants.tax_rounding_method IS '端数処理方法: HALF_UP(四捨五入), ROUND_DOWN(切り捨て), ROUND_UP(切り上げ)';
COMMENT ON COLUMN orders.tax_amount IS '消費税額（税率ごとに計算）';
COMMENT ON COLUMN orders.tax_rate IS '消費税率（0.10 = 10%, 0.08 = 8%）';
COMMENT ON COLUMN orders.tax_excluded_amount IS '税抜金額（合計金額から消費税を除いた額）';
COMMENT ON COLUMN orders.invoice_issued_at IS '請求書発行日時（インボイス制度対応）';

