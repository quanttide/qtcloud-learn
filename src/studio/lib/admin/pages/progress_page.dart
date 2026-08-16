// 后台进度查看（学员上报的进度记录表格）。

import 'package:flutter/material.dart';

import '../admin_api.dart';
import '../../models/progress.dart';

class ProgressPage extends StatefulWidget {
  const ProgressPage({super.key, required this.api});

  final AdminApi api;

  @override
  State<ProgressPage> createState() => _ProgressPageState();
}

class _ProgressPageState extends State<ProgressPage> {
  List<Progress>? _progress;
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final progress = await widget.api.fetchProgress();
      if (mounted) {
        setState(() {
          _progress = progress;
          _error = null;
        });
      }
    } catch (_) {
      if (mounted) {
        setState(() => _error = '加载进度失败，请稍后重试');
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
    final progress = _progress;
    if (progress == null) {
      return const Center(child: CircularProgressIndicator());
    }
    return ListView(
      padding: const EdgeInsets.all(24),
      children: [
        Text('进度上报', style: Theme.of(context).textTheme.headlineSmall),
        const SizedBox(height: 16),
        DataTable(
          columns: const [
            DataColumn(label: Text('学员 ID')),
            DataColumn(label: Text('班级')),
            DataColumn(label: Text('进度')),
            DataColumn(label: Text('状态')),
            DataColumn(label: Text('更新时间')),
          ],
          rows: [
            for (final p in progress)
              DataRow(cells: [
                DataCell(Text(p.studentId, style: const TextStyle(fontSize: 12))),
                DataCell(Text(p.classId)),
                DataCell(
                  Row(
                    children: [
                      SizedBox(
                        width: 120,
                        child: LinearProgressIndicator(value: p.percent),
                      ),
                      const SizedBox(width: 8),
                      Text('${(p.percent * 100).toStringAsFixed(0)}%'),
                    ],
                  ),
                ),
                DataCell(Text(p.finished ? '已完成' : '学习中')),
                DataCell(Text(p.updatedAt ?? '-', style: const TextStyle(fontSize: 12))),
              ]),
          ],
        ),
        if (progress.isEmpty) const Padding(
          padding: EdgeInsets.all(16),
          child: Center(child: Text('暂无进度记录')),
        ),
      ],
    );
  }
}
