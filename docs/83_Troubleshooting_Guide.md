# トラブルシューティングガイド: Suit-MBTI統合

**作成日**: 2025-01  
**目的**: よくある問題と解決方法

---

## 🔍 PostgreSQL接続エラー

### 問題: PostgreSQL接続失敗

**エラーメッセージ**:
```
❌ PostgreSQL接続失敗
接続情報を確認してください。
```

---

### 解決方法

#### 1. PostgreSQLが起動しているか確認

```bash
# PostgreSQLの状態を確認
pg_isready

# または
psql --version
```

**PostgreSQLが起動していない場合**:
```bash
# macOS (Homebrew)
brew services start postgresql@14
# または
brew services start postgresql

# Linux
sudo systemctl start postgresql
# または
sudo service postgresql start
```

---

#### 2. データベースが存在するか確認

```bash
# データベース一覧を表示
psql -h localhost -U postgres -l

# または、特定のデータベースに接続を試みる
psql -h localhost -U postgres -d tailorcloud -c "SELECT 1"
```

**データベースが存在しない場合**:
```bash
# データベースを作成
createdb -h localhost -U postgres tailorcloud

# または、psqlで作成
psql -h localhost -U postgres -c "CREATE DATABASE tailorcloud;"
```

---

#### 3. ユーザーが存在するか確認

```bash
# ユーザー一覧を表示
psql -h localhost -U postgres -d postgres -c "\du"
```

**ユーザーが存在しない場合**:
```bash
# ユーザーを作成
psql -h localhost -U postgres -d postgres -c "CREATE USER tailorcloud WITH PASSWORD 'your_password';"

# データベースへの権限を付与
psql -h localhost -U postgres -d postgres -c "GRANT ALL PRIVILEGES ON DATABASE tailorcloud TO tailorcloud;"
```

---

#### 4. 環境変数の確認

```bash
# 現在の環境変数を確認
echo "POSTGRES_HOST: $POSTGRES_HOST"
echo "POSTGRES_PORT: $POSTGRES_PORT"
echo "POSTGRES_USER: $POSTGRES_USER"
echo "POSTGRES_DB: $POSTGRES_DB"
```

**環境変数を設定**:
```bash
# .env.localファイルを読み込む
source .env.local

# または、直接設定
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=your_user
export POSTGRES_PASSWORD=your_password
export POSTGRES_DB=tailorcloud
```

---

#### 5. 接続テスト

```bash
# パスワードを環境変数に設定
export PGPASSWORD=$POSTGRES_PASSWORD

# 接続テスト
psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -c "SELECT 1"
```

---

### よくある接続エラーパターン

#### エラー: "connection refused"

**原因**: PostgreSQLが起動していない、またはポートが間違っている

**解決方法**:
```bash
# PostgreSQLが起動しているか確認
pg_isready -h localhost -p 5432

# 起動していない場合は起動
brew services start postgresql  # macOS
# または
sudo systemctl start postgresql  # Linux
```

---

#### エラー: "database does not exist"

**原因**: データベースが存在しない

**解決方法**:
```bash
# データベースを作成
createdb -h localhost -U postgres tailorcloud
```

---

#### エラー: "password authentication failed"

**原因**: パスワードが間違っている、またはユーザーが存在しない

**解決方法**:
```bash
# パスワードを確認・設定
export POSTGRES_PASSWORD='your_password'

# または、.env.localファイルに設定
echo "POSTGRES_PASSWORD=your_password" >> .env.local
```

---

#### エラー: "permission denied"

**原因**: ユーザーにデータベースへのアクセス権限がない

**解決方法**:
```bash
# 権限を付与
psql -h localhost -U postgres -d postgres -c "GRANT ALL PRIVILEGES ON DATABASE tailorcloud TO tailorcloud;"
```

---

## 🔧 マイグレーションエラー

### 問題: マイグレーション実行エラー

**エラーメッセージ**:
```
❌ 失敗
ERROR: relation "diagnoses" already exists
```

**解決方法**:
- テーブルが既に存在する場合は、マイグレーションはスキップされます（正常）
- エラーが出ても続行できるように、スクリプトで対応済み

---

### 問題: 外部キー制約エラー

**エラーメッセージ**:
```
ERROR: relation "tenants" does not exist
```

**解決方法**:
- `tenants`テーブルが存在しない場合は、先に既存のマイグレーションを実行してください
- または、外部キー制約を一時的にコメントアウトしてください

---

## 🐛 コンパイルエラー

### 問題: Goのコンパイルエラー

**エラーメッセージ**:
```
undefined: domain.Diagnosis
```

**解決方法**:
```bash
# モジュールを再ダウンロード
cd tailor-cloud-backend
go mod tidy
go mod download

# 再コンパイル
go build ./cmd/api/main.go
```

---

## 📝 よくある質問

### Q1: PostgreSQLがインストールされていない

**解決方法**:
```bash
# macOS (Homebrew)
brew install postgresql@14
brew services start postgresql@14

# Linux (Ubuntu/Debian)
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib
sudo systemctl start postgresql
```

---

### Q2: 環境変数が読み込まれない

**解決方法**:
```bash
# .env.localファイルを読み込む
source .env.local

# または、スクリプト内で読み込む
export $(cat .env.local | grep -v '^#' | xargs)
```

---

### Q3: マイグレーションが失敗する

**解決方法**:
1. PostgreSQLに接続できるか確認
2. データベースが存在するか確認
3. ユーザーに権限があるか確認
4. 既存のマイグレーションが実行されているか確認

---

## 📞 サポート

問題が解決しない場合は、以下の情報を含めて報告してください:

1. エラーメッセージ（全文）
2. 実行したコマンド
3. 環境情報:
   - OS
   - PostgreSQLバージョン
   - Goバージョン
   - 環境変数の設定（パスワードは除く）

---

**最終更新日**: 2025-01

