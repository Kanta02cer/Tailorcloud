-- TailorCloud: テストデータ準備スクリプト
-- このスクリプトは、ローカル開発環境でテスト用のデータを準備します

-- 注意: このスクリプトはPostgreSQLが利用可能な場合にのみ実行してください
-- PostgreSQLが利用できない場合は、Firestoreモードで動作します

-- テナント作成
INSERT INTO tenants (id, name, invoice_registration_no, tax_rounding_method, created_at, updated_at)
VALUES 
  ('tenant-123', 'Regalis Yotsuya Salon', 'T1234567890123', 'HALF_UP', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 顧客作成
INSERT INTO customers (id, tenant_id, name, name_kana, phone, email, created_at, updated_at)
VALUES 
  ('customer-001', 'tenant-123', '田中 太郎', 'タナカ タロウ', '090-1234-5678', 'tanaka@example.com', NOW(), NOW()),
  ('customer-002', 'tenant-123', '佐藤 花子', 'サトウ ハナコ', '090-2345-6789', 'sato@example.com', NOW(), NOW()),
  ('customer-003', 'tenant-123', '鈴木 一郎', 'スズキ イチロウ', '090-3456-7890', 'suzuki@example.com', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 生地作成
INSERT INTO fabrics (id, tenant_id, brand, name, sku, color, pattern, price_per_meter, stock_quantity, status, created_at, updated_at)
VALUES 
  ('fabric-001', 'tenant-123', 'V.B.C', 'Perennial Navy', 'VBC-001-NV', 'Navy', 'Solid', 12000, 150.0, 'Available', NOW(), NOW()),
  ('fabric-002', 'tenant-123', 'Zegna', 'Trofeo Charcoal Grid', 'ZEG-002-CG', 'Charcoal', 'Grid', 18000, 45.0, 'Available', NOW(), NOW()),
  ('fabric-003', 'tenant-123', 'Loro Piana', 'Tasmanian', 'LP-003-TA', 'Gray', 'Solid', 25000, 80.0, 'Available', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 反物（Fabric Roll）作成
INSERT INTO fabric_rolls (id, tenant_id, fabric_id, roll_number, initial_length, current_length, status, location, created_at, updated_at)
VALUES 
  ('roll-001', 'tenant-123', 'fabric-001', 'VBC-001-NV-001', 50.0, 50.0, 'Available', '倉庫A', NOW(), NOW()),
  ('roll-002', 'tenant-123', 'fabric-001', 'VBC-001-NV-002', 50.0, 50.0, 'Available', '倉庫A', NOW(), NOW()),
  ('roll-003', 'tenant-123', 'fabric-002', 'ZEG-002-CG-001', 30.0, 30.0, 'Available', '倉庫B', NOW(), NOW()),
  ('roll-004', 'tenant-123', 'fabric-003', 'LP-003-TA-001', 40.0, 40.0, 'Available', '倉庫B', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- アンバサダー作成
INSERT INTO ambassadors (id, tenant_id, user_id, name, commission_rate, created_at, updated_at)
VALUES 
  ('ambassador-001', 'tenant-123', 'user-001', 'アンバサダーA', 5.0, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

SELECT '✅ テストデータの準備が完了しました' AS status;

