import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../services/api_client.dart';
import '../config/app_config.dart';

part 'api_client_provider.g.dart';

/// APIクライアントプロバイダー
@riverpod
ApiClient apiClient(ApiClientRef ref) {
  return ApiClient(baseUrl: AppConfig.baseUrl);
}
