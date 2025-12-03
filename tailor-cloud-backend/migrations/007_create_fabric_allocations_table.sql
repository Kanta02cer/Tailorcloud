-- ============================================================================
-- TailorCloud Enterprise: 反物引当（Fabric Allocation）テーブル作成
-- ============================================================================
-- 目的: 発注時にどの反物（Roll）のどの部分を使用するかを記録
-- これにより「どの反物を使ったか」を完全に追跡可能にする
-- ============================================================================

-- 反物引当テーブル
-- 発注（Order）と反物（Roll）の紐付け
CREATE TABLE IF NOT EXISTS fabric_allocations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    order_id UUID NOT NULL,
    order_item_id UUID, -- order_itemsテーブルのID（将来の拡張用）
    fabric_roll_id UUID NOT NULL,
    allocated_length DECIMAL(10, 2) NOT NULL, -- 引当数量（メートル）
    actual_used_length DECIMAL(10, 2), -- 実際に使用した数量（メートル、裁断後に記録）
    remnant_length DECIMAL(10, 2), -- 端尺（キレ）の長さ（メートル）
    allocation_status VARCHAR(50) NOT NULL DEFAULT 'RESERVED', -- 'RESERVED', 'CONFIRMED', 'CUT', 'CANCELLED'
    allocated_at TIMESTAMPTZ DEFAULT NOW(), -- 引当日時
    confirmed_at TIMESTAMPTZ, -- 確定日時
    cut_at TIMESTAMPTZ, -- 裁断日時
    notes TEXT, -- 備考
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (fabric_roll_id) REFERENCES fabric_rolls(id) ON DELETE RESTRICT,
    CONSTRAINT fabric_allocations_status_check CHECK (allocation_status IN ('RESERVED', 'CONFIRMED', 'CUT', 'CANCELLED')),
    CONSTRAINT fabric_allocations_length_check CHECK (allocated_length > 0)
);

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_fabric_allocations_tenant_id ON fabric_allocations(tenant_id);
CREATE INDEX IF NOT EXISTS idx_fabric_allocations_order_id ON fabric_allocations(order_id);
CREATE INDEX IF NOT EXISTS idx_fabric_allocations_fabric_roll_id ON fabric_allocations(fabric_roll_id);
CREATE INDEX IF NOT EXISTS idx_fabric_allocations_status ON fabric_allocations(allocation_status);

-- コメント追加
COMMENT ON TABLE fabric_allocations IS '反物引当テーブル: 発注時にどの反物（Roll）をどのくらい使用するかを記録';
COMMENT ON COLUMN fabric_allocations.allocated_length IS '引当数量（メートル）。発注時に確保された数量';
COMMENT ON COLUMN fabric_allocations.actual_used_length IS '実際に使用した数量（メートル）。裁断後に記録';
COMMENT ON COLUMN fabric_allocations.remnant_length IS '端尺（キレ）の長さ（メートル）。キレが発生した場合に記録';
COMMENT ON COLUMN fabric_allocations.allocation_status IS '引当状態: RESERVED(予約済み), CONFIRMED(確定), CUT(裁断済み), CANCELLED(キャンセル)';

-- 発注時には複数の反物から引当できるため、1つのorder_idに対して複数のfabric_roll_idが紐付く可能性がある

