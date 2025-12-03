-- TailorCloud PostgreSQL マイグレーション
-- アンバサダー・成果報酬テーブル作成（Phase 1: Ambassador ID管理機能）

-- アンバサダーテーブル
CREATE TABLE IF NOT EXISTS ambassadors (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    status VARCHAR(50) NOT NULL DEFAULT 'Active',
    commission_rate DECIMAL(5,4) NOT NULL DEFAULT 0.1000,
    total_sales BIGINT NOT NULL DEFAULT 0,
    total_commission BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 成果報酬テーブル
CREATE TABLE IF NOT EXISTS commissions (
    id VARCHAR(255) PRIMARY KEY,
    order_id VARCHAR(255) NOT NULL,
    ambassador_id VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(255) NOT NULL,
    order_amount BIGINT NOT NULL,
    commission_rate DECIMAL(5,4) NOT NULL,
    commission_amount BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'Pending',
    paid_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- インデックス作成（パフォーマンス向上）
CREATE INDEX IF NOT EXISTS idx_ambassadors_tenant_id ON ambassadors(tenant_id);
CREATE INDEX IF NOT EXISTS idx_ambassadors_user_id ON ambassadors(user_id);
CREATE INDEX IF NOT EXISTS idx_ambassadors_status ON ambassadors(status);

CREATE INDEX IF NOT EXISTS idx_commissions_order_id ON commissions(order_id);
CREATE INDEX IF NOT EXISTS idx_commissions_ambassador_id ON commissions(ambassador_id);
CREATE INDEX IF NOT EXISTS idx_commissions_tenant_id ON commissions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_commissions_status ON commissions(status);
CREATE INDEX IF NOT EXISTS idx_commissions_created_at ON commissions(created_at DESC);

-- コメント追加
COMMENT ON TABLE ambassadors IS 'アンバサダーテーブル（学生アンバサダー情報・成果報酬管理）';
COMMENT ON COLUMN ambassadors.id IS 'アンバサダーID（UUID）';
COMMENT ON COLUMN ambassadors.user_id IS 'Firebase AuthのUserID（一意制約）';
COMMENT ON COLUMN ambassadors.commission_rate IS '成果報酬率（例: 0.1000 = 10%）';
COMMENT ON COLUMN ambassadors.total_sales IS '累計売上（円）';
COMMENT ON COLUMN ambassadors.total_commission IS '累計報酬（円）';

COMMENT ON TABLE commissions IS '成果報酬テーブル（注文ごとの成果報酬記録）';
COMMENT ON COLUMN commissions.id IS '成果報酬ID（UUID）';
COMMENT ON COLUMN commissions.order_id IS '注文ID';
COMMENT ON COLUMN commissions.ambassador_id IS 'アンバサダーID';
COMMENT ON COLUMN commissions.status IS '成果報酬ステータス（Pending, Approved, Paid, Cancelled）';
COMMENT ON COLUMN commissions.paid_at IS '支払日（nil = 未払い）';

