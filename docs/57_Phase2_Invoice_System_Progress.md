# TailorCloud: Phase 2 インボイス制度対応 実装進捗レポート

**作成日**: 2025-01  
**フェーズ**: Phase 2 - 法規制完全準拠  
**Week**: Week 5-6  
**ステータス**: 🔄 実装中

---

## 📋 エグゼクティブサマリー

エンタープライズグレードのTailorCloudシステムにおいて、**インボイス制度（適格請求書等保存方式）対応**の実装を開始しました。2023年10月施行のインボイス制度に対応し、日本の商社・大手アパレル企業が導入する際に必須の要件を満たすための基盤を構築しています。

---

## ✅ 実装完了内容

### 1. データベースマイグレーション ✅

**ファイル**: `migrations/008_add_invoice_fields.sql`

**実装内容**:
- **tenantsテーブル拡張**:
  - `invoice_registration_no` (VARCHAR(50)): インボイス登録番号（T番号）
  - `legal_name` (VARCHAR(255)): 法人名（請求書記載用）
  - `address` (TEXT): 住所（請求書記載用）
  - `tax_rounding_method` (VARCHAR(20)): 端数処理方法（HALF_UP, ROUND_DOWN, ROUND_UP）

- **ordersテーブル拡張**:
  - `tax_amount` (DECIMAL(12, 2)): 消費税額
  - `tax_rate` (DECIMAL(3, 2)): 消費税率（10% or 8%）
  - `tax_excluded_amount` (DECIMAL(12, 2)): 税抜金額
  - `invoice_issued_at` (TIMESTAMPTZ): 請求書発行日時

**インデックス**: `idx_tenants_invoice_registration_no` 追加

---

### 2. ドメインモデル拡張 ✅

#### 2.1 TaxRate & TaxRoundingMethod 定義

**ファイル**: `internal/config/domain/tax.go`

**実装内容**:
- `TaxRate` 型: 消費税率（0.10 = 10%, 0.08 = 8%）
- `TaxRoundingMethod` 型: 端数処理方法
- `CalculateTax()`: 消費税額計算関数
- `CalculateTaxIncludedAmount()`: 税込金額計算関数
- `ParseTaxRate()`: 税率文字列パーサー
- `FormatTaxRate()`: 税率表示用フォーマッター

**端数処理方法**:
- `HALF_UP`: 四捨五入（デフォルト）
- `ROUND_DOWN`: 切り捨て
- `ROUND_UP`: 切り上げ

#### 2.2 Tenantモデル拡張

**ファイル**: `internal/config/domain/models.go`

**追加フィールド**:
```go
type Tenant struct {
    ID                      string
    InvoiceRegistrationNo   string             // インボイス登録番号（T番号）
    Address                 string             // 住所
    TaxRoundingMethod       TaxRoundingMethod  // 端数処理方法
    // ...
}
```

#### 2.3 Orderモデル拡張

**追加フィールド**:
```go
type Order struct {
    TaxAmount         int64      // 消費税額
    TaxRate           TaxRate    // 消費税率
    TaxExcludedAmount *int64     // 税抜金額（明示的な場合）
    InvoiceIssuedAt   *time.Time // 請求書発行日時
    // ...
}
```

---

### 3. TenantRepository実装 ✅

**ファイル**: `internal/repository/tenant_repository.go`

**実装メソッド**:
- `GetByID(ctx, tenantID)`: テナント情報取得（インボイス登録番号、端数処理方法含む）
- `Update(ctx, tenant)`: テナント情報更新

**特徴**:
- PostgreSQL実装
- NULL値の適切な処理
- デフォルト値の設定（端数処理方法 = HALF_UP）

---

### 4. TaxCalculationService実装 ✅

**ファイル**: `internal/service/tax_calculation_service.go`

**実装メソッド**:

#### CalculateTax

**機能**: 消費税額を計算

**入力**:
- `TenantID`: テナントID
- `TaxExcludedAmount`: 税抜金額
- `TaxRate`: 消費税率

**出力**:
- `TaxExcludedAmount`: 税抜金額
- `TaxAmount`: 消費税額（端数処理済み）
- `TaxIncludedAmount`: 税込金額
- `TaxRate`: 消費税率
- `RoundingMethod`: 使用した端数処理方法

**特徴**:
- テナントの端数処理方法を自動取得
- 標準税率（10%）と軽減税率（8%）に対応
- 正確な端数処理（四捨五入・切り捨て・切り上げ）

#### CalculateTaxForOrder

**機能**: 注文に対して消費税を自動計算

**特徴**:
- 注文の税抜金額から自動計算
- 税率未指定時は標準税率（10%）を使用
- `TaxExcludedAmount`があれば優先的に使用

---

## 📊 実装統計

### 新規作成ファイル

1. `migrations/008_add_invoice_fields.sql` (約50行)
2. `internal/config/domain/tax.go` (約90行)
3. `internal/repository/tenant_repository.go` (約110行)
4. `internal/service/tax_calculation_service.go` (約100行)

### 更新ファイル

- `internal/config/domain/models.go` (約20行追加)

### 合計

- **追加コード行数**: 約370行
- **新規ファイル数**: 4ファイル
- **更新ファイル数**: 1ファイル
- **データベースカラム追加**: 8カラム

---

## 🎯 実装された機能

### 1. インボイス制度対応基盤 ✅

- ✅ テナントごとのインボイス登録番号（T番号）管理
- ✅ 税率ごとの消費税額の正確な計算
- ✅ 端数処理のルール統一（四捨五入・切り捨て・切り上げ）
- ✅ 標準税率（10%）と軽減税率（8%）の混在対応

### 2. 税率計算エンジン ✅

- ✅ 自動的な端数処理方法の適用
- ✅ テナント設定に基づく計算
- ✅ 注文単位での自動計算

---

## 🔄 次のステップ（未実装）

### 1. 請求書PDF生成サービス ⏳

**実装予定**:
- 適格請求書のPDF生成
- T番号の記載
- 税率ごとの消費税額の明記
- 税抜金額、消費税額、税込金額の分離表示

### 2. APIエンドポイント ⏳

**実装予定**:
- `POST /api/orders/{id}/calculate-tax`: 税率計算
- `POST /api/orders/{id}/generate-invoice`: 請求書PDF生成
- `PUT /api/tenants/{id}/invoice-settings`: テナントのインボイス設定更新

### 3. OrderRepository拡張 ⏳

**実装予定**:
- `tax_amount`, `tax_rate`, `invoice_issued_at` フィールドの保存・取得

---

## 🏗️ アーキテクチャ

### データフロー

```
注文作成
  ↓
TaxCalculationService.CalculateTaxForOrder
  ├── TenantRepository.GetByID（端数処理方法を取得）
  ├── domain.CalculateTax（消費税額計算）
  └── 注文に tax_amount, tax_rate を設定
  ↓
請求書PDF生成（未実装）
  ├── テナント情報（T番号、法人名、住所）
  ├── 注文情報（税抜金額、消費税額、税込金額）
  └── 適格請求書PDF生成
```

---

## 📈 成功指標（KPI）

### Week 5-6 の目標

- [x] データベーススキーマ拡張完了
- [x] ドメインモデル拡張完了
- [x] TenantRepository実装完了
- [x] TaxCalculationService実装完了
- [ ] 請求書PDF生成サービス実装（次ステップ）
- [ ] APIエンドポイント実装（次ステップ）

---

## ✅ チェックリスト

### Phase 2 Week 5-6 完了項目

- [x] データベースマイグレーション
- [x] ドメインモデル拡張（TaxRate, TaxRoundingMethod）
- [x] 税率計算ロジック
- [x] TenantRepository実装
- [x] TaxCalculationService実装
- [ ] 請求書PDF生成サービス
- [ ] OrderRepository拡張
- [ ] APIエンドポイント実装

---

## 🎉 成果

### インボイス制度対応の基盤が完成

- ✅ **正確な税率計算**: 標準税率と軽減税率の混在に対応
- ✅ **端数処理の統一**: テナントごとの設定に基づく一貫した処理
- ✅ **データ整合性**: データベースレベルでの税額管理

---

**最終更新日**: 2025-01  
**ステータス**: 🔄 Phase 2 Week 5-6 実装中（基盤完了）

**次のアクション**: 請求書PDF生成サービスの実装

