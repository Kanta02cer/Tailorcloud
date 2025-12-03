import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:firebase_auth/firebase_auth.dart';
import '../config/app_config.dart';

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
  final FirebaseAuth _auth = FirebaseAuth.instance;

  ApiClient({String? baseUrl})
      : baseUrl = baseUrl ?? AppConfig.baseUrl;

  /// IDトークンを取得
  Future<String?> _getIdToken() async {
    try {
      final user = _auth.currentUser;
      if (user == null) return null;
      return await user.getIdToken();
    } catch (e) {
      if (AppConfig.enableDebugLogging) {
        print('Failed to get ID token: $e');
      }
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
    final uri = Uri.parse('$baseUrl$path')
        .replace(queryParameters: queryParameters);

    if (AppConfig.enableDebugLogging) {
      print('GET $uri');
    }

    final response = await http
        .get(
          uri,
          headers: await _headers(),
        )
        .timeout(AppConfig.apiTimeout);

    return _handleResponse<T>(response, fromJson);
  }

  /// POSTリクエスト
  Future<T> post<T>(
    String path,
    Map<String, dynamic> body, {
    T Function(Map<String, dynamic>)? fromJson,
  }) async {
    if (AppConfig.enableDebugLogging) {
      print('POST $baseUrl$path');
      print('Body: ${jsonEncode(body)}');
    }

    final response = await http
        .post(
          Uri.parse('$baseUrl$path'),
          headers: await _headers(),
          body: jsonEncode(body),
        )
        .timeout(AppConfig.apiTimeout);

    return _handleResponse<T>(response, fromJson);
  }

  /// PUTリクエスト
  Future<T> put<T>(
    String path,
    Map<String, dynamic> body, {
    T Function(Map<String, dynamic>)? fromJson,
  }) async {
    if (AppConfig.enableDebugLogging) {
      print('PUT $baseUrl$path');
      print('Body: ${jsonEncode(body)}');
    }

    final response = await http
        .put(
          Uri.parse('$baseUrl$path'),
          headers: await _headers(),
          body: jsonEncode(body),
        )
        .timeout(AppConfig.apiTimeout);

    return _handleResponse<T>(response, fromJson);
  }

  /// レスポンスを処理
  T _handleResponse<T>(
    http.Response response,
    T Function(Map<String, dynamic>)? fromJson,
  ) {
    if (response.statusCode >= 200 && response.statusCode < 300) {
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
    } else {
      String message = response.body;
      try {
        final errorData = jsonDecode(response.body);
        message = errorData['error'] ?? response.body;
      } catch (e) {
        // JSONパースエラーはそのままbodyを使用
      }

      throw ApiException(
        statusCode: response.statusCode,
        message: message,
      );
    }
  }
}

