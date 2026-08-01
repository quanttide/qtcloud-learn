import 'enums.dart';

/// 班级（统一 `qtcloud-course` Class / ClassTeaching）。
class ClassModel {
  final String id;
  final String name;
  final String slug;
  final String refName;
  final String refType;
  final String refId;
  final ClassStatus status;
  final String startDate;
  final String endDate;
  final int studentCount;
  final double progress;
  final List<String> teacherIds;
  final List<String> studentIds;

  const ClassModel({
    required this.id,
    required this.name,
    this.slug = '',
    required this.refName,
    this.refType = 'program',
    required this.refId,
    this.status = ClassStatus.preparing,
    required this.startDate,
    required this.endDate,
    this.studentCount = 0,
    this.progress = 0.0,
    this.teacherIds = const [],
    this.studentIds = const [],
  });

  factory ClassModel.fromJson(Map<String, dynamic> json) {
    return ClassModel(
      id: json['id'] as String,
      name: json['name'] as String,
      slug: json['slug'] as String? ?? '',
      refName: json['refName'] as String? ?? '',
      refType: json['refType'] as String? ?? 'program',
      refId: json['refId'] as String? ?? '',
      status: ClassStatus.fromString(json['status'] as String? ?? 'preparing'),
      startDate: json['startDate'] as String? ?? '',
      endDate: json['endDate'] as String? ?? '',
      studentCount: (json['studentCount'] as num?)?.toInt() ?? 0,
      progress: (json['progress'] as num?)?.toDouble() ?? 0.0,
      teacherIds: (json['teacherIds'] as List<dynamic>?)
              ?.map((e) => e as String)
              .toList() ??
          [],
      studentIds: (json['studentIds'] as List<dynamic>?)
              ?.map((e) => e as String)
              .toList() ??
          [],
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'name': name,
        'refName': refName,
        'refType': refType,
        'refId': refId,
        'status': status.name,
        'startDate': startDate,
        'endDate': endDate,
        'studentCount': studentCount,
        'progress': progress,
        'teacherIds': teacherIds,
        'studentIds': studentIds,
      };
}
