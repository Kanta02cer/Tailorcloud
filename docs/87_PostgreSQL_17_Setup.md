# PostgreSQL 17 セットアップガイド

**作成日**: 2025-01  
**状況**: PostgreSQL 17が既に起動しており、ポート5432を使用中

---

## 🔍 現在の状況

- ✅ PostgreSQL 17が起動中（`/Library/PostgreSQL/17/bin/postgres`）
- ✅ ポート5432で稼働中
- ❌ パスワード認証が必要

---

## ✅ データベースとユーザーの作成

### 方法1: postgresユーザーで接続（パスワードが必要）

PostgreSQL 17に接続するには、パスワードが必要です。

```bash
# PostgreSQL 17のpsqlを使用
/Library/PostgreSQL/17/bin/psql -d postgres -U postgres

# パスワードを入力後、以下を実行:
CREATE USER tailorcloud WITH PASSWORD 'your_password';
CREATE DATABASE tailorcloud OWNER tailorcloud;
GRANT ALL PRIVILEGES ON DATABASE tailorcloud TO tailorcloud;
```

---

### 方法2: 認証設定を変更（開発環境のみ）

PostgreSQL 17の認証設定を変更して、パスワードなしで接続できるようにします。

```bash
# pg_hba.confファイルを編集
sudo nano /Library/PostgreSQL/17/data/pg_hba.conf

# ファイルの最後に以下を追加:
# local   all             all                                     trust
# host    all             all             127.0.0.1/32            trust
# host    all             all             ::1/128                 trust

# PostgreSQL 17を再起動
# macOSの場合、PostgreSQL 17のサービス管理方法を確認してください
```

---

### 方法3: 環境変数でパスワードを設定

パスワードを知っている場合、環境変数で設定できます。

```bash
# .env.localファイルにパスワードを設定
export POSTGRES_PASSWORD='your_password'

# 接続
/Library/PostgreSQL/17/bin/psql -d postgres -U postgres -h localhost
```

---

## 🎯 推奨手順

1. **PostgreSQL 17のパスワードを確認またはリセット**
   - インストール時に設定したパスワードを確認
   - または、postgresユーザーのパスワードをリセット

2. **データベースとユーザーを作成**

3. **.env.localファイルを更新**
   ```bash
   POSTGRES_HOST=localhost
   POSTGRES_PORT=5432
   POSTGRES_USER=tailorcloud
   POSTGRES_PASSWORD=your_password
   POSTGRES_DB=tailorcloud
   ```

4. **マイグレーション実行**

---

## 📝 注意事項

- PostgreSQL 17は既に起動しているため、HomebrewのPostgreSQL@15は不要です
- PostgreSQL 17を使用する場合、環境変数やスクリプトでPostgreSQL 17のパスを指定してください
- 本番環境では必ずパスワード認証を使用してください

---

**最終更新日**: 2025-01

