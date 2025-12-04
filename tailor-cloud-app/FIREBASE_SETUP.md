# Firebase セットアップガイド

このドキュメントでは、TailorCloud FlutterアプリでFirebaseを使用するためのセットアップ手順を説明します。

## 📋 前提条件

- Firebaseプロジェクトが作成されていること
- Firebase Consoleへのアクセス権限があること
- Flutter SDK 3.2.0以上がインストールされていること

## 🔧 Firebaseプロジェクトの作成

### 1. Firebase Consoleでプロジェクトを作成

1. [Firebase Console](https://console.firebase.google.com/) にアクセス
2. 「プロジェクトを追加」をクリック
3. プロジェクト名を入力（例: `tailorcloud-production`）
4. Google Analyticsの設定（オプション）
5. プロジェクトを作成

### 2. Webアプリを追加

1. Firebase Consoleでプロジェクトを選択
2. プロジェクト設定（⚙️アイコン）をクリック
3. 「アプリを追加」から「Web」を選択
4. アプリのニックネームを入力（例: `TailorCloud Web App`）
5. 「このアプリのFirebase Hostingも設定しますか？」は任意
6. 「アプリを登録」をクリック

### 3. Firebase設定情報を取得

プロジェクト設定 > 全般 > マイアプリから、以下の情報を取得します：

- **API Key**: `AIzaSyXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX`
- **App ID**: `1:123456789:web:abcdef123456`
- **Project ID**: `tailorcloud-production`
- **Messaging Sender ID**: `123456789012`

## 🔐 Firebase Authenticationの設定

### 1. Authenticationを有効化

1. Firebase Consoleで「Authentication」を選択
2. 「始める」をクリック
3. 「Sign-in method」タブを選択
4. 「メール/パスワード」を有効化

### 2. 認証ドメインの設定（Web環境）

1. Authentication > 設定 > 承認済みドメイン
2. アプリのドメインを追加（例: `tailorcloud.com`）
3. 開発環境では `localhost` が自動的に追加されています

## 📝 環境変数の設定

### 開発環境

`tailor-cloud-app/config/development.env` を編集：

```bash
ENV=development
API_BASE_URL=http://localhost:8080

# Firebase設定
ENABLE_FIREBASE=true
FIREBASE_API_KEY=AIzaSyXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
FIREBASE_APP_ID=1:123456789:web:abcdef123456
FIREBASE_PROJECT_ID=tailorcloud-production
FIREBASE_MESSAGING_SENDER_ID=123456789012

DEFAULT_TENANT_ID=tenant-123
```

### 本番環境

`tailor-cloud-app/config/production.env` を編集：

```bash
ENV=production
API_BASE_URL=https://api.tailorcloud.com

# Firebase設定
ENABLE_FIREBASE=true
FIREBASE_API_KEY=AIzaSyXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
FIREBASE_APP_ID=1:123456789:web:abcdef123456
FIREBASE_PROJECT_ID=tailorcloud-production
FIREBASE_MESSAGING_SENDER_ID=123456789012

DEFAULT_TENANT_ID=tenant-production-001
```

**重要**: 本番環境の設定ファイルには機密情報が含まれるため、Gitにコミットしないでください。

## 🚀 アプリの起動

### Firebase有効で起動

```bash
# 開発環境
./scripts/start_flutter.sh development chrome

# 本番環境
./scripts/start_flutter.sh production chrome
```

### Firebase無効で起動（デフォルト）

Firebase設定を省略すると、アプリはFirebaseなしで動作します：

```bash
# Firebase無効（デフォルト）
./scripts/start_flutter.sh development chrome
```

## 🔍 Firebase初期化の確認

アプリ起動時に、コンソールに以下のようなログが表示されます：

### Firebase有効時

```
[INFO] Starting TailorCloud App
[DEBUG] Environment: development
[DEBUG] API Base URL: http://localhost:8080
[DEBUG] Firebase Enabled: true
[INFO] Firebase: Initialized successfully.
```

### Firebase無効時

```
[INFO] Starting TailorCloud App
[DEBUG] Environment: development
[DEBUG] API Base URL: http://localhost:8080
[DEBUG] Firebase Enabled: false
[DEBUG] Firebase: Disabled. Skipping initialization.
```

### Firebase設定不完全時

```
[WARNING] Firebase: Configuration incomplete. Skipping initialization.
[DEBUG] Firebase: Required settings - API_KEY, APP_ID, PROJECT_ID
```

## 🧪 テストユーザーの作成

### Firebase Consoleから作成

1. Authentication > ユーザー
2. 「ユーザーを追加」をクリック
3. メールアドレスとパスワードを入力
4. 「ユーザーを追加」をクリック

### アプリから作成（開発用）

現在の実装では、Firebase Consoleからユーザーを作成する必要があります。将来的にユーザー登録機能を追加する予定です。

## 🔒 セキュリティ設定

### Firebase Security Rules

Firebase Authenticationを使用する場合、以下のセキュリティ設定を確認してください：

1. **Authentication > 設定 > 承認済みドメイン**: アプリのドメインのみを許可
2. **Firestore/Realtime Database**: 使用する場合は適切なセキュリティルールを設定

### 環境変数の管理

- 本番環境のFirebase設定は機密情報として扱う
- CI/CDパイプラインでは環境変数をシークレットとして設定
- `.gitignore`に環境変数ファイルが含まれていることを確認

## 🐛 トラブルシューティング

### Firebase初期化エラー

**エラー**: `FirebaseOptions cannot be null when creating the default app.`

**解決方法**:
1. 環境変数が正しく設定されているか確認
2. `ENABLE_FIREBASE=true` が設定されているか確認
3. すべてのFirebase設定（API_KEY, APP_ID, PROJECT_ID）が設定されているか確認

### 認証エラー

**エラー**: `Firebase Auth is not enabled`

**解決方法**:
1. `ENABLE_FIREBASE=true` を設定
2. Firebase設定が完全か確認
3. Firebase ConsoleでAuthenticationが有効になっているか確認

### ネットワークエラー

**エラー**: `ネットワークエラーが発生しました`

**解決方法**:
1. インターネット接続を確認
2. Firebase Consoleでプロジェクトが有効か確認
3. 承認済みドメインに現在のドメインが含まれているか確認

## 📚 関連ドキュメント

- [Firebase公式ドキュメント](https://firebase.google.com/docs)
- [FlutterFire公式ドキュメント](https://firebase.flutter.dev/)
- [環境変数設定ガイド](ENVIRONMENT_SETUP.md)
- [本番環境セットアップガイド](PRODUCTION_SETUP.md)

## 🔄 Firebase機能の追加

現在実装されているFirebase機能：

- ✅ Firebase Authentication（メール/パスワード）
- ✅ Firebase初期化と設定管理
- ✅ 認証状態の監視
- ✅ ログイン/ログアウト

将来追加予定の機能：

- ⏳ Firebase Cloud Messaging（プッシュ通知）
- ⏳ Firebase Storage（ファイルアップロード）
- ⏳ Firebase Firestore（リアルタイムデータベース）
- ⏳ Firebase Analytics（分析）

## 💡 ベストプラクティス

1. **環境別設定**: 開発環境と本番環境で異なるFirebaseプロジェクトを使用
2. **セキュリティ**: APIキーなどの機密情報は環境変数で管理
3. **エラーハンドリング**: Firebaseが無効でもアプリが動作するように実装
4. **ログ**: デバッグモードでFirebase初期化のログを確認
5. **テスト**: Firebase有効/無効の両方のシナリオをテスト

