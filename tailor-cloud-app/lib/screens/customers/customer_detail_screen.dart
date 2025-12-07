import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../config/enterprise_theme.dart';
import '../../models/customer.dart';
import '../../providers/customer_provider.dart';
import 'customer_edit_screen.dart';
import '../orders/order_list_screen.dart';

/// 顧客詳細画面
class CustomerDetailScreen extends ConsumerWidget {
  final String customerId;

  const CustomerDetailScreen({
    super.key,
    required this.customerId,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final customerAsync = ref.watch(customerProvider(customerId));

    return Scaffold(
      backgroundColor: EnterpriseColors.deepBlack,
      appBar: AppBar(
        backgroundColor: EnterpriseColors.surfaceGray,
        title: const Text(
          '顧客詳細',
          style: TextStyle(
            color: EnterpriseColors.textPrimary,
            fontSize: 20,
            fontWeight: FontWeight.bold,
          ),
        ),
        elevation: 0,
        actions: [
          customerAsync.when(
            data: (customer) => IconButton(
              icon: const Icon(
                Icons.edit,
                color: EnterpriseColors.primaryBlue,
              ),
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (context) => CustomerEditScreen(
                      customer: customer,
                    ),
                  ),
                ).then((_) {
                  // 編集後、顧客情報を更新
                  ref.invalidate(customerProvider(customerId));
                });
              },
            ),
            loading: () => const SizedBox.shrink(),
            error: (_, __) => const SizedBox.shrink(),
          ),
        ],
      ),
      body: customerAsync.when(
        data: (customer) => SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // 顧客情報カード
              _buildCustomerInfoCard(context, customer),

              const SizedBox(height: 24),

              if (customer.interactionNotes.isNotEmpty)
                _buildInteractionTimeline(customer),

              if (customer.interactionNotes.isNotEmpty)
                const SizedBox(height: 24),

              // 注文履歴セクション
              _buildOrdersSection(context, customer),
            ],
          ),
        ),
        loading: () => const Center(
          child: CircularProgressIndicator(
            color: EnterpriseColors.primaryBlue,
          ),
        ),
        error: (error, stack) => Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(
                Icons.error_outline,
                color: EnterpriseColors.errorRed,
                size: 48,
              ),
              const SizedBox(height: 16),
              Text(
                'エラーが発生しました',
                style: const TextStyle(
                  color: EnterpriseColors.textPrimary,
                  fontSize: 16,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                error.toString(),
                style: const TextStyle(
                  color: EnterpriseColors.textSecondary,
                  fontSize: 12,
                ),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: () {
                  ref.invalidate(customerProvider(customerId));
                },
                style: ElevatedButton.styleFrom(
                  backgroundColor: EnterpriseColors.primaryBlue,
                ),
                child: const Text('再試行'),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildCustomerInfoCard(BuildContext context, Customer customer) {
    return Card(
      color: EnterpriseColors.surfaceGray,
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 顧客名
            Row(
              children: [
                CircleAvatar(
                  radius: 32,
                  backgroundColor: EnterpriseColors.primaryBlue,
                  child: Text(
                    customer.name.isNotEmpty
                        ? customer.name[0].toUpperCase()
                        : '?',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        customer.name,
                        style: const TextStyle(
                          color: EnterpriseColors.textPrimary,
                          fontSize: 24,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      if (customer.nameKana != null) ...[
                        const SizedBox(height: 4),
                        Text(
                          customer.nameKana!,
                          style: const TextStyle(
                            color: EnterpriseColors.textSecondary,
                            fontSize: 14,
                          ),
                        ),
                      ],
                      const SizedBox(height: 8),
                      Wrap(
                        spacing: 8,
                        crossAxisAlignment: WrapCrossAlignment.center,
                        children: [
                          _buildStatusChip(customer.status ?? 'lead'),
                          if (customer.vipRank != null)
                            Chip(
                              backgroundColor: EnterpriseColors.surfaceGray,
                              label: Text('VIP ${customer.vipRank}'),
                            ),
                          if (customer.lifetimeValue != null)
                            Chip(
                              backgroundColor: EnterpriseColors.surfaceGray,
                              label: Text(
                                'LTV ¥${customer.lifetimeValue!.toStringAsFixed(0)}',
                              ),
                            ),
                        ],
                      ),
                    ],
                  ),
                ),
              ],
            ),

            const Divider(
              color: EnterpriseColors.borderGray,
              height: 32,
            ),

            if (customer.tags.isNotEmpty) ...[
              const Text(
                'タグ',
                style: TextStyle(
                  color: EnterpriseColors.textSecondary,
                  fontSize: 12,
                ),
              ),
              const SizedBox(height: 8),
              Wrap(
                spacing: 8,
                runSpacing: 8,
                children: customer.tags
                    .map(
                      (tag) => Chip(
                        label: Text(tag),
                        backgroundColor: EnterpriseColors.deepBlack,
                      ),
                    )
                    .toList(),
              ),
              const SizedBox(height: 16),
            ],

            // 連絡先情報
            if (customer.email != null) ...[
              _buildInfoRow(
                icon: Icons.email,
                label: 'メールアドレス',
                value: customer.email!,
              ),
              const SizedBox(height: 16),
            ],
            if (customer.phone != null) ...[
              _buildInfoRow(
                icon: Icons.phone,
                label: '電話番号',
                value: customer.phone!,
              ),
              const SizedBox(height: 16),
            ],
            if (customer.address != null) ...[
              _buildInfoRow(
                icon: Icons.location_on,
                label: '住所',
                value: customer.address!,
              ),
              const SizedBox(height: 16),
            ],

            // 登録日時
            _buildInfoRow(
              icon: Icons.calendar_today,
              label: '登録日',
              value: _formatDate(customer.createdAt),
            ),
            if (customer.notes != null && customer.notes!.isNotEmpty) ...[
              const SizedBox(height: 16),
              _buildInfoRow(
                icon: Icons.note_alt,
                label: '社内メモ',
                value: customer.notes!,
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildInfoRow({
    required IconData icon,
    required String label,
    required String value,
  }) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Icon(
          icon,
          color: EnterpriseColors.primaryBlue,
          size: 20,
        ),
        const SizedBox(width: 12),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                label,
                style: const TextStyle(
                  color: EnterpriseColors.textSecondary,
                  fontSize: 12,
                ),
              ),
              const SizedBox(height: 4),
              Text(
                value,
                style: const TextStyle(
                  color: EnterpriseColors.textPrimary,
                  fontSize: 16,
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildStatusChip(String status) {
    Color color;
    switch (status) {
      case 'vip':
        color = EnterpriseColors.metallicGold;
        break;
      case 'active':
        color = EnterpriseColors.statusAvailable;
        break;
      case 'prospect':
        color = EnterpriseColors.primaryBlue;
        break;
      case 'inactive':
        color = EnterpriseColors.textSecondary;
        break;
      default:
        color = EnterpriseColors.statusLowStock;
    }
    return Chip(
      label: Text(status.toUpperCase()),
      backgroundColor: color.withOpacity(0.2),
      labelStyle: TextStyle(
        color: color,
        fontWeight: FontWeight.bold,
      ),
    );
  }

  Widget _buildInteractionTimeline(Customer customer) {
    return Card(
      color: EnterpriseColors.surfaceGray,
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'コンタクト履歴',
              style: TextStyle(
                color: EnterpriseColors.textPrimary,
                fontSize: 18,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16),
            ...customer.interactionNotes.map(
              (note) => ListTile(
                contentPadding: EdgeInsets.zero,
                leading: Icon(
                  Icons.bubble_chart,
                  color: EnterpriseColors.primaryBlue,
                  size: 20,
                ),
                title: Text(
                  note.note,
                  style: const TextStyle(color: EnterpriseColors.textPrimary),
                ),
                subtitle: Text(
                  '${note.type.toUpperCase()} • ${_formatDate(note.timestamp)}'
                  '${note.staff != null ? ' • ${note.staff}' : ''}',
                  style: const TextStyle(
                    color: EnterpriseColors.textSecondary,
                    fontSize: 12,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildOrdersSection(BuildContext context, Customer customer) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            const Text(
              '注文履歴',
              style: TextStyle(
                color: EnterpriseColors.textPrimary,
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            TextButton(
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (context) => OrderListScreen(
                      customerId: customer.id,
                    ),
                  ),
                );
              },
              child: const Text(
                'すべて見る',
                style: TextStyle(
                  color: EnterpriseColors.primaryBlue,
                ),
              ),
            ),
          ],
        ),
        const SizedBox(height: 16),
        // TODO: 注文履歴の実装（OrderListScreenと統合）
        Card(
          color: EnterpriseColors.surfaceGray,
          child: Padding(
            padding: const EdgeInsets.all(20),
            child: Center(
              child: Text(
                '注文履歴は注文管理画面から確認できます',
                style: const TextStyle(
                  color: EnterpriseColors.textSecondary,
                  fontSize: 14,
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }

  String _formatDate(DateTime date) {
    final local = date.toLocal();
    return '${local.year}年${local.month}月${local.day}日';
  }
}
