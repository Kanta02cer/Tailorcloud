-- TailorCloud PostgreSQL マイグレーション
-- 監査ログテーブル作成（仕様書要件: 法的証拠能力のため）

-- 監査ログテーブル
CREATE TABLE IF NOT EXISTS audit_logs (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    action VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id VARCHAR(255) NOT NULL,
    old_value TEXT,
    new_value TEXT,
    changed_fields JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 契約書閲覧ログテーブル
CREATE TABLE IF NOT EXISTS compliance_document_view_logs (
    id VARCHAR(255) PRIMARY KEY,
    order_id VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    document_url TEXT NOT NULL,
    document_hash VARCHAR(255) NOT NULL,
    viewed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(45),
    user_agent TEXT
);

-- インデックス作成（パフォーマンス向上）
CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_id ON audit_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);

-- 契約書閲覧ログのインデックス
CREATE INDEX IF NOT EXISTS idx_compliance_view_logs_order_id ON compliance_document_view_logs(order_id);
CREATE INDEX IF NOT EXISTS idx_compliance_view_logs_tenant_id ON compliance_document_view_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_compliance_view_logs_viewed_at ON compliance_document_view_logs(viewed_at DESC);

-- コメント追加
COMMENT ON TABLE audit_logs IS '監査ログテーブル（誰が・いつ・何を変更したかの完全なログ）';
COMMENT ON COLUMN audit_logs.id IS '監査ログID（UUID）';
COMMENT ON COLUMN audit_logs.tenant_id IS 'テナントID（マルチテナント分離用）';
COMMENT ON COLUMN audit_logs.user_id IS '実行ユーザーID';
COMMENT ON COLUMN audit_logs.action IS 'アクションタイプ（CREATE, UPDATE, DELETE, VIEW, CONFIRM, STATUS_CHANGE）';
COMMENT ON COLUMN audit_logs.resource_type IS 'リソースタイプ（order, customer, fabric等）';
COMMENT ON COLUMN audit_logs.resource_id IS 'リソースID';
COMMENT ON COLUMN audit_logs.old_value IS '変更前の値（JSON形式）';
COMMENT ON COLUMN audit_logs.new_value IS '変更後の値（JSON形式）';
COMMENT ON COLUMN audit_logs.changed_fields IS '変更されたフィールド名のリスト（JSON配列）';
COMMENT ON COLUMN audit_logs.ip_address IS '実行IPアドレス';
COMMENT ON COLUMN audit_logs.user_agent IS '実行ユーザーエージェント';

COMMENT ON TABLE compliance_document_view_logs IS '契約書閲覧ログテーブル（いつ契約書を閲覧したかの完全なログ）';
COMMENT ON COLUMN compliance_document_view_logs.id IS '閲覧ログID（UUID）';
COMMENT ON COLUMN compliance_document_view_logs.order_id IS '注文ID';
COMMENT ON COLUMN compliance_document_view_logs.document_hash IS '契約書のハッシュ値（改ざん検出用）';

