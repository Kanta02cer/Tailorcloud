import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';
import '../../config/enterprise_theme.dart';
import '../../models/customer.dart';
import '../../models/fabric.dart';
import '../../models/order.dart';
import '../../providers/customer_provider.dart';
import '../../providers/fabric_provider.dart';
import '../../providers/order_provider.dart';
import '../../providers/measurement_validation_provider.dart';
import '../customers/customer_create_screen.dart';
import '../orders/order_confirm_screen.dart';

/// クイック発注画面（3ステップフロー改善版）
class QuickOrderScreen extends ConsumerStatefulWidget {
  const QuickOrderScreen({super.key});

  @override
  ConsumerState<QuickOrderScreen> createState() => _QuickOrderScreenState();
}

class _QuickOrderScreenState extends ConsumerState<QuickOrderScreen> {
  // ステップ管理
  int _currentStep = 0; // 0: 顧客選択, 1: 生地選択, 2: 注文情報入力

  // フォームコントローラー
  final _formKey = GlobalKey<FormState>();
  final _amountController = TextEditingController();
  final _deliveryDateController = TextEditingController();
  final _jacketLengthController = TextEditingController();
  final _sleeveController = TextEditingController();
  final _chestController = TextEditingController();

  // 選択状態
  Customer? _selectedCustomer;
  Fabric? _selectedFabric;
  DateTime? _deliveryDate;

  // バリデーション状態
  ValidationResponse? _validationResult;
  bool _isValidating = false;

  @override
  void dispose() {
    _amountController.dispose();
    _deliveryDateController.dispose();
    _jacketLengthController.dispose();
    _sleeveController.dispose();
    _chestController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    // TODO: 認証からテナントIDを取得
    const tenantId = 'tenant-123'; // 仮のテナントID

    return Scaffold(
      backgroundColor: EnterpriseColors.deepBlack,
      appBar: AppBar(
        backgroundColor: EnterpriseColors.surfaceGray,
        title: const Text(
          'クイック発注',
          style: TextStyle(
            color: EnterpriseColors.textPrimary,
            fontSize: 20,
            fontWeight: FontWeight.bold,
          ),
        ),
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: EnterpriseColors.textPrimary),
          onPressed: () => Navigator.pop(context),
        ),
      ),
      body: Column(
        children: [
          // ステップインジケーター
          _buildStepIndicator(),
          
          // メインコンテンツ
          Expanded(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(16),
              child: Form(
                key: _formKey,
                child: _buildStepContent(tenantId),
              ),
            ),
          ),
          
          // ナビゲーションボタン
          _buildNavigationButtons(),
        ],
      ),
    );
  }

  Widget _buildStepIndicator() {
    return Container(
      padding: const EdgeInsets.all(16),
      color: EnterpriseColors.surfaceGray,
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          _buildStepItem(0, '顧客選択', Icons.person),
          _buildStepConnector(),
          _buildStepItem(1, '生地選択', Icons.inventory_2),
          _buildStepConnector(),
          _buildStepItem(2, '注文情報', Icons.edit),
        ],
      ),
    );
  }

  Widget _buildStepItem(int step, String label, IconData icon) {
    final isActive = _currentStep == step;
    final isCompleted = _currentStep > step;

    return Column(
      children: [
        Container(
          width: 48,
          height: 48,
          decoration: BoxDecoration(
            color: isActive || isCompleted
                ? EnterpriseColors.primaryBlue
                : EnterpriseColors.borderGray,
            shape: BoxShape.circle,
          ),
          child: Icon(
            isCompleted ? Icons.check : icon,
            color: Colors.white,
            size: 24,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          label,
          style: TextStyle(
            color: isActive || isCompleted
                ? EnterpriseColors.primaryBlue
                : EnterpriseColors.textSecondary,
            fontSize: 12,
            fontWeight: isActive ? FontWeight.bold : FontWeight.normal,
          ),
        ),
      ],
    );
  }

  Widget _buildStepConnector() {
    return Container(
      width: 40,
      height: 2,
      color: _currentStep > 0
          ? EnterpriseColors.primaryBlue
          : EnterpriseColors.borderGray,
    );
  }

  Widget _buildStepContent(String tenantId) {
    switch (_currentStep) {
      case 0:
        return _buildStep1CustomerSelection(tenantId);
      case 1:
        return _buildStep2FabricSelection(tenantId);
      case 2:
        return _buildStep3OrderInfo();
      default:
        return const SizedBox.shrink();
    }
  }

  Widget _buildStep1CustomerSelection(String tenantId) {
    final customersAsync = ref.watch(customerListProvider(tenantId));

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'ステップ1: 顧客を選択',
          style: TextStyle(
            color: EnterpriseColors.textPrimary,
            fontSize: 20,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 24),
        customersAsync.when(
          data: (customers) {
            if (customers.isEmpty) {
              return Card(
                color: EnterpriseColors.surfaceGray,
                child: Padding(
                  padding: const EdgeInsets.all(20),
                  child: Column(
                    children: [
                      const Text(
                        '顧客が登録されていません',
                        style: TextStyle(
                          color: EnterpriseColors.textSecondary,
                          fontSize: 16,
                        ),
                      ),
                      const SizedBox(height: 16),
                      ElevatedButton.icon(
                        onPressed: () => _navigateToNewCustomer(context, tenantId),
                        icon: const Icon(Icons.person_add),
                        label: const Text('新規顧客を登録'),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: EnterpriseColors.primaryBlue,
                        ),
                      ),
                    ],
                  ),
                ),
              );
            }

            return Column(
              children: [
                // 既存顧客選択
                DropdownButtonFormField<Customer>(
                  value: _selectedCustomer,
                  decoration: InputDecoration(
                    labelText: '顧客を選択',
                    labelStyle: const TextStyle(
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
                  dropdownColor: EnterpriseColors.surfaceGray,
                  style: const TextStyle(
                    color: EnterpriseColors.textPrimary,
                  ),
                  items: customers.map((customer) {
                    return DropdownMenuItem<Customer>(
                      value: customer,
                      child: Text(customer.name),
                    );
                  }).toList(),
                  onChanged: (customer) {
                    setState(() {
                      _selectedCustomer = customer;
                    });
                  },
                  validator: (value) {
                    if (value == null) {
                      return '顧客を選択してください';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: 16),
                // または新規登録
                OutlinedButton.icon(
                  onPressed: () => _navigateToNewCustomer(context, tenantId),
                  icon: const Icon(Icons.person_add),
                  label: const Text('新規顧客を登録'),
                  style: OutlinedButton.styleFrom(
                    foregroundColor: EnterpriseColors.primaryBlue,
                    side: const BorderSide(color: EnterpriseColors.primaryBlue),
                    padding: const EdgeInsets.symmetric(vertical: 16),
                  ),
                ),
              ],
            );
          },
          loading: () => const Center(
            child: CircularProgressIndicator(
              color: EnterpriseColors.primaryBlue,
            ),
          ),
          error: (error, stack) => Card(
            color: EnterpriseColors.surfaceGray,
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                children: [
                  const Icon(
                    Icons.error_outline,
                    color: EnterpriseColors.errorRed,
                    size: 48,
                  ),
                  const SizedBox(height: 16),
                  Text(
                    'エラー: $error',
                    style: const TextStyle(
                      color: EnterpriseColors.textPrimary,
                    ),
                    textAlign: TextAlign.center,
                  ),
                ],
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildStep2FabricSelection(String tenantId) {
    final fabricsAsync = ref.watch(fabricListProvider(
      FabricListParams(tenantId: tenantId),
    ));

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'ステップ2: 生地を選択',
          style: TextStyle(
            color: EnterpriseColors.textPrimary,
            fontSize: 20,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 24),
        fabricsAsync.when(
          data: (fabrics) {
            if (fabrics.isEmpty) {
              return Card(
                color: EnterpriseColors.surfaceGray,
                child: const Padding(
                  padding: EdgeInsets.all(20),
                  child: Text(
                    '生地が登録されていません',
                    style: TextStyle(
                      color: EnterpriseColors.textSecondary,
                      fontSize: 16,
                    ),
                  ),
                ),
              );
            }

            return DropdownButtonFormField<Fabric>(
              value: _selectedFabric,
              decoration: InputDecoration(
                labelText: '生地を選択',
                labelStyle: const TextStyle(
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
              dropdownColor: EnterpriseColors.surfaceGray,
              style: const TextStyle(
                color: EnterpriseColors.textPrimary,
              ),
              items: fabrics.map((fabric) {
                final displayName = fabric.brand != null && fabric.brand!.isNotEmpty
                    ? '${fabric.brand} - ${fabric.name}'
                    : fabric.name;
                return DropdownMenuItem<Fabric>(
                  value: fabric,
                  child: Text(displayName),
                );
              }).toList(),
              onChanged: (fabric) {
                setState(() {
                  _selectedFabric = fabric;
                });
              },
              validator: (value) {
                if (value == null) {
                  return '生地を選択してください';
                }
                return null;
              },
            );
          },
          loading: () => const Center(
            child: CircularProgressIndicator(
              color: EnterpriseColors.primaryBlue,
            ),
          ),
          error: (error, stack) => Card(
            color: EnterpriseColors.surfaceGray,
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                children: [
                  const Icon(
                    Icons.error_outline,
                    color: EnterpriseColors.errorRed,
                    size: 48,
                  ),
                  const SizedBox(height: 16),
                  Text(
                    'エラー: $error',
                    style: const TextStyle(
                      color: EnterpriseColors.textPrimary,
                    ),
                    textAlign: TextAlign.center,
                  ),
                ],
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildStep3OrderInfo() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'ステップ3: 注文情報を入力',
          style: TextStyle(
            color: EnterpriseColors.textPrimary,
            fontSize: 20,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 24),
        
        // 金額入力
        _buildTextField(
          controller: _amountController,
          label: '金額',
          hintText: '135000',
          icon: Icons.attach_money,
          keyboardType: TextInputType.number,
          prefixText: '¥',
          validator: (value) {
            if (value == null || value.isEmpty) {
              return '金額を入力してください';
            }
            if (int.tryParse(value) == null) {
              return '有効な数値を入力してください';
            }
            return null;
          },
        ),
        
        const SizedBox(height: 16),
        
        // 納期選択
        TextField(
          controller: _deliveryDateController,
          readOnly: true,
          style: const TextStyle(color: EnterpriseColors.textPrimary),
          decoration: InputDecoration(
            labelText: '納期',
            labelStyle: const TextStyle(
              color: EnterpriseColors.textSecondary,
            ),
            hintText: '納期を選択',
            hintStyle: const TextStyle(
              color: EnterpriseColors.textSecondary,
            ),
            prefixIcon: const Icon(
              Icons.calendar_today,
              color: EnterpriseColors.primaryBlue,
            ),
            suffixIcon: const Icon(
              Icons.arrow_drop_down,
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
          onTap: () => _selectDeliveryDate(context),
        ),
        
        const SizedBox(height: 24),
        
        // 採寸データ（簡易版）
        const Text(
          '採寸データ（任意）',
          style: TextStyle(
            color: EnterpriseColors.textPrimary,
            fontSize: 16,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 12),
        const Text(
          '主要寸法のみ入力してください',
          style: TextStyle(
            color: EnterpriseColors.textSecondary,
            fontSize: 12,
          ),
        ),
        const SizedBox(height: 16),
        
        Row(
          children: [
            Expanded(
              child: _buildTextField(
                controller: _jacketLengthController,
                label: 'ジャケット長',
                hintText: '72.5',
                icon: Icons.straighten,
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
                suffixText: 'cm',
                onChanged: _onMeasurementChanged,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: _buildTextField(
                controller: _sleeveController,
                label: '袖長',
                hintText: '60.0',
                icon: Icons.straighten,
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
                suffixText: 'cm',
                onChanged: _onMeasurementChanged,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: _buildTextField(
                controller: _chestController,
                label: '胸囲',
                hintText: '100.0',
                icon: Icons.straighten,
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
                suffixText: 'cm',
                onChanged: _onMeasurementChanged,
              ),
            ),
          ],
        ),

        // バリデーション結果表示
        if (_validationResult != null) ...[
          const SizedBox(height: 16),
          _buildValidationAlerts(_validationResult!),
        ],
      ],
    );
  }

  Widget _buildTextField({
    required TextEditingController controller,
    required String label,
    required String hintText,
    required IconData icon,
    TextInputType? keyboardType,
    String? prefixText,
    String? suffixText,
    String? Function(String?)? validator,
    void Function(String)? onChanged,
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
          ],
        ),
        const SizedBox(height: 8),
        TextFormField(
          controller: controller,
          keyboardType: keyboardType,
          style: const TextStyle(color: EnterpriseColors.textPrimary),
          onChanged: onChanged,
          validator: validator,
          decoration: InputDecoration(
            hintText: hintText,
            hintStyle: const TextStyle(
              color: EnterpriseColors.textSecondary,
            ),
            prefixText: prefixText,
            suffixText: suffixText,
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

  Widget _buildNavigationButtons() {
    return Container(
      padding: const EdgeInsets.all(16),
      color: EnterpriseColors.surfaceGray,
      child: Row(
        children: [
          // 戻るボタン
          if (_currentStep > 0)
            Expanded(
              child: OutlinedButton(
                onPressed: () {
                  setState(() {
                    _currentStep--;
                  });
                },
                style: OutlinedButton.styleFrom(
                  foregroundColor: EnterpriseColors.textPrimary,
                  side: BorderSide(color: EnterpriseColors.borderGray),
                  padding: const EdgeInsets.symmetric(vertical: 16),
                ),
                child: const Text('戻る'),
              ),
            ),
          if (_currentStep > 0) const SizedBox(width: 12),
          
          // 次へ/確定ボタン
          Expanded(
            flex: _currentStep > 0 ? 1 : 1,
            child: ElevatedButton(
              onPressed: _handleNext,
              style: ElevatedButton.styleFrom(
                backgroundColor: EnterpriseColors.primaryBlue,
                padding: const EdgeInsets.symmetric(vertical: 16),
              ),
              child: Text(
                _currentStep < 2 ? '次へ' : '注文を作成',
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  void _handleNext() {
    if (_currentStep < 2) {
      // バリデーション
      if (_currentStep == 0) {
        if (_selectedCustomer == null) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('顧客を選択してください'),
              backgroundColor: EnterpriseColors.errorRed,
            ),
          );
          return;
        }
      } else if (_currentStep == 1) {
        if (_selectedFabric == null) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('生地を選択してください'),
              backgroundColor: EnterpriseColors.errorRed,
            ),
          );
          return;
        }
      }

      setState(() {
        _currentStep++;
      });
    } else {
      // 最終ステップ: 注文作成
      _createOrder();
    }
  }

  Future<void> _createOrder() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    if (_selectedCustomer == null || _selectedFabric == null || _deliveryDate == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('すべての項目を入力してください'),
          backgroundColor: EnterpriseColors.errorRed,
        ),
      );
      return;
    }

    try {
      // 採寸データを準備
      Map<String, dynamic>? measurementData;
      if (_jacketLengthController.text.isNotEmpty ||
          _sleeveController.text.isNotEmpty ||
          _chestController.text.isNotEmpty) {
        measurementData = {};
        if (_jacketLengthController.text.isNotEmpty) {
          measurementData['jacket_length'] = double.tryParse(_jacketLengthController.text);
        }
        if (_sleeveController.text.isNotEmpty) {
          measurementData['sleeve'] = double.tryParse(_sleeveController.text);
        }
        if (_chestController.text.isNotEmpty) {
          measurementData['chest'] = double.tryParse(_chestController.text);
        }
      }

      // 注文を作成
      const tenantId = 'tenant-123'; // TODO: 認証から取得
      final orderRequest = CreateOrderRequest(
        customerId: _selectedCustomer!.id,
        fabricId: _selectedFabric!.id,
        totalAmount: int.parse(_amountController.text),
        deliveryDate: _deliveryDate!,
        details: OrderDetails(
          description: 'オーダースーツ縫製',
          measurementData: measurementData,
        ),
      );

      final order = await ref.read(createOrderProvider(orderRequest).future);

      // 注文確認画面に遷移
      if (mounted) {
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(
            builder: (context) => OrderConfirmScreen(
              order: order,
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
      }
    }
  }

  // 採寸データ変更時のバリデーション
  Future<void> _onMeasurementChanged(String value) async {
    // 顧客が選択されていない場合はスキップ
    if (_selectedCustomer == null) {
      return;
    }

    // 採寸データを準備
    Map<String, dynamic> measurementData = {};
    if (_jacketLengthController.text.isNotEmpty) {
      final val = double.tryParse(_jacketLengthController.text);
      if (val != null) measurementData['jacket_length'] = val;
    }
    if (_sleeveController.text.isNotEmpty) {
      final val = double.tryParse(_sleeveController.text);
      if (val != null) measurementData['sleeve'] = val;
    }
    if (_chestController.text.isNotEmpty) {
      final val = double.tryParse(_chestController.text);
      if (val != null) measurementData['chest'] = val;
    }

    // 採寸データが空の場合はスキップ
    if (measurementData.isEmpty) {
      setState(() {
        _validationResult = null;
      });
      return;
    }

    // バリデーション実行
    setState(() {
      _isValidating = true;
    });

    try {
      const tenantId = 'tenant-123'; // TODO: 認証から取得
      final request = ValidateMeasurementsRequest(
        customerId: _selectedCustomer!.id,
        currentMeasurements: measurementData,
      );

      final result = await ref.read(validateMeasurementsProvider(request).future);
      
      if (mounted) {
        setState(() {
          _validationResult = result;
          _isValidating = false;
        });
      }
    } catch (e) {
      // バリデーションエラーは無視（前回データがない場合など）
      if (mounted) {
        setState(() {
          _validationResult = null;
          _isValidating = false;
        });
      }
    }
  }

  // バリデーションアラート表示
  Widget _buildValidationAlerts(ValidationResponse result) {
    if (result.alerts.isEmpty) {
      return const SizedBox.shrink();
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        for (final alert in result.alerts)
          Container(
            margin: const EdgeInsets.only(bottom: 8),
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: alert.severity == 'error'
                  ? EnterpriseColors.errorRed.withOpacity(0.1)
                  : Colors.orange.withOpacity(0.1),
              borderRadius: BorderRadius.circular(8),
              border: Border.all(
                color: alert.severity == 'error'
                    ? EnterpriseColors.errorRed
                    : Colors.orange,
                width: 1,
              ),
            ),
            child: Row(
              children: [
                Icon(
                  alert.severity == 'error' ? Icons.error : Icons.warning,
                  color: alert.severity == 'error'
                      ? EnterpriseColors.errorRed
                      : Colors.orange,
                  size: 20,
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: Text(
                    alert.message,
                    style: TextStyle(
                      color: alert.severity == 'error'
                          ? EnterpriseColors.errorRed
                          : Colors.orange,
                      fontSize: 12,
                    ),
                  ),
                ),
              ],
            ),
          ),
      ],
    );
  }

  Future<void> _navigateToNewCustomer(BuildContext context, String tenantId) async {
    final result = await Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => const CustomerCreateScreen(),
      ),
    );

    if (result == true && mounted) {
      // 顧客一覧を再読み込み
      ref.invalidate(customerListProvider(tenantId));
      
      // 最新の顧客一覧を取得して、最後に作成された顧客を選択
      final customersAsync = ref.read(customerListProvider(tenantId).future);
      final customers = await customersAsync;
      if (customers.isNotEmpty) {
        setState(() {
          _selectedCustomer = customers.last;
        });
      }
    }
  }

  Future<void> _selectDeliveryDate(BuildContext context) async {
    final picked = await showDatePicker(
      context: context,
      initialDate: DateTime.now().add(const Duration(days: 14)),
      firstDate: DateTime.now(),
      lastDate: DateTime.now().add(const Duration(days: 365)),
      builder: (context, child) {
        return Theme(
          data: Theme.of(context).copyWith(
            colorScheme: const ColorScheme.dark(
              primary: EnterpriseColors.primaryBlue,
              onPrimary: Colors.white,
              surface: EnterpriseColors.surfaceGray,
              onSurface: EnterpriseColors.textPrimary,
            ),
          ),
          child: child!,
        );
      },
    );

    if (picked != null) {
      setState(() {
        _deliveryDate = picked;
        _deliveryDateController.text = DateFormat('yyyy年MM月dd日').format(picked);
      });
    }
  }
}
