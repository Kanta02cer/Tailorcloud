# 認証システム ドキュメント

## 概要

TailorCloudバックエンドは、Firebase Authenticationを使用した認証システムを実装しています。Google Workspaceとの互換性を保ちながら、マルチテナント対応のユーザー管理を行います。

## アーキテクチャ

### 認証フロー

1. **Flutter側（クライアント）**
   - Googleサインインを実行
   - Firebase IDトークンを取得
   - すべてのAPIリクエストに`Authorization: Bearer <token>`ヘッダーを自動付与

2. **バックエンド（サーバー）**
   - `FirebaseAuthMiddleware`でトークンを検証
   - ユーザー情報をコンテキストに注入
   - 各エンドポイントで認証済みユーザー情報を利用

3. **ユーザー永続化**
   - `POST /api/auth/verify`エンドポイントでトークン検証
   - Firebase UIDでユーザーを検索
   - 存在しない場合は自動的にユーザーを作成
   - テナントIDとロールを設定

## エンドポイント

### POST /api/auth/verify

Firebase IDトークンを検証し、ユーザー情報を取得または作成します。

**リクエスト:**
```json
{
  "id_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6..."
}
```

**レスポンス（成功）:**
```json
{
  "user_id": "uuid-of-user",
  "email": "user@example.com",
  "tenant_id": "tenant-uuid",
  "role": "Staff",
  "verified": true
}
```

**レスポンス（失敗）:**
```json
{
  "verified": false
}
```

## 環境変数

### バックエンド

- `DEFAULT_TENANT_ID`: デフォルトテナントID（新規ユーザー作成時に使用）
- `GCP_PROJECT_ID`: FirebaseプロジェクトID
- `GOOGLE_APPLICATION_CREDENTIALS`: Firebase認証情報ファイルのパス

### Flutterアプリ

- `ENABLE_FIREBASE`: Firebaseを有効にするか（`true`/`false`）
- `FIREBASE_API_KEY`: Firebase Web API Key
- `FIREBASE_APP_ID`: Firebase App ID
- `FIREBASE_PROJECT_ID`: Firebase Project ID
- `FIREBASE_MESSAGING_SENDER_ID`: Firebase Messaging Sender ID
- `GOOGLE_WORKSPACE_DOMAIN`: 許可するGoogle Workspaceドメイン（オプション）

## ユーザーモデル

```go
type User struct {
    ID        string    `json:"id"`
    TenantID  string    `json:"tenant_id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Role      UserRole  `json:"role"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### ユーザーロール

- `Owner`: 契約、決済、全データ閲覧権限
- `Staff`: 接客、発注作成権限
- `Factory_Manager`: 受注承認、工程管理権限
- `Worker`: 作業完了チェックのみ

## テナント管理

### デフォルトテナント

新規ユーザー作成時、テナントIDが指定されていない場合は、環境変数`DEFAULT_TENANT_ID`から取得します。

### 将来の拡張

- メールドメインからテナントを自動決定
- テナントリポジトリに`GetByDomain`メソッドを実装
- Google Workspaceドメインとテナントのマッピング

## セキュリティ

### トークン検証

- Firebase IDトークンはサーバー側で検証
- トークンの有効期限をチェック
- 署名を検証して改ざんを防止

### マルチテナント分離

- すべてのデータは`tenant_id`で分離
- ユーザーは自分のテナントのデータのみアクセス可能
- RBAC（Role-Based Access Control）で権限管理

## トラブルシューティング

### ユーザー作成に失敗する

1. テナントが存在するか確認
   ```sql
   SELECT * FROM tenants WHERE id = 'your-tenant-id';
   ```

2. 環境変数`DEFAULT_TENANT_ID`が設定されているか確認
   ```bash
   echo $DEFAULT_TENANT_ID
   ```

3. データベース接続を確認
   ```bash
   ./scripts/check_postgres_connection.sh
   ```

### トークン検証に失敗する

1. Firebase設定が正しいか確認
   - `GCP_PROJECT_ID`が設定されているか
   - `GOOGLE_APPLICATION_CREDENTIALS`が正しいパスか

2. Firebase Consoleで認証プロバイダーが有効か確認
   - Google認証プロバイダーが有効
   - 承認済みドメインが設定されている

## 開発時の注意事項

### 開発環境

- `FirebaseAuthMiddleware.OptionalAuth`を使用（認証が失敗しても通す）
- 本番環境では`FirebaseAuthMiddleware.Authenticate`を使用

### テスト

```bash
# バックエンドを起動
go run cmd/api/main.go

# 別ターミナルで認証エンドポイントをテスト
curl -X POST http://localhost:8080/api/auth/verify \
  -H "Content-Type: application/json" \
  -d '{"id_token": "your-firebase-id-token"}'
```

## 関連ファイル

- `internal/handler/auth_handler.go`: 認証ハンドラー
- `internal/service/user_service.go`: ユーザーサービス
- `internal/repository/user_repository.go`: ユーザーリポジトリ
- `internal/middleware/auth.go`: 認証ミドルウェア

