import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../models/enums.dart';
import '../services/learn_data_service.dart';

/// 课堂 / 考勤页（自 `qtclass` classroom_screen.dart 迁移，
/// 学员列表与考勤改为服务端 Session.Attendances 数据）。
class ClassroomScreen extends StatelessWidget {
  final String sessionId;

  const ClassroomScreen({super.key, required this.sessionId});

  @override
  Widget build(BuildContext context) {
    return Consumer<LearnDataService>(
      builder: (context, service, _) {
        final sessions =
            service.sessions.where((s) => s.id == sessionId).toList();
        if (sessions.isEmpty) {
          return const Scaffold(body: Center(child: Text('课次不存在')));
        }
        final session = sessions.first;
        final className = service.classNameOf(session.classId) ?? '';
        final isInProgress = session.status == SessionStatus.inProgress;

        return Scaffold(
          appBar: AppBar(
            title: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(session.lessonTitle,
                    style: const TextStyle(fontSize: 16)),
                Text(className,
                    style: const TextStyle(fontSize: 12)),
              ],
            ),
          ),
          body: isInProgress
              ? _InSessionView(sessionId: sessionId)
              : _ReviewView(sessionId: sessionId),
        );
      },
    );
  }
}

class _InSessionView extends StatelessWidget {
  final String sessionId;

  const _InSessionView({required this.sessionId});

  @override
  Widget build(BuildContext context) {
    final service = context.watch<LearnDataService>();
    final session = service.sessions.firstWhere((s) => s.id == sessionId);
    final progress = _calculateProgress(session);

    return Column(
      children: [
        LinearProgressIndicator(value: progress, minHeight: 4),
        Expanded(
          child: Row(
            children: [
              Expanded(flex: 3, child: _ContentArea(sessionId: sessionId)),
              _StudentSidebar(sessionId: sessionId),
            ],
          ),
        ),
        _Toolbar(sessionId: sessionId),
      ],
    );
  }

  double _calculateProgress(dynamic session) {
    final now = DateTime.now();
    final elapsed = now.difference(session.startTime).inMinutes.toDouble();
    return (elapsed / session.durationMinutes).clamp(0.0, 1.0);
  }
}

class _ContentArea extends StatelessWidget {
  final String sessionId;

  const _ContentArea({required this.sessionId});

  @override
  Widget build(BuildContext context) {
    final service = context.watch<LearnDataService>();
    final session = service.sessions.firstWhere((s) => s.id == sessionId);
    return Container(
      margin: const EdgeInsets.all(8),
      decoration: BoxDecoration(
        color: Colors.grey.shade100,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.slideshow, size: 64, color: Colors.grey.shade400),
            const SizedBox(height: 16),
            Text('课件展示区',
                style: TextStyle(fontSize: 18, color: Colors.grey.shade600)),
            const SizedBox(height: 8),
            Text(session.lessonTitle,
                style: TextStyle(color: Colors.grey.shade500)),
          ],
        ),
      ),
    );
  }
}

/// 学员侧栏：由 Session.Attendances 驱动（学员名自服务解析）。
class _StudentSidebar extends StatelessWidget {
  final String sessionId;

  const _StudentSidebar({required this.sessionId});

  @override
  Widget build(BuildContext context) {
    final service = context.watch<LearnDataService>();
    final session = service.sessions.firstWhere((s) => s.id == sessionId);
    final rows = _attendances(service, session);

    return Container(
      width: 200,
      margin: const EdgeInsets.all(8),
      decoration: BoxDecoration(
        border: Border.all(color: Colors.grey.shade300),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Column(
        children: [
          Container(
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(
              color: Colors.grey.shade100,
              borderRadius: const BorderRadius.vertical(top: Radius.circular(7)),
            ),
            child: Row(
              children: [
                const Icon(Icons.people, size: 16),
                const SizedBox(width: 4),
                Text('学员 (${rows.length})',
                    style: const TextStyle(fontWeight: FontWeight.w600)),
              ],
            ),
          ),
          Expanded(
            child: ListView(
              children: rows
                  .map((r) => _StudentTile(
                        sessionId: sessionId,
                        studentId: r.$1,
                        name: r.$2,
                        status: r.$3,
                      ))
                  .toList(),
            ),
          ),
        ],
      ),
    );
  }

  List<(String, String, AttendanceStatus)> _attendances(
      LearnDataService service, dynamic session) {
    if (session.attendances.isNotEmpty) {
      return session.attendances.map((a) {
        final name = service.students
            .where((s) => s.id == a.studentId)
            .fold<String>(a.studentId, (_, s) => s.name);
        return (a.studentId, name, a.status);
      }).toList();
    }
    // 兜底：班级成员列表
    final classStudents = service.getStudentsByClass(session.classId);
    return classStudents
        .map((s) => (s.id, s.name, AttendanceStatus.unknown))
        .toList();
  }
}

class _StudentTile extends StatelessWidget {
  final String sessionId;
  final String studentId;
  final String name;
  final AttendanceStatus status;

  const _StudentTile({
    required this.sessionId,
    required this.studentId,
    required this.name,
    required this.status,
  });

  @override
  Widget build(BuildContext context) {
    final service = context.read<LearnDataService>();
    Color bgColor;
    Widget? trailing;
    switch (status) {
      case AttendanceStatus.present:
        bgColor = Colors.green.shade50;
        trailing = const Icon(Icons.check_circle, color: Colors.green, size: 18);
      case AttendanceStatus.late:
        bgColor = Colors.orange.shade50;
        trailing = const Icon(Icons.access_time, color: Colors.orange, size: 18);
      case AttendanceStatus.absent:
        bgColor = Colors.red.shade50;
        trailing = const Icon(Icons.cancel, color: Colors.red, size: 18);
      case AttendanceStatus.unknown:
        bgColor = Colors.transparent;
        trailing = null;
    }

    return InkWell(
      onTap: () => _showStatusMenu(context, service),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 6),
        color: bgColor,
        child: Row(
          children: [
            CircleAvatar(
              radius: 14,
              child: Text(name.isNotEmpty ? name[0] : '?',
                  style: const TextStyle(fontSize: 12)),
            ),
            const SizedBox(width: 8),
            Expanded(child: Text(name, style: const TextStyle(fontSize: 13))),
            ?trailing,
          ],
        ),
      ),
    );
  }

  void _showStatusMenu(BuildContext context, LearnDataService service) {
    showModalBottomSheet(
      context: context,
      builder: (ctx) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Padding(
              padding: const EdgeInsets.all(16),
              child: Text('标记 $name 的出勤',
                  style: const TextStyle(
                      fontWeight: FontWeight.bold, fontSize: 16)),
            ),
            for (final s in AttendanceStatus.values)
              if (s != AttendanceStatus.unknown)
                ListTile(
                  leading: Icon(_icon(s), color: _color(s)),
                  title: Text(s.label),
                  onTap: () {
                    service.markAttendance(
                        sessionId, studentId, s);
                    Navigator.pop(ctx);
                  },
                ),
          ],
        ),
      ),
    );
  }

  IconData _icon(AttendanceStatus s) => switch (s) {
        AttendanceStatus.present => Icons.check_circle,
        AttendanceStatus.late => Icons.access_time,
        AttendanceStatus.absent => Icons.cancel,
        AttendanceStatus.unknown => Icons.help,
      };

  Color _color(AttendanceStatus s) => switch (s) {
        AttendanceStatus.present => Colors.green,
        AttendanceStatus.late => Colors.orange,
        AttendanceStatus.absent => Colors.red,
        AttendanceStatus.unknown => Colors.grey,
      };
}

class _Toolbar extends StatelessWidget {
  final String sessionId;

  const _Toolbar({required this.sessionId});

  void _showAttendanceMenu(BuildContext context) {
    final service = context.read<LearnDataService>();
    showModalBottomSheet(
      context: context,
      builder: (ctx) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Padding(
              padding: EdgeInsets.all(16),
              child: Text('出勤标记',
                  style:
                      TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
            ),
            ListTile(
              leading: const Icon(Icons.check_circle, color: Colors.green),
              title: const Text('全部正常'),
              onTap: () {
                service.markAllAttendance(sessionId, AttendanceStatus.present);
                Navigator.pop(ctx);
              },
            ),
            ListTile(
              leading: const Icon(Icons.access_time, color: Colors.orange),
              title: const Text('全部迟到'),
              onTap: () {
                service.markAllAttendance(sessionId, AttendanceStatus.late);
                Navigator.pop(ctx);
              },
            ),
            ListTile(
              leading: const Icon(Icons.cancel, color: Colors.red),
              title: const Text('全部缺勤'),
              onTap: () {
                service.markAllAttendance(sessionId, AttendanceStatus.absent);
                Navigator.pop(ctx);
              },
            ),
            const Divider(),
            const Padding(
              padding: EdgeInsets.all(16),
              child: Text('点击学员可单独标记',
                  style: TextStyle(color: Colors.grey, fontSize: 12)),
            ),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      decoration: BoxDecoration(
        color: Colors.white,
        border: Border(top: BorderSide(color: Colors.grey.shade300)),
      ),
      child: Row(
        children: [
          _ToolbarButton(
            icon: Icons.edit_note,
            label: '出勤',
            onTap: () => _showAttendanceMenu(context),
          ),
          const SizedBox(width: 16),
          _ToolbarButton(
            icon: Icons.question_answer,
            label: '提问',
            onTap: () => ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text('提问功能开发中')),
            ),
          ),
          const SizedBox(width: 16),
          _ToolbarButton(
            icon: Icons.poll,
            label: '投票',
            onTap: () => ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text('投票功能开发中')),
            ),
          ),
          const Spacer(),
          ElevatedButton.icon(
            onPressed: () => showDialog(
              context: context,
              builder: (ctx) => AlertDialog(
                title: const Text('结束课程'),
                content: const Text('确认结束本节课吗？'),
                actions: [
                  TextButton(
                      onPressed: () => Navigator.pop(ctx),
                      child: const Text('取消')),
                  ElevatedButton(
                      onPressed: () => Navigator.pop(ctx),
                      child: const Text('下课')),
                ],
              ),
            ),
            icon: const Icon(Icons.stop),
            label: const Text('下课'),
            style: ElevatedButton.styleFrom(
                backgroundColor: Colors.red, foregroundColor: Colors.white),
          ),
        ],
      ),
    );
  }
}

class _ToolbarButton extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;

  const _ToolbarButton(
      {required this.icon, required this.label, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(8),
      child: Padding(
        padding: const EdgeInsets.all(8),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(icon, size: 20),
            const SizedBox(height: 2),
            Text(label, style: const TextStyle(fontSize: 12)),
          ],
        ),
      ),
    );
  }
}

class _ReviewView extends StatelessWidget {
  final String sessionId;

  const _ReviewView({required this.sessionId});

  @override
  Widget build(BuildContext context) {
    final service = context.watch<LearnDataService>();
    final session = service.sessions.firstWhere((s) => s.id == sessionId);
    final teacherName = service.teacherById(session.teacherId)?.name ?? '';
    final rows = _attendanceRows(service, session);

    int count(AttendanceStatus s) =>
        rows.where((r) => r.$2 == s).length;

    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(session.lessonTitle,
                    style: const TextStyle(
                        fontSize: 18, fontWeight: FontWeight.bold)),
                const SizedBox(height: 8),
                Text(service.classNameOf(session.classId) ?? ''),
                Text('授课教师: $teacherName'),
                Text(
                    '${session.startTime.year}-${session.startTime.month.toString().padLeft(2, '0')}-${session.startTime.day.toString().padLeft(2, '0')} · ${session.location}'),
              ],
            ),
          ),
        ),
        const SizedBox(height: 16),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text('出勤统计',
                    style:
                        TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                const SizedBox(height: 12),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceAround,
                  children: [
                    _StatItem(label: '应到', count: rows.length, color: Colors.grey),
                    _StatItem(
                        label: '正常',
                        count: count(AttendanceStatus.present),
                        color: Colors.green),
                    _StatItem(
                        label: '迟到',
                        count: count(AttendanceStatus.late),
                        color: Colors.orange),
                    _StatItem(
                        label: '缺勤',
                        count: count(AttendanceStatus.absent),
                        color: Colors.red),
                    _StatItem(
                        label: '未记',
                        count: count(AttendanceStatus.unknown),
                        color: Colors.grey.shade400),
                  ],
                ),
              ],
            ),
          ),
        ),
        const SizedBox(height: 16),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text('学员列表',
                    style:
                        TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                const SizedBox(height: 8),
                ...rows.map((r) => ListTile(
                      dense: true,
                      leading: CircleAvatar(
                          child: Text(r.$1.isNotEmpty ? r.$1[0] : '?')),
                      title: Text(r.$1),
                      trailing: Chip(
                          label: Text(r.$2.label,
                              style: const TextStyle(fontSize: 12))),
                    )),
              ],
            ),
          ),
        ),
      ],
    );
  }

  List<(String, AttendanceStatus)> _attendanceRows(
      LearnDataService service, dynamic session) {
    return session.attendances.isEmpty
        ? service
            .getStudentsByClass(session.classId)
            .map((s) => (s.name, AttendanceStatus.unknown))
            .toList()
        : session.attendances.map((a) {
            final name = service.students
                .where((s) => s.id == a.studentId)
                .fold<String>(a.studentId, (_, s) => s.name);
            return (name, a.status);
          }).toList();
  }
}

class _StatItem extends StatelessWidget {
  final String label;
  final int count;
  final Color color;

  const _StatItem(
      {required this.label, required this.count, required this.color});

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Text('$count',
            style: TextStyle(
                fontSize: 24, fontWeight: FontWeight.bold, color: color)),
        Text(label, style: const TextStyle(fontSize: 13)),
      ],
    );
  }
}
