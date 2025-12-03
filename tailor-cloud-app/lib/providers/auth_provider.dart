import 'package:firebase_auth/firebase_auth.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'auth_provider.g.dart';

/// Firebase Authインスタンスプロバイダー
@riverpod
FirebaseAuth firebaseAuth(FirebaseAuthRef ref) {
  return FirebaseAuth.instance;
}

/// 現在のユーザープロバイダー
@riverpod
Stream<User?> authStateChanges(AuthStateChangesRef ref) {
  final auth = ref.watch(firebaseAuthProvider);
  return auth.authStateChanges();
}

/// 現在のユーザー（同期版）
@riverpod
User? currentUser(CurrentUserRef ref) {
  final auth = ref.watch(firebaseAuthProvider);
  return auth.currentUser;
}

/// ログイン関数
@riverpod
Future<void> signIn(
  SignInRef ref, {
  required String email,
  required String password,
}) async {
  final auth = ref.watch(firebaseAuthProvider);
  await auth.signInWithEmailAndPassword(
    email: email,
    password: password,
  );
}

/// ログアウト関数
@riverpod
Future<void> signOut(SignOutRef ref) async {
  final auth = ref.watch(firebaseAuthProvider);
  await auth.signOut();
}

