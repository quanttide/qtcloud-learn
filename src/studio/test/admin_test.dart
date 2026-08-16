// LMS 后台测试：AdminApi 解析学员/立项 + 后台表格展示（注入内存替身，不发网络请求）。

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
    if (path.endsWith('/api/learners')) {
      return http.Response(
        jsonEncode({
          'learners': [
            {
              'id': 'lea-1',
              'name': '张三',
              'course': '生产实习',
              'progressMax': 3,
              'progressTotal': 5,
              'activeAt': '2026-08-16T10:00:00Z',
              'status': '在读',
              'projectName': '选课助手',
            },
          ],
        }),
        200,
        headers: {'content-type': 'application/json'},
      );
    }
    if (path.endsWith('/api/proposals/history')) {
      return http.Response(jsonEncode({'history': []}), 200,
          headers: {'content-type': 'application/json'});
    }
    if (path.contains('/api/proposals/')) {
      return http.Response('{"success":true}', 200,
          headers: {'content-type': 'application/json'});
    }
    if (path.endsWith('/api/proposals')) {
      return http.Response(
        jsonEncode({
          'proposals': [
            {
              'id': 'appl-1',
              'projectName': '选课助手',
              'opportunity': '选课信息分散',
              'fit': '量潮有数据能力',
              'hypothesis': '聚合查询提升效率',
              'demo': '查询页原型',
              'directionType': '内容',
              'teamMode': 'partner',
              'teamLeader': '张三',
              'teamMember': '李四、王五',
              'studentName': '张三',
              'status': '已提交',
              'submittedAt': '2026-08-16T10:00:00Z',
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
  test('AdminApi 解析学员档案/立项列表/历史', () async {
    final api = AdminApi(client: _mockClient(), baseUrl: 'http://fake');

    final learners = await api.fetchLearners();
    expect(learners.length, 1);
    expect(learners.first.name, '张三');
    expect(learners.first.progressLabel, '3/5');
    expect(learners.first.projectName, '选课助手');

    final proposals = await api.fetchProposals();
    expect(proposals.length, 1);
    expect(proposals.first.projectName, '选课助手');
    expect(proposals.first.teamLabel, '队长：张三 / 队员：李四、王五');

    final history = await api.fetchHistory();
    expect(history, isEmpty);

    await api.deleteProposal('appl-1'); // 不抛异常即通过
  });

  test('Application.teamLabel：个人独立 / 搭档队长+队员', () {
    const personal = Application(
      id: 'a1',
      projectName: 'p',
      opportunity: '',
      fit: '',
      hypothesis: '',
      demo: '',
      directionType: '',
      teamMode: 'personal',
      teamLeader: '张三',
      teamMember: '',
      studentName: '张三',
      status: '已提交',
    );
    expect(personal.teamLabel, '张三');

    const partner = Application(
      id: 'a2',
      projectName: 'p',
      opportunity: '',
      fit: '',
      hypothesis: '',
      demo: '',
      directionType: '',
      teamMode: 'partner',
      teamLeader: '李四',
      teamMember: '王五、赵六',
      studentName: '李四',
      status: '已提交',
    );
    expect(partner.teamLabel, '队长：李四 / 队员：王五、赵六');
  });

  testWidgets('AdminShell 侧边栏三板块 + 学员表进度/立项 + 立项表组队', (WidgetTester tester) async {
    final api = AdminApi(client: _mockClient(), baseUrl: 'http://fake');
    await tester.pumpWidget(AdminApp(api: api));
    await tester.pumpAndSettle();

    // 侧边栏
    expect(find.text('概览'), findsWidgets);
    expect(find.text('学员'), findsWidgets);
    expect(find.text('立项'), findsWidgets);

    // 切到学员 tab：姓名/进度/立项✓
    await tester.tap(find.descendant(
      of: find.byType(NavigationRail),
      matching: find.text('学员'),
    ));
    await tester.pumpAndSettle();
    expect(find.text('张三'), findsOneWidget);
    expect(find.text('3/5 模块'), findsOneWidget);
    expect(find.text('✓ 选课助手'), findsOneWidget);

    // 立项表：组队姓名栏 + 历史开关
    await tester.tap(find.descendant(
      of: find.byType(NavigationRail),
      matching: find.text('立项'),
    ));
    await tester.pumpAndSettle();
    expect(find.text('选课助手'), findsOneWidget);
    expect(find.text('队长：张三 / 队员：李四、王五'), findsOneWidget);
    expect(find.textContaining('历史记录'), findsOneWidget);
  });
}
