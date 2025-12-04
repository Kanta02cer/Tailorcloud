import 'package:freezed_annotation/freezed_annotation.dart';

part 'ambassador.freezed.dart';
part 'ambassador.g.dart';

/// アンバサダーステータス
enum AmbassadorStatus {
  @JsonValue('Active')
  active,
  @JsonValue('Inactive')
  inactive,
  @JsonValue('Suspended')
  suspended,
}

/// 成果報酬ステータス
enum CommissionStatus {
  @JsonValue('Pending')
  pending,
  @JsonValue('Approved')
  approved,
  @JsonValue('Paid')
  paid,
  @JsonValue('Cancelled')
  cancelled,
}

/// アンバサダーモデル
@freezed
class Ambassador with _$Ambassador {
  const factory Ambassador({
    required String id,
    @JsonKey(name: 'tenant_id') required String tenantId,
    @JsonKey(name: 'user_id') required String userId,
    required String name,
    required String email,
    String? phone,
    required AmbassadorStatus status,
    @JsonKey(name: 'commission_rate') required double commissionRate,
    @JsonKey(name: 'total_sales') @Default(0) int totalSales,
    @JsonKey(name: 'total_commission') @Default(0) int totalCommission,
    @JsonKey(name: 'created_at') required DateTime createdAt,
    @JsonKey(name: 'updated_at') required DateTime updatedAt,
  }) = _Ambassador;

  factory Ambassador.fromJson(Map<String, dynamic> json) =>
      _$AmbassadorFromJson(json);
}

/// 成果報酬モデル
@freezed
class Commission with _$Commission {
  const factory Commission({
    required String id,
    @JsonKey(name: 'order_id') required String orderId,
    @JsonKey(name: 'ambassador_id') required String ambassadorId,
    @JsonKey(name: 'tenant_id') required String tenantId,
    @JsonKey(name: 'order_amount') required int orderAmount,
    @JsonKey(name: 'commission_rate') required double commissionRate,
    @JsonKey(name: 'commission_amount') required int commissionAmount,
    required CommissionStatus status,
    @JsonKey(name: 'paid_at') DateTime? paidAt,
    @JsonKey(name: 'created_at') required DateTime createdAt,
    @JsonKey(name: 'updated_at') required DateTime updatedAt,
  }) = _Commission;

  factory Commission.fromJson(Map<String, dynamic> json) =>
      _$CommissionFromJson(json);
}

/// アンバサダー拡張メソッド
extension AmbassadorExtension on Ambassador {
  /// 成果報酬率の表示（%）
  String get commissionRateDisplay {
    return '${(commissionRate * 100).toStringAsFixed(1)}%';
  }

  /// 累計売上の表示（円）
  String get totalSalesDisplay {
    return '¥${totalSales.toString().replaceAllMapped(
          RegExp(r'(\d)(?=(\d{3})+(?!\d))'),
          (Match m) => '${m[1]},',
        )}';
  }

  /// 累計報酬の表示（円）
  String get totalCommissionDisplay {
    return '¥${totalCommission.toString().replaceAllMapped(
          RegExp(r'(\d)(?=(\d{3})+(?!\d))'),
          (Match m) => '${m[1]},',
        )}';
  }
}

/// 成果報酬拡張メソッド
extension CommissionExtension on Commission {
  /// 注文金額の表示（円）
  String get orderAmountDisplay {
    return '¥${orderAmount.toString().replaceAllMapped(
          RegExp(r'(\d)(?=(\d{3})+(?!\d))'),
          (Match m) => '${m[1]},',
        )}';
  }

  /// 報酬額の表示（円）
  String get commissionAmountDisplay {
    return '¥${commissionAmount.toString().replaceAllMapped(
          RegExp(r'(\d)(?=(\d{3})+(?!\d))'),
          (Match m) => '${m[1]},',
        )}';
  }

  /// ステータスの表示ラベル
  String get statusLabel {
    switch (status) {
      case CommissionStatus.pending:
        return '未確定';
      case CommissionStatus.approved:
        return '確定';
      case CommissionStatus.paid:
        return '支払済み';
      case CommissionStatus.cancelled:
        return 'キャンセル';
    }
  }
}
