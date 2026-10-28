import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import '../theme/app_theme.dart';
import '../config/app_config.dart';
import '../auth/auth_service.dart';
import '../api/core_api_client.dart';
import '../widgets/app_nav_bar.dart';
import '../widgets/architecture_card.dart';
import '../widgets/footer.dart';

class HomeScreen extends StatefulWidget {
  final AppConfig config;
  final AuthService authService;
  final CoreApiClient apiClient;

  const HomeScreen({
    super.key,
    required this.config,
    required this.authService,
    required this.apiClient,
  });

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  bool _isBackendHealthy = false;
  Map<String, dynamic>? _healthData;
  bool _isCheckingHealth = false;

  // Live Email Verification Test state
  final TextEditingController _emailController = TextEditingController();
  bool _isVerifyingEmail = false;
  Map<String, dynamic>? _verificationResult;

  @override
  void initState() {
    super.initState();
    _checkApiHealth();
  }

  Future<void> _checkApiHealth() async {
    setState(() => _isCheckingHealth = true);
    final data = await widget.apiClient.checkHealth();
    setState(() {
      _healthData = data;
      _isBackendHealthy = data['status'] == 'healthy';
      _isCheckingHealth = false;
    });
  }

  Future<void> _verifyEmail() async {
    final email = _emailController.text.trim();
    if (email.isEmpty) return;

    setState(() {
      _isVerifyingEmail = true;
      _verificationResult = null;
    });

    final res = await widget.apiClient.verifyEmail(email);

    setState(() {
      _isVerifyingEmail = false;
      _verificationResult = res;
    });
  }

  @override
  void dispose() {
    _emailController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final isDesktop = MediaQuery.of(context).size.width >= 900;

    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            AppNavBar(
              appName: widget.config.appName,
              isBackendHealthy: _isBackendHealthy,
              authService: widget.authService,
            ),

            Center(
              child: ConstrainedBox(
                constraints: const BoxConstraints(maxWidth: 1100),
                child: Padding(
                  padding: EdgeInsets.symmetric(
                    horizontal: 24,
                    vertical: isDesktop ? 60 : 36,
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.center,
                    children: [
                      // Tagline badge
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 6),
                        decoration: BoxDecoration(
                          color: AppTheme.surface,
                          borderRadius: BorderRadius.circular(20),
                          border: Border.all(color: AppTheme.cardBorder),
                        ),
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            const Icon(Icons.bolt, size: 16, color: AppTheme.primary),
                            const SizedBox(width: 8),
                            Text(
                              'Open Source Foundation • Micro-SaaS Engine',
                              style: GoogleFonts.inter(
                                fontSize: 12,
                                fontWeight: FontWeight.w600,
                                color: AppTheme.textPrimary,
                              ),
                            ),
                          ],
                        ),
                      ),

                      const SizedBox(height: 24),

                      // Main Hero Heading
                      Text(
                        'Modular Monolith Engine\nfor Micro-SaaS Portfolios',
                        textAlign: TextAlign.center,
                        style: Theme.of(context).textTheme.displayLarge?.copyWith(
                              fontSize: isDesktop ? 48 : 32,
                            ),
                      ),

                      const SizedBox(height: 20),

                      // Subtitle
                      ConstrainedBox(
                        constraints: const BoxConstraints(maxWidth: 720),
                        child: Text(
                          'A high-performance Go backend with dynamic PostgreSQL schema routing, '
                          'paired with a Flutter Web PWA on Cloudflare Pages, Supabase OIDC authentication, '
                          'and unified Sentry error tracking.',
                          textAlign: TextAlign.center,
                          style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                                fontSize: isDesktop ? 18 : 15,
                              ),
                        ),
                      ),

                      const SizedBox(height: 36),

                      // CTA Actions
                      Wrap(
                        spacing: 16,
                        runSpacing: 12,
                        alignment: WrapAlignment.center,
                        children: [
                          ElevatedButton.icon(
                            onPressed: _checkApiHealth,
                            icon: _isCheckingHealth
                                ? const SizedBox(
                                    width: 16,
                                    height: 16,
                                    child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
                                  )
                                : const Icon(Icons.sync_rounded, size: 18),
                            label: Text(_isCheckingHealth ? 'Checking...' : 'Ping Go API Backend'),
                          ),
                          OutlinedButton.icon(
                            onPressed: () {
                              ScaffoldMessenger.of(context).showSnackBar(
                                const SnackBar(
                                  content: Text('Repository: github.com/danysz/toolgridlab-core'),
                                ),
                              );
                            },
                            icon: const Icon(Icons.code_rounded, size: 18),
                            label: const Text('View Architecture Spec'),
                          ),
                        ],
                      ),

                      // Live Health Status Banner
                      if (_healthData != null) ...[
                        const SizedBox(height: 24),
                        Container(
                          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
                          decoration: BoxDecoration(
                            color: _isBackendHealthy
                                ? AppTheme.accentGreen.withValues(alpha: 0.12)
                                : AppTheme.accentRed.withValues(alpha: 0.12),
                            borderRadius: BorderRadius.circular(10),
                            border: Border.all(
                              color: _isBackendHealthy
                                  ? AppTheme.accentGreen.withValues(alpha: 0.4)
                                  : AppTheme.accentRed.withValues(alpha: 0.4),
                            ),
                          ),
                          child: Text(
                            _isBackendHealthy
                                ? '✓ Connected to Go Core API: Uptime ${_healthData?['uptime'] ?? ''}'
                                : '⚠️ Go backend offline at ${widget.config.apiBaseUrl}. Run `cd backend && go run ./cmd/server` to launch.',
                            style: GoogleFonts.inter(
                              fontSize: 13,
                              fontWeight: FontWeight.w500,
                              color: _isBackendHealthy ? AppTheme.accentGreen : AppTheme.accentRed,
                            ),
                          ),
                        ),
                      ],

                      const SizedBox(height: 64),

                      // 4 Architecture Pillars Grid
                      Align(
                        alignment: Alignment.centerLeft,
                        child: Text(
                          'Core Architectural Pillars',
                          style: GoogleFonts.inter(
                            fontSize: 22,
                            fontWeight: FontWeight.w700,
                            color: AppTheme.textPrimary,
                            letterSpacing: -0.5,
                          ),
                        ),
                      ),
                      const SizedBox(height: 8),
                      Align(
                        alignment: Alignment.centerLeft,
                        child: Text(
                          'Designed for high isolation, zero platform lock-in, and independent spin-off capability.',
                          style: GoogleFonts.inter(
                            fontSize: 14,
                            color: AppTheme.textSecondary,
                          ),
                        ),
                      ),
                      const SizedBox(height: 24),

                      GridView.count(
                        crossAxisCount: isDesktop ? 2 : 1,
                        shrinkWrap: true,
                        physics: const NeverScrollableScrollPhysics(),
                        mainAxisSpacing: 20,
                        crossAxisSpacing: 20,
                        childAspectRatio: isDesktop ? 1.45 : 1.15,
                        children: const [
                          ArchitectureCard(
                            icon: Icons.alt_route_rounded,
                            title: 'Dynamic Schema Routing',
                            description:
                                'All tools live in a single Supabase PostgreSQL database while maintaining strict schema isolation via `SET search_path TO schema_<tool>;`. Any tool can be exported without lock-in.',
                            badge: 'PostgreSQL Isolation',
                            accentColor: AppTheme.primary,
                          ),
                          ArchitectureCard(
                            icon: Icons.layers_outlined,
                            title: 'Modular Monolith (Go)',
                            description:
                                'Unified Go binary on Google Cloud Run with strict `/api/v1/<tool>/` versioning. High performance, zero baseline idle cost, and a shared core connection pool.',
                            badge: 'Go 1.22+ on Cloud Run',
                            accentColor: AppTheme.primaryGradientEnd,
                          ),
                          ArchitectureCard(
                            icon: Icons.shield_outlined,
                            title: 'Decoupled Supabase OIDC',
                            description:
                                'Google & Apple OAuth handled centrally. Global pointer maps external Auth IDs to internal profiles (`schema_core.user_profiles`) with zero local password storage.',
                            badge: 'Central Global Auth',
                            accentColor: AppTheme.accentGreen,
                          ),
                          ArchitectureCard(
                            icon: Icons.cloud_done_outlined,
                            title: 'Cloudflare Pages & Sentry',
                            description:
                                'Flutter Web PWA deployed globally via Cloudflare Pages edge caching, with end-to-end unified Sentry error tracking across frontend and backend services.',
                            badge: 'Edge PWA & Observability',
                            accentColor: AppTheme.accentAmber,
                          ),
                        ],
                      ),

                      const SizedBox(height: 64),

                      // Live Interactive Core Utility: Email & DNS MX Verifier
                      Container(
                        padding: const EdgeInsets.all(32),
                        decoration: BoxDecoration(
                          color: AppTheme.cardSurface,
                          borderRadius: BorderRadius.circular(20),
                          border: Border.all(color: AppTheme.cardBorder, width: 1.5),
                        ),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Container(
                                  padding: const EdgeInsets.all(8),
                                  decoration: BoxDecoration(
                                    color: AppTheme.primary.withValues(alpha: 0.15),
                                    borderRadius: BorderRadius.circular(8),
                                  ),
                                  child: const Icon(Icons.verified_user_rounded, color: AppTheme.primary, size: 20),
                                ),
                                const SizedBox(width: 12),
                                Flexible(
                                  child: Text(
                                    'Interactive Core Utility: Strict Email Verifier',
                                    style: GoogleFonts.inter(
                                      fontSize: 18,
                                      fontWeight: FontWeight.w700,
                                      color: AppTheme.textPrimary,
                                    ),
                                  ),
                                ),
                              ],
                            ),
                            const SizedBox(height: 8),
                            Text(
                              'Test the Go backend verification engine in real time: executes (1) RFC 5322 regex validation, (2) disposable domain blacklist filtering, and (3) DNS `net.LookupMX` mail server checks.',
                              style: GoogleFonts.inter(fontSize: 14, color: AppTheme.textSecondary, height: 1.5),
                            ),
                            const SizedBox(height: 24),
                            Row(
                              children: [
                                Expanded(
                                  child: TextField(
                                    controller: _emailController,
                                    decoration: InputDecoration(
                                      hintText: 'Enter an email (e.g., founder@toolgridlab.com or test@tempmail.com)',
                                      hintStyle: GoogleFonts.inter(color: AppTheme.textMuted, fontSize: 14),
                                      filled: true,
                                      fillColor: AppTheme.surface,
                                      border: OutlineInputBorder(
                                        borderRadius: BorderRadius.circular(12),
                                        borderSide: const BorderSide(color: AppTheme.cardBorder),
                                      ),
                                      enabledBorder: OutlineInputBorder(
                                        borderRadius: BorderRadius.circular(12),
                                        borderSide: const BorderSide(color: AppTheme.cardBorder),
                                      ),
                                      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                                    ),
                                    onSubmitted: (_) => _verifyEmail(),
                                  ),
                                ),
                                const SizedBox(width: 12),
                                ElevatedButton(
                                  onPressed: _isVerifyingEmail ? null : _verifyEmail,
                                  child: _isVerifyingEmail
                                      ? const SizedBox(
                                          width: 16,
                                          height: 16,
                                          child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
                                        )
                                      : const Text('Verify Email'),
                                ),
                              ],
                            ),

                            if (_verificationResult != null) ...[
                              const SizedBox(height: 16),
                              Container(
                                padding: const EdgeInsets.all(16),
                                decoration: BoxDecoration(
                                  color: (_verificationResult?['valid'] == true)
                                      ? AppTheme.accentGreen.withValues(alpha: 0.12)
                                      : AppTheme.accentRed.withValues(alpha: 0.12),
                                  borderRadius: BorderRadius.circular(12),
                                  border: Border.all(
                                    color: (_verificationResult?['valid'] == true)
                                        ? AppTheme.accentGreen.withValues(alpha: 0.4)
                                        : AppTheme.accentRed.withValues(alpha: 0.4),
                                  ),
                                ),
                                child: Row(
                                  children: [
                                    Icon(
                                      (_verificationResult?['valid'] == true)
                                          ? Icons.check_circle_rounded
                                          : Icons.error_rounded,
                                      color: (_verificationResult?['valid'] == true)
                                          ? AppTheme.accentGreen
                                          : AppTheme.accentRed,
                                      size: 20,
                                    ),
                                    const SizedBox(width: 12),
                                    Expanded(
                                      child: Text(
                                        'Result: ${_verificationResult?['message'] ?? ''} (${_verificationResult?['email'] ?? ''})',
                                        style: GoogleFonts.inter(
                                          fontSize: 14,
                                          fontWeight: FontWeight.w600,
                                          color: (_verificationResult?['valid'] == true)
                                              ? AppTheme.accentGreen
                                              : AppTheme.accentRed,
                                        ),
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ],
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),

            const AppFooter(),
          ],
        ),
      ),
    );
  }
}
