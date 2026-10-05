class LearningTask {
  const LearningTask({
    required this.id,
    required this.title,
    required this.description,
  });

  final String id;
  final String title;
  final String description;

  factory LearningTask.fromJson(Map<String, dynamic> json) {
    return LearningTask(
      id: json['id'] as String? ?? '',
      title: json['title'] as String? ?? '',
      description: json['description'] as String? ?? '',
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'title': title,
    'description': description,
  };
}

class LearningSchedule {
  const LearningSchedule({
    required this.id,
    required this.title,
    this.description,
    required this.tasks,
  });

  final String id;
  final String title;
  final String? description;
  final List<LearningTask> tasks;

  factory LearningSchedule.fromJson(Map<String, dynamic> json) {
    return LearningSchedule(
      id: json['id'] as String? ?? '',
      title: json['title'] as String? ?? '',
      description: json['description'] as String?,
      tasks: (json['tasks'] as List<dynamic>? ?? [])
          .map((e) => LearningTask.fromJson(e as Map<String, dynamic>))
          .toList(),
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'title': title,
    if (description != null && description!.isNotEmpty)
      'description': description,
    'tasks': tasks.map((task) => task.toJson()).toList(),
  };
}
