import 'package:flutter/material.dart';
import '../config/enterprise_theme.dart';

/// KPIカードウィジェット
class KPICard extends StatelessWidget {
  final String title;
  final String value;
  final String? subtitle;
  final IconData? icon;
  final Color? trendColor;
  final String? trendText;
  final List<Color>? progressColors; // 進捗バーの色（複数）

  const KPICard({
    super.key,
    required this.title,
    required this.value,
    this.subtitle,
    this.icon,
    this.trendColor,
    this.trendText,
    this.progressColors,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Container(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // タイトル
            Text(
              title.toUpperCase(),
              style: Theme.of(context).textTheme.labelSmall?.copyWith(
                    letterSpacing: 2,
                  ),
            ),
            const SizedBox(height: 8),
            
            // 値
            Text(
              value,
              style: Theme.of(context).textTheme.displayLarge?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
            ),
            
            // サブタイトルまたはトレンド
            if (subtitle != null) ...[
              const SizedBox(height: 8),
              Text(
                subtitle!,
                style: Theme.of(context).textTheme.bodySmall,
              ),
            ],
            
            // トレンド表示
            if (trendText != null && trendColor != null) ...[
              const SizedBox(height: 8),
              Row(
                children: [
                  Icon(
                    Icons.trending_up,
                    size: 12,
                    color: trendColor,
                  ),
                  const SizedBox(width: 4),
                  Text(
                    trendText!,
                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: trendColor,
                        ),
                  ),
                ],
              ),
            ],
            
            // 進捗バー（複数の色）
            if (progressColors != null && progressColors!.isNotEmpty) ...[
              const SizedBox(height: 12),
              Row(
                children: progressColors!
                    .map((color) => Expanded(
                          child: Container(
                            height: 6,
                            margin: const EdgeInsets.only(right: 2),
                            decoration: BoxDecoration(
                              color: color,
                              borderRadius: BorderRadius.circular(3),
                            ),
                          ),
                        ))
                    .toList(),
              ),
            ],
            
          ],
        ),
      ),
    );
  }
}

