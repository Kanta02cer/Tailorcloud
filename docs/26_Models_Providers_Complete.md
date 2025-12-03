# TailorCloud モデル・プロバイダー実装完了

**作成日**: 2025-01  
**ステータス**: ✅ モデルクラス・プロバイダー実装完了

---

## ✅ 実装完了内容

### モデルクラス

1. **Orderモデル** (`lib/models/order.dart`)
   - Order, OrderDetails
   - CreateOrderRequest, ConfirmOrderRequest
   - OrderStatus enum
   - 拡張メソッド（statusLabel, amountDisplay）

2. **Ambassadorモデル** (`lib/models/ambassador.dart`)
   - Ambassador, Commission
   - AmbassadorStatus, CommissionStatus enum
   - 拡張メソッド（金額表示、ステータスラベル）

3. **Fabricモデル** (`lib/models/fabric.dart`)
   - 既に実装済み

### プロバイダー（Riverpod）

1. **APIクライアントプロバイダー** (`lib/providers/api_client_provider.dart`)
   - ApiClientのシングルトン提供

2. **認証プロバイダー** (`lib/providers/auth_provider.dart`)
   - Firebase Auth状態管理
   - ログイン・ログアウト関数

3. **生地プロバイダー** (`lib/providers/fabric_provider.dart`)
   - 生地一覧取得
   - 生地詳細取得
   - 生地確保

4. **注文プロバイダー** (`lib/providers/order_provider.dart`)
   - 注文作成
   - 注文確定
   - 注文取得・一覧

---

## 📦 実装ファイル一覧

### モデルクラス

- ✅ `lib/models/fabric.dart`
- ✅ `lib/models/order.dart`
- ✅ `lib/models/ambassador.dart`

### プロバイダー

- ✅ `lib/providers/api_client_provider.dart`
- ✅ `lib/providers/auth_provider.dart`
- ✅ `lib/providers/fabric_provider.dart`
- ✅ `lib/providers/order_provider.dart`

---

## 🔄 次のステップ

### コード生成

これらのファイルはFreezedとRiverpod Generatorを使用しているため、コード生成が必要です：

```bash
cd tailor-cloud-app
flutter pub get
flutter pub run build_runner build --delete-conflicting-outputs
```

### 実装チェックリスト

- [x] モデルクラス実装（Fabric, Order, Ambassador）
- [x] プロバイダー実装（Auth, Fabric, Order）
- [ ] 画面実装
  - [ ] Home画面
  - [ ] Inventory画面
  - [ ] Visual Ordering画面

---

## 📝 使用例

### 生地一覧を取得

```dart
final params = FabricListParams(
  tenantId: 'tenant-123',
  status: 'available',
);
final fabrics = await ref.read(fabricListProvider(params).future);
```

### 注文を作成

```dart
final request = CreateOrderRequest(
  customerId: 'customer-123',
  fabricId: 'fabric-456',
  totalAmount: 45000,
  deliveryDate: DateTime.now().add(Duration(days: 30)),
);
final order = await ref.read(createOrderProvider(request).future);
```

---

**最終更新日**: 2025-01

