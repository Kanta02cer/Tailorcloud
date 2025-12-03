# TailorCloud: システム概要

**最終更新日**: 2025-01  
**バージョン**: 2.0.0  
**ステータス**: ✅ シード調達用MVP完成、エンタープライズ機能実装完了

---

## 🎯 システムの核心価値

**「フリーランス保護法対応の発注書をスマホで3分で作れる」**

TailorCloudは、オーダースーツ業界向けの**マルチテナント型ERPシステム**です。

---

## 📊 システム規模

### コード統計

- **バックエンド**: 約15,000行（Go）
- **フロントエンド**: 約5,000行（Dart）
- **合計**: 約20,000行

### データベース

- **テーブル数**: 12テーブル
- **インデックス数**: 50+インデックス
- **マイグレーション**: 12ファイル

### API

- **エンドポイント数**: 30+エンドポイント
- **認証不要**: 2（/health, /api/metrics）
- **認証必須**: 28+

### ファイル構成

- **バックエンド**: 54ファイル
- **フロントエンド**: 30ファイル
- **マイグレーション**: 12ファイル
- **ドキュメント**: 74ファイル
- **スクリプト**: 5ファイル
- **合計**: 175ファイル

---

## 🏗️ アーキテクチャ

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
│  │      Port: 8080                            │              │
│  └─────────────┬─────────────────┬───────────┘              │
│                │                 │                          │
│      ┌─────────┴─────────┐  ┌───┴──────────┐              │
│      │   PostgreSQL      │  │  Firestore   │              │
│      │   (Primary DB)    │  │ (Secondary)  │              │
│      └───────────────────┘  └──────────────┘              │
│                │                                           │
│      ┌─────────┴─────────┐                               │
│      │  Cloud Storage    │                               │
│      │  (PDF Documents)  │                               │
│      └───────────────────┘                               │
│                                                           │
└───────────────────────────────────────────────────────────┘
```

---

## 📋 主要機能

### 1. クイック発注（シード調達用MVP） ✅

- **3分で発注書作成**
- **ステップ1**: 顧客選択（既存 or 新規登録）
- **ステップ2**: 生地選択
- **ステップ3**: 金額・納期入力
- **発注書生成**: フリーランス保護法対応PDF

---

### 2. Roll管理システム（エンタープライズ） ✅

- **反物単位の在庫管理**
- **在庫引当システム**
- **楽観的ロック**（SELECT FOR UPDATE SKIP LOCKED）
- **同時実行時の安全性**

---

### 3. コンプライアンスエンジン ✅

- **下請法対応発注書PDF生成**
- **日本語フォント対応**（Noto Sans JP）
- **修正発注書の履歴管理**
- **PDFハッシュによる改ざん検出**

---

### 4. インボイス制度対応 ✅

- **適格インボイスPDF生成**
- **T番号（インボイス登録番号）対応**
- **消費税計算**（10%・8%）
- **端数処理**（half-up/down/up）

---

### 5. RBAC（権限管理） ✅

- **細かい権限管理**
- **リソース単位の権限設定**
- **動的権限チェック**
- **ロール**: Owner, Staff, Factory_Manager, Worker

---

### 6. 監査ログ ✅

- **全操作の記録**
- **IPアドレス・ユーザーエージェント記録**
- **変更前後の値記録**
- **SHA-256ハッシュによる改ざん検出**
- **アーカイブ機能**（WORMストレージ）
- **5年間の保持期間**

---

### 7. 監視・運用基盤 ✅

- **構造化ログ**（JSON形式）
- **トレースID付与**
- **メトリクス収集**
- **アラート機能**
- **データベース接続監視**

---

## 🔌 APIエンドポイント概要

### 注文管理

- `POST /api/orders` - 注文作成
- `POST /api/orders/confirm` - 注文確定
- `GET /api/orders` - 注文取得（単一・一覧）

### コンプライアンス文書

- `POST /api/orders/{id}/generate-document` - 発注書生成
- `POST /api/orders/{id}/generate-amendment` - 修正発注書生成

### 顧客管理（CRM）

- `POST /api/customers` - 顧客作成
- `GET /api/customers` - 顧客一覧取得
- `PUT /api/customers/{id}` - 顧客更新
- `DELETE /api/customers/{id}` - 顧客削除

### 生地管理

- `GET /api/fabrics` - 生地一覧取得
- `GET /api/fabrics/detail` - 生地詳細取得
- `POST /api/fabrics/reserve` - 生地予約

### 反物管理（Roll Management）

- `POST /api/fabric-rolls` - 反物作成
- `GET /api/fabric-rolls` - 反物一覧取得
- `PUT /api/fabric-rolls/{id}` - 反物更新

### 在庫引当

- `POST /api/inventory/allocate` - 在庫引当
- `POST /api/inventory/release` - 在庫解放

### インボイス

- `POST /api/orders/{id}/generate-invoice` - インボイス生成

### アンバサダー管理

- `POST /api/ambassadors` - アンバサダー作成
- `GET /api/ambassadors` - アンバサダー一覧取得
- `GET /api/ambassadors/commissions` - 成果報酬一覧取得

### 権限管理（RBAC）

- `POST /api/permissions` - 権限作成
- `GET /api/permissions` - 権限一覧取得
- `POST /api/permissions/check` - 権限チェック

### 監視・運用

- `GET /api/metrics` - メトリクス取得
- `GET /health` - ヘルスチェック

**詳細**: [APIリファレンス](./docs/73_API_Reference.md)

---

## 📱 フロントエンド

### 画面構成

1. **ホーム画面（Dashboard）** ✅
   - KPI表示
   - タスクリスト
   - ミル更新フィード
   - クイック発注ボタン

2. **在庫画面（Inventory）** ✅
   - 生地一覧表示
   - 検索・フィルター
   - 在庫ステータス表示

3. **クイック発注画面** ✅
   - ステップ1: 顧客選択
   - ステップ2: 生地選択
   - ステップ3: 金額・納期入力
   - 発注書生成

4. **視覚的発注画面（Visual Ordering）** ⏳
   - 人体図による採寸入力（プレースホルダー）

---

## 🗄️ データベース構成

### 主要テーブル

1. **tenants** - テナント情報
2. **customers** - 顧客情報
3. **orders** - 注文情報
4. **fabrics** - 生地情報
5. **fabric_rolls** - 反物情報
6. **fabric_allocations** - 生地割当
7. **ambassadors** - アンバサダー情報
8. **commissions** - 成果報酬
9. **compliance_documents** - コンプライアンス文書
10. **permissions** - 権限情報
11. **audit_logs** - 監査ログ
12. **audit_log_archives** - 監査ログアーカイブ

**詳細**: [完全システム仕様書](./docs/72_Complete_System_Specification.md#データベース設計)

---

## 🔐 セキュリティ

### 認証・認可

- **Firebase Authentication** (JWT)
- **RBAC**（Role-Based Access Control）
- **マルチテナントデータ分離**

### 監査ログ

- **全操作の記録**
- **改ざん検出**（SHA-256ハッシュ）
- **5年間の保持期間**

---

## 📈 監視・運用

### 構造化ログ

- **JSON形式**
- **トレースID付与**
- **ログレベル**: INFO, WARNING, ERROR

### メトリクス

- **リクエスト数・エラー数**
- **エラー率**
- **平均レイテンシー**
- **データベース接続数**

### アラート

- **エラー率**: 5%以上
- **レイテンシー**: 1秒以上
- **DB接続**: 80%使用

---

## 🚀 デプロイメント

### バックエンド

- **プラットフォーム**: Google Cloud Run
- **データベース**: Cloud SQL (PostgreSQL)
- **ストレージ**: Cloud Storage

### フロントエンド

- **プラットフォーム**: iOS App Store, Google Play Store
- **環境変数**: `API_BASE_URL`

---

## 📚 ドキュメント

### 最重要ドキュメント

1. **[完全システム仕様書](./docs/72_Complete_System_Specification.md)** ⭐
2. **[実装状況サマリー](./docs/74_Implementation_Status.md)** ⭐
3. **[APIリファレンス](./docs/73_API_Reference.md)** ⭐

### 開発・運用

4. **[システム起動ガイド](./docs/67_System_Startup_Guide.md)**
5. **[完全起動手順書](./docs/70_Complete_Startup_Guide.md)**
6. **[次のアクションプラン](./docs/71_Next_Action_Plan.md)**

### ビジネス・戦略

7. **[今後の開発計画](./docs/68_Future_Development_Plan.md)**
8. **[LOI獲得戦略](./docs/69_LOI_Acquisition_Strategy.md)**

**全ドキュメント**: [docs/README.md](./docs/README.md)

---

## 🎯 次のステップ

### 即座に実行

1. **システムの起動テスト**
   ```bash
   ./scripts/start_backend.sh  # バックエンド起動
   ./scripts/start_flutter.sh  # Flutterアプリ起動
   ```

2. **デモの練習**（最低3回）

3. **LOI獲得活動の開始**

---

## 📊 完成度

- **シード調達用MVP**: 100% ✅
- **エンタープライズ基盤**: 90% ✅
- **認証機能**: 50% ⏳

---

**最終更新日**: 2025-01  
**バージョン**: 2.0.0

