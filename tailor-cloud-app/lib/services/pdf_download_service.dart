import 'dart:io';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:path_provider/path_provider.dart';
import 'package:permission_handler/permission_handler.dart';
import '../config/app_config.dart';

/// PDFダウンロードサービス
class PdfDownloadService {
  /// PDFをダウンロードしてローカルに保存
  /// 戻り値: 保存されたファイルのパス
  static Future<String> downloadPdf({
    required String url,
    required String fileName,
  }) async {
    try {
      // ストレージ権限をリクエスト
      final status = await Permission.storage.request();
      if (!status.isGranted) {
        throw Exception('ストレージ権限が必要です');
      }

      // PDFをダウンロード
      final response = await http.get(Uri.parse(url)).timeout(
        const Duration(seconds: 30),
      );

      if (response.statusCode != 200) {
        throw Exception('PDFのダウンロードに失敗しました: ${response.statusCode}');
      }

      // ダウンロードディレクトリを取得
      Directory? directory;
      if (Platform.isAndroid) {
        // Android: Downloadsディレクトリ
        directory = Directory('/storage/emulated/0/Download');
        if (!await directory.exists()) {
          // フォールバック: アプリのドキュメントディレクトリ
          directory = await getApplicationDocumentsDirectory();
        }
      } else if (Platform.isIOS) {
        // iOS: アプリのドキュメントディレクトリ
        directory = await getApplicationDocumentsDirectory();
      } else {
        // その他のプラットフォーム
        directory = await getApplicationDocumentsDirectory();
      }

      // ファイル名をクリーンアップ（無効な文字を削除）
      final cleanFileName = fileName.replaceAll(RegExp(r'[<>:"/\\|?*]'), '_');
      final filePath = '${directory.path}/$cleanFileName';

      // ファイルに保存
      final file = File(filePath);
      await file.writeAsBytes(response.bodyBytes);

      if (AppConfig.enableDebugLogging) {
        debugPrint('PDF downloaded to: $filePath');
      }

      return filePath;
    } catch (e) {
      if (AppConfig.enableDebugLogging) {
        debugPrint('PDF download error: $e');
      }
      rethrow;
    }
  }

  /// PDFを開く（外部アプリで開く）
  static Future<void> openPdf(String filePath) async {
    // url_launcherを使用してPDFを開く
    // 実装は呼び出し側で行う（url_launcherを使用）
  }
}

