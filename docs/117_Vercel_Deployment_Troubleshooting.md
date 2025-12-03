# Vercelデプロイトラブルシューティングガイド

**作成日**: 2025-01  
**目的**: Vercelデプロイで発生する一般的な問題と解決方法

---

## 🔍 現在の状況確認

### デプロイ履歴の確認

```bash
cd suit-mbti-web-app
vercel ls
```

### ログイン状態の確認

```bash
vercel whoami
```

### プロジェクトのリンク状態確認

```bash
cd suit-mbti-web-app
vercel link
```

---

## 🐛 よくある問題と解決方法

### 問題1: 認証エラー

**エラーメッセージ**:
```
Error: No existing credentials found. Please run `vercel login`
```

**解決方法**:
```bash
vercel login
```

---

### 問題2: プロジェクトが見つからない

**エラーメッセージ**:
```
Error: Project not found
```

**解決方法**:
```bash
cd suit-mbti-web-app
vercel link
```

既存のプロジェクトにリンクするか、新規プロジェクトを作成します。

---

### 問題3: ビルドエラー

**エラーメッセージ**:
```
Error: Build failed
```

**解決方法**:

1. **ローカルでビルドをテスト**
   ```bash
   cd suit-mbti-web-app
   npm run build
   ```

2. **エラーを確認**
   - TypeScriptエラー
   - 依存関係のエラー
   - 環境変数のエラー

3. **修正後、再デプロイ**
   ```bash
   vercel --prod
   ```

---

### 問題4: 環境変数が反映されない

**症状**: デプロイ後、環境変数が `undefined` になる

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

**注意**: 環境変数を追加・変更した後は、必ず再デプロイが必要です。

---

### 問題5: Root Directoryの設定エラー

**エラーメッセージ**:
```
Error: cd: suit-mbti-web-app: No such file or directory
```

**解決方法**:

1. **Vercelダッシュボードで確認**
   - プロジェクトの Settings → General
   - Root Directory が `suit-mbti-web-app` に設定されているか確認

2. **vercel.jsonの確認**
   - `cd suit-mbti-web-app` が含まれていないか確認
   - Root Directoryが設定されている場合、`vercel.json` では相対パスを使用

---

### 問題6: デプロイは成功するが、サイトが表示されない

**症状**: デプロイは成功するが、404エラーが表示される

**解決方法**:

1. **vercel.jsonのrewrites設定を確認**
   ```json
   {
     "rewrites": [
       {
         "source": "/(.*)",
         "destination": "/index.html"
       }
     ]
   }
   ```

2. **React Routerのbasename設定を確認**
   - `src/main.tsx` で `basename` が正しく設定されているか確認

---

### 問題7: CORSエラー

**症状**: APIリクエストでCORSエラーが発生

**解決方法**:

1. **バックエンドAPIのCORS設定を確認**
   - VercelのURLを許可リストに追加

2. **vercel.jsonのheaders設定を確認**
   ```json
   {
     "headers": [
       {
         "source": "/api/(.*)",
         "headers": [
           {
             "key": "Access-Control-Allow-Origin",
             "value": "*"
           }
         ]
       }
     ]
   }
   ```

---

## 🔧 デバッグ手順

### ステップ1: ローカルでビルドをテスト

```bash
cd suit-mbti-web-app
npm install
npm run build
```

### ステップ2: Vercel CLIでビルドをテスト

```bash
cd suit-mbti-web-app
vercel build
```

### ステップ3: プレビューデプロイを実行

```bash
cd suit-mbti-web-app
vercel
```

### ステップ4: 本番デプロイを実行

```bash
cd suit-mbti-web-app
vercel --prod
```

---

## 📊 デプロイログの確認

### Vercelダッシュボードで確認

1. https://vercel.com にログイン
2. プロジェクトを選択
3. "Deployments" タブをクリック
4. デプロイメントを選択
5. "Build Logs" を確認

### Vercel CLIで確認

```bash
vercel logs [deployment-url]
```

---

## ✅ チェックリスト

デプロイ前に以下を確認：

- [ ] ローカルでビルドが成功する
- [ ] Vercelにログインしている（`vercel whoami`）
- [ ] プロジェクトがリンクされている（`vercel link`）
- [ ] 環境変数が設定されている（`vercel env ls`）
- [ ] `vercel.json` の設定が正しい
- [ ] `package.json` のビルドスクリプトが正しい

---

## 🎯 推奨デプロイフロー

### 1. ローカルでテスト

```bash
cd suit-mbti-web-app
npm run build
npm run preview
```

### 2. プレビューデプロイ

```bash
vercel
```

### 3. プレビューを確認

- デプロイされたURLにアクセス
- 動作を確認

### 4. 本番デプロイ

```bash
vercel --prod
```

---

## 📚 関連ドキュメント

- **[Vercel手動デプロイガイド](./116_Vercel_Manual_Deployment.md)**
- **[Vercelビルドエラー修正ガイド](./110_Vercel_Build_Error_Fix.md)**
- **[Vercelデプロイメントガイド](./107_Vercel_Deployment_Guide.md)**

---

**最終更新日**: 2025-01

