import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../config/enterprise_theme.dart';
import '../../models/order.dart';
import '../../providers/order_provider.dart';
import '../../providers/compliance_provider.dart';
import 'compliance_document_screen.dart';

/// 注文確認画面
class OrderConfirmScreen extends ConsumerStatefulWidget {
  final Order order;

  const OrderConfirmScreen({
    super.key,
    required this.order,
  });

  @override
  ConsumerState<OrderConfirmScreen> createState() => _OrderConfirmScreenState();
}

class _OrderConfirmScreenState extends ConsumerState<OrderConfirmScreen> {
  final _principalNameController = TextEditingController(
    text: 'Regalis Societas', // デフォルト値
  );
  bool _isLoading = false;

  @override
  void dispose() {
    _principalNameController.dispose();
    super.dispose();
  }

  Future<void> _confirmOrder() async {
    if (_principalNameController.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('委託をする者の氏名を入力してください'),
          backgroundColor: EnterpriseColors.errorRed,
        ),
      );
      return;
    }

    setState(() {
      _isLoading = true;
    });

    try {
      // 注文確定APIを呼び出し
      final confirmRequest = ConfirmOrderRequest(
        orderId: widget.order.id,
        principalName: _principalNameController.text.trim(),
      );

      final updatedOrder = await ref.read(
        confirmOrderProvider(confirmRequest).future,
      );

      // 発注書PDFを生成
      final docResponse = await ref.read(
        generateComplianceDocumentProvider(updatedOrder.id).future,
      );

      if (mounted) {
        // 発注書PDF表示画面に遷移
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(
            builder: (context) => ComplianceDocumentScreen(
              orderId: updatedOrder.id,
              documentUrl: docResponse.docUrl,
            ),
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('エラー: ${e.toString()}'),
            backgroundColor: EnterpriseColors.errorRed,
          ),
        );
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: EnterpriseColors.deepBlack,
      appBar: AppBar(
        backgroundColor: EnterpriseColors.surfaceGray,
        title: const Text(
          '注文確認',
          style: TextStyle(
            color: EnterpriseColors.textPrimary,
            fontSize: 20,
            fontWeight: FontWeight.bold,
          ),
        ),
        elevation: 0,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 注文情報カード
            _buildOrderInfoCard(),

            const SizedBox(height: 24),

            // 委託をする者の氏名入力
            _buildPrincipalNameInput(),

            const SizedBox(height: 32),

            // 注意事項
            _buildNoticeCard(),

            const SizedBox(height: 32),

            // 確定ボタン
            ElevatedButton(
              onPressed: _isLoading ? null : _confirmOrder,
              style: ElevatedButton.styleFrom(
                backgroundColor: EnterpriseColors.successGreen,
                padding: const EdgeInsets.symmetric(vertical: 16),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
              ),
              child: _isLoading
                  ? const SizedBox(
                      height: 20,
                      width: 20,
                      child: CircularProgressIndicator(
                        strokeWidth: 2,
                        color: Colors.white,
                      ),
                    )
                  : const Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(Icons.check_circle, color: Colors.white),
                        SizedBox(width: 8),
                        Text(
                          '注文を確定する',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ],
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildOrderInfoCard() {
    return Card(
      color: EnterpriseColors.surfaceGray,
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              '注文情報',
              style: TextStyle(
                color: EnterpriseColors.textPrimary,
                fontSize: 18,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16),
            _buildInfoRow('注文ID', widget.order.id),
            const SizedBox(height: 12),
            _buildInfoRow('顧客ID', widget.order.customerId),
            const SizedBox(height: 12),
            _buildInfoRow('生地ID', widget.order.fabricId),
            const SizedBox(height: 12),
            _buildInfoRow('金額', widget.order.amountDisplay),
            const SizedBox(height: 12),
            _buildInfoRow('納期', _formatDate(widget.order.deliveryDate)),
            const SizedBox(height: 12),
            _buildInfoRow('支払期日', _formatDate(widget.order.paymentDueDate)),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoRow(String label, String value) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          label,
          style: const TextStyle(
            color: EnterpriseColors.textSecondary,
            fontSize: 14,
          ),
        ),
        Text(
          value,
          style: const TextStyle(
            color: EnterpriseColors.textPrimary,
            fontSize: 14,
            fontWeight: FontWeight.bold,
          ),
        ),
      ],
    );
  }

  Widget _buildPrincipalNameInput() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Row(
          children: [
            Icon(
              Icons.business,
              color: EnterpriseColors.primaryBlue,
              size: 20,
            ),
            SizedBox(width: 8),
            Text(
              '委託をする者の氏名',
              style: TextStyle(
                color: EnterpriseColors.textPrimary,
                fontSize: 14,
                fontWeight: FontWeight.bold,
              ),
            ),
            SizedBox(width: 4),
            Text(
              '*',
              style: TextStyle(
                color: EnterpriseColors.errorRed,
                fontSize: 14,
              ),
            ),
          ],
        ),
        const SizedBox(height: 8),
        TextField(
          controller: _principalNameController,
          style: const TextStyle(color: EnterpriseColors.textPrimary),
          decoration: InputDecoration(
            hintText: '例: Regalis Societas',
            hintStyle: const TextStyle(
              color: EnterpriseColors.textSecondary,
            ),
            filled: true,
            fillColor: EnterpriseColors.surfaceGray,
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
        ),
        const SizedBox(height: 8),
        const Text(
          '下請法3条書面に記載される委託をする者の氏名です',
          style: TextStyle(
            color: EnterpriseColors.textSecondary,
            fontSize: 12,
          ),
        ),
      ],
    );
  }

  Widget _buildNoticeCard() {
    return Card(
      color: EnterpriseColors.surfaceGray,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Row(
              children: [
                Icon(
                  Icons.info_outline,
                  color: EnterpriseColors.primaryBlue,
                  size: 20,
                ),
                SizedBox(width: 8),
                Text(
                  '注意事項',
                  style: TextStyle(
                    color: EnterpriseColors.textPrimary,
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            const Text(
              '• 注文を確定すると、法的拘束力が発生します\n'
              '• 下請法3条書面（発注書）が自動生成されます\n'
              '• 発注書はPDF形式でダウンロードできます\n'
              '• 注文確定後は、修正発注書の生成が必要です',
              style: TextStyle(
                color: EnterpriseColors.textSecondary,
                fontSize: 14,
                height: 1.5,
              ),
            ),
          ],
        ),
      ),
    );
  }

  String _formatDate(DateTime date) {
    return '${date.year}年${date.month}月${date.day}日';
  }
}
