# PostgreSQL接続問題の解決方法

**作成日**: 2025-01  
**問題**: PostgreSQL接続時にパスワード認証エラー

---

## 🔍 現在の状況

- ✅ PostgreSQL@15 インストール済み
- ✅ PATH設定済み（psqlコマンド使用可能）
- ❌ 接続時にパスワード認証エラー

---

## ✅ 解決方法

### 方法1: macOSユーザー名で接続（推奨）

HomebrewでインストールしたPostgreSQLは、通常macOSユーザー名でパスワードなしで接続できます。

```bash
# PATHを設定
source ./scripts/fix_postgresql_path.sh

# データベースを作成（現在のユーザー名で）
createdb tailorcloud

# 接続確認
psql tailorcloud -c "SELECT 1;"
```

---

### 方法2: PostgreSQLサービスを再初期化

PostgreSQLサービスが正しく起動していない可能性があります。

```bash
# サービスを停止
brew services stop postgresql@15

# データディレクトリを初期化（必要な場合）
initdb /opt/homebrew/var/postgresql@15

# サービスを再起動
brew services start postgresql@15
```

---

### 方法3: 認証設定を確認

PostgreSQLの認証設定ファイルを確認します。

```bash
# pg_hba.confファイルの場所を確認
psql -h localhost -U postgres -d postgres -c "SHOW hba_file;"

# 通常は以下の場所:
# /opt/homebrew/var/postgresql@15/pg_hba.conf

# ファイルを編集（trust方式に変更）
nano /opt/homebrew/var/postgresql@15/pg_hba.conf

# 以下の行を追加または変更:
# local   all             all                                     trust
# host    all             all             127.0.0.1/32            trust
# host    all             all             ::1/128                 trust

# PostgreSQLを再起動
brew services restart postgresql@15
```

---

### 方法4: パスワードを設定

PostgreSQLユーザーにパスワードを設定します。

```bash
# まず、postgresユーザーで接続（パスワードなしまたは設定済み）
psql postgres

# パスワードを設定
ALTER USER postgres WITH PASSWORD 'your_password';

# tailorcloudユーザーを作成
CREATE USER tailorcloud WITH PASSWORD 'your_password';

# データベースを作成
CREATE DATABASE tailorcloud OWNER tailorcloud;

# .env.localファイルにパスワードを設定
echo "POSTGRES_PASSWORD=your_password" >> .env.local
```

---

## 🎯 クイック解決手順

### Step 1: PATH設定

```bash
source ./scripts/fix_postgresql_path.sh
```

### Step 2: データベース作成を試行

```bash
createdb tailorcloud
```

### Step 3: 接続確認

```bash
psql tailorcloud -c "SELECT version();"
```

---

## 📞 サポート

問題が解決しない場合は、以下の情報を含めて報告してください:

1. PostgreSQLのバージョン: `psql --version`
2. サービス状態: `brew services list | grep postgresql`
3. エラーメッセージ全文
4. macOSバージョン

---

**最終更新日**: 2025-01

