// LMS 管理后台入口（独立构建目标：flutter build web --target=lib/main_admin.dart）。
// 与学员端（main.dart）分离部署：前台后台各一个 HTML 文件。

import 'package:flutter/material.dart';

import 'admin/admin_api.dart';
import 'admin/admin_shell.dart';

void main() {
  runApp(const AdminApp());
}

class AdminApp extends StatelessWidget {
  const AdminApp({super.key, this.api});

  /// 测试注入替身；生产按 QTCLOUD_LEARN_API_URL 创建。
  final AdminApi? api;

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: '量潮学习云 · LMS 后台',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.teal),
        useMaterial3: true,
      ),
      home: AdminShell(api: api ?? AdminApi()),
    );
  }
}
