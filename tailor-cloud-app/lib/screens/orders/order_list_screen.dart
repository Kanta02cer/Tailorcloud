import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../config/enterprise_theme.dart';
import '../../models/order.dart';
import '../../providers/order_provider.dart';
import 'order_detail_screen.dart';

/// 注文一覧画面
class OrderListScreen extends ConsumerStatefulWidget {
  final String? customerId; // 顧客IDでフィルター（オプション）

  const OrderListScreen({
    super.key,
    this.customerId,
  });

  @override
  ConsumerState<OrderListScreen> createState() => _OrderListScreenState();
}

class _OrderListScreenState extends ConsumerState<OrderListScreen> {
  final TextEditingController _searchController = TextEditingController();
  String _searchQuery = '';
  String? _selectedStatus; // null = すべて

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    // TODO: 認証からテナントIDを取得（JWTトークンのカスタムクレームから取得）
    // 現在は開発用の仮のテナントIDを使用
    const tenantId = 'tenant-123'; // 仮のテナントID

    final ordersAsync = ref.watch(orderListProvider(tenantId));

    return Scaffold(
      backgroundColor: EnterpriseColors.deepBlack,
      appBar: AppBar(
        backgroundColor: EnterpriseColors.surfaceGray,
        title: Text(
          widget.customerId != null ? '顧客の注文履歴' : '注文一覧',
          style: const TextStyle(
            color: EnterpriseColors.textPrimary,
            fontSize: 20,
            fontWeight: FontWeight.bold,
          ),
        ),
        elevation: 0,
      ),
      body: Column(
        children: [
          // フィルター・検索バー
          _buildFilterBar(),

          // 注文リスト
          Expanded(
            child: ordersAsync.when(
              data: (orders) {
                // フィルター適用
                var filteredOrders = orders;

                // 顧客IDフィルター
                if (widget.customerId != null) {
                  filteredOrders = filteredOrders
                      .where((order) => order.customerId == widget.customerId)
                      .toList();
                }

                // ステータスフィルター
                if (_selectedStatus != null) {
                  filteredOrders = filteredOrders
                      .where((order) => order.status.name == _selectedStatus)
                      .toList();
                }

                // 検索フィルター
                if (_searchQuery.isNotEmpty) {
                  filteredOrders = filteredOrders.where((order) {
                    // 注文ID、顧客ID、金額などで検索
                    final query = _searchQuery.toLowerCase();
                    return order.id.toLowerCase().contains(query) ||
                        order.customerId.toLowerCase().contains(query) ||
                        order.totalAmount.toString().contains(query);
                  }).toList();
                }

                if (filteredOrders.isEmpty) {
                  return Center(
                    child: Text(
                      _searchQuery.isNotEmpty || _selectedStatus != null
                          ? '検索結果が見つかりませんでした'
                          : '注文が登録されていません',
                      style: const TextStyle(
                        color: EnterpriseColors.textSecondary,
                        fontSize: 16,
                      ),
                    ),
                  );
                }

                return RefreshIndicator(
                  onRefresh: () async {
                    ref.invalidate(orderListProvider(tenantId));
                  },
                  color: EnterpriseColors.primaryBlue,
                  child: ListView.builder(
                    padding: const EdgeInsets.all(16),
                    itemCount: filteredOrders.length,
                    itemBuilder: (context, index) {
                      final order = filteredOrders[index];
                      return _buildOrderCard(context, order);
                    },
                  ),
                );
              },
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
                    const Text(
                      'エラーが発生しました',
                      style: TextStyle(
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
                        ref.invalidate(orderListProvider(tenantId));
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
          ),
        ],
      ),
    );
  }

  Widget _buildFilterBar() {
    return Container(
      padding: const EdgeInsets.all(16),
      color: EnterpriseColors.surfaceGray,
      child: Column(
        children: [
          // 検索バー
          TextField(
            controller: _searchController,
            style: const TextStyle(color: EnterpriseColors.textPrimary),
            decoration: InputDecoration(
              hintText: '注文ID、顧客ID、金額で検索',
              hintStyle: const TextStyle(
                color: EnterpriseColors.textSecondary,
              ),
              prefixIcon: const Icon(
                Icons.search,
                color: EnterpriseColors.textSecondary,
              ),
              suffixIcon: _searchQuery.isNotEmpty
                  ? IconButton(
                      icon: const Icon(
                        Icons.clear,
                        color: EnterpriseColors.textSecondary,
                      ),
                      onPressed: () {
                        _searchController.clear();
                        setState(() {
                          _searchQuery = '';
                        });
                      },
                    )
                  : null,
              filled: true,
              fillColor: EnterpriseColors.deepBlack,
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(8),
                borderSide: BorderSide(
                  color: EnterpriseColors.borderGray,
                ),
              ),
              enabledBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(8),
                borderSide: BorderSide(
                  color: EnterpriseColors.borderGray,
                ),
              ),
              focusedBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(8),
                borderSide: const BorderSide(
                  color: EnterpriseColors.primaryBlue,
                  width: 2,
                ),
              ),
            ),
            onChanged: (value) {
              setState(() {
                _searchQuery = value;
              });
            },
          ),

          const SizedBox(height: 12),

          // ステータスフィルター
          SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            child: Row(
              children: [
                _buildStatusChip('すべて', null),
                const SizedBox(width: 8),
                _buildStatusChip('下書き', 'draft'),
                const SizedBox(width: 8),
                _buildStatusChip('確定', 'confirmed'),
                const SizedBox(width: 8),
                _buildStatusChip('製作中', 'inProduction'),
                const SizedBox(width: 8),
                _buildStatusChip('完了', 'completed'),
                const SizedBox(width: 8),
                _buildStatusChip('キャンセル', 'cancelled'),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatusChip(String label, String? status) {
    final isSelected = _selectedStatus == status;
    return FilterChip(
      label: Text(label),
      selected: isSelected,
      onSelected: (selected) {
        setState(() {
          _selectedStatus = selected ? status : null;
        });
      },
      selectedColor: EnterpriseColors.primaryBlue,
      checkmarkColor: Colors.white,
      labelStyle: TextStyle(
        color: isSelected ? Colors.white : EnterpriseColors.textPrimary,
        fontSize: 12,
      ),
      backgroundColor: EnterpriseColors.deepBlack,
      side: BorderSide(
        color: isSelected
            ? EnterpriseColors.primaryBlue
            : EnterpriseColors.borderGray,
      ),
    );
  }

  Widget _buildOrderCard(BuildContext context, Order order) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      color: EnterpriseColors.surfaceGray,
      child: InkWell(
        onTap: () {
          Navigator.push(
            context,
            MaterialPageRoute(
              builder: (context) => OrderDetailScreen(
                orderId: order.id,
                tenantId: order.tenantId,
              ),
            ),
          );
        },
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // ヘッダー（注文ID、ステータス）
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Expanded(
                    child: Text(
                      '注文ID: ${order.id}',
                      style: const TextStyle(
                        color: EnterpriseColors.textPrimary,
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  _buildStatusBadge(order.status),
                ],
              ),

              const SizedBox(height: 12),

              // 顧客ID
              Text(
                '顧客ID: ${order.customerId}',
                style: const TextStyle(
                  color: EnterpriseColors.textSecondary,
                  fontSize: 14,
                ),
              ),

              const SizedBox(height: 8),

              // 金額
              Text(
                order.amountDisplay,
                style: const TextStyle(
                  color: EnterpriseColors.textPrimary,
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                ),
              ),

              const SizedBox(height: 8),

              // 納期・支払期日
              Row(
                children: [
                  const Icon(
                    Icons.calendar_today,
                    size: 16,
                    color: EnterpriseColors.textSecondary,
                  ),
                  const SizedBox(width: 4),
                  Text(
                    '納期: ${_formatDate(order.deliveryDate)}',
                    style: const TextStyle(
                      color: EnterpriseColors.textSecondary,
                      fontSize: 12,
                    ),
                  ),
                  const SizedBox(width: 16),
                  Text(
                    '支払期日: ${_formatDate(order.paymentDueDate)}',
                    style: const TextStyle(
                      color: EnterpriseColors.textSecondary,
                      fontSize: 12,
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildStatusBadge(OrderStatus status) {
    Color backgroundColor;
    Color textColor;
    String label;

    switch (status) {
      case OrderStatus.draft:
        backgroundColor = EnterpriseColors.textSecondary;
        textColor = Colors.white;
        label = '下書き';
        break;
      case OrderStatus.confirmed:
        backgroundColor = EnterpriseColors.primaryBlue;
        textColor = Colors.white;
        label = '確定';
        break;
      case OrderStatus.inProduction:
        backgroundColor = EnterpriseColors.statusLowStock;
        textColor = Colors.white;
        label = '製作中';
        break;
      case OrderStatus.completed:
        backgroundColor = EnterpriseColors.successGreen;
        textColor = Colors.white;
        label = '完了';
        break;
      case OrderStatus.cancelled:
        backgroundColor = EnterpriseColors.errorRed;
        textColor = Colors.white;
        label = 'キャンセル';
        break;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(4),
      ),
      child: Text(
        label,
        style: TextStyle(
          color: textColor,
          fontSize: 12,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }

  String _formatDate(DateTime date) {
    return '${date.year}/${date.month}/${date.day}';
  }
}
