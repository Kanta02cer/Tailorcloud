import 'package:freezed_annotation/freezed_annotation.dart';

part 'analytics.freezed.dart';
part 'analytics.g.dart';

@freezed
class AnalyticsSummary with _$AnalyticsSummary {
  const factory AnalyticsSummary({
    @JsonKey(name: 'tenant_id') required String tenantId,
    @JsonKey(name: 'range_days') required int rangeDays,
    @JsonKey(name: 'generated_at') required DateTime generatedAt,
    @JsonKey(name: 'total_orders') required int totalOrders,
    @JsonKey(name: 'total_revenue') required int totalRevenue,
    @JsonKey(name: 'average_order_value') required double averageOrderValue,
    @JsonKey(name: 'active_customers') required int activeCustomers,
    @JsonKey(name: 'status_breakdown')
    @Default(<String, int>{})
    Map<String, int> statusBreakdown,
    @JsonKey(name: 'top_tags')
    @Default(<AnalyticsTagCount>[])
    List<AnalyticsTagCount> topTags,
  }) = _AnalyticsSummary;

  factory AnalyticsSummary.fromJson(Map<String, dynamic> json) =>
      _$AnalyticsSummaryFromJson(json);
}

@freezed
class AnalyticsTagCount with _$AnalyticsTagCount {
  const factory AnalyticsTagCount({
    required String tag,
    required int count,
  }) = _AnalyticsTagCount;

  factory AnalyticsTagCount.fromJson(Map<String, dynamic> json) =>
      _$AnalyticsTagCountFromJson(json);
}

