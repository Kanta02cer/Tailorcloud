# バックエンドAPI 動作確認・改善計画

**作成日**: 2025-01  
**ステータス**: 🔄 確認中

---

## 📋 現状確認

### 1. 注文作成API (`POST /api/orders`)

#### 現状
- ✅ 基本的な注文作成機能は実装済み
- ✅ `fabric_id`で1つの生地を指定
- ✅ `OrderDetails`に`measurement_data`を含めることができる

#### 技術仕様書との差異
- ❌ **order_items配列対応が未実装**
  - 技術仕様書ではレスポンスに`order_items`配列が含まれる想定
  - スキーマには`order_items`テーブルが存在するが、未使用
  - 現在は`fabric_id`が1つだけ

#### 改善が必要な点
1. **order_itemsテーブルの活用**
   - 注文作成時に`order_items`レコードを作成
   - レスポンスに`order_items`配列を含める
   - MVPでは1つのアイテムのみでも、将来的な拡張に対応

2. **レスポンス形式の統一**
   - 技術仕様書に合わせて`order_items`配列を含む形式に変更

---

### 2. PDFダウンロードエンドポイント

#### 現状
- ✅ `POST /api/orders/{id}/generate-document`でPDF生成・アップロード
- ✅ Cloud StorageのURLを返す
- ❌ **直接ダウンロードエンドポイントが未実装**

#### 改善が必要な点
1. **PDFダウンロードエンドポイントの追加**
   - `GET /api/orders/{id}/compliance-document/download`
   - Cloud StorageからPDFを取得して直接ダウンロード
   - 認証チェックを含む

2. **認証・認可の実装**
   - テナントIDの一致確認
   - 適切な権限チェック

---

## 🎯 改善計画

### Phase 1: order_items配列対応（優先度: 🟡 中）

#### 実装内容
1. **OrderItemモデルの追加**
   ```go
   type OrderItem struct {
       ID                  string          `json:"id"`
       OrderID             string          `json:"order_id"`
       ItemType            string          `json:"item_type"` // "SUIT", "SHIRT", etc.
       FabricID            string          `json:"fabric_id"`
       Measurements        json.RawMessage `json:"measurements"`
       Options             json.RawMessage `json:"options"`
       RequiredFabricLength float64        `json:"required_fabric_length"`
       UnitPrice           int64           `json:"unit_price"`
       Quantity            int             `json:"quantity"`
   }
   ```

2. **OrderRepositoryの拡張**
   - `CreateOrderItem`メソッドの追加
   - `GetOrderItemsByOrderID`メソッドの追加

3. **OrderServiceの修正**
   - 注文作成時に`order_items`レコードを作成
   - レスポンスに`order_items`配列を含める

4. **HTTPハンドラーの修正**
   - レスポンス形式を技術仕様書に合わせる

#### 見積もり
- **時間**: 2-3時間
- **難易度**: 🟡 中

---

### Phase 2: PDFダウンロードエンドポイント（優先度: 🟡 中）

#### 実装内容
1. **ComplianceHandlerにダウンロードエンドポイント追加**
   ```go
   // DownloadDocument GET /api/orders/{id}/compliance-document/download
   func (h *ComplianceHandler) DownloadDocument(w http.ResponseWriter, r *http.Request)
   ```

2. **実装ロジック**
   - 注文IDから最新のコンプライアンス文書を取得
   - Cloud StorageからPDFをダウンロード
   - 認証・認可チェック
   - PDFをレスポンスとして返す

3. **エラーハンドリング**
   - 注文が見つからない場合
   - PDFが見つからない場合
   - 認証エラー

#### 見積もり
- **時間**: 1-2時間
- **難易度**: 🟢 低

---

## 📝 実装詳細

### order_items配列対応の実装

#### 1. OrderItemモデルの追加
`internal/config/domain/models.go`に追加

#### 2. OrderItemRepositoryの作成
`internal/repository/order_item_repository.go`を作成

#### 3. OrderServiceの修正
- `CreateOrder`メソッドを修正して`order_items`を作成
- `GetOrder`メソッドを修正して`order_items`を含める

#### 4. HTTPハンドラーの修正
- レスポンス形式を技術仕様書に合わせる

---

### PDFダウンロードエンドポイントの実装

#### 1. ComplianceHandlerにメソッド追加
```go
func (h *ComplianceHandler) DownloadDocument(w http.ResponseWriter, r *http.Request) {
    // 1. 注文IDを取得
    // 2. 認証チェック
    // 3. 最新のコンプライアンス文書を取得
    // 4. Cloud StorageからPDFをダウンロード
    // 5. PDFをレスポンスとして返す
}
```

#### 2. ルーティングの追加
`cmd/api/main.go`にルートを追加

---

## 🚀 実装順序

1. **order_items配列対応**（優先度: 🟡 中）
   - データモデルの拡張
   - リポジトリの実装
   - サービスの修正
   - ハンドラーの修正

2. **PDFダウンロードエンドポイント**（優先度: 🟡 中）
   - ハンドラーの実装
   - ルーティングの追加
   - テスト

---

## ✅ 完了条件

### order_items配列対応
- [ ] OrderItemモデルの追加
- [ ] OrderItemRepositoryの実装
- [ ] OrderServiceの修正
- [ ] HTTPハンドラーの修正
- [ ] レスポンス形式の確認

### PDFダウンロードエンドポイント
- [ ] DownloadDocumentメソッドの実装
- [ ] ルーティングの追加
- [ ] 認証・認可の実装
- [ ] エラーハンドリング
- [ ] テスト

---

## 📊 進捗状況

- **order_items配列対応**: 0% (未着手)
- **PDFダウンロードエンドポイント**: 0% (未着手)

---

**最終更新日**: 2025-01  
**ステータス**: 🔄 確認中

