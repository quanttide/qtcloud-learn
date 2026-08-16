// 后台概览：学员数 / 进度完成数 / 立项数统计卡（对齐规格 dashboard.md 最小版）。

import 'package:flutter/material.dart';

import '../admin_api.dart';
import '../../models/application.dart';
import '../../models/progress.dart';
import '../../models/student.dart';

class DashboardPage extends StatefulWidget {
  const DashboardPage({super.key, required this.api});

  final AdminApi api;

  @override
  State<DashboardPage> createState() => _DashboardPageState();
}

class _DashboardPageState extends State<DashboardPage> {
  List<Student>? _students;
  List<Progress>? _progress;
  List<Application>? _applications;
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final results = await Future.wait([
        widget.api.fetchStudents(),
        widget.api.fetchProgress(),
        widget.api.fetchApplications(),
      ]);
      if (mounted) {
        setState(() {
          _students = results[0] as List<Student>;
          _progress = results[1] as List<Progress>;
          _applications = results[2] as List<Application>;
          _error = null;
        });
      }
    } catch (_) {
      if (mounted) {
        setState(() => _error = '加载概览失败，请稍后重试');
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_error != null) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(_error!),
            const SizedBox(height: 12),
            FilledButton(onPressed: _load, child: const Text('重试')),
          ],
        ),
      );
    }
    final students = _students;
    final progress = _progress;
    final applications = _applications;
    if (students == null || progress == null || applications == null) {
      return const Center(child: CircularProgressIndicator());
    }
    final finished = progress.where((p) => p.finished).length;
    return ListView(
      padding: const EdgeInsets.all(24),
      children: [
        Text('概览', style: Theme.of(context).textTheme.headlineSmall),
        const SizedBox(height: 16),
        Row(
          children: [
            _StatCard(
              icon: Icons.person_outline,
              label: '学员数',
              value: '${students.length}',
            ),
            const SizedBox(width: 16),
            _StatCard(
              icon: Icons.trending_up,
              label: '进度记录',
              value: '${progress.length}（完成 $finished）',
            ),
            const SizedBox(width: 16),
            _StatCard(
              icon: Icons.rocket_launch_outlined,
              label: '立项申请',
              value: '${applications.length}',
            ),
          ],
        ),
      ],
    );
  }
}

class _StatCard extends StatelessWidget {
  const _StatCard({required this.icon, required this.label, required this.value});

  final IconData icon;
  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Icon(icon),
            const SizedBox(height: 8),
            Text(value, style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
            Text(label, style: const TextStyle(color: Colors.grey, fontSize: 12)),
          ],
        ),
      ),
    );
  }
}
