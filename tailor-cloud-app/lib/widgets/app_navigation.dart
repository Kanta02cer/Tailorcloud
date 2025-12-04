import 'package:flutter/material.dart';
import '../config/enterprise_theme.dart';

/// アプリケーションのナビゲーションタイプ
enum NavigationItem {
  dashboard,
  customers,
  orders,
  inventory,
  ordering,
  settings,
}

/// サイドバーナビゲーションウィジェット
class AppNavigation extends StatelessWidget {
  final NavigationItem selected;
  final Function(NavigationItem) onItemSelected;

  const AppNavigation({
    super.key,
    required this.selected,
    required this.onItemSelected,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 80,
      decoration: const BoxDecoration(
        color: EnterpriseColors.surfaceGray,
        border: Border(
          right: BorderSide(
            color: EnterpriseColors.borderGray,
            width: 1,
          ),
        ),
      ),
      child: Column(
        children: [
          const SizedBox(height: 24),

          // ロゴ
          Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: EnterpriseColors.metallicGold.withOpacity(0.1),
              borderRadius: BorderRadius.circular(20),
              border: Border.all(
                color: EnterpriseColors.metallicGold,
                width: 1,
              ),
            ),
            child: const Center(
              child: Text(
                'R',
                style: TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                  color: EnterpriseColors.metallicGold,
                ),
              ),
            ),
          ),

          const SizedBox(height: 32),

          // ナビゲーションアイテム
          _buildNavItem(
            context,
            icon: Icons.dashboard,
            label: 'Home',
            item: NavigationItem.dashboard,
          ),

          const SizedBox(height: 24),

          _buildNavItem(
            context,
            icon: Icons.people,
            label: 'Customer',
            item: NavigationItem.customers,
          ),

          const SizedBox(height: 24),

          _buildNavItem(
            context,
            icon: Icons.shopping_cart,
            label: 'Orders',
            item: NavigationItem.orders,
          ),

          const SizedBox(height: 24),

          _buildNavItem(
            context,
            icon: Icons.inventory_2,
            label: 'Fabric',
            item: NavigationItem.inventory,
          ),

          const SizedBox(height: 24),

          _buildNavItem(
            context,
            icon: Icons.create,
            label: 'Order',
            item: NavigationItem.ordering,
          ),

          const Spacer(),

          // 設定
          _buildNavItem(
            context,
            icon: Icons.settings,
            label: '',
            item: NavigationItem.settings,
          ),

          const SizedBox(height: 24),
        ],
      ),
    );
  }

  Widget _buildNavItem(
    BuildContext context, {
    required IconData icon,
    required String label,
    required NavigationItem item,
  }) {
    final isSelected = selected == item;

    return InkWell(
      onTap: () => onItemSelected(item),
      child: Column(
        children: [
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(12),
              color: isSelected
                  ? EnterpriseColors.regalisBlack
                  : Colors.transparent,
            ),
            child: Icon(
              icon,
              size: 24,
              color: isSelected
                  ? EnterpriseColors.metallicGold
                  : EnterpriseColors.textSecondary,
            ),
          ),
          if (label.isNotEmpty) ...[
            const SizedBox(height: 4),
            Text(
              label.toUpperCase(),
              style: TextStyle(
                fontSize: 9,
                letterSpacing: 2,
                color: isSelected
                    ? EnterpriseColors.metallicGold
                    : EnterpriseColors.textSecondary,
              ),
            ),
          ],
        ],
      ),
    );
  }
}
