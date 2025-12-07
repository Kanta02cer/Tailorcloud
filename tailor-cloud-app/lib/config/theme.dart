import 'package:flutter/material.dart';

/// TailorCloud カラーパレット
/// 参考: docs/11_UI_Design_Specifications.md
class AppColors {
  // Primary Navy
  static const Color primaryNavy = Color(0xFF1A1F3A);
  static const Color primaryNavyDark = Color(0xFF0F1424);
  static const Color primaryNavyLight = Color(0xFF2A3054);

  // Accent Colors
  static const Color accentGold = Color(0xFFD4AF37);
  static const Color accentCream = Color(0xFFF5F1E8);

  // Status Colors
  static const Color statusAvailable = Color(0xFF10B981); // 緑
  static const Color statusLimited = Color(0xFFF59E0B); // 黄色
  static const Color statusSoldOut = Color(0xFFEF4444); // 赤
  static const Color statusWarning = Color(0xFFF59E0B);

  // Neutral Colors
  static const Color neutralGray100 = Color(0xFFF3F4F6);
  static const Color neutralGray200 = Color(0xFFE5E7EB);
  static const Color neutralGray500 = Color(0xFF6B7280);
  static const Color neutralGray900 = Color(0xFF111827);

  // Background
  static const Color backgroundWhite = Color(0xFFFFFFFF);
  static const Color backgroundCream = Color(0xFFF5F1E8);
}

/// TailorCloud テキストスタイル
class AppTextStyles {
  // H1 - 大見出し（ダッシュボードKPI）
  static const TextStyle h1 = TextStyle(
    fontSize: 32,
    fontWeight: FontWeight.bold,
    letterSpacing: -0.5,
    height: 1.2,
    color: AppColors.primaryNavy,
  );

  // H2 - 中見出し（画面タイトル）
  static const TextStyle h2 = TextStyle(
    fontSize: 24,
    fontWeight: FontWeight.bold,
    letterSpacing: -0.3,
    height: 1.3,
    color: AppColors.primaryNavy,
  );

  // Body - 本文
  static const TextStyle body = TextStyle(
    fontSize: 16,
    fontWeight: FontWeight.normal,
    letterSpacing: 0,
    height: 1.5,
    color: AppColors.primaryNavy,
  );

  // KPI Number - 数字表示用
  static const TextStyle kpiNumber = TextStyle(
    fontSize: 28,
    fontWeight: FontWeight.bold,
    letterSpacing: -1,
    color: AppColors.primaryNavy,
  );

  // Caption - キャプション
  static const TextStyle caption = TextStyle(
    fontSize: 12,
    fontWeight: FontWeight.normal,
    letterSpacing: 0.2,
    color: AppColors.neutralGray500,
  );
}

/// TailorCloud テーマ
ThemeData appTheme() {
  return ThemeData(
    useMaterial3: true,
    primaryColor: AppColors.primaryNavy,
    scaffoldBackgroundColor: AppColors.backgroundCream,
    fontFamily: 'NotoSansJP',
    colorScheme: ColorScheme.fromSeed(
      seedColor: AppColors.primaryNavy,
      primary: AppColors.primaryNavy,
      secondary: AppColors.accentGold,
      surface: AppColors.backgroundWhite,
    ),
    textTheme: const TextTheme(
      headlineLarge: AppTextStyles.h1,
      headlineMedium: AppTextStyles.h2,
      bodyLarge: AppTextStyles.body,
      bodySmall: AppTextStyles.caption,
    ),
    cardTheme: CardThemeData(
      elevation: 2,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
      ),
      color: AppColors.backgroundWhite,
    ),
    elevatedButtonTheme: ElevatedButtonThemeData(
      style: ElevatedButton.styleFrom(
        backgroundColor: AppColors.primaryNavy,
        foregroundColor: AppColors.backgroundWhite,
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
      ),
    ),
  );
}
