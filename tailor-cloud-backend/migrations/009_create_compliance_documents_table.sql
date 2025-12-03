-- ============================================================================
-- TailorCloud Enterprise: コンプライアンス文書履歴管理テーブル
-- ============================================================================
-- 目的: 下請法第3条に基づく発注書の修正履歴を管理
-- 上書き禁止: 元のPDFは保持し、修正時は新しいPDFを作成
-- ============================================================================

-- コンプライアンス文書テーブル
CREATE TABLE IF NOT EXISTS compliance_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    document_type VARCHAR(50) NOT NULL, -- 'INITIAL' (初回発注書), 'AMENDMENT' (修正発注書)
    parent_document_id UUID, -- 修正元の文書ID（修正発注書の場合）
    pdf_url TEXT NOT NULL, -- Cloud Storage上のPDFのURL
    pdf_hash VARCHAR(255) NOT NULL, -- SHA-256ハッシュ値（改ざん防止）
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    generated_by UUID NOT NULL, -- ユーザーID（発行者のID）
    amendment_reason TEXT, -- 修正理由（修正発注書の場合のみ）
    version INTEGER NOT NULL DEFAULT 1, -- バージョン番号（同一注文の連番）
    tenant_id UUID NOT NULL, -- テナントID（マルチテナント対応）
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT compliance_documents_order_fk FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    CONSTRAINT compliance_documents_parent_fk FOREIGN KEY (parent_document_id) REFERENCES compliance_documents(id) ON DELETE SET NULL,
    CONSTRAINT compliance_documents_tenant_fk FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT compliance_documents_type_check CHECK (document_type IN ('INITIAL', 'AMENDMENT'))
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_compliance_documents_order_id ON compliance_documents(order_id);
CREATE INDEX IF NOT EXISTS idx_compliance_documents_tenant_id ON compliance_documents(tenant_id);
CREATE INDEX IF NOT EXISTS idx_compliance_documents_parent_id ON compliance_documents(parent_document_id);
CREATE INDEX IF NOT EXISTS idx_compliance_documents_generated_at ON compliance_documents(generated_at DESC);

-- コメント
COMMENT ON TABLE compliance_documents IS 'コンプライアンス文書（下請法第3条に基づく発注書）の履歴管理';
COMMENT ON COLUMN compliance_documents.document_type IS '文書タイプ: INITIAL（初回発注書）, AMENDMENT（修正発注書）';
COMMENT ON COLUMN compliance_documents.parent_document_id IS '修正元の文書ID。修正発注書の場合のみ設定';
COMMENT ON COLUMN compliance_documents.pdf_hash IS 'PDFのSHA-256ハッシュ値（改ざん検知用）';
COMMENT ON COLUMN compliance_documents.amendment_reason IS '修正理由（修正発注書の場合のみ）';
COMMENT ON COLUMN compliance_documents.version IS 'バージョン番号（同一注文の連番。初回は1、修正ごとに+1）';

