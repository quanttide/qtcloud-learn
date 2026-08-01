import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../models/student.dart';
import '../services/assessment_service.dart';
import '../services/learn_data_service.dart';
import '../widgets/cards.dart';
import 'my_courses_screen.dart';
import 'schedule_screen.dart';

/// 仪表盘 —— 学员侧指标（自 `qtcloud-course` dashboard_screen.dart 迁移，
/// 去除课程制作侧 ProgramService 依赖，改为学员视角）。
class DashboardScreen extends StatelessWidget {
  final Student student;
  final VoidCallback onLogout;

  const DashboardScreen({super.key, required this.student, required this.onLogout});

  @override
  Widget build(BuildContext context) {
    final dataService = context.watch<LearnDataService>();
    final assessmentService = context.watch<AssessmentService>();

    final myClasses = dataService.getMyClasses(student.id);
    final unsubmittedCount = assessmentService.assessments
        .where((a) => myClasses.any((c) => c.id == a.classId))
        .length;
    final avgProgress = myClasses.isEmpty
        ? 0.0
        : myClasses.fold(0.0, (sum, c) => sum + c.progress) / myClasses.length;

    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Expanded(
                child: Text('你好，${student.name}',
                    style: Theme.of(context)
                        .textTheme
                        .headlineMedium
                        ?.copyWith(fontWeight: FontWeight.bold)),
              ),
              IconButton(
                tooltip: '退出登录',
                icon: const Icon(Icons.logout),
                onPressed: onLogout,
              ),
            ],
          ),
          const SizedBox(height: 4),
          Text('权益：${student.plan.label} · 学号 ${student.id}',
              style: TextStyle(color: Colors.grey[600])),
          const SizedBox(height: 24),
          Row(
            children: [
              MetricCard(label: '我的课程', value: '${myClasses.length}', trend: '已选'),
              const SizedBox(width: 16),
              MetricCard(
                  label: '进行中班级',
                  value: '${dataService.activeClasses}',
                  trend: '平台'),
              const SizedBox(width: 16),
              MetricCard(
                  label: '待完成考核',
                  value: '$unsubmittedCount',
                  trend: unsubmittedCount > 0 ? '⚠' : '✓'),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              MetricCard(
                  label: '平均进度',
                  value: '${(avgProgress * 100).toInt()}%',
                  trend: avgProgress > 0 ? '继续加油' : '—'),
              const SizedBox(width: 16),
              MetricCard(
                  label: '即将上课',
                  value: '${dataService.upcomingSessions.length}',
                  trend: '本周'),
              const SizedBox(width: 16),
              MetricCard(
                  label: '学员数',
                  value: '${dataService.totalStudents}',
                  trend: '平台'),
            ],
          ),
          const SizedBox(height: 24),
          Row(
            children: [
              Expanded(
                child: SizedBox(
                  height: 260,
                  child: _buildSectionPanel(
                    context,
                    title: '我的课程',
                    emptyHint: '还未选课，去报名吧',
                    items: myClasses
                        .map((c) => _SimpleItem(
                              name: c.name,
                              subtitle: '${c.refName} · ${c.status.label} · '
                                  '${(c.progress * 100).toInt()}%',
                            ))
                        .toList(),
                    onViewAll: () => Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (_) => MyCoursesScreen(student: student),
                      ),
                    ),
                  ),
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: SizedBox(
                  height: 260,
                  child: _buildSectionPanel(
                    context,
                    title: '课表',
                    emptyHint: '暂无课次',
                    items: dataService.upcomingSessions
                        .take(5)
                        .map((s) => _SimpleItem(
                              name: s.lessonTitle,
                              subtitle:
                                  '${dataService.classNameOf(s.classId) ?? ''} · ${s.location}',
                            ))
                        .toList(),
                    onViewAll: () => Navigator.push(
                      context,
                      MaterialPageRoute(builder: (_) => const ScheduleScreen()),
                    ),
                  ),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _SimpleItem {
  final String name;
  final String subtitle;

  const _SimpleItem({required this.name, required this.subtitle});
}

Widget _buildSectionPanel(BuildContext context,
    {required String title,
    required String emptyHint,
    required List<_SimpleItem> items,
    VoidCallback? onViewAll}) {
  return Card(
    child: Padding(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Expanded(
                child: Text(title,
                    overflow: TextOverflow.ellipsis,
                    style: Theme.of(context)
                        .textTheme
                        .titleMedium
                        ?.copyWith(fontWeight: FontWeight.bold)),
              ),
              if (onViewAll != null)
                TextButton(
                  onPressed: onViewAll,
                  child: const Text('查看全部 →'),
                ),
            ],
          ),
          const SizedBox(height: 12),
          Expanded(
            child: items.isEmpty
                ? Center(
                    child: Text(emptyHint,
                        style: TextStyle(color: Colors.grey[500])))
                : ListView.separated(
                    itemCount: items.length,
                    separatorBuilder: (_, _) => const Divider(height: 1),
                    itemBuilder: (_, i) => ListTile(
                      dense: true,
                      title: Text(items[i].name),
                      subtitle: Text(items[i].subtitle,
                          style:
                              TextStyle(color: Colors.grey[600], fontSize: 12)),
                    ),
                  ),
          ),
        ],
      ),
    ),
  );
}
