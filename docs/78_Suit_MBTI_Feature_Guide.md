# Suit-MBTI統合機能: 機能ガイド & 手動テストマニュアル

**作成日**: 2025-01  
**バージョン**: 1.0  
**ステータス**: Phase 1 Week 1 実装完了

---

## 📋 目次

1. [実装済み機能一覧](#実装済み機能一覧)
2. [ターゲット情報](#ターゲット情報)
3. [API使用方法](#api使用方法)
4. [手動テストガイド](#手動テストガイド)
5. [トラブルシューティング](#トラブルシューティング)

---

## ✅ 実装済み機能一覧

### 1. 診断管理機能（Diagnosis）

**機能概要**: Suit-MBTI診断結果をTailorCloudに統合し、診断履歴を管理

**実装内容**:
- ✅ 診断結果の登録（アーキタイプ、プランタイプ）
- ✅ 診断履歴の取得（ユーザー別、テナント別）
- ✅ 診断結果のフィルター検索

**データモデル**:
- アーキタイプ（RATタイプ）: Classic, Modern, Elegant, Sporty, Casual
- プランタイプ: Best Value, Authentic
- 診断結果詳細（JSON形式）

---

### 2. 予約管理機能（Appointment）

**機能概要**: フィッティング予約の管理と空き状況チェック

**実装内容**:
- ✅ 予約の作成（日時、フィッター、顧客）
- ✅ 予約の取得・一覧表示
- ✅ 予約の更新・キャンセル
- ✅ 空き状況チェック（時間重複防止）
- ✅ 期間フィルター（開始日・終了日）

**データモデル**:
- ステータス: Pending, Confirmed, Cancelled, Completed, NoShow
- デポジット管理（Stripe連携準備済み）
- キャンセル理由の記録

---

## 🎯 ターゲット情報

### 主要ターゲット

#### 1. テーラー（Tailor）

**属性**:
- 個人テーラー、独立系テーラー
- 年商3,000万〜1億円
- 事務作業の効率化を希望

**利用シーン**:
- 顧客の診断結果を確認して適切なプランを提案
- フィッティング予約の管理
- 顧客の診断履歴を参照してリピート提案

**使用機能**:
- 診断履歴の閲覧
- 予約カレンダーの管理
- 顧客カルテ（診断結果含む）

---

#### 2. スタッフ（Staff）

**属性**:
- テーラーの従業員
- 接客業務を担当

**利用シーン**:
- 接客時に診断結果を確認
- フィッティング予約の作成・変更

**使用機能**:
- 診断履歴の閲覧
- 予約の作成・更新
- 空き状況の確認

---

#### 3. 顧客（Customer/User）

**属性**:
- Suit-MBTI診断を受けた顧客
- オーダースーツを購入検討中

**利用シーン**:
- 診断結果の確認
- フィッティング予約の確認・変更

**使用機能**:
- 自分の診断履歴の閲覧
- 自分の予約の確認・キャンセル

---

## 🔌 API使用方法

### ベースURL

```
http://localhost:8080
```

### 認証

すべてのAPIは認証が必要です（開発環境ではOptionalAuthが有効）。

**ヘッダー**:
```
Authorization: Bearer <JWT_TOKEN>
```

開発環境では、認証なしでも`tenant_id`をクエリパラメータで指定すれば動作します。

---

### 1. 診断管理API

#### 1.1 診断結果の登録

**エンドポイント**: `POST /api/diagnoses`

**リクエストボディ**:
```json
{
  "user_id": "user_12345",
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
    "recommendations": ["Classic", "Elegant"]
  }
}
```

**レスポンス**: `201 Created`
```json
{
  "id": "diagnosis_001",
  "user_id": "user_12345",
  "tenant_id": "tenant_001",
  "archetype": "Classic",
  "plan_type": "Best Value",
  "diagnosis_result": {...},
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

**使用例（curl）**:
```bash
curl -X POST http://localhost:8080/api/diagnoses?tenant_id=tenant_001 \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_12345",
    "archetype": "Classic",
    "plan_type": "Best Value",
    "diagnosis_result": {"scores": {"classic": 85}}
  }'
```

---

#### 1.2 診断結果の取得（単一）

**エンドポイント**: `GET /api/diagnoses/{id}`

**クエリパラメータ**:
- `tenant_id` (必須)

**レスポンス**: `200 OK`
```json
{
  "id": "diagnosis_001",
  "user_id": "user_12345",
  "tenant_id": "tenant_001",
  "archetype": "Classic",
  "plan_type": "Best Value",
  "diagnosis_result": {...},
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

**使用例（curl）**:
```bash
curl "http://localhost:8080/api/diagnoses/diagnosis_001?tenant_id=tenant_001"
```

---

#### 1.3 診断履歴の一覧取得

**エンドポイント**: `GET /api/diagnoses`

**クエリパラメータ**:
- `tenant_id` (必須)
- `user_id` (オプション): ユーザーIDでフィルター
- `archetype` (オプション): アーキタイプでフィルター
- `plan_type` (オプション): プランタイプでフィルター
- `limit` (オプション): 件数制限（デフォルト: 20、最大: 100）
- `offset` (オプション): オフセット（デフォルト: 0）

**レスポンス**: `200 OK`
```json
{
  "data": [
    {
      "id": "diagnosis_001",
      "user_id": "user_12345",
      "archetype": "Classic",
      "plan_type": "Best Value",
      ...
    },
    {
      "id": "diagnosis_002",
      "user_id": "user_12346",
      "archetype": "Modern",
      "plan_type": "Authentic",
      ...
    }
  ],
  "total": 2
}
```

**使用例（curl）**:
```bash
# 全診断履歴を取得
curl "http://localhost:8080/api/diagnoses?tenant_id=tenant_001"

# ユーザー別に取得
curl "http://localhost:8080/api/diagnoses?tenant_id=tenant_001&user_id=user_12345"

# アーキタイプでフィルター
curl "http://localhost:8080/api/diagnoses?tenant_id=tenant_001&archetype=Classic"

# ページネーション
curl "http://localhost:8080/api/diagnoses?tenant_id=tenant_001&limit=10&offset=0"
```

---

### 2. 予約管理API

#### 2.1 予約の作成

**エンドポイント**: `POST /api/appointments`

**リクエストボディ**:
```json
{
  "user_id": "user_12345",
  "fitter_id": "fitter_001",
  "appointment_datetime": "2025-01-15T14:00:00Z",
  "duration_minutes": 60,
  "notes": "初回フィッティング"
}
```

**レスポンス**: `201 Created`
```json
{
  "id": "appointment_001",
  "user_id": "user_12345",
  "tenant_id": "tenant_001",
  "fitter_id": "fitter_001",
  "appointment_datetime": "2025-01-15T14:00:00Z",
  "duration_minutes": 60,
  "status": "Pending",
  "notes": "初回フィッティング",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

**バリデーション**:
- 予約日時は未来である必要がある
- 時間重複チェック（同じフィッター・同じ時間帯）
- 予約時間は最大8時間（480分）

**使用例（curl）**:
```bash
curl -X POST http://localhost:8080/api/appointments?tenant_id=tenant_001 \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_12345",
    "fitter_id": "fitter_001",
    "appointment_datetime": "2025-01-15T14:00:00Z",
    "duration_minutes": 60,
    "notes": "初回フィッティング"
  }'
```

---

#### 2.2 予約の取得（単一）

**エンドポイント**: `GET /api/appointments/{id}`

**クエリパラメータ**:
- `tenant_id` (必須)

**レスポンス**: `200 OK`
```json
{
  "id": "appointment_001",
  "user_id": "user_12345",
  "tenant_id": "tenant_001",
  "fitter_id": "fitter_001",
  "appointment_datetime": "2025-01-15T14:00:00Z",
  "duration_minutes": 60,
  "status": "Confirmed",
  ...
}
```

**使用例（curl）**:
```bash
curl "http://localhost:8080/api/appointments/appointment_001?tenant_id=tenant_001"
```

---

#### 2.3 予約一覧の取得

**エンドポイント**: `GET /api/appointments`

**クエリパラメータ**:
- `tenant_id` (必須)
- `user_id` (オプション): ユーザーIDでフィルター
- `fitter_id` (オプション): フィッターIDでフィルター
- `start_date` (オプション): 開始日（ISO 8601形式）
- `end_date` (オプション): 終了日（ISO 8601形式）

**レスポンス**: `200 OK`
```json
{
  "data": [
    {
      "id": "appointment_001",
      "user_id": "user_12345",
      "fitter_id": "fitter_001",
      "appointment_datetime": "2025-01-15T14:00:00Z",
      "status": "Confirmed",
      ...
    }
  ],
  "total": 1
}
```

**使用例（curl）**:
```bash
# 全予約を取得
curl "http://localhost:8080/api/appointments?tenant_id=tenant_001"

# ユーザー別に取得
curl "http://localhost:8080/api/appointments?tenant_id=tenant_001&user_id=user_12345"

# フィッター別に取得
curl "http://localhost:8080/api/appointments?tenant_id=tenant_001&fitter_id=fitter_001"

# 期間指定
curl "http://localhost:8080/api/appointments?tenant_id=tenant_001&start_date=2025-01-01T00:00:00Z&end_date=2025-01-31T23:59:59Z"
```

---

#### 2.4 予約の更新

**エンドポイント**: `PUT /api/appointments/{id}`

**リクエストボディ**:
```json
{
  "fitter_id": "fitter_002",
  "appointment_datetime": "2025-01-15T15:00:00Z",
  "duration_minutes": 90,
  "status": "Confirmed",
  "notes": "時間変更"
}
```

**レスポンス**: `200 OK`
```json
{
  "id": "appointment_001",
  ...
  "appointment_datetime": "2025-01-15T15:00:00Z",
  "status": "Confirmed",
  ...
}
```

**使用例（curl）**:
```bash
curl -X PUT http://localhost:8080/api/appointments/appointment_001?tenant_id=tenant_001 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "Confirmed",
    "notes": "確認済み"
  }'
```

---

#### 2.5 予約のキャンセル

**エンドポイント**: `DELETE /api/appointments/{id}`

**リクエストボディ** (オプション):
```json
{
  "reason": "急用ができたため"
}
```

**レスポンス**: `204 No Content`

**使用例（curl）**:
```bash
curl -X DELETE http://localhost:8080/api/appointments/appointment_001?tenant_id=tenant_001 \
  -H "Content-Type: application/json" \
  -d '{"reason": "急用ができたため"}'
```

---

## 🧪 手動テストガイド

### テスト環境の準備

#### 1. データベースマイグレーション実行

```bash
cd tailor-cloud-backend

# PostgreSQL接続情報を環境変数で設定
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=your_user
export POSTGRES_PASSWORD=your_password
export POSTGRES_DB=tailorcloud

# マイグレーション実行
psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -f migrations/013_create_diagnoses_table.sql
psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -f migrations/014_create_appointments_table.sql
psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -f migrations/015_extend_customers_for_suit_mbti.sql
```

#### 2. バックエンドサーバー起動

```bash
cd tailor-cloud-backend
./scripts/start_backend.sh
# または
go run cmd/api/main.go
```

サーバーは `http://localhost:8080` で起動します。

---

### テストシナリオ

#### シナリオ 1: 診断結果の登録と取得

**目的**: 診断結果を登録し、正常に取得できることを確認

**ステップ**:

1. **診断結果を登録**
   ```bash
   curl -X POST http://localhost:8080/api/diagnoses?tenant_id=tenant_test \
     -H "Content-Type: application/json" \
     -d '{
       "user_id": "user_test_001",
       "archetype": "Classic",
       "plan_type": "Best Value",
       "diagnosis_result": {
         "scores": {"classic": 85, "modern": 20},
         "recommendations": ["Classic"]
       }
     }'
   ```
   
   **期待結果**: `201 Created`、診断IDが返される

2. **診断結果を取得**
   ```bash
   # ステップ1で取得した診断IDを使用
   curl "http://localhost:8080/api/diagnoses/{diagnosis_id}?tenant_id=tenant_test"
   ```
   
   **期待結果**: `200 OK`、登録した診断結果が返される

3. **診断履歴一覧を取得**
   ```bash
   curl "http://localhost:8080/api/diagnoses?tenant_id=tenant_test&user_id=user_test_001"
   ```
   
   **期待結果**: `200 OK`、ステップ1で登録した診断が含まれる

---

#### シナリオ 2: 予約の作成と空き状況チェック

**目的**: 予約を作成し、空き状況チェックが正常に動作することを確認

**ステップ**:

1. **予約を作成**
   ```bash
   curl -X POST http://localhost:8080/api/appointments?tenant_id=tenant_test \
     -H "Content-Type: application/json" \
     -d '{
       "user_id": "user_test_001",
       "fitter_id": "fitter_test_001",
       "appointment_datetime": "2025-01-20T14:00:00Z",
       "duration_minutes": 60,
       "notes": "初回フィッティング"
     }'
   ```
   
   **期待結果**: `201 Created`、予約IDが返される

2. **重複予約を試みる（エラーテスト）**
   ```bash
   # 同じフィッター・同じ時間で予約を試みる
   curl -X POST http://localhost:8080/api/appointments?tenant_id=tenant_test \
     -H "Content-Type: application/json" \
     -d '{
       "user_id": "user_test_002",
       "fitter_id": "fitter_test_001",
       "appointment_datetime": "2025-01-20T14:00:00Z",
       "duration_minutes": 60
     }'
   ```
   
   **期待結果**: `409 Conflict`、エラーメッセージで「not available」と返される

3. **予約一覧を取得（期間指定）**
   ```bash
   curl "http://localhost:8080/api/appointments?tenant_id=tenant_test&fitter_id=fitter_test_001&start_date=2025-01-01T00:00:00Z&end_date=2025-01-31T23:59:59Z"
   ```
   
   **期待結果**: `200 OK`、ステップ1で作成した予約が含まれる

---

#### シナリオ 3: 予約の更新とキャンセル

**目的**: 予約の更新とキャンセルが正常に動作することを確認

**ステップ**:

1. **予約を更新**
   ```bash
   # ステップ1で取得した予約IDを使用
   curl -X PUT http://localhost:8080/api/appointments/{appointment_id}?tenant_id=tenant_test \
     -H "Content-Type: application/json" \
     -d '{
       "status": "Confirmed",
       "notes": "確認済み"
     }'
   ```
   
   **期待結果**: `200 OK`、更新された予約情報が返される

2. **予約をキャンセル**
   ```bash
   curl -X DELETE http://localhost:8080/api/appointments/{appointment_id}?tenant_id=tenant_test \
     -H "Content-Type: application/json" \
     -d '{
       "reason": "テストキャンセル"
     }'
   ```
   
   **期待結果**: `204 No Content`

3. **キャンセルされた予約を取得**
   ```bash
   curl "http://localhost:8080/api/appointments/{appointment_id}?tenant_id=tenant_test"
   ```
   
   **期待結果**: `200 OK`、ステータスが「Cancelled」になっている

---

#### シナリオ 4: エラーハンドリングテスト

**目的**: バリデーションエラーとエラーハンドリングが正常に動作することを確認

**ステップ**:

1. **必須パラメータ欠如（診断登録）**
   ```bash
   curl -X POST http://localhost:8080/api/diagnoses?tenant_id=tenant_test \
     -H "Content-Type: application/json" \
     -d '{
       "archetype": "Classic"
     }'
   ```
   
   **期待結果**: `400 Bad Request`、エラーメッセージで「user_id is required」と返される

2. **無効なアーキタイプ（診断登録）**
   ```bash
   curl -X POST http://localhost:8080/api/diagnoses?tenant_id=tenant_test \
     -H "Content-Type: application/json" \
     -d '{
       "user_id": "user_test_001",
       "archetype": "InvalidType"
     }'
   ```
   
   **期待結果**: `400 Bad Request`、エラーメッセージで「invalid archetype」と返される

3. **過去の日時で予約作成**
   ```bash
   curl -X POST http://localhost:8080/api/appointments?tenant_id=tenant_test \
     -H "Content-Type: application/json" \
     -d '{
       "user_id": "user_test_001",
       "fitter_id": "fitter_test_001",
       "appointment_datetime": "2020-01-01T14:00:00Z",
       "duration_minutes": 60
     }'
   ```
   
   **期待結果**: `400 Bad Request`、エラーメッセージで「must be in the future」と返される

4. **存在しない予約を取得**
   ```bash
   curl "http://localhost:8080/api/appointments/nonexistent_id?tenant_id=tenant_test"
   ```
   
   **期待結果**: `404 Not Found`、エラーメッセージで「not found」と返される

---

## 🔍 トラブルシューティング

### よくある問題

#### 1. データベース接続エラー

**症状**: 
```
ERROR: Failed to connect to PostgreSQL
```

**解決方法**:
- PostgreSQLが起動しているか確認
- 環境変数（POSTGRES_HOST, POSTGRES_USER等）が正しく設定されているか確認
- マイグレーションが実行されているか確認

```bash
# PostgreSQL接続確認
psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -c "SELECT 1"
```

---

#### 2. テーブルが存在しないエラー

**症状**:
```
ERROR: relation "diagnoses" does not exist
```

**解決方法**:
- マイグレーションを実行

```bash
cd tailor-cloud-backend
psql -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB -f migrations/013_create_diagnoses_table.sql
```

---

#### 3. 認証エラー

**症状**:
```
ERROR: Authentication required
```

**解決方法**:
- 開発環境では`tenant_id`をクエリパラメータで指定
- 本番環境ではJWTトークンを`Authorization`ヘッダーに設定

---

#### 4. 予約時間重複エラー

**症状**:
```
ERROR: appointment time slot is not available
```

**解決方法**:
- 異なる時間帯を指定
- 異なるフィッターを指定
- 既存の予約をキャンセル

---

## 📊 テストチェックリスト

### 診断API

- [ ] 診断結果の登録（正常系）
- [ ] 診断結果の取得（正常系）
- [ ] 診断履歴一覧の取得（正常系）
- [ ] ユーザー別フィルター
- [ ] アーキタイプフィルター
- [ ] プランタイプフィルター
- [ ] ページネーション
- [ ] 必須パラメータ欠如エラー
- [ ] 無効なアーキタイプエラー

### 予約API

- [ ] 予約の作成（正常系）
- [ ] 予約の取得（正常系）
- [ ] 予約一覧の取得（正常系）
- [ ] 予約の更新（正常系）
- [ ] 予約のキャンセル（正常系）
- [ ] ユーザー別フィルター
- [ ] フィッター別フィルター
- [ ] 期間フィルター
- [ ] 時間重複チェック
- [ ] 過去の日時エラー
- [ ] 存在しない予約取得エラー

---

**最終更新日**: 2025-01  
**バージョン**: 1.0

