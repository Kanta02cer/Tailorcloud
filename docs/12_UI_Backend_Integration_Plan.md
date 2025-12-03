# TailorCloud UI・バックエンド統合実装計画

**作成日**: 2025-01  
**デザインバージョン**: 1.0.0

---

## 📋 概要

Figmaデザイン仕様に基づき、フロントエンド（Flutter）とバックエンド（Go API）の統合に必要な実装項目を整理しました。

---

## 🎯 実装優先順位

### Phase 1: UI対応のためのバックエンド拡張（最優先）

#### 1.1 Fabricモデルの拡張 ✅

**ステータス**: 実装完了

- [x] `image_url` フィールド追加（生地画像表示用）
- [x] `minimum_order` フィールド追加（最小発注数量 = 3.2m）

#### 1.2 Inventory API実装

**ステータス**: 実装必要

- [ ] `GET /api/fabrics` - 生地一覧取得（フィルター対応）
  - クエリパラメータ: `tenant_id`, `status`, `search`
- [ ] `GET /api/fabrics/{id}` - 生地詳細取得
- [ ] `POST /api/fabrics/{id}/reserve` - 在庫確保（発注フロー開始）

#### 1.3 Dashboard API実装

**ステータス**: 実装必要

- [ ] `GET /api/dashboard?tenant_id={id}` - ダッシュボードデータ取得
  - KPIデータ（月間売上、注文件数等）
  - タスクリスト（承認待ち等）

---

## 📊 データモデル拡張

### Fabricモデル拡張 ✅

```go
type Fabric struct {
    ID           string      `json:"id"`
    SupplierID   string      `json:"supplier_id"`
    Name         string      `json:"name"`
    StockAmount  float64     `json:"stock_amount"`
    Price        int64       `json:"price"`
    StockStatus  StockStatus `json:"stock_status"`
    ImageURL     string      `json:"image_url"`          // ✅ 追加
    MinimumOrder float64     `json:"minimum_order"`      // ✅ 追加（デフォルト3.2m）
    CreatedAt    time.Time   `json:"created_at"`
    UpdatedAt    time.Time   `json:"updated_at"`
}
```

### Taskモデル新規実装（必要）

```go
type Task struct {
    ID          string     `json:"id"`
    TenantID    string     `json:"tenant_id"`
    Type        TaskType   `json:"type"`
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Status      TaskStatus `json:"status"`
    Priority    int        `json:"priority"`
    CreatedAt   time.Time  `json:"created_at"`
}

type TaskType string

const (
    TaskTypeComplianceApproval TaskType = "compliance_approval"  // 下請法対応書類の承認待ち
    TaskTypeInventoryCheck     TaskType = "inventory_check"      // 生地在庫の確認が必要
    TaskTypeOrderApproval      TaskType = "order_approval"       // 新規注文の承認待ち
    TaskTypeFactoryReply       TaskType = "factory_reply"        // 工場からの返信待ち
)
```

---

## 🔌 APIエンドポイント仕様

### Inventory API

#### GET /api/fabrics

**説明**: 生地一覧を取得（フィルター・検索対応）

**クエリパラメータ**:
- `tenant_id` (必須): テナントID
- `status` (オプション): フィルター (`all`, `available`, `limited`, `soldout`)
- `search` (オプション): 検索キーワード（生地名で検索）

**レスポンス例**:
```json
{
  "fabrics": [
    {
      "id": "fabric-1",
      "name": "Premium Navy Wool",
      "supplier_id": "supplier-1",
      "price": 4500,
      "stock_amount": 5.2,
      "stock_status": "Available",
      "image_url": "https://storage.googleapis.com/.../fabric-1.jpg",
      "minimum_order": 3.2,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    },
    {
      "id": "fabric-2",
      "name": "Classic Black",
      "price": 4200,
      "stock_amount": 2.5,
      "stock_status": "Limited",
      "image_url": "https://storage.googleapis.com/.../fabric-2.jpg",
      "minimum_order": 3.2,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  ],
  "total": 2
}
```

#### POST /api/fabrics/{fabric_id}/reserve

**説明**: 在庫確保（発注フロー開始）

**リクエスト**:
```json
{
  "tenant_id": "tenant-123",
  "amount": 3.2
}
```

**レスポンス**:
```json
{
  "reservation_id": "reservation-1",
  "fabric_id": "fabric-1",
  "amount": 3.2,
  "status": "reserved",
  "expires_at": "2025-01-01T01:00:00Z"
}
```

### Dashboard API

#### GET /api/dashboard?tenant_id={tenant_id}

**説明**: ダッシュボードデータ取得

**レスポンス例**:
```json
{
  "kpis": {
    "monthly_revenue": {
      "value": 2450000,
      "trend": "+15",
      "trend_direction": "up"
    },
    "monthly_orders": {
      "value": 42,
      "trend": "+8",
      "trend_direction": "up"
    },
    "pending_tasks_count": {
      "value": 3,
      "trend": null
    }
  },
  "pending_tasks": [
    {
      "id": "task-1",
      "type": "compliance_approval",
      "title": "下請法対応書類の承認待ち",
      "icon": "warning",
      "count": 3,
      "created_at": "2025-01-01T00:00:00Z"
    },
    {
      "id": "task-2",
      "type": "inventory_check",
      "title": "生地在庫の確認が必要",
      "icon": "warning",
      "count": 1,
      "created_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

---

## 🏗️ 実装タスク

### Phase 1.1: Fabricモデル拡張 ✅

- [x] `image_url` フィールド追加
- [x] `minimum_order` フィールド追加
- [ ] PostgreSQLマイグレーション（fabricsテーブル拡張）

### Phase 1.2: Inventory API実装

- [ ] FabricRepository実装（PostgreSQL版）
- [ ] FabricService実装
  - 一覧取得（フィルター・検索対応）
  - 在庫ステータス計算
  - 在庫確保機能
- [ ] FabricHandler実装
  - `GET /api/fabrics`
  - `GET /api/fabrics/{id}`
  - `POST /api/fabrics/{id}/reserve`

### Phase 1.3: Dashboard API実装

- [ ] Taskモデル実装
- [ ] DashboardService実装
  - KPI集計（月間売上、注文件数等）
  - タスクリスト取得
- [ ] DashboardHandler実装
  - `GET /api/dashboard`

---

## 📝 PostgreSQLマイグレーション

### fabricsテーブル拡張

```sql
ALTER TABLE fabrics 
ADD COLUMN IF NOT EXISTS image_url TEXT,
ADD COLUMN IF NOT EXISTS minimum_order DECIMAL(10,2) DEFAULT 3.2;

COMMENT ON COLUMN fabrics.image_url IS '生地画像URL（UI表示用）';
COMMENT ON COLUMN fabrics.minimum_order IS '最小発注数量（メートル、デフォルト3.2m = スーツ1着分）';
```

### tasksテーブル新規作成

```sql
CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL,
    priority INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tasks_tenant_id ON tasks(tenant_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
```

---

## 🎨 UI・バックエンド連携フロー

### Inventory画面フロー

```
1. ユーザーがInventory画面を開く
   ↓
2. GET /api/fabrics?tenant_id={id}&status=all
   ↓
3. レスポンスを受信
   ├─ 在庫ステータスに基づいてバッジ表示
   │  - Available: 緑色 "◎"
   │  - Limited: 黄色 "⚠ 残り2.5m"
   │  - SoldOut: 赤色 "SOLD OUT"
   └─ 画像URLから生地画像を表示
   ↓
4. ユーザーが"在庫残りわずか"アイテムをタップ
   ↓
5. ボトムシート表示（在庫詳細）
   ↓
6. ユーザーが"この生地を確保して発注"をタップ
   ↓
7. POST /api/fabrics/{id}/reserve
   ↓
8. 発注フロー開始
```

### Dashboard画面フロー

```
1. ユーザーがHome画面を開く
   ↓
2. GET /api/dashboard?tenant_id={id}
   ↓
3. レスポンスを受信
   ├─ KPIカード表示
   │  - 月間売上: ¥2,450,000 (↑ +15%)
   │  - 月間注文件数: 42件 (↑ +8%)
   └─ タスクリスト表示
      - 承認・確認が必要なタスク (Badge: 3)
   ↓
4. ユーザーがタスクをタップ
   ↓
5. 詳細画面へ遷移
```

---

## 🔄 次の実装ステップ

### 即座に実装（Phase 1.1完了後）

1. **FabricRepository実装** - PostgreSQL版
2. **Inventory API実装** - 一覧取得、フィルター、在庫確保
3. **Dashboard API実装** - KPI集計、タスクリスト

### 並行作業

- Flutter側のUI実装準備
- デザインシステム（カラーパレット、タイポグラフィ）の実装

---

**最終更新日**: 2025-01  
**実装者**: AI Assistant (Auto)

