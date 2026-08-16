// 后台学员列表（对齐规格 student-management.md 最小版：表格查看）。

import 'package:flutter/material.dart';

import '../admin_api.dart';
import '../../models/student.dart';

class StudentsPage extends StatefulWidget {
  const StudentsPage({super.key, required this.api});

  final AdminApi api;

  @override
  State<StudentsPage> createState() => _StudentsPageState();
}

class _StudentsPageState extends State<StudentsPage> {
  List<Student>? _students;
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final students = await widget.api.fetchStudents();
      if (mounted) {
        setState(() {
          _students = students;
          _error = null;
        });
      }
    } catch (_) {
      if (mounted) {
        setState(() => _error = '加载学员列表失败，请稍后重试');
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
    if (students == null) {
      return const Center(child: CircularProgressIndicator());
    }
    return ListView(
      padding: const EdgeInsets.all(24),
      children: [
        Text('学员管理', style: Theme.of(context).textTheme.headlineSmall),
        const SizedBox(height: 16),
        DataTable(
          columns: const [
            DataColumn(label: Text('姓名')),
            DataColumn(label: Text('邮箱')),
            DataColumn(label: Text('ID')),
          ],
          rows: [
            for (final s in students)
              DataRow(cells: [
                DataCell(Text(s.name)),
                DataCell(Text(s.email)),
                DataCell(Text(s.id, style: const TextStyle(fontSize: 12))),
              ]),
          ],
        ),
        if (students.isEmpty) const Padding(
          padding: EdgeInsets.all(16),
          child: Center(child: Text('暂无学员')),
        ),
      ],
    );
  }
}
