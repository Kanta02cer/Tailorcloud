-- ============================================================================
-- Suit-MBTI統合用テストデータ準備スクリプト
-- ============================================================================
-- 目的: 診断・予約機能のテスト用データを作成
-- ============================================================================

-- 注意: このスクリプトを実行する前に、以下が必要です:
-- 1. 既存のテストデータ（tenants, customers）が存在すること
-- 2. マイグレーション（013, 014, 015）が実行されていること

-- テスト用テナントの確認（存在しない場合は作成）
INSERT INTO tenants (id, type, legal_name, address, created_at, updated_at)
SELECT 'tenant_test_suit_mbti', 'Tailor', 'テストテーラー', '東京都渋谷区', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM tenants WHERE id = 'tenant_test_suit_mbti');

-- テスト用顧客の作成
INSERT INTO customers (id, tenant_id, name, email, phone, created_at, updated_at)
VALUES
  ('customer_test_001', 'tenant_test_suit_mbti', 'テスト顧客1', 'customer1@test.com', '090-1111-1111', NOW(), NOW()),
  ('customer_test_002', 'tenant_test_suit_mbti', 'テスト顧客2', 'customer2@test.com', '090-2222-2222', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- テスト用診断データの作成
INSERT INTO diagnoses (id, user_id, tenant_id, archetype, plan_type, diagnosis_result, created_at, updated_at)
VALUES
  (
    'diagnosis_test_001',
    'user_test_001',
    'tenant_test_suit_mbti',
    'Classic',
    'Best Value',
    '{"scores": {"classic": 85, "modern": 20, "elegant": 70, "sporty": 30, "casual": 45}, "recommendations": ["Classic", "Elegant"]}'::jsonb,
    NOW() - INTERVAL '7 days',
    NOW() - INTERVAL '7 days'
  ),
  (
    'diagnosis_test_002',
    'user_test_002',
    'tenant_test_suit_mbti',
    'Modern',
    'Authentic',
    '{"scores": {"classic": 30, "modern": 90, "elegant": 40, "sporty": 70, "casual": 60}, "recommendations": ["Modern", "Sporty"]}'::jsonb,
    NOW() - INTERVAL '5 days',
    NOW() - INTERVAL '5 days'
  ),
  (
    'diagnosis_test_003',
    'user_test_001',
    'tenant_test_suit_mbti',
    'Elegant',
    'Authentic',
    '{"scores": {"classic": 70, "modern": 30, "elegant": 95, "sporty": 20, "casual": 40}, "recommendations": ["Elegant"]}'::jsonb,
    NOW() - INTERVAL '3 days',
    NOW() - INTERVAL '3 days'
  )
ON CONFLICT (id) DO NOTHING;

-- テスト用予約データの作成
INSERT INTO appointments (id, user_id, tenant_id, fitter_id, appointment_datetime, duration_minutes, status, notes, created_at, updated_at)
VALUES
  (
    'appointment_test_001',
    'user_test_001',
    'tenant_test_suit_mbti',
    'fitter_test_001',
    NOW() + INTERVAL '7 days' + INTERVAL '14 hours', -- 7日後の14:00
    60,
    'Pending',
    '初回フィッティング、Classicスタイル希望',
    NOW() - INTERVAL '2 days',
    NOW() - INTERVAL '2 days'
  ),
  (
    'appointment_test_002',
    'user_test_002',
    'tenant_test_suit_mbti',
    'fitter_test_001',
    NOW() + INTERVAL '10 days' + INTERVAL '15 hours', -- 10日後の15:00
    90,
    'Confirmed',
    '再フィッティング、Modernスタイル',
    NOW() - INTERVAL '1 day',
    NOW() - INTERVAL '1 day'
  ),
  (
    'appointment_test_003',
    'user_test_001',
    'tenant_test_suit_mbti',
    'fitter_test_002',
    NOW() + INTERVAL '14 days' + INTERVAL '16 hours', -- 14日後の16:00
    60,
    'Pending',
    'Elegantスタイル希望',
    NOW(),
    NOW()
  )
ON CONFLICT (id) DO NOTHING;

-- 顧客テーブルの拡張フィールドを更新
UPDATE customers
SET 
  occupation = '会社員',
  annual_income_range = '500-1000万円',
  preferred_archetype = 'Classic',
  diagnosis_count = 2
WHERE id = 'customer_test_001';

UPDATE customers
SET 
  occupation = '自営業',
  annual_income_range = '1000-1500万円',
  preferred_archetype = 'Modern',
  diagnosis_count = 1
WHERE id = 'customer_test_002';

-- テストデータ確認用クエリ
-- SELECT * FROM diagnoses WHERE tenant_id = 'tenant_test_suit_mbti';
-- SELECT * FROM appointments WHERE tenant_id = 'tenant_test_suit_mbti';
-- SELECT * FROM customers WHERE tenant_id = 'tenant_test_suit_mbti';

