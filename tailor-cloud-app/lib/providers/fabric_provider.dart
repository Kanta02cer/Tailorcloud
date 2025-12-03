import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../models/fabric.dart';
import '../services/api_client.dart';
import 'api_client_provider.dart';

part 'fabric_provider.freezed.dart';
part 'fabric_provider.g.dart';

/// 生地一覧取得パラメータ
@freezed
class FabricListParams with _$FabricListParams {
  const factory FabricListParams({
    required String tenantId,
    String? status,
    String? search,
  }) = _FabricListParams;
}

/// 生地一覧プロバイダー
@riverpod
Future<List<Fabric>> fabricList(
  FabricListRef ref,
  FabricListParams params,
) async {
  final apiClient = ref.watch(apiClientProvider);

  final queryParameters = <String, String>{
    'tenant_id': params.tenantId,
    if (params.status != null && params.status != 'all')
      'status': params.status!,
    if (params.search != null && params.search!.isNotEmpty)
      'search': params.search!,
  };

  final response = await apiClient.get<Map<String, dynamic>>(
    '/api/fabrics',
    queryParameters: queryParameters,
  );

  final fabricsList = response['fabrics'] as List;
  return fabricsList
      .map((json) => Fabric.fromJson(json as Map<String, dynamic>))
      .toList();
}

/// 生地詳細プロバイダー
@riverpod
Future<Fabric> fabric(
  FabricRef ref,
  String fabricId, {
  required String tenantId,
}) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.get<Map<String, dynamic>>(
    '/api/fabrics/detail',
    queryParameters: {
      'fabric_id': fabricId,
      'tenant_id': tenantId,
    },
  );

  return Fabric.fromJson(response);
}

/// 生地確保リクエスト
@freezed
class ReserveFabricRequest with _$ReserveFabricRequest {
  const factory ReserveFabricRequest({
    @JsonKey(name: 'fabric_id') required String fabricId,
    required double amount,
  }) = _ReserveFabricRequest;

  factory ReserveFabricRequest.fromJson(Map<String, dynamic> json) =>
      _$ReserveFabricRequestFromJson(json);
}

/// 生地確保プロバイダー
@riverpod
Future<void> reserveFabric(
  ReserveFabricRef ref,
  ReserveFabricRequest request,
) async {
  final apiClient = ref.watch(apiClientProvider);

  await apiClient.post(
    '/api/fabrics/reserve',
    request.toJson(),
  );
}

