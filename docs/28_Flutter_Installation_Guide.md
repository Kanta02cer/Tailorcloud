# Flutter インストールガイド（macOS）

**作成日**: 2025-01  
**対象OS**: macOS

---

## 🚀 Flutter SDK インストール手順

### 方法1: 公式サイトからダウンロード（推奨）

1. **Flutter公式サイトからダウンロード**
   - https://docs.flutter.dev/get-started/install/macos

2. **SDKを解凍**
   ```bash
   cd ~
   unzip ~/Downloads/flutter_macos_*.zip
   ```

3. **PATHに追加**
   ```bash
   # ~/.zshrc に追加
   export PATH="$PATH:$HOME/flutter/bin"
   
   # 反映
   source ~/.zshrc
   ```

4. **動作確認**
   ```bash
   flutter doctor
   ```

### 方法2: Homebrewを使用（簡単）

```bash
# Homebrewでインストール
brew install --cask flutter

# 動作確認
flutter doctor
```

---

## ✅ 必要な前提条件

### Xcode（iOS開発用）

```bash
# Xcodeをインストール
# App Storeからインストール、または：
xcode-select --install

# Xcodeライセンスに同意
sudo xcodebuild -license accept

# コマンドラインツールをインストール
xcode-select --install
```

### Android Studio（Android開発用、オプション）

```bash
# Homebrewでインストール
brew install --cask android-studio
```

---

## 🔧 セットアップ後の確認

### Flutter Doctorで確認

```bash
flutter doctor
```

すべての項目が ✅ になっていることを確認してください。

### プロジェクトでの動作確認

```bash
cd /Users/wantan/teiloroud-ERPSystem/tailor-cloud-app

# 依存パッケージをインストール
flutter pub get

# コード生成（Freezed, Riverpod Generator）
flutter pub run build_runner build --delete-conflicting-outputs

# 動作確認（iOSシミュレーター）
flutter doctor -v
```

---

## ⚠️ よくある問題と解決策

### 問題1: "command not found: flutter"

**原因**: PATHにFlutterのbinディレクトリが追加されていない

**解決策**:
```bash
# 1. Flutterのインストール場所を確認
which flutter  # 何も表示されない場合は未インストール

# 2. PATHに追加（~/.zshrc）
echo 'export PATH="$PATH:$HOME/flutter/bin"' >> ~/.zshrc
source ~/.zshrc

# 3. 再度確認
flutter --version
```

### 問題2: "CocoaPods not installed"

**解決策**:
```bash
# CocoaPodsをインストール
sudo gem install cocoapods

# Podをセットアップ
cd tailor-cloud-app/ios
pod setup
```

### 問題3: "No devices available"

**解決策**:
```bash
# iOSシミュレーターを起動
open -a Simulator

# または、利用可能なデバイスを確認
flutter devices
```

---

## 📝 セットアップ後の次のステップ

### 1. 依存パッケージインストール

```bash
cd /Users/wantan/teiloroud-ERPSystem/tailor-cloud-app
flutter pub get
```

### 2. コード生成

```bash
flutter pub run build_runner build --delete-conflicting-outputs
```

### 3. 動作確認

```bash
# iOSシミュレーターで起動
flutter run -d ios

# または、利用可能なデバイスを確認
flutter devices
```

---

## 🔗 参考リンク

- **Flutter公式ドキュメント**: https://docs.flutter.dev/get-started/install/macos
- **Flutter Doctor**: https://docs.flutter.dev/get-started/install/macos#verify-setup
- **iOS開発セットアップ**: https://docs.flutter.dev/get-started/install/macos#ios-setup
- **Android開発セットアップ**: https://docs.flutter.dev/get-started/install/macos#android-setup

---

## 💡 トラブルシューティング

### すべてのエラーを確認

```bash
# 詳細な診断情報を表示
flutter doctor -v
```

### Flutter SDKのバージョン確認

```bash
flutter --version
```

### チャンネル確認・変更

```bash
# 現在のチャンネルを確認
flutter channel

# 安定版に切り替え
flutter channel stable
flutter upgrade
```

---

**最終更新日**: 2025-01

