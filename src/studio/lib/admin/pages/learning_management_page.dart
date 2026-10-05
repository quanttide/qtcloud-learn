import 'package:flutter/material.dart';

import '../../models/learning.dart';
import '../admin_api.dart';

class LearningManagementPage extends StatefulWidget {
  const LearningManagementPage({super.key, required this.api});

  final AdminApi api;

  @override
  State<LearningManagementPage> createState() => _LearningManagementPageState();
}

class _LearningManagementPageState extends State<LearningManagementPage> {
  final _taskId = TextEditingController();
  final _taskTitle = TextEditingController();
  final _taskDescription = TextEditingController();
  final _scheduleId = TextEditingController();
  final _scheduleTitle = TextEditingController();
  final _scheduleDescription = TextEditingController();

  List<LearningTask>? _tasks;
  List<LearningSchedule>? _schedules;
  final Set<String> _selectedTaskIds = {};
  String? _error;
  bool _saving = false;

  @override
  void initState() {
    super.initState();
    _load();
  }

  @override
  void dispose() {
    _taskId.dispose();
    _taskTitle.dispose();
    _taskDescription.dispose();
    _scheduleId.dispose();
    _scheduleTitle.dispose();
    _scheduleDescription.dispose();
    super.dispose();
  }

  Future<void> _load() async {
    try {
      final results = await Future.wait([
        widget.api.fetchTasks(),
        widget.api.fetchSchedules(),
      ]);
      if (!mounted) return;
      setState(() {
        _tasks = results[0] as List<LearningTask>;
        _schedules = results[1] as List<LearningSchedule>;
        _error = null;
      });
    } catch (_) {
      if (mounted) setState(() => _error = '加载训练营/任务失败，请稍后重试');
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
    final tasks = _tasks;
    final schedules = _schedules;
    if (tasks == null || schedules == null) {
      return const Center(child: CircularProgressIndicator());
    }
    return DefaultTabController(
      length: 2,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(24, 24, 24, 0),
            child: Row(
              children: [
                Text(
                  '训练营与任务',
                  style: Theme.of(context).textTheme.headlineSmall,
                ),
                const Spacer(),
                IconButton(
                  tooltip: '刷新',
                  onPressed: _load,
                  icon: const Icon(Icons.refresh),
                ),
              ],
            ),
          ),
          const TabBar(
            tabs: [
              Tab(icon: Icon(Icons.route_outlined), text: '训练营'),
              Tab(icon: Icon(Icons.task_alt_outlined), text: '任务'),
            ],
          ),
          Expanded(
            child: TabBarView(
              children: [
                _ScheduleTab(
                  schedules: schedules,
                  tasks: tasks,
                  selectedTaskIds: _selectedTaskIds,
                  id: _scheduleId,
                  title: _scheduleTitle,
                  description: _scheduleDescription,
                  saving: _saving,
                  onSelect: _selectSchedule,
                  onToggleTask: _toggleTask,
                  onClear: _clearSchedule,
                  onSave: _saveSchedule,
                  onDelete: _deleteSchedule,
                ),
                _TaskTab(
                  tasks: tasks,
                  id: _taskId,
                  title: _taskTitle,
                  description: _taskDescription,
                  saving: _saving,
                  onSelect: _selectTask,
                  onClear: _clearTask,
                  onSave: _saveTask,
                  onDelete: _deleteTask,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  void _selectTask(LearningTask task) {
    setState(() {
      _taskId.text = task.id;
      _taskTitle.text = task.title;
      _taskDescription.text = task.description;
    });
  }

  void _clearTask() {
    setState(() {
      _taskId.clear();
      _taskTitle.clear();
      _taskDescription.clear();
    });
  }

  Future<void> _saveTask() async {
    final task = LearningTask(
      id: _taskId.text.trim(),
      title: _taskTitle.text.trim(),
      description: _taskDescription.text.trim(),
    );
    await _mutate(() => widget.api.saveTask(task));
    _clearTask();
  }

  Future<void> _deleteTask(String id) async =>
      _mutate(() => widget.api.deleteTask(id));

  void _selectSchedule(LearningSchedule schedule) {
    setState(() {
      _scheduleId.text = schedule.id;
      _scheduleTitle.text = schedule.title;
      _scheduleDescription.text = schedule.description ?? '';
      _selectedTaskIds
        ..clear()
        ..addAll(schedule.tasks.map((task) => task.id));
    });
  }

  void _clearSchedule() {
    setState(() {
      _scheduleId.clear();
      _scheduleTitle.clear();
      _scheduleDescription.clear();
      _selectedTaskIds.clear();
    });
  }

  void _toggleTask(String id, bool selected) {
    setState(() {
      if (selected) {
        _selectedTaskIds.add(id);
      } else {
        _selectedTaskIds.remove(id);
      }
    });
  }

  Future<void> _saveSchedule() async {
    final taskLookup = {
      for (final task in _tasks ?? <LearningTask>[]) task.id: task,
    };
    final tasks = [
      for (final id in _selectedTaskIds)
        if (taskLookup[id] != null) taskLookup[id]!,
    ];
    final schedule = LearningSchedule(
      id: _scheduleId.text.trim(),
      title: _scheduleTitle.text.trim(),
      description: _scheduleDescription.text.trim(),
      tasks: tasks,
    );
    await _mutate(() => widget.api.saveSchedule(schedule));
    _clearSchedule();
  }

  Future<void> _deleteSchedule(String id) async =>
      _mutate(() => widget.api.deleteSchedule(id));

  Future<void> _mutate(Future<Object?> Function() operation) async {
    setState(() => _saving = true);
    try {
      await operation();
      await _load();
    } catch (_) {
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(const SnackBar(content: Text('保存失败')));
      }
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }
}

class _TaskTab extends StatelessWidget {
  const _TaskTab({
    required this.tasks,
    required this.id,
    required this.title,
    required this.description,
    required this.saving,
    required this.onSelect,
    required this.onClear,
    required this.onSave,
    required this.onDelete,
  });

  final List<LearningTask> tasks;
  final TextEditingController id;
  final TextEditingController title;
  final TextEditingController description;
  final bool saving;
  final ValueChanged<LearningTask> onSelect;
  final VoidCallback onClear;
  final Future<void> Function() onSave;
  final Future<void> Function(String id) onDelete;

  @override
  Widget build(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.all(24),
      children: [
        _TaskForm(
          id: id,
          title: title,
          description: description,
          saving: saving,
          onClear: onClear,
          onSave: onSave,
        ),
        const SizedBox(height: 24),
        DataTable(
          columns: const [
            DataColumn(label: Text('ID')),
            DataColumn(label: Text('标题')),
            DataColumn(label: Text('描述')),
            DataColumn(label: Text('操作')),
          ],
          rows: [
            for (final task in tasks)
              DataRow(
                cells: [
                  DataCell(Text(task.id)),
                  DataCell(Text(task.title)),
                  DataCell(
                    SizedBox(
                      width: 420,
                      child: Text(task.description, maxLines: 2),
                    ),
                  ),
                  DataCell(
                    Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        IconButton(
                          tooltip: '编辑',
                          onPressed: () => onSelect(task),
                          icon: const Icon(Icons.edit_outlined),
                        ),
                        IconButton(
                          tooltip: '删除',
                          onPressed: () => onDelete(task.id),
                          icon: const Icon(Icons.delete_outline),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
          ],
        ),
      ],
    );
  }
}

class _ScheduleTab extends StatelessWidget {
  const _ScheduleTab({
    required this.schedules,
    required this.tasks,
    required this.selectedTaskIds,
    required this.id,
    required this.title,
    required this.description,
    required this.saving,
    required this.onSelect,
    required this.onToggleTask,
    required this.onClear,
    required this.onSave,
    required this.onDelete,
  });

  final List<LearningSchedule> schedules;
  final List<LearningTask> tasks;
  final Set<String> selectedTaskIds;
  final TextEditingController id;
  final TextEditingController title;
  final TextEditingController description;
  final bool saving;
  final ValueChanged<LearningSchedule> onSelect;
  final void Function(String id, bool selected) onToggleTask;
  final VoidCallback onClear;
  final Future<void> Function() onSave;
  final Future<void> Function(String id) onDelete;

  @override
  Widget build(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.all(24),
      children: [
        _ScheduleForm(
          id: id,
          title: title,
          description: description,
          tasks: tasks,
          selectedTaskIds: selectedTaskIds,
          saving: saving,
          onToggleTask: onToggleTask,
          onClear: onClear,
          onSave: onSave,
        ),
        const SizedBox(height: 24),
        DataTable(
          columns: const [
            DataColumn(label: Text('ID')),
            DataColumn(label: Text('标题')),
            DataColumn(label: Text('任务数')),
            DataColumn(label: Text('描述')),
            DataColumn(label: Text('操作')),
          ],
          rows: [
            for (final schedule in schedules)
              DataRow(
                cells: [
                  DataCell(Text(schedule.id)),
                  DataCell(Text(schedule.title)),
                  DataCell(Text('${schedule.tasks.length}')),
                  DataCell(
                    SizedBox(
                      width: 360,
                      child: Text(schedule.description ?? '', maxLines: 2),
                    ),
                  ),
                  DataCell(
                    Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        IconButton(
                          tooltip: '编辑',
                          onPressed: () => onSelect(schedule),
                          icon: const Icon(Icons.edit_outlined),
                        ),
                        IconButton(
                          tooltip: '删除',
                          onPressed: () => onDelete(schedule.id),
                          icon: const Icon(Icons.delete_outline),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
          ],
        ),
      ],
    );
  }
}

class _TaskForm extends StatelessWidget {
  const _TaskForm({
    required this.id,
    required this.title,
    required this.description,
    required this.saving,
    required this.onClear,
    required this.onSave,
  });

  final TextEditingController id;
  final TextEditingController title;
  final TextEditingController description;
  final bool saving;
  final VoidCallback onClear;
  final Future<void> Function() onSave;

  @override
  Widget build(BuildContext context) {
    return _Panel(
      child: Column(
        children: [
          Row(
            children: [
              Expanded(
                child: TextField(
                  controller: id,
                  decoration: const InputDecoration(labelText: '任务 ID'),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: TextField(
                  controller: title,
                  decoration: const InputDecoration(labelText: '标题'),
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          TextField(
            controller: description,
            minLines: 3,
            maxLines: 5,
            decoration: const InputDecoration(labelText: '描述'),
          ),
          const SizedBox(height: 12),
          _Actions(saving: saving, onClear: onClear, onSave: onSave),
        ],
      ),
    );
  }
}

class _ScheduleForm extends StatelessWidget {
  const _ScheduleForm({
    required this.id,
    required this.title,
    required this.description,
    required this.tasks,
    required this.selectedTaskIds,
    required this.saving,
    required this.onToggleTask,
    required this.onClear,
    required this.onSave,
  });

  final TextEditingController id;
  final TextEditingController title;
  final TextEditingController description;
  final List<LearningTask> tasks;
  final Set<String> selectedTaskIds;
  final bool saving;
  final void Function(String id, bool selected) onToggleTask;
  final VoidCallback onClear;
  final Future<void> Function() onSave;

  @override
  Widget build(BuildContext context) {
    return _Panel(
      child: Column(
        children: [
          Row(
            children: [
              Expanded(
                child: TextField(
                  controller: id,
                  decoration: const InputDecoration(labelText: '训练营 ID'),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: TextField(
                  controller: title,
                  decoration: const InputDecoration(labelText: '标题'),
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          TextField(
            controller: description,
            minLines: 2,
            maxLines: 4,
            decoration: const InputDecoration(labelText: '描述'),
          ),
          const SizedBox(height: 12),
          Align(
            alignment: Alignment.centerLeft,
            child: Wrap(
              spacing: 8,
              runSpacing: 8,
              children: [
                for (final task in tasks)
                  FilterChip(
                    label: Text(task.title),
                    selected: selectedTaskIds.contains(task.id),
                    onSelected: (selected) => onToggleTask(task.id, selected),
                  ),
              ],
            ),
          ),
          const SizedBox(height: 12),
          _Actions(saving: saving, onClear: onClear, onSave: onSave),
        ],
      ),
    );
  }
}

class _Actions extends StatelessWidget {
  const _Actions({
    required this.saving,
    required this.onClear,
    required this.onSave,
  });

  final bool saving;
  final VoidCallback onClear;
  final Future<void> Function() onSave;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.end,
      children: [
        OutlinedButton.icon(
          onPressed: saving ? null : onClear,
          icon: const Icon(Icons.clear),
          label: const Text('清空'),
        ),
        const SizedBox(width: 8),
        FilledButton.icon(
          onPressed: saving ? null : onSave,
          icon: const Icon(Icons.save_outlined),
          label: const Text('保存'),
        ),
      ],
    );
  }
}

class _Panel extends StatelessWidget {
  const _Panel({required this.child});

  final Widget child;

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
      decoration: BoxDecoration(
        border: Border.all(color: Theme.of(context).dividerColor),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Padding(padding: const EdgeInsets.all(16), child: child),
    );
  }
}
