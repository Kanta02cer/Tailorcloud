# Flutter未インストール時の作業フロー

**作成日**: 2025-01  
**状況**: Flutter SDK未インストール

---

## 📋 現在の状況

Flutter SDKがインストールされていないため、以下の作業は**まだ実行できません**：

- ❌ `flutter pub get` - 依存パッケージインストール
- ❌ `flutter pub run build_runner` - コード生成
- ❌ `flutter run` - アプリ実行

---

## ✅ 現在できること

### 1. コード実装の継続

Flutterがインストールされていなくても、**コードの実装は継続できます**：

- ✅ モデルクラスの追加
- ✅ プロバイダーの実装
- ✅ 画面（Widget）の実装
- ✅ ドキュメント作成

### 2. 実装準備

以下は既に完了しています：

- ✅ プロジェクト構造の作成
- ✅ 設定ファイル（pubspec.yaml等）
- ✅ モデルクラス実装
- ✅ プロバイダー実装

---

## 🔄 推奨作業フロー

### オプション1: Flutterをインストールしてから続行

1. **Flutter SDKをインストール**
   - 参考: `docs/28_Flutter_Installation_Guide.md`

2. **依存パッケージインストール・コード生成**
   ```bash
   cd tailor-cloud-app
   flutter pub get
   flutter pub run build_runner build --delete-conflicting-outputs
   ```

3. **画面実装を開始**

### オプション2: Flutterなしでコード実装を継続

1. **画面（Widget）のテンプレート実装**
   - Dartファイルのみの実装
   - コード生成は後回し

2. **ドキュメント整備**
   - 画面仕様書
   - 実装ガイド

3. **Flutterインストール後**
   - コード生成を一括実行
   - 動作確認

---

## 📝 次の実装項目（Flutterなしでも可能）

### 画面テンプレート実装

以下のファイルは、Flutterがなくても実装できます：

1. **Home画面** (`lib/screens/home/home_screen.dart`)
   - UI構造の定義
   - プロバイダーとの連携ロジック

2. **Inventory画面** (`lib/screens/inventory/inventory_screen.dart`)
   - 生地一覧表示UI
   - フィルター・検索UI

3. **Widgetコンポーネント**
   - `lib/widgets/fabric_card.dart`
   - `lib/widgets/kpi_card.dart`

---

## 🎯 推奨アクション

### 今すぐできること

1. **画面テンプレートの実装**
   - Home画面の基本構造
   - Inventory画面の基本構造

2. **Widgetコンポーネントの実装**
   - 再利用可能なコンポーネント

3. **ドキュメント整備**
   - 画面フロー図
   - コンポーネント仕様

### Flutterインストール後にやること

1. 依存パッケージインストール
2. コード生成実行
3. 動作確認・デバッグ

---

## 💡 メリット

### Flutterなしで進めるメリット

- ✅ 実装を継続できる
- ✅ レビュー可能なコードを準備できる
- ✅ 設計・構造を検討できる

### 注意点

- ⚠️ コード生成ファイル（.g.dart, .freezed.dart）は作成されない
- ⚠️ 構文エラーのチェックができない
- ⚠️ 動作確認ができない

---

**最終更新日**: 2025-01

