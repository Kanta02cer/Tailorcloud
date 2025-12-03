# Suit-MBTI Web Application

TailorCloudバックエンドと統合されたSuit-MBTI診断ツール・管理画面

## 📋 概要

このアプリケーションは、以下の機能を提供します：

- **診断ツール**: Suit-MBTI診断の実行と結果表示
- **管理画面**: 顧客管理、予約管理、診断履歴管理
- **CRM機能**: 顧客情報の一元管理

## 🏗️ 技術スタック

- **Framework**: React 18+ with TypeScript
- **State Management**: React Query (TanStack Query) または Zustand
- **UI Library**: Material-UI (MUI) または Tailwind CSS
- **HTTP Client**: Axios
- **Routing**: React Router v6
- **Build Tool**: Vite

## 📁 プロジェクト構造

```
suit-mbti-web-app/
├── src/
│   ├── api/              # APIクライアント
│   │   ├── client.ts     # Axios設定
│   │   ├── diagnoses.ts  # 診断API
│   │   └── appointments.ts # 予約API
│   ├── components/       # 再利用可能なコンポーネント
│   │   ├── Diagnosis/    # 診断関連コンポーネント
│   │   ├── Appointment/  # 予約関連コンポーネント
│   │   └── Customer/     # 顧客関連コンポーネント
│   ├── pages/            # ページコンポーネント
│   │   ├── DiagnosisPage.tsx
│   │   ├── AppointmentPage.tsx
│   │   └── CustomerPage.tsx
│   ├── hooks/            # カスタムフック
│   ├── types/            # TypeScript型定義
│   ├── utils/            # ユーティリティ関数
│   └── App.tsx           # メインアプリコンポーネント
├── package.json
├── tsconfig.json
└── vite.config.ts
```

## 🚀 セットアップ

### 1. 依存関係のインストール

```bash
npm install
# または
yarn install
# または
pnpm install
```

### 2. 環境変数の設定

`.env.local`ファイルを作成：

```env
VITE_API_BASE_URL=http://localhost:8080
VITE_FIREBASE_API_KEY=your_firebase_api_key
VITE_FIREBASE_AUTH_DOMAIN=your_firebase_auth_domain
VITE_FIREBASE_PROJECT_ID=your_firebase_project_id
```

### 3. 開発サーバーの起動

```bash
npm run dev
# または
yarn dev
# または
pnpm dev
```

アプリケーションは `http://localhost:3000` で起動します。

## 🔗 TailorCloudバックエンドとの連携

### APIエンドポイント

- **診断API**: `/api/diagnoses`
- **予約API**: `/api/appointments`
- **顧客API**: `/api/customers`

詳細は `docs/78_Suit_MBTI_Feature_Guide.md` を参照してください。

## 📚 関連ドキュメント

- [機能ガイド](../docs/78_Suit_MBTI_Feature_Guide.md)
- [手動テストガイド](../docs/79_Manual_Testing_Guide.md)
- [統合マスタープラン](../docs/75_Suit_MBTI_Integration_Master_Plan.md)

