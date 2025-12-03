# TailorCloud: API リファレンス

**作成日**: 2025-01  
**バージョン**: 2.0.0  
**ベースURL**: `http://localhost:8080` (開発環境)

---

## 📋 目次

1. [認証](#認証)
2. [エンドポイント一覧](#エンドポイント一覧)
3. [共通レスポンス](#共通レスポンス)
4. [エラーハンドリング](#エラーハンドリング)

---

## 🔐 認証

### 認証方法

すべてのAPI（`/health`と`/api/metrics`を除く）は、Firebase AuthenticationのJWTトークンを使用します。

### 認証ヘッダー

```
Authorization: Bearer <JWT_TOKEN>
```

### 認証フロー

1. Firebase Authenticationでログイン
2. IDトークンを取得
3. すべてのAPIリクエストに `Authorization` ヘッダーを付与

### 開発環境

開発環境では `OptionalAuth` ミドルウェアが使用され、認証が失敗してもリクエストが通ります（本番環境では無効化）。

---

## 📡 エンドポイント一覧

### ヘルスチェック

#### GET /health

サーバーの状態を確認します。

**認証**: 不要  
**レスポンス**: `200 OK`

```json
"OK"
```

---

### 注文管理

#### POST /api/orders

注文を作成します。

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
  "tax_excluded_amount": 122727,
  "tax_amount": 12273,
  "tax_rate": 0.10,
  "payment_due_date": "2025-03-02T00:00:00Z",
  "delivery_date": "2025-12-31T00:00:00Z",
  "details": {
    "description": "オーダースーツ縫製",
    "measurement_data": {},
    "adjustments": {}
  },
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z",
  "created_by": "user-001"
}
```

---

#### POST /api/orders/confirm

注文を確定します（法的拘束力が発生）。

**認証**: 必須（Owner/Staff）  
**リクエストボディ**:
```json
{
  "order_id": "order-001",
  "principal_name": "Regalis Societas"
}
```

**レスポンス**: `200 OK`
```json
{
  "id": "order-001",
  "status": "Confirmed",
  ...
}
```

---

#### GET /api/orders

注文を取得します（単一または一覧）。

**認証**: 必須  
**クエリパラメータ**:
- `order_id` (オプション): 単一注文取得
- `tenant_id` (必須): テナントID
- `page` (オプション): ページ番号（デフォルト: 1）
- `limit` (オプション): 1ページあたりの件数（デフォルト: 20）

**レスポンス**: `200 OK`

単一取得:
```json
{
  "id": "order-001",
  ...
}
```

一覧取得:
```json
{
  "data": [
    {
      "id": "order-001",
      ...
    },
    {
      "id": "order-002",
      ...
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

---

### コンプライアンス文書

#### POST /api/orders/{id}/generate-document

下請法対応の発注書PDFを生成します。

**認証**: 必須（Owner/Staff）  
**パスパラメータ**:
- `id`: 注文ID

**レスポンス**: `200 OK`
```json
{
  "order_id": "order-001",
  "doc_url": "https://storage.googleapis.com/bucket/path/to/document.pdf",
  "doc_hash": "sha256:abc123...",
  "generated_at": "2025-01-01T00:00:00Z"
}
```

**特徴**:
- 日本語フォント対応（Noto Sans JP）
- フリーランス保護法・下請法完全準拠
- Cloud Storageに保存
- SHA-256ハッシュによる改ざん検出

---

#### POST /api/orders/{id}/generate-amendment

修正発注書PDFを生成します。

**認証**: 必須（Owner/Staff）  
**パスパラメータ**:
- `id`: 注文ID

**リクエストボディ**:
```json
{
  "amendment_reason": "納期変更のため"
}
```

**レスポンス**: `200 OK`
```json
{
  "order_id": "order-001",
  "doc_url": "https://storage.googleapis.com/...",
  "doc_hash": "sha256:...",
  "generated_at": "2025-01-01T00:00:00Z",
  "parent_document_id": "doc-001",
  "version": 2
}
```

**特徴**:
- 親文書へのリンク
- バージョン管理
- 修正理由の記録

---

### 顧客管理（CRM）

#### POST /api/customers

顧客を作成します。

**認証**: 必須（Owner/Staff）  
**リクエストボディ**:
```json
{
  "name": "田中 太郎",
  "name_kana": "タナカ タロウ",
  "email": "tanaka@example.com",
  "phone": "090-1234-5678",
  "address": "東京都渋谷区..."
}
```

**レスポンス**: `201 Created`
```json
{
  "id": "customer-001",
  "tenant_id": "tenant-123",
  "name": "田中 太郎",
  "name_kana": "タナカ タロウ",
  "email": "tanaka@example.com",
  "phone": "090-1234-5678",
  "address": "東京都渋谷区...",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

---

#### GET /api/customers/{id}

顧客を取得します。

**認証**: 必須

---

#### GET /api/customers

顧客一覧を取得します。

**認証**: 必須  
**クエリパラメータ**:
- `tenant_id` (必須): テナントID

---

#### PUT /api/customers/{id}

顧客を更新します。

**認証**: 必須（Owner/Staff）

---

#### DELETE /api/customers/{id}

顧客を削除します。

**認証**: 必須（Ownerのみ）

---

#### GET /api/customers/{id}/orders

顧客の注文一覧を取得します。

**認証**: 必須

---

### 生地管理

#### GET /api/fabrics

生地一覧を取得します。

**認証**: 必須  
**クエリパラメータ**:
- `tenant_id` (必須): テナントID
- `status` (オプション): Available/Limited/SoldOut
- `search` (オプション): 検索キーワード

**レスポンス**: `200 OK`
```json
[
  {
    "id": "fabric-001",
    "tenant_id": "tenant-123",
    "brand": "V.B.C",
    "name": "Perennial Navy",
    "sku": "VBC-001-NV",
    "color": "Navy",
    "pattern": "Solid",
    "price_per_meter": 12000,
    "stock_quantity": 150.0,
    "status": "Available",
    "image_url": "https://...",
    "created_at": "2025-01-01T00:00:00Z"
  }
]
```

---

#### GET /api/fabrics/detail

生地詳細を取得します。

**認証**: 必須  
**クエリパラメータ**:
- `id` (必須): 生地ID

---

#### POST /api/fabrics/reserve

生地を予約します。

**認証**: 必須  
**リクエストボディ**:
```json
{
  "fabric_id": "fabric-001",
  "order_id": "order-001",
  "amount": 3.2
}
```

---

### 反物管理（Roll Management）

#### POST /api/fabric-rolls

反物を作成します。

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

**レスポンス**: `201 Created`
```json
{
  "id": "roll-001",
  "tenant_id": "tenant-123",
  "fabric_id": "fabric-001",
  "roll_number": "VBC-001-NV-001",
  "initial_length": 50.0,
  "current_length": 50.0,
  "status": "Available",
  "location": "倉庫A",
  "created_at": "2025-01-01T00:00:00Z"
}
```

---

#### GET /api/fabric-rolls/{id}

反物を取得します。

**認証**: 必須

---

#### GET /api/fabric-rolls

反物一覧を取得します。

**認証**: 必須  
**クエリパラメータ**:
- `tenant_id` (必須): テナントID
- `fabric_id` (オプション): 生地ID

---

#### PUT /api/fabric-rolls/{id}

反物を更新します。

**認証**: 必須（Owner/Staff）

---

### 在庫引当

#### POST /api/inventory/allocate

在庫を引当します。

**認証**: 必須（Owner/Staff）  
**リクエストボディ**:
```json
{
  "order_id": "order-001",
  "fabric_roll_id": "roll-001",
  "length": 3.2
}
```

**レスポンス**: `200 OK`
```json
{
  "allocation_id": "allocation-001",
  "order_id": "order-001",
  "fabric_roll_id": "roll-001",
  "allocated_length": 3.2,
  "remaining_length": 46.8,
  "status": "Allocated"
}
```

**特徴**:
- 楽観的ロック（SELECT FOR UPDATE SKIP LOCKED）
- トランザクション管理
- 同時実行時の安全性

---

#### POST /api/inventory/release

在庫引当を解放します。

**認証**: 必須（Owner/Staff）  
**リクエストボディ**:
```json
{
  "allocation_id": "allocation-001"
}
```

---

### インボイス

#### POST /api/orders/{id}/generate-invoice

インボイスPDFを生成します。

**認証**: 必須（Owner/Staff）  
**パスパラメータ**:
- `id`: 注文ID

**レスポンス**: `200 OK`
```json
{
  "order_id": "order-001",
  "doc_url": "https://storage.googleapis.com/...",
  "doc_hash": "sha256:...",
  "generated_at": "2025-01-01T00:00:00Z"
}
```

**特徴**:
- 適格インボイス対応
- T番号（インボイス登録番号）の表示
- 消費税の正確な計算（10%・8%）
- 端数処理（half-up/down/up）

---

### アンバサダー管理

#### POST /api/ambassadors

アンバサダーを作成します。

**認証**: 必須（Ownerのみ）  
**リクエストボディ**:
```json
{
  "user_id": "user-001",
  "name": "アンバサダーA",
  "email": "ambassador@example.com",
  "phone": "090-1234-5678",
  "commission_rate": 5.0
}
```

---

#### GET /api/ambassadors/me

自分のアンバサダー情報を取得します。

**認証**: 必須

---

#### GET /api/ambassadors

アンバサダー一覧を取得します。

**認証**: 必須  
**クエリパラメータ**:
- `tenant_id` (必須): テナントID

---

#### GET /api/ambassadors/commissions

成果報酬一覧を取得します。

**認証**: 必須  
**クエリパラメータ**:
- `ambassador_id` (必須): アンバサダーID

---

### 権限管理（RBAC）

#### POST /api/permissions

権限を作成します。

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

---

#### GET /api/permissions

権限一覧を取得します。

**認証**: 必須  
**クエリパラメータ**:
- `tenant_id` (必須): テナントID

---

#### POST /api/permissions/check

権限をチェックします。

**認証**: 必須  
**リクエストボディ**:
```json
{
  "resource_type": "Order",
  "resource_id": "order-001",
  "action": "update"
}
```

**レスポンス**: `200 OK`
```json
{
  "granted": true,
  "reason": "User has Owner role"
}
```

---

### 監視・運用

#### GET /api/metrics

メトリクスを取得します。

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

## 📋 共通レスポンス

### 成功レスポンス

- `200 OK`: 成功
- `201 Created`: リソース作成成功
- `204 No Content`: 成功（コンテンツなし）

### エラーレスポンス

- `400 Bad Request`: リクエストエラー
- `401 Unauthorized`: 認証エラー
- `403 Forbidden`: 権限エラー
- `404 Not Found`: リソースが見つからない
- `409 Conflict`: 競合エラー
- `500 Internal Server Error`: サーバーエラー

**エラーレスポンス形式**:
```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {}
}
```

---

## 🔄 エラーハンドリング

### エラーコード一覧

| コード | 説明 |
|--------|------|
| `INVALID_REQUEST` | リクエストが無効 |
| `UNAUTHORIZED` | 認証が必要 |
| `FORBIDDEN` | 権限が不足 |
| `NOT_FOUND` | リソースが見つからない |
| `CONFLICT` | 競合（例: 在庫不足） |
| `INTERNAL_ERROR` | サーバー内部エラー |

### エラーハンドリング例

```json
{
  "error": "Insufficient fabric inventory",
  "code": "CONFLICT",
  "details": {
    "fabric_id": "fabric-001",
    "required": 3.2,
    "available": 2.5
  }
}
```

---

**最終更新日**: 2025-01  
**バージョン**: 2.0.0

