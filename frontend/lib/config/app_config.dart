class AppConfig {
  final String appName;
  final String apiBaseUrl;
  final String supabaseUrl;
  final String supabaseAnonKey;
  final String sentryDsn;
  final String environment;
  final bool isProduction;

  const AppConfig({
    required this.appName,
    required this.apiBaseUrl,
    required this.supabaseUrl,
    required this.supabaseAnonKey,
    required this.sentryDsn,
    this.environment = 'development',
    required this.isProduction,
  });

  // White-label configuration loader via compile-time environment flags or local defaults
  factory AppConfig.fromEnvironment() {
    const appName = String.fromEnvironment('APP_NAME', defaultValue: 'ToolGridLab');
    const apiBaseUrl = String.fromEnvironment('API_BASE_URL', defaultValue: 'http://localhost:8080');
    const supabaseUrl = String.fromEnvironment('SUPABASE_URL', defaultValue: '');
    const supabaseAnonKey = String.fromEnvironment('SUPABASE_ANON_KEY', defaultValue: '');
    const sentryDsn = String.fromEnvironment('SENTRY_DSN', defaultValue: '');
    const env = String.fromEnvironment('APP_ENV', defaultValue: 'development');

    return const AppConfig(
      appName: appName,
      apiBaseUrl: apiBaseUrl,
      supabaseUrl: supabaseUrl,
      supabaseAnonKey: supabaseAnonKey,
      sentryDsn: sentryDsn,
      environment: env,
      isProduction: env == 'production',
    );
  }
}
