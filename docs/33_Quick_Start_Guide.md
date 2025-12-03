# TailorCloud Flutter App クイックスタートガイド

**作成日**: 2025-01  
**対象**: Flutterインストール後、すぐに開発を始めたい方

---

## 🚀 5分で始める

### Step 1: Flutterをインストール

```bash
# Homebrewを使用（推奨）
brew install --cask flutter

# または公式サイトから
# https://docs.flutter.dev/get-started/install/macos
```

### Step 2: セットアップスクリプトを実行

```bash
cd /Users/wantan/teiloroud-ERPSystem/tailor-cloud-app
./setup.sh
```

### Step 3: アプリを起動

```bash
# iOSシミュレーターを起動
open -a Simulator

# アプリを実行
flutter run
```

---

## ✅ セットアップ確認

### Flutterがインストールされているか確認

```bash
flutter --version
# 出力例: Flutter 3.16.0 • channel stable
```

### 環境を確認

```bash
flutter doctor
# 主要項目が ✅ になっていることを確認
```

### セットアップスクリプトが成功したか確認

```bash
# 生成されたファイルを確認
ls lib/models/*.g.dart
ls lib/providers/*.g.dart
```

---

## 📝 よく使うコマンド

### 開発中

```bash
# アプリを実行
flutter run

# 特定のデバイスで実行
flutter run -d "iPhone 15 Pro"

# 利用可能なデバイスを確認
flutter devices
```

### コード生成

```bash
# コード生成を再実行
flutter pub run build_runner build --delete-conflicting-outputs

# 変更を監視して自動生成
flutter pub run build_runner watch
```

### クリーンアップ

```bash
# ビルドキャッシュをクリア
flutter clean

# 依存パッケージを再取得
flutter pub get
```

---

## 🎯 次のステップ

セットアップが完了したら、画面実装に進みます：

1. **Home画面** - Dashboard UI
2. **Inventory画面** - 生地一覧表示
3. **Visual Ordering画面** - 注文作成フロー

---

## 📚 参考ドキュメント

- **インストールガイド**: `docs/28_Flutter_Installation_Guide.md`
- **セットアップチェックリスト**: `docs/31_Flutter_Setup_Checklist.md`
- **インストール後の作業**: `docs/32_After_Flutter_Installation.md`
- **API仕様書**: `docs/20_API_Specification_For_Flutter.md`

---

**最終更新日**: 2025-01

