-- ============================================================================
-- TailorCloud Enterprise Database Schema (PostgreSQL)
-- ============================================================================
-- Concept: Multi-Tenant, Event Sourcing Lite, Supply Chain Aware
-- Version: 1.0.0
-- Date: 2025-01
-- 
-- このスキーマは、日本の商社・大手アパレル企業向けのエンタープライズ要件を
-- 満たすために設計されています。
-- 
-- 主要な特徴:
-- 1. マルチテナント対応（テナント間のデータ完全分離）
-- 2. 反物（Roll）単位の在庫管理（キレ問題の解決）
-- 3. 下請法・インボイス制度対応
-- 4. 監査ログ（5年以上保存）
-- 5. Row Level Security (RLS) によるデータ分離
-- ============================================================================

-- ============================================================================
-- 1. Organization & Security (テナント管理・セキュリティ)
-- ============================================================================

-- テナントテーブル（企業・店舗・工場）
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL, -- "Regalis Group", "ABC Trading"
    legal_name VARCHAR(255), -- 法人名（下請法記載用）
    address TEXT, -- 住所（下請法記載用）
    invoice_registration_no VARCHAR(50), -- インボイス登録番号（T番号）
    tax_rounding_method VARCHAR(20) DEFAULT 'HALF_UP', -- 端数処理方法
    type VARCHAR(50) NOT NULL, -- 'TAILOR' (発注者), 'FACTORY' (受注者), 'SUPPLIER' (仕入先)
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT tenants_type_check CHECK (type IN ('TAILOR', 'FACTORY', 'SUPPLIER'))
);

-- ユーザーテーブル
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255), -- Firebase Auth使用時はNULL可
    firebase_uid VARCHAR(255) UNIQUE, -- Firebase Authentication UID
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL, -- 'OWNER', 'STORE_MGR', 'STAFF', 'FACTORY_MGR', 'WORKER'
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT users_role_check CHECK (role IN ('OWNER', 'STORE_MGR', 'STAFF', 'FACTORY_MGR', 'WORKER')),
    CONSTRAINT users_email_tenant_unique UNIQUE (email, tenant_id)
);

-- デバイステーブル（端末認証用）
CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    device_id VARCHAR(255) UNIQUE NOT NULL, -- デバイス固有ID（UUID）
    device_name VARCHAR(255), -- デバイス名（例: "店舗1のiPad"）
    ip_address VARCHAR(45), -- 最後に使用したIPアドレス
    is_active BOOLEAN DEFAULT TRUE,
    registered_at TIMESTAMPTZ DEFAULT NOW(),
    last_access_at TIMESTAMPTZ,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ============================================================================
-- 2. Master Data (マスタ: 生地・仕様)
-- ============================================================================

-- 生地マスタテーブル
CREATE TABLE IF NOT EXISTS fabrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    supplier_id UUID, -- 仕入先テナントID
    sku VARCHAR(100) NOT NULL, -- 品番
    brand VARCHAR(100) NOT NULL, -- "V.B.C", "Zegna"
    name VARCHAR(255) NOT NULL,
    composition VARCHAR(255), -- 混率 "Wool 100%"
    cost_price DECIMAL(10, 2), -- 原価
    sales_price DECIMAL(10, 2), -- 上代（販売価格）
    image_url TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (supplier_id) REFERENCES tenants(id) ON DELETE SET NULL,
    CONSTRAINT fabrics_tenant_sku_unique UNIQUE (tenant_id, sku) -- テナント内で品番はユニーク
);

-- ============================================================================
-- 3. Inventory (在庫: 反物管理) -> 最重要差別化ポイント
-- ============================================================================
-- 単なる「総量」ではなく「物理的な巻き(Roll)」単位で管理
-- これにより「キレ（端尺）問題」を解決

-- 反物（Roll）テーブル
CREATE TABLE IF NOT EXISTS fabric_rolls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    fabric_id UUID NOT NULL,
    location_id UUID, -- 倉庫ID or 店舗ID（将来拡張用）
    lot_number VARCHAR(100), -- ロット番号（例: "LOT-2025-001"）
    initial_length DECIMAL(10, 2) NOT NULL, -- 入荷時長さ (m)
    current_length DECIMAL(10, 2) NOT NULL, -- 現在長さ (m)
    status VARCHAR(50) DEFAULT 'AVAILABLE', -- AVAILABLE, RESERVED, CUT, EMPTY
    ownership_type VARCHAR(50) DEFAULT 'OWNER', -- OWNER, CONSIGNED, ON_CONSIGNMENT
    owner_tenant_id UUID, -- 所有権を持つテナント
    custodian_tenant_id UUID, -- 保管場所のテナント
    version INTEGER DEFAULT 1, -- 楽観的ロック用
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (fabric_id) REFERENCES fabrics(id) ON DELETE CASCADE,
    FOREIGN KEY (owner_tenant_id) REFERENCES tenants(id) ON DELETE SET NULL,
    FOREIGN KEY (custodian_tenant_id) REFERENCES tenants(id) ON DELETE SET NULL,
    CONSTRAINT fabric_rolls_status_check CHECK (status IN ('AVAILABLE', 'RESERVED', 'CUT', 'EMPTY')),
    CONSTRAINT fabric_rolls_ownership_check CHECK (ownership_type IN ('OWNER', 'CONSIGNED', 'ON_CONSIGNMENT')),
    CONSTRAINT fabric_rolls_length_check CHECK (current_length >= 0 AND current_length <= initial_length)
);

-- ============================================================================
-- 4. Order & Transaction (受注・発注)
-- ============================================================================

-- 顧客テーブル
CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    address TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

-- 注文テーブル
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    order_number VARCHAR(50) NOT NULL, -- 顧客に見せる注文番号（例: "ORD-2025-001"）
    customer_id UUID,
    store_id UUID, -- 店舗ID（将来拡張用）
    
    -- Status Management
    status VARCHAR(50) NOT NULL DEFAULT 'DRAFT', -- DRAFT, CONFIRMED, PRODUCTION, SHIPPED, DELIVERED, PAID, CANCELLED
    
    -- Pricing
    total_amount DECIMAL(12, 2) NOT NULL, -- 税抜合計金額
    tax_amount DECIMAL(12, 2) NOT NULL DEFAULT 0, -- 消費税額
    tax_rate DECIMAL(3, 2) NOT NULL DEFAULT 0.10, -- 消費税率（10% or 8%）
    currency VARCHAR(3) DEFAULT 'JPY',
    
    -- Dates
    ordered_at TIMESTAMPTZ,
    delivery_due_date DATE,
    payment_due_date DATE NOT NULL, -- 支払期日（下請法60日ルール）
    
    -- Compliance
    compliance_doc_url TEXT, -- 下請法3条書面PDFのURL
    compliance_doc_hash VARCHAR(255), -- PDFのハッシュ値（改ざん防止）
    invoice_issued_at TIMESTAMPTZ, -- 請求書発行日時
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID NOT NULL, -- 作成ユーザーID
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT orders_status_check CHECK (status IN ('DRAFT', 'CONFIRMED', 'PRODUCTION', 'SHIPPED', 'DELIVERED', 'PAID', 'CANCELLED')),
    CONSTRAINT orders_tenant_order_number_unique UNIQUE (tenant_id, order_number)
);

-- 注文明細テーブル
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    item_type VARCHAR(50), -- 'SUIT', 'SHIRT', 'COAT' など
    fabric_id UUID REFERENCES fabrics(id),
    
    -- Measurement Data (JSONB for flexibility)
    -- 寸法データは将来的な項目変更に備えJSONBで格納
    measurements JSONB NOT NULL, -- { "jacket_length": 72.5, "sleeve": 60.0 ... }
    options JSONB, -- { "lapel": "notch", "lining": "cupra_A" ... }
    
    required_fabric_length DECIMAL(5, 2) NOT NULL, -- 必要用尺 (例: 3.2m)
    unit_price DECIMAL(10, 2) NOT NULL, -- 単価
    quantity INTEGER NOT NULL DEFAULT 1, -- 数量
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (fabric_id) REFERENCES fabrics(id) ON DELETE RESTRICT
);

-- ============================================================================
-- 5. Allocations (在庫引当)
-- ============================================================================
-- どの注文が、どの反物(Roll)のどの部分を確保しているか

-- 在庫引当テーブル
CREATE TABLE IF NOT EXISTS fabric_allocations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_item_id UUID NOT NULL,
    fabric_roll_id UUID NOT NULL,
    allocated_length DECIMAL(5, 2) NOT NULL, -- 確保した長さ (m)
    allocated_at TIMESTAMPTZ DEFAULT NOW(),
    status VARCHAR(50) DEFAULT 'RESERVED', -- RESERVED, CUT (裁断済み), CANCELLED
    FOREIGN KEY (order_item_id) REFERENCES order_items(id) ON DELETE CASCADE,
    FOREIGN KEY (fabric_roll_id) REFERENCES fabric_rolls(id) ON DELETE RESTRICT,
    CONSTRAINT fabric_allocations_status_check CHECK (status IN ('RESERVED', 'CUT', 'CANCELLED'))
);

-- ============================================================================
-- 6. Compliance Documents (コンプライアンス文書)
-- ============================================================================

-- コンプライアンス文書テーブル（下請法3条書面、修正注文書など）
CREATE TABLE IF NOT EXISTS compliance_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    document_type VARCHAR(50) NOT NULL, -- 'INITIAL' (初回), 'AMENDMENT' (修正)
    parent_document_id UUID, -- 修正元の文書ID（修正注文書の場合）
    pdf_url TEXT NOT NULL,
    pdf_hash VARCHAR(255) NOT NULL, -- SHA-256
    generated_at TIMESTAMPTZ NOT NULL,
    generated_by UUID NOT NULL,
    amendment_reason TEXT, -- 修正理由（修正注文書の場合）
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_document_id) REFERENCES compliance_documents(id) ON DELETE SET NULL,
    FOREIGN KEY (generated_by) REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT compliance_documents_type_check CHECK (document_type IN ('INITIAL', 'AMENDMENT'))
);

-- ============================================================================
-- 7. Audit Log (監査ログ) -> エンタープライズ必須
-- ============================================================================

-- 監査ログテーブル
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID,
    user_id UUID,
    action VARCHAR(100) NOT NULL, -- "UPDATE_PRICE", "FORCE_ADJUST_STOCK", "CREATE_ORDER"
    target_table VARCHAR(50), -- "orders", "fabrics", "fabric_rolls" など
    target_id UUID,
    old_value JSONB,
    new_value JSONB,
    changed_fields JSONB, -- 変更されたフィールド名のリスト
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 契約書閲覧ログテーブル
CREATE TABLE IF NOT EXISTS compliance_document_view_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    document_id UUID NOT NULL, -- compliance_documents.id
    document_url TEXT NOT NULL,
    document_hash VARCHAR(255) NOT NULL,
    viewed_at TIMESTAMPTZ DEFAULT NOW(),
    ip_address VARCHAR(45),
    user_agent TEXT,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (document_id) REFERENCES compliance_documents(id) ON DELETE CASCADE
);

-- ============================================================================
-- 8. Ambassador & Commission (アンバサダー・成果報酬)
-- ============================================================================

-- アンバサダーテーブル
CREATE TABLE IF NOT EXISTS ambassadors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID UNIQUE, -- Firebase AuthのUserID
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, INACTIVE
    commission_rate DECIMAL(5, 4) NOT NULL DEFAULT 0.1000, -- 成果報酬率（例: 0.1000 = 10%）
    total_sales BIGINT NOT NULL DEFAULT 0, -- 累計売上（円）
    total_commission BIGINT NOT NULL DEFAULT 0, -- 累計報酬（円）
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT ambassadors_status_check CHECK (status IN ('ACTIVE', 'INACTIVE'))
);

-- 成果報酬テーブル
CREATE TABLE IF NOT EXISTS commissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    ambassador_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    order_amount BIGINT NOT NULL,
    commission_rate DECIMAL(5, 4) NOT NULL,
    commission_amount BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING', -- PENDING, APPROVED, PAID, CANCELLED
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (ambassador_id) REFERENCES ambassadors(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT commissions_status_check CHECK (status IN ('PENDING', 'APPROVED', 'PAID', 'CANCELLED'))
);

-- ============================================================================
-- 9. Indexes for Performance (パフォーマンス向上のためのインデックス)
-- ============================================================================

-- Tenants
CREATE INDEX IF NOT EXISTS idx_tenants_type ON tenants(type);
CREATE INDEX IF NOT EXISTS idx_tenants_active ON tenants(is_active) WHERE is_active = TRUE;

-- Users
CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_firebase_uid ON users(firebase_uid) WHERE firebase_uid IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Devices
CREATE INDEX IF NOT EXISTS idx_devices_device_id ON devices(device_id);
CREATE INDEX IF NOT EXISTS idx_devices_tenant_id ON devices(tenant_id);
CREATE INDEX IF NOT EXISTS idx_devices_user_id ON devices(user_id);

-- Fabrics
CREATE INDEX IF NOT EXISTS idx_fabrics_tenant_id ON fabrics(tenant_id);
CREATE INDEX IF NOT EXISTS idx_fabrics_sku ON fabrics(tenant_id, sku);
CREATE INDEX IF NOT EXISTS idx_fabrics_brand ON fabrics(brand);
CREATE INDEX IF NOT EXISTS idx_fabrics_active ON fabrics(is_active) WHERE is_active = TRUE;

-- Fabric Rolls (最重要: 在庫引当のパフォーマンス)
CREATE INDEX IF NOT EXISTS idx_fabric_rolls_fabric_id ON fabric_rolls(fabric_id);
CREATE INDEX IF NOT EXISTS idx_fabric_rolls_status ON fabric_rolls(status) WHERE status = 'AVAILABLE';
CREATE INDEX IF NOT EXISTS idx_fabric_rolls_tenant_fabric_status ON fabric_rolls(tenant_id, fabric_id, status);
CREATE INDEX IF NOT EXISTS idx_fabric_rolls_location ON fabric_rolls(location_id) WHERE location_id IS NOT NULL;

-- Customers
CREATE INDEX IF NOT EXISTS idx_customers_tenant_id ON customers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_customers_name ON customers(tenant_id, name);

-- Orders
CREATE INDEX IF NOT EXISTS idx_orders_tenant_id ON orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_orders_tenant_created ON orders(tenant_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_orders_order_number ON orders(tenant_id, order_number);

-- Order Items
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_fabric_id ON order_items(fabric_id);

-- Fabric Allocations (最重要: 在庫引当のパフォーマンス)
CREATE INDEX IF NOT EXISTS idx_fabric_allocations_order_item_id ON fabric_allocations(order_item_id);
CREATE INDEX IF NOT EXISTS idx_fabric_allocations_fabric_roll_id ON fabric_allocations(fabric_roll_id);
CREATE INDEX IF NOT EXISTS idx_fabric_allocations_status ON fabric_allocations(status);

-- Compliance Documents
CREATE INDEX IF NOT EXISTS idx_compliance_documents_order_id ON compliance_documents(order_id);
CREATE INDEX IF NOT EXISTS idx_compliance_documents_type ON compliance_documents(document_type);

-- Audit Logs
CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_id ON audit_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_target ON audit_logs(target_table, target_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
-- 5年以上の保存期間を考慮したパーティショニング（将来実装）
-- CREATE INDEX IF NOT EXISTS idx_audit_logs_year_month ON audit_logs(DATE_TRUNC('month', created_at));

-- Compliance Document View Logs
CREATE INDEX IF NOT EXISTS idx_compliance_view_logs_order_id ON compliance_document_view_logs(order_id);
CREATE INDEX IF NOT EXISTS idx_compliance_view_logs_tenant_id ON compliance_document_view_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_compliance_view_logs_viewed_at ON compliance_document_view_logs(viewed_at DESC);

-- Ambassadors
CREATE INDEX IF NOT EXISTS idx_ambassadors_tenant_id ON ambassadors(tenant_id);
CREATE INDEX IF NOT EXISTS idx_ambassadors_user_id ON ambassadors(user_id);
CREATE INDEX IF NOT EXISTS idx_ambassadors_status ON ambassadors(status);

-- Commissions
CREATE INDEX IF NOT EXISTS idx_commissions_order_id ON commissions(order_id);
CREATE INDEX IF NOT EXISTS idx_commissions_ambassador_id ON commissions(ambassador_id);
CREATE INDEX IF NOT EXISTS idx_commissions_tenant_id ON commissions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_commissions_status ON commissions(status);

-- ============================================================================
-- 10. Row Level Security (RLS) Policies (データ分離)
-- ============================================================================
-- テナント間のデータ完全分離を保証

-- RLSを有効化
ALTER TABLE tenants ENABLE ROW LEVEL SECURITY;
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE devices ENABLE ROW LEVEL SECURITY;
ALTER TABLE fabrics ENABLE ROW LEVEL SECURITY;
ALTER TABLE fabric_rolls ENABLE ROW LEVEL SECURITY;
ALTER TABLE customers ENABLE ROW LEVEL SECURITY;
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE order_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE fabric_allocations ENABLE ROW LEVEL SECURITY;
ALTER TABLE compliance_documents ENABLE ROW LEVEL SECURITY;
ALTER TABLE audit_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE compliance_document_view_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE ambassadors ENABLE ROW LEVEL SECURITY;
ALTER TABLE commissions ENABLE ROW LEVEL SECURITY;

-- テナント分離ポリシー（例: usersテーブル）
-- 注意: 実際の実装では、セッション変数やJWTからtenant_idを取得する必要があります
-- CREATE POLICY tenant_isolation_users ON users
--   USING (tenant_id = current_setting('app.current_tenant_id', TRUE)::UUID);

-- ============================================================================
-- 11. Triggers (自動更新)
-- ============================================================================

-- updated_atを自動更新するトリガー関数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- updated_atトリガーの設定（主要テーブル）
CREATE TRIGGER update_tenants_updated_at BEFORE UPDATE ON tenants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_fabrics_updated_at BEFORE UPDATE ON fabrics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_fabric_rolls_updated_at BEFORE UPDATE ON fabric_rolls
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_order_items_updated_at BEFORE UPDATE ON order_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ambassadors_updated_at BEFORE UPDATE ON ambassadors
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_commissions_updated_at BEFORE UPDATE ON commissions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- 12. Comments (テーブル・カラムの説明)
-- ============================================================================

COMMENT ON TABLE tenants IS 'テナントテーブル（企業・店舗・工場）';
COMMENT ON COLUMN tenants.invoice_registration_no IS 'インボイス登録番号（T番号）';
COMMENT ON COLUMN tenants.tax_rounding_method IS '端数処理方法（HALF_UP, DOWN, UP）';

COMMENT ON TABLE fabric_rolls IS '反物（Roll）テーブル - 物理的な巻き単位での在庫管理（キレ問題の解決）';
COMMENT ON COLUMN fabric_rolls.version IS '楽観的ロック用のバージョン番号';
COMMENT ON COLUMN fabric_rolls.ownership_type IS '所有権タイプ（OWNER: 自社所有, CONSIGNED: 委託在庫, ON_CONSIGNMENT: 委託出庫）';

COMMENT ON TABLE fabric_allocations IS '在庫引当テーブル - どの注文が、どの反物のどの部分を確保しているか';
COMMENT ON COLUMN fabric_allocations.allocated_length IS '確保した長さ（メートル）';

COMMENT ON TABLE compliance_documents IS 'コンプライアンス文書テーブル（下請法3条書面、修正注文書など）';
COMMENT ON COLUMN compliance_documents.document_type IS '文書タイプ（INITIAL: 初回, AMENDMENT: 修正）';
COMMENT ON COLUMN compliance_documents.pdf_hash IS 'PDFのハッシュ値（SHA-256、改ざん防止用）';

COMMENT ON TABLE audit_logs IS '監査ログテーブル - 5年以上保存（法的証拠能力のため）';
COMMENT ON COLUMN audit_logs.action IS 'アクションタイプ（例: "UPDATE_PRICE", "FORCE_ADJUST_STOCK"）';
COMMENT ON COLUMN audit_logs.changed_fields IS '変更されたフィールド名のリスト（JSON配列）';

-- ============================================================================
-- END OF SCHEMA
-- ============================================================================

