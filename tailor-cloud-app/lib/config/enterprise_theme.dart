import 'package:flutter/material.dart';

/// TailorCloud エンタープライズテーマ
/// "Digital Bespoke" - 伝統的な重厚感とAppleのような洗練された操作性の融合
class EnterpriseColors {
  // Background Colors
  static const Color deepBlack = Color(0xFF050505);
  static const Color regalisBlack = Color(0xFF151515);
  static const Color darkGray = Color(0xFF1E1E1E);
  static const Color surfaceGray = Color(0xFF111111);

  // Primary Colors
  static const Color metallicGold = Color(0xFFD4AF37);
  static const Color metallicGoldDark = Color(0xFFB8941F);

  // Status Colors
  static const Color statusAvailable = Color(0xFF10B981); // 緑
  static const Color statusLowStock = Color(0xFFF59E0B); // 黄色
  static const Color statusOutOfStock = Color(0xFFEF4444); // 赤
  
  // Action Colors
  static const Color primaryBlue = Color(0xFF3B82F6); // プライマリブルー
  static const Color errorRed = Color(0xFFEF4444); // エラー赤（statusOutOfStockと同じ）
  static const Color successGreen = Color(0xFF10B981); // 成功緑（statusAvailableと同じ）

  // Text Colors
  static const Color textPrimary = Color(0xFFE5E5E5);
  static const Color textSecondary = Color(0xFF9CA3AF);
  static const Color textTertiary = Color(0xFF6B7280);

  // Border Colors
  static const Color borderGray = Color(0xFF374151);
  static const Color borderLight = Color(0xFF4B5563);
}

/// エンタープライズテーマ
ThemeData enterpriseTheme() {
  return ThemeData(
    useMaterial3: true,
    brightness: Brightness.dark,
    scaffoldBackgroundColor: EnterpriseColors.deepBlack,
    
    colorScheme: ColorScheme.dark(
      primary: EnterpriseColors.metallicGold,
      secondary: EnterpriseColors.metallicGoldDark,
      surface: EnterpriseColors.regalisBlack,
      background: EnterpriseColors.deepBlack,
      error: EnterpriseColors.statusOutOfStock,
    ),
    
    textTheme: const TextTheme(
      displayLarge: TextStyle(
        fontSize: 32,
        fontWeight: FontWeight.bold,
        color: EnterpriseColors.textPrimary,
        letterSpacing: -0.5,
      ),
      headlineMedium: TextStyle(
        fontSize: 24,
        fontWeight: FontWeight.bold,
        color: EnterpriseColors.textPrimary,
        letterSpacing: -0.3,
      ),
      bodyLarge: TextStyle(
        fontSize: 16,
        color: EnterpriseColors.textPrimary,
        height: 1.5,
      ),
      bodySmall: TextStyle(
        fontSize: 12,
        color: EnterpriseColors.textSecondary,
      ),
      labelSmall: TextStyle(
        fontSize: 10,
        color: EnterpriseColors.textTertiary,
        letterSpacing: 0.2,
      ),
    ),
    
    cardTheme: CardThemeData(
      color: EnterpriseColors.regalisBlack,
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
        side: const BorderSide(
          color: EnterpriseColors.borderGray,
          width: 1,
        ),
      ),
    ),
    
    elevatedButtonTheme: ElevatedButtonThemeData(
      style: ElevatedButton.styleFrom(
        backgroundColor: EnterpriseColors.metallicGold,
        foregroundColor: Colors.black,
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        elevation: 0,
      ),
    ),
    
    inputDecorationTheme: InputDecorationTheme(
      filled: true,
      fillColor: EnterpriseColors.regalisBlack,
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(8),
        borderSide: const BorderSide(color: EnterpriseColors.borderGray),
      ),
      enabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(8),
        borderSide: const BorderSide(color: EnterpriseColors.borderGray),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(8),
        borderSide: const BorderSide(
          color: EnterpriseColors.metallicGold,
          width: 2,
        ),
      ),
      hintStyle: const TextStyle(color: EnterpriseColors.textTertiary),
    ),
  );
}

