import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../services/api_client.dart';
import 'api_client_provider.dart';

part 'compliance_provider.g.dart';

/// 発注書生成レスポンス
class ComplianceDocumentResponse {
  final String orderId;
  final String docUrl;
  final String docHash;
  final DateTime generatedAt;

  ComplianceDocumentResponse({
    required this.orderId,
    required this.docUrl,
    required this.docHash,
    required this.generatedAt,
  });

  factory ComplianceDocumentResponse.fromJson(Map<String, dynamic> json) {
    return ComplianceDocumentResponse(
      orderId: json['order_id'] as String,
      docUrl: json['doc_url'] as String,
      docHash: json['doc_hash'] as String,
      generatedAt: DateTime.parse(json['generated_at'] as String),
    );
  }
}

/// 発注書生成プロバイダー
@riverpod
Future<ComplianceDocumentResponse> generateComplianceDocument(
  GenerateComplianceDocumentRef ref,
  String orderId,
) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.post<Map<String, dynamic>>(
    '/api/orders/$orderId/generate-document',
    {},
  );

  return ComplianceDocumentResponse.fromJson(response);
}

