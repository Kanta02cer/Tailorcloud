import 'package:freezed_annotation/freezed_annotation.dart';

part 'fabric.freezed.dart';
part 'fabric.g.dart';

/// 在庫ステータス
enum StockStatus {
  @JsonValue('Available')
  available,
  @JsonValue('Limited')
  limited,
  @JsonValue('SoldOut')
  soldOut,
}

/// 生地モデル
@freezed
class Fabric with _$Fabric {
  const factory Fabric({
    required String id,
    @JsonKey(name: 'tenant_id') String? tenantId,
    @JsonKey(name: 'supplier_id') String? supplierId,
    String? brand,
    required String name,
    String? sku,
    @JsonKey(name: 'stock_amount') double? stockAmount,
    int? price,
    @JsonKey(name: 'cost_price') int? costPrice,
    @JsonKey(name: 'sales_price') int? salesPrice,
    @JsonKey(name: 'stock_status') StockStatus? stockStatus,
    @JsonKey(name: 'image_url') String? imageUrl,
    @JsonKey(name: 'minimum_order') @Default(3.2) double minimumOrder,
    @JsonKey(name: 'created_at') DateTime? createdAt,
    @JsonKey(name: 'updated_at') DateTime? updatedAt,
  }) = _Fabric;

  factory Fabric.fromJson(Map<String, dynamic> json) => _$FabricFromJson(json);
}

/// 生地拡張メソッド
extension FabricExtension on Fabric {
  /// 在庫ステータスの表示ラベル
  String get stockStatusLabel {
    if (stockStatus == null) return '不明';
    switch (stockStatus!) {
      case StockStatus.available:
        return '在庫あり';
      case StockStatus.limited:
        return '在庫残りわずか';
      case StockStatus.soldOut:
        return '在庫切れ';
    }
  }

  /// 在庫ステータスのカラーコード
  String get stockStatusColor {
    if (stockStatus == null) return '#6B7280'; // グレー
    switch (stockStatus!) {
      case StockStatus.available:
        return '#10B981'; // 緑
      case StockStatus.limited:
        return '#F59E0B'; // 黄色
      case StockStatus.soldOut:
        return '#EF4444'; // 赤
    }
  }

  /// 在庫数量の表示（Limitedの場合）
  String? get stockAmountDisplay {
    if (stockStatus == StockStatus.limited && stockAmount != null) {
      return '残り${stockAmount!.toStringAsFixed(1)}m';
    }
    return null;
  }
}
