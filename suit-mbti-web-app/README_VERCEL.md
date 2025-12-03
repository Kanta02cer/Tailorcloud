# Vercelデプロイメント用README

このディレクトリはVercelにデプロイするための設定が含まれています。

## クイックスタート

### 1. Vercelにプロジェクトをインポート

1. https://vercel.com にアクセス
2. "Add New Project" をクリック
3. GitHubリポジトリ `Kanta02cer/Tailorcloud` を選択
4. 以下の設定を入力：
   - **Framework Preset**: Vite
   - **Root Directory**: `suit-mbti-web-app`
   - **Build Command**: `npm run build`
   - **Output Directory**: `dist`

### 2. 環境変数の設定

Vercelのプロジェクト設定で以下の環境変数を設定：

```
VITE_API_BASE_URL=https://your-backend-api.com
VITE_TENANT_ID=tenant_test_suit_mbti
```

### 3. デプロイ

`main` ブランチにプッシュすると自動的にデプロイされます。

## 詳細な手順

詳細は [../docs/107_Vercel_Deployment_Guide.md](../docs/107_Vercel_Deployment_Guide.md) を参照してください。

