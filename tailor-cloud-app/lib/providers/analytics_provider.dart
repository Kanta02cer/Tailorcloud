import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../models/analytics.dart';
import 'api_client_provider.dart';

part 'analytics_provider.g.dart';

@riverpod
Future<AnalyticsSummary> analyticsSummary(
  AnalyticsSummaryRef ref, {
  String tenantId = 'tenant-123',
  int rangeDays = 30,
}) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.get<Map<String, dynamic>>(
    '/api/analytics/summary',
    queryParameters: {
      'tenant_id': tenantId,
      'range_days': rangeDays.toString(),
    },
  );

  return AnalyticsSummary.fromJson(response);
}

