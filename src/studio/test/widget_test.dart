import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:provider/provider.dart';

import 'package:qtcloud_learn_studio/main.dart';
import 'package:qtcloud_learn_studio/models/enums.dart';
import 'package:qtcloud_learn_studio/models/student.dart';
import 'package:qtcloud_learn_studio/screens/class_screen.dart';
import 'package:qtcloud_learn_studio/screens/dashboard_screen.dart';
import 'package:qtcloud_learn_studio/screens/my_courses_screen.dart';
import 'package:qtcloud_learn_studio/screens/schedule_screen.dart';
import 'package:qtcloud_learn_studio/services/assessment_service.dart';
import 'package:qtcloud_learn_studio/services/learn_data_service.dart';

/// 通过真实事件循环完成 assets 加载（widget 测试的 fake-async 区无法完成 rootBundle 异步）。
Future<LearnDataService> loadData(WidgetTester tester) async {
  final service = LearnDataService();
  await tester.runAsync(() => service.load());
  return service;
}

Future<AssessmentService> loadAssessment(WidgetTester tester) async {
  final service = AssessmentService();
  await tester.runAsync(() => service.load());
  return service;
}

/// 测试用 provider 外壳（注入预加载 service，不触发异步 load）。
Widget harness({
  required Widget child,
  LearnDataService? data,
  AssessmentService? assessment,
}) {
  return MultiProvider(
    providers: [
      ChangeNotifierProvider<LearnDataService>.value(
          value: data ?? LearnDataService()),
      ChangeNotifierProvider<AssessmentService>.value(
          value: assessment ?? AssessmentService()),
    ],
    child: MaterialApp(home: child),
  );
}

const demoStudent = Student(
  id: 'student-1',
  name: '张三',
  email: 'zhangsan@example.com',
);

void main() {
  group('服务层', () {
    test('从 assets 加载统一数据', () async {
      final service = LearnDataService();
      await service.load();
      expect(service.classes.length, 2);
      expect(service.students.length, 5);
      expect(service.teachers.length, 2);
      expect(service.sessions.length, 3);
      expect(service.enrollments.length, 3);
      expect(service.progress.length, 3);
      expect(service.upcomingSessions.length, 1);
      expect(service.inProgressSessions.length, 1);
      expect(service.completedSessions.length, 1);
    });

    test('选课 / 退课 / 进度上报', () async {
      final service = LearnDataService();
      await service.load();
      expect(service.isEnrolled('class-2', 'student-1'), isFalse);
      service.enroll('class-2', 'student-1');
      expect(service.isEnrolled('class-2', 'student-1'), isTrue);
      expect(service.getMyClasses('student-1').length, 2);

      service.reportProgress('class-1', 'student-1',
          percent: 0.9, finished: true);
      final p = service.getProgress('class-1', 'student-1');
      expect(p, isNotNull);
      expect(p!.percent, 0.9);
      expect(p.finished, isTrue);

      service.withdraw('class-2', 'student-1');
      expect(service.isEnrolled('class-2', 'student-1'), isFalse);
    });

    test('考勤标记', () async {
      final service = LearnDataService();
      await service.load();
      service.markAllAttendance('sess-1', AttendanceStatus.present);
      final session = service.sessions.firstWhere((s) => s.id == 'sess-1');
      expect(
        session.attendances.every((a) => a.status == AttendanceStatus.present),
        isTrue,
      );
    });
  });

  group('登录与主壳', () {
    testWidgets('未登录显示登录页（本地模拟）', (tester) async {
      final data = await loadData(tester);
      await tester.pumpWidget(QtcloudLearnApp(dataService: data));
      await tester.pumpAndSettle();
      expect(find.text('飞书登录（本地模拟）'), findsOneWidget);
    });

    testWidgets('登录后进入仪表盘', (tester) async {
      final data = await loadData(tester);
      final assessment = await loadAssessment(tester);
      await tester.pumpWidget(
          QtcloudLearnApp(dataService: data, assessmentService: assessment));
      await tester.pumpAndSettle();
      await tester.tap(find.text('飞书登录（本地模拟）'));
      await tester.pumpAndSettle();
      expect(find.textContaining('你好，张三'), findsOneWidget);
      expect(find.text('我的课程'), findsWidgets);
    });

    testWidgets('底部导航切换课表', (tester) async {
      final data = await loadData(tester);
      final assessment = await loadAssessment(tester);
      await tester.pumpWidget(
          QtcloudLearnApp(dataService: data, assessmentService: assessment));
      await tester.pumpAndSettle();
      await tester.tap(find.text('飞书登录（本地模拟）'));
      await tester.pumpAndSettle();
      await tester.tap(find.descendant(
        of: find.byType(NavigationBar),
        matching: find.text('课表'),
      ));
      await tester.pumpAndSettle();
      expect(find.text('Python 基础'), findsOneWidget);
    });
  });

  group('页面渲染', () {
    testWidgets('仪表盘显示学员指标', (tester) async {
      final data = await loadData(tester);
      await tester.pumpWidget(harness(
        data: data,
        child: DashboardScreen(student: demoStudent, onLogout: () {}),
      ));
      await tester.pumpAndSettle();
      expect(find.textContaining('你好，张三'), findsOneWidget);
      expect(find.text('我的课程'), findsWidgets);
    });

    testWidgets('我的课程页显示已选课程与报名区', (tester) async {
      final data = await loadData(tester);
      await tester.pumpWidget(harness(
        data: data,
        child: MyCoursesScreen(student: demoStudent),
      ));
      await tester.pumpAndSettle();
      expect(find.text('已选课程 (1)'), findsOneWidget);
      expect(find.text('选课报名'), findsOneWidget);
    });

    testWidgets('课表页按状态分组', (tester) async {
      final data = await loadData(tester);
      await tester.pumpWidget(harness(
        data: data,
        child: const ScheduleScreen(),
      ));
      await tester.pumpAndSettle();
      expect(find.text('正在上课'), findsOneWidget);
      expect(find.text('即将开始'), findsOneWidget);
      expect(find.text('历史课程'), findsOneWidget);
    });

    testWidgets('班级管理页显示班级列表', (tester) async {
      final data = await loadData(tester);
      await tester.pumpWidget(harness(
        data: data,
        child: const ClassScreen(),
      ));
      await tester.pumpAndSettle();
      expect(find.text('浙理班级'), findsOneWidget);
      expect(find.text('杭电班级'), findsOneWidget);
    });

    testWidgets('班级详情显示教师与学员', (tester) async {
      final data = await loadData(tester);
      final c = data.classes.first;
      await tester.pumpWidget(harness(
        data: data,
        child: ClassDetailScreen(classModel: c),
      ));
      await tester.pumpAndSettle();
      expect(find.text('授课教师'), findsOneWidget);
      expect(find.text('学员列表 (3)'), findsOneWidget);
      await tester.scrollUntilVisible(find.text('考核管理'), 200,
          scrollable: find.byType(Scrollable).first);
      expect(find.text('考核管理'), findsOneWidget);
    });
  });
}
