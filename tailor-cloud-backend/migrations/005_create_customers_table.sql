-- TailorCloud PostgreSQL マイグレーション
-- 顧客テーブル作成（CRM機能用）

CREATE TABLE IF NOT EXISTS customers (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    address TEXT,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- インデックス作成（パフォーマンス向上）
CREATE INDEX IF NOT EXISTS idx_customers_tenant_id ON customers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_customers_name ON customers(tenant_id, name);
CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(tenant_id, email);
CREATE INDEX IF NOT EXISTS idx_customers_phone ON customers(tenant_id, phone);
CREATE INDEX IF NOT EXISTS idx_customers_created_at ON customers(created_at DESC);

-- コメント追加
COMMENT ON TABLE customers IS '顧客テーブル（CRM機能用）';
COMMENT ON COLUMN customers.id IS '顧客ID（UUID）';
COMMENT ON COLUMN customers.tenant_id IS 'テナントID（マルチテナント分離用）';
COMMENT ON COLUMN customers.name IS '顧客名';
COMMENT ON COLUMN customers.email IS 'メールアドレス';
COMMENT ON COLUMN customers.phone IS '電話番号';
COMMENT ON COLUMN customers.address IS '住所';
COMMENT ON COLUMN customers.notes IS '備考';

