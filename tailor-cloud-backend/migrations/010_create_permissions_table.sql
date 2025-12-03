-- ============================================================================
-- TailorCloud Enterprise: 権限管理システム（細分化されたRBAC）
-- ============================================================================
-- 目的: リソースベースの細かい権限管理
-- ロールだけでなく、リソース単位での権限を管理
-- ============================================================================

-- 権限マトリクステーブル（リソース×操作の組み合わせ）
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    resource_type VARCHAR(100) NOT NULL, -- 'ORDER', 'CUSTOMER', 'FABRIC', 'INVOICE', 'COMPLIANCE_DOCUMENT'
    resource_id UUID, -- NULLの場合は全リソースに対する権限
    action VARCHAR(50) NOT NULL, -- 'CREATE', 'READ', 'UPDATE', 'DELETE', 'GENERATE', 'APPROVE'
    role VARCHAR(50) NOT NULL, -- 'Owner', 'Staff', 'Factory_Manager', 'Worker'
    granted BOOLEAN NOT NULL DEFAULT TRUE, -- 許可/拒否
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT permissions_tenant_fk FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT permissions_resource_type_check CHECK (resource_type IN ('ORDER', 'CUSTOMER', 'FABRIC', 'INVOICE', 'COMPLIANCE_DOCUMENT', 'FABRIC_ROLL', 'ALL')),
    CONSTRAINT permissions_action_check CHECK (action IN ('CREATE', 'READ', 'UPDATE', 'DELETE', 'GENERATE', 'APPROVE', 'VIEW', 'ALL')),
    CONSTRAINT permissions_role_check CHECK (role IN ('Owner', 'Staff', 'Factory_Manager', 'Worker'))
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_permissions_tenant_resource ON permissions(tenant_id, resource_type);
CREATE INDEX IF NOT EXISTS idx_permissions_role ON permissions(role);
CREATE INDEX IF NOT EXISTS idx_permissions_resource_action ON permissions(resource_type, action);

-- デフォルト権限設定（テナントごとに設定可能）
-- 注文管理のデフォルト権限
INSERT INTO permissions (tenant_id, resource_type, resource_id, action, role, granted) VALUES
    -- Owner: 全権限
    ('00000000-0000-0000-0000-000000000000', 'ALL', NULL, 'ALL', 'Owner', TRUE),
    -- Staff: 注文作成・閲覧・修正（削除不可、承認不可）
    ('00000000-0000-0000-0000-000000000000', 'ORDER', NULL, 'CREATE', 'Staff', TRUE),
    ('00000000-0000-0000-0000-000000000000', 'ORDER', NULL, 'READ', 'Staff', TRUE),
    ('00000000-0000-0000-0000-000000000000', 'ORDER', NULL, 'UPDATE', 'Staff', TRUE),
    ('00000000-0000-0000-0000-000000000000', 'ORDER', NULL, 'DELETE', 'Staff', FALSE),
    ('00000000-0000-0000-0000-000000000000', 'ORDER', NULL, 'APPROVE', 'Staff', FALSE),
    -- Factory_Manager: 受注承認、工程管理
    ('00000000-0000-0000-0000-000000000000', 'ORDER', NULL, 'READ', 'Factory_Manager', TRUE),
    ('00000000-0000-0000-0000-000000000000', 'ORDER', NULL, 'APPROVE', 'Factory_Manager', TRUE),
    ('00000000-0000-0000-0000-000000000000', 'ORDER', NULL, 'UPDATE', 'Factory_Manager', TRUE),
    -- Worker: 作業完了チェックのみ
    ('00000000-0000-0000-0000-000000000000', 'ORDER', NULL, 'READ', 'Worker', TRUE),
    ('00000000-0000-0000-0000-000000000000', 'ORDER', NULL, 'UPDATE', 'Worker', TRUE)
ON CONFLICT DO NOTHING; -- 既に存在する場合はスキップ

-- コメント
COMMENT ON TABLE permissions IS '権限マトリクス（リソース×操作×ロール）';
COMMENT ON COLUMN permissions.resource_type IS 'リソースタイプ: ORDER, CUSTOMER, FABRIC, INVOICE, COMPLIANCE_DOCUMENT, FABRIC_ROLL, ALL';
COMMENT ON COLUMN permissions.resource_id IS '特定リソースのID。NULLの場合は全リソースに対する権限';
COMMENT ON COLUMN permissions.action IS '操作: CREATE, READ, UPDATE, DELETE, GENERATE, APPROVE, VIEW, ALL';
COMMENT ON COLUMN permissions.role IS 'ロール: Owner, Staff, Factory_Manager, Worker';
COMMENT ON COLUMN permissions.granted IS '許可(TRUE)または拒否(FALSE)';

