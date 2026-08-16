// 后台立项列表（双创项目/立项申请，对齐规格 innovation-project.md 最小版：
// 表格展示项目名/姓名栏/方向/提交时间；v0.1 不做审批流）。

import 'package:flutter/material.dart';

import '../admin_api.dart';
import '../../models/application.dart';

class ApplicationsPage extends StatefulWidget {
  const ApplicationsPage({super.key, required this.api});

  final AdminApi api;

  @override
  State<ApplicationsPage> createState() => _ApplicationsPageState();
}

class _ApplicationsPageState extends State<ApplicationsPage> {
  List<Application>? _applications;
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final applications = await widget.api.fetchApplications();
      if (mounted) {
        setState(() {
          _applications = applications;
          _error = null;
        });
      }
    } catch (_) {
      if (mounted) {
        setState(() => _error = '加载立项列表失败，请稍后重试');
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
    final applications = _applications;
    if (applications == null) {
      return const Center(child: CircularProgressIndicator());
    }
    return ListView(
      padding: const EdgeInsets.all(24),
      children: [
        Text('立项申请', style: Theme.of(context).textTheme.headlineSmall),
        const SizedBox(height: 16),
        DataTable(
          columns: const [
            DataColumn(label: Text('项目名称')),
            DataColumn(label: Text('姓名栏')),
            DataColumn(label: Text('方向')),
            DataColumn(label: Text('申请人')),
            DataColumn(label: Text('提交时间')),
          ],
          rows: [
            for (final a in applications)
              DataRow(cells: [
                DataCell(SizedBox(width: 160, child: Text(a.projectName))),
                DataCell(SizedBox(width: 180, child: Text(a.memberLabel))),
                DataCell(Text(a.direction.isEmpty ? '-' : a.direction)),
                DataCell(Text(a.studentName)),
                DataCell(
                  Text(
                    a.createdAt == null ? '-' : a.createdAt!.substring(0, 16),
                    style: const TextStyle(fontSize: 12),
                  ),
                ),
              ]),
          ],
        ),
        if (applications.isEmpty) const Padding(
          padding: EdgeInsets.all(16),
          child: Center(child: Text('暂无立项申请')),
        ),
      ],
    );
  }
}
