# Vercelビルドエラー修正ガイド

**作成日**: 2025-01  
**問題**: `cd: suit-mbti-web-app: No such file or directory`

---

## 🔴 エラー内容

```
sh: line 1: cd: suit-mbti-web-app: No such file or directory
Error: Command "cd suit-mbti-web-app && npm install" exited with 1
```

---

## 🔍 原因

Vercelのプロジェクト設定で **Root Directory** を `suit-mbti-web-app` に設定している場合、Vercelは既にそのディレクトリ内でコマンドを実行します。

そのため、`vercel.json` 内で `cd suit-mbti-web-app` を実行しようとすると、既に `suit-mbti-web-app` ディレクトリ内にいるため、エラーが発生します。

---

## ✅ 解決方法

### 方法1: vercel.jsonの修正（推奨）

`vercel.json` から `cd suit-mbti-web-app` を削除：

**修正前**:
```json
{
  "buildCommand": "cd suit-mbti-web-app && npm run build",
  "outputDirectory": "suit-mbti-web-app/dist",
  "installCommand": "cd suit-mbti-web-app && npm install"
}
```

**修正後**:
```json
{
  "buildCommand": "npm run build",
  "outputDirectory": "dist",
  "installCommand": "npm install"
}
```

### 方法2: Vercelプロジェクト設定の確認

1. **Vercelダッシュボードにログイン**
2. **プロジェクトの Settings → General**
3. **Root Directory** が `suit-mbti-web-app` に設定されているか確認
4. 設定されていない場合は設定

---

## 📋 Vercelプロジェクト設定の確認手順

### ステップ1: プロジェクト設定を確認

1. Vercelダッシュボード → プロジェクトを選択
2. Settings → General
3. 以下の設定を確認：

| 設定項目 | 値 |
|---------|-----|
| **Root Directory** | `suit-mbti-web-app` |
| **Framework Preset** | Vite |
| **Build Command** | `npm run build` |
| **Output Directory** | `dist` |
| **Install Command** | `npm install` |

### ステップ2: Root Directoryの設定

**Root Directory** が空欄または `/` の場合：

1. **Root Directory** に `suit-mbti-web-app` を入力
2. **Save** をクリック
3. 再デプロイを実行

---

## 🔧 vercel.jsonの正しい設定

Root Directoryを `suit-mbti-web-app` に設定している場合の `vercel.json`:

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

**ポイント**:
- `cd suit-mbti-web-app` は不要（Root Directoryで既に設定済み）
- パスは相対パス（Root Directoryからの相対パス）

---

## 🧪 テスト方法

### 1. ローカルでビルドをテスト

```bash
cd suit-mbti-web-app
npm install
npm run build
```

ビルドが成功することを確認。

### 2. Vercel CLIでテスト

```bash
cd suit-mbti-web-app
vercel build
```

エラーが発生しないことを確認。

---

## 📚 関連ドキュメント

- **[Vercelデプロイメントガイド](./107_Vercel_Deployment_Guide.md)**
- **[Vercel GitHub Actions セットアップガイド](./108_Vercel_GitHub_Actions_Setup.md)**

---

**最終更新日**: 2025-01

