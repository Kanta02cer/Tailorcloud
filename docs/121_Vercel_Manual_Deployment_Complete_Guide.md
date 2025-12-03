# Vercel手動デプロイ完全ガイド

**作成日**: 2025-01  
**目的**: Vercel CLIを使用した手動デプロイの完全な手順

---

## 🎯 概要

このガイドでは、Vercel CLIを使用して手動でデプロイする方法を説明します。Git authorエラーなどの問題も含めて、すべての手順をカバーします。

---

## 📋 前提条件

- Node.js 20以上がインストールされていること
- Vercelアカウントがあること
- プロジェクトがVercelに既に作成されていること

---

## 🚀 方法1: Vercel CLIを使用（Git authorエラー回避版）

### ステップ1: Vercel CLIのインストール

```bash
npm install -g vercel
```

または、プロジェクトローカルにインストール：

```bash
cd suit-mbti-web-app
npm install --save-dev vercel
```

### ステップ2: Vercelにログイン

```bash
cd suit-mbti-web-app
vercel login
```

ブラウザが開き、GitHubアカウントでログインします。

### ステップ3: プロジェクトのリンク

```bash
cd suit-mbti-web-app
vercel link --project tailorcloud
```

### ステップ4: Git authorエラーを回避してデプロイ

Git authorエラーを回避するために、環境変数を設定：

```bash
cd suit-mbti-web-app

# Git authorを一時的に設定
export GIT_AUTHOR_EMAIL=$(vercel whoami | grep -o '[^@]*@[^@]*' || echo "kinouecertify@gmail.com")
export GIT_COMMITTER_EMAIL=$GIT_AUTHOR_EMAIL

# デプロイ
vercel --prod
```

または、より確実な方法：

```bash
cd suit-mbti-web-app

# 最新のコミットを確認
git log --format='%ae' -1

# メールアドレスが正しいことを確認後、デプロイ
vercel --prod
```

---

## 🖥️ 方法2: Vercelダッシュボードから直接デプロイ（最も簡単）

### ステップ1: Vercelダッシュボードにアクセス

1. https://vercel.com にログイン
2. プロジェクト `tailorcloud` を選択

### ステップ2: デプロイメント一覧から再デプロイ

1. **"Deployments" タブをクリック**
2. **最新のデプロイメントを選択**
3. **"..." メニュー → "Redeploy" をクリック**
4. **"Redeploy" を確認**

### ステップ3: 新しいデプロイの確認

- デプロイが開始されます
- 数分で完了します
- デプロイURLが表示されます

---

## 🔧 方法3: GitHubからトリガー（推奨）

### ステップ1: 変更をコミット

```bash
cd /Users/wantan/teiloroud-ERPSystem
git add -A
git commit -m "Your commit message"
```

### ステップ2: GitHubにプッシュ

```bash
git push
```

### ステップ3: GitHub Actionsが自動デプロイ

- `.github/workflows/deploy-vercel.yml` が自動的に実行されます
- Git authorエラーが発生しません

---

## 🛠️ トラブルシューティング

### 問題1: Git authorエラー

**エラーメッセージ**:
```
Error: Git author you@example.com must have access to the team
```

**解決方法**:

#### 方法A: Vercelダッシュボードからデプロイ（推奨）

1. Vercelダッシュボードにアクセス
2. プロジェクトを選択
3. "Redeploy" をクリック

#### 方法B: GitHub Actionsを使用

```bash
git push
```

GitHub Actionsが自動的にデプロイします。

#### 方法C: Git設定を確認

```bash
# Git設定を確認
git config --global user.email
git config --global user.name

# Vercelアカウントのメールアドレスに変更
git config --global user.email "your-vercel-email@example.com"

# 最新のコミットを修正
git commit --amend --reset-author --no-edit
```

---

### 問題2: Root Directoryエラー

**エラーメッセージ**:
```
Error: The provided path "~/teiloroud-ERPSystem/suit-mbti-web-app/suit-mbti-web-app" does not exist
```

**解決方法**:

1. **Vercelダッシュボードにアクセス**
   - https://vercel.com/kinouecertify-gmailcoms-projects/tailorcloud/settings

2. **Root Directoryを空にする**
   - Settings → General
   - "Root Directory" フィールドを空にする
   - Save をクリック

3. **再デプロイ**
   ```bash
   cd suit-mbti-web-app
   vercel --prod
   ```

---

### 問題3: ビルドエラー

**解決方法**:

1. **ローカルでビルドをテスト**
   ```bash
   cd suit-mbti-web-app
   npm install
   npm run build
   ```

2. **エラーを確認して修正**

3. **再デプロイ**
   ```bash
   vercel --prod
   ```

---

### 問題4: 環境変数が反映されない

**解決方法**:

1. **環境変数の確認**
   ```bash
   vercel env ls
   ```

2. **環境変数の追加**
   ```bash
   vercel env add VITE_API_BASE_URL
   # すべての環境（Production, Preview, Development）を選択
   # 値を入力
   ```

3. **再デプロイ**
   ```bash
   vercel --prod
   ```

---

## 📝 デプロイスクリプト

### 簡単なデプロイスクリプト

`scripts/deploy_vercel.sh` を作成：

```bash
#!/bin/bash
# Vercel手動デプロイスクリプト

set -e

echo "=== Vercel手動デプロイ ==="
echo ""

cd "$(dirname "$0")/../suit-mbti-web-app"

# Vercel CLIの確認
if ! command -v vercel &> /dev/null; then
    echo "❌ Vercel CLIがインストールされていません"
    echo "📦 インストール中..."
    npm install -g vercel
fi

# ログイン確認
if ! vercel whoami &> /dev/null; then
    echo "🔐 Vercelにログインしてください"
    vercel login
fi

# デプロイタイプの選択
echo "デプロイタイプを選択してください:"
echo "  1) プレビューデプロイ"
echo "  2) 本番環境デプロイ"
read -p "選択 (1 or 2): " deploy_type

case $deploy_type in
    1)
        echo ""
        echo "🚀 プレビューデプロイを開始..."
        vercel
        ;;
    2)
        echo ""
        echo "🚀 本番環境デプロイを開始..."
        vercel --prod
        ;;
    *)
        echo "❌ 無効な選択です"
        exit 1
        ;;
esac

echo ""
echo "=== デプロイ完了 ==="
```

### 使用方法

```bash
chmod +x scripts/deploy_vercel.sh
./scripts/deploy_vercel.sh
```

---

## 🎯 推奨ワークフロー

### 日常的なデプロイ

1. **変更をコミット**
   ```bash
   git add -A
   git commit -m "Update"
   ```

2. **GitHubにプッシュ**
   ```bash
   git push
   ```

3. **GitHub Actionsが自動デプロイ**
   - エラーが発生しません
   - 自動的にデプロイされます

### 緊急時の手動デプロイ

1. **Vercelダッシュボードから再デプロイ**
   - 最も簡単で確実

2. **または、GitHub Actionsを手動実行**
   - GitHub → Actions → "Deploy to Vercel" → "Run workflow"

---

## 📊 デプロイの確認

### デプロイ一覧の確認

```bash
cd suit-mbti-web-app
vercel ls
```

### デプロイログの確認

```bash
vercel logs [deployment-url]
```

### Vercelダッシュボードで確認

1. https://vercel.com にログイン
2. プロジェクトを選択
3. "Deployments" タブで確認

---

## 🔍 よくある質問

### Q: Git authorエラーを完全に回避する方法は？

**A**: 以下の方法があります：

1. **Vercelダッシュボードから再デプロイ**（最も簡単）
2. **GitHub Actionsを使用**（推奨）
3. **Git設定を正しく設定**（上記参照）

### Q: 手動デプロイと自動デプロイの違いは？

**A**:
- **手動デプロイ**: Vercel CLIまたはダッシュボードから実行
- **自動デプロイ**: GitHubにプッシュすると自動的に実行

### Q: どの方法が最も確実ですか？

**A**: **GitHub Actionsを使用する方法**が最も確実です。Git authorエラーも発生しません。

---

## 📚 関連ドキュメント

- **[Vercel手動デプロイガイド](./116_Vercel_Manual_Deployment.md)**
- **[Vercelデプロイトラブルシューティング](./117_Vercel_Deployment_Troubleshooting.md)**
- **[Vercel GitHub Actions セットアップガイド](./108_Vercel_GitHub_Actions_Setup.md)**

---

## 🎯 クイックリファレンス

### 最も簡単な方法（推奨）

```bash
# 1. GitHubにプッシュ
git push

# 2. GitHub Actionsが自動デプロイ
# または、Vercelダッシュボードから "Redeploy"
```

### Vercel CLIを使用する場合

```bash
cd suit-mbti-web-app
vercel --prod
```

### Vercelダッシュボードから

1. https://vercel.com にアクセス
2. プロジェクトを選択
3. "Redeploy" をクリック

---

**最終更新日**: 2025-01

