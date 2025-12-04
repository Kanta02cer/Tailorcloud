/// 環境変数管理クラス
///
/// 開発環境と本番環境の設定を管理します。
/// 環境変数は以下の方法で設定できます：
/// - ビルド時: `flutter run --dart-define=ENV=production --dart-define=API_BASE_URL=https://api.example.com`
/// - 実行時: 環境変数から読み込み
class Environment {
  /// 環境タイプ
  static const String env = String.fromEnvironment(
    'ENV',
    defaultValue: 'development',
  );

  /// 本番環境かどうか
  static bool get isProduction => env == 'production';

  /// 開発環境かどうか
  static bool get isDevelopment => !isProduction;

  /// ステージング環境かどうか
  static bool get isStaging => env == 'staging';

  /// APIベースURL
  static const String apiBaseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8080',
  );

  /// Firebase設定を有効にするかどうか
  static const bool enableFirebase = bool.fromEnvironment(
    'ENABLE_FIREBASE',
    defaultValue: false,
  );

  /// Firebase Web API Key（オプショナル）
  static const String firebaseApiKey = String.fromEnvironment(
    'FIREBASE_API_KEY',
    defaultValue: '',
  );

  /// Firebase App ID（オプショナル）
  static const String firebaseAppId = String.fromEnvironment(
    'FIREBASE_APP_ID',
    defaultValue: '',
  );

  /// Firebase Project ID（オプショナル）
  static const String firebaseProjectId = String.fromEnvironment(
    'FIREBASE_PROJECT_ID',
    defaultValue: '',
  );

  /// Firebase Messaging Sender ID（オプショナル）
  static const String firebaseMessagingSenderId = String.fromEnvironment(
    'FIREBASE_MESSAGING_SENDER_ID',
    defaultValue: '',
  );

  /// デフォルトテナントID
  static const String defaultTenantId = String.fromEnvironment(
    'DEFAULT_TENANT_ID',
    defaultValue: 'tenant-123',
  );

  /// デバッグログを有効にするかどうか
  static bool get enableDebugLogging => isDevelopment;

  /// 環境情報を文字列で取得
  static String get info => '''
Environment: $env
API Base URL: $apiBaseUrl
Firebase Enabled: $enableFirebase
Debug Logging: $enableDebugLogging
''';
}
