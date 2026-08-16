// 后台立项管理（对齐原型 qt-proposals）：
// 项目名称 / 学员（组队=队长）/ 方向 / 组队 / 核心假设 / 提交时间 / 状态 / 操作（删除=软删）。
// 顶部「📜 历史记录」开关展开软删除历史。

import 'package:flutter/material.dart';

import '../admin_api.dart';
import '../../models/application.dart';

class ProposalsPage extends StatefulWidget {
  const ProposalsPage({super.key, required this.api});

  final AdminApi api;

  @override
  State<ProposalsPage> createState() => _ProposalsPageState();
}

class _ProposalsPageState extends State<ProposalsPage> {
  List<Application>? _proposals;
  List<Application> _history = [];
  bool _showHistory = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final proposals = await widget.api.fetchProposals();
      final history = await widget.api.fetchHistory();
      if (mounted) {
        setState(() {
          _proposals = proposals;
          _history = history;
          _error = null;
        });
      }
    } catch (_) {
      if (mounted) {
        setState(() => _error = '加载立项列表失败，请稍后重试');
      }
    }
  }

  Future<void> _delete(Application app) async {
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('删除立项'),
        content: Text('确定删除「${app.projectName}」？删除后可在历史记录中找回。'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('取消')),
          FilledButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('删除')),
        ],
      ),
    );
    if (ok != true) return;
    try {
      await widget.api.deleteProposal(app.id);
      await _load();
    } catch (_) {
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(const SnackBar(content: Text('删除失败，请稍后重试')));
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
    final proposals = _proposals;
    if (proposals == null) {
      return const Center(child: CircularProgressIndicator());
    }
    return ListView(
      padding: const EdgeInsets.all(24),
      children: [
        Row(
          children: [
            Text('立项管理', style: Theme.of(context).textTheme.headlineSmall),
            const Spacer(),
            TextButton.icon(
              onPressed: () => setState(() => _showHistory = !_showHistory),
              icon: const Icon(Icons.history, size: 18),
              label: Text('📜 历史记录${_showHistory ? '' : '（${_history.length}）'}'),
            ),
          ],
        ),
        const SizedBox(height: 16),
        DataTable(
          columns: const [
            DataColumn(label: Text('项目名称')),
            DataColumn(label: Text('学员')),
            DataColumn(label: Text('方向')),
            DataColumn(label: Text('组队')),
            DataColumn(label: Text('核心假设')),
            DataColumn(label: Text('提交时间')),
            DataColumn(label: Text('状态')),
            DataColumn(label: Text('操作')),
          ],
          rows: [
            for (final a in proposals)
              DataRow(cells: [
                DataCell(SizedBox(width: 150, child: Text(a.projectName))),
                DataCell(Text(a.studentName)),
                DataCell(Text(a.directionType.isEmpty ? '-' : a.directionType)),
                DataCell(SizedBox(width: 170, child: Text(a.teamLabel))),
                DataCell(SizedBox(width: 160, child: Text(a.hypothesis, maxLines: 2, overflow: TextOverflow.ellipsis))),
                DataCell(
                  Text(
                    a.submittedAt == null ? '-' : a.submittedAt!.substring(0, 16),
                    style: const TextStyle(fontSize: 12),
                  ),
                ),
                DataCell(Text(a.status)),
                DataCell(
                  IconButton(
                    tooltip: '删除',
                    icon: const Icon(Icons.delete_outline),
                    onPressed: () => _delete(a),
                  ),
                ),
              ]),
          ],
        ),
        if (proposals.isEmpty) const Padding(
          padding: EdgeInsets.all(16),
          child: Center(child: Text('暂无立项申请')),
        ),
        if (_showHistory) ...[
          const SizedBox(height: 24),
          Text('历史记录（软删除）', style: Theme.of(context).textTheme.titleMedium),
          const SizedBox(height: 8),
          if (_history.isEmpty)
            const Padding(padding: EdgeInsets.all(8), child: Text('暂无历史记录'))
          else
            DataTable(
              columns: const [
                DataColumn(label: Text('项目名称')),
                DataColumn(label: Text('学员')),
                DataColumn(label: Text('方向')),
                DataColumn(label: Text('提交时间')),
                DataColumn(label: Text('删除时间')),
              ],
              rows: [
                for (final a in _history)
                  DataRow(cells: [
                    DataCell(Text(a.projectName)),
                    DataCell(Text(a.studentName)),
                    DataCell(Text(a.directionType)),
                    DataCell(
                      Text(
                        a.submittedAt == null ? '-' : a.submittedAt!.substring(0, 16),
                        style: const TextStyle(fontSize: 12),
                      ),
                    ),
                    DataCell(
                      Text(
                        a.deletedAt == null ? '-' : a.deletedAt!.substring(0, 16),
                        style: const TextStyle(fontSize: 12),
                      ),
                    ),
                  ]),
              ],
            ),
        ],
      ],
    );
  }
}
