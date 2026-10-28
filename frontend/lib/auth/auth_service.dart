import 'package:flutter/foundation.dart';
import 'package:supabase_flutter/supabase_flutter.dart';

class AuthService extends ChangeNotifier {
  final bool isSupabaseConfigured;
  User? _currentUser;
  String? _accessToken;
  bool _isLoading = false;

  AuthService({required this.isSupabaseConfigured}) {
    if (isSupabaseConfigured) {
      _initSupabaseListener();
    }
  }

  User? get currentUser => _currentUser;
  String? get accessToken => _accessToken;
  bool get isAuthenticated => _currentUser != null;
  bool get isLoading => _isLoading;

  void _initSupabaseListener() {
    Supabase.instance.client.auth.onAuthStateChange.listen((data) {
      final session = data.session;
      _currentUser = session?.user;
      _accessToken = session?.accessToken;
      notifyListeners();
    });
  }

  // Google OAuth sign-in via Supabase OIDC
  Future<void> signInWithGoogle() async {
    _isLoading = true;
    notifyListeners();

    try {
      if (isSupabaseConfigured) {
        await Supabase.instance.client.auth.signInWithOAuth(
          OAuthProvider.google,
          redirectTo: kIsWeb ? null : 'com.toolgridlab.core://login-callback',
        );
      } else {
        // Fallback simulation for local development without configured Supabase credentials
        await Future.delayed(const Duration(milliseconds: 600));
        _currentUser = const User(
          id: 'mock-founder-id',
          appMetadata: {'provider': 'google'},
          userMetadata: {'name': 'Demo Founder'},
          aud: 'authenticated',
          createdAt: '2026-10-26T00:00:00Z',
          email: 'founder@toolgridlab.com',
        );
        _accessToken = 'mock-jwt-token-local';
      }
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  // Apple OAuth sign-in via Supabase OIDC
  Future<void> signInWithApple() async {
    _isLoading = true;
    notifyListeners();

    try {
      if (isSupabaseConfigured) {
        await Supabase.instance.client.auth.signInWithOAuth(
          OAuthProvider.apple,
          redirectTo: kIsWeb ? null : 'com.toolgridlab.core://login-callback',
        );
      } else {
        await Future.delayed(const Duration(milliseconds: 600));
        _currentUser = const User(
          id: 'mock-apple-user-id',
          appMetadata: {'provider': 'apple'},
          userMetadata: {'name': 'Apple User'},
          aud: 'authenticated',
          createdAt: '2026-10-26T00:00:00Z',
          email: 'user@privaterelay.appleid.com',
        );
        _accessToken = 'mock-jwt-token-local';
      }
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> signOut() async {
    _isLoading = true;
    notifyListeners();

    try {
      if (isSupabaseConfigured) {
        await Supabase.instance.client.auth.signOut();
      }
      _currentUser = null;
      _accessToken = null;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }
}
