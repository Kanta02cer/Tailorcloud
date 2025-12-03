# TailorCloud API仕様書（Flutter開発者向け）

**作成日**: 2025-01  
**バージョン**: 1.0.0  
**ベースURL**: `http://localhost:8080` (開発環境)

---

## 📋 目次

1. [認証](#認証)
2. [エンドポイント一覧](#エンドポイント一覧)
3. [注文API](#注文api)
4. [生地API](#生地api)
5. [アンバサダーAPI](#アンバサダーapi)
6. [エラーハンドリング](#エラーハンドリング)

---

## 🔐 認証

### Firebase認証フロー

```dart
// 1. Firebase Authでログイン
final userCredential = await FirebaseAuth.instance.signInWithEmailAndPassword(
  email: email,
  password: password,
);

// 2. IDトークンを取得
final idToken = await userCredential.user?.getIdToken();

// 3. APIリクエスト時にAuthorizationヘッダーに追加
final headers = {
  'Authorization': 'Bearer $idToken',
  'Content-Type': 'application/json',
};
```

### 認証が必要なエンドポイント

- すべての `/api/*` エンドポイント
- `/health` は認証不要

---

## 📡 エンドポイント一覧

| メソッド | エンドポイント | 説明 | 認証 |
|---------|---------------|------|------|
| GET | `/health` | ヘルスチェック | ❌ |
| POST | `/api/orders` | 注文作成 | ✅ |
| POST | `/api/orders/confirm` | 注文確定 | ✅ |
| GET | `/api/orders` | 注文取得・一覧 | ✅ |
| GET | `/api/fabrics` | 生地一覧取得 | ✅ |
| GET | `/api/fabrics/detail` | 生地詳細取得 | ✅ |
| POST | `/api/fabrics/reserve` | 生地確保 | ✅ |
| POST | `/api/ambassadors` | アンバサダー作成 | ✅ |
| GET | `/api/ambassadors/me` | 自分のアンバサダー情報 | ✅ |
| GET | `/api/ambassadors` | アンバサダー一覧 | ✅ |
| GET | `/api/ambassadors/commissions` | 成果報酬一覧 | ✅ |

---

## 📦 注文API

### POST /api/orders - 注文作成

**説明**: 新しい注文を作成（Draftステータス）

**リクエスト**:
```json
{
  "customer_id": "customer-123",
  "fabric_id": "fabric-456",
  "total_amount": 45000,
  "delivery_date": "2025-12-31T00:00:00Z",
  "details": {
    "description": "オーダースーツ縫製",
    "measurement_data": {
      "chest": 100,
      "waist": 85,
      "hip": 95
    },
    "adjustments": {
      "shoulder": "standard",
      "sleeve_length": "custom"
    }
  }
}
```

**レスポンス** (201 Created):
```json
{
  "id": "order-789",
  "tenant_id": "tenant-123",
  "customer_id": "customer-123",
  "fabric_id": "fabric-456",
  "status": "Draft",
  "total_amount": 45000,
  "delivery_date": "2025-12-31T00:00:00Z",
  "payment_due_date": "2025-03-01T00:00:00Z",
  "details": {
    "description": "オーダースーツ縫製"
  },
  "created_at": "2025-01-01T00:00:00Z",
  "created_by": "user-123"
}
```

**エラー**:
- `400 Bad Request`: バリデーションエラー
- `401 Unauthorized`: 認証失敗
- `500 Internal Server Error`: サーバーエラー

---

### POST /api/orders/confirm - 注文確定

**説明**: 注文を確定（Confirmedステータスに変更）

**リクエスト**:
```json
{
  "order_id": "order-789",
  "principal_name": "株式会社テーラー"
}
```

**レスポンス** (200 OK):
```json
{
  "id": "order-789",
  "status": "Confirmed",
  ...
}
```

**エラー**:
- `400 Bad Request`: ステータスがDraftでない
- `401 Unauthorized`: 認証失敗
- `403 Forbidden`: OwnerまたはStaffロールが必要
- `404 Not Found`: 注文が見つからない

---

### GET /api/orders - 注文取得・一覧

**説明**: 注文を取得（単一または一覧）

**クエリパラメータ**:
- `order_id` (オプション): 注文ID（指定した場合は単一取得）
- `tenant_id` (必須): テナントID（認証ユーザーから自動取得される場合は不要）

**レスポンス** (200 OK):

単一取得:
```json
{
  "id": "order-789",
  "tenant_id": "tenant-123",
  ...
}
```

一覧取得:
```json
[
  {
    "id": "order-789",
    "tenant_id": "tenant-123",
    ...
  },
  {
    "id": "order-790",
    ...
  }
]
```

---

## 🧵 生地API

### GET /api/fabrics - 生地一覧取得

**説明**: 生地一覧を取得（フィルター・検索対応）

**クエリパラメータ**:
- `tenant_id` (必須): テナントID
- `status` (オプション): フィルター (`all`, `available`, `limited`, `soldout`)
- `search` (オプション): 検索キーワード（生地名で検索）

**リクエスト例**:
```
GET /api/fabrics?tenant_id=tenant-123&status=available&search=navy
```

**レスポンス** (200 OK):
```json
{
  "fabrics": [
    {
      "id": "fabric-1",
      "name": "Premium Navy Wool",
      "supplier_id": "supplier-1",
      "price": 4500,
      "stock_amount": 5.2,
      "stock_status": "Available",
      "image_url": "https://storage.googleapis.com/.../fabric-1.jpg",
      "minimum_order": 3.2,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  ],
  "total": 1
}
```

**在庫ステータス**:
- `Available`: 在庫あり（> 3.2m）
- `Limited`: 在庫残りわずか（0 < stock_amount ≤ 3.2m）
- `SoldOut`: 在庫切れ（= 0）

---

### GET /api/fabrics/detail - 生地詳細取得

**説明**: 生地詳細を取得

**クエリパラメータ**:
- `fabric_id` (必須): 生地ID
- `tenant_id` (必須): テナントID

**レスポンス** (200 OK):
```json
{
  "id": "fabric-1",
  "name": "Premium Navy Wool",
  "supplier_id": "supplier-1",
  "price": 4500,
  "stock_amount": 5.2,
  "stock_status": "Available",
  "image_url": "https://storage.googleapis.com/.../fabric-1.jpg",
  "minimum_order": 3.2,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

---

### POST /api/fabrics/reserve - 生地確保

**説明**: 生地を確保（発注フロー開始）

**リクエスト**:
```json
{
  "fabric_id": "fabric-1",
  "amount": 3.2
}
```

**レスポンス** (200 OK):
```json
{
  "message": "Fabric reservation successful",
  "fabric_id": "fabric-1",
  "amount": 3.2,
  "status": "reserved"
}
```

**エラー**:
- `400 Bad Request`: 在庫不足、最小発注数量未満

---

## 👤 アンバサダーAPI

### POST /api/ambassadors - アンバサダー作成

**説明**: アンバサダーを作成（Ownerのみ）

**リクエスト**:
```json
{
  "user_id": "firebase-user-123",
  "name": "山田太郎",
  "email": "yamada@example.com",
  "phone": "090-1234-5678",
  "commission_rate": 0.10
}
```

**レスポンス** (201 Created):
```json
{
  "id": "ambassador-1",
  "tenant_id": "tenant-123",
  "user_id": "firebase-user-123",
  "name": "山田太郎",
  "email": "yamada@example.com",
  "status": "Active",
  "commission_rate": 0.10,
  "total_sales": 0,
  "total_commission": 0,
  "created_at": "2025-01-01T00:00:00Z"
}
```

---

### GET /api/ambassadors/me - 自分のアンバサダー情報

**説明**: 認証ユーザーのアンバサダー情報を取得

**レスポンス** (200 OK):
```json
{
  "id": "ambassador-1",
  "name": "山田太郎",
  "total_sales": 450000,
  "total_commission": 45000,
  ...
}
```

---

### GET /api/ambassadors/commissions - 成果報酬一覧

**説明**: 成果報酬一覧を取得

**クエリパラメータ**:
- `ambassador_id` (オプション): アンバサダーID（省略時は自分の成果報酬）
- `limit` (オプション): 取得件数（デフォルト: 20）
- `offset` (オプション): オフセット（デフォルト: 0）

**レスポンス** (200 OK):
```json
{
  "commissions": [
    {
      "id": "commission-1",
      "order_id": "order-789",
      "ambassador_id": "ambassador-1",
      "order_amount": 45000,
      "commission_rate": 0.10,
      "commission_amount": 4500,
      "status": "Approved",
      "paid_at": null,
      "created_at": "2025-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "limit": 20,
  "offset": 0
}
```

**成果報酬ステータス**:
- `Pending`: 未確定（注文が確定していない）
- `Approved`: 確定（支払い待ち）
- `Paid`: 支払済み
- `Cancelled`: キャンセル

---

## ⚠️ エラーハンドリング

### エラーレスポンス形式

```json
{
  "error": "Error message here"
}
```

### HTTPステータスコード

| ステータスコード | 説明 |
|----------------|------|
| `200 OK` | 成功 |
| `201 Created` | 作成成功 |
| `400 Bad Request` | バリデーションエラー |
| `401 Unauthorized` | 認証失敗 |
| `403 Forbidden` | 権限不足 |
| `404 Not Found` | リソースが見つからない |
| `500 Internal Server Error` | サーバーエラー |

### Flutterでのエラーハンドリング例

```dart
try {
  final response = await http.get(
    Uri.parse('$baseUrl/api/orders?order_id=order-123'),
    headers: headers,
  );
  
  if (response.statusCode == 200) {
    final order = jsonDecode(response.body);
    // 処理
  } else if (response.statusCode == 401) {
    // 認証エラー - 再ログイン
    await _refreshAuthToken();
  } else if (response.statusCode == 404) {
    // リソースが見つからない
    showError('注文が見つかりません');
  } else {
    // その他のエラー
    showError('エラーが発生しました: ${response.statusCode}');
  }
} catch (e) {
  // ネットワークエラーなど
  showError('接続エラー: $e');
}
```

---

## 📝 Flutter実装のヒント

### 1. APIクライアントクラス

```dart
class TailorCloudApiClient {
  final String baseUrl;
  final String? idToken;
  
  TailorCloudApiClient({
    required this.baseUrl,
    this.idToken,
  });
  
  Map<String, String> get _headers => {
    'Content-Type': 'application/json',
    if (idToken != null) 'Authorization': 'Bearer $idToken',
  };
  
  Future<Map<String, dynamic>> get(String path) async {
    final response = await http.get(
      Uri.parse('$baseUrl$path'),
      headers: _headers,
    );
    return _handleResponse(response);
  }
  
  Future<Map<String, dynamic>> post(String path, Map<String, dynamic> body) async {
    final response = await http.post(
      Uri.parse('$baseUrl$path'),
      headers: _headers,
      body: jsonEncode(body),
    );
    return _handleResponse(response);
  }
  
  Map<String, dynamic> _handleResponse(http.Response response) {
    if (response.statusCode >= 200 && response.statusCode < 300) {
      return jsonDecode(response.body);
    } else {
      throw ApiException(
        statusCode: response.statusCode,
        message: response.body,
      );
    }
  }
}
```

### 2. モデルクラス

```dart
class Fabric {
  final String id;
  final String name;
  final int price;
  final double stockAmount;
  final String stockStatus;
  final String? imageUrl;
  
  Fabric({
    required this.id,
    required this.name,
    required this.price,
    required this.stockAmount,
    required this.stockStatus,
    this.imageUrl,
  });
  
  factory Fabric.fromJson(Map<String, dynamic> json) {
    return Fabric(
      id: json['id'],
      name: json['name'],
      price: json['price'],
      stockAmount: (json['stock_amount'] as num).toDouble(),
      stockStatus: json['stock_status'],
      imageUrl: json['image_url'],
    );
  }
}
```

---

**最終更新日**: 2025-01  
**バージョン**: 1.0.0

