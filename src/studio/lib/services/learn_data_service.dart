import 'dart:async';
import 'dart:convert';

import 'package:flutter/foundation.dart';
import 'package:flutter/services.dart';
import 'package:http/http.dart' as http;
import 'package:uuid/uuid.dart';

import '../models/class_model.dart';
import '../models/enrollment.dart';
import '../models/enums.dart';
import '../models/progress.dart';
import '../models/session.dart';
import '../models/student.dart';
import '../models/teacher.dart';

/// 统一数据服务：班级 / 学员 / 讲师 / 课次 / 选课 / 进度。
///
/// 自 `qtcloud-course` data_service.dart 移植，收拢 `qtclass` app_state.dart 的
/// sessions 加载，统一对接 `qtcloud-learn/provider`（`/api/v1`）。
class LearnDataService extends ChangeNotifier {
  List<ClassModel> _classes = [];
  List<Student> _students = [];
  List<Teacher> _teachers = [];
  List<Session> _sessions = [];
  List<Enrollment> _enrollments = [];
  List<Progress> _progress = [];
  bool _loaded = false;
  String? _error;
  bool _loading = false;
  bool _offlineFallback = false;

  final String? baseUrl;
  http.Client client;

  List<ClassModel> get classes => _classes;
  List<Student> get students => _students;
  List<Teacher> get teachers => _teachers;
  List<Session> get sessions => _sessions;
  List<Enrollment> get enrollments => _enrollments;
  List<Progress> get progress => _progress;
  bool get loaded => _loaded;
  String? get error => _error;
  bool get loading => _loading;
  bool get offlineFallback => _offlineFallback;
  bool get isApiMode => baseUrl != null;

  int get activeClasses =>
      _classes.where((c) => c.status.name == 'active').length;
  int get totalStudents => _classes.fold(0, (sum, c) => sum + c.studentCount);

  List<Session> get upcomingSessions =>
      _sessions.where((s) => s.status == SessionStatus.upcoming).toList();
  List<Session> get inProgressSessions =>
      _sessions.where((s) => s.status == SessionStatus.inProgress).toList();
  List<Session> get completedSessions =>
      _sessions.where((s) => s.status == SessionStatus.completed).toList();

  LearnDataService({this.baseUrl, http.Client? client})
      : client = client ?? http.Client();

  void markLoaded() {
    _loaded = true;
    notifyListeners();
  }

  Future<void> load() async {
    _loading = true;
    _error = null;
    _offlineFallback = false;
    notifyListeners();
    try {
      if (baseUrl != null) {
        try {
          await _loadFromApi();
        } catch (e) {
          debugPrint('API load failed ($e), falling back to local JSON');
          _offlineFallback = true;
          try {
            await _loadFromAssets();
          } catch (e2) {
            debugPrint('Fallback assets also failed: $e2');
          }
        }
      } else {
        await _loadFromAssets();
      }
      _loaded = true;
    } catch (e) {
      _error = e.toString();
    }
    _loading = false;
    notifyListeners();
  }

  Future<void> _loadFromAssets() async {
    _classes = await _decodeList('assets/classes.json', ClassModel.fromJson);
    _students = await _decodeList('assets/students.json', Student.fromJson);
    _teachers = await _decodeList('assets/teachers.json', Teacher.fromJson);
    _sessions = await _decodeList('assets/sessions.json', Session.fromJson);
    _enrollments = await _decodeList('assets/enrollments.json', Enrollment.fromJson);
    _progress = await _decodeList('assets/progress.json', Progress.fromJson);
  }

  Future<List<T>> _decodeList<T>(
      String asset, T Function(Map<String, dynamic>) fromJson) async {
    try {
      final str = await rootBundle.loadString(asset);
      return (json.decode(str) as List<dynamic>)
          .map((e) => fromJson(e as Map<String, dynamic>))
          .toList();
    } catch (e) {
      debugPrint('Failed to load $asset: $e');
      return [];
    }
  }

  Future<void> _loadFromApi() async {
    final api = _Api(baseUrl!, client);
    _classes = await api.getList('classes', ClassModel.fromJson);
    _students = await api.getList('students', Student.fromJson);
    _teachers = await api.getList('teachers', Teacher.fromJson);
    _sessions = await api.getList('sessions', Session.fromJson);
    _enrollments = await api.getList('enrollments', Enrollment.fromJson);
    _progress = await api.getList('progress', Progress.fromJson);
  }

  // ── API Sync Helpers ──

  void _apiPost(String path, Map<String, dynamic> body) {
    if (baseUrl == null) return;
    unawaited(
      client
          .post(
            Uri.parse('$baseUrl/api/v1$path'),
            headers: {'Content-Type': 'application/json'},
            body: json.encode(body),
          )
          .then((_) {})
          .catchError((_) {}),
    );
  }

  void _apiPut(String path, Map<String, dynamic> body) {
    if (baseUrl == null) return;
    unawaited(
      client
          .put(
            Uri.parse('$baseUrl/api/v1$path'),
            headers: {'Content-Type': 'application/json'},
            body: json.encode(body),
          )
          .then((_) {})
          .catchError((_) {}),
    );
  }

  void _apiDelete(String path) {
    if (baseUrl == null) return;
    unawaited(
      client
          .delete(Uri.parse('$baseUrl/api/v1$path'))
          .then((_) {})
          .catchError((_) {}),
    );
  }

  static const _uuid = Uuid();
  String _nextId() => _uuid.v4();

  // ---- 班级 CRUD ----

  void createClass({
    required String name,
    required String refName,
    required String refId,
    String refType = 'program',
    required String startDate,
    required String endDate,
  }) {
    final newClass = ClassModel(
      id: _nextId(),
      name: name,
      refName: refName,
      refType: refType,
      refId: refId,
      startDate: startDate,
      endDate: endDate,
    );
    _classes = [..._classes, newClass];
    _apiPost('/classes', newClass.toJson());
    notifyListeners();
  }

  void deleteClass(String id) {
    _classes = _classes.where((c) => c.id != id).toList();
    _apiDelete('/classes/$id');
    notifyListeners();
  }

  // ---- 选课 / 报名 ----

  bool isEnrolled(String classId, String studentId) =>
      _enrollments.any(
          (e) => e.classId == classId && e.studentId == studentId && e.isEnrolled);

  List<ClassModel> getMyClasses(String studentId) {
    final ids = _enrollments
        .where((e) => e.studentId == studentId && e.isEnrolled)
        .map((e) => e.classId)
        .toSet();
    return _classes.where((c) => ids.contains(c.id)).toList();
  }

  void enroll(String classId, String studentId) {
    if (isEnrolled(classId, studentId)) return;
    final e = Enrollment(
      id: _nextId(),
      classId: classId,
      studentId: studentId,
      enrolledAt: DateTime.now().toIso8601String(),
    );
    _enrollments = [..._enrollments, e];
    _apiPost('/enrollments', e.toJson());
    notifyListeners();
  }

  void withdraw(String classId, String studentId) {
    final index = _enrollments.indexWhere(
        (e) => e.classId == classId && e.studentId == studentId && e.isEnrolled);
    if (index == -1) return;
    final e = _enrollments[index];
    _enrollments[index] = Enrollment(
      id: e.id,
      classId: e.classId,
      studentId: e.studentId,
      status: 'withdrawn',
      enrolledAt: e.enrolledAt,
    );
    _enrollments = [..._enrollments];
    _apiPut('/enrollments/${e.id}', _enrollments[index].toJson());
    notifyListeners();
  }

  // ---- 学习进度 ----

  Progress? getProgress(String classId, String studentId) {
    final list = _progress
        .where((p) => p.classId == classId && p.studentId == studentId)
        .toList();
    return list.isEmpty ? null : list.first;
  }

  void reportProgress(String classId, String studentId,
      {double percent = 0.0, bool finished = false}) {
    final existing = getProgress(classId, studentId);
    if (existing != null) {
      _progress = [
        for (final p in _progress)
          if (p.id == existing.id)
            Progress(
              id: p.id,
              studentId: studentId,
              classId: classId,
              percent: percent,
              finished: finished,
              updatedAt: DateTime.now().toIso8601String(),
            )
          else
            p,
      ];
      _apiPut('/progress/${existing.id}',
          _progress.firstWhere((p) => p.id == existing.id).toJson());
    } else {
      final p = Progress(
        id: _nextId(),
        studentId: studentId,
        classId: classId,
        percent: percent,
        finished: finished,
        updatedAt: DateTime.now().toIso8601String(),
      );
      _progress = [..._progress, p];
      _apiPost('/progress', p.toJson());
    }
    notifyListeners();
  }

  // ---- 课次 / 考勤 ----

  void markAttendance(String sessionId, String studentId, AttendanceStatus status) {
    final index = _sessions.indexWhere((s) => s.id == sessionId);
    if (index == -1) return;
    final session = _sessions[index];
    final attendances = [
      for (final a in session.attendances)
        if (a.studentId == studentId) Attendance(studentId: studentId, status: status)
        else a,
    ];
    if (!attendances.any((a) => a.studentId == studentId)) {
      attendances.add(Attendance(studentId: studentId, status: status));
    }
    _updateSession(session, attendances: attendances);
  }

  void markAllAttendance(String sessionId, AttendanceStatus status) {
    final index = _sessions.indexWhere((s) => s.id == sessionId);
    if (index == -1) return;
    final session = _sessions[index];
    final attendances = session.attendances.isEmpty
        ? _students.map((s) => Attendance(studentId: s.id, status: status)).toList()
        : [
            for (final a in session.attendances)
              Attendance(studentId: a.studentId, status: status)
          ];
    _updateSession(session, attendances: attendances);
  }

  void _updateSession(Session session, {List<Attendance>? attendances}) {
    final updated = Session(
      id: session.id,
      classId: session.classId,
      lessonTitle: session.lessonTitle,
      teacherId: session.teacherId,
      startTime: session.startTime,
      durationMinutes: session.durationMinutes,
      location: session.location,
      status: session.status,
      attendances: attendances ?? session.attendances,
    );
    _sessions = [
      for (final s in _sessions)
        if (s.id == session.id) updated else s,
    ];
    _apiPut('/sessions/${session.id}', updated.toJson());
    notifyListeners();
  }

  // ---- 关联查询 ----

  List<Student> getStudentsByClass(String classId) {
    final idx = _classes.indexWhere((c) => c.id == classId);
    if (idx == -1) return [];
    return _students
        .where((s) => _classes[idx].studentIds.contains(s.id))
        .toList();
  }

  List<Teacher> getTeachersByClass(String classId) {
    final idx = _classes.indexWhere((c) => c.id == classId);
    if (idx == -1) return [];
    return _teachers
        .where((t) => _classes[idx].teacherIds.contains(t.id))
        .toList();
  }

  String? classNameOf(String classId) {
    for (final c in _classes) {
      if (c.id == classId) return c.name;
    }
    return null;
  }

  Teacher? teacherById(String teacherId) {
    for (final t in _teachers) {
      if (t.id == teacherId) return t;
    }
    return null;
  }
}

/// /api/v1 前缀的 GET 列表封装。
class _Api {
  final String baseUrl;
  final http.Client client;

  _Api(this.baseUrl, this.client);

  Future<List<T>> getList<T>(
      String resource, T Function(Map<String, dynamic>) fromJson) async {
    final uri = Uri.parse('$baseUrl/api/v1/$resource');
    final resp = await client.get(uri);
    if (resp.statusCode != 200) {
      throw Exception('Failed to load $resource: ${resp.statusCode}');
    }
    return (json.decode(resp.body) as List<dynamic>)
        .map((e) => fromJson(e as Map<String, dynamic>))
        .toList();
  }
}
