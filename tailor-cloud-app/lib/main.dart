import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:firebase_core/firebase_core.dart';
import 'config/enterprise_theme.dart';
import 'screens/main_screen.dart';
import 'screens/auth/login_screen.dart';
import 'providers/auth_provider.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Firebase初期化（オプショナル - 設定ファイルがない場合でも動作）
  try {
    await Firebase.initializeApp();
  } catch (e) {
    // Firebase設定ファイルがない場合は警告のみ
    debugPrint('Warning: Firebase initialization failed: $e');
    debugPrint('The app will run without Firebase features.');
  }
  
  runApp(
    const ProviderScope(
      child: TailorCloudApp(),
    ),
  );
}

class TailorCloudApp extends ConsumerWidget {
  const TailorCloudApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // 認証状態を監視
    final authState = ref.watch(authStateChangesProvider);

    return MaterialApp(
      title: 'TailorCloud',
      theme: enterpriseTheme(),
      debugShowCheckedModeBanner: false,
      home: authState.when(
        data: (user) {
          // ユーザーがログインしている場合はメイン画面、そうでなければログイン画面
          return user != null ? const MainScreen() : const LoginScreen();
        },
        loading: () => const Scaffold(
          backgroundColor: EnterpriseColors.deepBlack,
          body: Center(
            child: CircularProgressIndicator(
              valueColor: AlwaysStoppedAnimation<Color>(EnterpriseColors.primaryBlue),
            ),
          ),
        ),
        error: (error, stack) => const LoginScreen(),
      ),
    );
  }
}

