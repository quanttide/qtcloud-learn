/// 讲师（对齐 `qtcloud-course` teacher.dart + `qtclass` session.dart Teacher）。
class Teacher {
  final String id;
  final String name;
  final String email;
  final String? title;

  const Teacher({
    required this.id,
    required this.name,
    this.email = '',
    this.title,
  });

  factory Teacher.fromJson(Map<String, dynamic> json) {
    return Teacher(
      id: json['id'] as String,
      name: json['name'] as String,
      email: json['email'] as String? ?? '',
      title: json['title'] as String?,
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'name': name,
        'email': email,
        if (title != null) 'title': title,
      };
}
