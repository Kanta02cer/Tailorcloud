# Vercelデプロイメントガイド

**作成日**: 2025-01  
**目的**: TailorCloud WebアプリをVercelにデプロイして動的サイトとして公開

---

## 🎯 概要

TailorCloudのWebアプリ（suit-mbti-web-app）をVercelにデプロイし、動的サイトとして公開します。

**デプロイ先**: Vercel  
**URL**: `https://tailorcloud.vercel.app` (カスタムドメイン設定可能)

---

## 📋 前提条件

1. **Vercelアカウント**: https://vercel.com でアカウント作成
2. **GitHubリポジトリ**: https://github.com/Kanta02cer/Tailorcloud.git
3. **バックエンドAPI**: 別途デプロイ済み（例: Cloud Run, Heroku, Railway等）

---

## 🚀 デプロイ手順

### ステップ1: Vercelプロジェクトの作成

1. **Vercelにログイン**
   - https://vercel.com にアクセス
   - GitHubアカウントでログイン

2. **プロジェクトをインポート**
   - "Add New Project" をクリック
   - GitHubリポジトリ `Kanta02cer/Tailorcloud` を選択
   - "Import" をクリック

3. **プロジェクト設定**
   - **Framework Preset**: Vite
   - **Root Directory**: `suit-mbti-web-app`
   - **Build Command**: `npm run build`
   - **Output Directory**: `dist`
   - **Install Command**: `npm install`

### ステップ2: 環境変数の設定

Vercelのプロジェクト設定で以下の環境変数を設定：

| 環境変数名 | 説明 | 例 |
|-----------|------|-----|
| `VITE_API_BASE_URL` | バックエンドAPIのベースURL | `https://api.tailorcloud.com` |
| `VITE_TENANT_ID` | デフォルトテナントID（オプション） | `tenant_test_suit_mbti` |

**設定方法**:
1. プロジェクトの "Settings" → "Environment Variables" に移動
2. 各環境変数を追加（Production, Preview, Development）
3. "Save" をクリック

### ステップ3: デプロイの実行

1. **自動デプロイ**
   - `main` ブランチにプッシュすると自動的にデプロイされます
   - Vercelが変更を検知して自動ビルド・デプロイを実行

2. **手動デプロイ**
   - Vercelダッシュボードから "Deployments" タブを選択
   - "Redeploy" をクリック

### ステップ4: カスタムドメインの設定（オプション）

1. **ドメインを追加**
   - プロジェクトの "Settings" → "Domains" に移動
   - カスタムドメインを入力（例: `tailorcloud.com`）
   - DNS設定の指示に従う

2. **SSL証明書**
   - Vercelが自動的にSSL証明書を発行・更新

---

## 🔧 設定ファイル

### vercel.json

```json
{
  "version": 2,
  "buildCommand": "cd suit-mbti-web-app && npm run build",
  "outputDirectory": "suit-mbti-web-app/dist",
  "devCommand": "cd suit-mbti-web-app && npm run dev",
  "installCommand": "cd suit-mbti-web-app && npm install",
  "framework": "vite",
  "rewrites": [
    {
      "source": "/(.*)",
      "destination": "/index.html"
    }
  ]
}
```

### vite.config.ts

```typescript
export default defineConfig({
  plugins: [react()],
  // Vercelデプロイ時はbase pathを削除
  base: process.env.VERCEL ? '/' : '/',
  // ...
})
```

---

## 🔌 バックエンドAPI接続

### 環境変数の設定

本番環境では、バックエンドAPIのURLを環境変数で設定：

```env
VITE_API_BASE_URL=https://your-backend-api.com
```

### CORS設定

バックエンドAPIでCORSを許可する必要があります：

```go
// Goバックエンドの例
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "https://tailorcloud.vercel.app")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

---

## 🧪 ローカルでのテスト

### Vercel CLIを使用

```bash
# Vercel CLIをインストール
npm i -g vercel

# プロジェクトディレクトリに移動
cd suit-mbti-web-app

# ログイン
vercel login

# デプロイ（プレビュー）
vercel

# 本番環境にデプロイ
vercel --prod
```

### 環境変数の設定

```bash
# .env.local ファイルを作成
VITE_API_BASE_URL=http://localhost:8080
VITE_TENANT_ID=tenant_test_suit_mbti
```

---

## 📊 デプロイメントの確認

### Vercelダッシュボード

1. **デプロイメント一覧**
   - プロジェクトの "Deployments" タブで確認
   - 各デプロイメントのステータス、ログ、URLを確認

2. **ビルドログ**
   - デプロイメントをクリックしてビルドログを確認
   - エラーがある場合はログで確認

3. **プレビューデプロイメント**
   - プルリクエストごとに自動的にプレビューURLが生成される
   - 本番環境に影響を与えずにテスト可能

---

## 🔧 トラブルシューティング

### ビルドエラー

**問題**: ビルドが失敗する

**解決方法**:
1. ビルドログを確認
2. 依存関係のインストールエラーの場合: `package.json` を確認
3. TypeScriptエラーの場合: `tsconfig.json` を確認

### 環境変数が読み込まれない

**問題**: 環境変数が `undefined` になる

**解決方法**:
1. 環境変数名が `VITE_` で始まっているか確認
2. Vercelの環境変数設定を確認
3. デプロイ後に環境変数を再設定

### API接続エラー

**問題**: バックエンドAPIに接続できない

**解決方法**:
1. `VITE_API_BASE_URL` が正しく設定されているか確認
2. バックエンドAPIが起動しているか確認
3. CORS設定を確認

### ルーティングエラー（404）

**問題**: ページをリロードすると404エラー

**解決方法**:
1. `vercel.json` の `rewrites` 設定を確認
2. React Routerの `basename` 設定を確認

---

## 📚 関連ドキュメント

- **[Vercel公式ドキュメント](https://vercel.com/docs)**
- **[Viteデプロイメントガイド](https://vitejs.dev/guide/static-deploy.html#vercel)**
- **[React Router設定](https://reactrouter.com/en/main/start/overview)**

---

## 🎯 GitHub Pagesとの比較

| 項目 | GitHub Pages | Vercel |
|------|-------------|--------|
| **ホスティングタイプ** | 静的サイト | 動的サイト（SSR対応） |
| **環境変数** | 制限あり | 完全対応 |
| **自動デプロイ** | ✅ | ✅ |
| **プレビューデプロイ** | ❌ | ✅ |
| **カスタムドメイン** | ✅ | ✅ |
| **SSL証明書** | ✅ | ✅（自動更新） |
| **無料プラン** | ✅ | ✅ |

**推奨**: 動的サイトとして動作させる場合は **Vercel** を推奨

---

**最終更新日**: 2025-01

