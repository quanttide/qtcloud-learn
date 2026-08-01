import 'package:flutter/material.dart';

import '../models/enums.dart';

/// 班级状态标签。
class StatusChip extends StatelessWidget {
  final ClassStatus status;

  const StatusChip({super.key, required this.status});

  @override
  Widget build(BuildContext context) {
    final color = switch (status) {
      ClassStatus.preparing => Colors.orange,
      ClassStatus.active => Colors.green,
      ClassStatus.ended => Colors.grey,
    };
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Text(
        status.label,
        style: TextStyle(fontSize: 11, color: color),
      ),
    );
  }
}
