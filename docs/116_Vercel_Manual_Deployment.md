# Vercel手動デプロイガイド

**作成日**: 2025-01  
**目的**: Vercel CLIを使用した手動デプロイの手順

---

## 🚀 クイックスタート

### 前提条件

- Node.js 20以上がインストールされていること
- Vercelアカウントがあること
- プロジェクトがVercelに既に作成されていること

---

## 📋 手動デプロイ手順

### ステップ1: Vercel CLIのインストール

```bash
npm install -g vercel
```

または、プロジェクトローカルにインストール：

```bash
cd suit-mbti-web-app
npm install --save-dev vercel
```

---

### ステップ2: Vercelにログイン

```bash
cd suit-mbti-web-app
vercel login
```

ブラウザが開き、GitHubアカウントでログインします。

---

### ステップ3: プロジェクトのリンク（初回のみ）

既存のVercelプロジェクトにリンクする場合：

```bash
cd suit-mbti-web-app
vercel link
```

以下の情報を入力：
- **Set up and deploy?** → `Y`
- **Which scope?** → アカウントを選択
- **Link to existing project?** → `Y`
- **What's the name of your existing project?** → プロジェクト名を入力（例: `tailorcloud`）

または、新規プロジェクトを作成する場合：

```bash
cd suit-mbti-web-app
vercel link
```

- **Set up and deploy?** → `Y`
- **Which scope?** → アカウントを選択
- **Link to existing project?** → `N`
- **What's your project's name?** → プロジェクト名を入力

---

### ステップ4: 環境変数の設定

#### 方法1: Vercel CLIで設定

```bash
cd suit-mbti-web-app
vercel env add VITE_API_BASE_URL
# 値を入力（例: https://api.tailorcloud.com）

vercel env add VITE_TENANT_ID
# 値を入力（例: tenant_test_suit_mbti）
```

#### 方法2: Vercelダッシュボードで設定

1. https://vercel.com にアクセス
2. プロジェクトを選択
3. Settings → Environment Variables
4. 環境変数を追加：
   - `VITE_API_BASE_URL`
   - `VITE_TENANT_ID`

---

### ステップ5: ビルドとデプロイ

#### プレビューデプロイ（開発環境）

```bash
cd suit-mbti-web-app
vercel
```

#### 本番環境へのデプロイ

```bash
cd suit-mbti-web-app
vercel --prod
```

---

## 🔧 詳細なデプロイオプション

### 環境変数を指定してデプロイ

```bash
cd suit-mbti-web-app
vercel --prod \
  --env VITE_API_BASE_URL=https://api.tailorcloud.com \
  --env VITE_TENANT_ID=tenant_test_suit_mbti
```

### ビルドコマンドを指定

```bash
cd suit-mbti-web-app
vercel --prod --build-env NODE_ENV=production
```

### デプロイメッセージを追加

```bash
cd suit-mbti-web-app
vercel --prod --message "Manual deployment from local machine"
```

---

## 📝 デプロイスクリプトの作成

### package.jsonにスクリプトを追加

```json
{
  "scripts": {
    "deploy": "vercel --prod",
    "deploy:preview": "vercel"
  }
}
```

### 使用方法

```bash
cd suit-mbti-web-app
npm run deploy        # 本番環境にデプロイ
npm run deploy:preview # プレビュー環境にデプロイ
```

---

## 🔍 デプロイの確認

### デプロイ一覧の確認

```bash
cd suit-mbti-web-app
vercel ls
```

### デプロイログの確認

```bash
cd suit-mbti-web-app
vercel logs [deployment-url]
```

### デプロイ情報の確認

```bash
cd suit-mbti-web-app
vercel inspect [deployment-url]
```

---

## 🧪 ローカルでのテスト

### ローカルでビルドをテスト

```bash
cd suit-mbti-web-app
npm run build
```

### ローカルでプレビュー

```bash
cd suit-mbti-web-app
vercel dev
```

ローカルサーバーが起動し、Vercelの環境でテストできます。

---

## 🔄 デプロイのロールバック

### 前のデプロイにロールバック

```bash
cd suit-mbti-web-app
vercel rollback [deployment-url]
```

### デプロイ履歴の確認

```bash
cd suit-mbti-web-app
vercel ls
```

---

## ⚙️ 設定ファイル

### vercel.json

プロジェクトルート（`suit-mbti-web-app/`）に `vercel.json` が存在する場合、その設定が使用されます。

現在の設定：
```json
{
  "version": 2,
  "buildCommand": "npm run build",
  "outputDirectory": "dist",
  "devCommand": "npm run dev",
  "installCommand": "npm install",
  "framework": "vite",
  "rewrites": [
    {
      "source": "/(.*)",
      "destination": "/index.html"
    }
  ]
}
```

---

## 🐛 トラブルシューティング

### エラー: Not authenticated

```bash
vercel login
```

### エラー: Project not found

```bash
vercel link
```

プロジェクトを再リンクしてください。

### エラー: Build failed

1. ローカルでビルドをテスト：
   ```bash
   cd suit-mbti-web-app
   npm run build
   ```

2. エラーを確認して修正

### 環境変数が反映されない

1. Vercelダッシュボードで環境変数を確認
2. 再デプロイ：
   ```bash
   vercel --prod
   ```

---

## 📚 関連コマンド

### プロジェクト情報の確認

```bash
vercel whoami          # 現在のユーザー情報
vercel projects ls     # プロジェクト一覧
vercel domains ls      # ドメイン一覧
```

### 環境変数の管理

```bash
vercel env ls          # 環境変数一覧
vercel env add        # 環境変数を追加
vercel env rm         # 環境変数を削除
vercel env pull       # 環境変数をローカルにダウンロード
```

---

## 🎯 推奨ワークフロー

### 開発フロー

1. **ローカルで開発**
   ```bash
   cd suit-mbti-web-app
   npm run dev
   ```

2. **プレビューデプロイ**
   ```bash
   vercel
   ```

3. **本番デプロイ**
   ```bash
   vercel --prod
   ```

### CI/CDとの併用

- 通常のデプロイ: GitHub Actions（自動）
- 緊急時のデプロイ: Vercel CLI（手動）

---

## 📚 関連ドキュメント

- **[Vercelデプロイメントガイド](./107_Vercel_Deployment_Guide.md)**
- **[Vercel GitHub Actions セットアップガイド](./108_Vercel_GitHub_Actions_Setup.md)**
- **[Vercel公式CLIドキュメント](https://vercel.com/docs/cli)**

---

**最終更新日**: 2025-01

