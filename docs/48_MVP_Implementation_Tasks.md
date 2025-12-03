# TailorCloud MVP実装タスクリスト

**作成日**: 2025-01  
**目標**: シード調達（3000万〜5000万円）に向けたMVP開発  
**期間**: 6ヶ月

---

## 📋 実装タスク一覧

### 優先度定義

- 🔴 **Critical**: MVPの必須機能（デモに必要）
- 🟡 **High**: MVPで実装推奨（あると良い）
- 🟢 **Low**: Phase 2以降で実装（後回し可）

---

## 1. フリーランス保護法適用型発注書生成 ✅ 最優先

### 1.1 バックエンド実装

#### Task 1.1.1: PDF生成ライブラリの導入

**ファイル**: `tailor-cloud-backend/go.mod`

**作業内容**:
- PDF生成ライブラリの選定・導入
  - 候補1: `github.com/jung-kurt/gofpdf`（軽量、シンプル）
  - 候補2: `github.com/unidoc/unipdf`（高機能、有料版あり）
- `go.mod`に依存関係を追加

**優先度**: 🔴 Critical  
**工数**: 0.5日

---

#### Task 1.1.2: 発注書PDF生成サービスの実装

**ファイル**: `tailor-cloud-backend/internal/service/compliance_service.go`

**作業内容**:
- PDF生成メソッドの実装
- 下請法・フリーランス保護法に準拠したテンプレートの作成
- 必須項目の自動マッピング:
  - 委託をする者の氏名及び住所
  - 給付の内容
  - 報酬の額
  - 支払期日

**既存実装の活用**:
- ✅ `ComplianceRequirement`モデルは定義済み
- ✅ `ConfirmOrder`メソッドは実装済み（PDF生成のトリガー）

**優先度**: 🔴 Critical  
**工数**: 3日

---

#### Task 1.1.3: Cloud StorageへのPDF保存

**ファイル**: `tailor-cloud-backend/internal/service/storage_service.go` (新規)

**作業内容**:
- Cloud StorageへのPDFアップロード機能
- PDFのハッシュ値（SHA-256）計算
- タイムスタンプ付きファイル名の生成

**優先度**: 🔴 Critical  
**工数**: 2日

---

#### Task 1.1.4: 発注書生成APIエンドポイントの追加

**ファイル**: `tailor-cloud-backend/internal/handler/compliance_handler.go` (新規)

**作業内容**:
- `POST /api/orders/{id}/generate-document` エンドポイント
- PDF生成 → Storage保存 → Order更新の流れ

**優先度**: 🔴 Critical  
**工数**: 1日

---

### 1.2 フロントエンド実装

#### Task 1.2.1: 発注書生成画面の実装

**ファイル**: `tailor-cloud-app/lib/screens/order/document_generation_screen.dart` (新規)

**作業内容**:
- 発注確定画面の実装
- 「発注確定」ボタンを押すとPDFが生成される
- PDFプレビュー機能（WebView）
- PDFダウンロード機能

**優先度**: 🔴 Critical  
**工数**: 3日

---

#### Task 1.2.2: PDF表示・ダウンロード機能

**ファイル**: `tailor-cloud-app/lib/services/pdf_service.dart` (新規)

**作業内容**:
- PDFファイルのダウンロード
- PDFビューアの実装（`flutter_pdfview`パッケージ等）
- オフライン対応（キャッシュ）

**優先度**: 🔴 Critical  
**工数**: 2日

---

## 2. 簡易的CRM（顧客管理）✅ 最優先

### 2.1 バックエンド実装（拡張）

#### Task 2.1.1: 顧客管理APIの拡張

**ファイル**: `tailor-cloud-backend/internal/handler/customer_handler.go` (新規)

**作業内容**:
- `GET /api/customers` - 顧客一覧（検索・フィルター）
- `GET /api/customers/{id}` - 顧客詳細
- `POST /api/customers` - 顧客登録
- `PUT /api/customers/{id}` - 顧客更新
- `DELETE /api/customers/{id}` - 顧客削除

**既存実装の活用**:
- ✅ `Customer`モデルは定義済み
- ✅ 基本的なCRUDリポジトリは実装可能

**優先度**: 🔴 Critical  
**工数**: 2日

---

#### Task 2.1.2: 顧客の注文履歴取得API

**ファイル**: `tailor-cloud-backend/internal/handler/customer_handler.go`

**作業内容**:
- `GET /api/customers/{id}/orders` - 顧客の注文履歴一覧
- 注文一覧を時系列で返す

**優先度**: 🔴 Critical  
**工数**: 1日

---

#### Task 2.1.3: 採寸データの保存・取得API拡張

**ファイル**: `tailor-cloud-backend/internal/handler/order_handler.go`

**作業内容**:
- 注文作成時に採寸データ（JSONB）を保存
- 顧客ごとの採寸履歴取得

**既存実装の活用**:
- ✅ `OrderDetails.MeasurementData`はJSONBで保存可能
- ✅ 既存の注文作成APIを拡張

**優先度**: 🔴 Critical  
**工数**: 1日

---

### 2.2 フロントエンド実装

#### Task 2.2.1: 顧客一覧画面の実装

**ファイル**: `tailor-cloud-app/lib/screens/customer/customer_list_screen.dart` (新規)

**作業内容**:
- 顧客一覧のグリッド/リスト表示
- 検索バー（名前、メール、電話番号）
- フィルター機能
- 無限スクロール

**優先度**: 🔴 Critical  
**工数**: 2日

---

#### Task 2.2.2: 顧客詳細画面の実装

**ファイル**: `tailor-cloud-app/lib/screens/customer/customer_detail_screen.dart` (新規)

**作業内容**:
- 顧客情報の表示（名前、連絡先、住所）
- 注文履歴一覧
- 採寸データ履歴
- 編集ボタン

**優先度**: 🔴 Critical  
**工数**: 2日

---

#### Task 2.2.3: 顧客登録・編集画面の実装

**ファイル**: `tailor-cloud-app/lib/screens/customer/customer_edit_screen.dart` (新規)

**作業内容**:
- 顧客情報入力フォーム
- バリデーション
- 保存・更新処理

**優先度**: 🔴 Critical  
**工数**: 2日

---

## 3. スマホ対応（Flutter）✅ 必須

### 3.1 既存実装の最適化

#### Task 3.1.1: レスポンシブデザインの調整

**ファイル**: `tailor-cloud-app/lib/screens/**/*.dart`

**作業内容**:
- スマホ画面サイズに対応したレイアウト調整
- タブレット（iPad）とスマホの両対応
- フォントサイズ・余白の調整

**既存実装の活用**:
- ✅ エンタープライズテーマは実装済み
- ✅ 基本的な画面構造は実装済み

**優先度**: 🔴 Critical  
**工数**: 3日

---

#### Task 3.1.2: オフライン対応の基本実装

**ファイル**: `tailor-cloud-app/lib/services/offline_service.dart` (新規)

**作業内容**:
- Hiveによるローカルデータ保存
- オフライン時のデータ読み込み
- オンライン復帰時の自動同期

**優先度**: 🟡 High  
**工数**: 3日

---

#### Task 3.1.3: iOS/Androidビルド設定

**ファイル**: `tailor-cloud-app/ios/`, `tailor-cloud-app/android/`

**作業内容**:
- iOSビルド設定（Info.plist、AppDelegate等）
- Androidビルド設定（AndroidManifest.xml、build.gradle等）
- 実機テスト

**優先度**: 🔴 Critical  
**工数**: 2日

---

## 4. 統合・テスト

### 4.1 エンドツーエンドテスト

#### Task 4.1.1: 統合テストの作成

**ファイル**: `tailor-cloud-backend/internal/handler/*_test.go`

**作業内容**:
- 発注書生成のE2Eテスト
- 顧客管理のE2Eテスト
- エラーハンドリングのテスト

**優先度**: 🔴 Critical  
**工数**: 3日

---

#### Task 4.1.2: ユーザーテストの実施

**作業内容**:
- 5店舗程度でのユーザーテスト
- フィードバック収集
- バグ修正

**優先度**: 🔴 Critical  
**工数**: 1週間（テスト期間含む）

---

## 5. デモ準備

### 5.1 デモ環境の構築

#### Task 5.1.1: デモ環境のセットアップ

**作業内容**:
- ステージング環境の構築（GCP）
- デモデータの準備
- デモアカウントの作成

**優先度**: 🔴 Critical  
**工数**: 2日

---

#### Task 5.1.2: デモ動画の作成

**作業内容**:
- 発注書生成のデモ動画（3分程度）
- 顧客管理のデモ動画（2分程度）
- スクリーンショットの準備

**優先度**: 🔴 Critical  
**工数**: 1日

---

## 6. 開発スケジュール（6ヶ月）

### Month 1-2: バックエンド基盤

**Week 1-2**: 発注書生成API
- Task 1.1.1: PDF生成ライブラリ導入 ✅
- Task 1.1.2: PDF生成サービス実装 ✅
- Task 1.1.3: Cloud Storage保存 ✅

**Week 3-4**: 顧客管理API拡張
- Task 2.1.1: 顧客管理API拡張 ✅
- Task 2.1.2: 注文履歴取得 ✅
- Task 2.1.3: 採寸データ保存 ✅

**Week 5-6**: API統合テスト
- Task 4.1.1: 統合テスト作成 ✅

**Week 7-8**: バックエンド完成
- バグ修正・ドキュメント整備 ✅

---

### Month 3-4: フロントエンド開発

**Week 9-10**: 発注書生成画面
- Task 1.2.1: 発注書生成画面 ✅
- Task 1.2.2: PDF表示・ダウンロード ✅

**Week 11-12**: 顧客管理画面
- Task 2.2.1: 顧客一覧画面 ✅
- Task 2.2.2: 顧客詳細画面 ✅
- Task 2.2.3: 顧客登録・編集画面 ✅

**Week 13-14**: 画面統合
- Task 3.1.1: レスポンシブデザイン調整 ✅

**Week 15-16**: フロントエンド完成
- Task 3.1.3: iOS/Androidビルド設定 ✅
- バグ修正 ✅

---

### Month 5-6: 統合・テスト・デモ準備

**Week 17-18**: 統合テスト
- Task 4.1.2: ユーザーテスト実施 ✅

**Week 19-20**: デモ準備
- Task 5.1.1: デモ環境構築 ✅
- Task 5.1.2: デモ動画作成 ✅

**Week 21-22**: 投資家向け資料作成
- Pitch Deck作成 ✅

**Week 23-24**: 最終調整
- バグ修正・パフォーマンス最適化 ✅

---

## 7. 既存実装の活用状況

### 活用できる既存実装

1. ✅ **バックエンド基盤**
   - Goプロジェクト構造
   - Firebase認証
   - PostgreSQL接続
   - 基本的なCRUD API

2. ✅ **データモデル**
   - `Customer`モデル
   - `Order`モデル
   - `ComplianceRequirement`モデル

3. ✅ **フロントエンド基盤**
   - Flutterプロジェクト構造
   - エンタープライズテーマ
   - APIクライアント
   - Riverpod状態管理

### 拡張が必要な既存実装

1. ⚠️ **PDF生成機能**
   - 構造のみ定義済み
   - 実装が必要

2. ⚠️ **顧客管理API**
   - モデルは定義済み
   - ハンドラーは未実装

3. ⚠️ **スマホ対応**
   - 基本的な画面は実装済み
   - レスポンシブ調整が必要

---

## 8. リスク管理

### 技術的リスク

- **PDF生成の品質**: 法的要件を満たすPDFが生成できるか？
  - **対策**: 法務専門家とのレビューを実施

- **スマホ対応の品質**: iOS/Androidで正常に動作するか？
  - **対策**: 実機テストを早期に実施

### ビジネスリスク

- **ユーザーのニーズ**: 本当に必要とされる機能か？
  - **対策**: ユーザーテストを早期に実施

---

## 9. 次のアクション

### 即座に実行すべきこと

1. **開発チームの編成**
   - バックエンドエンジニア（Go経験者）1名
   - フロントエンドエンジニア（Flutter経験者）1名

2. **Week 1の開発開始**
   - Task 1.1.1: PDF生成ライブラリ導入
   - Task 1.1.2: PDF生成サービス実装

3. **ユーザーテストの準備**
   - 5店舗程度の協力店舗を確保

---

**最終更新日**: 2025-01  
**ステータス**: ✅ 実装タスクリスト完了

**次のアクション**: Week 1の開発開始（PDF生成ライブラリ導入）

