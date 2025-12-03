# Vercel GitHub Actions セットアップガイド

**作成日**: 2025-01  
**目的**: GitHub ActionsからVercelに自動デプロイするためのセットアップ手順

---

## 📋 前提条件

1. **Vercelアカウント**: https://vercel.com でアカウント作成済み
2. **GitHubリポジトリ**: https://github.com/Kanta02cer/Tailorcloud.git
3. **Vercelプロジェクト**: Vercelでプロジェクトを作成済み

---

## 🚀 セットアップ手順

### ステップ1: Vercelプロジェクトの作成

1. **Vercelにログイン**
   - https://vercel.com にアクセス
   - GitHubアカウントでログイン

2. **プロジェクトをインポート**
   - "Add New Project" をクリック
   - GitHubリポジトリ `Kanta02cer/Tailorcloud` を選択
   - 設定を入力：
     - **Framework Preset**: Vite
     - **Root Directory**: `suit-mbti-web-app`
     - **Build Command**: `npm run build`
     - **Output Directory**: `dist`
   - "Deploy" をクリック

3. **初回デプロイを完了**
   - 初回デプロイが完了するまで待つ
   - デプロイが成功したことを確認

### ステップ2: VercelトークンとIDの取得

1. **Vercelトークンの取得**
   - Vercelダッシュボード → Settings → Tokens
   - "Create Token" をクリック
   - トークン名を入力（例: `github-actions`）
   - スコープ: Full Account
   - "Create" をクリック
   - **トークンをコピー**（一度しか表示されません）

2. **Organization IDの取得**
   - Vercelダッシュボード → Settings → General
   - "Organization ID" をコピー

3. **Project IDの取得**
   - プロジェクトの Settings → General
   - "Project ID" をコピー

### ステップ3: GitHub Secretsの設定

1. **GitHubリポジトリに移動**
   - https://github.com/Kanta02cer/Tailorcloud
   - Settings → Secrets and variables → Actions

2. **以下のSecretsを追加**

| Secret名 | 説明 | 値 |
|---------|------|-----|
| `VERCEL_TOKEN` | Vercelトークン | ステップ2で取得したトークン |
| `VERCEL_ORG_ID` | Organization ID | ステップ2で取得したOrganization ID |
| `VERCEL_PROJECT_ID` | Project ID | ステップ2で取得したProject ID |
| `VITE_API_BASE_URL` | バックエンドAPI URL | `https://your-api.com` |
| `VITE_TENANT_ID` | デフォルトテナントID | `tenant_test_suit_mbti` |

**設定方法**:
1. "New repository secret" をクリック
2. Nameに上記のSecret名を入力
3. Secretに値を入力
4. "Add secret" をクリック

### ステップ4: 環境変数の設定（Vercel側）

Vercelのプロジェクト設定でも環境変数を設定：

1. **プロジェクトの Settings → Environment Variables**
2. 以下の環境変数を追加：

| 環境変数名 | 値 | 環境 |
|-----------|-----|------|
| `VITE_API_BASE_URL` | `https://your-api.com` | Production, Preview, Development |
| `VITE_TENANT_ID` | `tenant_test_suit_mbti` | Production, Preview, Development |

---

## 🔧 ワークフローの動作

### 自動デプロイの流れ

1. **mainブランチにプッシュ**
   - GitHub Actionsが自動的にトリガーされる
   - `.github/workflows/deploy-vercel.yml` が実行される

2. **ビルドプロセス**
   - Node.js環境のセットアップ
   - 依存関係のインストール
   - 環境変数を使用してビルド

3. **Vercelへのデプロイ**
   - Vercel Actionがビルド成果物をデプロイ
   - 本番環境（Production）にデプロイ

### プルリクエスト時の動作

- プルリクエストが作成されると、プレビューデプロイが自動的に作成されます
- プレビューデプロイのURLがプルリクエストにコメントされます

---

## 🧪 テスト方法

### 1. 手動でワークフローを実行

1. GitHubリポジトリの "Actions" タブに移動
2. "Deploy to Vercel" ワークフローを選択
3. "Run workflow" をクリック
4. ブランチを選択（`main`）
5. "Run workflow" をクリック

### 2. デプロイの確認

1. **GitHub Actions**
   - "Actions" タブでデプロイの進行状況を確認
   - 緑色のチェックマークが表示されれば成功

2. **Vercelダッシュボード**
   - Vercelダッシュボードでデプロイメントを確認
   - デプロイメントのステータス、ログ、URLを確認

---

## 🔧 トラブルシューティング

### エラー: No existing credentials found

**原因**: Vercelトークンが設定されていない

**解決方法**:
1. GitHub Secretsに `VERCEL_TOKEN` が設定されているか確認
2. トークンが正しいか確認（Vercelで再生成）

### エラー: Project not found

**原因**: `VERCEL_PROJECT_ID` が間違っている

**解決方法**:
1. Vercelプロジェクトの Settings → General でProject IDを確認
2. GitHub Secretsの `VERCEL_PROJECT_ID` を更新

### エラー: Organization not found

**原因**: `VERCEL_ORG_ID` が間違っている

**解決方法**:
1. Vercel Settings → General でOrganization IDを確認
2. GitHub Secretsの `VERCEL_ORG_ID` を更新

### ビルドエラー

**原因**: 環境変数が設定されていない、またはビルドコマンドが失敗

**解決方法**:
1. GitHub Actionsのログを確認
2. 環境変数が正しく設定されているか確認
3. ローカルで `npm run build` が成功するか確認

---

## 📚 関連ドキュメント

- **[Vercelデプロイメントガイド](./107_Vercel_Deployment_Guide.md)**
- **[Vercel公式ドキュメント](https://vercel.com/docs)**
- **[GitHub Actions公式ドキュメント](https://docs.github.com/en/actions)**

---

## 🎯 次のステップ

1. **カスタムドメインの設定**（オプション）
   - Vercelプロジェクトの Settings → Domains
   - カスタムドメインを追加

2. **環境変数の追加**
   - 必要に応じて追加の環境変数を設定

3. **監視とアラートの設定**
   - Vercelの監視機能を有効化
   - エラー通知の設定

---

**最終更新日**: 2025-01

