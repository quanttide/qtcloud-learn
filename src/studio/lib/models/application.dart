/// 立项申请（生产实习第五步·微型创业）。
/// teamMode: personal（个人独立，memberNames=[个人姓名]）/ partner（搭档，[队长, 队员]）。
class Application {
  final String id;
  final String projectName;
  final String blindSpot;
  final String demoPlan;
  final String direction;
  final String teamMode;
  final List<String> memberNames;
  final String studentId;
  final String studentName;
  final String status;
  final String? createdAt;

  const Application({
    required this.id,
    required this.projectName,
    required this.blindSpot,
    required this.demoPlan,
    required this.direction,
    required this.teamMode,
    required this.memberNames,
    required this.studentId,
    required this.studentName,
    required this.status,
    this.createdAt,
  });

  factory Application.fromJson(Map<String, dynamic> json) => Application(
    id: json['id'] as String? ?? '',
    projectName: json['projectName'] as String? ?? '',
    blindSpot: json['blindSpot'] as String? ?? '',
    demoPlan: json['demoPlan'] as String? ?? '',
    direction: json['direction'] as String? ?? '',
    teamMode: json['teamMode'] as String? ?? '',
    memberNames:
        (json['memberNames'] as List<dynamic>? ?? []).map((e) => e as String).toList(),
    studentId: json['studentId'] as String? ?? '',
    studentName: json['studentName'] as String? ?? '',
    status: json['status'] as String? ?? '',
    createdAt: json['createdAt'] as String? ?? null,
  );

  /// 姓名栏展示（个人独立 → 个人姓名；搭档 → 队长+队员）。
  String get memberLabel {
    final names = memberNames.where((n) => n.isNotEmpty).toList();
    if (names.isEmpty) return '-';
    if (teamMode == 'partner' && names.length > 1) {
      return '队长：${names.first} / 队员：${names.sublist(1).join('、')}';
    }
    return names.join('、');
  }
}
