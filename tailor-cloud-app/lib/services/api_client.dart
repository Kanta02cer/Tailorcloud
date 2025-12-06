import 'dart:async';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:firebase_auth/firebase_auth.dart';
import '../config/app_config.dart';
import '../services/logger.dart';

/// API例外クラス
class ApiException implements Exception {
  final int statusCode;
  final String message;

  ApiException({required this.statusCode, required this.message});

  @override
  String toString() => 'ApiException($statusCode): $message';
}

/// TailorCloud APIクライアント
class ApiClient {
  final String baseUrl;
  FirebaseAuth? _auth;

  ApiClient({String? baseUrl}) : baseUrl = baseUrl ?? AppConfig.baseUrl {
    // Firebaseが有効な場合のみ初期化
    if (AppConfig.enableFirebase) {
      try {
        _auth = FirebaseAuth.instance;
      } catch (e) {
        Logger.warning('Firebase Auth not available: $e');
        _auth = null;
      }
    }
  }

  /// IDトークンを取得
  Future<String?> _getIdToken() async {
    if (!AppConfig.enableFirebase || _auth == null) {
      return null;
    }

    try {
      final user = _auth!.currentUser;
      if (user == null) return null;
      return await user.getIdToken();
    } catch (e) {
      Logger.warning('Failed to get ID token: $e');
      return null;
    }
  }

  /// リクエストヘッダーを取得
  Future<Map<String, String>> _headers() async {
    final token = await _getIdToken();
    return {
      'Content-Type': 'application/json',
      if (token != null) 'Authorization': 'Bearer $token',
    };
  }

  /// GETリクエスト
  Future<T> get<T>(
    String path, {
    Map<String, String>? queryParameters,
    T Function(Map<String, dynamic>)? fromJson,
  }) async {
    final uri =
        Uri.parse('$baseUrl$path').replace(queryParameters: queryParameters);

    Logger.logApiRequest('GET', uri.toString());

    try {
      final response = await http
          .get(
            uri,
            headers: await _headers(),
          )
          .timeout(AppConfig.apiTimeout);

      Logger.logApiResponse('GET', uri.toString(), response.statusCode);
      return _handleResponse<T>(response, fromJson);
    } on http.ClientException catch (e) {
      Logger.logApiError('GET', uri.toString(), 0, 'ネットワークエラー: ${e.message}');
      throw ApiException(
        statusCode: 0,
        message: 'ネットワークエラーが発生しました。接続を確認してください。',
      );
    } on TimeoutException catch (e) {
      Logger.logApiError('GET', uri.toString(), 0, 'タイムアウト: ${e.message}');
      throw ApiException(
        statusCode: 0,
        message: 'リクエストがタイムアウトしました。しばらく待ってから再試行してください。',
      );
    } catch (e) {
      if (e is ApiException) rethrow;
      Logger.logApiError('GET', uri.toString(), 0, e.toString());
      throw ApiException(
        statusCode: 0,
        message: '予期しないエラーが発生しました: ${e.toString()}',
      );
    }
  }

  /// POSTリクエスト
  Future<T> post<T>(
    String path,
    Map<String, dynamic> body, {
    T Function(Map<String, dynamic>)? fromJson,
  }) async {
    final uri = Uri.parse('$baseUrl$path');
    Logger.logApiRequest('POST', uri.toString(), body: body);

    try {
      final response = await http
          .post(
            uri,
            headers: await _headers(),
            body: jsonEncode(body),
          )
          .timeout(AppConfig.apiTimeout);

      Logger.logApiResponse('POST', uri.toString(), response.statusCode);
      return _handleResponse<T>(response, fromJson);
    } on http.ClientException catch (e) {
      Logger.logApiError('POST', uri.toString(), 0, 'ネットワークエラー: ${e.message}');
      throw ApiException(
        statusCode: 0,
        message: 'ネットワークエラーが発生しました。接続を確認してください。',
      );
    } on TimeoutException catch (e) {
      Logger.logApiError('POST', uri.toString(), 0, 'タイムアウト: ${e.message}');
      throw ApiException(
        statusCode: 0,
        message: 'リクエストがタイムアウトしました。しばらく待ってから再試行してください。',
      );
    } catch (e) {
      if (e is ApiException) rethrow;
      Logger.logApiError('POST', uri.toString(), 0, e.toString());
      throw ApiException(
        statusCode: 0,
        message: '予期しないエラーが発生しました: ${e.toString()}',
      );
    }
  }

  /// PUTリクエスト
  Future<T> put<T>(
    String path,
    Map<String, dynamic> body, {
    T Function(Map<String, dynamic>)? fromJson,
  }) async {
    final uri = Uri.parse('$baseUrl$path');
    Logger.logApiRequest('PUT', uri.toString(), body: body);

    try {
      final response = await http
          .put(
            uri,
            headers: await _headers(),
            body: jsonEncode(body),
          )
          .timeout(AppConfig.apiTimeout);

      Logger.logApiResponse('PUT', uri.toString(), response.statusCode);
      return _handleResponse<T>(response, fromJson);
    } on http.ClientException catch (e) {
      Logger.logApiError('PUT', uri.toString(), 0, 'ネットワークエラー: ${e.message}');
      throw ApiException(
        statusCode: 0,
        message: 'ネットワークエラーが発生しました。接続を確認してください。',
      );
    } on TimeoutException catch (e) {
      Logger.logApiError('PUT', uri.toString(), 0, 'タイムアウト: ${e.message}');
      throw ApiException(
        statusCode: 0,
        message: 'リクエストがタイムアウトしました。しばらく待ってから再試行してください。',
      );
    } catch (e) {
      if (e is ApiException) rethrow;
      Logger.logApiError('PUT', uri.toString(), 0, e.toString());
      throw ApiException(
        statusCode: 0,
        message: '予期しないエラーが発生しました: ${e.toString()}',
      );
    }
  }

  /// レスポンスを処理
  T _handleResponse<T>(
    http.Response response,
    T Function(Map<String, dynamic>)? fromJson,
  ) {
    if (response.statusCode >= 200 && response.statusCode < 300) {
      try {
        final data = jsonDecode(response.body);

        if (fromJson != null) {
          if (data is Map<String, dynamic>) {
            return fromJson(data);
          } else if (data is List) {
            // リストの場合はエラー
            throw ApiException(
              statusCode: response.statusCode,
              message: 'Expected object but got array',
            );
          }
        }

        return data as T;
      } catch (e) {
        Logger.error('Failed to parse response: $e');
        throw ApiException(
          statusCode: response.statusCode,
          message: 'Failed to parse response: $e',
        );
      }
    } else {
      String message = response.body;
      try {
        final errorData = jsonDecode(response.body);
        message = errorData['error'] ?? response.body;
      } catch (e) {
        // JSONパースエラーはそのままbodyを使用
      }

      Logger.logApiError('API', response.request?.url.toString() ?? '',
          response.statusCode, message);
      throw ApiException(
        statusCode: response.statusCode,
        message: message,
      );
    }
  }
}
