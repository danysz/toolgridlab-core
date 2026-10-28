import 'package:flutter/material.dart';
import 'package:sentry_flutter/sentry_flutter.dart';
import 'package:supabase_flutter/supabase_flutter.dart';

import 'config/app_config.dart';
import 'theme/app_theme.dart';
import 'auth/auth_service.dart';
import 'api/core_api_client.dart';
import 'screens/home_screen.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();

  final config = AppConfig.fromEnvironment();

  // Initialize Supabase if credentials are provided
  var isSupabaseConfigured = false;
  if (config.supabaseUrl.isNotEmpty && config.supabaseAnonKey.isNotEmpty) {
    try {
      await Supabase.initialize(
        url: config.supabaseUrl,
        // ignore: deprecated_member_use
        anonKey: config.supabaseAnonKey,
      );
      isSupabaseConfigured = true;
    } catch (e) {
      debugPrint('[SUPABASE] Initialization skipped: $e');
    }
  }

  final authService = AuthService(isSupabaseConfigured: isSupabaseConfigured);
  final apiClient = CoreApiClient(
    baseUrl: config.apiBaseUrl,
    getAccessToken: () => authService.accessToken,
  );

  final app = ToolGridLabApp(
    config: config,
    authService: authService,
    apiClient: apiClient,
  );

  // Initialize Sentry if configured, or run directly
  if (config.sentryDsn.isNotEmpty) {
    await SentryFlutter.init(
      (options) {
        options.dsn = config.sentryDsn;
        options.environment = config.environment;
        options.tracesSampleRate = 0.2;
      },
      appRunner: () => runApp(app),
    );
  } else {
    runApp(app);
  }
}

class ToolGridLabApp extends StatelessWidget {
  final AppConfig config;
  final AuthService authService;
  final CoreApiClient apiClient;

  const ToolGridLabApp({
    super.key,
    required this.config,
    required this.authService,
    required this.apiClient,
  });

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: config.appName,
      debugShowCheckedModeBanner: false,
      theme: AppTheme.darkTheme(),
      home: HomeScreen(
        config: config,
        authService: authService,
        apiClient: apiClient,
      ),
    );
  }
}
