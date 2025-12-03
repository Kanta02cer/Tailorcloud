# TailorCloud: インストールガイド

**作成日**: 2025-01  
**目的**: システムの完全なセットアップ手順

---

## 📋 前提条件

### 必須ソフトウェア

1. **Go 1.24+** (バックエンドAPI用)
   ```bash
   # macOS (Homebrew)
   brew install go
   
   # 確認
   go version
   ```

2. **Flutter 3.16.0+** (モバイルアプリ用)
   ```bash
   # macOS (Homebrew)
   brew install --cask flutter
   
   # 確認
   flutter --version
   ```

3. **Node.js 18+** (Webアプリ用、オプション)
   ```bash
   # macOS (Homebrew)
   brew install node
   
   # 確認
   node --version
   npm --version
   ```

### オプションソフトウェア

- **PostgreSQL 17+** (オプション - Firestoreモードでも動作可能)
- **Docker** (オプション - コンテナ化デプロイ用)
- **Firebase CLI** (オプション - Firebase機能を使用する場合)

---

## 🚀 インストール手順

### ステップ1: リポジトリのクローン

```bash
git clone https://github.com/Kanta02cer/Tailorcloud.git
cd Tailorcloud
```

### ステップ2: 環境変数のセットアップ

```bash
# 環境変数テンプレートをコピー
cp .env.example .env.local

# 必要に応じて .env.local を編集
# デフォルト設定で動作する場合は編集不要
```

または、セットアップスクリプトを使用:

```bash
./scripts/setup_local_environment.sh
```

### ステップ3: バックエンド依存関係のインストール

```bash
cd tailor-cloud-backend
go mod download
go mod tidy
cd ..
```

### ステップ4: Flutterアプリ依存関係のインストール

```bash
cd tailor-cloud-app
flutter pub get

# コード生成（Freezed, Riverpod Generator）
flutter pub run build_runner build --delete-conflicting-outputs
cd ..
```

### ステップ5: Webアプリ依存関係のインストール（オプション）

```bash
cd suit-mbti-web-app
npm install
cd ..
```

### ステップ6: システム状態の確認

```bash
./scripts/check_system.sh
```

**確認項目**:
- ✅ Goがインストールされている
- ✅ Flutterがインストールされている
- ✅ 環境変数ファイルが存在する
- ✅ バックエンド・Flutterアプリのディレクトリが存在する

---

## 🗄️ データベースセットアップ（オプション）

### PostgreSQLを使用する場合

```bash
# PostgreSQLのセットアップ
./scripts/setup_postgresql_user_db.sh

# または PostgreSQL 17 を使用する場合
./scripts/create_database_postgresql17.sh

# マイグレーションの実行
cd tailor-cloud-backend
# マイグレーションファイルを手動で実行するか、
# データベース管理ツールを使用
```

### Firestoreを使用する場合

1. Firebase Consoleでプロジェクトを作成
2. Firestoreデータベースを作成
3. `.env.local` に以下を設定:
   ```bash
   GCP_PROJECT_ID=your-project-id
   GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json
   ```

---

## 🧪 動作確認

### バックエンドAPIの起動確認

```bash
# バックエンドを起動
./scripts/start_backend.sh

# 別のターミナルでヘルスチェック
curl http://localhost:8080/health
# 期待される結果: OK
```

### Flutterアプリの起動確認

```bash
# Flutterアプリを起動
./scripts/start_flutter.sh

# エミュレーターまたは実機でアプリが起動することを確認
```

### Webアプリの起動確認（オプション）

```bash
cd suit-mbti-web-app
npm run dev

# ブラウザで http://localhost:5173 にアクセス
```

---

## 🔧 トラブルシューティング

### Goの依存関係エラー

```bash
cd tailor-cloud-backend
go mod download
go mod tidy
```

### Flutterの依存関係エラー

```bash
cd tailor-cloud-app
flutter clean
flutter pub get
flutter pub run build_runner build --delete-conflicting-outputs
```

### 環境変数の読み込みエラー

```bash
# 環境変数ファイルを確認
cat .env.local

# 環境変数を手動で読み込む
export $(cat .env.local | grep -v '^#' | xargs)
```

### PostgreSQL接続エラー

```bash
# PostgreSQL接続を確認
./scripts/check_postgres_connection.sh

# PostgreSQLが起動していることを確認
# macOS
brew services list | grep postgresql
```

---

## 📚 次のステップ

インストールが完了したら、以下を参照してください:

- **[システム起動ガイド](./docs/67_System_Startup_Guide.md)**
- **[完全起動手順書](./docs/70_Complete_Startup_Guide.md)**
- **[APIリファレンス](./docs/73_API_Reference.md)**

---

**最終更新日**: 2025-01

