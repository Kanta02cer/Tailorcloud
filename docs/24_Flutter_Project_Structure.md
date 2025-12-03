# TailorCloud Flutterプロジェクト構造

**作成日**: 2025-01  
**ステータス**: プロジェクトテンプレート作成完了

---

## ✅ 作成されたファイル

### プロジェクト設定

- ✅ `tailor-cloud-app/README.md` - プロジェクト説明とセットアップ手順
- ✅ `tailor-cloud-app/pubspec.yaml` - 依存パッケージ定義
- ✅ `tailor-cloud-app/.gitignore` - Git除外設定

### アプリケーションコード

- ✅ `lib/main.dart` - エントリーポイント
- ✅ `lib/config/app_config.dart` - アプリケーション設定
- ✅ `lib/config/theme.dart` - テーマ・カラーパレット・テキストスタイル
- ✅ `lib/services/api_client.dart` - APIクライアント実装
- ✅ `lib/models/fabric.dart` - 生地モデル（テンプレート）

---

## 📁 プロジェクト構造

```
tailor-cloud-app/
├── lib/
│   ├── main.dart                           ✅ 作成済み
│   ├── config/
│   │   ├── app_config.dart                 ✅ 作成済み
│   │   └── theme.dart                      ✅ 作成済み
│   ├── services/
│   │   ├── api_client.dart                 ✅ 作成済み
│   │   ├── auth_service.dart               ⏳ 実装予定
│   │   └── storage_service.dart            ⏳ 実装予定
│   ├── models/
│   │   ├── fabric.dart                     ✅ 作成済み
│   │   ├── order.dart                      ⏳ 実装予定
│   │   └── ambassador.dart                 ⏳ 実装予定
│   ├── providers/
│   │   ├── auth_provider.dart              ⏳ 実装予定
│   │   ├── fabric_provider.dart            ⏳ 実装予定
│   │   └── order_provider.dart             ⏳ 実装予定
│   ├── screens/
│   │   ├── home/
│   │   │   └── home_screen.dart            ⏳ 実装予定
│   │   ├── inventory/
│   │   │   └── inventory_screen.dart       ⏳ 実装予定
│   │   └── order/
│   │       └── order_create_screen.dart    ⏳ 実装予定
│   └── widgets/
│       ├── fabric_card.dart                ⏳ 実装予定
│       ├── order_card.dart                 ⏳ 実装予定
│       └── kpi_card.dart                   ⏳ 実装予定
├── assets/
│   ├── images/                             ⏳ 準備予定
│   ├── icons/                              ⏳ 準備予定
│   └── fonts/                              ⏳ 準備予定
└── pubspec.yaml                            ✅ 作成済み
```

---

## 🚀 次のステップ

### 1. Flutterインストール・セットアップ

```bash
# Flutter SDKをインストール
# https://docs.flutter.dev/get-started/install/macos

# プロジェクトディレクトリに移動
cd tailor-cloud-app

# 依存パッケージをインストール
flutter pub get

# コード生成（モデルクラス）
flutter pub run build_runner build --delete-conflicting-outputs
```

### 2. Firebase設定

1. Firebase Consoleでプロジェクトを作成
2. iOS用の`GoogleService-Info.plist`を`ios/Runner/`に配置
3. Android用の`google-services.json`を`android/app/`に配置

### 3. 実装の続き

以下の順序で実装を進める：

1. **モデルクラス実装**
   - `lib/models/order.dart`
   - `lib/models/ambassador.dart`

2. **プロバイダー実装**
   - `lib/providers/auth_provider.dart`
   - `lib/providers/fabric_provider.dart`

3. **画面実装**
   - `lib/screens/home/home_screen.dart`
   - `lib/screens/inventory/inventory_screen.dart`

---

## 📝 実装チェックリスト

### プロジェクト基盤 ✅

- [x] プロジェクト構造作成
- [x] 設定ファイル作成
- [x] テーマ・カラーパレット実装
- [x] APIクライアント実装

### 次の実装項目

- [ ] モデルクラス実装（Order, Ambassador）
- [ ] プロバイダー実装
- [ ] 画面実装

---

**最終更新日**: 2025-01

