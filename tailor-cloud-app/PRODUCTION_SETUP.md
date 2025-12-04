# TailorCloud 本番環境セットアップガイド

このドキュメントでは、TailorCloud Flutterアプリを本番環境で動作させるための手順を説明します。

## 📋 前提条件

- Flutter SDK 3.2.0以上がインストールされていること
- バックエンドAPIサーバーが起動していること
- 本番環境のAPI URLが分かっていること

## 🔧 セットアップ手順

### 1. 環境変数ファイルの作成

`config/production.env.example`をコピーして`config/production.env`を作成し、実際の値を設定してください。

```bash
cd tailor-cloud-app
cp config/production.env.example config/production.env
```

### 2. 環境変数の設定

`config/production.env`ファイルを編集して、以下の値を設定します：

```bash
# 環境タイプ
ENV=production

# APIベースURL（本番環境のAPIサーバー）
API_BASE_URL=https://api.tailorcloud.com

# Firebase設定（使用する場合）
ENABLE_FIREBASE=false
FIREBASE_API_KEY=
FIREBASE_APP_ID=
FIREBASE_PROJECT_ID=
FIREBASE_MESSAGING_SENDER_ID=

# デフォルトテナントID
DEFAULT_TENANT_ID=tenant-production-001
```

**重要**: `config/production.env`ファイルには機密情報が含まれるため、Gitにコミットしないでください。`.gitignore`に追加されていることを確認してください。

### 3. Firebase設定（オプション）

Firebaseを使用する場合は、Firebase Consoleから以下を取得して設定してください：

1. Firebase Console (https://console.firebase.google.com/) にアクセス
2. プロジェクトを選択
3. プロジェクト設定 > 全般 から以下を取得：
   - Web API Key
   - App ID
   - Project ID
   - Messaging Sender ID

### 4. アプリの起動

#### 開発環境で起動（ホットリロード対応）

```bash
# プロジェクトルートから実行
./scripts/start_flutter.sh development chrome
```

#### 本番環境で起動

```bash
# プロジェクトルートから実行
./scripts/start_flutter.sh production chrome
```

#### Webサーバーとして起動（本番環境用）

```bash
# プロジェクトルートから実行
./scripts/start_flutter.sh production web-server
```

## 🏗️ 本番環境へのビルド

### Webアプリのビルド

```bash
cd tailor-cloud-app
./scripts/build_production.sh
```

ビルド成果物は`build/web`ディレクトリに出力されます。

### ビルド成果物のデプロイ

`build/web`ディレクトリの内容をWebサーバーにデプロイしてください。

#### Nginxの設定例

```nginx
server {
    listen 80;
    server_name tailorcloud.example.com;
    
    root /path/to/build/web;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # APIプロキシ設定
    location /api {
        proxy_pass https://api.tailorcloud.com;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 🔒 セキュリティ設定

### 環境変数の管理

- 本番環境の環境変数ファイルは機密情報を含むため、安全に管理してください
- CI/CDパイプラインを使用する場合は、環境変数をシークレットとして設定してください

### CORS設定

バックエンドAPIサーバーで、フロントエンドのドメインからのリクエストを許可するようにCORSを設定してください。

### HTTPSの使用

本番環境では必ずHTTPSを使用してください。

## 🐛 トラブルシューティング

### Firebase初期化エラー

Firebaseが無効な場合は、`ENABLE_FIREBASE=false`に設定してください。アプリはFirebaseなしで正常に動作します。

### API接続エラー

- バックエンドAPIサーバーが起動していることを確認してください
- `API_BASE_URL`が正しく設定されていることを確認してください
- ブラウザの開発者ツールでネットワークエラーを確認してください

### ビルドエラー

- Flutter SDKのバージョンを確認してください（3.2.0以上が必要）
- `flutter pub get`を実行して依存関係を更新してください
- `flutter clean`を実行してから再度ビルドしてください

## 📚 関連ドキュメント

- [環境設定ガイド](ENVIRONMENT_SETUP.md)
- [開発ガイド](README.md)
- [API仕様](../docs/API.md)

