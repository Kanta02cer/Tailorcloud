import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:firebase_core/firebase_core.dart';
import 'config/enterprise_theme.dart';
import 'screens/main_screen.dart';

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

class TailorCloudApp extends StatelessWidget {
  const TailorCloudApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'TailorCloud',
      theme: enterpriseTheme(),
      debugShowCheckedModeBanner: false,
      home: const MainScreen(),
    );
  }
}

