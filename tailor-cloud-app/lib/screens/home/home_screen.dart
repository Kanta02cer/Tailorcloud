import 'package:flutter/material.dart';
import '../../config/enterprise_theme.dart';
import '../../widgets/kpi_card.dart';
import '../../widgets/task_item.dart';
import '../order/quick_order_screen.dart';

/// Home画面（Dashboard）
class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: EnterpriseColors.deepBlack,
      body: SafeArea(
        child: Column(
          children: [
            // ヘッダー
            _buildHeader(context),
            
            // メインコンテンツ
            Expanded(
              child: SingleChildScrollView(
                padding: const EdgeInsets.all(32),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // KPIカード
                    _buildKPIRow(context),
                    
                    const SizedBox(height: 40),
                    
                    // タスクリストとフィード
                    _buildContentRow(context),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildHeader(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 24),
      decoration: const BoxDecoration(
        color: EnterpriseColors.surfaceGray,
        border: Border(
          bottom: BorderSide(
            color: EnterpriseColors.borderGray,
            width: 1,
          ),
        ),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'STORE DASHBOARD',
                style: Theme.of(context).textTheme.labelSmall?.copyWith(
                      color: EnterpriseColors.metallicGold,
                      letterSpacing: 2,
                    ),
              ),
              const SizedBox(height: 4),
              Text(
                'Regalis Yotsuya Salon',
                style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
              ),
            ],
          ),
          Row(
            children: [
              OutlinedButton.icon(
                onPressed: () {},
                icon: const Icon(Icons.calendar_today, size: 16),
                label: const Text('Schedule'),
                style: OutlinedButton.styleFrom(
                  foregroundColor: EnterpriseColors.textPrimary,
                  side: const BorderSide(color: EnterpriseColors.borderGray),
                  padding: const EdgeInsets.symmetric(
                    horizontal: 16,
                    vertical: 12,
                  ),
                ),
              ),
              const SizedBox(width: 12),
              ElevatedButton.icon(
                onPressed: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => const QuickOrderScreen(),
                    ),
                  );
                },
                icon: const Icon(Icons.add, size: 16),
                label: const Text('クイック発注'),
                style: ElevatedButton.styleFrom(
                  backgroundColor: EnterpriseColors.metallicGold,
                  foregroundColor: Colors.black,
                  padding: const EdgeInsets.symmetric(
                    horizontal: 16,
                    vertical: 12,
                  ),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildKPIRow(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final isWide = constraints.maxWidth > 1200;
        final crossAxisCount = isWide ? 4 : 2;
        
        return GridView.count(
          shrinkWrap: true,
          physics: const NeverScrollableScrollPhysics(),
          crossAxisCount: crossAxisCount,
          crossAxisSpacing: 24,
          mainAxisSpacing: 24,
          childAspectRatio: 1.2,
          children: [
            KPICard(
              title: 'Total Sales (Today)',
              value: '¥480,000',
              trendColor: EnterpriseColors.statusAvailable,
              trendText: '+15% vs Last Week',
              icon: Icons.attach_money,
            ),
            KPICard(
              title: 'Active Orders',
              value: '24',
              subtitle: 'suits',
              progressColors: [
                EnterpriseColors.statusAvailable,
                EnterpriseColors.statusLowStock,
                EnterpriseColors.statusLowStock,
              ],
            ),
            KPICard(
              title: 'Stock Alerts',
              value: '3',
              subtitle: 'items',
              trendColor: EnterpriseColors.statusOutOfStock,
              trendText: 'VBC Navy Low Stock',
              icon: Icons.warning_amber,
            ),
            KPICard(
              title: 'Factory Status',
              value: 'Running',
              subtitle: 'Avg Lead Time: 21 Days',
            ),
          ],
        );
      },
    );
  }

  Widget _buildContentRow(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final isWide = constraints.maxWidth > 1000;
        
        if (isWide) {
          return Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // タスクリスト
              Expanded(
                flex: 2,
                child: _buildTaskList(context),
              ),
              const SizedBox(width: 32),
              // フィード
              Expanded(
                child: _buildFeed(context),
              ),
            ],
          );
        } else {
          return Column(
            children: [
              _buildTaskList(context),
              const SizedBox(height: 32),
              _buildFeed(context),
            ],
          );
        }
      },
    );
  }

  Widget _buildTaskList(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            Icon(
              Icons.check_box,
              size: 20,
              color: EnterpriseColors.metallicGold,
            ),
            const SizedBox(width: 8),
            Text(
              'Action Required',
              style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
            ),
          ],
        ),
        const SizedBox(height: 16),
        TaskItem(
          title: 'Initial Fitting: K. Tanaka',
          subtitle: '14:00 - 15:30 / Milestone Line',
          icon: Icons.straighten,
          borderColor: EnterpriseColors.metallicGold,
          iconColor: EnterpriseColors.metallicGold,
          actionLabel: 'View Profile',
          onActionTap: () {},
        ),
        const SizedBox(height: 12),
        TaskItem(
          title: 'Fabric Out of Stock',
          subtitle: 'Order #1024 - Miyuki 3055 is unavailable.',
          icon: Icons.warning,
          borderColor: EnterpriseColors.statusOutOfStock,
          iconColor: EnterpriseColors.statusOutOfStock,
          actionLabel: 'Change Fabric',
          onActionTap: () {},
        ),
      ],
    );
  }

  Widget _buildFeed(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Mill Updates',
          style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
        ),
        const SizedBox(height: 16),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(20),
            child: Column(
              children: [
                _buildMillUpdateItem(
                  context,
                  'VBC',
                  'New Collection',
                  EnterpriseColors.metallicGold,
                  '"Revenge" 2025AW collection is now available for order.',
                  '10 mins ago',
                ),
                const SizedBox(height: 24),
                _buildMillUpdateItem(
                  context,
                  'LP',
                  'Price Revision',
                  EnterpriseColors.statusOutOfStock,
                  'Loro Piana Australis series price will increase by 5% from next month.',
                  '2 hours ago',
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildMillUpdateItem(
    BuildContext context,
    String avatar,
    String title,
    Color titleColor,
    String content,
    String time,
  ) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Container(
          width: 32,
          height: 32,
          decoration: BoxDecoration(
            color: EnterpriseColors.surfaceGray,
            borderRadius: BorderRadius.circular(16),
          ),
          child: Center(
            child: Text(
              avatar,
              style: const TextStyle(
                fontSize: 10,
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                title,
                style: TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.bold,
                  color: titleColor,
                ),
              ),
              const SizedBox(height: 4),
              Text(
                content,
                style: Theme.of(context).textTheme.bodySmall,
              ),
              const SizedBox(height: 8),
              Text(
                time,
                style: Theme.of(context).textTheme.labelSmall,
              ),
            ],
          ),
        ),
      ],
    );
  }
}

