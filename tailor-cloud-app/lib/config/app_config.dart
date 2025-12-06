import 'environment.dart';

/// TailorCloud アプリケーション設定
class AppConfig {
  // API Base URL（環境変数から取得）
  static String get baseUrl => Environment.apiBaseUrl;

  // 環境設定
  static bool get isProduction => Environment.isProduction;
  static bool get isDevelopment => Environment.isDevelopment;
  static bool get isStaging => Environment.isStaging;

  // デバッグ設定
  static bool get enableDebugLogging => Environment.enableDebugLogging;

  // APIタイムアウト
  static const Duration apiTimeout = Duration(seconds: 30);

  // 画像キャッシュ期間
  static const Duration imageCacheDuration = Duration(days: 7);

  // オフラインストレージ設定
  static const String fabricBoxName = 'fabrics';
  static const String orderBoxName = 'orders';

  // デフォルトテナントID
  static String get defaultTenantId => Environment.defaultTenantId;

  // Firebase設定
  static bool get enableFirebase => Environment.enableFirebase;

  /// アプリケーション設定情報を取得
  static String get info => '''
App Configuration:
${Environment.info}
API Timeout: ${apiTimeout.inSeconds}s
Image Cache Duration: ${imageCacheDuration.inDays} days
''';
}
