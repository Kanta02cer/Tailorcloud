# GitHub Pages デプロイメントガイド

**作成日**: 2025-01  
**目的**: TailorCloud WebアプリをGitHub Pagesで公開する手順

---

## 🎯 概要

TailorCloudのWebアプリ（suit-mbti-web-app）をGitHub Pagesで静的サイトとして公開します。

**公開URL**: `https://Kanta02cer.github.io/Tailorcloud/`

---

## 📋 前提条件

1. **GitHubリポジトリ**: https://github.com/Kanta02cer/Tailorcloud.git
2. **GitHub Pagesが有効**: リポジトリ設定でGitHub Pagesを有効化
3. **GitHub Actions**: 自動デプロイ用（ワークフローが含まれています）

---

## 🚀 自動デプロイ（推奨）

### セットアップ手順

1. **GitHubリポジトリの設定**

   - リポジトリページにアクセス
   - `Settings` → `Pages` に移動
   - `Source` を `GitHub Actions` に設定
   - `Save` をクリック

2. **自動デプロイの確認**

   - `main` ブランチにプッシュすると自動的にデプロイされます
   - `.github/workflows/deploy-pages.yml` が自動実行されます
   - デプロイ完了後、`Actions` タブで確認できます

3. **公開URLの確認**

   - デプロイ完了後、以下のURLでアクセス可能:
     ```
     https://Kanta02cer.github.io/Tailorcloud/
     ```

---

## 🔧 手動デプロイ

### ローカルでビルドしてデプロイ

```bash
# 1. Webアプリディレクトリに移動
cd suit-mbti-web-app

# 2. 依存関係のインストール
npm install

# 3. ビルド
npm run build

# 4. distディレクトリの内容を確認
ls -la dist/
```

### GitHub Pagesに手動でデプロイ

1. **gh-pagesブランチを作成**（初回のみ）

   ```bash
   # distディレクトリの内容をgh-pagesブランチにプッシュ
   cd suit-mbti-web-app
   npm install -g gh-pages
   gh-pages -d dist
   ```

2. **GitHubリポジトリの設定**

   - `Settings` → `Pages` に移動
   - `Source` を `gh-pages` ブランチに設定
   - `Save` をクリック

---

## 📝 設定ファイル

### vite.config.ts

```typescript
export default defineConfig({
  base: process.env.NODE_ENV === 'production' ? '/Tailorcloud/' : '/',
  // ...
})
```

### main.tsx

```typescript
const basename = import.meta.env.PROD ? '/Tailorcloud' : '';
<BrowserRouter basename={basename}>
  <App />
</BrowserRouter>
```

---

## 🔌 API接続設定

### 本番環境でのAPI接続

GitHub Pagesで公開されたWebアプリは、バックエンドAPIに接続する必要があります。

**オプション1: 環境変数を使用（推奨）**

1. `.env.production` ファイルを作成:

   ```bash
   cd suit-mbti-web-app
   echo "VITE_API_BASE_URL=https://your-backend-api.com" > .env.production
   ```

2. ビルド時に環境変数が読み込まれます

**オプション2: ビルド後の設定**

- `dist/index.html` に直接API URLを埋め込む
- または、設定ファイルを読み込む仕組みを実装

**注意**: GitHub Pagesは静的サイトホスティングのため、バックエンドAPIは別途デプロイする必要があります。

---

## 🧪 ローカルでのテスト

### 本番環境と同じ設定でテスト

```bash
cd suit-mbti-web-app

# 本番環境と同じbase pathでビルド
npm run build

# ビルド結果をプレビュー
npm run preview
```

または、Viteのプレビューサーバーを使用:

```bash
# ビルド
npm run build

# プレビュー（base pathを指定）
npx vite preview --base /Tailorcloud/
```

---

## 📊 デプロイメントの確認

### GitHub Actionsの確認

1. リポジトリページの `Actions` タブにアクセス
2. `Deploy to GitHub Pages` ワークフローを確認
3. 緑色のチェックマークが表示されれば成功

### デプロイメントログの確認

- `Actions` タブ → 最新のワークフロー実行をクリック
- `build` ジョブと `deploy` ジョブのログを確認

### 公開URLの確認

- `https://Kanta02cer.github.io/Tailorcloud/` にアクセス
- ページが正常に表示されることを確認

---

## 🔧 トラブルシューティング

### 404エラーが発生する

**原因**: base pathの設定が正しくない

**解決方法**:
1. `vite.config.ts` の `base` 設定を確認
2. `main.tsx` の `basename` 設定を確認
3. ビルドを再実行

### API接続エラーが発生する

**原因**: バックエンドAPIが起動していない、またはURLが間違っている

**解決方法**:
1. バックエンドAPIが起動していることを確認
2. `.env.production` の `VITE_API_BASE_URL` を確認
3. CORS設定を確認（バックエンド側）

### デプロイが失敗する

**原因**: GitHub Actionsの権限設定が不足している

**解決方法**:
1. リポジトリの `Settings` → `Actions` → `General` に移動
2. `Workflow permissions` を `Read and write permissions` に設定
3. `Allow GitHub Actions to create and approve pull requests` を有効化

---

## 📚 関連ドキュメント

- **[システム起動ガイド](./67_System_Startup_Guide.md)**
- **[APIリファレンス](./73_API_Reference.md)**
- **[完全システム仕様書](./72_Complete_System_Specification.md)**

---

## 🎯 次のステップ

1. **バックエンドAPIのデプロイ**
   - Cloud Run、Heroku、Railwayなどにデプロイ
   - CORS設定を確認

2. **環境変数の設定**
   - 本番環境のAPI URLを設定
   - セキュリティ設定を確認

3. **カスタムドメインの設定**（オプション）
   - GitHub Pagesでカスタムドメインを設定
   - DNS設定を確認

---

**最終更新日**: 2025-01

