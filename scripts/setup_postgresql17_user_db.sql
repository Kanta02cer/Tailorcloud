-- PostgreSQL 17: tailorcloudユーザーとデータベース作成スクリプト
-- 使用方法: PostgreSQL 17に接続後、このファイルを読み込む
-- \i scripts/setup_postgresql17_user_db.sql
-- または、psql -d postgres -U postgres -f scripts/setup_postgresql17_user_db.sql

-- ユーザーの存在確認と作成（DOブロックは使える）
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'tailorcloud') THEN
        CREATE USER tailorcloud WITH PASSWORD 'tailorcloud_dev_password';
        RAISE NOTICE 'ユーザー tailorcloud を作成しました';
    ELSE
        RAISE NOTICE 'ユーザー tailorcloud は既に存在します';
    END IF;
END
$$;

-- データベースの作成（DOブロック内では実行できないため、直接実行）
-- 既に存在する場合はエラーになるので、エラーを無視
CREATE DATABASE tailorcloud OWNER tailorcloud;

-- 権限の付与
GRANT ALL PRIVILEGES ON DATABASE tailorcloud TO tailorcloud;

-- tailorcloudデータベースに接続して権限を付与
\c tailorcloud

-- スキーマへの権限を付与
GRANT ALL ON SCHEMA public TO tailorcloud;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO tailorcloud;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO tailorcloud;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO tailorcloud;

-- 接続確認
SELECT 'データベースとユーザーのセットアップが完了しました！' AS status;
SELECT current_database() AS current_database, current_user AS current_user;

