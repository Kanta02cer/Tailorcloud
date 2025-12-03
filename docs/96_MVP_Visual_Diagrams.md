# TailorCloud MVP 視覚化図表

**作成日**: 2025-01  
**バージョン**: 1.0.0  
**補足資料**: `95_MVP_Technical_Specification.md` の視覚化版

---

## 1. ER図（Mermaid形式）

### 1.1 完全ER図

```mermaid
erDiagram
    tenants ||--o{ users : "has"
    tenants ||--o{ customers : "has"
    tenants ||--o{ orders : "has"
    tenants ||--o{ fabrics : "has"
    
    users ||--o{ orders : "creates"
    users ||--o{ compliance_documents : "generates"
    
    customers ||--o{ orders : "places"
    
    orders ||--o{ order_items : "contains"
    orders ||--o{ compliance_documents : "has"
    
    fabrics ||--o{ order_items : "used_in"
    
    tenants {
        uuid id PK
        varchar name
        varchar legal_name
        varchar address
        varchar invoice_registration_no
        varchar type
        boolean is_active
        timestamptz created_at
        timestamptz updated_at
    }
    
    users {
        uuid id PK
        uuid tenant_id FK
        varchar email
        varchar firebase_uid
        varchar name
        varchar role
        boolean is_active
        timestamptz last_login_at
        timestamptz created_at
        timestamptz updated_at
    }
    
    customers {
        uuid id PK
        uuid tenant_id FK
        varchar name
        varchar email
        varchar phone
        text address
        timestamptz created_at
        timestamptz updated_at
    }
    
    orders {
        uuid id PK
        uuid tenant_id FK
        varchar order_number
        uuid customer_id FK
        uuid created_by FK
        varchar status
        decimal total_amount
        decimal tax_amount
        decimal tax_rate
        date delivery_due_date
        date payment_due_date
        text compliance_doc_url
        varchar compliance_doc_hash
        timestamptz invoice_issued_at
        timestamptz created_at
        timestamptz updated_at
    }
    
    order_items {
        uuid id PK
        uuid order_id FK
        uuid fabric_id FK
        varchar item_type
        jsonb measurements
        jsonb options
        decimal required_fabric_length
        decimal unit_price
        integer quantity
        timestamptz created_at
        timestamptz updated_at
    }
    
    fabrics {
        uuid id PK
        uuid tenant_id FK
        uuid supplier_id FK
        varchar sku
        varchar brand
        varchar name
        varchar composition
        decimal cost_price
        decimal sales_price
        text image_url
        boolean is_active
        timestamptz created_at
        timestamptz updated_at
    }
    
    compliance_documents {
        uuid id PK
        uuid order_id FK
        uuid parent_document_id FK
        uuid generated_by FK
        varchar document_type
        text pdf_url
        varchar pdf_hash
        text amendment_reason
        timestamptz generated_at
    }
```

### 1.2 データフロー図

```mermaid
flowchart TD
    A[顧客] -->|注文| B[注文作成]
    B --> C[注文確定]
    C --> D[発注書PDF生成]
    D --> E[Cloud Storage保存]
    E --> F[compliance_documentsテーブル記録]
    
    G[テーラー] -->|発注書確認| H[PDFダウンロード]
    H --> I[工場へ送付]
    
    J[修正が必要] -->|修正発注書生成| K[amendment_document作成]
    K --> L[親文書との紐付け]
    L --> M[PDF生成・保存]
```

---

## 2. 画面遷移図（Mermaid形式）

### 2.1 全体フロー

```mermaid
flowchart TD
    Start[ログイン画面] -->|Firebase認証| Home[ホーム画面]
    
    Home -->|顧客タブ| CustomerList[顧客一覧]
    Home -->|生地タブ| FabricList[生地一覧]
    Home -->|注文タブ| OrderList[注文一覧]
    Home -->|クイック発注| QuickOrder[クイック発注]
    
    CustomerList -->|顧客選択| CustomerDetail[顧客詳細]
    CustomerList -->|新規登録| CustomerCreate[新規顧客登録]
    CustomerDetail -->|編集| CustomerEdit[顧客編集]
    CustomerDetail -->|注文履歴| OrderList
    
    FabricList -->|生地選択| FabricDetail[生地詳細]
    
    OrderList -->|注文選択| OrderDetail[注文詳細]
    OrderDetail -->|発注書表示| ComplianceDoc[発注書PDF表示]
    OrderDetail -->|確定| OrderConfirm[注文確定確認]
    OrderDetail -->|修正発注書| AmendmentDoc[修正発注書生成]
    
    QuickOrder -->|ステップ1| Step1[顧客選択]
    Step1 -->|既存顧客| Step2[生地選択]
    Step1 -->|新規登録| CustomerCreate
    CustomerCreate -->|登録完了| Step2
    
    Step2 -->|生地選択| Step3[注文情報入力]
    Step3 -->|保存| OrderDraft[下書き保存]
    Step3 -->|確定| OrderConfirm
    OrderConfirm -->|発注書生成| ComplianceDoc
    
    OrderDraft -->|編集| Step3
    OrderDraft -->|確定| OrderConfirm
```

### 2.2 クイック発注フロー（詳細）

```mermaid
flowchart LR
    A[ステップ1: 顧客選択] -->|既存顧客選択| B[ステップ2: 生地選択]
    A -->|新規顧客登録| C[新規顧客登録画面]
    C -->|登録完了| B
    B -->|生地選択| D[ステップ3: 注文情報入力]
    D -->|保存| E[下書き保存]
    D -->|確定| F[注文確定]
    F -->|発注書生成| G[発注書PDF表示]
    E -->|編集| D
```

### 2.3 注文ステータス遷移

```mermaid
stateDiagram-v2
    [*] --> Draft: 注文作成
    Draft --> Confirmed: 注文確定（発注書生成）
    Confirmed --> Production: 工場受注
    Production --> Shipped: 発送
    Shipped --> Delivered: 納品
    Delivered --> Paid: 支払い完了
    Draft --> Cancelled: キャンセル
    Confirmed --> Cancelled: キャンセル
    Cancelled --> [*]
    Paid --> [*]
```

---

## 3. APIフロー図

### 3.1 注文作成フロー

```mermaid
sequenceDiagram
    participant Client as Flutter App
    participant API as Backend API
    participant Auth as Firebase Auth
    participant DB as PostgreSQL
    participant Storage as Cloud Storage
    
    Client->>Auth: ログイン（Email/Password）
    Auth-->>Client: JWT Token
    
    Client->>API: POST /api/orders (JWT)
    API->>Auth: JWT検証
    Auth-->>API: ユーザー情報
    
    API->>DB: 顧客情報取得
    DB-->>API: 顧客データ
    
    API->>DB: 生地情報取得
    DB-->>API: 生地データ
    
    API->>DB: 注文作成（Draft）
    DB-->>API: 注文ID
    
    API-->>Client: 注文データ（Draft）
    
    Client->>API: POST /api/orders/{id}/confirm
    API->>DB: 注文ステータス更新（Confirmed）
    API->>Storage: PDF生成・アップロード
    Storage-->>API: PDF URL
    API->>DB: compliance_documents記録
    API-->>Client: 発注書URL
```

### 3.2 認証フロー

```mermaid
sequenceDiagram
    participant Client as Flutter App
    participant Auth as Firebase Auth
    participant API as Backend API
    participant DB as PostgreSQL
    
    Client->>Auth: ログイン（Email/Password）
    Auth-->>Client: JWT Token + Refresh Token
    
    Client->>API: APIリクエスト（JWT）
    API->>Auth: JWT検証
    Auth-->>API: ユーザー情報（UID）
    
    API->>DB: ユーザー情報取得（firebase_uid）
    DB-->>API: ユーザーデータ（tenant_id, role）
    
    API->>API: 権限チェック（RBAC）
    API-->>Client: レスポンス
```

---

## 4. インフラ構成図（Mermaid形式）

### 4.1 全体構成

```mermaid
graph TB
    subgraph Client["クライアント層"]
        Flutter[Flutter App<br/>iOS/Android]
        Web[Web Portal<br/>React - 将来実装]
    end
    
    subgraph API["API層"]
        CloudRun[Cloud Run<br/>Go Backend API<br/>Auto Scaling]
    end
    
    subgraph Auth["認証層"]
        Firebase[Firebase Authentication<br/>JWT Token]
    end
    
    subgraph Data["データ層"]
        CloudSQL[Cloud SQL<br/>PostgreSQL<br/>Primary Database]
        CloudStorage[Cloud Storage<br/>PDF Documents<br/>Images]
    end
    
    subgraph Monitoring["監視層"]
        Logging[Cloud Logging<br/>構造化ログ]
        Monitoring[Cloud Monitoring<br/>メトリクス・アラート]
    end
    
    Flutter -->|HTTPS + JWT| CloudRun
    Web -->|HTTPS + JWT| CloudRun
    
    CloudRun -->|認証検証| Firebase
    CloudRun -->|データ取得| CloudSQL
    CloudRun -->|PDF保存| CloudStorage
    
    CloudRun -->|ログ送信| Logging
    CloudRun -->|メトリクス送信| Monitoring
```

### 4.2 データベース接続フロー

```mermaid
sequenceDiagram
    participant App as Cloud Run
    participant Proxy as Cloud SQL Proxy
    participant DB as Cloud SQL PostgreSQL
    
    App->>Proxy: 接続要求（接続名）
    Proxy->>DB: 暗号化接続
    DB-->>Proxy: 接続確立
    Proxy-->>App: 接続完了
    
    App->>Proxy: SQLクエリ
    Proxy->>DB: クエリ実行
    DB-->>Proxy: 結果
    Proxy-->>App: 結果返却
```

### 4.3 PDF生成・保存フロー

```mermaid
sequenceDiagram
    participant API as Backend API
    participant PDF as PDF Generator
    participant Storage as Cloud Storage
    participant DB as PostgreSQL
    
    API->>PDF: 注文データ + テンプレート
    PDF->>PDF: PDF生成（下請法3条書面）
    PDF->>PDF: SHA-256ハッシュ計算
    
    PDF->>Storage: PDFアップロード
    Storage-->>PDF: PDF URL
    
    PDF->>DB: compliance_documents記録
    DB-->>PDF: 記録完了
    
    PDF-->>API: PDF URL + Hash
    API-->>Client: 発注書URL
```

---

## 5. セキュリティフロー

### 5.1 マルチテナントデータ分離

```mermaid
flowchart TD
    A[APIリクエスト] -->|JWT検証| B[ユーザー情報取得]
    B -->|tenant_id抽出| C[クエリにtenant_id追加]
    C -->|WHERE tenant_id = ?| D[データ取得]
    D -->|テナント分離| E[レスポンス返却]
    
    F[不正アクセス試行] -->|異なるtenant_id| G[403 Forbidden]
```

### 5.2 RBAC（ロールベースアクセス制御）

```mermaid
flowchart TD
    A[APIリクエスト] -->|JWT検証| B[ユーザー情報取得]
    B -->|role取得| C{ロールチェック}
    
    C -->|OWNER| D[全権限]
    C -->|STAFF| E[発注・顧客管理]
    C -->|FACTORY_MGR| F[受注管理]
    C -->|WORKER| G[作業完了チェック]
    
    D -->|権限OK| H[処理実行]
    E -->|権限OK| H
    F -->|権限OK| H
    G -->|権限OK| H
    
    C -->|権限不足| I[403 Forbidden]
```

---

## 6. デプロイメントフロー

### 6.1 CI/CDパイプライン

```mermaid
flowchart LR
    A[コードコミット] -->|GitHub| B[GitHub Actions]
    B -->|テスト実行| C{テスト成功?}
    C -->|Yes| D[Dockerイメージビルド]
    C -->|No| E[エラー通知]
    D -->|GCRプッシュ| F[Cloud Runデプロイ]
    F -->|ヘルスチェック| G{デプロイ成功?}
    G -->|Yes| H[本番環境]
    G -->|No| I[ロールバック]
```

### 6.2 バージョン管理

```mermaid
graph LR
    A[開発環境<br/>v1.0.0-dev] -->|テスト| B[ステージング環境<br/>v1.0.0-staging]
    B -->|承認| C[本番環境<br/>v1.0.0]
    C -->|問題発生| D[ロールバック<br/>v0.9.0]
```

---

## 7. 監視・アラートフロー

### 7.1 メトリクス収集

```mermaid
flowchart TD
    A[Cloud Run] -->|メトリクス送信| B[Cloud Monitoring]
    B -->|閾値チェック| C{閾値超過?}
    C -->|Yes| D[アラート発火]
    C -->|No| E[正常]
    D -->|通知| F[Slack/Email]
    D -->|自動対応| G[スケールアウト]
```

### 7.2 ログ分析

```mermaid
flowchart LR
    A[アプリケーション] -->|構造化ログ| B[Cloud Logging]
    B -->|集計| C[ログ分析]
    C -->|エラー検出| D[アラート]
    C -->|トレンド分析| E[ダッシュボード]
```

---

## 8. データフロー（注文確定から発注書生成まで）

### 8.1 詳細フロー

```mermaid
sequenceDiagram
    participant User as ユーザー
    participant App as Flutter App
    participant API as Backend API
    participant Service as Order Service
    participant Compliance as Compliance Service
    participant PDF as PDF Generator
    participant Storage as Cloud Storage
    participant DB as PostgreSQL
    
    User->>App: 注文確定ボタン
    App->>API: POST /api/orders/{id}/confirm
    API->>Service: 注文確定処理
    
    Service->>DB: 注文ステータス更新（DRAFT→CONFIRMED）
    Service->>DB: 注文データ取得
    DB-->>Service: 注文データ
    
    Service->>Compliance: 発注書生成要求
    Compliance->>DB: テナント情報取得（法人名、住所等）
    Compliance->>DB: 顧客情報取得
    Compliance->>DB: 注文明細取得
    
    Compliance->>PDF: PDF生成（下請法3条書面テンプレート）
    PDF->>PDF: SHA-256ハッシュ計算
    PDF-->>Compliance: PDFバイナリ + ハッシュ
    
    Compliance->>Storage: PDFアップロード
    Storage-->>Compliance: PDF URL
    
    Compliance->>DB: compliance_documents記録
    DB-->>Compliance: 記録完了
    
    Compliance->>DB: orders.compliance_doc_url更新
    DB-->>Compliance: 更新完了
    
    Compliance-->>Service: 発注書URL + ハッシュ
    Service-->>API: レスポンス
    API-->>App: 発注書URL
    App->>User: 発注書PDF表示
```

---

## 9. エラーハンドリングフロー

### 9.1 エラー処理の階層

```mermaid
flowchart TD
    A[APIリクエスト] -->|エラー発生| B{エラータイプ}
    
    B -->|認証エラー| C[401 Unauthorized]
    B -->|権限エラー| D[403 Forbidden]
    B -->|リソース不存在| E[404 Not Found]
    B -->|バリデーションエラー| F[400 Bad Request]
    B -->|サーバーエラー| G[500 Internal Server Error]
    
    C -->|ログ記録| H[Cloud Logging]
    D -->|ログ記録| H
    E -->|ログ記録| H
    F -->|ログ記録| H
    G -->|ログ記録| H
    G -->|アラート発火| I[Cloud Monitoring]
    
    H -->|エラー分析| J[ダッシュボード]
    I -->|通知| K[Slack/Email]
```

---

## 10. パフォーマンス最適化フロー

### 10.1 データベースクエリ最適化

```mermaid
flowchart LR
    A[クエリ実行] -->|インデックス使用| B[高速検索]
    A -->|インデックス未使用| C[フルスキャン]
    C -->|パフォーマンス問題| D[インデックス追加]
    D -->|最適化| B
    B -->|結果返却| E[レスポンス]
```

### 10.2 キャッシュ戦略（将来実装）

```mermaid
flowchart TD
    A[APIリクエスト] -->|キャッシュ確認| B{キャッシュヒット?}
    B -->|Yes| C[キャッシュから返却]
    B -->|No| D[データベース取得]
    D -->|キャッシュ保存| E[Redis/Memorystore]
    D -->|レスポンス返却| F[クライアント]
    C --> F
```

---

## まとめ

このドキュメントは、`95_MVP_Technical_Specification.md`の視覚化版です。Mermaid形式の図表により、システムの構造とフローを直感的に理解できます。

**使用方法**:
- GitHubやMarkdownビューアーでMermaid図を表示
- 設計レビュー時の資料として活用
- 開発チームへの説明資料として使用

---

**最終更新日**: 2025-01  
**バージョン**: 1.0.0  
**ステータス**: ✅ 視覚化図表完成

