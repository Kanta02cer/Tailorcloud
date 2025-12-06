# クイックスタートガイド

## 認証システムの動作確認

### 前提条件

1. PostgreSQLが起動している
2. データベーススキーマが作成されている
3. Firebaseプロジェクトが設定されている

### ステップ1: データベースの準備

```bash
cd tailor-cloud-backend

# セットアップスクリプトを実行
./scripts/setup_auth.sh
```

このスクリプトは以下を実行します:
- データベース接続の確認
- デフォルトテナントの作成
- 環境変数の確認

### ステップ2: 環境変数の設定

```bash
# デフォルトテナントIDを設定
export DEFAULT_TENANT_ID="00000000-0000-0000-0000-000000000001"

# Firebase設定（必要に応じて）
export GCP_PROJECT_ID="your-firebase-project-id"
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account-key.json"
```

### ステップ3: バックエンドの起動

```bash
cd tailor-cloud-backend
go run cmd/api/main.go
```

正常に起動すると、以下のようなログが表示されます:
```
User repository initialized
User service initialized
Auth handler initialized
Auth endpoints registered
TailorCloud Backend running on port 8080
```

### ステップ4: ヘルスチェック

別のターミナルで:

```bash
curl http://localhost:8080/health
```

期待されるレスポンス: `OK`

### ステップ5: Flutterアプリの起動

```bash
cd tailor-cloud-app
./scripts/start_flutter.sh development chrome
```

### ステップ6: Googleサインインのテスト

1. ブラウザでFlutterアプリが開きます
2. ログイン画面で「Google アカウントでログイン」をクリック
3. Googleアカウントを選択
4. ログイン成功後、メイン画面に遷移することを確認

### トラブルシューティング

#### データベース接続エラー

```bash
# データベース接続を確認
psql -U tailorcloud -d tailorcloud -c "SELECT 1;"

# 接続できない場合は、PostgreSQLが起動しているか確認
# macOSの場合:
brew services list | grep postgresql
```

#### テナントが見つからない

```bash
# テナントを手動で作成
psql -U tailorcloud -d tailorcloud -f scripts/create_default_tenant.sql

# テナントが存在するか確認
psql -U tailorcloud -d tailorcloud -c "SELECT * FROM tenants WHERE id = '00000000-0000-0000-0000-000000000001';"
```

#### Firebase初期化エラー

- `GCP_PROJECT_ID`が設定されているか確認
- `GOOGLE_APPLICATION_CREDENTIALS`が正しいパスか確認
- Firebase Consoleで認証プロバイダーが有効か確認

## 次のステップ

詳細なドキュメント:
- [認証システム ドキュメント](./docs/AUTHENTICATION.md)
- [動作確認ガイド](./docs/AUTH_TESTING.md)

