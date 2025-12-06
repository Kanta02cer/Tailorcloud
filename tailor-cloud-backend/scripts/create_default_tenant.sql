-- デフォルトテナント作成スクリプト
-- 使用方法: psql -U tailorcloud -d tailorcloud -f scripts/create_default_tenant.sql

-- デフォルトテナントを作成（既に存在する場合はスキップ）
INSERT INTO tenants (
    id,
    name,
    type,
    legal_name,
    address,
    invoice_registration_no,
    tax_rounding_method,
    is_active,
    created_at,
    updated_at
) VALUES (
    '00000000-0000-0000-0000-000000000001', -- 固定UUID（環境変数で使用可能）
    'Default Tenant',
    'TAILOR',
    'デフォルトテナント株式会社',
    '東京都渋谷区...',
    NULL, -- インボイス登録番号（必要に応じて設定）
    'HALF_UP',
    TRUE,
    NOW(),
    NOW()
)
ON CONFLICT (id) DO NOTHING;

-- 作成されたテナントIDを表示
SELECT 
    id,
    name,
    type,
    legal_name,
    created_at
FROM tenants
WHERE id = '00000000-0000-0000-0000-000000000001';

-- 環境変数に設定する値
-- export DEFAULT_TENANT_ID="00000000-0000-0000-0000-000000000001"

