import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../models/student.dart';
import '../services/learn_data_service.dart';

/// 登录页 —— 飞书登录（本地模拟）。
///
/// 认证 API 未排期（见根 ROADMAP「版本依赖」），v0.3 先本地模拟：
/// 默认使用第一个学员作为当前用户登录。
class LoginScreen extends StatelessWidget {
  final void Function(Student student) onLogin;

  const LoginScreen({super.key, required this.onLogin});

  @override
  Widget build(BuildContext context) {
    final service = context.watch<LearnDataService>();
    final students = service.students;

    return Scaffold(
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(32),
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 420),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.school, size: 72, color: Colors.teal),
                  const SizedBox(height: 16),
                  Text('量潮学习云',
                      style: Theme.of(context)
                          .textTheme
                          .headlineMedium
                          ?.copyWith(fontWeight: FontWeight.bold)),
                  const SizedBox(height: 8),
                  const Text('AI 原生学员学习中心',
                      style: TextStyle(color: Colors.grey)),
                  const SizedBox(height: 48),
                  if (students.isNotEmpty)
                    DropdownButtonFormField<String>(
                      initialValue: students.first.id,
                      decoration: const InputDecoration(
                        labelText: '学员',
                        border: OutlineInputBorder(),
                      ),
                      items: students
                          .map((s) => DropdownMenuItem(
                                value: s.id,
                                child: Text('${s.name}（${s.plan.label}）'),
                              ))
                          .toList(),
                      onChanged: (_) {},
                    ),
                  const SizedBox(height: 24),
                  SizedBox(
                    width: double.infinity,
                    height: 48,
                    child: FilledButton.icon(
                      onPressed: students.isEmpty
                          ? null
                          : () {
                              final selectedId = students.first.id;
                              onLogin(
                                  students.firstWhere((s) => s.id == selectedId));
                            },
                      icon: const Icon(Icons.login),
                      label: const Text('飞书登录（本地模拟）'),
                    ),
                  ),
                  const SizedBox(height: 12),
                  const Text('认证 API 未排期，当前为本地模拟登录',
                      style: TextStyle(fontSize: 12, color: Colors.grey)),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
