# PostgreSQL セットアップガイド（macOS）

**作成日**: 2025-01  
**目的**: macOSでPostgreSQLをセットアップする手順

---

## 🔍 問題の確認

現在、`psql`コマンドが見つからない状態です。PostgreSQLがインストールされていないか、PATHに含まれていない可能性があります。

---

## ✅ 解決方法

### Option 1: HomebrewでPostgreSQLをインストール（推奨）

#### Step 1: Homebrewがインストールされているか確認

```bash
brew --version
```

Homebrewがインストールされていない場合は、以下を実行:
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

---

#### Step 2: PostgreSQLをインストール

```bash
# PostgreSQL 14をインストール
brew install postgresql@14

# または最新版
brew install postgresql
```

---

#### Step 3: PostgreSQLを起動

```bash
# PostgreSQLサービスを開始
brew services start postgresql@14

# または
brew services start postgresql
```

---

#### Step 4: PATHに追加

インストール後、PATHに追加する必要があります:

```bash
# HomebrewのPostgreSQLのパスを確認
brew --prefix postgresql@14
# または
brew --prefix postgresql

# 一時的にPATHに追加（現在のシェルセッションのみ）
export PATH="/opt/homebrew/opt/postgresql@14/bin:$PATH"
# または
export PATH="/usr/local/opt/postgresql@14/bin:$PATH"

# 永続的に追加する場合は、~/.zshrcに追加
echo 'export PATH="/opt/homebrew/opt/postgresql@14/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

---

#### Step 5: PostgreSQLの状態確認

```bash
# バージョン確認
psql --version

# PostgreSQLが起動しているか確認
brew services list | grep postgresql
```

---

### Option 2: PostgreSQL.appを使用（GUI）

macOS用のPostgreSQLアプリを使用することもできます:

1. [Postgres.app](https://postgresapp.com/) をダウンロード
2. アプリケーションにインストール
3. 起動すると自動的にPostgreSQLが起動

---

## 🗄️ データベースとユーザーの作成

### Step 1: デフォルトのpostgresユーザーで接続

```bash
# PostgreSQLに接続（パスワードなしの場合）
psql postgres

# または、パスワードが必要な場合
psql -U postgres -d postgres
```

---

### Step 2: ユーザーを作成

```sql
-- tailorcloudユーザーを作成
CREATE USER tailorcloud WITH PASSWORD 'your_password';

-- または、パスワードなしで作成（開発環境用）
CREATE USER tailorcloud;
```

---

### Step 3: データベースを作成

```sql
-- データベースを作成
CREATE DATABASE tailorcloud;

-- ユーザーに権限を付与
GRANT ALL PRIVILEGES ON DATABASE tailorcloud TO tailorcloud;

-- データベースに接続
\c tailorcloud

-- スキーマへの権限も付与
GRANT ALL ON SCHEMA public TO tailorcloud;
```

---

### Step 4: 接続確認

```bash
# 作成したユーザーで接続
psql -U tailorcloud -d tailorcloud -h localhost

# または、パスワードを環境変数で指定
export PGPASSWORD='your_password'
psql -U tailorcloud -d tailorcloud -h localhost
```

---

## 🔧 環境変数の設定

### .env.localファイルを更新

```bash
# .env.localファイルを編集
nano .env.local
```

以下のように設定:

```bash
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=tailorcloud
POSTGRES_PASSWORD=your_password
POSTGRES_DB=tailorcloud
```

---

## ✅ セットアップ確認

### 1. PostgreSQL接続確認スクリプトを実行

```bash
./scripts/check_postgres_connection.sh
```

### 2. マイグレーション実行

```bash
./scripts/run_migrations_suit_mbti.sh
```

---

## 🐛 トラブルシューティング

### 問題: psqlコマンドが見つからない

**解決方法**:
1. PostgreSQLがインストールされているか確認
2. PATHにPostgreSQLのbinディレクトリが含まれているか確認
3. `brew --prefix postgresql`でインストール先を確認

---

### 問題: PostgreSQLが起動しない

**解決方法**:
```bash
# サービスを再起動
brew services restart postgresql@14

# ログを確認
brew services info postgresql@14
```

---

### 問題: 接続が拒否される

**解決方法**:
1. PostgreSQLが起動しているか確認
2. ポート5432が使用されているか確認
3. ユーザーが存在するか確認
4. パスワードが正しいか確認

---

### 問題: 権限エラー

**解決方法**:
```sql
-- ユーザーに権限を付与
GRANT ALL PRIVILEGES ON DATABASE tailorcloud TO tailorcloud;

-- スキーマへの権限も付与
\c tailorcloud
GRANT ALL ON SCHEMA public TO tailorcloud;
```

---

## 📚 参考リンク

- [PostgreSQL公式サイト](https://www.postgresql.org/)
- [Postgres.app](https://postgresapp.com/)
- [Homebrew PostgreSQL](https://formulae.brew.sh/formula/postgresql)

---

**最終更新日**: 2025-01

