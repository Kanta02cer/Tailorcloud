import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../config/enterprise_theme.dart';
import '../widgets/app_navigation.dart';
import '../providers/auth_provider.dart';
import 'home/home_screen.dart';
import 'customers/customer_list_screen.dart';
import 'orders/order_list_screen.dart';
import 'inventory/inventory_screen.dart';
import 'order/visual_ordering_screen.dart';
import 'order/quick_order_screen.dart';

/// メイン画面（ナビゲーション付き）
class MainScreen extends ConsumerStatefulWidget {
  const MainScreen({super.key});

  @override
  ConsumerState<MainScreen> createState() => _MainScreenState();
}

class _MainScreenState extends ConsumerState<MainScreen> {
  NavigationItem _selectedItem = NavigationItem.dashboard;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: EnterpriseColors.deepBlack,
      body: Row(
        children: [
          // サイドバーナビゲーション
          AppNavigation(
            selected: _selectedItem,
            onItemSelected: (item) {
              setState(() {
                _selectedItem = item;
              });
            },
          ),
          
          // メインコンテンツ
          Expanded(
            child: _buildContent(),
          ),
        ],
      ),
    );
  }

  Widget _buildContent() {
    switch (_selectedItem) {
      case NavigationItem.dashboard:
        return const HomeScreen();
      case NavigationItem.customers:
        return const CustomerListScreen();
      case NavigationItem.orders:
        return const OrderListScreen();
      case NavigationItem.inventory:
        return const InventoryScreen();
      case NavigationItem.ordering:
        return const VisualOrderingScreen();
      case NavigationItem.settings:
        return _buildSettingsScreen();
    }
  }

  Widget _buildSettingsScreen() {
    final currentUser = ref.watch(currentUserProvider);

    return Scaffold(
      backgroundColor: EnterpriseColors.deepBlack,
      appBar: AppBar(
        backgroundColor: EnterpriseColors.surfaceGray,
        title: const Text(
          '設定',
          style: TextStyle(
            color: EnterpriseColors.textPrimary,
            fontSize: 20,
            fontWeight: FontWeight.bold,
          ),
        ),
        elevation: 0,
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // ユーザー情報セクション
          Card(
            color: EnterpriseColors.surfaceGray,
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    'ユーザー情報',
                    style: TextStyle(
                      color: EnterpriseColors.textSecondary,
                      fontSize: 12,
                    ),
                  ),
                  const SizedBox(height: 8),
                  if (currentUser != null) ...[
                    Text(
                      currentUser.email ?? 'メールアドレス未設定',
                      style: const TextStyle(
                        color: EnterpriseColors.textPrimary,
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      'UID: ${currentUser.uid}',
                      style: const TextStyle(
                        color: EnterpriseColors.textTertiary,
                        fontSize: 12,
                      ),
                    ),
                  ],
                ],
              ),
            ),
          ),
          const SizedBox(height: 24),

          // ログアウトボタン
          ElevatedButton(
            onPressed: () async {
              // 確認ダイアログ
              final confirmed = await showDialog<bool>(
                context: context,
                builder: (context) => AlertDialog(
                  backgroundColor: EnterpriseColors.surfaceGray,
                  title: const Text(
                    'ログアウト',
                    style: TextStyle(color: EnterpriseColors.textPrimary),
                  ),
                  content: const Text(
                    'ログアウトしますか？',
                    style: TextStyle(color: EnterpriseColors.textSecondary),
                  ),
                  actions: [
                    TextButton(
                      onPressed: () => Navigator.of(context).pop(false),
                      child: const Text(
                        'キャンセル',
                        style: TextStyle(color: EnterpriseColors.textSecondary),
                      ),
                    ),
                    TextButton(
                      onPressed: () => Navigator.of(context).pop(true),
                      style: TextButton.styleFrom(
                        foregroundColor: EnterpriseColors.errorRed,
                      ),
                      child: const Text('ログアウト'),
                    ),
                  ],
                ),
              );

              if (confirmed == true && mounted) {
                try {
                  await ref.read(signOutProvider.future);
                  // 認証状態の変更により自動的にログイン画面に遷移
                  if (mounted) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(
                        content: Text('ログアウトしました'),
                        backgroundColor: EnterpriseColors.successGreen,
                      ),
                    );
                  }
                } catch (e) {
                  if (mounted) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(
                        content: Text('ログアウトに失敗しました: $e'),
                        backgroundColor: EnterpriseColors.errorRed,
                      ),
                    );
                  }
                }
              }
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: EnterpriseColors.errorRed,
              foregroundColor: Colors.white,
              padding: const EdgeInsets.symmetric(vertical: 16),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8),
              ),
            ),
            child: const Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(Icons.logout),
                SizedBox(width: 8),
                Text(
                  'ログアウト',
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

