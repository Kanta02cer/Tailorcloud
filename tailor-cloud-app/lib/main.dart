import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'config/enterprise_theme.dart';
import 'config/firebase_config.dart';
import 'config/app_config.dart';
import 'services/logger.dart';
import 'screens/main_screen.dart';
import 'screens/auth/login_screen.dart';
import 'providers/auth_provider.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // 環境設定をログ出力
  Logger.info('Starting TailorCloud App');
  if (AppConfig.enableDebugLogging) {
    Logger.debug(AppConfig.info);
  }

  // Firebase初期化（環境変数で制御）
  final firebaseInitialized = await FirebaseConfig.initialize();
  if (!firebaseInitialized && AppConfig.enableFirebase) {
    Logger.warning(
        'Firebase initialization failed. App will run without Firebase features.');
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
              valueColor:
                  AlwaysStoppedAnimation<Color>(EnterpriseColors.primaryBlue),
            ),
          ),
        ),
        error: (error, stack) => const LoginScreen(),
      ),
    );
  }
}
