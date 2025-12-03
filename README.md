# TailorCloud ERP System

**オーダースーツ業界向けマルチテナント型ERPシステム**

[![License](https://img.shields.io/badge/license-Proprietary-red.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org)
[![Flutter Version](https://img.shields.io/badge/flutter-3.16+-blue.svg)](https://flutter.dev)

---

## 🎯 プロジェクト概要

TailorCloudは、オーダースーツ業界向けの**マルチテナント型ERPシステム**です。

**核心価値**: 「フリーランス保護法対応の発注書をスマホで3分で作れる」

### 主要機能

- ✅ **クイック発注**: 3分で発注書作成
- ✅ **Roll管理システム**: 反物単位の在庫管理
- ✅ **コンプライアンスエンジン**: 下請法対応発注書PDF生成
- ✅ **インボイス制度対応**: 適格インボイスPDF生成
- ✅ **RBAC（権限管理）**: 細かい権限管理
- ✅ **監査ログ**: 全操作の記録
- ✅ **監視・運用基盤**: 構造化ログ、メトリクス収集

---

## 📁 プロジェクト構成

```
teiloroud-ERPSystem/
├── tailor-cloud-backend/     # バックエンドAPI (Go)
├── tailor-cloud-app/         # Flutterモバイルアプリ
├── suit-mbti-web-app/        # React Webアプリ
├── docs/                     # ドキュメント
├── scripts/                  # セットアップ・起動スクリプト
└── migrations/               # データベースマイグレーション
```

---

## 🌐 GitHub Pagesで公開

TailorCloud WebアプリはGitHub Pagesで静的サイトとして公開できます。

**公開URL**: https://Kanta02cer.github.io/Tailorcloud/

### 自動デプロイ

1. リポジトリの `Settings` → `Pages` に移動
2. `Source` を `GitHub Actions` に設定
3. `main` ブランチにプッシュすると自動的にデプロイされます

詳細は [GitHub Pages デプロイメントガイド](./docs/99_GitHub_Pages_Deployment.md) を参照してください。

---

## 🚀 クイックスタート

### 前提条件

- **Go 1.24+** (バックエンドAPI用)
- **Flutter 3.16.0+** (モバイルアプリ用)
- **PostgreSQL 17+** (オプション - Firestoreモードでも動作可能)
- **Node.js 18+** (Webアプリ用、オプション)

### 1. リポジトリのクローン

```bash
git clone https://github.com/Kanta02cer/Tailorcloud.git
cd Tailorcloud
```

### 2. システム状態の確認

```bash
./scripts/check_system.sh
```

### 3. 環境変数のセットアップ

```bash
./scripts/setup_local_environment.sh
```

これにより `.env.local` ファイルが作成されます。

### 4. バックエンドAPIの起動

```bash
./scripts/start_backend.sh
```

または手動で:

```bash
cd tailor-cloud-backend
export PORT=8080
go run cmd/api/main.go
```

**確認**: http://localhost:8080/health にアクセスして "OK" が返ってくることを確認

### 5. Flutterアプリの起動

```bash
./scripts/start_flutter.sh
```

または手動で:

```bash
cd tailor-cloud-app
export API_BASE_URL=http://localhost:8080
flutter run --dart-define=API_BASE_URL=http://localhost:8080
```

---

## 📚 詳細ドキュメント

### 最重要ドキュメント

1. **[システム概要](./SYSTEM_OVERVIEW.md)** ⭐
2. **[完全システム仕様書](./docs/72_Complete_System_Specification.md)** ⭐
3. **[実装状況サマリー](./docs/74_Implementation_Status.md)** ⭐
4. **[APIリファレンス](./docs/73_API_Reference.md)** ⭐

### 開発・運用

- **[システム起動ガイド](./docs/67_System_Startup_Guide.md)**
- **[完全起動手順書](./docs/70_Complete_Startup_Guide.md)**
- **[次のアクションプラン](./docs/71_Next_Action_Plan.md)**

### ビジネス・戦略

- **[今後の開発計画](./docs/68_Future_Development_Plan.md)**
- **[LOI獲得戦略](./docs/69_LOI_Acquisition_Strategy.md)**

**全ドキュメント**: [docs/README.md](./docs/README.md)

---

## 🏗️ アーキテクチャ

```
┌─────────────────────────────────────────────────────────────┐
│                    TailorCloud System                        │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌─────────────────┐              ┌──────────────────┐      │
│  │  Flutter App    │              │   Web Portal     │      │
│  │  (Mobile/Tablet)│              │   (React)        │      │
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

### 監視・運用

- `GET /api/metrics` - メトリクス取得
- `GET /health` - ヘルスチェック

**詳細**: [APIリファレンス](./docs/73_API_Reference.md)

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

## 🧪 テスト

### バックエンド

```bash
cd tailor-cloud-backend
go test ./...
go test -cover ./...
```

### Flutterアプリ

```bash
cd tailor-cloud-app
flutter test
```

---

## 📊 システム規模

### コード統計

- **バックエンド**: 約15,000行（Go）
- **フロントエンド**: 約5,000行（Dart）
- **合計**: 約20,000行

### データベース

- **テーブル数**: 12テーブル
- **インデックス数**: 50+インデックス
- **マイグレーション**: 15ファイル

### API

- **エンドポイント数**: 30+エンドポイント
- **認証不要**: 2（/health, /api/metrics）
- **認証必須**: 28+

---

## 🎯 完成度

- **シード調達用MVP**: 100% ✅
- **エンタープライズ基盤**: 90% ✅
- **認証機能**: 50% ⏳

---

## 📝 ライセンス

Copyright © 2025 Regalis Japan Group. All Rights Reserved.

---

## 📞 サポート

詳細なドキュメントは [docs/](./docs/) ディレクトリを参照してください。

**最終更新日**: 2025-01  
**バージョン**: 2.0.0

