# 認証システム 動作確認ガイド

## 前提条件

1. PostgreSQLデータベースが起動している
2. データベーススキーマが作成されている
3. Firebaseプロジェクトが設定されている
4. 環境変数が適切に設定されている

## セットアップ手順

### 1. データベースの準備

```bash
# データベースに接続
psql -U tailorcloud -d tailorcloud

# または、接続スクリプトを使用
./scripts/check_postgres_connection.sh
```

### 2. デフォルトテナントの作成

```bash
# SQLスクリプトを実行
psql -U tailorcloud -d tailorcloud -f scripts/create_default_tenant.sql

# または、psql内で実行
\i scripts/create_default_tenant.sql
```

作成されるテナントID: `00000000-0000-0000-0000-000000000001`

### 3. 環境変数の設定

**バックエンド:**
```bash
export DEFAULT_TENANT_ID="00000000-0000-0000-0000-000000000001"
export GCP_PROJECT_ID="your-firebase-project-id"
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account-key.json"
export PORT="8080"
```

**Flutterアプリ:**
`tailor-cloud-app/config/development.env`に以下を設定:
```bash
ENABLE_FIREBASE=true
FIREBASE_API_KEY=your-api-key
FIREBASE_APP_ID=your-app-id
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_MESSAGING_SENDER_ID=your-sender-id
GOOGLE_WORKSPACE_DOMAIN=your-domain.com  # オプション
```

## 動作確認手順

### ステップ1: バックエンドの起動

```bash
cd tailor-cloud-backend

# 環境変数を設定
export DEFAULT_TENANT_ID="00000000-0000-0000-0000-000000000001"
export GCP_PROJECT_ID="your-firebase-project-id"

# バックエンドを起動
go run cmd/api/main.go
```

正常に起動すると、以下のログが表示されます:
```
User repository initialized
User service initialized
Auth handler initialized
Auth endpoints registered
TailorCloud Backend running on port 8080
```

### ステップ2: ヘルスチェック

```bash
curl http://localhost:8080/health
```

期待されるレスポンス: `OK`

### ステップ3: Flutterアプリの起動

```bash
cd tailor-cloud-app
./scripts/start_flutter.sh development chrome
```

### ステップ4: Googleサインインのテスト

1. Flutterアプリが起動したら、ログイン画面が表示されます
2. 「Google アカウントでログイン」ボタンをクリック
3. Googleアカウントを選択（Google Workspaceドメインが設定されている場合は、そのドメインのアカウントのみ）
4. ログイン成功後、メイン画面に遷移することを確認

### ステップ5: 認証エンドポイントの直接テスト

Firebase IDトークンを取得して、認証エンドポイントを直接テスト:

```bash
# Firebase IDトークンを取得（Flutterアプリのコンソールから、またはFirebase SDKで取得）
ID_TOKEN="your-firebase-id-token"

# 認証エンドポイントをテスト
./scripts/test_auth.sh "$ID_TOKEN"
```

期待されるレスポンス:
```json
{
  "user_id": "uuid-of-user",
  "email": "user@example.com",
  "tenant_id": "00000000-0000-0000-0000-000000000001",
  "role": "Staff",
  "verified": true
}
```

### ステップ6: データベースでの確認

```bash
# データベースに接続
psql -U tailorcloud -d tailorcloud

# 作成されたユーザーを確認
SELECT 
    id,
    email,
    name,
    role,
    tenant_id,
    firebase_uid,
    created_at,
    last_login_at
FROM users
ORDER BY created_at DESC
LIMIT 10;

# テナント情報を確認
SELECT * FROM tenants WHERE id = '00000000-0000-0000-0000-000000000001';
```

## トラブルシューティング

### エラー: "tenant not found"

**原因:** デフォルトテナントが作成されていない、または`DEFAULT_TENANT_ID`が正しく設定されていない

**解決方法:**
1. デフォルトテナントを作成:
   ```bash
   psql -U tailorcloud -d tailorcloud -f scripts/create_default_tenant.sql
   ```

2. 環境変数を確認:
   ```bash
   echo $DEFAULT_TENANT_ID
   ```

3. テナントが存在するか確認:
   ```sql
   SELECT * FROM tenants WHERE id = '00000000-0000-0000-0000-000000000001';
   ```

### エラー: "Firebase initialization failed"

**原因:** Firebase設定が不完全

**解決方法:**
1. `GCP_PROJECT_ID`が設定されているか確認
2. `GOOGLE_APPLICATION_CREDENTIALS`が正しいパスか確認
3. Firebase Consoleで認証プロバイダーが有効か確認

### エラー: "Invalid token"

**原因:** Firebase IDトークンが無効または期限切れ

**解決方法:**
1. Flutterアプリで再度ログインしてトークンを取得
2. トークンの有効期限を確認（通常1時間）
3. Firebase Consoleで認証設定を確認

### エラー: "Failed to create user"

**原因:** データベース接続エラーまたは制約違反

**解決方法:**
1. データベース接続を確認:
   ```bash
   ./scripts/check_postgres_connection.sh
   ```

2. ユーザーテーブルの制約を確認:
   ```sql
   SELECT * FROM users WHERE email = 'user@example.com';
   ```

3. ログを確認して詳細なエラーメッセージを確認

## 期待される動作

### 初回ログイン時

1. Googleサインインが成功
2. Firebase IDトークンが取得される
3. `POST /api/auth/verify`が呼び出される（内部的に）
4. ユーザーがデータベースに作成される
5. ユーザーID、テナントID、ロールが返される
6. メイン画面に遷移

### 2回目以降のログイン時

1. Googleサインインが成功
2. Firebase IDトークンが取得される
3. 既存のユーザー情報が取得される
4. `last_login_at`が更新される
5. メイン画面に遷移

## 関連ドキュメント

- [認証システム ドキュメント](./AUTHENTICATION.md)
- [環境設定ガイド](../tailor-cloud-app/ENVIRONMENT_SETUP.md)
- [Firebase設定ガイド](../tailor-cloud-app/FIREBASE_SETUP.md)

