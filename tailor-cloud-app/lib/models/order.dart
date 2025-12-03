import 'package:freezed_annotation/freezed_annotation.dart';

part 'order.freezed.dart';
part 'order.g.dart';

/// 注文ステータス
enum OrderStatus {
  @JsonValue('Draft')
  draft,
  @JsonValue('Confirmed')
  confirmed,
  @JsonValue('InProduction')
  inProduction,
  @JsonValue('Completed')
  completed,
  @JsonValue('Cancelled')
  cancelled,
}

/// 注文詳細情報
@freezed
class OrderDetails with _$OrderDetails {
  const factory OrderDetails({
    String? description,
    @JsonKey(name: 'measurement_data') Map<String, dynamic>? measurementData,
    Map<String, dynamic>? adjustments,
  }) = _OrderDetails;

  factory OrderDetails.fromJson(Map<String, dynamic> json) =>
      _$OrderDetailsFromJson(json);
}

/// 注文モデル
@freezed
class Order with _$Order {
  const factory Order({
    required String id,
    @JsonKey(name: 'tenant_id') required String tenantId,
    @JsonKey(name: 'customer_id') required String customerId,
    @JsonKey(name: 'fabric_id') required String fabricId,
    required OrderStatus status,
    @JsonKey(name: 'total_amount') required int totalAmount,
    @JsonKey(name: 'payment_due_date') required DateTime paymentDueDate,
    @JsonKey(name: 'delivery_date') required DateTime deliveryDate,
    OrderDetails? details,
    @JsonKey(name: 'compliance_doc_url') String? complianceDocUrl,
    @JsonKey(name: 'compliance_doc_hash') String? complianceDocHash,
    @JsonKey(name: 'created_at') required DateTime createdAt,
    @JsonKey(name: 'updated_at') required DateTime updatedAt,
    @JsonKey(name: 'created_by') required String createdBy,
  }) = _Order;

  factory Order.fromJson(Map<String, dynamic> json) => _$OrderFromJson(json);
}

/// 注文拡張メソッド
extension OrderExtension on Order {
  /// ステータスの表示ラベル
  String get statusLabel {
    switch (status) {
      case OrderStatus.draft:
        return '下書き';
      case OrderStatus.confirmed:
        return '確定';
      case OrderStatus.inProduction:
        return '製作中';
      case OrderStatus.completed:
        return '完了';
      case OrderStatus.cancelled:
        return 'キャンセル';
    }
  }

  /// 金額の表示（円）
  String get amountDisplay {
    return '¥${totalAmount.toString().replaceAllMapped(
          RegExp(r'(\d)(?=(\d{3})+(?!\d))'),
          (Match m) => '${m[1]},',
        )}';
  }
}

/// 注文作成リクエスト
@freezed
class CreateOrderRequest with _$CreateOrderRequest {
  const factory CreateOrderRequest({
    @JsonKey(name: 'customer_id') required String customerId,
    @JsonKey(name: 'fabric_id') required String fabricId,
    @JsonKey(name: 'total_amount') required int totalAmount,
    @JsonKey(name: 'delivery_date') required DateTime deliveryDate,
    OrderDetails? details,
  }) = _CreateOrderRequest;

  factory CreateOrderRequest.fromJson(Map<String, dynamic> json) =>
      _$CreateOrderRequestFromJson(json);
}

/// 注文確定リクエスト
@freezed
class ConfirmOrderRequest with _$ConfirmOrderRequest {
  const factory ConfirmOrderRequest({
    @JsonKey(name: 'order_id') required String orderId,
    @JsonKey(name: 'principal_name') required String principalName,
  }) = _ConfirmOrderRequest;

  factory ConfirmOrderRequest.fromJson(Map<String, dynamic> json) =>
      _$ConfirmOrderRequestFromJson(json);
}

