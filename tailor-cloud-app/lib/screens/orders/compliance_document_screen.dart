import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';
import '../../config/enterprise_theme.dart';

/// 発注書PDF表示画面
class ComplianceDocumentScreen extends StatelessWidget {
  final String orderId;
  final String documentUrl;

  const ComplianceDocumentScreen({
    super.key,
    required this.orderId,
    required this.documentUrl,
  });

  Future<void> _openPdf(BuildContext context) async {
    final uri = Uri.parse(documentUrl);
    if (await canLaunchUrl(uri)) {
      await launchUrl(
        uri,
        mode: LaunchMode.externalApplication,
      );
    } else {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('PDFを開けませんでした'),
            backgroundColor: EnterpriseColors.errorRed,
          ),
        );
      }
    }
  }

  Future<void> _downloadPdf(BuildContext context) async {
    // TODO: PDFダウンロード機能の実装
    // 現在は外部ブラウザで開く
    await _openPdf(context);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: EnterpriseColors.deepBlack,
      appBar: AppBar(
        backgroundColor: EnterpriseColors.surfaceGray,
        title: const Text(
          '発注書PDF',
          style: TextStyle(
            color: EnterpriseColors.textPrimary,
            fontSize: 20,
            fontWeight: FontWeight.bold,
          ),
        ),
        elevation: 0,
        actions: [
          IconButton(
            icon: const Icon(
              Icons.download,
              color: EnterpriseColors.primaryBlue,
            ),
            onPressed: () => _downloadPdf(context),
            tooltip: 'ダウンロード',
          ),
        ],
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(
              Icons.description,
              size: 80,
              color: EnterpriseColors.primaryBlue,
            ),
            const SizedBox(height: 24),
            const Text(
              '発注書PDF',
              style: TextStyle(
                color: EnterpriseColors.textPrimary,
                fontSize: 24,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              '注文ID: $orderId',
              style: const TextStyle(
                color: EnterpriseColors.textSecondary,
                fontSize: 14,
              ),
            ),
            const SizedBox(height: 32),
            ElevatedButton.icon(
              onPressed: () => _openPdf(context),
              icon: const Icon(Icons.open_in_new),
              label: const Text('PDFを開く'),
              style: ElevatedButton.styleFrom(
                backgroundColor: EnterpriseColors.primaryBlue,
                padding: const EdgeInsets.symmetric(
                  horizontal: 24,
                  vertical: 16,
                ),
              ),
            ),
            const SizedBox(height: 16),
            OutlinedButton.icon(
              onPressed: () => _downloadPdf(context),
              icon: const Icon(Icons.download),
              label: const Text('ダウンロード'),
              style: OutlinedButton.styleFrom(
                foregroundColor: EnterpriseColors.primaryBlue,
                side: const BorderSide(color: EnterpriseColors.primaryBlue),
                padding: const EdgeInsets.symmetric(
                  horizontal: 24,
                  vertical: 16,
                ),
              ),
            ),
            const SizedBox(height: 32),
            Card(
              color: EnterpriseColors.surfaceGray,
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      '注意事項',
                      style: TextStyle(
                        color: EnterpriseColors.textPrimary,
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 8),
                    const Text(
                      '• この発注書は下請法3条書面として法的効力を持ちます\n'
                      '• PDFは外部ブラウザで開きます\n'
                      '• 発注書は必ず保存してください',
                      style: TextStyle(
                        color: EnterpriseColors.textSecondary,
                        fontSize: 12,
                        height: 1.5,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

