# TailorCloud Firebase認証統合ガイド

**作成日**: 2025-01  
**実装フェーズ**: Phase 1.1

---

## ✅ 実装完了内容

### 1. Firebase認証ミドルウェア実装

**ファイル**: `internal/middleware/auth.go`

#### 実装内容

- **FirebaseAuthMiddleware** - Firebase認証ミドルウェア
  - JWTトークン検証
  - ユーザー情報のコンテキスト注入
  - OptionalAuth（開発環境用）

- **AuthUser** - 認証済みユーザー情報
  - ID, TenantID, Role, Email

- **コンテキストキー**
  - UserIDKey, TenantIDKey, RoleKey

#### 機能

1. **Authenticate()** - 必須認証ミドルウェア
   - AuthorizationヘッダーからBearerトークンを取得
   - Firebase AuthでJWTトークンを検証
   - カスタムクレームから`tenant_id`と`role`を取得
   - ユーザー情報をコンテキストに注入

2. **OptionalAuth()** - オプショナル認証ミドルウェア
   - トークンがない場合もリクエストを通す（開発環境用）
   - トークンがある場合は検証を試みる

3. **GetUserFromContext()** - コンテキストからユーザー情報を取得

### 2. RBAC（ロールベースアクセス制御）実装

**ファイル**: `internal/middleware/rbac.go`

#### 実装内容

- **RBACMiddleware** - ロールベースアクセス制御ミドルウェア
  - RequireRole() - 特定のロールを要求
  - RequireOwnerOrStaff() - OwnerまたはStaffロールを要求
  - RequireOwnerOnly() - Ownerのみ許可
  - CheckTenantAccess() - テナントアクセスチェック

### 3. HTTPハンドラー統合

**ファイル**: `internal/handler/http_handler.go`

- コンテキストから認証済みユーザー情報を取得
- テナントIDの自動設定（認証ユーザーから）
- ユーザーIDの自動設定（認証ユーザーから）

### 4. main.go統合

**ファイル**: `cmd/api/main.go`

- Firebase Authミドルウェアの初期化
- RBACミドルウェアの初期化
- ルーティングへの認証ミドルウェア適用

---

## 🔐 Firebase Auth設定

### カスタムクレームの設定

Firebase Authで、ユーザートークンに以下のカスタムクレームを設定する必要があります：

- `tenant_id`: ユーザーが所属するテナントID
- `role`: ユーザーのロール（Owner, Staff, Factory_Manager, Worker）

#### 設定方法（Firebase Admin SDK）

```javascript
// Firebase Admin SDKを使用してカスタムクレームを設定
const admin = require('firebase-admin');

async function setCustomClaims(uid, tenantId, role) {
  await admin.auth().setCustomUserClaims(uid, {
    tenant_id: tenantId,
    role: role
  });
}

// 使用例
await setCustomClaims('user-123', 'tenant-456', 'Staff');
```

---

## 📡 API使用例

### 認証なし（開発環境用）

```bash
# OptionalAuthが有効な場合、認証なしでもリクエスト可能
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "tenant-123",
    "customer_id": "customer-456",
    "fabric_id": "fabric-789",
    "total_amount": 45000,
    "delivery_date": "2025-12-31T00:00:00Z",
    "details": {
      "description": "オーダースーツ縫製"
    },
    "created_by": "user-001"
  }'
```

### 認証あり（本番環境）

```bash
# Firebase Authで取得したIDトークンを使用
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <ID_TOKEN>" \
  -d '{
    "customer_id": "customer-456",
    "fabric_id": "fabric-789",
    "total_amount": 45000,
    "delivery_date": "2025-12-31T00:00:00Z",
    "details": {
      "description": "オーダースーツ縫製"
    }
  }'
```

**注意**: 認証ありの場合、`tenant_id`と`created_by`は自動的に設定されます（カスタムクレームから取得）。

---

## 🔄 認証フロー

```
1. クライアント（Flutter App）
   ↓
2. Firebase Authでログイン
   ↓
3. IDトークンを取得
   ↓
4. APIリクエスト時にAuthorizationヘッダーに追加
   Authorization: Bearer <ID_TOKEN>
   ↓
5. FirebaseAuthMiddleware.Authenticate()
   ├─ JWTトークン検証
   ├─ カスタムクレーム取得
   └─ コンテキストにユーザー情報注入
   ↓
6. HTTPハンドラー
   ├─ GetUserFromContext()でユーザー情報取得
   ├─ tenant_id, user_idを自動設定
   └─ ビジネスロジック実行
```

---

## 🛡️ セキュリティ実装

### 実装済み

- ✅ JWTトークン検証
- ✅ カスタムクレームの取得
- ✅ コンテキストへのユーザー情報注入
- ✅ テナントIDの自動設定
- ✅ ロールベースアクセス制御

### 次の実装項目

- [ ] 本番環境での必須認証（OptionalAuthをAuthenticateに切り替え）
- [ ] トークンリフレッシュ機能
- [ ] セッション管理
- [ ] レート制限（認証済みユーザーごと）

---

## 📝 ロール別権限

| ロール | 注文作成 | 注文確定 | 注文閲覧 | 注文一覧 |
|--------|---------|---------|---------|---------|
| **Owner** | ✅ | ✅ | ✅ | ✅ |
| **Staff** | ✅ | ✅ | ✅ | ✅ |
| **Factory_Manager** | ❌ | ✅ | ✅ | ✅ |
| **Worker** | ❌ | ❌ | ✅ | ❌ |

---

## 🔧 開発環境設定

### 環境変数

```bash
# Firebase設定（既存）
GCP_PROJECT_ID=your-gcp-project-id
GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json

# 認証モード（開発環境用）
AUTH_MODE=optional  # optional または required
```

### OptionalAuthの使用

開発環境では、`OptionalAuth`を使用することで、認証なしでもリクエストをテストできます。

```go
// main.goでの設定
if authMiddleware != nil {
    // 開発環境: OptionalAuth
    authHandler = authMiddleware.OptionalAuth
}
```

本番環境では、`Authenticate()`を使用してください：

```go
// 本番環境: 必須認証
authHandler = authMiddleware.Authenticate
```

---

## 🧪 テスト方法

### 1. Firebase Authでユーザー作成

```bash
# Firebase ConsoleまたはAdmin SDKを使用
```

### 2. カスタムクレーム設定

```javascript
await admin.auth().setCustomUserClaims(uid, {
  tenant_id: 'tenant-123',
  role: 'Staff'
});
```

### 3. IDトークン取得

FlutterアプリまたはFirebase Admin SDKでIDトークンを取得

### 4. APIリクエスト

```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Authorization: Bearer <ID_TOKEN>" \
  ...
```

---

## 📊 実装状況

### ✅ 実装完了

- [x] Firebase認証ミドルウェア
- [x] RBACミドルウェア
- [x] HTTPハンドラー統合
- [x] main.go統合

### ⚠️ 実装必要

- [ ] 本番環境での必須認証切り替え
- [ ] トークンリフレッシュ機能
- [ ] エラーハンドリング強化

---

**最終更新日**: 2025-01  
**実装者**: AI Assistant (Auto)

