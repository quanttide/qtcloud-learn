import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../models/enums.dart';
import '../services/learn_data_service.dart';
import 'classroom_screen.dart';

/// 课表页（自 `qtclass` schedule_screen.dart 迁移，数据源改为
/// `qtcloud-learn` 服务端 Session；className 由 classId 关联派生）。
class ScheduleScreen extends StatelessWidget {
  const ScheduleScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Consumer<LearnDataService>(
      builder: (context, service, _) {
        if (service.loading) {
          return const Scaffold(
              body: Center(child: CircularProgressIndicator()));
        }
        if (service.error != null && service.sessions.isEmpty) {
          return Scaffold(
            body: Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text('加载失败: ${service.error}'),
                  const SizedBox(height: 16),
                  ElevatedButton(
                      onPressed: () => service.load(),
                      child: const Text('重试')),
                ],
              ),
            ),
          );
        }
        return Scaffold(
          appBar: AppBar(title: const Text('课表')),
          body: ListView(
            children: [
              if (service.inProgressSessions.isNotEmpty) ...[
                _sectionHeader(context, '正在上课'),
                ...service.inProgressSessions
                    .map((s) => _SessionCard(sessionId: s.id, isHighlighted: true)),
              ],
              if (service.upcomingSessions.isNotEmpty) ...[
                _sectionHeader(context, '即将开始'),
                ...service.upcomingSessions
                    .map((s) => _SessionCard(sessionId: s.id)),
              ],
              if (service.completedSessions.isNotEmpty) ...[
                _sectionHeader(context, '历史课程'),
                ...service.completedSessions
                    .map((s) => _SessionCard(sessionId: s.id)),
              ],
              if (service.sessions.isEmpty)
                const Padding(
                  padding: EdgeInsets.all(32),
                  child: Center(child: Text('暂无课次')),
                ),
            ],
          ),
        );
      },
    );
  }

  Widget _sectionHeader(BuildContext context, String title) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 4),
      child: Text(title,
          style: Theme.of(context)
              .textTheme
              .titleMedium
              ?.copyWith(fontWeight: FontWeight.bold)),
    );
  }
}

class _SessionCard extends StatelessWidget {
  final String sessionId;
  final bool isHighlighted;

  const _SessionCard({required this.sessionId, this.isHighlighted = false});

  @override
  Widget build(BuildContext context) {
    final service = context.watch<LearnDataService>();
    final sessions =
        service.sessions.where((s) => s.id == sessionId).toList();
    if (sessions.isEmpty) return const SizedBox.shrink();
    final session = sessions.first;
    final className = service.classNameOf(session.classId) ?? '';
    final teacherName = service.teacherById(session.teacherId)?.name ?? '';
    final timeStr =
        '${session.startTime.hour.toString().padLeft(2, '0')}:${session.startTime.minute.toString().padLeft(2, '0')}';
    final isInProgress = session.status == SessionStatus.inProgress;

    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
      color: isHighlighted
          ? Theme.of(context).colorScheme.primaryContainer
          : null,
      child: ListTile(
        leading: CircleAvatar(
          backgroundColor: isInProgress ? Colors.green : Colors.grey.shade300,
          child: Icon(
            isInProgress ? Icons.play_arrow : Icons.check,
            color: Colors.white,
          ),
        ),
        title: Text(session.lessonTitle,
            style: const TextStyle(fontWeight: FontWeight.w600)),
        subtitle: Text(
            '$className · $timeStr · ${session.location}\n授课：$teacherName'),
        isThreeLine: true,
        trailing: isInProgress
            ? ElevatedButton(
                onPressed: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                        builder: (_) => ClassroomScreen(sessionId: session.id)),
                  );
                },
                child: const Text('进入课堂'),
              )
            : null,
        onTap: () {
          Navigator.push(
            context,
            MaterialPageRoute(
                builder: (_) => ClassroomScreen(sessionId: session.id)),
          );
        },
      ),
    );
  }
}
