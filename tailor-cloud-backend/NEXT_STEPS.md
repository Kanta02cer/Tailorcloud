# 次のステップ - 動作確認ガイド

## 🎯 現在の状況

✅ **実装完了:**
- ユーザーリポジトリ・サービス
- 認証ハンドラー（`POST /api/auth/verify`）
- テナント管理
- エラーハンドリング

✅ **バックエンド:**
- ビルド成功
- PostgreSQLなしでもFirebase認証は動作（警告のみ）

✅ **Flutterアプリ:**
- Firebase設定済み（`development.env`）

## 🚀 動作確認手順

### ステップ1: バックエンドの起動

```bash
cd tailor-cloud-backend

# 開発環境用起動スクリプトを使用
./scripts/start_backend_dev.sh
```

または手動で:

```bash
export GCP_PROJECT_ID="regalis-erp"
export DEFAULT_TENANT_ID="00000000-0000-0000-0000-000000000001"
go run cmd/api/main.go
```

**注意:** PostgreSQLがない場合、以下の警告が表示されますが、Firebase認証は動作します:
```
WARNING: Failed to connect to PostgreSQL: ...
WARNING: Continuing with Firestore only mode
```

### ステップ2: ヘルスチェック

別のターミナルで:

```bash
curl http://localhost:8080/health
```

期待されるレスポンス: `OK` または `WARNING: PostgreSQL not connected`

### ステップ3: Flutterアプリの起動

```bash
cd tailor-cloud-app
./scripts/start_flutter.sh development chrome
```

### ステップ4: Googleサインインのテスト

1. ブラウザでFlutterアプリが開きます
2. ログイン画面で「Google アカウントでログイン」をクリック
3. Googleアカウントを選択
4. ログイン成功後、メイン画面に遷移することを確認

### ステップ5: 認証エンドポイントの直接テスト（オプション）

Firebase IDトークンを取得してテスト:

```bash
# FlutterアプリのコンソールからIDトークンを取得
# または、Firebase SDKで取得

./scripts/test_auth.sh <firebase-id-token>
```

## ⚠️ 注意事項

### PostgreSQLがない場合

- バックエンドは起動しますが、警告が表示されます
- Firebase認証（トークン検証）は動作します
- **ユーザー情報の永続化はできません**（PostgreSQLが必要）
- 初回ログイン時のユーザー作成は失敗しますが、トークン検証は成功します

### PostgreSQLがある場合

1. データベーススキーマを作成
2. `./scripts/setup_auth.sh` を実行してデフォルトテナントを作成
3. バックエンドを起動
4. 完全な認証フロー（ユーザー作成含む）が動作します

## 📊 動作確認チェックリスト

- [ ] バックエンドが起動する
- [ ] `/health`エンドポイントが応答する
- [ ] Flutterアプリが起動する
- [ ] Googleサインインボタンが表示される
- [ ] Googleサインインが成功する
- [ ] メイン画面に遷移する
- [ ] （PostgreSQLがある場合）ユーザーがDBに作成される

## 🔍 トラブルシューティング

### バックエンドが起動しない

- `GCP_PROJECT_ID`が設定されているか確認
- ポート8080が使用されていないか確認: `lsof -i :8080`

### Firebase認証エラー

- `development.env`のFirebase設定を確認
- Firebase Consoleで認証プロバイダーが有効か確認
- ブラウザのコンソールでエラーメッセージを確認

### ユーザー作成エラー（PostgreSQLがある場合）

- データベース接続を確認: `./scripts/check_postgres_connection.sh`
- デフォルトテナントが存在するか確認
- `DEFAULT_TENANT_ID`環境変数が正しく設定されているか確認

## 📚 関連ドキュメント

- [認証システム ドキュメント](./docs/AUTHENTICATION.md)
- [動作確認ガイド](./docs/AUTH_TESTING.md)
- [クイックスタートガイド](./QUICK_START.md)

