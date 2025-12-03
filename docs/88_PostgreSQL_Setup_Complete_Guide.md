# PostgreSQL セットアップ完了ガイド

**作成日**: 2025-01  
**状況**: PostgreSQL 17を使用中

---

## 🔍 現在の状況

- ✅ PostgreSQL 17が起動中
- ✅ ユーザー `tailorcloud` は作成済み
- ⚠️  データベース `tailorcloud` の作成と権限設定が必要

---

## ✅ セットアップ手順

### Step 1: PostgreSQL 17に接続

```bash
/Library/PostgreSQL/17/bin/psql -d postgres -U postgres
```

パスワードを入力すると、`postgres=#` プロンプトが表示されます。

---

### Step 2: データベースを作成

psql内で以下を実行:

```sql
-- postgresデータベースに接続（既に接続している場合は不要）
\c postgres

-- データベースを作成
CREATE DATABASE tailorcloud OWNER tailorcloud;
```

---

### Step 3: データベースに接続して権限を付与

```sql
-- tailorcloudデータベースに接続
\c tailorcloud

-- 権限を付与
GRANT ALL PRIVILEGES ON DATABASE tailorcloud TO tailorcloud;
GRANT ALL ON SCHEMA public TO tailorcloud;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO tailorcloud;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO tailorcloud;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO tailorcloud;
```

---

### Step 4: 確認

```sql
-- 現在のデータベースとユーザーを確認
SELECT current_database(), current_user;

-- データベース一覧を確認
\l

-- ユーザー一覧を確認
\du
```

---

### Step 5: psqlを終了

```sql
\q
```

---

## 🎯 完全なSQLコマンドセット

psqlに接続後、以下を順番に実行:

```sql
\c postgres
CREATE DATABASE tailorcloud OWNER tailorcloud;
\c tailorcloud
GRANT ALL PRIVILEGES ON DATABASE tailorcloud TO tailorcloud;
GRANT ALL ON SCHEMA public TO tailorcloud;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO tailorcloud;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO tailorcloud;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO tailorcloud;
SELECT current_database(), current_user;
```

---

## 📝 環境変数の設定

`.env.local`ファイルに以下を設定:

```bash
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=tailorcloud
POSTGRES_PASSWORD=tailorcloud_dev_password
POSTGRES_DB=tailorcloud
```

---

## ✅ 次のステップ

セットアップ完了後:

1. マイグレーション実行: `./scripts/run_migrations_suit_mbti.sh`
2. テストデータ準備: `psql -f scripts/prepare_test_data_suit_mbti.sql`
3. バックエンドサーバー起動: `./scripts/start_backend.sh`

---

**最終更新日**: 2025-01

