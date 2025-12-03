import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../config/enterprise_theme.dart';
import '../../models/customer.dart';
import '../../providers/customer_provider.dart';
import 'customer_detail_screen.dart';
import 'customer_create_screen.dart';

/// 顧客一覧画面
class CustomerListScreen extends ConsumerStatefulWidget {
  const CustomerListScreen({super.key});

  @override
  ConsumerState<CustomerListScreen> createState() => _CustomerListScreenState();
}

class _CustomerListScreenState extends ConsumerState<CustomerListScreen> {
  final TextEditingController _searchController = TextEditingController();
  String _searchQuery = '';

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

    final customersAsync = ref.watch(customerListProvider(tenantId));

    return Scaffold(
      backgroundColor: EnterpriseColors.deepBlack,
      appBar: AppBar(
        backgroundColor: EnterpriseColors.surfaceGray,
        title: const Text(
          '顧客一覧',
          style: TextStyle(
            color: EnterpriseColors.textPrimary,
            fontSize: 20,
            fontWeight: FontWeight.bold,
          ),
        ),
        elevation: 0,
      ),
      body: Column(
        children: [
          // 検索バー
          _buildSearchBar(),
          
          // 顧客リスト
          Expanded(
            child: customersAsync.when(
              data: (customers) {
                // 検索フィルター適用
                final filteredCustomers = _searchQuery.isEmpty
                    ? customers
                    : customers.where((customer) {
                        final query = _searchQuery.toLowerCase();
                        return customer.name.toLowerCase().contains(query) ||
                            (customer.email?.toLowerCase().contains(query) ?? false) ||
                            (customer.phone?.contains(_searchQuery) ?? false);
                      }).toList();

                if (filteredCustomers.isEmpty) {
                  return Center(
                    child: Text(
                      _searchQuery.isEmpty
                          ? '顧客が登録されていません'
                          : '検索結果が見つかりませんでした',
                      style: const TextStyle(
                        color: EnterpriseColors.textSecondary,
                        fontSize: 16,
                      ),
                    ),
                  );
                }

                return RefreshIndicator(
                  onRefresh: () async {
                    ref.invalidate(customerListProvider(tenantId));
                  },
                  color: EnterpriseColors.primaryBlue,
                  child: ListView.builder(
                    padding: const EdgeInsets.all(16),
                    itemCount: filteredCustomers.length,
                    itemBuilder: (context, index) {
                      final customer = filteredCustomers[index];
                      return _buildCustomerCard(context, customer);
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
                        ref.invalidate(customerListProvider(tenantId));
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
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          Navigator.push(
            context,
            MaterialPageRoute(
              builder: (context) => const CustomerCreateScreen(),
            ),
          ).then((_) {
            // 顧客作成後、リストを更新
            ref.invalidate(customerListProvider(tenantId));
          });
        },
        backgroundColor: EnterpriseColors.primaryBlue,
        child: const Icon(Icons.add, color: Colors.white),
      ),
    );
  }

  Widget _buildSearchBar() {
    return Container(
      padding: const EdgeInsets.all(16),
      color: EnterpriseColors.surfaceGray,
      child: TextField(
        controller: _searchController,
        style: const TextStyle(color: EnterpriseColors.textPrimary),
        decoration: InputDecoration(
          hintText: '顧客名、メール、電話番号で検索',
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
    );
  }

  Widget _buildCustomerCard(BuildContext context, Customer customer) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      color: EnterpriseColors.surfaceGray,
      child: InkWell(
        onTap: () {
          Navigator.push(
            context,
            MaterialPageRoute(
              builder: (context) => CustomerDetailScreen(
                customerId: customer.id,
              ),
            ),
          );
        },
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              // アバター
              CircleAvatar(
                radius: 24,
                backgroundColor: EnterpriseColors.primaryBlue,
                child: Text(
                  customer.name.isNotEmpty
                      ? customer.name[0].toUpperCase()
                      : '?',
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
              const SizedBox(width: 16),
              
              // 顧客情報
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      customer.name,
                      style: const TextStyle(
                        color: EnterpriseColors.textPrimary,
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    if (customer.email != null) ...[
                      const SizedBox(height: 4),
                      Text(
                        customer.email!,
                        style: const TextStyle(
                          color: EnterpriseColors.textSecondary,
                          fontSize: 14,
                        ),
                      ),
                    ],
                    if (customer.phone != null) ...[
                      const SizedBox(height: 4),
                      Text(
                        customer.phone!,
                        style: const TextStyle(
                          color: EnterpriseColors.textSecondary,
                          fontSize: 14,
                        ),
                      ),
                    ],
                  ],
                ),
              ),
              
              // 矢印アイコン
              const Icon(
                Icons.chevron_right,
                color: EnterpriseColors.textSecondary,
              ),
            ],
          ),
        ),
      ),
    );
  }
}

