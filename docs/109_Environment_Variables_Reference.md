# 環境変数リファレンス

**作成日**: 2025-01  
**目的**: TailorCloudプロジェクトで使用する環境変数の一覧と設定値

---

## 📋 環境変数一覧

### React Webアプリ（suit-mbti-web-app）

#### VITE_API_BASE_URL

**説明**: バックエンドAPIのベースURL

**設定値**:

| 環境 | 値 | 説明 |
|------|-----|------|
| **開発環境（ローカル）** | `http://localhost:8080` | ローカルでバックエンドを起動している場合 |
| **本番環境（Vercel）** | `https://your-backend-api.com` | 本番環境のバックエンドAPI URL |

**例**:
- ローカル開発: `http://localhost:8080`
- Cloud Run: `https://tailorcloud-api-xxxxx.run.app`
- Heroku: `https://tailorcloud-api.herokuapp.com`
- Railway: `https://tailorcloud-api.railway.app`

**デフォルト値**: `http://localhost:8080`（環境変数が設定されていない場合）

---

#### VITE_TENANT_ID

**説明**: デフォルトテナントID（Suit-MBTI診断用）

**設定値**:

| 環境 | 値 | 説明 |
|------|-----|------|
| **開発環境** | `tenant_test_suit_mbti` | Suit-MBTIテスト用テナントID |
| **本番環境** | 実際のテナントID | 本番環境で使用するテナントID |

**使用されているテナントID**:

1. **`tenant_test_suit_mbti`** (Suit-MBTI Webアプリ用)
   - 診断・予約機能で使用
   - ファイル: `src/pages/DiagnosisPage.tsx`, `src/pages/AppointmentPage.tsx`

2. **`tenant-123`** (Flutterアプリ用)
   - 注文・顧客管理機能で使用
   - テストデータスクリプト: `scripts/prepare_test_data.sql`

**デフォルト値**: `tenant_test_suit_mbti`（環境変数が設定されていない場合）

---

## 🔧 設定方法

### Vercelでの設定

1. **Vercelダッシュボードにログイン**
   - https://vercel.com

2. **プロジェクトを選択**
   - デプロイ済みのプロジェクトを選択

3. **Settings → Environment Variables**
   - 各環境変数を追加
   - 環境を選択（Production, Preview, Development）

4. **環境変数を追加**
   ```
   VITE_API_BASE_URL = https://your-backend-api.com
   VITE_TENANT_ID = tenant_test_suit_mbti
   ```

5. **再デプロイ**
   - 環境変数を追加・変更した後は再デプロイが必要

### GitHub Secretsでの設定（GitHub Actions用）

1. **GitHubリポジトリに移動**
   - https://github.com/Kanta02cer/Tailorcloud

2. **Settings → Secrets and variables → Actions**

3. **New repository secret をクリック**

4. **以下のSecretsを追加**:

| Secret名 | 値 |
|---------|-----|
| `VITE_API_BASE_URL` | `https://your-backend-api.com` |
| `VITE_TENANT_ID` | `tenant_test_suit_mbti` |

### ローカル開発環境での設定

`.env.local` ファイルを作成（`suit-mbti-web-app/` ディレクトリ内）:

```env
VITE_API_BASE_URL=http://localhost:8080
VITE_TENANT_ID=tenant_test_suit_mbti
```

**注意**: `.env.local` ファイルは `.gitignore` に含まれているため、Gitにコミットされません。

---

## 📊 環境別の設定例

### 開発環境（ローカル）

```env
VITE_API_BASE_URL=http://localhost:8080
VITE_TENANT_ID=tenant_test_suit_mbti
```

### ステージング環境

```env
VITE_API_BASE_URL=https://staging-api.tailorcloud.com
VITE_TENANT_ID=tenant_staging
```

### 本番環境

```env
VITE_API_BASE_URL=https://api.tailorcloud.com
VITE_TENANT_ID=tenant_production
```

---

## 🔍 環境変数の確認方法

### ブラウザの開発者ツールで確認

1. ブラウザの開発者ツールを開く（F12）
2. Consoleタブを選択
3. 以下のコマンドを実行:

```javascript
console.log('API Base URL:', import.meta.env.VITE_API_BASE_URL);
console.log('Tenant ID:', import.meta.env.VITE_TENANT_ID);
```

### コード内での確認

```typescript
// src/config/api.ts
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

// src/pages/DiagnosisPage.tsx
const TENANT_ID = import.meta.env.VITE_TENANT_ID || 'tenant_test_suit_mbti';
```

---

## ⚠️ 注意事項

### 1. 環境変数名のプレフィックス

- Viteでは、環境変数名は **`VITE_`** で始まる必要があります
- `VITE_` で始まらない環境変数は、クライアント側のコードからアクセスできません

### 2. ビルド時の環境変数

- 環境変数は**ビルド時**に埋め込まれます
- デプロイ後に環境変数を変更した場合は、**再ビルド・再デプロイ**が必要です

### 3. セキュリティ

- **機密情報（APIキー、パスワード等）は環境変数に含めない**
- バックエンドAPIのURLのみを設定
- 認証トークンはクライアント側で管理（Firebase Auth等）

---

## 📚 関連ドキュメント

- **[Vercelデプロイメントガイド](./107_Vercel_Deployment_Guide.md)**
- **[Vercel GitHub Actions セットアップガイド](./108_Vercel_GitHub_Actions_Setup.md)**
- **[Vite環境変数ドキュメント](https://vitejs.dev/guide/env-and-mode.html)**

---

**最終更新日**: 2025-01

