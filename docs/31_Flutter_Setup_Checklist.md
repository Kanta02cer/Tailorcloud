# Flutter セットアップ チェックリスト

**作成日**: 2025-01  
**ステータス**: Flutterインストール後のセットアップ手順

---

## ✅ セットアップ手順

### 1. Flutter SDKのインストール

#### 方法A: Homebrew（推奨・簡単）

```bash
brew install --cask flutter
```

#### 方法B: 公式サイトからダウンロード

1. https://docs.flutter.dev/get-started/install/macos にアクセス
2. Flutter SDKをダウンロード
3. 解凍して任意の場所に配置
4. PATHに追加:
   ```bash
   echo 'export PATH="$PATH:$HOME/flutter/bin"' >> ~/.zshrc
   source ~/.zshrc
   ```

### 2. Flutter環境の確認

```bash
# Flutterのバージョンを確認
flutter --version

# 環境を診断
flutter doctor

# 詳細な診断
flutter doctor -v
```

**確認ポイント**:
- ✅ Flutter (Channel stable)
- ✅ Xcode（iOS開発用）
- ✅ Android toolchain（オプション）
- ✅ VS CodeまたはAndroid Studio

### 3. プロジェクトのセットアップ

```bash
cd /Users/wantan/teiloroud-ERPSystem/tailor-cloud-app

# セットアップスクリプトを実行（自動）
./setup.sh
```

**または、手動で実行**:

```bash
# 依存パッケージをインストール
flutter pub get

# コード生成（Freezed, Riverpod Generator）
flutter pub run build_runner build --delete-conflicting-outputs
```

### 4. iOSシミュレーターの準備（オプション）

```bash
# シミュレーターを起動
open -a Simulator

# または、利用可能なデバイスを確認
flutter devices
```

### 5. 動作確認

```bash
# アプリを実行（iOS）
flutter run -d ios

# または、利用可能なデバイスで実行
flutter run
```

---

## 🔍 トラブルシューティング

### 問題: "command not found: flutter"

**解決策**:
```bash
# PATHを確認
echo $PATH

# FlutterのbinディレクトリをPATHに追加
export PATH="$PATH:$HOME/flutter/bin"

# 永続的に追加（~/.zshrcに追加）
echo 'export PATH="$PATH:$HOME/flutter/bin"' >> ~/.zshrc
source ~/.zshrc
```

### 問題: "CocoaPods not installed"

**解決策**:
```bash
sudo gem install cocoapods

# Podをセットアップ
cd tailor-cloud-app/ios
pod setup
pod install
```

### 問題: "No devices available"

**解決策**:
```bash
# iOSシミュレーターを起動
open -a Simulator

# デバイス一覧を確認
flutter devices

# 特定のデバイスで実行
flutter run -d "iPhone 15 Pro"
```

---

## ✅ セットアップ完了チェックリスト

- [ ] Flutter SDKがインストールされている
- [ ] `flutter --version` が動作する
- [ ] `flutter doctor` で主要項目が ✅
- [ ] プロジェクトディレクトリに移動できる
- [ ] `flutter pub get` が成功する
- [ ] コード生成が成功する（.g.dart, .freezed.dart ファイルが生成される）
- [ ] iOSシミュレーターが起動できる（オプション）
- [ ] `flutter run` が実行できる（オプション）

---

## 📝 セットアップ後の確認

### 生成されたファイルを確認

```bash
cd tailor-cloud-app/lib

# コード生成ファイルが存在するか確認
ls -la models/*.g.dart
ls -la models/*.freezed.dart
ls -la providers/*.g.dart
```

### プロジェクト構造の確認

```bash
tree -L 3 lib/  # treeコマンドがインストールされている場合
# または
find lib -name "*.dart" | head -20
```

---

## 🚀 セットアップ後の次のステップ

1. **画面実装を開始**
   - Home画面
   - Inventory画面
   - Visual Ordering画面

2. **動作確認**
   - 各画面の表示確認
   - API連携のテスト

3. **開発を継続**
   - 機能追加
   - バグ修正
   - UI改善

---

**最終更新日**: 2025-01

