/// 学员档案（对齐原型 qt-students）：上报进度/提交立项自动建档。
class Learner {
  final String id;
  final String name;
  final String course;
  final int progressMax;
  final int progressTotal;
  final String? activeAt;
  final String status;
  final String? projectName;

  const Learner({
    required this.id,
    required this.name,
    required this.course,
    required this.progressMax,
    required this.progressTotal,
    this.activeAt,
    required this.status,
    this.projectName,
  });

  factory Learner.fromJson(Map<String, dynamic> json) => Learner(
    id: json['id'] as String? ?? '',
    name: json['name'] as String? ?? '',
    course: json['course'] as String? ?? '',
    progressMax: json['progressMax'] as int? ?? 0,
    progressTotal: json['progressTotal'] as int? ?? 0,
    activeAt: json['activeAt'] as String? ?? null,
    status: json['status'] as String? ?? '',
    projectName: json['projectName'] as String? ?? null,
  );

  /// 进度文案（X/5 模块）。
  String get progressLabel => '$progressMax/$progressTotal';
}
