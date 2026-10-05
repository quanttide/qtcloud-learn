// 后台数据 API：学员档案/立项（qtcloud-learn，对齐原型契约 /api/learners + /api/proposals）。

import 'dart:convert';

import 'package:http/http.dart' as http;

import '../models/application.dart';
import '../models/learner.dart';
import '../models/learning.dart';

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

  Future<List<LearningTask>> fetchTasks() async {
    final body = await _getList('/tasks');
    return body
        .map((e) => LearningTask.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<LearningTask> saveTask(LearningTask task) async {
    final path = '/tasks/${task.id}';
    final exists = await _exists(path);
    final body = exists
        ? await _put(path, task.toJson())
        : await _post('/tasks', task.toJson());
    return LearningTask.fromJson(body);
  }

  Future<void> deleteTask(String id) => _delete('/tasks/$id');

  Future<List<LearningSchedule>> fetchSchedules() async {
    final body = await _getList('/schedules');
    return body
        .map((e) => LearningSchedule.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<LearningSchedule> saveSchedule(LearningSchedule schedule) async {
    final path = '/schedules/${schedule.id}';
    final exists = await _exists(path);
    final body = exists
        ? await _put(path, schedule.toJson())
        : await _post('/schedules', schedule.toJson());
    return LearningSchedule.fromJson(body);
  }

  Future<void> deleteSchedule(String id) => _delete('/schedules/$id');

  /// 软删除立项。
  Future<void> deleteProposal(String id) async {
    final resp = await _client
        .delete(Uri.parse('$baseUrl/api/proposals/$id'))
        .timeout(const Duration(seconds: 15));
    if (resp.statusCode != 200) {
      throw AdminApiException('删除失败（HTTP ${resp.statusCode}）');
    }
  }

  Future<bool> _exists(String path) async {
    final resp = await _client
        .get(Uri.parse('$baseUrl$path'))
        .timeout(const Duration(seconds: 15));
    if (resp.statusCode == 200) return true;
    if (resp.statusCode == 404) return false;
    throw AdminApiException('HTTP ${resp.statusCode}');
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

  Future<List<dynamic>> _getList(String path) async {
    final resp = await _client
        .get(Uri.parse('$baseUrl$path'))
        .timeout(const Duration(seconds: 15));
    if (resp.statusCode != 200) {
      throw AdminApiException('HTTP ${resp.statusCode}');
    }
    return jsonDecode(utf8.decode(resp.bodyBytes)) as List<dynamic>;
  }

  Future<Map<String, dynamic>> _post(
    String path,
    Map<String, dynamic> body,
  ) async {
    final resp = await _client
        .post(
          Uri.parse('$baseUrl$path'),
          headers: {'content-type': 'application/json'},
          body: jsonEncode(body),
        )
        .timeout(const Duration(seconds: 15));
    if (resp.statusCode != 201 && resp.statusCode != 200) {
      throw AdminApiException('HTTP ${resp.statusCode}');
    }
    return jsonDecode(utf8.decode(resp.bodyBytes)) as Map<String, dynamic>;
  }

  Future<Map<String, dynamic>> _put(
    String path,
    Map<String, dynamic> body,
  ) async {
    final resp = await _client
        .put(
          Uri.parse('$baseUrl$path'),
          headers: {'content-type': 'application/json'},
          body: jsonEncode(body),
        )
        .timeout(const Duration(seconds: 15));
    if (resp.statusCode != 200) {
      throw AdminApiException('HTTP ${resp.statusCode}');
    }
    return jsonDecode(utf8.decode(resp.bodyBytes)) as Map<String, dynamic>;
  }

  Future<void> _delete(String path) async {
    final resp = await _client
        .delete(Uri.parse('$baseUrl$path'))
        .timeout(const Duration(seconds: 15));
    if (resp.statusCode != 204 && resp.statusCode != 200) {
      throw AdminApiException('HTTP ${resp.statusCode}');
    }
  }
}
