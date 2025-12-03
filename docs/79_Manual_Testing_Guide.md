# 手動テストガイド: Suit-MBTI統合機能

**作成日**: 2025-01  
**目的**: 実装済み機能の手動テスト手順

---

## 📋 テスト前の準備

### 1. データベースマイグレーション実行

```bash
cd /Users/wantan/teiloroud-ERPSystem
./scripts/run_migrations_suit_mbti.sh
```

または手動で実行:

```bash
cd tailor-cloud-backend
psql -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB -f migrations/013_create_diagnoses_table.sql
psql -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB -f migrations/014_create_appointments_table.sql
psql -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB -f migrations/015_extend_customers_for_suit_mbti.sql
```

### 2. バックエンドサーバー起動

```bash
cd tailor-cloud-backend
./scripts/start_backend.sh
```

サーバーは `http://localhost:8080` で起動します。

### 3. テストデータの準備（オプション）

```bash
# テスト用テナントと顧客を作成（既存スクリプトを使用）
psql -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB -f scripts/prepare_test_data.sql
```

---

## 🧪 テストシナリオ

### テストシナリオ 1: 診断結果の登録と取得フロー

#### ステップ 1.1: 診断結果の登録

```bash
curl -X POST "http://localhost:8080/api/diagnoses?tenant_id=tenant_test_001" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_test_001",
    "archetype": "Classic",
    "plan_type": "Best Value",
    "diagnosis_result": {
      "scores": {
        "classic": 85,
        "modern": 20,
        "elegant": 70,
        "sporty": 30,
        "casual": 45
      },
      "recommendations": ["Classic", "Elegant"],
      "notes": "クラシックスタイルを好む傾向が強い"
    }
  }'
```

**期待結果**:
- ステータス: `201 Created`
- レスポンスに診断ID (`id`) が含まれる
- `archetype` が "Classic" になっている
- `plan_type` が "Best Value" になっている

**記録**: 診断IDを保存（例: `diagnosis_001`）

---

#### ステップ 1.2: 診断結果の取得

```bash
# ステップ1.1で取得した診断IDを使用
curl "http://localhost:8080/api/diagnoses/diagnosis_001?tenant_id=tenant_test_001"
```

**期待結果**:
- ステータス: `200 OK`
- 登録した診断結果が正しく返される
- 診断結果詳細（`diagnosis_result`）が含まれる

---

#### ステップ 1.3: 診断履歴一覧の取得（ユーザー別）

```bash
curl "http://localhost:8080/api/diagnoses?tenant_id=tenant_test_001&user_id=user_test_001"
```

**期待結果**:
- ステータス: `200 OK`
- `data` 配列にステップ1.1で登録した診断が含まれる
- 最新の診断が先頭に来る（`created_at` DESC順）

---

#### ステップ 1.4: 診断履歴のフィルター検索

```bash
# アーキタイプでフィルター
curl "http://localhost:8080/api/diagnoses?tenant_id=tenant_test_001&archetype=Classic"

# プランタイプでフィルター
curl "http://localhost:8080/api/diagnoses?tenant_id=tenant_test_001&plan_type=Best%20Value"
```

**期待結果**:
- ステータス: `200 OK`
- 指定したフィルター条件に一致する診断のみが返される

---

### テストシナリオ 2: 予約管理フロー

#### ステップ 2.1: 予約の作成

```bash
curl -X POST "http://localhost:8080/api/appointments?tenant_id=tenant_test_001" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_test_001",
    "fitter_id": "fitter_test_001",
    "appointment_datetime": "2025-02-01T14:00:00Z",
    "duration_minutes": 60,
    "notes": "初回フィッティング、Classicスタイル希望"
  }'
```

**期待結果**:
- ステータス: `201 Created`
- レスポンスに予約ID (`id`) が含まれる
- `status` が "Pending" になっている
- `appointment_datetime` が正しく設定されている

**記録**: 予約IDを保存（例: `appointment_001`）

---

#### ステップ 2.2: 予約の取得

```bash
curl "http://localhost:8080/api/appointments/appointment_001?tenant_id=tenant_test_001"
```

**期待結果**:
- ステータス: `200 OK`
- ステップ2.1で作成した予約情報が正しく返される

---

#### ステップ 2.3: 予約一覧の取得（フィッター別・期間指定）

```bash
# フィッター別に取得
curl "http://localhost:8080/api/appointments?tenant_id=tenant_test_001&fitter_id=fitter_test_001"

# 期間指定で取得
curl "http://localhost:8080/api/appointments?tenant_id=tenant_test_001&start_date=2025-02-01T00:00:00Z&end_date=2025-02-28T23:59:59Z"
```

**期待結果**:
- ステータス: `200 OK`
- `data` 配列にステップ2.1で作成した予約が含まれる
- 期間指定時は、指定期間内の予約のみが返される

---

#### ステップ 2.4: 時間重複チェック（エラーテスト）

```bash
# 同じフィッター・同じ時間帯で予約を試みる
curl -X POST "http://localhost:8080/api/appointments?tenant_id=tenant_test_001" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_test_002",
    "fitter_id": "fitter_test_001",
    "appointment_datetime": "2025-02-01T14:00:00Z",
    "duration_minutes": 60
  }'
```

**期待結果**:
- ステータス: `409 Conflict`
- エラーメッセージに "not available" が含まれる

---

#### ステップ 2.5: 予約の更新

```bash
curl -X PUT "http://localhost:8080/api/appointments/appointment_001?tenant_id=tenant_test_001" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "Confirmed",
    "notes": "確認済み、準備完了"
  }'
```

**期待結果**:
- ステータス: `200 OK`
- `status` が "Confirmed" に更新されている
- `notes` が更新されている
- `updated_at` が更新されている

---

#### ステップ 2.6: 予約のキャンセル

```bash
curl -X DELETE "http://localhost:8080/api/appointments/appointment_001?tenant_id=tenant_test_001" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "顧客都合によりキャンセル"
  }'
```

**期待結果**:
- ステータス: `204 No Content`

**確認**:
```bash
# キャンセルされた予約を取得
curl "http://localhost:8080/api/appointments/appointment_001?tenant_id=tenant_test_001"
```

**期待結果**:
- ステータス: `200 OK`
- `status` が "Cancelled" になっている
- `cancelled_at` が設定されている
- `cancelled_reason` が "顧客都合によりキャンセル" になっている

---

### テストシナリオ 3: エラーハンドリング

#### ステップ 3.1: 必須パラメータ欠如

```bash
# user_idがない診断登録
curl -X POST "http://localhost:8080/api/diagnoses?tenant_id=tenant_test_001" \
  -H "Content-Type: application/json" \
  -d '{
    "archetype": "Classic"
  }'
```

**期待結果**:
- ステータス: `400 Bad Request`
- エラーメッセージに "user_id is required" が含まれる

---

#### ステップ 3.2: 無効なアーキタイプ

```bash
curl -X POST "http://localhost:8080/api/diagnoses?tenant_id=tenant_test_001" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_test_001",
    "archetype": "InvalidType"
  }'
```

**期待結果**:
- ステータス: `400 Bad Request`
- エラーメッセージに "invalid archetype" が含まれる

---

#### ステップ 3.3: 過去の日時で予約作成

```bash
curl -X POST "http://localhost:8080/api/appointments?tenant_id=tenant_test_001" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_test_001",
    "fitter_id": "fitter_test_001",
    "appointment_datetime": "2020-01-01T14:00:00Z",
    "duration_minutes": 60
  }'
```

**期待結果**:
- ステータス: `400 Bad Request`
- エラーメッセージに "must be in the future" が含まれる

---

#### ステップ 3.4: 存在しないリソースの取得

```bash
# 存在しない診断ID
curl "http://localhost:8080/api/diagnoses/nonexistent_id?tenant_id=tenant_test_001"

# 存在しない予約ID
curl "http://localhost:8080/api/appointments/nonexistent_id?tenant_id=tenant_test_001"
```

**期待結果**:
- ステータス: `404 Not Found`
- エラーメッセージに "not found" が含まれる

---

## 📊 テスト結果記録

### テスト結果テンプレート

```markdown
## テスト実施日: YYYY-MM-DD

### シナリオ 1: 診断結果の登録と取得フロー
- [ ] ステップ 1.1: 診断結果の登録 ✅/❌
- [ ] ステップ 1.2: 診断結果の取得 ✅/❌
- [ ] ステップ 1.3: 診断履歴一覧の取得 ✅/❌
- [ ] ステップ 1.4: 診断履歴のフィルター検索 ✅/❌

### シナリオ 2: 予約管理フロー
- [ ] ステップ 2.1: 予約の作成 ✅/❌
- [ ] ステップ 2.2: 予約の取得 ✅/❌
- [ ] ステップ 2.3: 予約一覧の取得 ✅/❌
- [ ] ステップ 2.4: 時間重複チェック ✅/❌
- [ ] ステップ 2.5: 予約の更新 ✅/❌
- [ ] ステップ 2.6: 予約のキャンセル ✅/❌

### シナリオ 3: エラーハンドリング
- [ ] ステップ 3.1: 必須パラメータ欠如 ✅/❌
- [ ] ステップ 3.2: 無効なアーキタイプ ✅/❌
- [ ] ステップ 3.3: 過去の日時で予約作成 ✅/❌
- [ ] ステップ 3.4: 存在しないリソースの取得 ✅/❌

### 発見された問題
1. [問題の説明]
2. [問題の説明]

### 次のアクション
1. [対応が必要な項目]
2. [対応が必要な項目]
```

---

**最終更新日**: 2025-01  
**バージョン**: 1.0

