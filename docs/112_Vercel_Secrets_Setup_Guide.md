# Vercel GitHub Secrets 設定ガイド

**作成日**: 2025-01  
**目的**: Vercelデプロイに必要なGitHub Secretsの取得方法と設定手順

---

## 📋 必要なGitHub Secrets一覧

| Secret名 | 必須 | 説明 |
|---------|------|------|
| `VERCEL_TOKEN` | ✅ 必須 | Vercel APIトークン |
| `VERCEL_ORG_ID` | ✅ 必須 | Vercel Organization ID |
| `VERCEL_PROJECT_ID` | ✅ 必須 | Vercel Project ID |
| `VITE_API_BASE_URL` | ⚠️ オプション | バックエンドAPI URL（デフォルト: `http://localhost:8080`） |
| `VITE_TENANT_ID` | ⚠️ オプション | テナントID（デフォルト: `tenant_test_suit_mbti`） |

---

## 🔑 1. VERCEL_TOKEN（Vercelトークン）

### 説明
Vercel APIにアクセスするための認証トークンです。GitHub ActionsからVercelにデプロイするために必要です。

### 取得方法

1. **Vercelダッシュボードにログイン**
   - https://vercel.com にアクセス
   - GitHubアカウントでログイン

2. **Settings → Tokens に移動**
   - 右上のプロフィールアイコンをクリック
   - "Settings" を選択
   - 左メニューから "Tokens" を選択

3. **トークンを作成**
   - "Create Token" ボタンをクリック
   - **Token Name**: `github-actions-tailorcloud` など分かりやすい名前を入力
   - **Expiration**: 必要に応じて有効期限を設定（推奨: 無期限）
   - **Scope**: "Full Account" を選択
   - "Create" をクリック

4. **トークンをコピー**
   - ⚠️ **重要**: トークンは一度しか表示されません
   - 表示されたトークンをコピーして安全な場所に保存

### 設定例
```
VERCEL_TOKEN = vercel_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

---

## 🏢 2. VERCEL_ORG_ID（Organization ID）

### 説明
VercelのOrganization（組織）を識別するIDです。個人アカウントの場合もOrganizationとして扱われます。

### 取得方法

1. **Vercelダッシュボードにログイン**
   - https://vercel.com にアクセス

2. **Settings → General に移動**
   - 右上のプロフィールアイコンをクリック
   - "Settings" を選択
   - 左メニューから "General" を選択

3. **Organization IDをコピー**
   - "Organization ID" セクションを探す
   - IDをコピー（例: `team_xxxxxxxxxxxxxxxx`）

### 設定例
```
VERCEL_ORG_ID = team_xxxxxxxxxxxxxxxx
```

**注意**: 個人アカウントの場合、`team_` で始まるIDが表示されます。

---

## 📦 3. VERCEL_PROJECT_ID（Project ID）

### 説明
Vercelプロジェクトを識別するIDです。プロジェクトごとに異なるIDが割り当てられます。

### 取得方法

#### 方法1: Vercelダッシュボードから取得

1. **プロジェクトを選択**
   - Vercelダッシュボードでプロジェクトを選択
   - または、新規プロジェクトを作成

2. **Settings → General に移動**
   - プロジェクトの "Settings" タブをクリック
   - 左メニューから "General" を選択

3. **Project IDをコピー**
   - "Project ID" セクションを探す
   - IDをコピー（例: `prj_xxxxxxxxxxxxxxxx`）

#### 方法2: プロジェクトがまだ作成されていない場合

1. **Vercelでプロジェクトを作成**
   - "Add New Project" をクリック
   - GitHubリポジトリ `Kanta02cer/Tailorcloud` を選択
   - 設定を入力：
     - **Framework Preset**: Vite
     - **Root Directory**: `suit-mbti-web-app`
     - **Build Command**: `npm run build`
     - **Output Directory**: `dist`
   - "Deploy" をクリック

2. **初回デプロイ完了後、Project IDを取得**
   - プロジェクトの Settings → General からProject IDをコピー

### 設定例
```
VERCEL_PROJECT_ID = prj_xxxxxxxxxxxxxxxx
```

---

## 🌐 4. VITE_API_BASE_URL（バックエンドAPI URL）

### 説明
バックエンドAPIのベースURLです。本番環境では実際のAPI URLを設定します。

### 設定値

| 環境 | 値 | 説明 |
|------|-----|------|
| **開発環境** | `http://localhost:8080` | ローカル開発用（デフォルト値） |
| **本番環境** | `https://your-api.com` | 本番環境のAPI URL |

### 設定例

**ローカル開発**:
```
VITE_API_BASE_URL = http://localhost:8080
```

**本番環境（例）**:
```
VITE_API_BASE_URL = https://api.tailorcloud.com
```

**Cloud Runの場合**:
```
VITE_API_BASE_URL = https://tailorcloud-api-xxxxx.run.app
```

**Herokuの場合**:
```
VITE_API_BASE_URL = https://tailorcloud-api.herokuapp.com
```

**Railwayの場合**:
```
VITE_API_BASE_URL = https://tailorcloud-api.railway.app
```

### 注意事項
- ⚠️ バックエンドAPIがまだデプロイされていない場合は、後で設定できます
- デフォルト値（`http://localhost:8080`）が使用されますが、本番環境では動作しません

---

## 🏷️ 5. VITE_TENANT_ID（テナントID）

### 説明
デフォルトのテナントIDです。Suit-MBTI診断機能で使用されます。

### 設定値

| 環境 | 値 | 説明 |
|------|-----|------|
| **開発環境** | `tenant_test_suit_mbti` | テスト用テナントID（デフォルト値） |
| **本番環境** | 実際のテナントID | 本番環境で使用するテナントID |

### 設定例

**開発環境**:
```
VITE_TENANT_ID = tenant_test_suit_mbti
```

**本番環境**:
```
VITE_TENANT_ID = tenant_production_001
```

### 注意事項
- テストデータスクリプト（`scripts/prepare_test_data_suit_mbti.sql`）で `tenant_test_suit_mbti` が使用されています
- 本番環境では、実際のテナントIDを設定してください

---

## 🔧 GitHub Secretsの設定手順

### ステップ1: GitHubリポジトリに移動

1. https://github.com/Kanta02cer/Tailorcloud にアクセス
2. "Settings" タブをクリック

### ステップ2: Secrets and variables → Actions に移動

1. 左メニューから "Secrets and variables" → "Actions" を選択

### ステップ3: 各Secretを追加

1. **"New repository secret" をクリック**
2. **Name** にSecret名を入力（例: `VERCEL_TOKEN`）
3. **Secret** に値を入力（上記で取得した値）
4. **"Add secret" をクリック**

### ステップ4: 全Secretを追加

以下の順序で追加することを推奨：

1. ✅ `VERCEL_TOKEN`（必須）
2. ✅ `VERCEL_ORG_ID`（必須）
3. ✅ `VERCEL_PROJECT_ID`（必須）
4. ⚠️ `VITE_API_BASE_URL`（オプション、後で設定可能）
5. ⚠️ `VITE_TENANT_ID`（オプション、後で設定可能）

---

## ✅ 設定確認

### 1. GitHub Secretsの確認

GitHubリポジトリ → Settings → Secrets and variables → Actions で、以下のSecretが追加されていることを確認：

- ✅ `VERCEL_TOKEN`
- ✅ `VERCEL_ORG_ID`
- ✅ `VERCEL_PROJECT_ID`
- ⚠️ `VITE_API_BASE_URL`（オプション）
- ⚠️ `VITE_TENANT_ID`（オプション）

### 2. ワークフローの実行

1. **GitHubリポジトリ → Actions タブ**
2. **"Deploy to Vercel" ワークフローを選択**
3. **"Run workflow" をクリック**
4. デプロイが成功することを確認

---

## 🔍 トラブルシューティング

### エラー: Invalid token

**原因**: `VERCEL_TOKEN` が間違っている、または期限切れ

**解決方法**:
1. Vercelで新しいトークンを作成
2. GitHub Secretsの `VERCEL_TOKEN` を更新

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

### エラー: API接続エラー

**原因**: `VITE_API_BASE_URL` が正しく設定されていない、またはバックエンドAPIが起動していない

**解決方法**:
1. バックエンドAPIがデプロイされているか確認
2. `VITE_API_BASE_URL` に正しいURLを設定
3. CORS設定を確認

---

## 📚 関連ドキュメント

- **[Vercel GitHub Actions セットアップガイド](./108_Vercel_GitHub_Actions_Setup.md)**
- **[Vercelデプロイメントガイド](./107_Vercel_Deployment_Guide.md)**
- **[環境変数リファレンス](./109_Environment_Variables_Reference.md)**

---

## 🎯 クイックリファレンス

### 必須Secrets（最低限これらを設定）

```
VERCEL_TOKEN = vercel_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
VERCEL_ORG_ID = team_xxxxxxxxxxxxxxxx
VERCEL_PROJECT_ID = prj_xxxxxxxxxxxxxxxx
```

### オプションSecrets（後で設定可能）

```
VITE_API_BASE_URL = https://your-api.com
VITE_TENANT_ID = tenant_test_suit_mbti
```

---

**最終更新日**: 2025-01

