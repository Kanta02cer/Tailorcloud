-- ============================================================================
-- TailorCloud: Suit-MBTI統合 - 予約管理テーブル作成
-- ============================================================================
-- 目的: フィッティング予約を管理
-- ============================================================================

-- Appointments (予約) テーブル
CREATE TABLE IF NOT EXISTS appointments (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL, -- 顧客ID
    tenant_id VARCHAR(255) NOT NULL,
    fitter_id VARCHAR(255), -- フィッターID（スタッフ）
    appointment_datetime TIMESTAMPTZ NOT NULL,
    duration_minutes INTEGER DEFAULT 60, -- 予約時間（分）
    status VARCHAR(20) NOT NULL DEFAULT 'Pending', -- Confirmed/Cancelled/Completed
    deposit_amount BIGINT, -- デポジット金額（円）
    deposit_payment_intent_id VARCHAR(255), -- Stripe Payment Intent ID
    deposit_status VARCHAR(20), -- pending/succeeded/failed/refunded
    notes TEXT, -- メモ
    cancelled_at TIMESTAMPTZ, -- キャンセル日時
    cancelled_reason TEXT, -- キャンセル理由
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
    -- 注意: tenantsテーブルが存在する場合のみ外部キー制約を有効化
    -- FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT appointments_status_check CHECK (status IN ('Pending', 'Confirmed', 'Cancelled', 'Completed', 'NoShow')),
    CONSTRAINT appointments_deposit_status_check CHECK (deposit_status IN ('pending', 'succeeded', 'failed', 'refunded') OR deposit_status IS NULL)
);

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_appointments_user_id ON appointments(user_id);
CREATE INDEX IF NOT EXISTS idx_appointments_tenant_id ON appointments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_appointments_fitter_id ON appointments(fitter_id);
CREATE INDEX IF NOT EXISTS idx_appointments_datetime ON appointments(appointment_datetime);
CREATE INDEX IF NOT EXISTS idx_appointments_status ON appointments(status);
CREATE INDEX IF NOT EXISTS idx_appointments_tenant_datetime ON appointments(tenant_id, appointment_datetime);
CREATE INDEX IF NOT EXISTS idx_appointments_tenant_status ON appointments(tenant_id, status);

-- コメント追加
COMMENT ON TABLE appointments IS 'フィッティング予約テーブル';
COMMENT ON COLUMN appointments.fitter_id IS 'フィッターID（スタッフ）';
COMMENT ON COLUMN appointments.appointment_datetime IS '予約日時';
COMMENT ON COLUMN appointments.deposit_amount IS 'デポジット金額（円）';
COMMENT ON COLUMN appointments.deposit_payment_intent_id IS 'Stripe Payment Intent ID';

