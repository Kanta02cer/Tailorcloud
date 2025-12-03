# TailorCloud: 完全起動手順書

**作成日**: 2025-01  
**目的**: シード調達用MVPのデモ環境を完全構築

---

## 🎯 クイックスタート（5分で起動）

### 1. システム状態確認（1分）

```bash
cd /Users/wantan/teiloroud-ERPSystem
./scripts/check_system.sh
```

**確認項目**:
- ✅ Goがインストールされている
- ✅ Flutterがインストールされている
- ✅ 環境変数ファイルが存在する

### 2. バックエンドAPI起動（2分）

**ターミナル1**:

```bash
./scripts/start_backend.sh
```

**確認**:
- ターミナルに `TailorCloud Backend running on port 8080` と表示される
- ブラウザで `http://localhost:8080/health` にアクセスして "OK" が返ってくる

### 3. Flutterアプリ起動（2分）

**ターミナル2**:

```bash
./scripts/start_flutter.sh
```

**確認**:
- エミュレーターまたは実機でアプリが起動する
- ホーム画面が表示される

---

## 📋 詳細手順

### ステップ1: 環境変数の確認

```bash
# 環境変数ファイルを確認
cat .env.local
```

**重要な設定**:
- `PORT=8080` - バックエンドAPIのポート
- `API_BASE_URL=http://localhost:8080` - Flutterアプリからバックエンドへの接続URL

### ステップ2: バックエンドAPIの起動

#### 方法1: スクリプトを使用（推奨）

```bash
./scripts/start_backend.sh
```

#### 方法2: 手動で起動

```bash
cd tailor-cloud-backend
export PORT=8080
go run cmd/api/main.go
```

#### 起動確認

```bash
# ヘルスチェック
curl http://localhost:8080/health

# 期待される結果: OK
```

### ステップ3: Flutterアプリの起動

#### 方法1: スクリプトを使用（推奨）

```bash
./scripts/start_flutter.sh
```

#### 方法2: 手動で起動

```bash
cd tailor-cloud-app
export API_BASE_URL=http://localhost:8080
flutter run --dart-define=API_BASE_URL=http://localhost:8080
```

#### デバイス選択

Flutter起動時にデバイスを選択:
- iOSシミュレーター
- Androidエミュレーター
- 実機（USB接続）

---

## 🗄️ テストデータの準備（オプション）

### PostgreSQLを使用する場合

```bash
# PostgreSQLに接続
psql -h localhost -U tailorcloud -d tailorcloud

# テストデータを投入
\i scripts/prepare_test_data.sql
```

### PostgreSQLを使用しない場合

アプリ内で直接データを登録:
1. クイック発注画面を開く
2. 「新規顧客を登録」ボタンで顧客を登録
3. 生地は既存データを使用（またはAPI経由で登録）

---

## 🧪 デモの実行

### デモフロー（5分）

1. **ホーム画面の表示**（10秒）
   - ダッシュボードが表示される
   - 「クイック発注」ボタンを確認

2. **クイック発注画面への遷移**（10秒）
   - 「クイック発注」ボタンをタップ
   - クイック発注画面が開く

3. **ステップ1: 顧客選択**（30秒）
   - 既存顧客を選択
   - または「新規顧客を登録」で新規登録

4. **ステップ2: 生地選択**（30秒）
   - 生地一覧から選択

5. **ステップ3: 金額・納期入力**（1分）
   - 金額を入力（例: `135000`）
   - 納期をカレンダーから選択

6. **発注書生成**（1分）
   - 「発注書を生成（フリーランス保護法対応）」ボタンをタップ
   - ローディング表示
   - 成功メッセージが表示される

7. **質問タイム**（2分）
   - 質問に答える
   - LOI獲得の質問

---

## ⚠️ トラブルシューティング

### バックエンドAPIが起動しない

**エラー**: `Failed to connect to PostgreSQL`

**解決策**:
- PostgreSQLはオプションです
- 警告メッセージが出ても、Firestoreモードで動作します
- 起動は続行されます

**エラー**: `Failed to initialize Firebase`

**解決策**:
- Firebaseはオプションです
- 警告メッセージが出ても動作します
- デモ用には不要です

**エラー**: `Port 8080 is already in use`

**解決策**:
```bash
# ポート8080を使用しているプロセスを確認
lsof -i :8080

# プロセスを終了
kill -9 <PID>

# または別のポートを使用
export PORT=8081
```

### Flutterアプリが起動しない

**エラー**: `API_BASE_URL not found`

**解決策**:
```bash
export API_BASE_URL=http://localhost:8080
flutter run --dart-define=API_BASE_URL=http://localhost:8080
```

**エラー**: `Failed to connect to API`

**解決策**:
1. バックエンドAPIが起動しているか確認
   ```bash
   curl http://localhost:8080/health
   ```
2. 環境変数が正しく設定されているか確認
   ```bash
   echo $API_BASE_URL
   ```

**エラー**: `No devices found`

**解決策**:
1. エミュレーターを起動
   - iOS: Xcodeからシミュレーターを起動
   - Android: Android Studioからエミュレーターを起動
2. または実機を接続
   - USB接続で実機を接続
   - 開発者モードを有効化

### テストデータが表示されない

**問題**: 顧客や生地が表示されない

**解決策**:
1. PostgreSQLを使用している場合、テストデータを投入
   ```bash
   psql -h localhost -U tailorcloud -d tailorcloud -f scripts/prepare_test_data.sql
   ```
2. PostgreSQLを使用していない場合、アプリ内で直接登録

---

## 🎯 デモ準備チェックリスト

### 事前準備

- [ ] システム状態確認スクリプトを実行
- [ ] 環境変数ファイルが存在することを確認
- [ ] バックエンドAPIを起動
- [ ] ヘルスチェックで動作確認
- [ ] Flutterアプリを起動
- [ ] アプリが正常に起動することを確認

### デモデータ準備

- [ ] テスト顧客を登録（またはアプリ内で登録）
- [ ] テスト生地を登録（またはアプリ内で登録）
- [ ] 発注書生成が動作することを確認

### デモ練習

- [ ] デモフローを3回以上練習
- [ ] 5分以内でデモが完了することを確認
- [ ] エラーが発生しないことを確認
- [ ] 質問への回答を準備

---

## 📊 システム構成図

```
┌─────────────────┐
│  Flutter App    │
│  (スマホ/タブレット) │
└────────┬────────┘
         │ HTTP/JSON
         │ API_BASE_URL
         ↓
┌─────────────────┐
│  Backend API    │
│  (Go)           │
│  Port: 8080     │
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
    ↓         ↓
┌────────┐ ┌──────────┐
│Firestore│ │PostgreSQL│
│(オプション)│ │(オプション)│
└────────┘ └──────────┘
```

---

## 🔧 環境変数の詳細

### バックエンド環境変数

| 変数名 | 説明 | デフォルト | 必須 |
|--------|------|-----------|------|
| `PORT` | サーバーポート | 8080 | いいえ |
| `POSTGRES_HOST` | PostgreSQLホスト | localhost | いいえ |
| `POSTGRES_PORT` | PostgreSQLポート | 5432 | いいえ |
| `POSTGRES_USER` | PostgreSQLユーザー | tailorcloud | いいえ |
| `POSTGRES_PASSWORD` | PostgreSQLパスワード | (空) | いいえ |
| `POSTGRES_DB` | PostgreSQLデータベース | tailorcloud | いいえ |
| `GCP_PROJECT_ID` | FirebaseプロジェクトID | (空) | いいえ |
| `GOOGLE_APPLICATION_CREDENTIALS` | Firebase認証情報 | (空) | いいえ |
| `GCS_BUCKET_NAME` | Cloud Storageバケット | (空) | いいえ |

### Flutterアプリ環境変数

| 変数名 | 説明 | デフォルト | 必須 |
|--------|------|-----------|------|
| `API_BASE_URL` | バックエンドAPIのURL | http://localhost:8080 | いいえ |

---

## 🚀 次のステップ

### 1. デモの練習（必須）

- デモフローを最低3回練習
- 5分以内で完了することを確認
- エラーが発生しないことを確認

### 2. LOI獲得活動の開始

- ターゲットリストの作成
- デモの実施
- LOI 10件の獲得

### 3. 資金調達

- J-KISS投資家へのアプローチ
- ピッチデッキの作成
- 3,000〜5,000万円の調達

---

**最終更新日**: 2025-01  
**ステータス**: ✅ 完全起動手順書完成

