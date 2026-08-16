// 后台学员管理（对齐原型 qt-students 学员表）：
// 学员 / 当前课程 / 进度(X/5 模块) / 最近活跃 / 立项(✓ 项目名) / 状态。

import 'package:flutter/material.dart';

import '../admin_api.dart';
import '../../models/learner.dart';

class LearnersPage extends StatefulWidget {
  const LearnersPage({super.key, required this.api});

  final AdminApi api;

  @override
  State<LearnersPage> createState() => _LearnersPageState();
}

class _LearnersPageState extends State<LearnersPage> {
  List<Learner>? _learners;
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final learners = await widget.api.fetchLearners();
      if (mounted) {
        setState(() {
          _learners = learners;
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
    final learners = _learners;
    if (learners == null) {
      return const Center(child: CircularProgressIndicator());
    }
    return ListView(
      padding: const EdgeInsets.all(24),
      children: [
        Text('学员管理', style: Theme.of(context).textTheme.headlineSmall),
        const SizedBox(height: 16),
        DataTable(
          columns: const [
            DataColumn(label: Text('学员')),
            DataColumn(label: Text('当前课程')),
            DataColumn(label: Text('进度')),
            DataColumn(label: Text('最近活跃')),
            DataColumn(label: Text('立项')),
            DataColumn(label: Text('状态')),
          ],
          rows: [
            for (final l in learners)
              DataRow(cells: [
                DataCell(Text(l.name)),
                DataCell(Text(l.course)),
                DataCell(Text('${l.progressMax}/${l.progressTotal} 模块')),
                DataCell(
                  Text(
                    l.activeAt == null ? '-' : l.activeAt!.substring(0, 16),
                    style: const TextStyle(fontSize: 12),
                  ),
                ),
                DataCell(
                  l.projectName == null || l.projectName!.isEmpty
                      ? const Text('—')
                      : Text('✓ ${l.projectName!}', style: const TextStyle(color: Colors.green)),
                ),
                DataCell(Text(l.status)),
              ]),
          ],
        ),
        if (learners.isEmpty) const Padding(
          padding: EdgeInsets.all(16),
          child: Center(child: Text('暂无学员（上报进度或提交立项后自动建档）')),
        ),
      ],
    );
  }
}
