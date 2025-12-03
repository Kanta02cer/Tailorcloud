import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../models/order.dart';
import '../services/api_client.dart';
import 'api_client_provider.dart';

part 'order_provider.g.dart';

/// 注文作成プロバイダー
@riverpod
Future<Order> createOrder(
  CreateOrderRef ref,
  CreateOrderRequest request,
) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.post<Map<String, dynamic>>(
    '/api/orders',
    request.toJson(),
  );

  return Order.fromJson(response);
}

/// 注文確定プロバイダー
@riverpod
Future<Order> confirmOrder(
  ConfirmOrderRef ref,
  ConfirmOrderRequest request,
) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.post<Map<String, dynamic>>(
    '/api/orders/confirm',
    request.toJson(),
  );

  return Order.fromJson(response);
}

/// 注文取得プロバイダー
@riverpod
Future<Order> order(
  OrderRef ref,
  String orderId, {
  required String tenantId,
}) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.get<Map<String, dynamic>>(
    '/api/orders',
    queryParameters: {
      'order_id': orderId,
      'tenant_id': tenantId,
    },
  );

  return Order.fromJson(response);
}

/// 注文一覧プロバイダー
@riverpod
Future<List<Order>> orderList(
  OrderListRef ref,
  String tenantId,
) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.get<List<dynamic>>(
    '/api/orders',
    queryParameters: {
      'tenant_id': tenantId,
    },
  );

  return response
      .map((json) => Order.fromJson(json as Map<String, dynamic>))
      .toList();
}

