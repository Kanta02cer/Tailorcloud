# TailorCloud: 完全システム仕様書

**作成日**: 2025-01  
**バージョン**: 2.0.0  
**ステータス**: シード調達用MVP完成、エンタープライズ機能実装完了

---

## 📋 エグゼクティブサマリー

TailorCloudは、オーダースーツ業界向けの**マルチテナント型ERPシステム**です。

**核心価値**:
- **フリーランス保護法対応**: 発注書をスマホで3分で作成
- **エンタープライズグレード**: 100店舗×10工場×年間10万発注に対応
- **法規制完全準拠**: 下請法、インボイス制度、監査ログ対応

**技術スタック**:
- **バックエンド**: Go 1.21+, PostgreSQL, Firestore
- **フロントエンド**: Flutter 3.16+, Dart
- **インフラ**: Google Cloud Platform (Cloud Run, Cloud SQL, Cloud Storage)

---

## 🏗️ システムアーキテクチャ

### 全体構成

```
┌─────────────────────────────────────────────────────────────┐
│                    TailorCloud System                        │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌─────────────────┐              ┌──────────────────┐      │
│  │  Flutter App    │              │   Web Portal     │      │
│  │  (Mobile/Tablet)│              │   (Future)       │      │
│  └────────┬────────┘              └────────┬─────────┘      │
│           │                                 │                │
│           │ HTTPS/REST API                  │                │
│           └────────────┬────────────────────┘                │
│                        │                                     │
│  ┌─────────────────────┴─────────────────────┐              │
│  │      Backend API (Go)                      │              │
│  │      - Authentication (Firebase Auth)      │              │
│  │      - Business Logic                      │              │
│  │      - Compliance Engine                   │              │
│  │      - PDF Generation                      │              │
│  └─────────────┬─────────────────┬───────────┘              │
│                │                 │                          │
│      ┌─────────┴─────────┐  ┌───┴──────────┐              │
│      │   PostgreSQL      │  │  Firestore   │              │
│      │   (Primary DB)    │  │ (Secondary)  │              │
│      │                   │  │              │              │
│      │ - Orders          │  │ - Chat Logs  │              │
│      │ - Customers       │  │ - UI Status  │              │
│      │ - Fabrics         │  │ - Notifications│            │
│      │ - Audit Logs      │  │              │              │
│      └───────────────────┘  └──────────────┘              │
│                │                                           │
│      ┌─────────┴─────────┐                               │
│      │  Cloud Storage    │                               │
│      │  (PDF Documents)  │                               │
│      └───────────────────┘                               │
│                                                           │
└───────────────────────────────────────────────────────────┘
```

### 技術スタック

#### バックエンド

- **言語**: Go 1.21+
- **フレームワーク**: 標準ライブラリ（net/http）
- **データベース**:
  - **PostgreSQL** (Primary): 注文、顧客、在庫、監査ログ
  - **Firestore** (Secondary): チャット、UI状態、通知
- **認証**: Firebase Authentication (JWT)
- **ストレージ**: Google Cloud Storage (PDF保存)
- **監視**: 構造化ログ、メトリクス、アラート

#### フロントエンド

- **言語**: Dart 3.2+
- **フレームワーク**: Flutter 3.16+
- **状態管理**: Riverpod
- **API通信**: HTTP + JSON
- **ローカルストレージ**: Hive (オフライン対応)

---

## 📊 データベース設計

### テーブル一覧

#### 1. tenants (テナント)

マルチテナント対応のテナント情報。

| カラム名 | 型 | 説明 |
|---------|-----|------|
| id | VARCHAR(50) PK | テナントID |
| type | VARCHAR(20) | Tailor / Factory |
| legal_name | VARCHAR(255) | 法人名 |
| address | TEXT | 住所 |
| invoice_registration_no | VARCHAR(50) | インボイス登録番号（T番号） |
| tax_rounding_method | VARCHAR(20) | 端数処理方法（HALF_UP/DOWN/UP） |
| created_at | TIMESTAMPTZ | 作成日時 |
| updated_at | TIMESTAMPTZ | 更新日時 |

#### 2. customers (顧客)

顧客情報（CRM）。

| カラム名 | 型 | 説明 |
|---------|-----|------|
| id | VARCHAR(50) PK | 顧客ID |
| tenant_id | VARCHAR(50) FK | テナントID |
| name | VARCHAR(255) | 顧客名 |
| name_kana | VARCHAR(255) | カナ名 |
| email | VARCHAR(255) | メールアドレス |
| phone | VARCHAR(50) | 電話番号 |
| address | TEXT | 住所 |
| created_at | TIMESTAMPTZ | 作成日時 |
| updated_at | TIMESTAMPTZ | 更新日時 |

**インデックス**:
- `idx_customers_tenant_id` (tenant_id)
- `idx_customers_tenant_name` (tenant_id, name)

#### 3. orders (注文)

注文情報（コンプライアンスエンジンの中心）。

| カラム名 | 型 | 説明 |
|---------|-----|------|
| id | VARCHAR(50) PK | 注文ID |
| tenant_id | VARCHAR(50) FK | テナントID |
| customer_id | VARCHAR(50) FK | 顧客ID |
| fabric_id | VARCHAR(50) FK | 生地ID |
| status | VARCHAR(20) | ステータス（Draft/Confirmed/...） |
| total_amount | BIGINT | 合計金額（税込） |
| tax_excluded_amount | BIGINT | 税抜金額 |
| tax_amount | BIGINT | 消費税額 |
| tax_rate | DECIMAL(3,2) | 税率（0.10/0.08） |
| payment_due_date | TIMESTAMPTZ | 支払期日 |
| delivery_date | TIMESTAMPTZ | 納期 |
| details | JSONB | 注文詳細（採寸データ等） |
| compliance_doc_url | TEXT | 発注書PDF URL |
| compliance_doc_hash | VARCHAR(64) | 発注書PDFハッシュ |
| invoice_issued_at | TIMESTAMPTZ | インボイス発行日時 |
| created_at | TIMESTAMPTZ | 作成日時 |
| updated_at | TIMESTAMPTZ | 更新日時 |
| created_by | VARCHAR(50) | 作成者ユーザーID |

**インデックス**:
- `idx_orders_tenant_id` (tenant_id)
- `idx_orders_customer_id` (customer_id)
- `idx_orders_status` (status)
- `idx_orders_tenant_status` (tenant_id, status)
- `idx_orders_delivery_date` (delivery_date)

#### 4. fabrics (生地)

生地情報。

| カラム名 | 型 | 説明 |
|---------|-----|------|
| id | VARCHAR(50) PK | 生地ID |
| tenant_id | VARCHAR(50) FK | テナントID |
| supplier_id | VARCHAR(50) | サプライヤーID |
| brand | VARCHAR(255) | ブランド名 |
| name | VARCHAR(255) | 生地名 |
| sku | VARCHAR(100) | SKU |
| color | VARCHAR(100) | 色 |
| pattern | VARCHAR(100) | パターン |
| price_per_meter | BIGINT | 単価（円/メートル） |
| stock_quantity | DECIMAL(10,2) | 在庫数量（メートル） |
| status | VARCHAR(20) | ステータス（Available/Limited/SoldOut） |
| image_url | TEXT | 画像URL |
| created_at | TIMESTAMPTZ | 作成日時 |
| updated_at | TIMESTAMPTZ | 更新日時 |

**インデックス**:
- `idx_fabrics_tenant_id` (tenant_id)
- `idx_fabrics_brand` (brand)
- `idx_fabrics_tenant_status` (tenant_id, status)

#### 5. fabric_rolls (反物)

反物（Roll）単位の在庫管理。

| カラム名 | 型 | 説明 |
|---------|-----|------|
| id | VARCHAR(50) PK | 反物ID |
| tenant_id | VARCHAR(50) FK | テナントID |
| fabric_id | VARCHAR(50) FK | 生地ID |
| roll_number | VARCHAR(100) | ロール番号 |
| initial_length | DECIMAL(10,2) | 初期長さ（メートル） |
| current_length | DECIMAL(10,2) | 現在の長さ（メートル） |
| status | VARCHAR(20) | ステータス（Available/Allocated/Consumed） |
| location | VARCHAR(255) | 保管場所 |
| owner_type | VARCHAR(20) | 所有者タイプ（Own/Consigned/Entrusted） |
| created_at | TIMESTAMPTZ | 作成日時 |
| updated_at | TIMESTAMPTZ | 更新日時 |

**インデックス**:
- `idx_fabric_rolls_tenant_id` (tenant_id)
- `idx_fabric_rolls_fabric_id` (fabric_id)
- `idx_fabric_rolls_status` (status)
- `idx_fabric_rolls_tenant_status` (tenant_id, status)

#### 6. fabric_allocations (生地割当)

注文に対する生地の割当。

| カラム名 | 型 | 説明 |
|---------|-----|------|
| id | VARCHAR(50) PK | 割当ID |
| order_id | VARCHAR(50) FK | 注文ID |
| fabric_roll_id | VARCHAR(50) FK | 反物ID |
| allocated_length | DECIMAL(10,2) | 割当長さ（メートル） |
| status | VARCHAR(20) | ステータス（Allocated/Released） |
| created_at | TIMESTAMPTZ | 作成日時 |
| updated_at | TIMESTAMPTZ | 更新日時 |

**インデックス**:
- `idx_fabric_allocations_order_id` (order_id)
- `idx_fabric_allocations_roll_id` (fabric_roll_id)
- `idx_fabric_allocations_status` (status)

#### 7. ambassadors (アンバサダー)

学生アンバサダー情報。

| カラム名 | 型 | 説明 |
|---------|-----|------|
| id | VARCHAR(50) PK | アンバサダーID |
| tenant_id | VARCHAR(50) FK | テナントID |
| user_id | VARCHAR(50) | ユーザーID（Firebase Auth） |
| name | VARCHAR(255) | 名前 |
| email | VARCHAR(255) | メールアドレス |
| phone | VARCHAR(50) | 電話番号 |
| status | VARCHAR(20) | ステータス |
| commission_rate | DECIMAL(5,2) | 成果報酬率（%） |
| total_sales | BIGINT | 累計売上 |
| total_commission | BIGINT | 累計報酬 |
| created_at | TIMESTAMPTZ | 作成日時 |
| updated_at | TIMESTAMPTZ | 更新日時 |

#### 8. commissions (成果報酬)

注文ごとの成果報酬。

| カラム名 | 型 | 説明 |
|---------|-----|------|
| id | VARCHAR(50) PK | 報酬ID |
| ambassador_id | VARCHAR(50) FK | アンバサダーID |
| order_id | VARCHAR(50) FK | 注文ID |
| commission_amount | BIGINT | 報酬額 |
| status | VARCHAR(20) | ステータス |
| created_at | TIMESTAMPTZ | 作成日時 |
| paid_at | TIMESTAMPTZ | 支払日時 |

#### 9. compliance_documents (コンプライアンス文書)

発注書PDFの履歴管理。

| カラム名 | 型 | 説明 |
|---------|-----|------|
| id | VARCHAR(50) PK | 文書ID |
| order_id | VARCHAR(50) FK | 注文ID |
| document_type | VARCHAR(20) | INITIAL / AMENDMENT |
| parent_document_id | VARCHAR(50) FK | 親文書ID（修正の場合） |
| pdf_url | TEXT | PDF URL |
| pdf_hash | VARCHAR(64) | PDFハッシュ |
| generated_at | TIMESTAMPTZ | 生成日時 |
| generated_by | VARCHAR(50) FK | 生成者ユーザーID |
| amendment_reason | TEXT | 修正理由 |
| version | INTEGER | バージョン番号 |

#### 10. permissions (権限)

細かい権限管理（RBAC）。

| カラム名 | 型 | 説明 |
|---------|-----|------|
| id | VARCHAR(50) PK | 権限ID |
| tenant_id | VARCHAR(50) FK | テナントID |
| resource_type | VARCHAR(50) | リソースタイプ（Order/Customer/...） |
| resource_id | VARCHAR(50) | リソースID（*は全件） |
| action | VARCHAR(50) | アクション（create/read/update/delete） |
| role | VARCHAR(50) | ロール（Owner/Staff/...） |
| granted | BOOLEAN | 許可/拒否 |
| created_at | TIMESTAMPTZ | 作成日時 |
| updated_at | TIMESTAMPTZ | 更新日時 |

**インデックス**:
- `idx_permissions_tenant_id` (tenant_id)
- `idx_permissions_resource_type` (resource_type)
- `idx_permissions_role` (role)
- `idx_permissions_action` (action)

#### 11. audit_logs (監査ログ)

全ての操作を記録（改ざん防止）。

| カラム名 | 型 | 説明 |
|---------|-----|------|
| id | VARCHAR(50) PK | ログID |
| tenant_id | VARCHAR(50) FK | テナントID |
| user_id | VARCHAR(50) | ユーザーID |
| action | VARCHAR(50) | アクション |
| resource_type | VARCHAR(50) | リソースタイプ |
| resource_id | VARCHAR(50) | リソースID |
| ip_address | VARCHAR(45) | IPアドレス |
| user_agent | TEXT | ユーザーエージェント |
| device_id | VARCHAR(50) | デバイスID |
| old_values | JSONB | 変更前の値 |
| new_values | JSONB | 変更後の値 |
| change_summary | TEXT | 変更サマリー |
| log_hash | VARCHAR(64) | ログハッシュ（改ざん検出用） |
| created_at | TIMESTAMPTZ | 作成日時 |
| archived_at | TIMESTAMPTZ | アーカイブ日時 |
| archive_location | TEXT | アーカイブ先 |

**インデックス**:
- `idx_audit_logs_tenant_id` (tenant_id)
- `idx_audit_logs_user_id` (user_id)
- `idx_audit_logs_created_at` (created_at)
- `idx_audit_logs_resource` (resource_type, resource_id)

#### 12. audit_log_archives (監査ログアーカイブ)

アーカイブされた監査ログのメタデータ。

| カラム名 | 型 | 説明 |
|---------|-----|------|
| id | VARCHAR(50) PK | アーカイブID |
| tenant_id | VARCHAR(50) FK | テナントID |
| archive_period_start | TIMESTAMPTZ | アーカイブ期間開始 |
| archive_period_end | TIMESTAMPTZ | アーカイブ期間終了 |
| file_url | TEXT | アーカイブファイルURL |
| file_hash | VARCHAR(64) | ファイルハッシュ |
| archived_by | VARCHAR(50) | アーカイブ実行者 |
| archived_at | TIMESTAMPTZ | アーカイブ日時 |

---

## 🔌 API仕様

### 認証

すべてのAPI（`/health`と`/api/metrics`を除く）は認証が必要です。

**認証方法**: Firebase Authentication JWT
**ヘッダー**: `Authorization: Bearer <JWT_TOKEN>`

### エンドポイント一覧

#### 1. ヘルスチェック

```
GET /health
```

**認証**: 不要  
**レスポンス**: `200 OK`
```json
"OK"
```

---

#### 2. 注文管理

##### 2.1 注文作成

```
POST /api/orders
```

**認証**: 必須  
**リクエストボディ**:
```json
{
  "customer_id": "customer-001",
  "fabric_id": "fabric-001",
  "total_amount": 135000,
  "delivery_date": "2025-12-31T00:00:00Z",
  "details": {
    "description": "オーダースーツ縫製",
    "measurement_data": {},
    "adjustments": {}
  }
}
```

**レスポンス**: `201 Created`
```json
{
  "id": "order-001",
  "tenant_id": "tenant-123",
  "customer_id": "customer-001",
  "fabric_id": "fabric-001",
  "status": "Draft",
  "total_amount": 135000,
  "delivery_date": "2025-12-31T00:00:00Z",
  "created_at": "2025-01-01T00:00:00Z"
}
```

##### 2.2 注文確定

```
POST /api/orders/confirm
```

**認証**: 必須（Owner/Staff）  
**リクエストボディ**:
```json
{
  "order_id": "order-001",
  "principal_name": "Regalis Societas"
}
```

##### 2.3 注文取得（単一）

```
GET /api/orders?order_id={order_id}&tenant_id={tenant_id}
```

**認証**: 必須

##### 2.4 注文一覧取得

```
GET /api/orders?tenant_id={tenant_id}&page=1&limit=20
```

**認証**: 必須  
**クエリパラメータ**:
- `tenant_id` (必須)
- `page` (オプション、デフォルト: 1)
- `limit` (オプション、デフォルト: 20)

---

#### 3. コンプライアンス文書

##### 3.1 発注書生成（下請法対応）

```
POST /api/orders/{id}/generate-document
```

**認証**: 必須（Owner/Staff）  
**レスポンス**: `200 OK`
```json
{
  "order_id": "order-001",
  "doc_url": "https://storage.googleapis.com/...",
  "doc_hash": "sha256:...",
  "generated_at": "2025-01-01T00:00:00Z"
}
```

##### 3.2 修正発注書生成

```
POST /api/orders/{id}/generate-amendment
```

**認証**: 必須（Owner/Staff）  
**リクエストボディ**:
```json
{
  "amendment_reason": "納期変更のため"
}
```

---

#### 4. 顧客管理（CRM）

##### 4.1 顧客作成

```
POST /api/customers
```

**認証**: 必須（Owner/Staff）  
**リクエストボディ**:
```json
{
  "name": "田中 太郎",
  "name_kana": "タナカ タロウ",
  "email": "tanaka@example.com",
  "phone": "090-1234-5678",
  "address": "東京都..."
}
```

##### 4.2 顧客取得

```
GET /api/customers/{id}
```

##### 4.3 顧客一覧取得

```
GET /api/customers?tenant_id={tenant_id}
```

##### 4.4 顧客更新

```
PUT /api/customers/{id}
```

**認証**: 必須（Owner/Staff）

##### 4.5 顧客削除

```
DELETE /api/customers/{id}
```

**認証**: 必須（Ownerのみ）

##### 4.6 顧客の注文一覧

```
GET /api/customers/{id}/orders
```

---

#### 5. 生地管理

##### 5.1 生地一覧取得

```
GET /api/fabrics?tenant_id={tenant_id}&status={status}
```

**クエリパラメータ**:
- `tenant_id` (必須)
- `status` (オプション: Available/Limited/SoldOut)
- `search` (オプション: 検索キーワード)

##### 5.2 生地詳細取得

```
GET /api/fabrics/detail?id={id}
```

##### 5.3 生地予約

```
POST /api/fabrics/reserve
```

**リクエストボディ**:
```json
{
  "fabric_id": "fabric-001",
  "order_id": "order-001",
  "amount": 3.2
}
```

---

#### 6. 反物管理（Roll Management）

##### 6.1 反物作成

```
POST /api/fabric-rolls
```

**認証**: 必須（Owner/Staff）  
**リクエストボディ**:
```json
{
  "fabric_id": "fabric-001",
  "roll_number": "VBC-001-NV-001",
  "initial_length": 50.0,
  "location": "倉庫A"
}
```

##### 6.2 反物取得

```
GET /api/fabric-rolls/{id}
```

##### 6.3 反物一覧取得

```
GET /api/fabric-rolls?tenant_id={tenant_id}&fabric_id={fabric_id}
```

##### 6.4 反物更新

```
PUT /api/fabric-rolls/{id}
```

---

#### 7. 在庫引当

##### 7.1 在庫引当

```
POST /api/inventory/allocate
```

**認証**: 必須（Owner/Staff）  
**リクエストボディ**:
```json
{
  "order_id": "order-001",
  "fabric_roll_id": "roll-001",
  "length": 3.2
}
```

##### 7.2 在庫解放

```
POST /api/inventory/release
```

**リクエストボディ**:
```json
{
  "allocation_id": "allocation-001"
}
```

---

#### 8. インボイス

##### 8.1 インボイス生成

```
POST /api/orders/{id}/generate-invoice
```

**認証**: 必須（Owner/Staff）  
**レスポンス**: `200 OK`
```json
{
  "order_id": "order-001",
  "doc_url": "https://storage.googleapis.com/...",
  "doc_hash": "sha256:...",
  "generated_at": "2025-01-01T00:00:00Z"
}
```

---

#### 9. アンバサダー管理

##### 9.1 アンバサダー作成

```
POST /api/ambassadors
```

**認証**: 必須（Ownerのみ）

##### 9.2 自分のアンバサダー情報取得

```
GET /api/ambassadors/me
```

##### 9.3 アンバサダー一覧取得

```
GET /api/ambassadors?tenant_id={tenant_id}
```

##### 9.4 成果報酬一覧取得

```
GET /api/ambassadors/commissions?ambassador_id={id}
```

---

#### 10. 権限管理（RBAC）

##### 10.1 権限作成

```
POST /api/permissions
```

**認証**: 必須（Ownerのみ）  
**リクエストボディ**:
```json
{
  "resource_type": "Order",
  "resource_id": "*",
  "action": "create",
  "role": "Staff",
  "granted": true
}
```

##### 10.2 権限一覧取得

```
GET /api/permissions?tenant_id={tenant_id}
```

##### 10.3 権限チェック

```
POST /api/permissions/check
```

**リクエストボディ**:
```json
{
  "resource_type": "Order",
  "resource_id": "order-001",
  "action": "update"
}
```

---

#### 11. 監視・運用

##### 11.1 メトリクス取得

```
GET /api/metrics
```

**認証**: 不要  
**レスポンス**: `200 OK`
```json
{
  "total_requests": 1000,
  "total_errors": 5,
  "error_rate": 0.005,
  "average_latency": 150000000,
  "request_count": 1000,
  "db_connections": 10,
  "db_connections_in_use": 5,
  "timestamp": "2025-01-01T00:00:00Z"
}
```

---

## 📱 フロントエンド仕様

### 画面構成

#### 1. ホーム画面（Dashboard）

**機能**:
- KPI表示（今日の売上、アクティブ注文数、在庫アラート）
- タスクリスト
- ミル更新フィード
- クイック発注ボタン

**ファイル**: `lib/screens/home/home_screen.dart`

#### 2. 在庫画面（Inventory）

**機能**:
- 生地一覧表示
- 検索・フィルター
- 在庫ステータス表示（Available/Limited/SoldOut）

**ファイル**: `lib/screens/inventory/inventory_screen.dart`

#### 3. クイック発注画面

**機能**:
- ステップ1: 顧客選択（既存顧客選択 or 新規登録）
- ステップ2: 生地選択
- ステップ3: 金額・納期入力
- 発注書生成（フリーランス保護法対応）

**ファイル**: `lib/screens/order/quick_order_screen.dart`

**目標**: 3分で発注書作成完了

#### 4. 視覚的発注画面（Visual Ordering）

**機能**:
- 人体図による採寸入力
- 仕様選択
- 注文確定

**ファイル**: `lib/screens/order/visual_ordering_screen.dart`

**ステータス**: プレースホルダー（将来実装）

---

### モデルクラス

#### 1. Order (注文)

```dart
class Order {
  final String id;
  final String tenantId;
  final String customerId;
  final String fabricId;
  final OrderStatus status;
  final int totalAmount;
  final DateTime paymentDueDate;
  final DateTime deliveryDate;
  final OrderDetails? details;
  final String? complianceDocUrl;
  final String? complianceDocHash;
}
```

#### 2. Customer (顧客)

```dart
class Customer {
  final String id;
  final String tenantId;
  final String name;
  final String? nameKana;
  final String? email;
  final String? phone;
  final String? address;
}
```

#### 3. Fabric (生地)

```dart
class Fabric {
  final String id;
  final String tenantId;
  final String brand;
  final String name;
  final String sku;
  final String color;
  final String pattern;
  final int pricePerMeter;
  final double stockQuantity;
  final StockStatus status;
}
```

#### 4. Ambassador (アンバサダー)

```dart
class Ambassador {
  final String id;
  final String tenantId;
  final String userId;
  final String name;
  final double commissionRate;
  final int totalSales;
  final int totalCommission;
}
```

---

## 🔐 セキュリティ

### 認証・認可

1. **Firebase Authentication**
   - JWT トークンベース認証
   - メール/パスワード認証
   - ユーザーIDの取得

2. **RBAC（Role-Based Access Control）**
   - Owner: 全権限
   - Staff: 発注作成・顧客管理
   - Factory_Manager: 受注管理
   - Worker: 作業完了チェック

3. **マルチテナントデータ分離**
   - すべてのクエリで`tenant_id`によるフィルタリング
   - リポジトリ層でのテナントID検証
   - データリークの防止

### 監査ログ

1. **全操作の記録**
   - 誰が、いつ、何を、どのように変更したか
   - IPアドレス、ユーザーエージェント、デバイスID
   - 変更前後の値（JSONB形式）

2. **改ざん検出**
   - SHA-256ハッシュによる整合性チェック
   - WORMストレージ（Write Once Read Many）

3. **アーカイブ**
   - 5年間の保持期間
   - Cloud Storageへの長期保存
   - メタデータの管理

---

## 📈 監視・運用

### 構造化ログ

- **フォーマット**: JSON
- **ログレベル**: INFO, WARNING, ERROR
- **トレースID**: 全リクエストに付与
- **サービス名**: tailorcloud-backend

### メトリクス収集

- **リクエスト数**: 総リクエスト数、エラー数
- **エラー率**: エラー率の計算
- **レイテンシー**: 平均レイテンシー
- **データベース接続**: 接続数、使用中接続数

### アラート

- **エラー率**: 5%以上でアラート
- **レイテンシー**: 1秒以上でアラート
- **DB接続**: 80%使用でアラート

---

## 🚀 デプロイメント

### バックエンド

**プラットフォーム**: Google Cloud Run

**環境変数**:
- `PORT`: サーバーポート（デフォルト: 8080）
- `POSTGRES_HOST`: PostgreSQLホスト
- `POSTGRES_USER`: PostgreSQLユーザー
- `POSTGRES_PASSWORD`: PostgreSQLパスワード
- `POSTGRES_DB`: PostgreSQLデータベース名
- `GCP_PROJECT_ID`: FirebaseプロジェクトID
- `GCS_BUCKET_NAME`: Cloud Storageバケット名

### フロントエンド

**プラットフォーム**: 
- iOS: App Store
- Android: Google Play Store
- Web: PWA（将来）

**環境変数**:
- `API_BASE_URL`: バックエンドAPIのURL

---

## 📋 実装済み機能

### Phase 1: エンタープライズ基盤 ✅

- ✅ Roll管理システム（反物単位の在庫管理）
- ✅ 在庫引当システム（楽観的ロック）
- ✅ 法規制完全準拠（下請法、インボイス制度）
- ✅ 修正発注書の履歴管理
- ✅ RBAC（細かい権限管理）
- ✅ 監査ログ強化（改ざん検出、アーカイブ）
- ✅ データベース最適化（インデックス、ページネーション）
- ✅ 監視・運用基盤（ログ、メトリクス、アラート）

### Phase 2: シード調達用MVP ✅

- ✅ クイック発注画面（3分で発注書作成）
- ✅ フリーランス保護法対応PDF生成
- ✅ 簡易顧客管理（CRM）
- ✅ スマホ対応UI

---

## 📊 統計

### コード統計

- **バックエンド**: 約15,000行（Go）
- **フロントエンド**: 約5,000行（Dart）
- **データベース**: 12テーブル、50+インデックス
- **APIエンドポイント**: 30+

### 実装ファイル数

- **バックエンド**: 54ファイル
- **フロントエンド**: 30ファイル
- **マイグレーション**: 12ファイル
- **ドキュメント**: 72ファイル

---

## 🔄 今後の開発計画

詳細は `docs/68_Future_Development_Plan.md` を参照。

### Phase 1: PMF達成（資金調達後〜3ヶ月）

1. **デモ改善**（1週間）
2. **認証機能の実装**（2週間）
3. **顧客管理の完全実装**（2週間）
4. **発注書履歴管理**（1週間）
5. **工場連携の基礎**（3週間）

### Phase 2: スケール準備（3〜6ヶ月）

6. **在庫管理の実装**（6週間）
7. **決済機能の実装**（4週間）
8. **レポート機能**（2週間）

### Phase 3: スケール（6〜12ヶ月）

9. **AI機能の実装**（8週間）
10. **マルチブランド対応**（4週間）
11. **API公開**（6週間）

---

**最終更新日**: 2025-01  
**バージョン**: 2.0.0  
**ステータス**: ✅ エンタープライズ機能実装完了、シード調達用MVP完成

