// 后台数据 API：学员/进度/立项（qtcloud-learn /api/v1，LMS 后台只读查看）。

import 'dart:convert';

import 'package:http/http.dart' as http;

import '../models/application.dart';
import '../models/progress.dart';
import '../models/student.dart';

/// 默认后台 API 地址（--dart-define=QTCLOUD_LEARN_API_URL=... 注入生产网关）。
String defaultAdminBaseUrl() {
  const String fromEnv = String.fromEnvironment('QTCLOUD_LEARN_API_URL');
  if (fromEnv.isNotEmpty) {
    return fromEnv;
  }
  return 'http://localhost:8080/api/v1';
}

class AdminApiException implements Exception {
  const AdminApiException(this.message);
  final String message;
  @override
  String toString() => message;
}

class AdminApi {
  AdminApi({http.Client? client, String? baseUrl})
    : _client = client ?? http.Client(),
      baseUrl = baseUrl ?? defaultAdminBaseUrl();

  final http.Client _client;
  final String baseUrl;

  Future<List<Student>> fetchStudents() async {
    final body = await _get('/students');
    return (body['students'] as List<dynamic>? ?? [])
        .map((e) => Student.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<List<Progress>> fetchProgress() async {
    final body = await _get('/progress');
    return (body['progress'] as List<dynamic>? ?? [])
        .map((e) => Progress.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<List<Application>> fetchApplications() async {
    final body = await _get('/applications');
    return (body['applications'] as List<dynamic>? ?? [])
        .map((e) => Application.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<Map<String, dynamic>> _get(String path) async {
    final resp = await _client
        .get(Uri.parse('$baseUrl$path'))
        .timeout(const Duration(seconds: 15));
    if (resp.statusCode != 200) {
      throw AdminApiException('HTTP ${resp.statusCode}');
    }
    return jsonDecode(utf8.decode(resp.bodyBytes)) as Map<String, dynamic>;
  }
}
