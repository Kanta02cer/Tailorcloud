# マイグレーション後の次のステップ

**作成日**: 2025-01  
**状況**: データベースマイグレーション完了、バックエンドサーバー起動

---

## ✅ 完了した作業

- ✅ データベースマイグレーション完了
  - `diagnoses`テーブル作成
  - `appointments`テーブル作成
  - `customers`テーブル作成・拡張

- ✅ バックエンドサーバー起動
  - サーバーURL: http://localhost:8080
  - ヘルスチェック: http://localhost:8080/health

---

## 🎯 次のステップ

### Step 1: テストデータの準備（オプション）

テストデータを準備して、API動作を確認します。

```bash
export PGPASSWORD=tailorcloud_dev_password
/Library/PostgreSQL/17/bin/psql -h localhost -U tailorcloud -d tailorcloud -f scripts/prepare_test_data_suit_mbti.sql
```

---

### Step 2: APIエンドポイントのテスト

#### 診断API（Diagnoses）

**診断を作成**
```bash
curl -X POST http://localhost:8080/api/diagnoses?tenant_id=YOUR_TENANT_ID \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "archetype": "Classic",
    "plan_type": "BestValue",
    "diagnosis_result": {"score": 85}
  }'
```

**診断一覧を取得**
```bash
curl "http://localhost:8080/api/diagnoses?tenant_id=YOUR_TENANT_ID"
```

**診断を取得**
```bash
curl "http://localhost:8080/api/diagnoses/DIAGNOSIS_ID?tenant_id=YOUR_TENANT_ID"
```

---

#### 予約API（Appointments）

**予約を作成**
```bash
curl -X POST http://localhost:8080/api/appointments?tenant_id=YOUR_TENANT_ID \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "appointment_datetime": "2025-12-15T10:00:00Z",
    "duration_minutes": 60,
    "notes": "初回カウンセリング"
  }'
```

**予約一覧を取得**
```bash
curl "http://localhost:8080/api/appointments?tenant_id=YOUR_TENANT_ID"
```

**予約を取得**
```bash
curl "http://localhost:8080/api/appointments/APPOINTMENT_ID?tenant_id=YOUR_TENANT_ID"
```

**空き状況を確認**
```bash
curl "http://localhost:8080/api/appointments/availability?tenant_id=YOUR_TENANT_ID&fitter_id=FITTER_ID&date=2025-12-15"
```

**予約ステータスを更新**
```bash
curl -X PATCH http://localhost:8080/api/appointments/APPOINTMENT_ID?tenant_id=YOUR_TENANT_ID \
  -H "Content-Type: application/json" \
  -d '{
    "status": "CONFIRMED"
  }'
```

---

### Step 3: 詳細なテストガイド

詳細なテスト手順は以下のドキュメントを参照してください:

- `docs/78_Suit_MBTI_Feature_Guide.md` - 機能説明 & API使用方法
- `docs/79_Manual_Testing_Guide.md` - 手動テストガイド

---

## 📝 環境変数の確認

`.env.local`ファイルに以下が設定されていることを確認:

```bash
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=tailorcloud
POSTGRES_PASSWORD=tailorcloud_dev_password
POSTGRES_DB=tailorcloud
```

---

## 🐛 トラブルシューティング

### 問題: PostgreSQL接続エラー

1. `.env.local`ファイルに`POSTGRES_PASSWORD`が設定されているか確認
2. PostgreSQL 17が起動しているか確認
3. データベースとユーザーが存在するか確認

詳細は `docs/83_Troubleshooting_Guide.md` を参照してください。

---

**最終更新日**: 2025-01

