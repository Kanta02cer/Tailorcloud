# TailorCloud × Suit-MBTI統合: クイックスタートガイド

**作成日**: 2025-01  
**目的**: システムをすぐに起動してテストするための最短手順

---

## 🚀 クイックスタート（5分で起動）

### Step 1: 環境変数の設定

```bash
# .env.localファイルを読み込む
source .env.local

# または、環境変数を直接設定
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=your_user
export POSTGRES_PASSWORD=your_password
export POSTGRES_DB=tailorcloud
```

---

### Step 2: PostgreSQL接続確認

```bash
# PostgreSQLが起動しているか確認
pg_isready

# データベースに接続できるか確認
psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -c "SELECT 1"
```

**接続できない場合**:
- PostgreSQLが起動していない場合は起動してください
- データベースが存在しない場合は作成してください:
  ```bash
  createdb -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER tailorcloud
  ```

---

### Step 3: データベースマイグレーション実行

```bash
./scripts/run_migrations_suit_mbti.sh
```

**期待結果**: 3つのマイグレーションが成功

---

### Step 4: テストデータ準備（オプション）

```bash
psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -f scripts/prepare_test_data_suit_mbti.sql
```

---

### Step 5: バックエンドサーバー起動

```bash
./scripts/start_backend.sh
```

サーバーは `http://localhost:8080` で起動します。

---

### Step 6: 動作確認

```bash
# ヘルスチェック
curl http://localhost:8080/health

# 診断APIテスト
curl -X POST "http://localhost:8080/api/diagnoses?tenant_id=tenant_test_suit_mbti" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_test_001",
    "archetype": "Classic",
    "plan_type": "Best Value",
    "diagnosis_result": {"scores": {"classic": 85}}
  }'
```

---

## 📚 詳細なドキュメント

- **[機能ガイド & テストガイド](./78_Suit_MBTI_Feature_Guide.md)** - 機能説明とAPI使用方法
- **[手動テストガイド](./79_Manual_Testing_Guide.md)** - 詳細なテスト手順
- **[現在のシステム状況](./81_Current_System_Status.md)** - システム全体の状況

---

**最終更新日**: 2025-01

