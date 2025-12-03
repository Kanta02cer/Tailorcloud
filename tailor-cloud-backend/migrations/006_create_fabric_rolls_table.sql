-- ============================================================================
-- TailorCloud Enterprise: 反物（Roll）管理テーブル作成
-- ============================================================================
-- 目的: 単なる「総量」管理ではなく「物理的な巻き(Roll)」単位で管理
-- これにより「キレ（端尺）問題」を根本的に解決
-- ============================================================================

-- 反物（Roll）テーブル
-- 1本の物理的な反物を1レコードとして管理
CREATE TABLE IF NOT EXISTS fabric_rolls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    fabric_id UUID NOT NULL,
    roll_number VARCHAR(100) NOT NULL, -- ロール番号（例: "VBC-2025-001"）
    initial_length DECIMAL(10, 2) NOT NULL, -- 初期長さ（メートル）
    current_length DECIMAL(10, 2) NOT NULL, -- 現在の残り長さ（メートル）
    width DECIMAL(5, 2), -- 幅（センチメートル、オプション）
    supplier_lot_no VARCHAR(100), -- 仕入先ロット番号
    received_at TIMESTAMPTZ, -- 入荷日
    location VARCHAR(255), -- 保管場所（例: "倉庫A-3F-12"）
    status VARCHAR(50) NOT NULL DEFAULT 'AVAILABLE', -- 'AVAILABLE', 'ALLOCATED', 'CONSUMED', 'DAMAGED'
    notes TEXT, -- 備考
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (fabric_id) REFERENCES fabrics(id) ON DELETE CASCADE,
    CONSTRAINT fabric_rolls_status_check CHECK (status IN ('AVAILABLE', 'ALLOCATED', 'CONSUMED', 'DAMAGED')),
    CONSTRAINT fabric_rolls_length_check CHECK (current_length >= 0 AND current_length <= initial_length),
    CONSTRAINT fabric_rolls_tenant_roll_unique UNIQUE (tenant_id, roll_number)
);

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_fabric_rolls_tenant_id ON fabric_rolls(tenant_id);
CREATE INDEX IF NOT EXISTS idx_fabric_rolls_fabric_id ON fabric_rolls(fabric_id);
CREATE INDEX IF NOT EXISTS idx_fabric_rolls_status ON fabric_rolls(status);
CREATE INDEX IF NOT EXISTS idx_fabric_rolls_available ON fabric_rolls(tenant_id, fabric_id, status) WHERE status = 'AVAILABLE';

-- コメント追加
COMMENT ON TABLE fabric_rolls IS '反物（Roll）管理テーブル: 物理的な1本の反物を1レコードとして管理';
COMMENT ON COLUMN fabric_rolls.roll_number IS 'ロール番号（テナント内でユニーク）';
COMMENT ON COLUMN fabric_rolls.initial_length IS '初期長さ（メートル）';
COMMENT ON COLUMN fabric_rolls.current_length IS '現在の残り長さ（メートル）。使用されるたびに減算される';
COMMENT ON COLUMN fabric_rolls.status IS '状態: AVAILABLE(利用可能), ALLOCATED(引当済み), CONSUMED(消費済み), DAMAGED(破損)';

