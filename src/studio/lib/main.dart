import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import 'models/student.dart';
import 'screens/class_screen.dart';
import 'screens/dashboard_screen.dart';
import 'screens/login_screen.dart';
import 'screens/my_courses_screen.dart';
import 'screens/schedule_screen.dart';
import 'services/assessment_service.dart';
import 'services/learn_data_service.dart';

void main() {
  runApp(const QtcloudLearnApp());
}

class QtcloudLearnApp extends StatelessWidget {
  final LearnDataService? dataService;
  final AssessmentService? assessmentService;

  const QtcloudLearnApp({
    super.key,
    this.dataService,
    this.assessmentService,
  });

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        dataService != null
            ? ChangeNotifierProvider<LearnDataService>.value(value: dataService!)
            : ChangeNotifierProvider<LearnDataService>(
                create: (_) => LearnDataService()..load()),
        assessmentService != null
            ? ChangeNotifierProvider<AssessmentService>.value(
                value: assessmentService!)
            : ChangeNotifierProvider<AssessmentService>(
                create: (_) => AssessmentService()..load()),
      ],
      child: MaterialApp(
        title: '量潮学习云',
        theme: ThemeData(
          colorScheme: ColorScheme.fromSeed(seedColor: Colors.teal),
          useMaterial3: true,
        ),
        home: const RootShell(),
      ),
    );
  }
}

/// 登录门控：未登录时展示登录页（本地模拟），登录后进入主壳。
class RootShell extends StatefulWidget {
  const RootShell({super.key});

  @override
  State<RootShell> createState() => _RootShellState();
}

class _RootShellState extends State<RootShell> {
  Student? _student;

  @override
  Widget build(BuildContext context) {
    final student = _student;
    if (student == null) {
      return LoginScreen(
        onLogin: (s) => setState(() => _student = s),
      );
    }
    return HomeShell(
      student: student,
      onLogout: () => setState(() => _student = null),
    );
  }
}

class HomeShell extends StatefulWidget {
  final Student student;
  final VoidCallback onLogout;

  const HomeShell({
    super.key,
    required this.student,
    required this.onLogout,
  });

  @override
  State<HomeShell> createState() => _HomeShellState();
}

class _HomeShellState extends State<HomeShell> {
  int _index = 0;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: IndexedStack(
        index: _index,
        children: [
          DashboardScreen(student: widget.student, onLogout: widget.onLogout),
          MyCoursesScreen(student: widget.student),
          const ScheduleScreen(),
          const ClassScreen(),
        ],
      ),
      bottomNavigationBar: NavigationBar(
        selectedIndex: _index,
        onDestinationSelected: (i) => setState(() => _index = i),
        destinations: const [
          NavigationDestination(
              icon: Icon(Icons.dashboard_outlined),
              selectedIcon: Icon(Icons.dashboard),
              label: '仪表盘'),
          NavigationDestination(
              icon: Icon(Icons.book_outlined),
              selectedIcon: Icon(Icons.book),
              label: '我的课程'),
          NavigationDestination(
              icon: Icon(Icons.calendar_month_outlined),
              selectedIcon: Icon(Icons.calendar_month),
              label: '课表'),
          NavigationDestination(
              icon: Icon(Icons.groups_outlined),
              selectedIcon: Icon(Icons.groups),
              label: '班级'),
        ],
      ),
    );
  }
}
