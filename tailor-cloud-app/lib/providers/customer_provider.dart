import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../models/customer.dart';
import '../services/api_client.dart';
import 'api_client_provider.dart';

part 'customer_provider.g.dart';

/// 顧客作成プロバイダー
@riverpod
Future<Customer> createCustomer(
  CreateCustomerRef ref,
  CreateCustomerRequest request,
) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.post<Map<String, dynamic>>(
    '/api/customers',
    request.toJson(),
  );

  return Customer.fromJson(response);
}

/// 顧客取得プロバイダー
@riverpod
Future<Customer> customer(
  CustomerRef ref,
  String customerId,
) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.get<Map<String, dynamic>>(
    '/api/customers/$customerId',
  );

  return Customer.fromJson(response);
}

/// 顧客一覧プロバイダー
@riverpod
Future<List<Customer>> customerList(
  CustomerListRef ref,
  String tenantId,
) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.get<List<dynamic>>(
    '/api/customers',
    queryParameters: {
      'tenant_id': tenantId,
    },
  );

  return response
      .map((json) => Customer.fromJson(json as Map<String, dynamic>))
      .toList();
}

/// 顧客更新プロバイダー
@riverpod
Future<Customer> updateCustomer(
  UpdateCustomerRef ref,
  String customerId,
  CreateCustomerRequest request,
) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.put<Map<String, dynamic>>(
    '/api/customers/$customerId',
    request.toJson(),
  );

  return Customer.fromJson(response);
}

