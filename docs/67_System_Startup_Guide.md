# TailorCloud: システム起動ガイド

**作成日**: 2025-01  
**目的**: シード調達用MVPのデモ環境構築

---

## 📋 前提条件

### 必要なソフトウェア

1. **Go 1.21+** (バックエンドAPI用)
   ```bash
   go version  # インストール確認
   ```

2. **Flutter 3.16.0+** (モバイルアプリ用)
   ```bash
   flutter --version  # インストール確認
   ```

3. **PostgreSQL** (オプション - なくても動作可能)
   - Firestoreモードで動作するため、必須ではありません

### オプション

- **Firebase/Firestore** (オプション - なくても動作可能)
- **Docker** (オプション - コンテナ化デプロイ用)

---

## 🚀 クイックスタート

### 1. システム状態の確認

```bash
cd /Users/wantan/teiloroud-ERPSystem
./scripts/check_system.sh
```

このスクリプトで以下を確認します:
- 必要なコマンドのインストール状態
- 環境変数ファイルの有無
- バックエンド・Flutterアプリの状態

### 2. 環境変数のセットアップ

```bash
./scripts/setup_local_environment.sh
```

これにより `.env.local` ファイルが作成されます。

### 3. バックエンドAPIの起動

**ターミナル1**:

```bash
./scripts/start_backend.sh
```

または手動で:

```bash
cd tailor-cloud-backend
export PORT=8080
go run cmd/api/main.go
```

**確認**:
- http://localhost:8080/health にアクセス
- "OK" が返ってくれば正常

### 4. Flutterアプリの起動

**ターミナル2**:

```bash
./scripts/start_flutter.sh
```

または手動で:

```bash
cd tailor-cloud-app
export API_BASE_URL=http://localhost:8080
flutter run --dart-define=API_BASE_URL=http://localhost:8080
```

---

## 🗄️ テストデータの準備（オプション）

PostgreSQLを使用する場合:

```bash
# PostgreSQLに接続
psql -h localhost -U tailorcloud -d tailorcloud

# テストデータを投入
\i scripts/prepare_test_data.sql
```

または、SQLファイルを直接実行:

```bash
psql -h localhost -U tailorcloud -d tailorcloud -f scripts/prepare_test_data.sql
```

### テストデータの内容

- **テナント**: `tenant-123` (Regalis Yotsuya Salon)
- **顧客**: 3名（田中太郎、佐藤花子、鈴木一郎）
- **生地**: 3種類（V.B.C、Zegna、Loro Piana）
- **反物**: 4ロール

---

## 📱 デモ手順

### 1. アプリ起動後の画面

1. **ホーム画面** が表示される
2. 右上の **「クイック発注」** ボタンをタップ

### 2. クイック発注画面での操作

#### ステップ1: 顧客選択
- 既存顧客をドロップダウンから選択
- または「新規顧客を登録」ボタンで新規登録

#### ステップ2: 生地選択
- 生地一覧からドロップダウンで選択

#### ステップ3: 金額・納期入力
- 金額を入力（例: `135000`）
- 納期をカレンダーから選択

#### 発注書生成
- **「発注書を生成（フリーランス保護法対応）」** ボタンをタップ
- 注文が作成され、発注書PDFが生成される
- 成功メッセージが表示される

---

## ⚠️ トラブルシューティング

### バックエンドAPIが起動しない

**問題**: `Failed to connect to PostgreSQL`

**解決策**:
- PostgreSQLはオプションです
- 警告メッセージが出ても、Firestoreモードで動作します
- または、PostgreSQLを起動するか、環境変数を設定してください

**問題**: `Failed to initialize Firebase`

**解決策**:
- Firebaseはオプションです
- 警告メッセージが出ても動作します
- デモ用には不要です

### Flutterアプリが起動しない

**問題**: `API_BASE_URL not found`

**解決策**:
```bash
export API_BASE_URL=http://localhost:8080
flutter run --dart-define=API_BASE_URL=http://localhost:8080
```

**問題**: `Failed to connect to API`

**解決策**:
1. バックエンドAPIが起動しているか確認
2. `http://localhost:8080/health` にアクセスして確認
3. ファイアウォール設定を確認

### テストデータが表示されない

**問題**: 顧客や生地が表示されない

**解決策**:
1. PostgreSQLを使用している場合、テストデータを投入してください
2. または、アプリ内で新規顧客・生地を登録してください

---

## 🔍 システム状態の確認

### バックエンドAPI

```bash
# ヘルスチェック
curl http://localhost:8080/health

# メトリクス（監視用）
curl http://localhost:8080/api/metrics
```

### Flutterアプリ

- アプリ内のログを確認
- デバッグコンソールでエラーメッセージを確認

---

## 📊 デモ用の推奨設定

### 最小構成（デモ用）

- ✅ バックエンドAPI起動（PostgreSQL不要）
- ✅ Flutterアプリ起動
- ✅ アプリ内で顧客・生地を登録
- ✅ 発注書生成機能のデモ

### 完全構成（本番準備）

- ✅ PostgreSQL起動・マイグレーション実行
- ✅ テストデータ投入
- ✅ Firebase設定（認証・ストレージ）
- ✅ 完全な機能デモ

---

## 🎯 デモ準備チェックリスト

### 事前準備

- [ ] システム状態確認スクリプトを実行
- [ ] 環境変数ファイルを作成
- [ ] バックエンドAPIを起動
- [ ] ヘルスチェックで動作確認
- [ ] Flutterアプリを起動
- [ ] アプリが正常に起動することを確認

### デモデータ準備

- [ ] テスト顧客を登録（またはアプリ内で登録）
- [ ] テスト生地を登録（またはアプリ内で登録）
- [ ] 発注書生成が動作することを確認

### デモ実行

- [ ] ホーム画面の表示確認
- [ ] クイック発注画面への遷移確認
- [ ] 顧客選択の動作確認
- [ ] 生地選択の動作確認
- [ ] 金額・納期入力の動作確認
- [ ] 発注書生成の動作確認
- [ ] 成功メッセージの表示確認

---

**最終更新日**: 2025-01  
**ステータス**: ✅ 準備完了

