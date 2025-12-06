import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../config/app_config.dart';
import '../../config/enterprise_theme.dart';
import '../../providers/auth_provider.dart';

/// ログイン画面
class LoginScreen extends ConsumerStatefulWidget {
  const LoginScreen({super.key});

  @override
  ConsumerState<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends ConsumerState<LoginScreen> {
  final _formKey = GlobalKey<FormState>();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  bool _isLoading = false;
  bool _obscurePassword = true;

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _handleLogin() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    setState(() {
      _isLoading = true;
    });

    try {
      await ref.read(signInProvider(
        email: _emailController.text.trim(),
        password: _passwordController.text,
      ).future);

      // ログイン成功 - 認証状態の変更により自動的にメイン画面に遷移
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('ログインに成功しました'),
            backgroundColor: EnterpriseColors.successGreen,
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        String errorMessage = 'ログインに失敗しました';

        // Firebase Authのエラーメッセージを解析
        final errorString = e.toString().toLowerCase();
        if (errorString.contains('user-not-found') ||
            errorString.contains('wrong-password')) {
          errorMessage = 'メールアドレスまたはパスワードが正しくありません';
        } else if (errorString.contains('invalid-email')) {
          errorMessage = 'メールアドレスの形式が正しくありません';
        } else if (errorString.contains('too-many-requests')) {
          errorMessage = 'ログイン試行回数が多すぎます。しばらく待ってから再試行してください';
        } else if (errorString.contains('network')) {
          errorMessage = 'ネットワークエラーが発生しました';
        }

        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(errorMessage),
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
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(24.0),
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 400),
              child: Form(
                key: _formKey,
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    // ロゴ・タイトル
                    Icon(
                      Icons.account_circle,
                      size: 80,
                      color: EnterpriseColors.metallicGold,
                    ),
                    const SizedBox(height: 24),
                    Text(
                      'TailorCloud',
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 32,
                        fontWeight: FontWeight.bold,
                        color: EnterpriseColors.metallicGold,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      'Factory to Wardrobe OS',
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 14,
                        color: EnterpriseColors.textSecondary,
                      ),
                    ),
                    const SizedBox(height: 48),

                    // メールアドレス入力
                    TextFormField(
                      controller: _emailController,
                      keyboardType: TextInputType.emailAddress,
                      textInputAction: TextInputAction.next,
                      style:
                          const TextStyle(color: EnterpriseColors.textPrimary),
                      decoration: InputDecoration(
                        labelText: 'メールアドレス',
                        labelStyle: const TextStyle(
                          color: EnterpriseColors.textSecondary,
                        ),
                        prefixIcon: const Icon(
                          Icons.email,
                          color: EnterpriseColors.textSecondary,
                        ),
                        filled: true,
                        fillColor: EnterpriseColors.surfaceGray,
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(8),
                          borderSide: const BorderSide(
                            color: EnterpriseColors.borderGray,
                          ),
                        ),
                        enabledBorder: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(8),
                          borderSide: const BorderSide(
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
                      ),
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'メールアドレスを入力してください';
                        }
                        if (!value.contains('@')) {
                          return '有効なメールアドレスを入力してください';
                        }
                        return null;
                      },
                    ),
                    const SizedBox(height: 16),

                    // パスワード入力
                    TextFormField(
                      controller: _passwordController,
                      obscureText: _obscurePassword,
                      textInputAction: TextInputAction.done,
                      onFieldSubmitted: (_) => _handleLogin(),
                      style:
                          const TextStyle(color: EnterpriseColors.textPrimary),
                      decoration: InputDecoration(
                        labelText: 'パスワード',
                        labelStyle: const TextStyle(
                          color: EnterpriseColors.textSecondary,
                        ),
                        prefixIcon: const Icon(
                          Icons.lock,
                          color: EnterpriseColors.textSecondary,
                        ),
                        suffixIcon: IconButton(
                          icon: Icon(
                            _obscurePassword
                                ? Icons.visibility
                                : Icons.visibility_off,
                            color: EnterpriseColors.textSecondary,
                          ),
                          onPressed: () {
                            setState(() {
                              _obscurePassword = !_obscurePassword;
                            });
                          },
                        ),
                        filled: true,
                        fillColor: EnterpriseColors.surfaceGray,
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(8),
                          borderSide: const BorderSide(
                            color: EnterpriseColors.borderGray,
                          ),
                        ),
                        enabledBorder: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(8),
                          borderSide: const BorderSide(
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
                      ),
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'パスワードを入力してください';
                        }
                        if (value.length < 6) {
                          return 'パスワードは6文字以上で入力してください';
                        }
                        return null;
                      },
                    ),
                    const SizedBox(height: 32),

                    // ログインボタン
                    ElevatedButton(
                      onPressed: _isLoading ? null : _handleLogin,
                      style: ElevatedButton.styleFrom(
                        backgroundColor: EnterpriseColors.primaryBlue,
                        foregroundColor: Colors.white,
                        padding: const EdgeInsets.symmetric(vertical: 16),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                        elevation: 0,
                      ),
                      child: _isLoading
                          ? const SizedBox(
                              height: 20,
                              width: 20,
                              child: CircularProgressIndicator(
                                strokeWidth: 2,
                                valueColor:
                                    AlwaysStoppedAnimation<Color>(Colors.white),
                              ),
                            )
                          : const Text(
                              'ログイン',
                              style: TextStyle(
                                fontSize: 16,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                    ),
                    const SizedBox(height: 16),

                    // Googleログイン（Firebase有効時のみ表示）
                    if (AppConfig.enableFirebase) ...[
                      Row(
                        children: const [
                          Expanded(
                            child: Divider(color: EnterpriseColors.borderGray),
                          ),
                          Padding(
                            padding: EdgeInsets.symmetric(horizontal: 8),
                            child: Text(
                              'または',
                              style: TextStyle(
                                fontSize: 12,
                                color: EnterpriseColors.textSecondary,
                              ),
                            ),
                          ),
                          Expanded(
                            child: Divider(color: EnterpriseColors.borderGray),
                          ),
                        ],
                      ),
                      const SizedBox(height: 16),
                      OutlinedButton.icon(
                        onPressed: _isLoading ? null : _handleGoogleLogin,
                        style: OutlinedButton.styleFrom(
                          foregroundColor: Colors.white,
                          side: const BorderSide(
                            color: EnterpriseColors.metallicGold,
                          ),
                          padding: const EdgeInsets.symmetric(vertical: 12),
                        ),
                        icon: const Icon(
                          Icons.login,
                          color: EnterpriseColors.metallicGold,
                        ),
                        label: const Text(
                          'Google アカウントでログイン',
                          style: TextStyle(fontSize: 14),
                        ),
                      ),
                      const SizedBox(height: 16),
                    ],

                    // 開発用メッセージ
                    Text(
                      AppConfig.enableFirebase
                          ? 'メール/パスワードまたはGoogleでログインできます'
                          : 'メール/パスワードでログインしてください',
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 12,
                        color: EnterpriseColors.textTertiary,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }

  Future<void> _handleGoogleLogin() async {
    setState(() {
      _isLoading = true;
    });

    try {
      await signInWithGoogle(ref);

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Googleログインに成功しました'),
            backgroundColor: EnterpriseColors.successGreen,
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        String errorMessage = 'Googleログインに失敗しました';
        final errorString = e.toString().toLowerCase();
        if (errorString.contains('domain')) {
          errorMessage = e.toString();
        } else if (errorString.contains('network')) {
          errorMessage = 'ネットワークエラーが発生しました';
        }

        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(errorMessage),
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
}
