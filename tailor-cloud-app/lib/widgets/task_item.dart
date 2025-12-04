import 'package:flutter/material.dart';
import '../config/enterprise_theme.dart';

/// タスクアイテムウィジェット
class TaskItem extends StatelessWidget {
  final String title;
  final String subtitle;
  final IconData icon;
  final Color borderColor;
  final Color iconColor;
  final String? actionLabel;
  final VoidCallback? onTap;
  final VoidCallback? onActionTap;

  const TaskItem({
    super.key,
    required this.title,
    required this.subtitle,
    required this.icon,
    required this.borderColor,
    required this.iconColor,
    this.actionLabel,
    this.onTap,
    this.onActionTap,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(16),
        child: Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(16),
            border: Border(
              left: BorderSide(
                color: borderColor,
                width: 4,
              ),
            ),
          ),
          child: Row(
            children: [
              // アイコン
              Container(
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: EnterpriseColors.surfaceGray,
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Icon(
                  icon,
                  size: 20,
                  color: iconColor,
                ),
              ),
              const SizedBox(width: 16),

              // テキスト情報
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      title,
                      style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                            fontWeight: FontWeight.bold,
                          ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      subtitle,
                      style: Theme.of(context).textTheme.bodySmall,
                    ),
                  ],
                ),
              ),

              // アクションボタン
              if (actionLabel != null)
                TextButton(
                  onPressed: onActionTap,
                  style: TextButton.styleFrom(
                    backgroundColor: EnterpriseColors.surfaceGray,
                    foregroundColor: EnterpriseColors.textSecondary,
                    padding: const EdgeInsets.symmetric(
                      horizontal: 12,
                      vertical: 6,
                    ),
                    minimumSize: Size.zero,
                    tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                  ),
                  child: Text(
                    actionLabel!,
                    style: const TextStyle(fontSize: 12),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}
