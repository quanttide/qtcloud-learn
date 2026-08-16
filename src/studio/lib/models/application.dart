/// 立项申请（生产实习第五步·微型创业 Sell Your Demo，对齐原型 qt-proposals）。
/// 5 问：机会/匹配/假设/Demo + 方向类型 + 组队方式姓名栏。
/// teamMode: personal（个人独立，TeamLeader=本人）/ partner（搭档，队长+队员）。
class Application {
  final String id;
  final String projectName;
  final String opportunity;
  final String fit;
  final String hypothesis;
  final String demo;
  final String directionType;
  final String teamMode;
  final String teamLeader;
  final String teamMember;
  final String studentName;
  final String status;
  final String? submittedAt;
  final String? deletedAt;

  const Application({
    required this.id,
    required this.projectName,
    required this.opportunity,
    required this.fit,
    required this.hypothesis,
    required this.demo,
    required this.directionType,
    required this.teamMode,
    required this.teamLeader,
    required this.teamMember,
    required this.studentName,
    required this.status,
    this.submittedAt,
    this.deletedAt,
  });

  factory Application.fromJson(Map<String, dynamic> json) => Application(
    id: json['id'] as String? ?? '',
    projectName: json['projectName'] as String? ?? '',
    opportunity: json['opportunity'] as String? ?? '',
    fit: json['fit'] as String? ?? '',
    hypothesis: json['hypothesis'] as String? ?? '',
    demo: json['demo'] as String? ?? '',
    directionType: json['directionType'] as String? ?? '',
    teamMode: json['teamMode'] as String? ?? '',
    teamLeader: json['teamLeader'] as String? ?? '',
    teamMember: json['teamMember'] as String? ?? '',
    studentName: json['studentName'] as String? ?? '',
    status: json['status'] as String? ?? '',
    submittedAt: json['submittedAt'] as String? ?? null,
    deletedAt: json['deletedAt'] as String? ?? null,
  );

  /// 组队展示：个人独立 → 个人姓名；搭档 → 队长+队员。
  String get teamLabel {
    if (teamMode == 'partner' && teamMember.isNotEmpty) {
      return '队长：$teamLeader / 队员：$teamMember';
    }
    return teamLeader;
  }
}
