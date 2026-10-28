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
    const rawApiBaseUrl = String.fromEnvironment('API_BASE_URL', defaultValue: 'http://localhost:8080');
    const supabaseUrl = String.fromEnvironment('SUPABASE_URL', defaultValue: '');
    const supabaseAnonKey = String.fromEnvironment('SUPABASE_ANON_KEY', defaultValue: '');
    const sentryDsn = String.fromEnvironment('SENTRY_DSN', defaultValue: '');
    const rawEnv = String.fromEnvironment('APP_ENV', defaultValue: 'development');

    // Runtime hostname awareness: if served on staging.toolgridlab.com or *.pages.dev preview,
    // ensure staging API and environment are used even if a production build artifact is served.
    final host = Uri.base.host.toLowerCase();
    final isStagingHost = host.startsWith('staging.') || host.endsWith('.pages.dev');
    final env = isStagingHost && rawEnv == 'production' ? 'staging' : rawEnv;
    final apiBaseUrl = (isStagingHost && rawApiBaseUrl == 'https://api.toolgridlab.com')
        ? 'https://api-staging.toolgridlab.com'
        : rawApiBaseUrl;

    return AppConfig(
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
