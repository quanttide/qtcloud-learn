import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../models/class_model.dart';
import '../models/student.dart';
import '../services/learn_data_service.dart';

/// 学习进度页 —— 展示并上报单个班级的学习进度。
/// 自 `qtclass` learning_record（localStorage）迁移：改为服务端进度数据。
class ProgressScreen extends StatefulWidget {
  final ClassModel classModel;
  final Student student;

  const ProgressScreen({
    super.key,
    required this.classModel,
    required this.student,
  });

  @override
  State<ProgressScreen> createState() => _ProgressScreenState();
}

class _ProgressScreenState extends State<ProgressScreen> {
  double _percent = 0.0;
  bool _finished = false;

  @override
  void initState() {
    super.initState();
    final service = context.read<LearnDataService>();
    final p = service.getProgress(widget.classModel.id, widget.student.id);
    _percent = p?.percent ?? widget.classModel.progress;
    _finished = p?.finished ?? false;
  }

  @override
  Widget build(BuildContext context) {
    final service = context.watch<LearnDataService>();
    final progress = service.getProgress(widget.classModel.id, widget.student.id);
    final percent = progress?.percent ?? _percent;
    final finished = progress?.finished ?? _finished;

    return Scaffold(
      appBar: AppBar(title: Text(widget.classModel.name)),
      body: ListView(
        padding: const EdgeInsets.all(24),
        children: [
          Card(
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(widget.classModel.refName,
                      style: Theme.of(context)
                          .textTheme
                          .titleMedium
                          ?.copyWith(fontWeight: FontWeight.bold)),
                  const SizedBox(height: 4),
                  Text('${widget.classModel.startDate} - ${widget.classModel.endDate}',
                      style: TextStyle(color: Colors.grey[600])),
                  const SizedBox(height: 16),
                  LinearProgressIndicator(value: percent, minHeight: 8),
                  const SizedBox(height: 8),
                  Text('${(percent * 100).toInt()}%',
                      style: const TextStyle(
                          fontSize: 32, fontWeight: FontWeight.bold)),
                  const SizedBox(height: 4),
                  if (finished)
                    Text('🎉 已完成',
                        style: TextStyle(color: Colors.green[700])),
                  if (progress?.updatedAt != null)
                    Text('最近更新：${progress!.updatedAt!.substring(0, 10)}',
                        style:
                            TextStyle(color: Colors.grey[500], fontSize: 12)),
                ],
              ),
            ),
          ),
          const SizedBox(height: 24),
          Card(
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('上报学习进度',
                      style:
                          TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                  const SizedBox(height: 8),
                  Text('原 `qtclass` history_service 的学习记录已迁移为服务端进度数据。',
                      style: TextStyle(color: Colors.grey[600], fontSize: 12)),
                  const SizedBox(height: 16),
                  Slider(
                    value: _percent,
                    min: 0,
                    max: 1,
                    divisions: 10,
                    label: '${(_percent * 100).toInt()}%',
                    onChanged: (v) => setState(() => _percent = v),
                  ),
                  Row(
                    children: [
                      Expanded(
                        child: Text('${(_percent * 100).toInt()}%',
                            style: const TextStyle(fontSize: 20)),
                      ),
                      CheckboxListTile(
                        contentPadding: EdgeInsets.zero,
                        title: const Text('已完成'),
                        value: _finished,
                        onChanged: (v) => setState(() => _finished = v ?? false),
                        controlAffinity: ListTileControlAffinity.leading,
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  SizedBox(
                    width: double.infinity,
                    child: FilledButton.icon(
                      onPressed: () {
                        service.reportProgress(
                          widget.classModel.id,
                          widget.student.id,
                          percent: _percent,
                          finished: _finished,
                        );
                        ScaffoldMessenger.of(context).showSnackBar(
                          const SnackBar(content: Text('进度已上报')),
                        );
                      },
                      icon: const Icon(Icons.upload),
                      label: const Text('上报进度'),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
