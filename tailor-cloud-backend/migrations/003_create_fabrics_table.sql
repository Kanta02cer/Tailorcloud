-- TailorCloud PostgreSQL マイグレーション
-- 生地テーブル作成

CREATE TABLE IF NOT EXISTS fabrics (
    id VARCHAR(255) PRIMARY KEY,
    supplier_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    stock_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    price BIGINT NOT NULL,
    image_url TEXT,
    minimum_order DECIMAL(10,2) NOT NULL DEFAULT 3.2,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- インデックス作成（パフォーマンス向上）
CREATE INDEX IF NOT EXISTS idx_fabrics_supplier_id ON fabrics(supplier_id);
CREATE INDEX IF NOT EXISTS idx_fabrics_name ON fabrics(name);
CREATE INDEX IF NOT EXISTS idx_fabrics_stock_amount ON fabrics(stock_amount);

-- 検索用の全文検索インデックス（将来的に拡張）
-- CREATE INDEX IF NOT EXISTS idx_fabrics_name_fts ON fabrics USING gin(to_tsvector('japanese', name));

-- コメント追加
COMMENT ON TABLE fabrics IS '生地テーブル（Inventory画面で使用）';
COMMENT ON COLUMN fabrics.id IS '生地ID（UUID）';
COMMENT ON COLUMN fabrics.supplier_id IS '仕入先ID（生地問屋ID）';
COMMENT ON COLUMN fabrics.name IS '生地名';
COMMENT ON COLUMN fabrics.stock_amount IS '在庫数量（メートル）';
COMMENT ON COLUMN fabrics.price IS '単価（円/メートル）';
COMMENT ON COLUMN fabrics.image_url IS '生地画像URL（UI表示用、Cloud Storage）';
COMMENT ON COLUMN fabrics.minimum_order IS '最小発注数量（メートル、デフォルト3.2m = スーツ1着分）';

