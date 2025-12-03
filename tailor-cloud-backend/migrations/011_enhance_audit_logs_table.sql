-- ============================================================================
-- TailorCloud Enterprise: 監査ログテーブルの強化
-- ============================================================================
-- 目的: 5年以上の監査ログ保存と改ざん防止対応
-- 全操作の記録、IPアドレス・デバイス情報、変更前後の値の記録
-- ============================================================================

-- 既存のaudit_logsテーブルを拡張
ALTER TABLE audit_logs
ADD COLUMN IF NOT EXISTS ip_address VARCHAR(45), -- IPアドレス（IPv4/IPv6対応）
ADD COLUMN IF NOT EXISTS user_agent TEXT, -- ユーザーエージェント（デバイス情報）
ADD COLUMN IF NOT EXISTS device_id VARCHAR(255), -- デバイスID（デバイス認証用）
ADD COLUMN IF NOT EXISTS old_values JSONB, -- 変更前の値（JSON形式）
ADD COLUMN IF NOT EXISTS new_values JSONB, -- 変更後の値（JSON形式）
ADD COLUMN IF NOT EXISTS change_summary TEXT, -- 変更サマリー（人間が読みやすい形式）
ADD COLUMN IF NOT EXISTS log_hash VARCHAR(255), -- ログのハッシュ値（改ざん検知用）
ADD COLUMN IF NOT EXISTS archived_at TIMESTAMPTZ, -- アーカイブ日時
ADD COLUMN IF NOT EXISTS archive_location TEXT; -- アーカイブ先（Cloud Storageパス）

-- インデックス追加
CREATE INDEX IF NOT EXISTS idx_audit_logs_ip_address ON audit_logs(ip_address);
CREATE INDEX IF NOT EXISTS idx_audit_logs_device_id ON audit_logs(device_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_archived_at ON audit_logs(archived_at) WHERE archived_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_resource ON audit_logs(tenant_id, resource_type, resource_id);

-- コメント追加
COMMENT ON COLUMN audit_logs.ip_address IS 'アクセス元IPアドレス（IPv4/IPv6対応）';
COMMENT ON COLUMN audit_logs.user_agent IS 'ユーザーエージェント（ブラウザ・デバイス情報）';
COMMENT ON COLUMN audit_logs.device_id IS 'デバイスID（デバイス認証用）';
COMMENT ON COLUMN audit_logs.old_values IS '変更前の値（JSON形式、UPDATE/DELETE時）';
COMMENT ON COLUMN audit_logs.new_values IS '変更後の値（JSON形式、CREATE/UPDATE時）';
COMMENT ON COLUMN audit_logs.change_summary IS '変更サマリー（人間が読みやすい形式）';
COMMENT ON COLUMN audit_logs.log_hash IS 'ログのハッシュ値（SHA-256、改ざん検知用）';
COMMENT ON COLUMN audit_logs.archived_at IS 'アーカイブ日時（1年以上のログをWORMストレージへ移行）';
COMMENT ON COLUMN audit_logs.archive_location IS 'アーカイブ先（Cloud Storageパス）';

-- アーカイブログテーブル（WORMストレージに移行したログのメタデータ）
CREATE TABLE IF NOT EXISTS audit_log_archives (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    archive_period_start TIMESTAMPTZ NOT NULL, -- アーカイブ期間（開始）
    archive_period_end TIMESTAMPTZ NOT NULL, -- アーカイブ期間（終了）
    log_count BIGINT NOT NULL, -- アーカイブしたログ数
    archive_location TEXT NOT NULL, -- Cloud Storageパス
    archive_hash VARCHAR(255) NOT NULL, -- アーカイブファイルのハッシュ値
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT audit_log_archives_tenant_fk FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_audit_log_archives_tenant_period ON audit_log_archives(tenant_id, archive_period_start, archive_period_end);

-- コメント
COMMENT ON TABLE audit_log_archives IS '監査ログアーカイブメタデータ（WORMストレージに移行したログ）';
COMMENT ON COLUMN audit_log_archives.archive_hash IS 'アーカイブファイル全体のハッシュ値（改ざん検知用）';

