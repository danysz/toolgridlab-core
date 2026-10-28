import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:frontend/main.dart';
import 'package:frontend/config/app_config.dart';
import 'package:frontend/auth/auth_service.dart';
import 'package:frontend/api/core_api_client.dart';

void main() {
  testWidgets('ToolGridLab Core smoke test', (WidgetTester tester) async {
    tester.view.physicalSize = const Size(1280, 900);
    tester.view.devicePixelRatio = 1.0;
    addTearDown(tester.view.resetPhysicalSize);
    addTearDown(tester.view.resetDevicePixelRatio);

    const config = AppConfig(
      appName: 'ToolGridLab',
      apiBaseUrl: 'http://localhost:8080',
      supabaseUrl: '',
      supabaseAnonKey: '',
      sentryDsn: '',
      isProduction: false,
    );

    final authService = AuthService(isSupabaseConfigured: false);
    final apiClient = CoreApiClient(
      baseUrl: config.apiBaseUrl,
      getAccessToken: () => null,
    );

    await tester.pumpWidget(
      ToolGridLabApp(
        config: config,
        authService: authService,
        apiClient: apiClient,
      ),
    );

    // Let any async frame stabilize
    await tester.pumpAndSettle();

    expect(find.text('ToolGridLab'), findsWidgets);
    expect(find.text('OPEN CORE'), findsOneWidget);
    expect(find.text('Modular Monolith Engine\nfor Micro-SaaS Portfolios'), findsOneWidget);
  });
}
