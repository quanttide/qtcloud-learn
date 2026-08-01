/// 学习进度（统一 `qtcloud-course` Class.progress 与 `qtclass` learning_record；
/// 学习记录转为服务端进度数据）。
class Progress {
  final String id;
  final String studentId;
  final String classId;
  final double percent;
  final bool finished;
  final String? updatedAt;

  const Progress({
    required this.id,
    required this.studentId,
    required this.classId,
    this.percent = 0.0,
    this.finished = false,
    this.updatedAt,
  });

  factory Progress.fromJson(Map<String, dynamic> json) {
    return Progress(
      id: json['id'] as String,
      studentId: json['studentId'] as String,
      classId: json['classId'] as String,
      percent: (json['percent'] as num?)?.toDouble() ?? 0.0,
      finished: json['finished'] as bool? ?? false,
      updatedAt: json['updatedAt'] as String?,
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'studentId': studentId,
        'classId': classId,
        'percent': percent,
        'finished': finished,
        if (updatedAt != null) 'updatedAt': updatedAt,
      };
}
