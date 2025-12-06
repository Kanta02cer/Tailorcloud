import 'package:flutter/foundation.dart';
import '../config/environment.dart';

/// ログレベル
enum LogLevel {
  debug,
  info,
  warning,
  error,
}

/// ログ管理サービス
///
/// 本番環境と開発環境で適切なログ出力を管理します。
class Logger {
  /// デバッグログを出力
  static void debug(String message, [Object? error, StackTrace? stackTrace]) {
    if (Environment.enableDebugLogging) {
      _log(LogLevel.debug, message, error, stackTrace);
    }
  }

  /// 情報ログを出力
  static void info(String message, [Object? error, StackTrace? stackTrace]) {
    _log(LogLevel.info, message, error, stackTrace);
  }

  /// 警告ログを出力
  static void warning(String message, [Object? error, StackTrace? stackTrace]) {
    _log(LogLevel.warning, message, error, stackTrace);
  }

  /// エラーログを出力
  static void error(String message, [Object? error, StackTrace? stackTrace]) {
    _log(LogLevel.error, message, error, stackTrace);
  }

  /// ログを出力
  static void _log(
    LogLevel level,
    String message,
    Object? error,
    StackTrace? stackTrace,
  ) {
    final timestamp = DateTime.now().toIso8601String();
    final levelString = level.name.toUpperCase();

    if (kDebugMode) {
      // デバッグモードでは詳細なログを出力
      final logMessage = '[$timestamp] [$levelString] $message';

      switch (level) {
        case LogLevel.debug:
          debugPrint(logMessage);
          break;
        case LogLevel.info:
          debugPrint(logMessage);
          break;
        case LogLevel.warning:
          debugPrint('⚠️ $logMessage');
          break;
        case LogLevel.error:
          debugPrint('❌ $logMessage');
          if (error != null) {
            debugPrint('Error: $error');
          }
          if (stackTrace != null) {
            debugPrint('Stack trace: $stackTrace');
          }
          break;
      }
    } else {
      // リリースモードでは重要なログのみ出力
      if (level == LogLevel.error || level == LogLevel.warning) {
        // 本番環境では外部ログサービスに送信することを推奨
        // 例: Sentry, Firebase Crashlytics, etc.
        debugPrint('[$levelString] $message');
        if (error != null) {
          debugPrint('Error: $error');
        }
      }
    }
  }

  /// APIリクエストログを出力
  static void logApiRequest(String method, String url,
      {Map<String, dynamic>? body}) {
    if (Environment.enableDebugLogging) {
      debug('API Request: $method $url');
      if (body != null) {
        debug('Request Body: $body');
      }
    }
  }

  /// APIレスポンスログを出力
  static void logApiResponse(String method, String url, int statusCode,
      {Object? body}) {
    if (Environment.enableDebugLogging) {
      debug('API Response: $method $url -> $statusCode');
      if (body != null) {
        debug('Response Body: $body');
      }
    }
  }

  /// APIエラーログを出力
  static void logApiError(
      String method, String url, int statusCode, String message) {
    error('API Error: $method $url -> $statusCode: $message');
  }
}
