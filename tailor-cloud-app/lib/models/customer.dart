import 'package:freezed_annotation/freezed_annotation.dart';

part 'customer.freezed.dart';
part 'customer.g.dart';

/// 顧客モデル
@freezed
class Customer with _$Customer {
  const factory Customer({
    required String id,
    @JsonKey(name: 'tenant_id') required String tenantId,
    required String name,
    @JsonKey(name: 'name_kana') String? nameKana,
    String? email,
    String? phone,
    String? address,
    @JsonKey(name: 'customer_status') String? status,
    @Default(<String>[]) List<String> tags,
    @JsonKey(name: 'vip_rank') int? vipRank,
    @JsonKey(name: 'ltv_score') double? ltvScore,
    @JsonKey(name: 'lifetime_value') double? lifetimeValue,
    @JsonKey(name: 'preferred_channel') String? preferredChannel,
    @JsonKey(name: 'lead_source') String? leadSource,
    String? notes,
    String? occupation,
    @JsonKey(name: 'annual_income_range') String? annualIncomeRange,
    @JsonKey(name: 'preferred_archetype') String? preferredArchetype,
    @JsonKey(name: 'diagnosis_count') int? diagnosisCount,
    @JsonKey(name: 'last_interaction_at') DateTime? lastInteractionAt,
    @JsonKey(name: 'interaction_notes')
    @Default(<CustomerInteraction>[])
    List<CustomerInteraction> interactionNotes,
    @JsonKey(name: 'created_at') required DateTime createdAt,
    @JsonKey(name: 'updated_at') required DateTime updatedAt,
  }) = _Customer;

  factory Customer.fromJson(Map<String, dynamic> json) =>
      _$CustomerFromJson(json);
}

/// 顧客作成リクエスト
@freezed
class CreateCustomerRequest with _$CreateCustomerRequest {
  const factory CreateCustomerRequest({
    required String name,
    @JsonKey(name: 'name_kana') String? nameKana,
    String? email,
    String? phone,
    String? address,
    @JsonKey(name: 'status') String? status,
    @JsonKey(name: 'tags') List<String>? tags,
    @JsonKey(name: 'vip_rank') int? vipRank,
    @JsonKey(name: 'ltv_score') double? ltvScore,
    @JsonKey(name: 'lifetime_value') double? lifetimeValue,
    @JsonKey(name: 'preferred_channel') String? preferredChannel,
    @JsonKey(name: 'lead_source') String? leadSource,
    String? notes,
    String? occupation,
    @JsonKey(name: 'annual_income_range') String? annualIncomeRange,
    @JsonKey(name: 'preferred_archetype') String? preferredArchetype,
    @JsonKey(name: 'diagnosis_count') int? diagnosisCount,
    @JsonKey(name: 'last_interaction_at') DateTime? lastInteractionAt,
    @JsonKey(name: 'interactions') List<CustomerInteraction>? interactions,
  }) = _CreateCustomerRequest;

  factory CreateCustomerRequest.fromJson(Map<String, dynamic> json) =>
      _$CreateCustomerRequestFromJson(json);
}

@freezed
class CustomerInteraction with _$CustomerInteraction {
  const factory CustomerInteraction({
    required String type,
    required String note,
    required DateTime timestamp,
    String? staff,
  }) = _CustomerInteraction;

  factory CustomerInteraction.fromJson(Map<String, dynamic> json) =>
      _$CustomerInteractionFromJson(json);
}
