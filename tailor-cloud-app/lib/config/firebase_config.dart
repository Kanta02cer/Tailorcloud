import 'package:firebase_core/firebase_core.dart';
import 'package:flutter/foundation.dart';
import '../services/logger.dart';
import 'environment.dart';

/// Firebase設定クラス
///
/// Firebaseの初期化設定を管理します。
/// 環境変数から設定を読み込み、Firebaseが有効な場合のみ初期化します。
class FirebaseConfig {
  /// Firebaseが有効かどうか
  static bool get isEnabled => Environment.enableFirebase;

  /// Firebase設定が完全かどうか
  static bool get isConfigured {
    if (!isEnabled) return false;

    return Environment.firebaseApiKey.isNotEmpty &&
        Environment.firebaseAppId.isNotEmpty &&
        Environment.firebaseProjectId.isNotEmpty;
  }

  /// Firebaseオプションを作成
  ///
  /// 環境変数からFirebase設定を読み込み、FirebaseOptionsを作成します。
  /// 設定が不完全な場合はnullを返します。
  static FirebaseOptions? createOptions() {
    if (!isEnabled || !isConfigured) {
      return null;
    }

    return FirebaseOptions(
      apiKey: Environment.firebaseApiKey,
      appId: Environment.firebaseAppId,
      messagingSenderId: Environment.firebaseMessagingSenderId,
      projectId: Environment.firebaseProjectId,
      // Web用の設定
      // 環境変数で指定されていない場合は自動生成
      authDomain: Environment.firebaseAuthDomain.isNotEmpty
          ? Environment.firebaseAuthDomain
          : '${Environment.firebaseProjectId}.firebaseapp.com',
      storageBucket: Environment.firebaseStorageBucket.isNotEmpty
          ? Environment.firebaseStorageBucket
          : '${Environment.firebaseProjectId}.appspot.com',
      // Analytics用（オプショナル）
      measurementId: Environment.firebaseMeasurementId.isNotEmpty
          ? Environment.firebaseMeasurementId
          : null,
    );
  }

  /// Firebaseを初期化
  ///
  /// 設定が有効な場合のみFirebaseを初期化します。
  /// 初期化に失敗した場合はfalseを返します。
  /// Web環境では既に初期化されている場合をチェックします。
  static Future<bool> initialize() async {
    // Firebaseが無効な場合はスキップ
    if (!isEnabled) {
      Logger.debug('Firebase: Disabled. Skipping initialization.');
      return false;
    }

    // 設定が不完全な場合はスキップ
    if (!isConfigured) {
      Logger.warning(
          'Firebase: Configuration incomplete. Skipping initialization.');
      Logger.debug('Firebase: Required settings - API_KEY, APP_ID, PROJECT_ID');
      return false;
    }

    try {
      // Web環境では既に初期化されているかチェック
      if (kIsWeb) {
        try {
          // 既に初期化されている場合は成功として扱う
          Firebase.app();
          Logger.info('Firebase: Already initialized.');
          return true;
        } catch (e) {
          // 初期化されていない場合は続行
        }
      }

      final options = createOptions();
      if (options == null) {
        Logger.warning(
            'Firebase: Options creation failed. Skipping initialization.');
        return false;
      }

      await Firebase.initializeApp(options: options);
      Logger.info('Firebase: Initialized successfully.');
      return true;
    } catch (e, stackTrace) {
      // エラーをログに記録するが、アプリの起動は続行
      Logger.warning('Firebase: Initialization failed: $e');
      if (kDebugMode) {
        Logger.debug('Firebase: Stack trace: $stackTrace');
      }
      return false;
    }
  }
}
