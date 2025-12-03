# TailorCloud Flutter App

**作成日**: 2025-01  
**バージョン**: 1.0.0

---

## 📋 プロジェクト概要

TailorCloudのモバイルアプリ（iPadアプリ）です。

Phase 1.2の目標: 「ド素人でも間違えない最強の入力インターフェース」を実装

---

## 🛠️ セットアップ

### 必要な環境

- Flutter SDK: 3.16.0以上
- Dart SDK: 3.2.0以上
- Xcode: 15.0以上（iOS開発用）
- Android Studio / VS Code

### Flutterインストール

#### 方法1: Homebrew（推奨）

```bash
brew install --cask flutter
```

#### 方法2: 公式サイトからダウンロード

```bash
# https://docs.flutter.dev/get-started/install/macos からダウンロード
# 解凍後、PATHに追加
export PATH="$PATH:$HOME/flutter/bin"
echo 'export PATH="$PATH:$HOME/flutter/bin"' >> ~/.zshrc
source ~/.zshrc
```

詳細は `../docs/28_Flutter_Installation_Guide.md` を参照してください。

### セットアップスクリプトを実行（推奨）

```bash
cd /Users/wantan/teiloroud-ERPSystem/tailor-cloud-app
./setup.sh
```

このスクリプトが以下を自動実行します：
- Flutter環境の確認
- 依存パッケージのインストール
- コード生成の実行

### 手動セットアップ

セットアップスクリプトを使用しない場合は、以下を手動で実行：

```bash
cd /Users/wantan/teiloroud-ERPSystem/tailor-cloud-app

# 依存パッケージをインストール
flutter pub get

# コード生成（Freezed, Riverpod Generator, JSON Serializable）
flutter pub run build_runner build --delete-conflicting-outputs
```

### Firebase設定

1. Firebase Consoleでプロジェクトを作成
2. iOS用の`GoogleService-Info.plist`を`ios/Runner/`に配置
3. Android用の`google-services.json`を`android/app/`に配置

---

## 📁 プロジェクト構成

```
tailor_cloud_app/
├── lib/
│   ├── main.dart
│   ├── config/
│   │   ├── app_config.dart
│   │   └── theme.dart
│   ├── models/
│   │   ├── fabric.dart
│   │   ├── order.dart
│   │   └── ambassador.dart
│   ├── services/
│   │   ├── api_client.dart
│   │   ├── auth_service.dart
│   │   └── storage_service.dart
│   ├── providers/
│   │   ├── auth_provider.dart
│   │   ├── fabric_provider.dart
│   │   └── order_provider.dart
│   ├── screens/
│   │   ├── home/
│   │   │   └── home_screen.dart
│   │   ├── inventory/
│   │   │   └── inventory_screen.dart
│   │   └── order/
│   │       └── order_create_screen.dart
│   ├── widgets/
│   │   ├── fabric_card.dart
│   │   ├── order_card.dart
│   │   └── kpi_card.dart
│   └── utils/
│       ├── constants.dart
│       └── validators.dart
├── assets/
│   ├── images/
│   └── icons/
└── pubspec.yaml
```

---

## 🚀 開発

### コード生成

```bash
# モデルクラスのコード生成
flutter pub run build_runner build --delete-conflicting-outputs
```

### アプリ起動

```bash
# iOSシミュレーターで起動
flutter run -d ios

# Androidエミュレーターで起動
flutter run -d android
```

---

## 📚 参考ドキュメント

- API仕様書: `../docs/20_API_Specification_For_Flutter.md`
- 開発ガイド: `../docs/21_Flutter_Development_Guide.md`
- 実装計画: `../docs/22_Phase1_2_Implementation_Plan.md`

---

**最終更新日**: 2025-01

