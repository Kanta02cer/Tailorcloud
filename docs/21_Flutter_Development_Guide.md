# TailorCloud Flutter開発ガイド

**作成日**: 2025-01  
**フェーズ**: Phase 1.2 - iPadアプリ開発準備

---

## 📋 目次

1. [開発環境セットアップ](#開発環境セットアップ)
2. [プロジェクト構成](#プロジェクト構成)
3. [デザインシステム](#デザインシステム)
4. [APIクライアント実装](#apiクライアント実装)
5. [状態管理](#状態管理)
6. [オフライン対応](#オフライン対応)

---

## 🛠️ 開発環境セットアップ

### 必要なツール

- Flutter SDK: 3.16.0以上
- Dart SDK: 3.2.0以上
- Xcode: 15.0以上（iOS開発用）
- Android Studio / VS Code

### プロジェクト作成

```bash
# Flutterプロジェクト作成
flutter create tailor_cloud_app
cd tailor_cloud_app

# 必要なパッケージを追加
flutter pub add \
  firebase_auth \
  firebase_core \
  http \
  riverpod \
  riverpod_annotation \
  freezed_annotation \
  json_annotation \
  hive \
  hive_flutter \
  cached_network_image \
  flutter_svg
```

### プロジェクト構成

```
tailor_cloud_app/
├── lib/
│   ├── main.dart
│   ├── config/
│   │   ├── app_config.dart
│   │   └── theme.dart
│   ├── models/
│   │   ├── fabric.dart
│   │   ├── order.dart
│   │   └── ambassador.dart
│   ├── services/
│   │   ├── api_client.dart
│   │   ├── auth_service.dart
│   │   └── storage_service.dart
│   ├── providers/
│   │   ├── auth_provider.dart
│   │   ├── fabric_provider.dart
│   │   └── order_provider.dart
│   ├── screens/
│   │   ├── home/
│   │   │   └── home_screen.dart
│   │   ├── inventory/
│   │   │   └── inventory_screen.dart
│   │   └── order/
│   │       └── order_create_screen.dart
│   ├── widgets/
│   │   ├── fabric_card.dart
│   │   ├── order_card.dart
│   │   └── kpi_card.dart
│   └── utils/
│       ├── constants.dart
│       └── validators.dart
├── assets/
│   ├── images/
│   └── icons/
└── pubspec.yaml
```

---

## 🎨 デザインシステム

### カラーパレット

**参考**: `docs/11_UI_Design_Specifications.md`

```dart
// lib/config/theme.dart
import 'package:flutter/material.dart';

class AppColors {
  // Primary Navy
  static const Color primaryNavy = Color(0xFF1A1F3A);
  static const Color primaryNavyDark = Color(0xFF0F1424);
  static const Color primaryNavyLight = Color(0xFF2A3054);
  
  // Accent Colors
  static const Color accentGold = Color(0xFFD4AF37);
  static const Color accentCream = Color(0xFFF5F1E8);
  
  // Status Colors
  static const Color statusAvailable = Color(0xFF10B981); // 緑
  static const Color statusLimited = Color(0xFFF59E0B); // 黄色
  static const Color statusSoldOut = Color(0xFFEF4444); // 赤
  static const Color statusWarning = Color(0xFFF59E0B);
  
  // Neutral Colors
  static const Color neutralGray100 = Color(0xFFF3F4F6);
  static const Color neutralGray200 = Color(0xFFE5E7EB);
  static const Color neutralGray500 = Color(0xFF6B7280);
  static const Color neutralGray900 = Color(0xFF111827);
  
  // Background
  static const Color backgroundWhite = Color(0xFFFFFFFF);
  static const Color backgroundCream = Color(0xFFF5F1E8);
}
```

### タイポグラフィ

```dart
class AppTextStyles {
  // H1 - 大見出し（ダッシュボードKPI）
  static const TextStyle h1 = TextStyle(
    fontSize: 32,
    fontWeight: FontWeight.bold,
    letterSpacing: -0.5,
    height: 1.2,
  );
  
  // H2 - 中見出し（画面タイトル）
  static const TextStyle h2 = TextStyle(
    fontSize: 24,
    fontWeight: FontWeight.bold,
    letterSpacing: -0.3,
    height: 1.3,
  );
  
  // Body - 本文
  static const TextStyle body = TextStyle(
    fontSize: 16,
    fontWeight: FontWeight.normal,
    letterSpacing: 0,
    height: 1.5,
  );
  
  // KPI Number - 数字表示用
  static const TextStyle kpiNumber = TextStyle(
    fontSize: 28,
    fontWeight: FontWeight.bold,
    letterSpacing: -1,
  );
  
  // Caption - キャプション
  static const TextStyle caption = TextStyle(
    fontSize: 12,
    fontWeight: FontWeight.normal,
    letterSpacing: 0.2,
  );
}
```

### テーマ設定

```dart
// lib/config/theme.dart
import 'package:flutter/material.dart';

ThemeData appTheme() {
  return ThemeData(
    primaryColor: AppColors.primaryNavy,
    scaffoldBackgroundColor: AppColors.backgroundCream,
    fontFamily: 'NotoSansJP', // 日本語フォント
    textTheme: TextTheme(
      headlineLarge: AppTextStyles.h1,
      headlineMedium: AppTextStyles.h2,
      bodyLarge: AppTextStyles.body,
      bodySmall: AppTextStyles.caption,
    ),
    colorScheme: ColorScheme.fromSeed(
      seedColor: AppColors.primaryNavy,
      primary: AppColors.primaryNavy,
      secondary: AppColors.accentGold,
    ),
  );
}
```

---

## 🔌 APIクライアント実装

### APIクライアントクラス

```dart
// lib/services/api_client.dart
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:firebase_auth/firebase_auth.dart';

class ApiException implements Exception {
  final int statusCode;
  final String message;
  
  ApiException({required this.statusCode, required this.message});
  
  @override
  String toString() => 'ApiException($statusCode): $message';
}

class ApiClient {
  final String baseUrl;
  final FirebaseAuth _auth = FirebaseAuth.instance;
  
  ApiClient({required this.baseUrl});
  
  Future<String?> _getIdToken() async {
    final user = _auth.currentUser;
    if (user == null) return null;
    return await user.getIdToken();
  }
  
  Future<Map<String, String>> _headers() async {
    final token = await _getIdToken();
    return {
      'Content-Type': 'application/json',
      if (token != null) 'Authorization': 'Bearer $token',
    };
  }
  
  Future<T> get<T>(
    String path, {
    Map<String, String>? queryParameters,
    T Function(Map<String, dynamic>)? fromJson,
  }) async {
    final uri = Uri.parse('$baseUrl$path')
        .replace(queryParameters: queryParameters);
    
    final response = await http.get(
      uri,
      headers: await _headers(),
    );
    
    return _handleResponse<T>(response, fromJson);
  }
  
  Future<T> post<T>(
    String path,
    Map<String, dynamic> body, {
    T Function(Map<String, dynamic>)? fromJson,
  }) async {
    final response = await http.post(
      Uri.parse('$baseUrl$path'),
      headers: await _headers(),
      body: jsonEncode(body),
    );
    
    return _handleResponse<T>(response, fromJson);
  }
  
  T _handleResponse<T>(
    http.Response response,
    T Function(Map<String, dynamic>)? fromJson,
  ) {
    if (response.statusCode >= 200 && response.statusCode < 300) {
      final data = jsonDecode(response.body) as Map<String, dynamic>;
      if (fromJson != null) {
        return fromJson(data);
      }
      return data as T;
    } else {
      throw ApiException(
        statusCode: response.statusCode,
        message: response.body,
      );
    }
  }
}
```

---

## 🔄 状態管理（Riverpod）

### 設定

```dart
// lib/main.dart
import 'package:flutter_riverpod/flutter_riverpod.dart';

void main() {
  runApp(
    const ProviderScope(
      child: TailorCloudApp(),
    ),
  );
}
```

### プロバイダー例

```dart
// lib/providers/fabric_provider.dart
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../models/fabric.dart';
import '../services/api_client.dart';

part 'fabric_provider.g.dart';

@riverpod
ApiClient apiClient(ApiClientRef ref) {
  return ApiClient(baseUrl: 'http://localhost:8080');
}

@riverpod
Future<List<Fabric>> fabricList(
  FabricListRef ref, {
  String? status,
  String? search,
}) async {
  final client = ref.watch(apiClientProvider);
  final response = await client.get<Map<String, dynamic>>(
    '/api/fabrics',
    queryParameters: {
      if (status != null) 'status': status,
      if (search != null) 'search': search,
    },
  );
  
  final fabrics = (response['fabrics'] as List)
      .map((json) => Fabric.fromJson(json))
      .toList();
  
  return fabrics;
}
```

---

## 📱 オフライン対応

### Hive設定

```dart
// lib/services/storage_service.dart
import 'package:hive_flutter/hive_flutter.dart';

class StorageService {
  static Future<void> init() async {
    await Hive.initFlutter();
    // ボックスの登録
    // await Hive.openBox<Fabric>('fabrics');
  }
}
```

### オフライン対応の実装パターン

```dart
@riverpod
Future<List<Fabric>> cachedFabricList(
  CachedFabricListRef ref,
) async {
  // 1. まずローカルストレージから取得
  final localBox = await Hive.openBox<Fabric>('fabrics');
  final localFabrics = localBox.values.toList();
  
  if (localFabrics.isNotEmpty) {
    // UIを更新（オフライン対応）
    ref.keepAlive();
  }
  
  try {
    // 2. サーバーから取得を試みる
    final remoteFabrics = await ref.watch(fabricListProvider().future);
    
    // 3. ローカルストレージを更新
    await localBox.clear();
    for (final fabric in remoteFabrics) {
      await localBox.put(fabric.id, fabric);
    }
    
    return remoteFabrics;
  } catch (e) {
    // ネットワークエラーの場合はローカルデータを返す
    if (localFabrics.isNotEmpty) {
      return localFabrics;
    }
    rethrow;
  }
}
```

---

## 📝 実装チェックリスト

### Phase 1.2準備

- [ ] プロジェクトセットアップ
- [ ] デザインシステム実装（カラーパレット、タイポグラフィ）
- [ ] APIクライアント実装
- [ ] モデルクラス実装（Fabric, Order, Ambassador）
- [ ] 状態管理設定（Riverpod）
- [ ] 認証統合（Firebase Auth）

### Phase 1.2実装

- [ ] Home（Dashboard）画面実装
- [ ] Inventory（生地一覧）画面実装
- [ ] Visual Ordering画面実装
- [ ] リアルタイム価格計算機能

---

**最終更新日**: 2025-01  
**バージョン**: 1.0.0

