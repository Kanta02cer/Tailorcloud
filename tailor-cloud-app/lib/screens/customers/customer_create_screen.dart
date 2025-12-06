import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../config/enterprise_theme.dart';
import '../../models/customer.dart';
import '../../providers/customer_provider.dart';
import '../../utils/responsive.dart';

/// 新規顧客登録画面
class CustomerCreateScreen extends ConsumerStatefulWidget {
  const CustomerCreateScreen({super.key});

  @override
  ConsumerState<CustomerCreateScreen> createState() =>
      _CustomerCreateScreenState();
}

class _CustomerCreateScreenState extends ConsumerState<CustomerCreateScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _nameKanaController = TextEditingController();
  final _emailController = TextEditingController();
  final _phoneController = TextEditingController();
  final _addressController = TextEditingController();
  final _vipRankController = TextEditingController();
  final _lifetimeValueController = TextEditingController();
  final _notesController = TextEditingController();
  final _tagInputController = TextEditingController();
  final _preferredChannelController = TextEditingController();
  final _leadSourceController = TextEditingController();
  final _interactionNoteController = TextEditingController();
  
  bool _isLoading = false;
  String _status = 'lead';
  String _interactionType = 'note';
  DateTime? _lastInteractionAt;
  final List<String> _tags = [];

  @override
  void dispose() {
    _nameController.dispose();
    _nameKanaController.dispose();
    _emailController.dispose();
    _phoneController.dispose();
    _addressController.dispose();
    _vipRankController.dispose();
    _lifetimeValueController.dispose();
    _notesController.dispose();
    _tagInputController.dispose();
    _preferredChannelController.dispose();
    _leadSourceController.dispose();
    _interactionNoteController.dispose();
    super.dispose();
  }

  Future<void> _saveCustomer() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    setState(() {
      _isLoading = true;
    });

    try {
      // TODO: 認証からテナントIDを取得（JWTトークンのカスタムクレームから取得）
      // 現在は開発用の仮のテナントIDを使用
      // const tenantId = 'tenant-123'; // 仮のテナントID（未使用のためコメントアウト）

      final request = CreateCustomerRequest(
        name: _nameController.text.trim(),
        nameKana: _nameKanaController.text.trim().isEmpty
            ? null
            : _nameKanaController.text.trim(),
        email: _emailController.text.trim().isEmpty
            ? null
            : _emailController.text.trim(),
        phone: _phoneController.text.trim().isEmpty
            ? null
            : _phoneController.text.trim(),
        address: _addressController.text.trim().isEmpty
            ? null
            : _addressController.text.trim(),
        status: _status,
        tags: _tags.isEmpty ? null : _tags,
        vipRank: int.tryParse(_vipRankController.text.trim()),
        lifetimeValue: double.tryParse(
          _lifetimeValueController.text.trim(),
        ),
        ltvScore: double.tryParse(_lifetimeValueController.text.trim()),
        preferredChannel: _preferredChannelController.text.trim().isEmpty
            ? null
            : _preferredChannelController.text.trim(),
        leadSource: _leadSourceController.text.trim().isEmpty
            ? null
            : _leadSourceController.text.trim(),
        notes: _notesController.text.trim().isEmpty
            ? null
            : _notesController.text.trim(),
        lastInteractionAt: _lastInteractionAt,
        interactions: _composeInteractionPayload(),
      );

      await ref.read(createCustomerProvider(request).future);

      if (!mounted) return;

        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('顧客を登録しました'),
            backgroundColor: EnterpriseColors.successGreen,
          ),
        );
        Navigator.pop(context, true);
    } catch (e) {
      if (!mounted) return;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('エラー: ${e.toString()}'),
            backgroundColor: EnterpriseColors.errorRed,
          ),
        );
    } finally {
      if (!mounted) return;
        setState(() {
          _isLoading = false;
        });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: EnterpriseColors.deepBlack,
      appBar: AppBar(
        backgroundColor: EnterpriseColors.surfaceGray,
        title: const Text(
          '新規顧客登録',
          style: TextStyle(
            color: EnterpriseColors.textPrimary,
            fontSize: 20,
            fontWeight: FontWeight.bold,
          ),
        ),
        elevation: 0,
      ),
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Center(
            child: ConstrainedBox(
              constraints: BoxConstraints(
                maxWidth: Responsive.formMaxWidth(context),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  // 顧客名（必須）
                  _buildTextField(
                    controller: _nameController,
                    label: '顧客名',
                    hintText: '田中 太郎',
                    icon: Icons.person,
                    required: true,
                    validator: (value) {
                      if (value == null || value.trim().isEmpty) {
                        return '顧客名を入力してください';
                      }
                      return null;
                    },
                  ),
                  const SizedBox(height: 16),
                  // カナ名（任意）
                  _akibuildTextField(
                    controller: _nameKanaController,
                    label: 'カナ名',
                    hintText: 'タナカ タロウ',
                    icon: Icons.text_fields,
                    required: false,
                  ),
                  const SizedBox(height: 16),
                  // メールアドレス（任意）
                  _buildTextField(
                    controller: _emailController,
                    label: 'メールアドレス',
                    hintText: 'tanaka@example.com',
                    icon: Icons.email,
                    keyboardType: TextInputType.emailAddress,
                    required: false,
                    validator: (value) {
                      if (value != null && value.trim().isNotEmpty) {
                        if (!value.contains('@')) {
                          return '有効なメールアドレスを入力してください';
                        }
                      }
                      return null;
                    },
                  ),
                  const SizedBox(height: 16),
                  // 電話番号（任意）
                  _buildTextField(
                    controller: _phoneController,
                    label: '電話番号',
                    hintText: '090-1234-5678',
                    icon: Icons.phone,
                    keyboardType: TextInputType.phone,
                    required: false,
                  ),
                  const SizedBox(height: 16),
                  // 住所（任意）
                  _buildTextField(
                    controller: _addressController,
                    label: '住所',
                    hintText: '東京都渋谷区...',
                    icon: Icons.location_on,
                    maxLines: 3,
                    required: false,
                  ),
                  const SizedBox(height: 16),
                  _buildStatusSelector(),
                  const SizedBox(height: 16),
                  _buildTagEditor(),
                  const SizedBox(height: 16),
                  _buildMetricFields(),
                  const SizedBox(height: 16),
                  _buildChannelFields(),
                  const SizedBox(height: 16),
                  _buildInteractionSection(),
                  const SizedBox(height: 16),
                  _buildNotesField(),
                  const SizedBox(height: 32),
                  // 保存ボタン
                  ElevatedButton(
                    onPressed: _isLoading ? null : _saveCustomer,
                    style: ElevatedButton.styleFrom(
                      backgroundColor: EnterpriseColors.primaryBlue,
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
                        : const Text(
                            '保存',
                            style: TextStyle(
                              color: Colors.white,
                              fontSize: 16,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildTextField({
    required TextEditingController controller,
    required String label,
    required String hintText,
    required IconData icon,
    bool required = false,
    TextInputType? keyboardType,
    int maxLines = 1,
    String? Function(String?)? validator,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            Icon(
              icon,
              color: EnterpriseColors.primaryBlue,
              size: 20,
            ),
            const SizedBox(width: 8),
            Text(
              label,
              style: const TextStyle(
                color: EnterpriseColors.textPrimary,
                fontSize: 14,
                fontWeight: FontWeight.bold,
              ),
            ),
            if (required) ...[
              const SizedBox(width: 4),
              const Text(
                '*',
                style: TextStyle(
                  color: EnterpriseColors.errorRed,
                  fontSize: 14,
                ),
              ),
            ],
          ],
        ),
        const SizedBox(height: 8),
        TextFormField(
          controller: controller,
          keyboardType: keyboardType,
          maxLines: maxLines,
          style: const TextStyle(color: EnterpriseColors.textPrimary),
          decoration: InputDecoration(
            hintText: hintText,
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
            errorBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(8),
              borderSide: const BorderSide(
                color: EnterpriseColors.errorRed,
              ),
            ),
            focusedErrorBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(8),
              borderSide: const BorderSide(
                color: EnterpriseColors.errorRed,
                width: 2,
              ),
            ),
          ),
          validator: validator,
        ),
      ],
    );
  }

  Widget _buildStatusSelector() {
    const statuses = ['lead', 'prospect', 'active', 'inactive', 'vip'];
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '顧客ステータス',
          style: TextStyle(
            color: EnterpriseColors.textPrimary,
            fontSize: 14,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 8),
        Wrap(
          spacing: 8,
          children: statuses.map((status) {
            final isSelected = _status == status;
            return ChoiceChip(
              label: Text(status.toUpperCase()),
              selected: isSelected,
              onSelected: (_) {
                setState(() => _status = status);
              },
              selectedColor: EnterpriseColors.primaryBlue,
              labelStyle: TextStyle(
                color: isSelected ? Colors.black : EnterpriseColors.textPrimary,
                fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
              ),
            );
          }).toList(),
        ),
      ],
    );
  }

  Widget _buildTagEditor() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'タグ',
          style: TextStyle(
            color: EnterpriseColors.textPrimary,
            fontSize: 14,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 8),
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children: _tags
              .map(
                (tag) => Chip(
                  label: Text(tag),
                  backgroundColor: EnterpriseColors.surfaceGray,
                  deleteIconColor: EnterpriseColors.textSecondary,
                  onDeleted: () => setState(() => _tags.remove(tag)),
                ),
              )
              .toList(),
        ),
        const SizedBox(height: 8),
        TextField(
          controller: _tagInputController,
          style: const TextStyle(color: EnterpriseColors.textPrimary),
          decoration: InputDecoration(
            hintText: 'タグを入力してEnterで追加',
            filled: true,
            fillColor: EnterpriseColors.surfaceGray,
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(8),
              borderSide: BorderSide(color: EnterpriseColors.borderGray),
            ),
          ),
          onSubmitted: (value) {
            final tag = value.trim().toLowerCase();
            if (tag.isEmpty) return;
            if (!_tags.contains(tag)) {
              setState(() => _tags.add(tag));
            }
            _tagInputController.clear();
          },
        ),
      ],
    );
  }

  Widget _buildMetricFields() {
    return Row(
      children: [
        Expanded(
          child: _buildTextField(
            controller: _vipRankController,
            label: 'VIPランク',
            hintText: '0',
            icon: Icons.star,
            keyboardType: TextInputType.number,
          ),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: _buildTextField(
            controller: _lifetimeValueController,
            label: 'LTV (円)',
            hintText: '1200000',
            icon: Icons.payments,
            keyboardType: const TextInputType.numberWithOptions(decimal: true),
          ),
        ),
      ],
    );
  }

  Widget _buildChannelFields() {
    return Column(
      children: [
        _buildTextField(
          controller: _preferredChannelController,
          label: 'Preferred Channel',
          hintText: 'LINE / Email / Phone',
          icon: Icons.forum,
        ),
        const SizedBox(height: 12),
        _buildTextField(
          controller: _leadSourceController,
          label: 'Lead Source',
          hintText: '紹介 / Web広告',
          icon: Icons.campaign,
        ),
      ],
    );
  }

  Widget _buildInteractionSection() {
    return Card(
      color: EnterpriseColors.surfaceGray,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.timeline, color: EnterpriseColors.primaryBlue),
                const SizedBox(width: 8),
                const Text(
                  '最新コンタクト',
                  style: TextStyle(
                    color: EnterpriseColors.textPrimary,
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const Spacer(),
                TextButton(
                  onPressed: _pickInteractionDate,
                  child: Text(
                    _lastInteractionAt == null
                        ? '日時を選択'
                        : _lastInteractionAt!.toLocal().toString(),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            DropdownButtonFormField<String>(
              value: _interactionType,
              dropdownColor: EnterpriseColors.surfaceGray,
              decoration: const InputDecoration(
                labelText: '種別',
                labelStyle: TextStyle(color: EnterpriseColors.textSecondary),
                enabledBorder: OutlineInputBorder(
                  borderSide: BorderSide(color: EnterpriseColors.borderGray),
                ),
                focusedBorder: OutlineInputBorder(
                  borderSide: BorderSide(
                    color: EnterpriseColors.primaryBlue,
                    width: 2,
                  ),
                ),
              ),
              items: const [
                DropdownMenuItem(value: 'note', child: Text('メモ')),
                DropdownMenuItem(value: 'meeting', child: Text('面談')),
                DropdownMenuItem(value: 'call', child: Text('電話')),
                DropdownMenuItem(value: 'fitting', child: Text('採寸')),
              ],
              onChanged: (value) {
                if (value == null) return;
                setState(() => _interactionType = value);
              },
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _interactionNoteController,
              maxLines: 3,
              style: const TextStyle(color: EnterpriseColors.textPrimary),
              decoration: InputDecoration(
                hintText: 'メモを入力',
                filled: true,
                fillColor: EnterpriseColors.deepBlack,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(8),
                  borderSide: BorderSide(color: EnterpriseColors.borderGray),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildNotesField() {
    return _buildTextField(
      controller: _notesController,
      label: '社内メモ',
      hintText: 'VIP顧客。〇〇ブランドを好む 等',
      icon: Icons.sticky_note_2,
      maxLines: 4,
    );
  }

  Future<void> _pickInteractionDate() async {
    final now = DateTime.now();
    final pickedDate = await showDatePicker(
      context: context,
      initialDate: _lastInteractionAt ?? now,
      firstDate: DateTime(now.year - 5),
      lastDate: DateTime(now.year + 1),
    );
    if (pickedDate == null) return;
    final pickedTime = await showTimePicker(
      context: context,
      initialTime: TimeOfDay.fromDateTime(_lastInteractionAt ?? now),
    );
    setState(() {
      _lastInteractionAt = DateTime(
        pickedDate.year,
        pickedDate.month,
        pickedDate.day,
        pickedTime?.hour ?? 0,
        pickedTime?.minute ?? 0,
      );
    });
  }

  List<CustomerInteraction>? _composeInteractionPayload() {
    final note = _interactionNoteController.text.trim();
    if (note.isEmpty) {
      return null;
    }
    return [
      CustomerInteraction(
        type: _interactionType,
        note: note,
        timestamp: (_lastInteractionAt ?? DateTime.now()).toUtc(),
      ),
    ];
  }
}
