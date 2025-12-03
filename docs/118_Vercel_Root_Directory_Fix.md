# Vercel Root Directory設定エラー修正

**作成日**: 2025-01  
**問題**: `Error: The provided path "~/teiloroud-ERPSystem/suit-mbti-web-app/suit-mbti-web-app" does not exist`

---

## 🔴 エラー内容

```
Error: The provided path "~/teiloroud-ERPSystem/suit-mbti-web-app/suit-mbti-web-app" does not exist.
To change your Project Settings, go to https://vercel.com/kinouecertify-gmailcoms-projects/tailorcloud/settings
```

---

## 🔍 原因

Vercelのプロジェクト設定で **Root Directory** が `suit-mbti-web-app` に設定されているが、Vercel CLIは既に `suit-mbti-web-app` ディレクトリ内で実行されているため、パスが二重になっています。

---

## ✅ 解決方法

### 方法1: VercelダッシュボードでRoot Directoryを削除（推奨）

1. **Vercelダッシュボードにアクセス**
   - https://vercel.com/kinouecertify-gmailcoms-projects/tailorcloud/settings
   - または、プロジェクト → Settings → General

2. **Root Directoryを空にする**
   - "Root Directory" フィールドを空にする
   - または `/` に設定

3. **Save** をクリック

4. **再デプロイ**
   ```bash
   cd suit-mbti-web-app
   vercel --prod
   ```

### 方法2: プロジェクトルートからデプロイ

プロジェクトルート（`teiloroud-ERPSystem`）からデプロイする場合：

1. **プロジェクトルートに移動**
   ```bash
   cd /Users/wantan/teiloroud-ERPSystem
   ```

2. **Vercelプロジェクトを再リンク**
   ```bash
   vercel link --project tailorcloud
   ```

3. **Root Directoryを指定してデプロイ**
   ```bash
   vercel --prod --cwd suit-mbti-web-app
   ```

---

## 🔧 推奨設定

### 現在のディレクトリ構造

```
teiloroud-ERPSystem/
├── suit-mbti-web-app/    ← ここで vercel コマンドを実行
│   ├── vercel.json
│   ├── package.json
│   └── dist/
```

### Vercelプロジェクト設定

| 設定項目 | 値 |
|---------|-----|
| **Root Directory** | （空）または `/` |
| **Framework Preset** | Vite |
| **Build Command** | `npm run build` |
| **Output Directory** | `dist` |
| **Install Command** | `npm install` |

---

## 📋 確認手順

### 1. 現在の設定を確認

Vercelダッシュボードで確認：
- https://vercel.com/kinouecertify-gmailcoms-projects/tailorcloud/settings

### 2. Root Directoryを確認

- **Root Directory** が `suit-mbti-web-app` になっている場合 → 空にする
- **Root Directory** が空の場合 → そのまま使用

### 3. 再デプロイ

```bash
cd suit-mbti-web-app
vercel --prod
```

---

## 🎯 クイックフィックス

### 最も簡単な解決方法

1. **Vercelダッシュボードにアクセス**
   - https://vercel.com/kinouecertify-gmailcoms-projects/tailorcloud/settings

2. **Root Directoryを削除**
   - "Root Directory" フィールドを空にする

3. **Save** をクリック

4. **再デプロイ**
   ```bash
   cd suit-mbti-web-app
   vercel --prod
   ```

---

## 📚 関連ドキュメント

- **[Vercelデプロイトラブルシューティング](./117_Vercel_Deployment_Troubleshooting.md)**
- **[Vercelビルドエラー修正ガイド](./110_Vercel_Build_Error_Fix.md)**

---

**最終更新日**: 2025-01

