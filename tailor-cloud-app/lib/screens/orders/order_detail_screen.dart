import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../config/enterprise_theme.dart';
import '../../models/order.dart';
import '../../providers/order_provider.dart';
import '../../providers/compliance_provider.dart';
import '../../providers/invoice_provider.dart';
import '../../services/pdf_download_service.dart';
import 'order_confirm_screen.dart';
import 'compliance_document_screen.dart';

/// 注文詳細画面
class OrderDetailScreen extends ConsumerWidget {
  final String orderId;
  final String tenantId;

  const OrderDetailScreen({
    super.key,
    required this.orderId,
    required this.tenantId,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final orderAsync = ref.watch(orderProvider(orderId, tenantId: tenantId));

    return Scaffold(
      backgroundColor: EnterpriseColors.deepBlack,
      appBar: AppBar(
        backgroundColor: EnterpriseColors.surfaceGray,
        title: const Text(
          '注文詳細',
          style: TextStyle(
            color: EnterpriseColors.textPrimary,
            fontSize: 20,
            fontWeight: FontWeight.bold,
          ),
        ),
        elevation: 0,
      ),
      body: orderAsync.when(
        data: (order) => SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // 注文情報カード
              _buildOrderInfoCard(context, order),
              
              const SizedBox(height: 24),
              
              // 注文明細セクション
              _buildOrderDetailsSection(context, order),
              
              const SizedBox(height: 24),
              
              // アクションボタン
              _buildActionButtons(context, ref, order),
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
                  ref.invalidate(orderProvider(orderId, tenantId: tenantId));
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

  Widget _buildOrderInfoCard(BuildContext context, Order order) {
    return Card(
      color: EnterpriseColors.surfaceGray,
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // ヘッダー（注文ID、ステータス）
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        '注文ID',
                        style: const TextStyle(
                          color: EnterpriseColors.textSecondary,
                          fontSize: 12,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        order.id,
                        style: const TextStyle(
                          color: EnterpriseColors.textPrimary,
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                ),
                _buildStatusBadge(order.status),
              ],
            ),
            
            const Divider(
              color: EnterpriseColors.borderGray,
              height: 32,
            ),
            
            // 顧客ID
            _buildInfoRow(
              icon: Icons.person,
              label: '顧客ID',
              value: order.customerId,
            ),
            
            const SizedBox(height: 16),
            
            // 生地ID
            _buildInfoRow(
              icon: Icons.inventory_2,
              label: '生地ID',
              value: order.fabricId,
            ),
            
            const SizedBox(height: 16),
            
            // 金額
            _buildInfoRow(
              icon: Icons.attach_money,
              label: '金額',
              value: order.amountDisplay,
            ),
            
            const SizedBox(height: 16),
            
            // 納期
            _buildInfoRow(
              icon: Icons.calendar_today,
              label: '納期',
              value: _formatDate(order.deliveryDate),
            ),
            
            const SizedBox(height: 16),
            
            // 支払期日
            _buildInfoRow(
              icon: Icons.payment,
              label: '支払期日',
              value: _formatDate(order.paymentDueDate),
            ),
            
            if (order.complianceDocUrl != null) ...[
              const SizedBox(height: 16),
              _buildInfoRow(
                icon: Icons.description,
                label: '発注書',
                value: '発行済み',
                isLink: true,
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
    bool isLink = false,
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
                style: TextStyle(
                  color: isLink
                      ? EnterpriseColors.primaryBlue
                      : EnterpriseColors.textPrimary,
                  fontSize: 16,
                  decoration: isLink ? TextDecoration.underline : null,
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildOrderDetailsSection(BuildContext context, Order order) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '注文明細',
          style: TextStyle(
            color: EnterpriseColors.textPrimary,
            fontSize: 20,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 16),
        Card(
          color: EnterpriseColors.surfaceGray,
          child: Padding(
            padding: const EdgeInsets.all(20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                if (order.details?.description != null) ...[
                  _buildDetailRow('説明', order.details!.description!),
                  const SizedBox(height: 16),
                ],
                if (order.details?.measurementData != null) ...[
                  _buildDetailRow(
                    '採寸データ',
                    order.details!.measurementData.toString(),
                  ),
                  const SizedBox(height: 16),
                ],
                if (order.details?.adjustments != null) ...[
                  _buildDetailRow(
                    '補正情報',
                    order.details!.adjustments.toString(),
                  ),
                ],
                if (order.details == null ||
                    (order.details?.description == null &&
                        order.details?.measurementData == null &&
                        order.details?.adjustments == null)) ...[
                  const Center(
                    child: Text(
                      '明細情報がありません',
                      style: TextStyle(
                        color: EnterpriseColors.textSecondary,
                        fontSize: 14,
                      ),
                    ),
                  ),
                ],
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildDetailRow(String label, String value) {
    return Column(
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
            fontSize: 14,
          ),
        ),
      ],
    );
  }

  Widget _buildActionButtons(
    BuildContext context,
    WidgetRef ref,
    Order order,
  ) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        // 発注書PDF生成ボタン（確定済みで未生成の場合）
        if (order.status == OrderStatus.confirmed &&
            order.complianceDocUrl == null) ...[
          ElevatedButton.icon(
            onPressed: () async {
              _generateComplianceDocument(context, ref, order);
            },
            icon: const Icon(Icons.description),
            label: const Text('発注書PDFを生成'),
            style: ElevatedButton.styleFrom(
              backgroundColor: EnterpriseColors.primaryBlue,
              padding: const EdgeInsets.symmetric(vertical: 16),
            ),
          ),
          const SizedBox(height: 12),
        ],
        
        // 発注書PDF表示ボタン（確定済みで生成済みの場合）
        if (order.status == OrderStatus.confirmed &&
            order.complianceDocUrl != null) ...[
          ElevatedButton.icon(
            onPressed: () {
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (context) => ComplianceDocumentScreen(
                    orderId: order.id,
                    documentUrl: order.complianceDocUrl!,
                  ),
                ),
              );
            },
            icon: const Icon(Icons.description),
            label: const Text('発注書PDFを表示'),
            style: ElevatedButton.styleFrom(
              backgroundColor: EnterpriseColors.primaryBlue,
              padding: const EdgeInsets.symmetric(vertical: 16),
            ),
          ),
          const SizedBox(height: 12),
        ],
        
        // インボイスPDF生成ボタン（確定済みの場合）
        if (order.status == OrderStatus.confirmed) ...[
          ElevatedButton.icon(
            onPressed: () async {
              _generateInvoice(context, ref, order);
            },
            icon: const Icon(Icons.receipt),
            label: const Text('インボイスPDFを生成'),
            style: ElevatedButton.styleFrom(
              backgroundColor: EnterpriseColors.metallicGold,
              padding: const EdgeInsets.symmetric(vertical: 16),
            ),
          ),
          const SizedBox(height: 12),
        ],
        
        // 注文確定ボタン（下書きの場合）
        if (order.status == OrderStatus.draft) ...[
          ElevatedButton.icon(
            onPressed: () {
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (context) => OrderConfirmScreen(
                    order: order,
                  ),
                ),
              ).then((_) {
                // 注文確定後、注文情報を更新
                ref.invalidate(orderProvider(order.id, tenantId: tenantId));
              });
            },
            icon: const Icon(Icons.check_circle),
            label: const Text('注文を確定'),
            style: ElevatedButton.styleFrom(
              backgroundColor: EnterpriseColors.successGreen,
              padding: const EdgeInsets.symmetric(vertical: 16),
            ),
          ),
          const SizedBox(height: 12),
        ],
        
        // 修正発注書生成ボタン（確定済みの場合）
        if (order.status == OrderStatus.confirmed) ...[
          OutlinedButton.icon(
            onPressed: () {
              // TODO: 修正発注書生成画面の実装
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Text('修正発注書生成機能は実装中です'),
                  backgroundColor: EnterpriseColors.textSecondary,
                ),
              );
            },
            icon: const Icon(Icons.edit),
            label: const Text('修正発注書を生成'),
            style: OutlinedButton.styleFrom(
              foregroundColor: EnterpriseColors.primaryBlue,
              side: const BorderSide(color: EnterpriseColors.primaryBlue),
              padding: const EdgeInsets.symmetric(vertical: 16),
            ),
          ),
        ],
      ],
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
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(4),
      ),
      child: Text(
        label,
        style: TextStyle(
          color: textColor,
          fontSize: 14,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }

  String _formatDate(DateTime date) {
    return '${date.year}年${date.month}月${date.day}日';
  }

  Future<void> _generateComplianceDocument(
    BuildContext context,
    WidgetRef ref,
    Order order,
  ) async {
    // ローディングダイアログを表示
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => const Center(
        child: CircularProgressIndicator(
          color: EnterpriseColors.primaryBlue,
        ),
      ),
    );

    try {
      final response = await ref.read(
        generateComplianceDocumentProvider(order.id).future,
      );

      if (context.mounted) {
        Navigator.pop(context); // ローディングダイアログを閉じる

        // 成功メッセージを表示
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('発注書PDFを生成しました'),
            backgroundColor: EnterpriseColors.successGreen,
          ),
        );

        // 注文情報を更新
        ref.invalidate(orderProvider(order.id, tenantId: tenantId));

        // PDF表示画面に遷移
        Navigator.push(
          context,
          MaterialPageRoute(
            builder: (context) => ComplianceDocumentScreen(
              orderId: order.id,
              documentUrl: response.docUrl,
            ),
          ),
        );
      }
    } catch (e) {
      if (context.mounted) {
        Navigator.pop(context); // ローディングダイアログを閉じる

        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('発注書PDFの生成に失敗しました: ${e.toString()}'),
            backgroundColor: EnterpriseColors.errorRed,
          ),
        );
      }
    }
  }

  Future<void> _generateInvoice(
    BuildContext context,
    WidgetRef ref,
    Order order,
  ) async {
    // ローディングダイアログを表示
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => const Center(
        child: CircularProgressIndicator(
          color: EnterpriseColors.primaryBlue,
        ),
      ),
    );

    try {
      final response = await ref.read(
        generateInvoiceProvider(order.id).future,
      );

      if (context.mounted) {
        Navigator.pop(context); // ローディングダイアログを閉じる

        // ダウンロード確認ダイアログ
        final shouldDownload = await showDialog<bool>(
          context: context,
          builder: (context) => AlertDialog(
            backgroundColor: EnterpriseColors.surfaceGray,
            title: const Text(
              'インボイスPDF生成完了',
              style: TextStyle(color: EnterpriseColors.textPrimary),
            ),
            content: const Text(
              'インボイスPDFを生成しました。ダウンロードしますか？',
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
              ElevatedButton(
                onPressed: () => Navigator.of(context).pop(true),
                style: ElevatedButton.styleFrom(
                  backgroundColor: EnterpriseColors.primaryBlue,
                ),
                child: const Text('ダウンロード'),
              ),
            ],
          ),
        );

        if (shouldDownload == true) {
          // PDFをダウンロード
          try {
            final fileName = 'invoice_${order.id}_${DateTime.now().millisecondsSinceEpoch}.pdf';
            final filePath = await PdfDownloadService.downloadPdf(
              url: response.invoiceUrl,
              fileName: fileName,
            );

            if (context.mounted) {
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(
                  content: Text('PDFをダウンロードしました: $filePath'),
                  backgroundColor: EnterpriseColors.successGreen,
                  duration: const Duration(seconds: 3),
                ),
              );
            }
          } catch (e) {
            if (context.mounted) {
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(
                  content: Text('PDFのダウンロードに失敗しました: ${e.toString()}'),
                  backgroundColor: EnterpriseColors.errorRed,
                ),
              );
            }
          }
        }
      }
    } catch (e) {
      if (context.mounted) {
        Navigator.pop(context); // ローディングダイアログを閉じる

        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('インボイスPDFの生成に失敗しました: ${e.toString()}'),
            backgroundColor: EnterpriseColors.errorRed,
          ),
        );
      }
    }
  }
}

