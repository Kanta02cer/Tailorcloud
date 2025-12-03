/// TailorCloud アプリケーション設定
class AppConfig {
  // API Base URL
  static const String baseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8080',
  );

  // 環境設定
  static const bool isProduction = bool.fromEnvironment('dart.vm.product');
  static const bool isDevelopment = !isProduction;

  // デバッグ設定
  static const bool enableDebugLogging = isDevelopment;

  // APIタイムアウト
  static const Duration apiTimeout = Duration(seconds: 30);

  // 画像キャッシュ期間
  static const Duration imageCacheDuration = Duration(days: 7);

  // オフラインストレージ設定
  static const String fabricBoxName = 'fabrics';
  static const String orderBoxName = 'orders';
}

