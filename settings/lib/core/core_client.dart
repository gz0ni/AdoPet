import 'dart:convert';

import 'package:http/http.dart' as http;

/// Talks to pet-core's local HTTP API (127.0.0.1 only).
class CoreClient {
  CoreClient({required this.port, required this.token});

  final int port;
  final String token;

  Uri _u(String path) => Uri.parse('http://127.0.0.1:$port$path');
  Map<String, String> get _h => {'X-Pet-Token': token};

  Future<bool> health() async {
    try {
      final r = await http.get(_u('/health'), headers: _h).timeout(const Duration(seconds: 1));
      return r.statusCode == 200;
    } catch (_) {
      return false;
    }
  }

  Future<Map<String, dynamic>> state() async =>
      jsonDecode((await http.get(_u('/state'), headers: _h)).body) as Map<String, dynamic>;

  Future<void> preview(String animation) =>
      http.post(_u('/preview'), headers: _h, body: jsonEncode({'animation': animation}));

  Future<void> reload() => http.post(_u('/reload'), headers: _h);
}
