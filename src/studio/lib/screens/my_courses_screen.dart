import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../models/enums.dart';
import '../models/student.dart';
import '../services/learn_data_service.dart';
import '../widgets/status_chip.dart';
import 'progress_screen.dart';

/// 我的课程 / 选课报名页。
class MyCoursesScreen extends StatelessWidget {
  final Student student;

  const MyCoursesScreen({super.key, required this.student});

  @override
  Widget build(BuildContext context) {
    final service = context.watch<LearnDataService>();
    final myClasses = service.getMyClasses(student.id);
    final availableClasses = service.classes
        .where((c) => !service.isEnrolled(c.id, student.id))
        .toList();

    return Scaffold(
      appBar: AppBar(title: const Text('我的课程')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          if (myClasses.isNotEmpty) ...[
            Text('已选课程 (${myClasses.length})',
                style: Theme.of(context)
                    .textTheme
                    .titleMedium
                    ?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            ...myClasses.map((c) => Card(
                  margin: const EdgeInsets.only(bottom: 12),
                  child: ListTile(
                    title: Text(c.name,
                        style: const TextStyle(fontWeight: FontWeight.bold)),
                    subtitle: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text('${c.refName} · ${c.startDate} - ${c.endDate}'),
                        const SizedBox(height: 4),
                        LinearProgressIndicator(
                          value: c.progress,
                          backgroundColor: Colors.grey[200],
                        ),
                        const SizedBox(height: 2),
                        Text('${(c.progress * 100).toInt()}%',
                            style: TextStyle(
                                color: Colors.grey[600], fontSize: 12)),
                      ],
                    ),
                    trailing: const Icon(Icons.chevron_right),
                    onTap: () => Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (_) =>
                            ProgressScreen(classModel: c, student: student),
                      ),
                    ),
                  ),
                )),
            const SizedBox(height: 16),
          ],
          Text('选课报名',
              style: Theme.of(context)
                  .textTheme
                  .titleMedium
                  ?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          if (availableClasses.isEmpty)
            Card(
              child: Padding(
                padding: const EdgeInsets.all(24),
                child: Center(
                    child:
                        Text('暂无可以报名的班级', style: TextStyle(color: Colors.grey[500]))),
              ),
            )
          else
            ...availableClasses.map((c) => Card(
                  margin: const EdgeInsets.only(bottom: 12),
                  child: ListTile(
                    leading: Icon(
                      switch (c.status) {
                        ClassStatus.active => Icons.play_circle,
                        ClassStatus.preparing => Icons.schedule,
                        ClassStatus.ended => Icons.check_circle,
                      },
                      color: Colors.teal,
                    ),
                    title: Text(c.name,
                        style: const TextStyle(fontWeight: FontWeight.bold)),
                    subtitle: Text(
                        '${c.refName}\n${c.startDate} - ${c.endDate} · ${c.studentCount}人'),
                    isThreeLine: true,
                    trailing: StatusChip(status: c.status),
                    onTap: () => service.enroll(c.id, student.id),
                  ),
                )),
        ],
      ),
    );
  }
}
