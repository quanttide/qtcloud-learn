import 'enums.dart';

/// 学员（对齐 `qtcloud-course` student.dart + `qtclass` session.dart Student）。
class Student {
  final String id;
  final String name;
  final String email;
  final String? avatar;
  final PlanType plan;

  const Student({
    required this.id,
    required this.name,
    required this.email,
    this.avatar,
    this.plan = PlanType.free,
  });

  factory Student.fromJson(Map<String, dynamic> json) {
    return Student(
      id: json['id'] as String,
      name: json['name'] as String,
      email: json['email'] as String? ?? '',
      avatar: json['avatar'] as String?,
      plan: PlanType.fromString(json['plan'] as String? ?? 'free'),
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'name': name,
        'email': email,
        if (avatar != null) 'avatar': avatar,
        'plan': plan.name,
      };
}
