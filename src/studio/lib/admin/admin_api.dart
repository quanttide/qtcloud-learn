// 后台数据 API：学员档案/立项（qtcloud-learn，对齐原型契约 /api/learners + /api/proposals）。

import 'dart:convert';

import 'package:http/http.dart' as http;

import '../models/application.dart';
import '../models/learner.dart';

/// 默认后台 API 地址（--dart-define=QTCLOUD_LEARN_API_URL=... 注入生产网关）。
String defaultAdminBaseUrl() {
  const String fromEnv = String.fromEnvironment('QTCLOUD_LEARN_API_URL');
  if (fromEnv.isNotEmpty) {
    return fromEnv;
  }
  return 'http://localhost:8080';
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

  /// 学员档案（后台学员表）。
  Future<List<Learner>> fetchLearners() async {
    final body = await _get('/api/learners');
    return (body['learners'] as List<dynamic>? ?? [])
        .map((e) => Learner.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  /// 立项列表（不含已删除）。
  Future<List<Application>> fetchProposals() async {
    final body = await _get('/api/proposals');
    return (body['proposals'] as List<dynamic>? ?? [])
        .map((e) => Application.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  /// 软删除历史。
  Future<List<Application>> fetchHistory() async {
    final body = await _get('/api/proposals/history');
    return (body['history'] as List<dynamic>? ?? [])
        .map((e) => Application.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  /// 软删除立项。
  Future<void> deleteProposal(String id) async {
    final resp = await _client
        .delete(Uri.parse('$baseUrl/api/proposals/$id'))
        .timeout(const Duration(seconds: 15));
    if (resp.statusCode != 200) {
      throw AdminApiException('删除失败（HTTP ${resp.statusCode}）');
    }
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
