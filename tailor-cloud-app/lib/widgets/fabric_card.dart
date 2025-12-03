import 'package:flutter/material.dart';
import '../config/enterprise_theme.dart';
import '../models/fabric.dart';

/// 生地カードウィジェット
class FabricCard extends StatelessWidget {
  final Fabric fabric;
  final VoidCallback? onTap;
  final bool isRecommended;

  const FabricCard({
    super.key,
    required this.fabric,
    this.onTap,
    this.isRecommended = false,
  });

  @override
  Widget build(BuildContext context) {
    final statusColor = _getStatusColor(fabric.stockStatus);
    
    return Card(
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(16),
        child: AspectRatio(
          aspectRatio: 3 / 4,
          child: Stack(
            children: [
              // 画像
              if (fabric.imageUrl != null)
                ClipRRect(
                  borderRadius: BorderRadius.circular(16),
                  child: Image.network(
                    fabric.imageUrl!,
                    width: double.infinity,
                    height: double.infinity,
                    fit: BoxFit.cover,
                    errorBuilder: (context, error, stackTrace) {
                      return Container(
                        color: EnterpriseColors.surfaceGray,
                        child: const Icon(
                          Icons.image_not_supported,
                          size: 48,
                          color: EnterpriseColors.textTertiary,
                        ),
                      );
                    },
                  ),
                )
              else
                Container(
                  color: EnterpriseColors.surfaceGray,
                  child: const Icon(
                    Icons.image_not_supported,
                    size: 48,
                    color: EnterpriseColors.textTertiary,
                  ),
                ),
              
              // 推奨バッジ
              if (isRecommended)
                Positioned(
                  top: 0,
                  left: 0,
                  right: 0,
                  child: Container(
                    padding: const EdgeInsets.symmetric(vertical: 4),
                    decoration: const BoxDecoration(
                      color: EnterpriseColors.metallicGold,
                      borderRadius: BorderRadius.only(
                        topLeft: Radius.circular(16),
                        topRight: Radius.circular(16),
                      ),
                    ),
                    child: const Text(
                      'AI RECOMMEND',
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 10,
                        fontWeight: FontWeight.bold,
                        color: Colors.black,
                      ),
                    ),
                  ),
                ),
              
              // 在庫ステータスバッジ
              Positioned(
                top: isRecommended ? 24 : 12,
                left: 12,
                child: Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 8,
                    vertical: 4,
                  ),
                  decoration: BoxDecoration(
                    color: Colors.black.withOpacity(0.6),
                    borderRadius: BorderRadius.circular(4),
                    border: Border.all(
                      color: Colors.white.withOpacity(0.1),
                      width: 1,
                    ),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Container(
                        width: 8,
                        height: 8,
                        decoration: BoxDecoration(
                          color: statusColor,
                          shape: BoxShape.circle,
                          boxShadow: [
                            BoxShadow(
                              color: statusColor.withOpacity(0.8),
                              blurRadius: 8,
                              spreadRadius: 0,
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(width: 6),
                      Text(
                        _getStatusText(fabric.stockStatus),
                        style: const TextStyle(
                          fontSize: 10,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              
              // 在庫数量表示（Limitedの場合）
              if (fabric.stockStatus == StockStatus.limited)
                Positioned(
                  top: isRecommended ? 48 : 36,
                  left: 12,
                  child: Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 4,
                    ),
                    decoration: BoxDecoration(
                      color: EnterpriseColors.statusLowStock.withOpacity(0.9),
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: Text(
                      fabric.stockAmountDisplay ?? '',
                      style: const TextStyle(
                        fontSize: 10,
                        fontWeight: FontWeight.bold,
                        color: Colors.white,
                      ),
                    ),
                  ),
                ),
              
              // 在庫切れオーバーレイ
              if (fabric.stockStatus == StockStatus.soldOut)
                Positioned.fill(
                  child: Container(
                    decoration: BoxDecoration(
                      color: Colors.black.withOpacity(0.7),
                      borderRadius: BorderRadius.circular(16),
                    ),
                    child: Center(
                      child: Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 12,
                          vertical: 6,
                        ),
                        decoration: const BoxDecoration(
                          color: Colors.black,
                          borderRadius: BorderRadius.all(Radius.circular(8)),
                        ),
                        child: const Text(
                          'Out of Stock',
                          style: TextStyle(
                            fontSize: 12,
                            color: EnterpriseColors.textSecondary,
                          ),
                        ),
                      ),
                    ),
                  ),
                ),
              
              // 下部情報
              Positioned(
                bottom: 0,
                left: 0,
                right: 0,
                child: Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    gradient: LinearGradient(
                      begin: Alignment.topCenter,
                      end: Alignment.bottomCenter,
                      colors: [
                        Colors.transparent,
                        Colors.black.withOpacity(0.8),
                        Colors.black,
                      ],
                    ),
                    borderRadius: const BorderRadius.only(
                      bottomLeft: Radius.circular(16),
                      bottomRight: Radius.circular(16),
                    ),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text(
                        'V.B.C - Super 110\'s', // TODO: サプライヤー名を取得
                        style: const TextStyle(
                          fontSize: 12,
                          color: EnterpriseColors.textSecondary,
                        ),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                      const SizedBox(height: 4),
                      Text(
                        fabric.name,
                        style: const TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                      const SizedBox(height: 4),
                      Text(
                        '¥${fabric.price.toString().replaceAllMapped(RegExp(r'(\d)(?=(\d{3})+(?!\d))'), (m) => '${m[1]},')}~',
                        style: const TextStyle(
                          fontSize: 14,
                          color: EnterpriseColors.metallicGold,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Color _getStatusColor(StockStatus status) {
    switch (status) {
      case StockStatus.available:
        return EnterpriseColors.statusAvailable;
      case StockStatus.limited:
        return EnterpriseColors.statusLowStock;
      case StockStatus.soldOut:
        return EnterpriseColors.statusOutOfStock;
    }
  }

  String _getStatusText(StockStatus status) {
    switch (status) {
      case StockStatus.available:
        return 'Available';
      case StockStatus.limited:
        return 'Low Stock';
      case StockStatus.soldOut:
        return 'Out of Stock';
    }
  }
}

