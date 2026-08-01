/// 选课 / 报名。
class Enrollment {
  final String id;
  final String classId;
  final String studentId;
  final String status; // "enrolled" / "withdrawn"
  final String? enrolledAt;

  const Enrollment({
    required this.id,
    required this.classId,
    required this.studentId,
    this.status = 'enrolled',
    this.enrolledAt,
  });

  bool get isEnrolled => status == 'enrolled';

  factory Enrollment.fromJson(Map<String, dynamic> json) {
    return Enrollment(
      id: json['id'] as String,
      classId: json['classId'] as String,
      studentId: json['studentId'] as String,
      status: json['status'] as String? ?? 'enrolled',
      enrolledAt: json['enrolledAt'] as String?,
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'classId': classId,
        'studentId': studentId,
        'status': status,
        if (enrolledAt != null) 'enrolledAt': enrolledAt,
      };
}
