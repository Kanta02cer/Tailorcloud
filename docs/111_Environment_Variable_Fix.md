# 環境変数設定エラー修正ガイド

**作成日**: 2025-01  
**問題**: VITE_API_BASE_URLの重複定義によるエラー

---

## 🔴 問題の原因

`VITE_API_BASE_URL` が複数のファイルで重複定義されていました：

1. `src/api/client.ts` - 直接定義
2. `src/config/api.ts` - 一元管理用の定義

この重複により、以下の問題が発生していました：
- 型定義の不整合
- ビルド時の警告
- 環境変数の参照が統一されていない

---

## ✅ 修正内容

### 1. 一元管理への統一

**修正前** (`src/api/client.ts`):
```typescript
// 環境変数からAPIベースURLを取得
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
```

**修正後** (`src/api/client.ts`):
```typescript
import { API_BASE_URL } from '../config/api';
```

### 2. 一元管理ファイル

`src/config/api.ts` で一元管理：

```typescript
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
```

---

## 📋 環境変数の使用箇所

### 正しい使用方法

1. **API設定ファイルからインポート**（推奨）
   ```typescript
   import { API_BASE_URL } from '../config/api';
   ```

2. **直接参照**（必要な場合のみ）
   ```typescript
   const apiUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
   ```

### 使用箇所一覧

| ファイル | 使用方法 | 状態 |
|---------|---------|------|
| `src/config/api.ts` | 定義・エクスポート | ✅ 一元管理 |
| `src/api/client.ts` | インポート | ✅ 修正済み |
| `vite.config.ts` | ビルド時設定 | ✅ 問題なし |

---

## 🔧 環境変数の設定

### Vercelでの設定

1. **Vercelダッシュボード → Settings → Environment Variables**
2. 以下の環境変数を追加：

```
VITE_API_BASE_URL = https://your-backend-api.com
VITE_TENANT_ID = tenant_test_suit_mbti
```

### GitHub Secretsでの設定

1. **GitHubリポジトリ → Settings → Secrets and variables → Actions**
2. 以下のSecretsを追加：

```
VITE_API_BASE_URL = https://your-backend-api.com
VITE_TENANT_ID = tenant_test_suit_mbti
```

### ローカル開発環境での設定

`.env.local` ファイルを作成（`suit-mbti-web-app/` ディレクトリ内）：

```env
VITE_API_BASE_URL=http://localhost:8080
VITE_TENANT_ID=tenant_test_suit_mbti
```

---

## 🧪 動作確認

### 1. ビルドの確認

```bash
cd suit-mbti-web-app
npm run build
```

エラーが発生しないことを確認。

### 2. TypeScriptの型チェック

```bash
cd suit-mbti-web-app
npx tsc --noEmit
```

型エラーが発生しないことを確認。

### 3. 環境変数の確認

ブラウザの開発者ツール（F12）で確認：

```javascript
console.log('API Base URL:', import.meta.env.VITE_API_BASE_URL);
```

---

## 📚 関連ドキュメント

- **[環境変数リファレンス](./109_Environment_Variables_Reference.md)**
- **[Vercelデプロイメントガイド](./107_Vercel_Deployment_Guide.md)**

---

**最終更新日**: 2025-01

