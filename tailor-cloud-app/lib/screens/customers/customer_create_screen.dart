import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../config/enterprise_theme.dart';
import '../../models/customer.dart';
import '../../providers/customer_provider.dart';

/// 新規顧客登録画面
class CustomerCreateScreen extends ConsumerStatefulWidget {
  const CustomerCreateScreen({super.key});

  @override
  ConsumerState<CustomerCreateScreen> createState() => _CustomerCreateScreenState();
}

class _CustomerCreateScreenState extends ConsumerState<CustomerCreateScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _nameKanaController = TextEditingController();
  final _emailController = TextEditingController();
  final _phoneController = TextEditingController();
  final _addressController = TextEditingController();
  
  bool _isLoading = false;

  @override
  void dispose() {
    _nameController.dispose();
    _nameKanaController.dispose();
    _emailController.dispose();
    _phoneController.dispose();
    _addressController.dispose();
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
      );

      await ref.read(createCustomerProvider(request).future);

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('顧客を登録しました'),
            backgroundColor: EnterpriseColors.successGreen,
          ),
        );
        Navigator.pop(context, true);
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('エラー: ${e.toString()}'),
            backgroundColor: EnterpriseColors.errorRed,
          ),
        );
      }
    } finally {
      if (mounted) {
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
              _buildTextField(
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
}

