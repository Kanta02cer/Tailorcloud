-- TailorCloud PostgreSQL マイグレーション
-- 注文テーブル作成（Primary DB: PostgreSQL）

CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    customer_id VARCHAR(255) NOT NULL,
    fabric_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    compliance_doc_url TEXT,
    compliance_doc_hash VARCHAR(255),
    total_amount BIGINT NOT NULL,
    payment_due_date TIMESTAMP NOT NULL,
    delivery_date TIMESTAMP NOT NULL,
    measurement_data JSONB,
    adjustments JSONB,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL
);

-- インデックス作成（パフォーマンス向上）
CREATE INDEX IF NOT EXISTS idx_orders_tenant_id ON orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);

-- テナントIDと作成日の複合インデックス（よく使うクエリパターン）
CREATE INDEX IF NOT EXISTS idx_orders_tenant_created ON orders(tenant_id, created_at DESC);

-- コメント追加
COMMENT ON TABLE orders IS '注文テーブル（Primary DB: PostgreSQL）';
COMMENT ON COLUMN orders.id IS '注文ID（UUID）';
COMMENT ON COLUMN orders.tenant_id IS 'テナントID（マルチテナント分離用）';
COMMENT ON COLUMN orders.status IS '注文ステータス（Draft, Confirmed, Material_Secured, Cutting, Sewing, Inspection, Shipped, Delivered, Paid, Cancelled）';
COMMENT ON COLUMN orders.compliance_doc_url IS 'コンプライアンスドキュメント（契約書PDF）のURL';
COMMENT ON COLUMN orders.compliance_doc_hash IS 'コンプライアンスドキュメントのハッシュ値（改ざん防止用）';
COMMENT ON COLUMN orders.total_amount IS '合計金額（税抜、単位：円）';
COMMENT ON COLUMN orders.payment_due_date IS '支払期日（下請法60日ルール準拠）';
COMMENT ON COLUMN orders.delivery_date IS '納期';
COMMENT ON COLUMN orders.measurement_data IS '採寸データ（JSON形式）';
COMMENT ON COLUMN orders.adjustments IS '補正情報（JSON形式）';
COMMENT ON COLUMN orders.description IS '給付の内容（コンプライアンス用）';

