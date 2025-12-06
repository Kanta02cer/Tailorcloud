import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../config/app_config.dart';
import '../config/environment.dart';
import '../services/logger.dart';

part 'auth_provider.g.dart';

/// Firebase Authインスタンスプロバイダー
///
/// Firebaseが無効な場合はnullを返します。
@riverpod
FirebaseAuth? firebaseAuth(FirebaseAuthRef ref) {
  if (!AppConfig.enableFirebase) {
    return null;
  }
  try {
    return FirebaseAuth.instance;
  } catch (e) {
    Logger.warning('Firebase Auth not available: $e');
    return null;
  }
}

/// 現在のユーザープロバイダー
///
/// Firebaseが無効な場合は常にnullを返すストリームを返します。
@riverpod
Stream<User?> authStateChanges(AuthStateChangesRef ref) {
  final auth = ref.watch(firebaseAuthProvider);
  if (auth == null) {
    // Firebaseが無効な場合は、nullを返すストリームを返す
    return Stream.value(null);
  }
  try {
    return auth.authStateChanges();
  } catch (e) {
    Logger.warning('Failed to get auth state changes: $e');
    return Stream.value(null);
  }
}

/// 現在のユーザー（同期版）
///
/// Firebaseが無効な場合はnullを返します。
@riverpod
User? currentUser(CurrentUserRef ref) {
  final auth = ref.watch(firebaseAuthProvider);
  if (auth == null) {
    return null;
  }
  try {
    return auth.currentUser;
  } catch (e) {
    Logger.warning('Failed to get current user: $e');
    return null;
  }
}

/// ログイン関数
///
/// Firebaseが無効な場合はエラーをスローします。
@riverpod
Future<void> signIn(
  SignInRef ref, {
  required String email,
  required String password,
}) async {
  final auth = ref.watch(firebaseAuthProvider);
  if (auth == null) {
    throw Exception(
        'Firebase Auth is not enabled. Please enable Firebase to use authentication.');
  }

  try {
    await auth.signInWithEmailAndPassword(
      email: email,
      password: password,
    );
    Logger.info('User signed in successfully: $email');
  } catch (e) {
    Logger.error('Sign in failed: $e');
    rethrow;
  }
}

/// Googleサインイン（Web向け）
///
/// - Firebase Auth を利用して Google アカウントでログイン
/// - `GOOGLE_WORKSPACE_DOMAIN` が設定されている場合、そのドメインのみ許可
Future<void> signInWithGoogle(WidgetRef ref) async {
  if (!AppConfig.enableFirebase) {
    throw Exception(
      'Firebase Auth is not enabled. Please enable Firebase to use Google sign-in.',
    );
  }

  final auth = ref.read(firebaseAuthProvider);
  if (auth == null) {
    throw Exception('Firebase Auth client is not available.');
  }

  try {
    if (!kIsWeb) {
      throw UnsupportedError(
        'Googleサインインは現在Web環境のみサポートされています。',
      );
    }

    final provider = GoogleAuthProvider();
    await auth.signInWithPopup(provider);

    final user = auth.currentUser;
    final allowedDomain = Environment.googleWorkspaceDomain;

    if (allowedDomain.isNotEmpty && user?.email != null) {
      final email = user!.email!;
      final domain = email.split('@').last;
      if (domain.toLowerCase() != allowedDomain.toLowerCase()) {
        await auth.signOut();
        throw Exception(
          'このシステムは $allowedDomain ドメインのGoogleアカウントのみ利用できます。',
        );
      }
    }

    Logger.info('User signed in with Google: ${user?.email}');
  } catch (e) {
    Logger.error('Google sign in failed: $e');
    rethrow;
  }
}

/// ログアウト関数
///
/// Firebaseが無効な場合は何もしません。
@riverpod
Future<void> signOut(SignOutRef ref) async {
  final auth = ref.watch(firebaseAuthProvider);
  if (auth == null) {
    Logger.info('Firebase Auth is not enabled. Sign out skipped.');
    return;
  }

  try {
    await auth.signOut();
    Logger.info('User signed out successfully');
  } catch (e) {
    Logger.error('Sign out failed: $e');
    rethrow;
  }
}
