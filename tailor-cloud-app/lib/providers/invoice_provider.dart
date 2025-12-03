import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../services/api_client.dart';
import 'api_client_provider.dart';

part 'invoice_provider.g.dart';

/// インボイス生成レスポンス
class InvoiceResponse {
  final String orderId;
  final String invoiceUrl;
  final String invoiceHash;
  final DateTime generatedAt;

  InvoiceResponse({
    required this.orderId,
    required this.invoiceUrl,
    required this.invoiceHash,
    required this.generatedAt,
  });

  factory InvoiceResponse.fromJson(Map<String, dynamic> json) {
    return InvoiceResponse(
      orderId: json['order_id'] as String,
      invoiceUrl: json['invoice_url'] as String,
      invoiceHash: json['invoice_hash'] as String,
      generatedAt: DateTime.parse(json['generated_at'] as String),
    );
  }
}

/// インボイス生成プロバイダー
@riverpod
Future<InvoiceResponse> generateInvoice(
  GenerateInvoiceRef ref,
  String orderId,
) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.post<Map<String, dynamic>>(
    '/api/orders/$orderId/generate-invoice',
    {},
  );

  return InvoiceResponse.fromJson(response);
}

