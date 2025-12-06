# 環境変数設定ガイド

**作成日**: 2025-01  
**目的**: TailorCloud Flutterアプリの環境変数設定方法

---

## 📋 環境変数の種類

### 必須設定

#### `ENV`
- **説明**: 実行環境（development, staging, production）
- **デフォルト値**: `development`
- **例**: `production`, `staging`, `development`

#### `API_BASE_URL`
- **説明**: バックエンドAPIのベースURL
- **デフォルト値**: `http://localhost:8080`
- **例**: 
  - 開発環境: `http://localhost:8080`
  - 本番環境: `https://api.tailorcloud.com`

### オプション設定

#### `ENABLE_FIREBASE`
- **説明**: Firebaseを有効にするかどうか
- **デフォルト値**: `false`
- **例**: `true`, `false`

#### `FIREBASE_API_KEY`
- **説明**: Firebase Web API Key（Firebase有効時のみ必要）
- **デフォルト値**: 空文字列
- **例**: `AIzaSyXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX`

#### `FIREBASE_APP_ID`
- **説明**: Firebase App ID（Firebase有効時のみ必要）
- **デフォルト値**: 空文字列
- **例**: `1:123456789:web:abcdef123456`

#### `FIREBASE_PROJECT_ID`
- **説明**: Firebase Project ID（Firebase有効時のみ必要）
- **デフォルト値**: 空文字列
- **例**: `tailorcloud-production`

#### `FIREBASE_MESSAGING_SENDER_ID`
- **説明**: Firebase Messaging Sender ID（Firebase有効時のみ必要）
- **デフォルト値**: 空文字列
- **例**: `123456789012`

#### `DEFAULT_TENANT_ID`
- **説明**: デフォルトテナントID
- **デフォルト値**: `tenant-123`
- **例**: `tenant-production-001`

---

## 🔧 設定方法

### 方法1: ビルド時設定（推奨）

`flutter run` または `flutter build` コマンドで `--dart-define` オプションを使用：

```bash
# 開発環境
flutter run -d chrome

# 本番環境（Firebase無効）
flutter run -d chrome \
  --dart-define=ENV=production \
  --dart-define=API_BASE_URL=https://api.tailorcloud.com

# 本番環境（Firebase有効）
flutter run -d chrome \
  --dart-define=ENV=production \
  --dart-define=API_BASE_URL=https://api.tailorcloud.com \
  --dart-define=ENABLE_FIREBASE=true \
  --dart-define=FIREBASE_API_KEY=AIzaSyXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX \
  --dart-define=FIREBASE_APP_ID=1:123456789:web:abcdef123456 \
  --dart-define=FIREBASE_PROJECT_ID=tailorcloud-production \
  --dart-define=FIREBASE_MESSAGING_SENDER_ID=123456789012
```

### 方法2: スクリプトで設定

`scripts/start_flutter.sh` を編集して環境変数を設定：

```bash
#!/bin/bash
flutter run -d chrome \
  --dart-define=ENV=production \
  --dart-define=API_BASE_URL=https://api.tailorcloud.com \
  --dart-define=ENABLE_FIREBASE=false
```

### 方法3: VS Code の launch.json

`.vscode/launch.json` を作成：

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Flutter (Development)",
      "type": "dart",
      "request": "launch",
      "program": "lib/main.dart",
      "args": [
        "--dart-define=ENV=development",
        "--dart-define=API_BASE_URL=http://localhost:8080",
        "--dart-define=ENABLE_FIREBASE=false"
      ]
    },
    {
      "name": "Flutter (Production)",
      "type": "dart",
      "request": "launch",
      "program": "lib/main.dart",
      "args": [
        "--dart-define=ENV=production",
        "--dart-define=API_BASE_URL=https://api.tailorcloud.com",
        "--dart-define=ENABLE_FIREBASE=false"
      ]
    }
  ]
}
```

---

## 🚀 デプロイ時の設定

### Web（Firebase Hosting）

`firebase.json` でビルドコマンドを設定：

```json
{
  "hosting": {
    "public": "build/web",
    "ignore": [
      "firebase.json",
      "**/.*",
      "**/node_modules/**"
    ],
    "rewrites": [
      {
        "source": "**",
        "destination": "/index.html"
      }
    ],
    "predeploy": [
      "flutter build web --dart-define=ENV=production --dart-define=API_BASE_URL=https://api.tailorcloud.com"
    ]
  }
}
```

### iOS/Android

ビルド時に環境変数を設定：

```bash
# iOS
flutter build ios \
  --dart-define=ENV=production \
  --dart-define=API_BASE_URL=https://api.tailorcloud.com

# Android
flutter build apk \
  --dart-define=ENV=production \
  --dart-define=API_BASE_URL=https://api.tailorcloud.com
```

---

## 🔍 環境変数の確認

アプリ起動時に環境変数が正しく設定されているか確認するには、デバッグログを確認してください。

開発環境では、アプリ起動時に以下のようなログが出力されます：

```
[INFO] Starting TailorCloud App
[DEBUG] Environment: production
[DEBUG] API Base URL: https://api.tailorcloud.com
[DEBUG] Firebase Enabled: false
[DEBUG] Debug Logging: false
```

---

## ⚠️ 注意事項

1. **Firebase設定**: Firebaseを使用する場合は、すべてのFirebase関連環境変数を設定する必要があります。
2. **セキュリティ**: 本番環境のAPIキーや認証情報は、Gitリポジトリにコミットしないでください。
3. **デフォルト値**: 環境変数が設定されていない場合は、デフォルト値が使用されます。

---

## 📝 環境別設定例

### 開発環境（ローカル）

```bash
flutter run -d chrome \
  --dart-define=ENV=development \
  --dart-define=API_BASE_URL=http://localhost:8080 \
  --dart-define=ENABLE_FIREBASE=false
```

### ステージング環境

```bash
flutter run -d chrome \
  --dart-define=ENV=staging \
  --dart-define=API_BASE_URL=https://staging-api.tailorcloud.com \
  --dart-define=ENABLE_FIREBASE=false
```

### 本番環境

```bash
flutter build web \
  --dart-define=ENV=production \
  --dart-define=API_BASE_URL=https://api.tailorcloud.com \
  --dart-define=ENABLE_FIREBASE=false
```

