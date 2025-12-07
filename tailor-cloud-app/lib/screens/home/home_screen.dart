import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';
import '../../config/enterprise_theme.dart';
import '../../models/analytics.dart';
import '../../providers/analytics_provider.dart';
import '../../utils/responsive.dart';
import '../../widgets/kpi_card.dart';
import '../../widgets/task_item.dart';
import '../order/quick_order_screen.dart';

/// Home画面（Dashboard）
class HomeScreen extends ConsumerWidget {
  const HomeScreen({super.key});

  static final NumberFormat _currencyFormatter =
      NumberFormat.currency(locale: 'ja_JP', symbol: '¥', decimalDigits: 0);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final analyticsAsync = ref.watch(analyticsSummaryProvider(rangeDays: 30));

    return Scaffold(
      backgroundColor: EnterpriseColors.deepBlack,
      body: SafeArea(
        child: Column(
          children: [
            // ヘッダー
            _buildHeader(context),

            // メインコンテンツ
            Expanded(
              child: LayoutBuilder(
                builder: (context, constraints) {
                  final horizontalPadding =
                      constraints.maxWidth < 800 ? 16.0 : 32.0;

                  return SingleChildScrollView(
                    padding: EdgeInsets.all(horizontalPadding),
                    child: Center(
                      child: ConstrainedBox(
                        constraints: BoxConstraints(
                          maxWidth: Responsive.formMaxWidth(context),
                        ),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            // KPIカード
                            analyticsAsync.when(
                              data: (summary) =>
                                  _buildKPIRow(context, summary),
                              loading: () => _buildKPIRow(context, null),
                              error: (_, __) => _buildKPIRow(context, null),
                            ),

                            const SizedBox(height: 40),

                            analyticsAsync.when(
                              data: (summary) =>
                                  _buildAnalyticsInsights(context, summary),
                              loading: () => const SizedBox.shrink(),
                              error: (_, __) => const SizedBox.shrink(),
                            ),

                            const SizedBox(height: 40),

                            // タスクリストとフィード
                            _buildContentRow(context),
                          ],
                        ),
                      ),
                    ),
                  );
                },
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

  Widget _buildKPIRow(BuildContext context, AnalyticsSummary? summary) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final isDesktop = Responsive.isDesktop(context);
        final isTablet = Responsive.isTablet(context);
        final crossAxisCount = isDesktop
            ? 4
            : isTablet
                ? 3
                : 2;

        final revenue =
            summary != null ? _formatCurrency(summary.totalRevenue) : '読み込み中';
        final orders = summary != null ? summary.totalOrders.toString() : '--';
        final avgOrder = summary != null
            ? _formatCurrency(summary.averageOrderValue.round())
            : '--';
        final activeCustomers =
            summary != null ? summary.activeCustomers.toString() : '--';

        return GridView.count(
          shrinkWrap: true,
          physics: const NeverScrollableScrollPhysics(),
          crossAxisCount: crossAxisCount,
          crossAxisSpacing: 24,
          mainAxisSpacing: 24,
          childAspectRatio: 1.2,
          children: [
            KPICard(
              title: '売上 (過去${summary?.rangeDays ?? 30}日)',
              value: revenue,
              trendColor: EnterpriseColors.statusAvailable,
              trendText: summary != null ? '平均: $avgOrder' : null,
              icon: Icons.attach_money,
            ),
            KPICard(
              title: '受注数',
              value: orders,
              subtitle: 'orders',
              progressColors: [
                EnterpriseColors.statusAvailable,
                EnterpriseColors.statusLowStock,
                EnterpriseColors.statusLowStock,
              ],
            ),
            KPICard(
              title: '稼働顧客',
              value: activeCustomers,
              subtitle: 'customers',
              icon: Icons.people_alt,
            ),
            KPICard(
              title: '平均単価',
              value: summary != null ? avgOrder : '--',
              subtitle: 'per order',
              icon: Icons.trending_up,
            ),
          ],
        );
      },
    );
  }

  Widget _buildAnalyticsInsights(
      BuildContext context, AnalyticsSummary summary) {
    return Card(
      color: EnterpriseColors.surfaceGray,
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.ssid_chart,
                    color: EnterpriseColors.metallicGold),
                const SizedBox(width: 8),
                Text(
                  'CRM Insights',
                  style: Theme.of(context).textTheme.titleLarge?.copyWith(
                        color: EnterpriseColors.textPrimary,
                        fontWeight: FontWeight.bold,
                      ),
                ),
              ],
            ),
            const SizedBox(height: 24),
            LayoutBuilder(
              builder: (context, constraints) {
                final isWide = constraints.maxWidth > 600;
                return isWide
                    ? Row(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Expanded(child: _buildStatusBreakdown(summary)),
                          const SizedBox(width: 24),
                          Expanded(child: _buildTopTags(summary)),
                        ],
                      )
                    : Column(
                        children: [
                          _buildStatusBreakdown(summary),
                          const SizedBox(height: 24),
                          _buildTopTags(summary),
                        ],
                      );
              },
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStatusBreakdown(AnalyticsSummary summary) {
    final entries = summary.statusBreakdown.entries.toList()
      ..sort((a, b) => b.value.compareTo(a.value));
    final totalCustomers = entries.fold<int>(
        0, (previousValue, element) => previousValue + element.value);
    final denominator = totalCustomers == 0 ? 1 : totalCustomers;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '顧客ステータス',
          style: TextStyle(
            color: EnterpriseColors.textSecondary,
            fontSize: 12,
          ),
        ),
        const SizedBox(height: 12),
        ...entries.map(
          (entry) => Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: Row(
              children: [
                Expanded(
                  child: LinearProgressIndicator(
                    value: entry.value / denominator,
                    backgroundColor: EnterpriseColors.deepBlack,
                    valueColor: AlwaysStoppedAnimation<Color>(
                      entry.key == 'vip'
                          ? EnterpriseColors.metallicGold
                          : EnterpriseColors.primaryBlue,
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Text(
                  '${entry.key.toUpperCase()} (${entry.value})',
                  style: const TextStyle(
                    color: EnterpriseColors.textPrimary,
                  ),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  String _formatCurrency(num value) {
    return _currencyFormatter.format(value);
  }

  Widget _buildTopTags(AnalyticsSummary summary) {
    if (summary.topTags.isEmpty) {
      return const Text(
        'タグのデータがまだありません',
        style: TextStyle(color: EnterpriseColors.textSecondary),
      );
    }
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '人気タグ',
          style: TextStyle(
            color: EnterpriseColors.textSecondary,
            fontSize: 12,
          ),
        ),
        const SizedBox(height: 12),
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children: summary.topTags
              .map(
                (tag) => Chip(
                  label: Text('${tag.tag} (${tag.count})'),
                  backgroundColor: EnterpriseColors.deepBlack,
                ),
              )
              .toList(),
        ),
      ],
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
