import 'package:flutter/material.dart';
import '../config/enterprise_theme.dart';
import '../widgets/app_navigation.dart';
import 'home/home_screen.dart';
import 'customers/customer_list_screen.dart';
import 'orders/order_list_screen.dart';
import 'inventory/inventory_screen.dart';
import 'order/visual_ordering_screen.dart';
import 'order/quick_order_screen.dart';

/// メイン画面（ナビゲーション付き）
class MainScreen extends StatefulWidget {
  const MainScreen({super.key});

  @override
  State<MainScreen> createState() => _MainScreenState();
}

class _MainScreenState extends State<MainScreen> {
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
        return const Scaffold(
          body: Center(
            child: Text(
              'Settings Screen\n(Coming Soon)',
              textAlign: TextAlign.center,
              style: TextStyle(color: EnterpriseColors.textPrimary),
            ),
          ),
        );
    }
  }
}

