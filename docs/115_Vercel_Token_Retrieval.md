# Vercelトークンの取得方法

**作成日**: 2025-01  
**目的**: Vercelトークンの取得・再生成手順

---

## 🔑 Vercelトークンの取得手順

### ステップ1: Vercelダッシュボードにログイン

1. **Vercelにアクセス**
   - https://vercel.com にアクセス
   - GitHubアカウントでログイン

### ステップ2: Settings → Tokens に移動

1. **右上のプロフィールアイコンをクリック**
   - アカウントメニューが表示されます

2. **"Settings" を選択**
   - 左メニューから "Settings" をクリック

3. **"Tokens" を選択**
   - 左メニューから "Tokens" をクリック
   - または、直接 https://vercel.com/account/tokens にアクセス

### ステップ3: 既存のトークンを確認

1. **トークン一覧を確認**
   - 既存のトークンがある場合は、名前と作成日が表示されます
   - ⚠️ **重要**: トークンの値は一度しか表示されないため、既存のトークンの値は確認できません

2. **既存のトークンを使用する場合**
   - 既存のトークンがある場合は、そのトークンを使用できます
   - ただし、トークンの値が分からない場合は、新しいトークンを作成する必要があります

### ステップ4: 新しいトークンを作成

1. **"Create Token" ボタンをクリック**

2. **トークン情報を入力**
   - **Token Name**: `github-actions-tailorcloud` など分かりやすい名前を入力
   - **Expiration**: 
     - "No expiration"（無期限）を選択（推奨）
     - または、有効期限を設定

3. **"Create" をクリック**

4. **トークンをコピー**
   - ⚠️ **重要**: トークンは一度しか表示されません
   - 表示されたトークンをコピーして安全な場所に保存
   - 形式: `vercel_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx` または `ecPNqw78Pi96895JcfZLQl0N` のような形式

### ステップ5: トークンをGitHub Secretsに設定

1. **GitHubリポジトリに移動**
   - https://github.com/Kanta02cer/Tailorcloud
   - Settings → Secrets and variables → Actions

2. **VERCEL_TOKENを追加または更新**
   - 既存の `VERCEL_TOKEN` がある場合: "Update" をクリック
   - 新規の場合: "New repository secret" をクリック
   - **Name**: `VERCEL_TOKEN`
   - **Secret**: コピーしたトークンを貼り付け
   - "Add secret" または "Update secret" をクリック

---

## 🔍 以前に提供されたトークン情報

以前の会話で以下のトークンが提供されました：

```
VERCEL_TOKEN = ecPNqw78Pi96895JcfZLQl0N
```

このトークンがまだ有効な場合は、そのまま使用できます。

### トークンの有効性を確認

1. **GitHub Actionsでテスト**
   - GitHubリポジトリ → Actions
   - "Deploy to Vercel" ワークフローを実行
   - エラーが発生しない場合は、トークンは有効です

2. **Vercel CLIでテスト**
   ```bash
   vercel login
   vercel whoami
   ```

---

## 🔄 トークンの再生成

### 既存のトークンが無効になった場合

1. **Vercelダッシュボード → Settings → Tokens**
2. **既存のトークンを削除**（オプション）
   - セキュリティのため、古いトークンを削除することを推奨
3. **新しいトークンを作成**（上記のステップ4を参照）
4. **GitHub Secretsを更新**

---

## ⚠️ セキュリティ注意事項

1. **トークンは機密情報**
   - トークンはGitHub Secretsにのみ保存してください
   - コードやドキュメントに直接記載しないでください

2. **トークンの漏洩を防ぐ**
   - トークンが漏洩した可能性がある場合は、すぐに再生成してください
   - 古いトークンは無効化してください

3. **トークンの有効期限**
   - 無期限のトークンを使用する場合は、定期的に再生成することを推奨

---

## 🧪 トークンの動作確認

### 方法1: GitHub Actionsで確認

1. **GitHubリポジトリ → Actions**
2. **"Deploy to Vercel" ワークフローを実行**
3. **デプロイが成功することを確認**

### 方法2: Vercel CLIで確認

```bash
# Vercel CLIをインストール（まだの場合）
npm install -g vercel

# ログイン（トークンを使用）
vercel login

# または、トークンを直接使用
vercel --token YOUR_TOKEN whoami
```

---

## 📚 関連ドキュメント

- **[Vercel GitHub Secrets 設定ガイド](./112_Vercel_Secrets_Setup_Guide.md)**
- **[Vercel認証情報設定ガイド](./113_Vercel_Credentials_Setup.md)**
- **[Vercel GitHub Actions セットアップガイド](./108_Vercel_GitHub_Actions_Setup.md)**

---

## 🎯 クイックリファレンス

### トークン取得の直接リンク

- **Vercel Tokens**: https://vercel.com/account/tokens
- **Vercel Settings**: https://vercel.com/account

### 必要な情報

1. ✅ **VERCEL_TOKEN**: Vercelトークン
2. ✅ **VERCEL_ORG_ID**: Organization ID（Settings → General）
3. ✅ **VERCEL_PROJECT_ID**: Project ID（プロジェクトの Settings → General）

---

**最終更新日**: 2025-01

