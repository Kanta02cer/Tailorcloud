# TailorCloud Backend API

オーダースーツ業界向けERPシステム「TailorCloud」のバックエンドAPIサーバー

## 🏗️ アーキテクチャ

- **Runtime**: Go 1.21+
- **Framework**: 標準ライブラリ（net/http）
- **Database**: 
  - Firestore (NoSQL) - 注文データ、リアルタイム同期
  - PostgreSQL (RDBMS) - 決済トランザクション（Phase 3で実装予定）
- **Cloud**: Google Cloud Platform
  - Cloud Run (APIサーバー)
  - Cloud Firestore
  - Cloud Storage (契約書PDF保存用)

## 📁 ディレクトリ構成

```
tailor-cloud-backend/
├── cmd/
│   └── api/
│       └── main.go          # エントリーポイント
├── internal/
│   ├── config/
│   │   └── domain/
│   │       ├── models.go        # データモデル定義
│   │       └── compliance.go    # コンプライアンス要件
│   ├── handler/
│   │   └── http_handler.go      # HTTPハンドラー
│   ├── service/
│   │   ├── order_service.go     # 注文ビジネスロジック
│   │   └── compliance_service.go # コンプライアンスエンジン
│   └── repository/
│       └── firestore.go         # Firestoreリポジトリ
├── pkg/                         # 公開パッケージ（今後追加）
├── go.mod
├── go.sum
└── Dockerfile
```

## 🚀 セットアップ

### 1. 環境変数の設定

```bash
export GCP_PROJECT_ID="your-gcp-project-id"
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account-key.json"
export DEFAULT_TENANT_ID="00000000-0000-0000-0000-000000000001"  # デフォルトテナントID
export PORT="8080"  # オプション（デフォルト: 8080）
```

### 1.1. 認証システムのセットアップ

```bash
# デフォルトテナントの作成と環境確認
./scripts/setup_auth.sh
```

### 2. 依存関係のインストール

```bash
go mod download
```

### 3. ローカル実行

```bash
go run cmd/api/main.go
```

### 4. 認証エンドポイントのテスト

```bash
# 認証エンドポイントをテスト（Firebase IDトークンが必要）
./scripts/test_auth.sh <firebase-id-token>
```

詳細は [認証システム ドキュメント](./docs/AUTHENTICATION.md) と [動作確認ガイド](./docs/AUTH_TESTING.md) を参照してください。

### 5. Dockerビルド

```bash
docker build -t tailor-cloud-backend .
docker run -p 8080:8080 \
  -e GCP_PROJECT_ID=your-project-id \
  -e DEFAULT_TENANT_ID=00000000-0000-0000-0000-000000000001 \
  -e GOOGLE_APPLICATION_CREDENTIALS=/path/to/key.json \
  tailor-cloud-backend
```

## 📡 API エンドポイント

### Health Check

```
GET /health
```

**Response**: `200 OK` with body `"OK"`

### 注文作成

```
POST /api/orders
```

**Request Body**:
```json
{
  "tenant_id": "tenant-123",
  "customer_id": "customer-456",
  "fabric_id": "fabric-789",
  "total_amount": 45000,
  "delivery_date": "2025-12-31T00:00:00Z",
  "details": {
    "measurement_data": {},
    "adjustments": {},
    "description": "オーダースーツ縫製（仕様書ID: xxx）"
  },
  "created_by": "user-001"
}
```

**Response**: `201 Created` with Order object

### 注文確定

```
POST /api/orders/confirm
```

**Request Body**:
```json
{
  "order_id": "order-123",
  "tenant_id": "tenant-123",
  "principal_name": "Regalis Societas"
}
```

**Response**: `200 OK` with updated Order object

**Note**: 注文確定時にコンプライアンスエンジンが動作し、契約書PDFが生成されます（Phase 1では構造のみ定義）。

### 注文取得（単一）

```
GET /api/orders?order_id={order_id}&tenant_id={tenant_id}
```

**Response**: `200 OK` with Order object

### 注文一覧取得

```
GET /api/orders?tenant_id={tenant_id}
```

**Response**: `200 OK` with array of Order objects

## 🔐 セキュリティ

### マルチテナントデータ分離

- すべてのクエリで`tenant_id`によるフィルタリングを強制
- リポジトリ層でテナントIDの一致を検証
- データリークを防止する設計

### 監査ログ

- すべての注文変更操作で`updated_at`と`created_by`を記録
- 将来的に監査ログテーブルを実装予定

## 📝 コンプライアンスエンジン

### 下請法・フリーランス保護法への準拠

- **給付の内容**: 注文詳細の`description`フィールドから自動マッピング
- **報酬の額**: 注文の`total_amount`から自動マッピング
- **支払期日**: 納期から60日後を自動計算（下請法60日ルール）

### PDF生成（Phase 1では構造のみ）

- 契約書PDFは`ComplianceService`で生成
- Cloud Storageに保存し、ハッシュ値を計算（改ざん防止）
- Phase 2で実際のPDF生成ライブラリを統合予定

## 🧪 テスト

```bash
# 単体テスト
go test ./...

# カバレッジ
go test -cover ./...
```

## 📚 参考資料

- [TailorCloud システム詳細仕様書](../docs/01_System_Specifications.md)
- [TailorCloud 開発ロードマップ](../docs/02_Development_Roadmap.md)
- [開発着手前チェックリスト](../docs/00_Pre-Development_Checklist.md)

## 🔄 開発フェーズ

### Phase 1: MVP - Compliance First (現在)

- [x] データモデルの実装
- [x] 注文APIの実装
- [x] コンプライアンスエンジンの構造定義
- [ ] PDF生成機能の実装
- [ ] Firebase認証との統合

### Phase 2: Engagement - Inventory & UX

- [ ] 在庫連携API
- [ ] チャット機能（Firestoreリアルタイム同期）

### Phase 3: Monetization - Fintech

- [ ] 決済トランザクション（PostgreSQL）
- [ ] ファクタリング機能

## 📄 ライセンス

Copyright © 2025 Regalis Japan Group. All Rights Reserved.

