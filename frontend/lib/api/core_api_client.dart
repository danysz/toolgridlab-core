import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:sentry_flutter/sentry_flutter.dart';

class CoreApiClient {
  final String baseUrl;
  final String? Function() getAccessToken;

  CoreApiClient({
    required this.baseUrl,
    required this.getAccessToken,
  });

  Map<String, String> _buildHeaders({String? schema}) {
    final headers = <String, String>{
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    };

    final token = getAccessToken();
    if (token != null && token.isNotEmpty) {
      headers['Authorization'] = 'Bearer $token';
    }

    if (schema != null && schema.isNotEmpty) {
      headers['X-ToolGrid-Schema'] = schema;
    }

    return headers;
  }

  // Health check endpoint
  Future<Map<String, dynamic>> checkHealth() async {
    final uri = Uri.parse('$baseUrl/api/v1/healthz');
    try {
      final response = await http.get(uri, headers: _buildHeaders());
      if (response.statusCode == 200) {
        return jsonDecode(response.body) as Map<String, dynamic>;
      }
      return {'status': 'error', 'statusCode': response.statusCode};
    } catch (e, stackTrace) {
      await Sentry.captureException(e, stackTrace: stackTrace);
      return {'status': 'offline', 'error': e.toString()};
    }
  }

  // Generic Email verification endpoint
  Future<Map<String, dynamic>> verifyEmail(String email) async {
    final uri = Uri.parse('$baseUrl/api/v1/core/verify-email');
    try {
      final response = await http.post(
        uri,
        headers: _buildHeaders(),
        body: jsonEncode({'email': email}),
      );
      return jsonDecode(response.body) as Map<String, dynamic>;
    } catch (e, stackTrace) {
      await Sentry.captureException(e, stackTrace: stackTrace);
      return {'valid': false, 'message': 'Network error: ${e.toString()}'};
    }
  }

  // Fetch authenticated user profile
  Future<Map<String, dynamic>?> fetchUserProfile() async {
    final uri = Uri.parse('$baseUrl/api/v1/core/me');
    try {
      final response = await http.get(uri, headers: _buildHeaders());
      if (response.statusCode == 200) {
        return jsonDecode(response.body) as Map<String, dynamic>;
      }
      return null;
    } catch (e, stackTrace) {
      await Sentry.captureException(e, stackTrace: stackTrace);
      return null;
    }
  }
}
