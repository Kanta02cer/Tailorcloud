# TailorCloud App セットアップ状況

## 現在の状況

- ✅ プロジェクト構造: 作成済み
- ✅ モデルクラス: 実装済み
- ✅ プロバイダー: 実装済み
- ⏳ Flutter SDK: 未インストール

## 次のステップ

### 1. Flutter SDKをインストール

```bash
# Homebrewを使用（推奨）
brew install --cask flutter

# または、公式サイトからダウンロード
# https://docs.flutter.dev/get-started/install/macos
```

詳細は `../docs/28_Flutter_Installation_Guide.md` を参照してください。

### 2. 依存パッケージをインストール

```bash
cd /Users/wantan/teiloroud-ERPSystem/tailor-cloud-app
flutter pub get
```

### 3. コード生成を実行

```bash
flutter pub run build_runner build --delete-conflicting-outputs
```

### 4. 動作確認

```bash
flutter doctor
flutter devices
flutter run
```

## トラブルシューティング

Flutterコマンドが見つからない場合:

```bash
# PATHに追加
echo 'export PATH="$PATH:$HOME/flutter/bin"' >> ~/.zshrc
source ~/.zshrc

# 確認
flutter --version
```

