import 'package:flutter/material.dart';
import '../../config/enterprise_theme.dart';

/// Visual Ordering画面（視覚的な採寸入力）
class VisualOrderingScreen extends StatelessWidget {
  const VisualOrderingScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: EnterpriseColors.deepBlack,
      body: Column(
        children: [
          // ヘッダー
          _buildHeader(context),

          // メインコンテンツ
          Expanded(
            child: Row(
              children: [
                // 左: 人体図（採寸入力）
                Expanded(
                  flex: 2,
                  child: _buildBodyMap(context),
                ),

                // 右: 仕様選択・オプション
                Expanded(
                  child: _buildSpecifications(context),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildHeader(BuildContext context) {
    return Container(
      height: 64,
      padding: const EdgeInsets.symmetric(horizontal: 24),
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
          Row(
            children: [
              IconButton(
                onPressed: () => Navigator.pop(context),
                icon: const Icon(Icons.arrow_back),
                color: EnterpriseColors.textSecondary,
              ),
              const SizedBox(width: 16),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    'New Order',
                    style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                  ),
                  Text(
                    'Customer: K. Tanaka (ID: 8821)',
                    style: Theme.of(context).textTheme.labelSmall,
                  ),
                ],
              ),
            ],
          ),
          Row(
            children: [
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    'Total Estimate',
                    style: Theme.of(context).textTheme.labelSmall,
                  ),
                  Text(
                    '¥135,000',
                    style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                          color: EnterpriseColors.metallicGold,
                          fontWeight: FontWeight.bold,
                        ),
                  ),
                ],
              ),
              const SizedBox(width: 16),
              ElevatedButton(
                onPressed: () {
                  // TODO: 注文確定
                },
                child: const Text('Confirm Order'),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildBodyMap(BuildContext context) {
    return Container(
      decoration: const BoxDecoration(
        color: EnterpriseColors.regalisBlack,
        border: Border(
          right: BorderSide(
            color: EnterpriseColors.borderGray,
            width: 1,
          ),
        ),
      ),
      child: Stack(
        children: [
          // 人体図のプレースホルダー
          Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(
                  '3D MEASUREMENT ASSIST',
                  style: Theme.of(context).textTheme.labelSmall?.copyWith(
                        letterSpacing: 2,
                      ),
                ),
                const SizedBox(height: 48),
                Container(
                  width: 300,
                  height: 600,
                  decoration: BoxDecoration(
                    border: Border.all(
                      color: EnterpriseColors.borderGray.withOpacity(0.3),
                      width: 2,
                    ),
                    borderRadius: BorderRadius.circular(16),
                  ),
                  child: const Center(
                    child: Text(
                      'Body Map\n(Coming Soon)',
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        color: EnterpriseColors.textTertiary,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),

          // 採寸入力モーダル（プレースホルダー）
          Positioned(
            top: 140,
            right: 32,
            child: Container(
              width: 256,
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: EnterpriseColors.regalisBlack.withOpacity(0.9),
                borderRadius: BorderRadius.circular(12),
                border: Border.all(
                  color: EnterpriseColors.metallicGold,
                  width: 1,
                ),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        'CHEST (バスト)',
                        style: Theme.of(context).textTheme.labelSmall?.copyWith(
                              color: EnterpriseColors.metallicGold,
                              fontWeight: FontWeight.bold,
                            ),
                      ),
                      IconButton(
                        icon: const Icon(Icons.close, size: 16),
                        onPressed: () {},
                        padding: EdgeInsets.zero,
                        constraints: const BoxConstraints(),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: [
                      SizedBox(
                        width: 96,
                        child: TextField(
                          controller: TextEditingController(text: '98.5'),
                          style: Theme.of(context)
                              .textTheme
                              .displayLarge
                              ?.copyWith(
                                fontSize: 32,
                              ),
                          keyboardType: TextInputType.number,
                          textAlign: TextAlign.center,
                          decoration: const InputDecoration(
                            border: UnderlineInputBorder(
                              borderSide: BorderSide(
                                color: EnterpriseColors.borderGray,
                              ),
                            ),
                            enabledBorder: UnderlineInputBorder(
                              borderSide: BorderSide(
                                color: EnterpriseColors.borderGray,
                              ),
                            ),
                            focusedBorder: UnderlineInputBorder(
                              borderSide: BorderSide(
                                color: EnterpriseColors.metallicGold,
                              ),
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(width: 8),
                      Text(
                        'cm',
                        style: Theme.of(context).textTheme.bodySmall,
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  // スライダー
                  Slider(
                    value: 98.5,
                    min: 90,
                    max: 110,
                    divisions: 40,
                    label: '98.5 cm',
                    onChanged: (value) {},
                    activeColor: EnterpriseColors.metallicGold,
                  ),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        child: OutlinedButton(
                          onPressed: () {},
                          child: const Text('+0.5'),
                        ),
                      ),
                      const SizedBox(width: 8),
                      Expanded(
                        child: OutlinedButton(
                          onPressed: () {},
                          child: const Text('-0.5'),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSpecifications(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: const BoxDecoration(
        color: EnterpriseColors.surfaceGray,
        border: Border(
          left: BorderSide(
            color: EnterpriseColors.borderGray,
            width: 1,
          ),
        ),
      ),
      child: SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'SPECIFICATION',
              style: Theme.of(context).textTheme.labelSmall?.copyWith(
                    letterSpacing: 2,
                  ),
            ),
            const SizedBox(height: 32),

            // 選択された生地
            Card(
              child: Padding(
                padding: const EdgeInsets.all(12),
                child: Row(
                  children: [
                    Container(
                      width: 48,
                      height: 48,
                      decoration: BoxDecoration(
                        color: EnterpriseColors.borderGray,
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: const Icon(
                        Icons.image,
                        color: EnterpriseColors.textTertiary,
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            'Selected Fabric',
                            style: Theme.of(context)
                                .textTheme
                                .labelSmall
                                ?.copyWith(
                                  color: EnterpriseColors.metallicGold,
                                ),
                          ),
                          Text(
                            'V.B.C Perennial Navy',
                            style:
                                Theme.of(context).textTheme.bodyLarge?.copyWith(
                                      fontWeight: FontWeight.bold,
                                    ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),

            const SizedBox(height: 32),

            // モデル選択
            Text(
              'Model',
              style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                Expanded(
                  child: OutlinedButton(
                    onPressed: () {},
                    style: OutlinedButton.styleFrom(
                      side: const BorderSide(
                          color: EnterpriseColors.metallicGold),
                      backgroundColor: EnterpriseColors.surfaceGray,
                      foregroundColor: EnterpriseColors.metallicGold,
                    ),
                    child: const Text('British (Modern)'),
                  ),
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: OutlinedButton(
                    onPressed: () {},
                    style: OutlinedButton.styleFrom(
                      side:
                          const BorderSide(color: EnterpriseColors.borderGray),
                      foregroundColor: EnterpriseColors.textSecondary,
                    ),
                    child: const Text('Italian (Classico)'),
                  ),
                ),
              ],
            ),

            const SizedBox(height: 32),

            // ラペル選択
            Text(
              'Lapel',
              style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
            ),
            const SizedBox(height: 12),
            RadioListTile<String>(
              value: 'notch',
              groupValue: 'notch',
              onChanged: (value) {},
              title: const Text('Notch Lapel (8.5cm)'),
              activeColor: EnterpriseColors.metallicGold,
              contentPadding: EdgeInsets.zero,
            ),
            RadioListTile<String>(
              value: 'peak',
              groupValue: 'notch',
              onChanged: (value) {},
              title: const Text('Peak Lapel (9.5cm)'),
              activeColor: EnterpriseColors.metallicGold,
              contentPadding: EdgeInsets.zero,
            ),

            const SizedBox(height: 32),

            // コンプライアンス確認
            Divider(
              color: EnterpriseColors.borderGray,
              height: 1,
            ),
            const SizedBox(height: 24),
            CheckboxListTile(
              value: false,
              onChanged: (value) {},
              title: Text(
                '本発注により、下請法第3条に基づく書面が自動生成され、提携工場へ送信されることに同意します。',
                style: Theme.of(context).textTheme.bodySmall,
              ),
              activeColor: EnterpriseColors.metallicGold,
              contentPadding: EdgeInsets.zero,
            ),
          ],
        ),
      ),
    );
  }
}
