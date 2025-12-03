import 'dart:convert';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../services/api_client.dart';
import 'api_client_provider.dart';

part 'measurement_validation_provider.g.dart';

/// バリデーションアラート
class ValidationAlert {
  final String field;
  final double? current;
  final double? previous;
  final double? difference;
  final double? threshold;
  final String severity; // "warning" or "error"
  final String message;

  ValidationAlert({
    required this.field,
    this.current,
    this.previous,
    this.difference,
    this.threshold,
    required this.severity,
    required this.message,
  });

  factory ValidationAlert.fromJson(Map<String, dynamic> json) {
    return ValidationAlert(
      field: json['field'] as String,
      current: json['current'] != null ? (json['current'] as num).toDouble() : null,
      previous: json['previous'] != null ? (json['previous'] as num).toDouble() : null,
      difference: json['difference'] != null ? (json['difference'] as num).toDouble() : null,
      threshold: json['threshold'] != null ? (json['threshold'] as num).toDouble() : null,
      severity: json['severity'] as String,
      message: json['message'] as String,
    );
  }
}

/// バリデーションレスポンス
class ValidationResponse {
  final bool isValid;
  final List<ValidationAlert> alerts;
  final bool hasWarnings;
  final bool hasErrors;
  final Map<String, dynamic>? previousData;

  ValidationResponse({
    required this.isValid,
    required this.alerts,
    required this.hasWarnings,
    required this.hasErrors,
    this.previousData,
  });

  factory ValidationResponse.fromJson(Map<String, dynamic> json) {
    return ValidationResponse(
      isValid: json['is_valid'] as bool? ?? true,
      alerts: (json['alerts'] as List<dynamic>?)
          ?.map((e) => ValidationAlert.fromJson(e as Map<String, dynamic>))
          .toList() ?? [],
      hasWarnings: json['has_warnings'] as bool? ?? false,
      hasErrors: json['has_errors'] as bool? ?? false,
      previousData: json['previous_data'] as Map<String, dynamic>?,
    );
  }
}

/// 採寸データバリデーションリクエスト
class ValidateMeasurementsRequest {
  final String customerId;
  final Map<String, dynamic> currentMeasurements;

  ValidateMeasurementsRequest({
    required this.customerId,
    required this.currentMeasurements,
  });

  Map<String, dynamic> toJson() {
    return {
      'customer_id': customerId,
      'current_measurements': currentMeasurements,
    };
  }
}

/// 採寸データバリデーションプロバイダー
@riverpod
Future<ValidationResponse> validateMeasurements(
  ValidateMeasurementsRef ref,
  ValidateMeasurementsRequest request,
) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.post<Map<String, dynamic>>(
    '/api/measurements/validate',
    request.toJson(),
  );

  return ValidationResponse.fromJson(response);
}

/// 採寸データ範囲バリデーションプロバイダー
@riverpod
Future<ValidationResponse> validateMeasurementRange(
  ValidateMeasurementRangeRef ref,
  Map<String, dynamic> measurements,
) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.post<Map<String, dynamic>>(
    '/api/measurements/validate-range',
    {'measurements': measurements},
  );

  return ValidationResponse.fromJson(response);
}

