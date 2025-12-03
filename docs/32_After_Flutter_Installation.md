# Flutterインストール後の作業

**作成日**: 2025-01  
**前提**: Flutter SDKがインストール済み

---

## 🎯 セットアップ完了後の流れ

### 1. セットアップスクリプトを実行

```bash
cd /Users/wantan/teiloroud-ERPSystem/tailor-cloud-app
./setup.sh
```

このスクリプトが以下を自動実行します：
- Flutter環境の確認
- 依存パッケージのインストール
- コード生成の実行

### 2. セットアップが成功した場合

✅ 以下のメッセージが表示されます：
```
✅ セットアップが完了しました！

次のステップ:
  1. iOSシミュレーターを起動: open -a Simulator
  2. アプリを実行: flutter run
  3. または、利用可能なデバイスを確認: flutter devices
```

### 3. アプリを実行

```bash
# iOSシミュレーターを起動
open -a Simulator

# アプリを実行
flutter run
```

---

## 📋 セットアップ後の確認項目

### 生成されたファイル

以下のファイルが自動生成されているか確認：

```bash
cd tailor-cloud-app/lib

# Freezedファイル
ls models/*.freezed.dart
# 例: fabric.freezed.dart, order.freezed.dart, ambassador.freezed.dart

# JSON Serializationファイル
ls models/*.g.dart
# 例: fabric.g.dart, order.g.dart, ambassador.g.dart

# Riverpod Generatorファイル
ls providers/*.g.dart
# 例: api_client_provider.g.dart, fabric_provider.g.dart, etc.
```

---

## 🔧 セットアップエラーの場合

### エラー1: "flutter pub get" でエラー

**原因**: 依存パッケージの解決に失敗

**解決策**:
```bash
# pubspec.yamlの依存関係を確認
cat pubspec.yaml

# パッケージをクリーンアップして再取得
flutter clean
flutter pub get
```

### エラー2: "build_runner" でエラー

**原因**: コード生成に失敗

**解決策**:
```bash
# 生成ファイルを削除して再生成
flutter pub run build_runner clean
flutter pub run build_runner build --delete-conflicting-outputs
```

### エラー3: "No such file or directory" エラー

**原因**: パスが間違っている

**解決策**:
```bash
# プロジェクトディレクトリにいることを確認
pwd
# 出力: /Users/wantan/teiloroud-ERPSystem/tailor-cloud-app

# 必要に応じて移動
cd /Users/wantan/teiloroud-ERPSystem/tailor-cloud-app
```

---

## ✅ セットアップ成功の確認

### 1. ファイル生成の確認

```bash
# 生成されたファイルを確認
find lib -name "*.g.dart" -o -name "*.freezed.dart" | wc -l
# 期待値: 10個以上
```

### 2. 構文チェック

```bash
# Dartの構文チェック
flutter analyze
```

### 3. ビルド確認

```bash
# ビルド可能か確認（実際に実行しない）
flutter build ios --debug --no-codesign 2>&1 | head -20
```

---

## 🚀 次の実装ステップ

セットアップが完了したら、以下の実装に進みます：

1. **Home画面実装**
   - Dashboard UI
   - KPIカード表示

2. **Inventory画面実装**
   - 生地一覧表示
   - フィルター・検索

3. **Visual Ordering画面実装**
   - 注文作成フロー

---

## 💡 開発のヒント

### ホットリロード

アプリ実行中にコードを変更すると、自動的に反映されます：
- `r` - ホットリロード
- `R` - ホットリスタート
- `q` - 終了

### デバッグログ

```dart
// デバッグログを出力
print('Debug message');
// または
debugPrint('Debug message');
```

### 利用可能なデバイス

```bash
# 利用可能なデバイスを確認
flutter devices

# 特定のデバイスで実行
flutter run -d "iPhone 15 Pro"
```

---

**最終更新日**: 2025-01

