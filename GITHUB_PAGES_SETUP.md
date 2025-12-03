# GitHub Pages セットアップ手順

**作成日**: 2025-01  
**目的**: GitHub PagesでWebアプリを公開するための簡単なセットアップ手順

---

## 🚀 クイックセットアップ（5分）

### ステップ1: GitHubリポジトリの設定

1. GitHubリポジトリにアクセス: https://github.com/Kanta02cer/Tailorcloud
2. `Settings` タブをクリック
3. 左メニューから `Pages` を選択
4. `Source` セクションで:
   - `Deploy from a branch` を選択
   - または `GitHub Actions` を選択（推奨）
5. `Save` をクリック

### ステップ2: 自動デプロイの確認

- `main` ブランチにプッシュすると、自動的にGitHub Actionsが実行されます
- `.github/workflows/deploy-pages.yml` が自動的にデプロイを実行します
- デプロイには数分かかります

### ステップ3: 公開URLの確認

デプロイ完了後、以下のURLでアクセスできます:

```
https://Kanta02cer.github.io/Tailorcloud/
```

---

## 📋 詳細手順

### GitHub Actionsを使用する場合（推奨）

1. **リポジトリの設定**
   - `Settings` → `Pages` → `Source` を `GitHub Actions` に設定

2. **ワークフローの確認**
   - `.github/workflows/deploy-pages.yml` が存在することを確認
   - `main` ブランチにプッシュすると自動実行されます

3. **デプロイの確認**
   - `Actions` タブでデプロイの進行状況を確認
   - 緑色のチェックマークが表示されれば成功

### 手動でデプロイする場合

```bash
# 1. Webアプリディレクトリに移動
cd suit-mbti-web-app

# 2. 依存関係のインストール
npm install

# 3. ビルド
npm run build:pages

# 4. gh-pagesパッケージをインストール（初回のみ）
npm install -g gh-pages

# 5. distディレクトリをgh-pagesブランチにデプロイ
gh-pages -d dist
```

---

## 🔧 トラブルシューティング

### 404エラーが発生する

**原因**: base pathの設定が正しくない

**解決方法**:
1. `vite.config.ts` の `base` 設定を確認
2. `main.tsx` の `basename` 設定を確認
3. ビルドを再実行

### デプロイが失敗する

**原因**: GitHub Actionsの権限設定が不足している

**解決方法**:
1. リポジトリの `Settings` → `Actions` → `General` に移動
2. `Workflow permissions` を `Read and write permissions` に設定
3. `Allow GitHub Actions to create and approve pull requests` を有効化

### API接続エラーが発生する

**原因**: バックエンドAPIが起動していない、またはURLが間違っている

**解決方法**:
1. バックエンドAPIが起動していることを確認
2. 環境変数 `VITE_API_BASE_URL` を設定
3. CORS設定を確認（バックエンド側）

---

## 📚 詳細ドキュメント

- **[GitHub Pages デプロイメントガイド](./docs/99_GitHub_Pages_Deployment.md)** - 詳細な設定手順
- **[システム起動ガイド](./docs/67_System_Startup_Guide.md)** - ローカル開発環境のセットアップ

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

