# TailorCloud: ギャップ分析書

**Status**: Prototype (Visual Only) → Production (Enterprise Ready)  
**Date**: 2025-01  
**Version**: 1.0.0

---

## 📋 エグゼクティブサマリー

現在のTailorCloudプロトタイプは「UI/UXの理想形」としては完成していますが、システムとしての「中身（バックエンド・データ構造）」は未完成です。

このドキュメントでは、**エンタープライズ導入可能なレベル（Production Ready）**に到達するために必要な機能と技術的ギャップを明確化します。

---

## 1. Frontend Gaps (フロントエンドの不足)

### 1.1 State Management (状態管理) ❌

**現状**:
- 静的HTMLプロトタイプ（`TailorCloud_Prototype_v2.html`）
- 画面をリロードすると入力内容が消える
- アプリケーション状態の永続化なし

**あるべき姿**:
- **Riverpod (Flutter)** または **Redux/Zustand (React)** による状態管理
- アプリ内のデータを永続化・同期
- オフライン対応（Hive/Drift for Flutter）
- セッション管理（ログイン状態の保持）

**優先度**: 🔴 Critical  
**工数見積**: 2週間

---

### 1.2 Error Handling (エラー処理) ❌

**現状**:
- 正常系（Happy Path）のみ実装
- エラー時のUIフロー未定義

**あるべき姿**:
- **通信圏外**時のエラーハンドリング
  - ネットワークエラー表示
  - オフラインモードへの切り替え
  - データの自動再同期
  
- **在庫引当失敗**時のエラーハンドリング
  - 「在庫が他の注文で確保されました」通知
  - 代替生地の提案
  
- **セッション切れ**時のエラーハンドリング
  - 自動ログアウト
  - ログイン画面への誘導

**優先度**: 🔴 Critical  
**工数見積**: 1週間

---

### 1.3 Performance Optimization (パフォーマンス最適化) ❌

**現状**:
- 少ないデータ量での描画のみ
- 画像の遅延読み込みなし
- 仮想スクロール未実装

**あるべき姿**:
- **仮想スクロール (Virtual Scrolling)**
  - 数万件の生地データをスクロールしてもカクつかない
  - Flutter: `ListView.builder` / React: `react-window`
  
- **画像キャッシュ戦略**
  - Cloud Storageからの画像をキャッシュ
  - サムネイルとフルサイズの段階的読み込み
  
- **データのページネーション**
  - APIからの取得をページ単位に分割
  - 無限スクロール実装

**優先度**: 🟡 High  
**工数見積**: 2週間

---

## 2. Backend Gaps (バックエンドの不足)

### 2.1 API Implementation (API実装) ⚠️ 部分実装

**現状**:
- ✅ 基本的なCRUD APIは実装済み（Order, Fabric, Ambassador）
- ❌ 在庫引当ロジック未実装
- ❌ 反物（Roll）管理API未実装
- ❌ 同時発注時の排他制御未実装

**あるべき姿**:

#### A. 在庫引当API

```go
POST /api/inventory/allocate
{
  "order_id": "uuid",
  "fabric_id": "uuid",
  "required_length": 3.2,
  "location_id": "uuid"
}

Response:
{
  "allocation_id": "uuid",
  "allocated_rolls": [
    {
      "roll_id": "uuid",
      "allocated_length": 3.2,
      "remaining_length": 46.8
    }
  ]
}
```

#### B. 反物管理API

```go
GET /api/fabric-rolls?fabric_id=uuid&status=AVAILABLE

POST /api/fabric-rolls
{
  "fabric_id": "uuid",
  "lot_number": "LOT-2025-001",
  "initial_length": 50.0,
  "location_id": "uuid"
}
```

**優先度**: 🔴 Critical  
**工数見積**: 3週間

---

### 2.2 Authentication & Authorization (認証・認可) ⚠️ 部分実装

**現状**:
- ✅ Firebase認証は実装済み
- ✅ RBACミドルウェア実装済み
- ❌ TenantIDに基づくデータ分離が不完全
- ❌ 店長は「発注可」、アルバイトは「閲覧のみ」の細かい権限制御未実装

**あるべき姿**:

#### A. テナント分離の強化

```go
// すべてのクエリにtenant_idを含める
SELECT * FROM orders WHERE tenant_id = $1 AND ...

// Row Level Security (RLS) ポリシーを設定
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON orders
  USING (tenant_id = current_setting('app.current_tenant_id')::UUID);
```

#### B. 権限の細分化

| ロール | 発注作成 | 発注確定 | 価格変更 | 在庫修正 | レポート閲覧 |
|--------|---------|---------|---------|---------|------------|
| Owner | ✅ | ✅ | ✅ | ✅ | ✅ |
| Store_Manager | ✅ | ✅ | ❌ | ❌ | ✅ |
| Staff | ✅ | ❌ | ❌ | ❌ | 自店のみ |
| Factory_Manager | ❌ | ✅ (受注) | ❌ | ❌ | 自工場のみ |
| Worker | ❌ | ❌ | ❌ | ❌ | ❌ |

**優先度**: 🔴 Critical  
**工数見積**: 2週間

---

### 2.3 Transaction Management (排他制御) ❌

**課題**:
- A店とB店が同時に「残り1着の生地」を発注した際、ダブルブッキングが起きる

**あるべき姿**:

#### A. 楽観的ロック（Optimistic Locking）

```sql
-- fabric_rollsテーブルにversionカラムを追加
ALTER TABLE fabric_rolls ADD COLUMN version INTEGER DEFAULT 1;

-- 更新時にversionをチェック
UPDATE fabric_rolls
SET current_length = current_length - $1,
    version = version + 1
WHERE id = $2
  AND version = $3  -- 楽観的ロック
  AND current_length >= $1;

-- 更新件数が0なら競合エラー
```

#### B. 悲観的ロック（Pessimistic Locking）

```sql
-- トランザクション内で行をロック
BEGIN;
SELECT * FROM fabric_rolls
WHERE id = $1
FOR UPDATE;  -- 行レベルロック

-- 在庫を減らす
UPDATE fabric_rolls SET current_length = current_length - $2;

COMMIT;
```

#### C. DBレベルのトランザクション分離レベル

```sql
-- PostgreSQLのデフォルト（READ COMMITTED）を維持
-- 必要に応じてSERIALIZABLEを使用
SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;
```

**優先度**: 🔴 Critical  
**工数見積**: 2週間

---

## 3. Infrastructure Gaps (インフラの不足)

### 3.1 CI/CD Pipeline ❌

**現状**:
- 手動でのデプロイ
- テスト自動化なし

**あるべき姿**:
- **GitHub Actions** または **Cloud Build** による自動デプロイ
- ユニットテスト、統合テストの自動実行
- ステージング環境への自動デプロイ
- 本番環境への承認フロー

**優先度**: 🟡 High  
**工数見積**: 1週間

---

### 3.2 Monitoring & Logging ❌

**現状**:
- エラーが起きても検知できない
- ログの集約なし

**あるべき姿**:
- **Datadog** または **Cloud Logging** の導入
- エラーアラートの設定
- パフォーマンスメトリクスの可視化
- ダッシュボードでのリアルタイム監視

**優先度**: 🟡 High  
**工数見積**: 1週間

---

## 4. Database Schema Gaps (データベーススキーマの不足)

### 4.1 反物管理（Roll Management）❌

**現状**:
- `fabrics`テーブルに`stock_amount`（総量）のみ
- 物理的な「巻き（Roll）」単位の管理なし

**あるべき姿**:
- `fabric_rolls`テーブルの追加
- 各巻きの在庫状況を個別管理
- キレ（端尺）問題の解決

**詳細**: `docs/44_jp_enterprise_requirements.md` 参照

**優先度**: 🔴 Critical  
**工数見積**: 3週間

---

### 4.2 インボイス制度対応 ❌

**現状**:
- インボイス登録番号の管理なし
- 消費税額の正確な計算なし

**あるべき姿**:
- `tenants`テーブルに`invoice_no`（登録番号）追加
- 請求書PDFに「T番号」と「税率ごとの消費税額」を記載
- 端数処理（切り捨て・切り上げ）のルール統一

**詳細**: `docs/44_jp_enterprise_requirements.md` 参照

**優先度**: 🔴 Critical  
**工数見積**: 2週間

---

## 5. Compliance & Legal Gaps (コンプライアンス・法務の不足)

### 5.1 下請法対応 ⚠️ 部分実装

**現状**:
- ✅ 3条書面の自動生成ロジックは定義済み
- ❌ PDF生成の実装未完了
- ❌ 修正注文書の履歴管理なし

**あるべき姿**:
- PDF生成ライブラリ（Go: `gofpdf` / `unidoc`）の導入
- 修正注文書は履歴として残す（上書き禁止）
- タイムスタンプ付きで保存

**優先度**: 🔴 Critical  
**工数見積**: 2週間

---

### 5.2 監査ログの強化 ⚠️ 部分実装

**現状**:
- ✅ 基本的な監査ログは実装済み
- ❌ 5年以上の保存期間設定なし
- ❌ WORMストレージ（改ざん不能）への移行なし

**あるべき姿**:
- 監査ログを5年以上保存
- Cloud StorageのWORM（Write Once Read Many）機能を使用
- 定期的なアーカイブ処理

**優先度**: 🟡 High  
**工数見積**: 1週間

---

## 6. 優先度マトリクス

| 機能 | 優先度 | 工数 | 依存関係 |
|------|--------|------|----------|
| 反物管理（Roll Management） | 🔴 Critical | 3週間 | DBスキーマ設計 |
| 在庫引当ロジック | 🔴 Critical | 3週間 | 反物管理 |
| 排他制御（Transaction） | 🔴 Critical | 2週間 | 在庫引当ロジック |
| インボイス制度対応 | 🔴 Critical | 2週間 | DBスキーマ設計 |
| 下請法PDF生成 | 🔴 Critical | 2週間 | インボイス制度 |
| 権限管理の細分化 | 🔴 Critical | 2週間 | 認証基盤 |
| State Management | 🔴 Critical | 2週間 | フロントエンド基盤 |
| Error Handling | 🔴 Critical | 1週間 | State Management |
| パフォーマンス最適化 | 🟡 High | 2週間 | State Management |
| CI/CD Pipeline | 🟡 High | 1週間 | なし |
| Monitoring & Logging | 🟡 High | 1週間 | インフラ |

---

## 7. 推奨開発ロードマップ

### Month 1: The Backbone (DB & API)

**Week 1-2: DBスキーマ設計**
- エンタープライズスキーマの確定（`schema/enterprise_schema.sql`）
- 反物管理テーブルの追加
- インボイス制度対応のフィールド追加

**Week 3-4: 在庫引当API**
- 反物管理APIの実装
- 在庫引当ロジックの実装
- 排他制御の実装

**検証**: 「A店で発注したら、B店から見て在庫が減っているか？」をAPIレベルでテスト

---

### Month 2: The Core (Logic & Auth)

**Week 1-2: 権限管理の細分化**
- ロールごとの権限制御の実装
- Row Level Security (RLS) の設定
- 権限テストの作成

**Week 3-4: コンプライアンス機能**
- 下請法PDF生成の実装
- インボイス制度対応
- 修正注文書の履歴管理

**検証**: 複数の端末から同時にアクセスしてもデータが壊れないか（負荷テスト）

---

### Month 3: The Integration (Frontend Connection)

**Week 1-2: State Management実装**
- Riverpodの導入（Flutter）
- オフライン対応
- エラーハンドリング

**Week 3-4: フロントエンド統合**
- 既存のHTMLプロトタイプをFlutterコードに書き換え
- 本物のAPIと接続
- パフォーマンス最適化

**検証**: 店舗での実運用テスト（PoC）

---

## 8. 結論

現在のTailorCloudプロトタイプは「見た目」としては完成していますが、「中身」は約30%程度の完成度です。

**最優先事項**:
1. **DBスキーマ設計の確定**（これさえ間違えなければ、後からいくらでも修正が効く）
2. **在庫引当ロジックの実装**（ビジネスの根幹）
3. **排他制御の実装**（データ整合性の保証）

この3つを完了すれば、エンタープライズ導入に向けた基盤が整います。

---

**次のアクション**:
1. `docs/44_jp_enterprise_requirements.md` を参照して日本固有の要件を確認
2. `schema/enterprise_schema.sql` を参照してDBスキーマを確認
3. バックエンドエンジニアを採用（またはアサイン）
4. Month 1の開発を開始

---

**最終更新日**: 2025-01  
**ステータス**: ✅ ギャップ分析完了

