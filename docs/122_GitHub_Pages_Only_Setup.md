# GitHub Pagesのみでのデプロイ設定

**作成日**: 2025-01  
**目的**: Vercelを削除し、GitHub Pagesのみでシステムを公開する設定

---

## 🎯 概要

Vercel関連の設定を削除し、GitHub Pagesのみでシステムを公開するように設定しました。

**デプロイ先**: GitHub Pages  
**URL**: `https://Kanta02cer.github.io/Tailorcloud/`

---

## ✅ 実施した変更

### 1. Vercel関連ファイルの削除

以下のファイルを削除しました：

- `.github/workflows/deploy-vercel.yml`
- `.github/workflows/deploy-netlify.yml`
- `.github/workflows/deploy-cloudflare-pages.yml`
- `suit-mbti-web-app/vercel.json`
- `suit-mbti-web-app/.vercelignore`
- `suit-mbti-web-app/README_VERCEL.md`
- `scripts/deploy_vercel_manual.sh`

### 2. コードからVercel関連の参照を削除

- `vite.config.ts`: Vercel関連の条件分岐を削除
- `src/main.tsx`: Vercel関連のbasename設定を削除
- `src/vite-env.d.ts`: VERCEL環境変数の型定義を削除
- `package.json`: vercelパッケージとデプロイスクリプトを削除

### 3. GitHub Pages設定の最適化

- `.github/workflows/deploy-pages.yml`: 環境変数をビルド時に埋め込む設定を追加

---

## 🚀 GitHub Pagesでのデプロイ

### 自動デプロイ（推奨）

1. **GitHubにプッシュ**
   ```bash
   git push
   ```

2. **GitHub Actionsが自動実行**
   - `.github/workflows/deploy-pages.yml` が実行されます
   - 数分でデプロイが完了します

3. **公開URL**
   ```
   https://Kanta02cer.github.io/Tailorcloud/
   ```

### 手動デプロイ

```bash
cd suit-mbti-web-app
npm install
npm run build:pages

# gh-pagesパッケージをインストール（初回のみ）
npm install -g gh-pages

# distディレクトリをgh-pagesブランチにデプロイ
gh-pages -d dist
```

---

## 🔧 環境変数の設定

GitHub Pagesは静的サイトのため、環境変数はビルド時に埋め込まれます。

### GitHub Secretsで設定

1. **GitHubリポジトリ → Settings → Secrets and variables → Actions**

2. **以下のSecretsを追加**:

| Secret名 | 説明 | デフォルト値 |
|---------|------|------------|
| `VITE_API_BASE_URL` | バックエンドAPI URL | `http://localhost:8080` |
| `VITE_TENANT_ID` | テナントID | `tenant_test_suit_mbti` |

### ビルド時に環境変数が埋め込まれる

`.github/workflows/deploy-pages.yml` で環境変数が設定され、ビルド時にコードに埋め込まれます。

---

## 📋 設定ファイル

### vite.config.ts

```typescript
export default defineConfig({
  plugins: [react()],
  // GitHub Pages用のbase path設定
  base: process.env.NODE_ENV === 'production' ? '/Tailorcloud/' : '/',
  // ...
})
```

### src/main.tsx

```typescript
// GitHub Pages用のbase path設定
const basename = import.meta.env.PROD ? '/Tailorcloud' : '';
```

---

## 🎯 アクセス方法

### 公開URL

```
https://Kanta02cer.github.io/Tailorcloud/
```

### ローカル開発

```bash
cd suit-mbti-web-app
npm run dev
```

http://localhost:3000 でアクセス

---

## ⚠️ 制限事項

### GitHub Pagesの制限

1. **静的サイトのみ**
   - サーバーサイドの処理はできません
   - 環境変数はビルド時に埋め込まれます

2. **環境変数の変更**
   - 環境変数を変更した場合は、再ビルド・再デプロイが必要です

3. **API接続**
   - バックエンドAPIは別途デプロイが必要です
   - CORS設定が必要です

---

## 🔄 WordPressとの比較

### WordPressを使用する場合

**メリット**:
- 管理画面が簡単
- プラグインが豊富

**デメリット**:
- Reactアプリをそのまま動かすのは困難
- 既存のReactコードを書き直す必要がある
- サーバーが必要

### GitHub Pagesを使用する場合（現在の設定）

**メリット**:
- 無料
- 既存のReactコードをそのまま使用可能
- 自動デプロイ
- サーバー不要

**デメリット**:
- 静的サイトのみ
- 環境変数の変更には再デプロイが必要

---

## 📚 関連ドキュメント

- **[GitHub Pages セットアップガイド](./GITHUB_PAGES_SETUP.md)**
- **[GitHub Pages デプロイメントガイド](./99_GitHub_Pages_Deployment.md)**

---

## 🎯 次のステップ

1. **GitHub Secretsの設定**（必要に応じて）
   - `VITE_API_BASE_URL`
   - `VITE_TENANT_ID`

2. **GitHubにプッシュ**
   ```bash
   git push
   ```

3. **デプロイの確認**
   - https://Kanta02cer.github.io/Tailorcloud/ にアクセス

---

**最終更新日**: 2025-01

