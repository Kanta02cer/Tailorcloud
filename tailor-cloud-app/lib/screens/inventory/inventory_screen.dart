import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../config/enterprise_theme.dart';
import '../../widgets/fabric_card.dart';
import '../../models/fabric.dart';
import '../../providers/fabric_provider.dart';

/// Inventory画面（生地一覧）
class InventoryScreen extends ConsumerStatefulWidget {
  const InventoryScreen({super.key});

  @override
  ConsumerState<InventoryScreen> createState() => _InventoryScreenState();
}

class _InventoryScreenState extends ConsumerState<InventoryScreen> {
  final TextEditingController _searchController = TextEditingController();
  String _selectedStatus = 'all';
  String? _searchQuery;

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    // TODO: 認証からテナントIDを取得
    const tenantId = 'tenant-123'; // 仮のテナントID

    final params = FabricListParams(
      tenantId: tenantId,
      status: _selectedStatus == 'all' ? null : _selectedStatus,
      search: _searchQuery?.isEmpty ?? true ? null : _searchQuery,
    );

    final fabricsAsync = ref.watch(fabricListProvider(params));

    return Scaffold(
      backgroundColor: EnterpriseColors.deepBlack,
      body: Column(
        children: [
          // 検索ヘッダー
          _buildSearchHeader(context),

          // グリッド
          Expanded(
            child: fabricsAsync.when(
              data: (fabrics) {
                if (fabrics.isEmpty) {
                  return Center(
                    child: Text(
                      '生地が見つかりませんでした',
                      style: Theme.of(context).textTheme.bodyLarge,
                    ),
                  );
                }

                return GridView.builder(
                  padding: const EdgeInsets.all(24),
                  gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                    crossAxisCount: 2,
                    childAspectRatio: 3 / 4,
                    crossAxisSpacing: 24,
                    mainAxisSpacing: 24,
                  ),
                  itemCount: fabrics.length,
                  itemBuilder: (context, index) {
                    final fabric = fabrics[index];
                    return FabricCard(
                      fabric: fabric,
                      isRecommended: index == 1, // 仮の推奨判定
                      onTap: () {
                        _showFabricDetail(context, fabric);
                      },
                    );
                  },
                );
              },
              loading: () => const Center(
                child: CircularProgressIndicator(
                  color: EnterpriseColors.metallicGold,
                ),
              ),
              error: (error, stack) => Center(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Icon(
                      Icons.error_outline,
                      size: 48,
                      color: EnterpriseColors.statusOutOfStock,
                    ),
                    const SizedBox(height: 16),
                    Text(
                      'エラーが発生しました',
                      style: Theme.of(context).textTheme.bodyLarge,
                    ),
                    const SizedBox(height: 8),
                    Text(
                      error.toString(),
                      style: Theme.of(context).textTheme.bodySmall,
                      textAlign: TextAlign.center,
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

  Widget _buildSearchHeader(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: const BoxDecoration(
        color: EnterpriseColors.surfaceGray,
        border: Border(
          bottom: BorderSide(
            color: EnterpriseColors.borderGray,
            width: 1,
          ),
        ),
      ),
      child: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                'Fabric Library',
                style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
              ),
              IconButton(
                onPressed: () {
                  // TODO: スキャン機能
                },
                icon: const Icon(Icons.qr_code_scanner),
                color: EnterpriseColors.textSecondary,
              ),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              // 検索バー
              Expanded(
                child: TextField(
                  controller: _searchController,
                  decoration: InputDecoration(
                    hintText: 'Search by SKU, Color, Brand...',
                    prefixIcon: const Icon(Icons.search),
                    suffixIcon: _searchController.text.isNotEmpty
                        ? IconButton(
                            icon: const Icon(Icons.clear),
                            onPressed: () {
                              _searchController.clear();
                              setState(() {
                                _searchQuery = null;
                              });
                            },
                          )
                        : null,
                  ),
                  style: const TextStyle(color: EnterpriseColors.textPrimary),
                  onSubmitted: (value) {
                    setState(() {
                      _searchQuery = value.isEmpty ? null : value;
                    });
                  },
                ),
              ),
              const SizedBox(width: 16),

              // フィルターボタン
              SingleChildScrollView(
                scrollDirection: Axis.horizontal,
                child: Row(
                  children: [
                    _buildFilterChip('All', 'all'),
                    const SizedBox(width: 8),
                    _buildFilterChip('Available', 'available'),
                    const SizedBox(width: 8),
                    _buildFilterChip('Limited', 'limited'),
                    const SizedBox(width: 8),
                    _buildFilterChip('Out of Stock', 'soldout'),
                  ],
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildFilterChip(String label, String value) {
    final isSelected = _selectedStatus == value;
    return FilterChip(
      label: Text(label),
      selected: isSelected,
      onSelected: (selected) {
        if (selected) {
          setState(() {
            _selectedStatus = value;
          });
        }
      },
      selectedColor: EnterpriseColors.metallicGold,
      labelStyle: TextStyle(
        color: isSelected ? Colors.black : EnterpriseColors.textSecondary,
        fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
      ),
      backgroundColor: EnterpriseColors.surfaceGray,
      side: BorderSide(
        color: isSelected
            ? EnterpriseColors.metallicGold
            : EnterpriseColors.borderGray,
      ),
    );
  }

  void _showFabricDetail(BuildContext context, Fabric fabric) {
    showModalBottomSheet(
      context: context,
      backgroundColor: EnterpriseColors.regalisBlack,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => Container(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // ヘッダー
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        fabric.name,
                        style: Theme.of(context).textTheme.headlineMedium,
                      ),
                      const SizedBox(height: 8),
                      Text(
                        '在庫: ${fabric.stockAmount?.toStringAsFixed(1) ?? "0.0"}m',
                        style: Theme.of(context).textTheme.bodySmall,
                      ),
                    ],
                  ),
                ),
                IconButton(
                  onPressed: () => Navigator.pop(context),
                  icon: const Icon(Icons.close),
                ),
              ],
            ),

            const SizedBox(height: 24),

            // アクションボタン
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: () {
                  // TODO: 生地確保機能
                  Navigator.pop(context);
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(
                      content: Text('生地を確保しました'),
                    ),
                  );
                },
                child: const Text('この生地を確保して発注'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
