# Suit-MBTI統合: Phase 1 実装タスク詳細

**作成日**: 2025-01  
**フェーズ**: Phase 1 - 管理画面 & CRM構築  
**期間**: 3〜4週間

---

## 📋 実装タスク一覧

### Week 1: データベース設計 & バックエンドAPI基盤

#### タスク 1.1: データベースマイグレーション作成 ✅

- [x] `013_create_diagnoses_table.sql` 作成
- [x] `014_create_appointments_table.sql` 作成
- [x] `015_extend_customers_for_suit_mbti.sql` 作成
- [ ] マイグレーション実行テスト
- [ ] ロールバックテスト

**ファイル**:
- `tailor-cloud-backend/migrations/013_create_diagnoses_table.sql`
- `tailor-cloud-backend/migrations/014_create_appointments_table.sql`
- `tailor-cloud-backend/migrations/015_extend_customers_for_suit_mbti.sql`

---

#### タスク 1.2: ドメインモデル定義

**診断（Diagnosis）モデル**

- [ ] `internal/config/domain/diagnosis.go` 作成
- [ ] `Diagnosis` 構造体定義
- [ ] `CreateDiagnosisRequest` 定義
- [ ] JSONマーシャリング対応

**予約（Appointment）モデル**

- [ ] `internal/config/domain/appointment.go` 作成
- [ ] `Appointment` 構造体定義
- [ ] `CreateAppointmentRequest` 定義
- [ ] `UpdateAppointmentRequest` 定義
- [ ] JSONマーシャリング対応

**ファイル**:
- `tailor-cloud-backend/internal/config/domain/diagnosis.go` (NEW)
- `tailor-cloud-backend/internal/config/domain/appointment.go` (NEW)

---

#### タスク 1.3: リポジトリ層実装

**診断リポジトリ**

- [ ] `internal/repository/diagnosis_repository.go` 作成
- [ ] `DiagnosisRepository` インターフェース定義
- [ ] `PostgreSQLDiagnosisRepository` 実装
  - [ ] `Create(ctx, diagnosis) error`
  - [ ] `GetByID(ctx, id) (*Diagnosis, error)`
  - [ ] `GetByUserID(ctx, userID) ([]*Diagnosis, error)`
  - [ ] `GetByTenantID(ctx, tenantID) ([]*Diagnosis, error)`
  - [ ] `List(ctx, filter) ([]*Diagnosis, error)`

**予約リポジトリ**

- [ ] `internal/repository/appointment_repository.go` 作成
- [ ] `AppointmentRepository` インターフェース定義
- [ ] `PostgreSQLAppointmentRepository` 実装
  - [ ] `Create(ctx, appointment) error`
  - [ ] `GetByID(ctx, id) (*Appointment, error)`
  - [ ] `GetByUserID(ctx, userID) ([]*Appointment, error)`
  - [ ] `GetByFitterID(ctx, fitterID, startDate, endDate) ([]*Appointment, error)`
  - [ ] `GetByTenantID(ctx, tenantID, startDate, endDate) ([]*Appointment, error)`
  - [ ] `Update(ctx, appointment) error`
  - [ ] `Cancel(ctx, id, reason) error`

**ファイル**:
- `tailor-cloud-backend/internal/repository/diagnosis_repository.go` (NEW)
- `tailor-cloud-backend/internal/repository/appointment_repository.go` (NEW)

---

#### タスク 1.4: サービス層実装

**診断サービス**

- [ ] `internal/service/diagnosis_service.go` 作成
- [ ] `DiagnosisService` 構造体定義
- [ ] `CreateDiagnosis(ctx, request) (*Diagnosis, error)`
- [ ] `GetDiagnosis(ctx, id) (*Diagnosis, error)`
- [ ] `GetDiagnosesByUser(ctx, userID) ([]*Diagnosis, error)`
- [ ] `GetDiagnosesByTenant(ctx, tenantID) ([]*Diagnosis, error)`

**予約サービス**

- [ ] `internal/service/appointment_service.go` 作成
- [ ] `AppointmentService` 構造体定義
- [ ] `CreateAppointment(ctx, request) (*Appointment, error)`
- [ ] `GetAppointment(ctx, id) (*Appointment, error)`
- [ ] `ListAppointments(ctx, filter) ([]*Appointment, error)`
- [ ] `UpdateAppointment(ctx, id, request) (*Appointment, error)`
- [ ] `CancelAppointment(ctx, id, reason) error`
- [ ] `CheckAvailability(ctx, fitterID, datetime, duration) (bool, error)`

**ファイル**:
- `tailor-cloud-backend/internal/service/diagnosis_service.go` (NEW)
- `tailor-cloud-backend/internal/service/appointment_service.go` (NEW)

---

### Week 2: HTTPハンドラー & API実装

#### タスク 2.1: 診断APIハンドラー

- [ ] `internal/handler/diagnosis_handler.go` 作成
- [ ] `DiagnosisHandler` 構造体定義
- [ ] `POST /api/diagnoses` - 診断作成
- [ ] `GET /api/diagnoses/{id}` - 診断取得
- [ ] `GET /api/diagnoses?user_id={id}` - ユーザーの診断一覧
- [ ] `GET /api/diagnoses?tenant_id={id}` - テナントの診断一覧

**ファイル**:
- `tailor-cloud-backend/internal/handler/diagnosis_handler.go` (NEW)

---

#### タスク 2.2: 予約APIハンドラー

- [ ] `internal/handler/appointment_handler.go` 作成
- [ ] `AppointmentHandler` 構造体定義
- [ ] `POST /api/appointments` - 予約作成
- [ ] `GET /api/appointments/{id}` - 予約取得
- [ ] `GET /api/appointments?tenant_id={id}&start_date={date}&end_date={date}` - 予約一覧取得
- [ ] `PUT /api/appointments/{id}` - 予約更新
- [ ] `DELETE /api/appointments/{id}` - 予約キャンセル

**ファイル**:
- `tailor-cloud-backend/internal/handler/appointment_handler.go` (NEW)

---

#### タスク 2.3: ルーティング統合

- [ ] `cmd/api/main.go` に診断APIルーティング追加
- [ ] `cmd/api/main.go` に予約APIルーティング追加
- [ ] 認証・認可ミドルウェア適用
- [ ] エラーハンドリング確認

**ファイル**:
- `tailor-cloud-backend/cmd/api/main.go` (UPDATE)

---

### Week 3: 顧客プロフィールAPI拡張 & 分析API

#### タスク 3.1: 顧客プロフィールAPI拡張

- [ ] `internal/handler/customer_handler.go` 拡張
- [ ] `GET /api/customers/{id}/profile` - 顧客プロフィール（診断結果含む）
- [ ] `GET /api/customers/{id}/diagnoses` - 診断履歴
- [ ] `POST /api/customers/{id}/notes` - メモ追加

**ファイル**:
- `tailor-cloud-backend/internal/handler/customer_handler.go` (UPDATE)

---

#### タスク 3.2: 分析API実装

- [ ] `internal/handler/analytics_handler.go` 作成
- [ ] `GET /api/analytics/sales?tenant_id={id}&period={month}` - 売上分析
- [ ] `GET /api/analytics/appointments?tenant_id={id}&period={month}` - 予約分析
- [ ] `GET /api/analytics/plan-distribution?tenant_id={id}` - プラン別構成比
- [ ] `GET /api/analytics/cpa?tenant_id={id}&period={month}` - CPA分析

**ファイル**:
- `tailor-cloud-backend/internal/handler/analytics_handler.go` (NEW)
- `tailor-cloud-backend/internal/service/analytics_service.go` (NEW)

---

### Week 4: フロントエンド統合準備 & テスト

#### タスク 4.1: API仕様書作成

- [ ] OpenAPI仕様書作成
- [ ] Postmanコレクション作成
- [ ] APIドキュメント更新

---

#### タスク 4.2: 統合テスト

- [ ] 診断API統合テスト
- [ ] 予約API統合テスト
- [ ] エラーハンドリングテスト
- [ ] パフォーマンステスト

---

#### タスク 4.3: Suit-MBTI Reactアプリ準備

- [ ] TailorCloudバックエンドとの接続確認
- [ ] APIクライアント実装準備
- [ ] 認証統合準備

---

## 📝 実装順序

1. **データベースマイグレーション** (Day 1-2)
2. **ドメインモデル定義** (Day 2-3)
3. **リポジトリ層実装** (Day 3-5)
4. **サービス層実装** (Day 5-7)
5. **HTTPハンドラー実装** (Day 8-10)
6. **ルーティング統合** (Day 10-11)
7. **顧客プロフィールAPI拡張** (Day 12-13)
8. **分析API実装** (Day 14-15)
9. **テスト & ドキュメント** (Day 16-20)

---

**最終更新日**: 2025-01

