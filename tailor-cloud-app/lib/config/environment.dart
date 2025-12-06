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

  /// Firebase Auth Domain（オプショナル、未設定時は自動生成）
  static const String firebaseAuthDomain = String.fromEnvironment(
    'FIREBASE_AUTH_DOMAIN',
    defaultValue: '',
  );

  /// Firebase Storage Bucket（オプショナル、未設定時は自動生成）
  static const String firebaseStorageBucket = String.fromEnvironment(
    'FIREBASE_STORAGE_BUCKET',
    defaultValue: '',
  );

  /// Firebase Measurement ID（Analytics用、オプショナル）
  static const String firebaseMeasurementId = String.fromEnvironment(
    'FIREBASE_MEASUREMENT_ID',
    defaultValue: '',
  );

  /// デフォルトテナントID
  static const String defaultTenantId = String.fromEnvironment(
    'DEFAULT_TENANT_ID',
    defaultValue: 'tenant-123',
  );

  /// 許可するGoogle Workspaceドメイン（例: example.com）
  /// 空文字列の場合はドメイン制限を行わない
  static const String googleWorkspaceDomain = String.fromEnvironment(
    'GOOGLE_WORKSPACE_DOMAIN',
    defaultValue: '',
  );

  /// デバッグログを有効にするかどうか
  static bool get enableDebugLogging => isDevelopment;

  /// 環境情報を文字列で取得
  static String get info => '''
Environment: $env
API Base URL: $apiBaseUrl
Firebase Enabled: $enableFirebase
Google Workspace Domain: $googleWorkspaceDomain
Debug Logging: $enableDebugLogging
''';
}
