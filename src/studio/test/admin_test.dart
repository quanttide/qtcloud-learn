// LMS 后台测试：AdminApi 解析 + AdminShell 导航 + 立项姓名栏展示（注入内存替身，不发网络请求）。

import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:http/testing.dart';

import 'package:qtcloud_learn_studio/main_admin.dart';
import 'package:qtcloud_learn_studio/admin/admin_api.dart';
import 'package:qtcloud_learn_studio/models/application.dart';

/// 内存 Mock 客户端：按路径返回预设 JSON。
MockClient _mockClient() {
  return MockClient((request) async {
    final path = request.url.path;
    if (path.endsWith('/students')) {
      return http.Response(
        jsonEncode({
          'students': [
            {'id': 'stu-1', 'name': '张三', 'email': 'zhangsan@example.com'},
          ],
        }),
        200,
        headers: {'content-type': 'application/json'},
      );
    }
    if (path.endsWith('/progress')) {
      return http.Response(
        jsonEncode({
          'progress': [
            {'id': 'prog-1', 'studentId': 'stu-1', 'classId': 'class-1', 'percent': 0.5, 'finished': false, 'updatedAt': '2026-08-16T10:00:00Z'},
          ],
        }),
        200,
        headers: {'content-type': 'application/json'},
      );
    }
    if (path.endsWith('/applications')) {
      return http.Response(
        jsonEncode({
          'applications': [
            {
              'id': 'appl-1',
              'projectName': '校园选课助手',
              'blindSpot': '选课信息分散',
              'demoPlan': '聚合查询页',
              'direction': '内容',
              'teamMode': 'partner',
              'memberNames': ['张三', '李四'],
              'studentId': 'stu-1',
              'studentName': '张三',
              'status': 'submitted',
              'createdAt': '2026-08-16T10:00:00Z',
            },
          ],
        }),
        200,
        headers: {'content-type': 'application/json'},
      );
    }
    return http.Response('{"error":"not found"}', 404);
  });
}

void main() {
  test('AdminApi 解析学员/进度/立项', () async {
    final api = AdminApi(client: _mockClient(), baseUrl: 'http://fake/api/v1');

    final students = await api.fetchStudents();
    expect(students.length, 1);
    expect(students.first.name, '张三');

    final progress = await api.fetchProgress();
    expect(progress.length, 1);
    expect(progress.first.percent, 0.5);

    final applications = await api.fetchApplications();
    expect(applications.length, 1);
    expect(applications.first.projectName, '校园选课助手');
  });

  test('Application.memberLabel：个人独立 / 搭档队长+队员', () {
    const personal = Application(
      id: 'a1',
      projectName: 'p',
      blindSpot: 'b',
      demoPlan: 'd',
      direction: '',
      teamMode: 'personal',
      memberNames: ['张三'],
      studentId: 's1',
      studentName: '张三',
      status: 'submitted',
    );
    expect(personal.memberLabel, '张三');

    const partner = Application(
      id: 'a2',
      projectName: 'p',
      blindSpot: 'b',
      demoPlan: 'd',
      direction: '',
      teamMode: 'partner',
      memberNames: ['李四', '王五'],
      studentId: 's2',
      studentName: '李四',
      status: 'submitted',
    );
    expect(partner.memberLabel, '队长：李四 / 队员：王五');
  });

  testWidgets('AdminShell 侧边栏四个板块 + 立项页展示姓名栏', (WidgetTester tester) async {
    final api = AdminApi(client: _mockClient(), baseUrl: 'http://fake/api/v1');
    await tester.pumpWidget(AdminApp(api: api));
    await tester.pumpAndSettle();

    // 侧边栏板块（页头标题同名，故 findsWidgets）
    expect(find.text('概览'), findsWidgets);
    expect(find.text('学员'), findsWidgets);
    expect(find.text('进度'), findsWidgets);
    expect(find.text('立项'), findsWidgets);

    // 切到立项页：项目名 + 姓名栏（队长/队员）
    await tester.tap(find.text('立项'));
    await tester.pumpAndSettle();
    expect(find.text('校园选课助手'), findsOneWidget);
    expect(find.text('队长：张三 / 队员：李四'), findsOneWidget);
  });
}
