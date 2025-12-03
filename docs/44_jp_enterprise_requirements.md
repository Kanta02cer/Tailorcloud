# TailorCloud: 日本型エンタープライズ要件定義書

**作成日**: 2025-01  
**対象**: 日本の大企業（商社・アパレルHD）向け  
**バージョン**: 1.0.0

---

## 📋 エグゼクティブサマリー

日本の商社・大手アパレル企業が導入する際に**必ずチェックされる「非機能要件」および「日本固有の業務要件」**を定義します。

この要件定義書に準拠することで、導入時の**法務審査・内部監査・セキュリティ審査**を通過できます。

---

## 1. Compliance & Legal (法務・コンプライアンス)

### 1.1 下請法 (Subcontract Act) 対応 🔴 Critical

#### A. 3条書面（発注書）の自動交付

**法的要件**:
> 下請法第3条: 委託をする者は、委託をした時までに、委託をする者の氏名及び住所、給付の内容、報酬の額、支払期日その他の事項を記載した書面を交付しなければならない。

**実装要件**:

1. **発注確定と同時のPDF自動生成**
   - 発注確定ボタンを押した瞬間に、法的要件を満たしたPDFを生成
   - タイムスタンプ付きで保存（改ざん防止）
   - Cloud Storage（WORM設定）に保存

2. **必須記載項目**:
   - ✅ 委託をする者の氏名及び住所（発注元企業情報）
   - ✅ 給付の内容（スーツ1着の縫製など）
   - ✅ 報酬の額（税抜金額 + 消費税額）
   - ✅ 支払期日（納期から60日以内）
   - ✅ 納期（納品希望日）

3. **修正注文書の履歴管理**
   - 後からデータを修正した場合、「修正注文書」を履歴として残す
   - **上書き禁止**（元のPDFは保持）
   - 修正理由を記録

**技術実装**:

```go
// 下請法PDF生成サービス
type ComplianceService interface {
    GenerateSubcontractDocument(orderID string) (*ComplianceDocument, error)
    GenerateAmendmentDocument(originalOrderID string, changes map[string]interface{}) (*ComplianceDocument, error)
}

// PDF生成（Go）
// 使用ライブラリ: unidoc/unipdf または gofpdf
```

**データモデル**:

```sql
CREATE TABLE compliance_documents (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    document_type VARCHAR(50) NOT NULL, -- 'INITIAL', 'AMENDMENT'
    parent_document_id UUID, -- 修正元の文書ID
    pdf_url TEXT NOT NULL,
    pdf_hash VARCHAR(255) NOT NULL, -- SHA-256
    generated_at TIMESTAMPTZ NOT NULL,
    generated_by UUID NOT NULL,
    amendment_reason TEXT -- 修正理由
);
```

**優先度**: 🔴 Critical  
**工数見積**: 2週間

---

### 1.2 インボイス制度 (Invoice System) 対応 🔴 Critical

#### A. 適格請求書発行

**法的要件**:
> インボイス制度: 2023年10月から、適格請求書等保存方式により消費税額を正確に計算する必要がある。

**実装要件**:

1. **登録番号（T番号）の管理**
   - テナントごとにインボイス登録番号を保存
   - 請求書PDFに「T番号」を記載

2. **税率ごとの消費税額の正確な計算**
   - 標準税率（10%）と軽減税率（8%）の混在に対応
   - 端数処理（切り捨て・切り上げ）のルール統一
   - **総額方式**と**個別計算方式**の選択可能

3. **請求書PDFの生成**
   - 税率ごとの消費税額を明記
   - 税抜金額、消費税額、税込金額を分けて表示

**技術実装**:

```go
// インボイス計算サービス
type InvoiceService interface {
    CalculateTax(orderID string, taxMethod string) (*TaxCalculation, error)
    GenerateInvoicePDF(orderID string) (*InvoiceDocument, error)
}

// 端数処理
type RoundingMethod string
const (
    RoundingDown RoundingMethod = "DOWN"    // 切り捨て
    RoundingUp   RoundingMethod = "UP"      // 切り上げ
    RoundingHalf RoundingMethod = "HALF_UP" // 四捨五入
)
```

**データモデル**:

```sql
-- tenantsテーブルに追加
ALTER TABLE tenants ADD COLUMN invoice_registration_no VARCHAR(50);
ALTER TABLE tenants ADD COLUMN tax_rounding_method VARCHAR(20) DEFAULT 'HALF_UP';

-- ordersテーブルに追加
ALTER TABLE orders ADD COLUMN tax_amount BIGINT NOT NULL DEFAULT 0;
ALTER TABLE orders ADD COLUMN tax_rate DECIMAL(3,2) NOT NULL DEFAULT 0.10;
ALTER TABLE orders ADD COLUMN invoice_issued_at TIMESTAMPTZ;
```

**優先度**: 🔴 Critical  
**工数見積**: 2週間

---

## 2. Inventory Management (在庫管理の厳密化)

### 2.1 理論在庫 vs 実在庫 🔴 Critical

#### A. 反物（Roll）管理

**課題**:
- 単なる「50メートル」という数値ではなく、「A反(30m) + B反(20m)」という物理的な内訳を管理する必要がある
- 「2.5mのスーツを作りたいが、A反の残りが2.0mしかない（B反を使わないといけない）」という**「キレ（端尺）問題」**を防ぐ

**実装要件**:

1. **物理的な巻き（Roll）単位の管理**
   - 各反物にロット番号を付与
   - 入荷時長さと現在長さを個別管理
   - 在庫の場所（倉庫ID or 店舗ID）を管理

2. **在庫引当ロジック**
   - 発注時に、必要な用尺を複数の反物から自動で組み合わせて確保
   - 1つの反物で足りない場合は、複数の反物から組み合わせ
   - 端尺（キレ）が発生しないように最適化

3. **在庫ステータス管理**
   - `AVAILABLE`: 利用可能
   - `RESERVED`: 発注で確保済み（未裁断）
   - `CUT`: 裁断済み
   - `EMPTY`: 使い切った

**技術実装**:

```go
// 在庫引当サービス
type AllocationService interface {
    AllocateFabric(orderID string, fabricID string, requiredLength float64) (*AllocationResult, error)
    ReleaseAllocation(allocationID string) error
}

// 引当結果
type AllocationResult struct {
    Allocations []FabricAllocation // どの反物から何m確保したか
    TotalAllocated float64         // 合計確保長さ
    RemainingRolls []FabricRoll    // 残りの反物リスト
}
```

**データモデル**:

```sql
-- 反物テーブル（新規追加）
CREATE TABLE fabric_rolls (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    fabric_id UUID NOT NULL,
    location_id UUID, -- 倉庫ID or 店舗ID
    lot_number VARCHAR(100), -- ロット番号
    initial_length DECIMAL(10, 2) NOT NULL, -- 入荷時長さ (m)
    current_length DECIMAL(10, 2) NOT NULL, -- 現在長さ (m)
    status VARCHAR(50) DEFAULT 'AVAILABLE', -- AVAILABLE, RESERVED, CUT, EMPTY
    version INTEGER DEFAULT 1, -- 楽観的ロック用
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (fabric_id) REFERENCES fabrics(id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- 在庫引当テーブル
CREATE TABLE fabric_allocations (
    id UUID PRIMARY KEY,
    order_item_id UUID NOT NULL,
    fabric_roll_id UUID NOT NULL,
    allocated_length DECIMAL(5, 2) NOT NULL,
    allocated_at TIMESTAMPTZ DEFAULT NOW(),
    status VARCHAR(50), -- RESERVED, CUT, CANCELLED
    FOREIGN KEY (order_item_id) REFERENCES order_items(id),
    FOREIGN KEY (fabric_roll_id) REFERENCES fabric_rolls(id)
);
```

**優先度**: 🔴 Critical  
**工数見積**: 3週間

---

#### B. 預かり在庫・委託在庫

**課題**:
- 商社から店舗への「委託在庫（消化仕入）」に対応
- 在庫の「所有権」と「保管場所」を分けて管理する必要がある

**実装要件**:

1. **所有権の管理**
   - `OWNER`: 自社所有
   - `CONSIGNED`: 委託在庫（他社所有、自社で保管）
   - `ON_CONSIGNMENT`: 委託出庫（自社所有、他社で保管）

2. **会計処理の対応**
   - 委託在庫は売上計上時にのみ在庫減少を記録
   - 消化仕入の仕訳処理に対応

**データモデル**:

```sql
-- fabric_rollsテーブルに追加
ALTER TABLE fabric_rolls ADD COLUMN ownership_type VARCHAR(50) DEFAULT 'OWNER';
ALTER TABLE fabric_rolls ADD COLUMN owner_tenant_id UUID; -- 所有権を持つテナント
ALTER TABLE fabric_rolls ADD COLUMN custodian_tenant_id UUID; -- 保管場所のテナント
```

**優先度**: 🟡 High  
**工数見積**: 1週間

---

## 3. Security & Governance (セキュリティ)

### 3.1 IPアドレス制限 & 端末認証 🔴 Critical

**要件**:
- 店舗のiPad以外（個人のスマホ等）からはアクセスできないようにする
- 「クライアント証明書」または「IP制限」機能

**実装要件**:

1. **IP制限（ネットワークレベル）**
   - Cloud Load Balancer（GCP）またはAPI GatewayでIP制限を設定
   - ホワイトリスト方式: 許可されたIPアドレスのみアクセス可能
   - VPN経由でのアクセスも対応

2. **端末認証（アプリケーションレベル）**
   - デバイスID（UUID）をFirebase Authenticationに紐付け
   - 初回ログイン時にデバイスを登録
   - 管理者がデバイス登録を承認

**技術実装**:

```go
// デバイス認証ミドルウェア
func DeviceAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        deviceID := r.Header.Get("X-Device-ID")
        if !isDeviceRegistered(deviceID) {
            http.Error(w, "Device not registered", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

**データモデル**:

```sql
-- デバイステーブル
CREATE TABLE devices (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    device_id VARCHAR(255) UNIQUE NOT NULL,
    device_name VARCHAR(255),
    ip_address VARCHAR(45),
    is_active BOOLEAN DEFAULT TRUE,
    registered_at TIMESTAMPTZ DEFAULT NOW(),
    last_access_at TIMESTAMPTZ,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

**優先度**: 🔴 Critical  
**工数見積**: 1週間

---

### 3.2 操作ログ (Audit Log) ⚠️ 部分実装

**要件**:
- 「誰が価格マスタを変更したか」「誰が在庫を強制修正したか」を**5年以上保存**
- 改ざん不能なストレージ（WORM）に保管

**実装要件**:

1. **監査ログの強化**
   - ✅ 基本的な監査ログは実装済み（`audit_logs`テーブル）
   - ❌ 5年以上の保存期間設定なし
   - ❌ WORMストレージへの移行なし

2. **保存期間の設定**
   - 監査ログは5年以上保存
   - 古いログはCloud Storage（WORM設定）にアーカイブ

3. **改ざん防止**
   - Cloud StorageのWORM（Write Once Read Many）機能を使用
   - ログのハッシュ値を計算して保存

**技術実装**:

```go
// 監査ログアーカイブサービス
type AuditLogArchiveService interface {
    ArchiveOldLogs(retentionPeriod time.Duration) error
    ExportToWORMStorage(logs []AuditLog) error
}
```

**優先度**: 🟡 High  
**工数見積**: 1週間

---

## 4. データ整合性 & トランザクション管理

### 4.1 排他制御（同時発注対策）🔴 Critical

**課題**:
- A店とB店が同時に「残り1着の生地」を発注した際、ダブルブッキングが起きる

**実装要件**:

1. **楽観的ロック（Optimistic Locking）**
   - `fabric_rolls`テーブルに`version`カラムを追加
   - 更新時にversionをチェック
   - 更新件数が0なら競合エラー

2. **悲観的ロック（Pessimistic Locking）**
   - トランザクション内で行をロック
   - `SELECT ... FOR UPDATE`を使用

3. **DBレベルのトランザクション分離レベル**
   - PostgreSQLのデフォルト（READ COMMITTED）を維持
   - 必要に応じてSERIALIZABLEを使用

**詳細**: `docs/43_gap_analysis.md` セクション2.3参照

**優先度**: 🔴 Critical  
**工数見積**: 2週間

---

## 5. 開発優先度 (Priority)

### Phase 1: Inventory Core (反物管理・引当ロジック) 🔴 Critical

**期間**: 3週間

- 反物（Roll）管理テーブルの実装
- 在庫引当ロジックの実装
- キレ（端尺）問題の解決

**検証**: 「複数の反物から最適に組み合わせて確保できるか？」

---

### Phase 2: Order Transaction (発注・排他制御) 🔴 Critical

**期間**: 2週間

- 排他制御の実装
- 同時発注時の競合処理
- トランザクション管理の強化

**検証**: 「複数の端末から同時にアクセスしてもデータが壊れないか？」

---

### Phase 3: PDF Generation (下請法対応) 🔴 Critical

**期間**: 2週間

- 下請法3条書面のPDF生成
- 修正注文書の履歴管理
- WORMストレージへの保存

**検証**: 「法的要件を満たしたPDFが生成されるか？」

---

### Phase 4: Audit Log (ログ基盤) 🟡 High

**期間**: 1週間

- 監査ログの5年保存
- WORMストレージへのアーカイブ
- ログ検索機能

**検証**: 「5年前のログが検索できるか？」

---

## 6. 要件実装チェックリスト

### コンプライアンス

- [ ] 下請法3条書面の自動生成
- [ ] 修正注文書の履歴管理
- [ ] インボイス登録番号の管理
- [ ] 税率ごとの消費税額計算
- [ ] 端数処理のルール統一

### 在庫管理

- [ ] 反物（Roll）単位の管理
- [ ] 在庫引当ロジック
- [ ] キレ（端尺）問題の解決
- [ ] 預かり在庫・委託在庫対応

### セキュリティ

- [ ] IP制限機能
- [ ] 端末認証
- [ ] 監査ログの5年保存
- [ ] WORMストレージへの移行

### データ整合性

- [ ] 楽観的ロック
- [ ] 悲観的ロック
- [ ] トランザクション分離レベルの設定

---

## 7. 結論

この要件定義書に準拠することで、日本の商社・大手アパレル企業が要求する**法務・コンプライアンス・セキュリティ要件**をすべて満たすことができます。

**次のアクション**:
1. `schema/enterprise_schema.sql` を参照してDBスキーマを確認
2. Phase 1の開発を開始（反物管理・引当ロジック）

---

**最終更新日**: 2025-01  
**ステータス**: ✅ 要件定義完了

