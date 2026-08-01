import 'enums.dart';

/// 课次（统一 `qtclass` session.dart；className / courseName 由 classId 关联 Class 派生）。
class Session {
  final String id;
  final String classId;
  final String lessonTitle;
  final String teacherId;
  final DateTime startTime;
  final int durationMinutes;
  final String location;
  final SessionStatus status;
  final List<Attendance> attendances;

  const Session({
    required this.id,
    required this.classId,
    required this.lessonTitle,
    required this.teacherId,
    required this.startTime,
    required this.durationMinutes,
    required this.location,
    required this.status,
    this.attendances = const [],
  });

  DateTime get endTime => startTime.add(Duration(minutes: durationMinutes));

  factory Session.fromJson(Map<String, dynamic> json) {
    return Session(
      id: json['id'] as String,
      classId: json['classId'] as String,
      lessonTitle: json['lessonTitle'] as String,
      teacherId: json['teacherId'] as String? ?? '',
      startTime: DateTime.parse(json['startTime'] as String),
      durationMinutes: (json['durationMinutes'] as num?)?.toInt() ?? 0,
      location: json['location'] as String? ?? '',
      status: SessionStatus.fromString(json['status'] as String? ?? 'upcoming'),
      attendances: (json['attendances'] as List<dynamic>?)
              ?.map((e) => Attendance.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'classId': classId,
        'lessonTitle': lessonTitle,
        'teacherId': teacherId,
        'startTime': startTime.toIso8601String(),
        'durationMinutes': durationMinutes,
        'location': location,
        'status': status.name,
        'attendances': attendances.map((a) => a.toJson()).toList(),
      };
}

/// 考勤记录（`qtclass` AttendanceStatus）。
class Attendance {
  final String studentId;
  final AttendanceStatus status;

  const Attendance({required this.studentId, required this.status});

  factory Attendance.fromJson(Map<String, dynamic> json) {
    return Attendance(
      studentId: json['studentId'] as String,
      status: AttendanceStatus.fromString(json['status'] as String? ?? 'unknown'),
    );
  }

  Map<String, dynamic> toJson() => {
        'studentId': studentId,
        'status': status.name,
      };
}
