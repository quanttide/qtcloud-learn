// LMS 管理后台框架页（对齐规格 admin.md）：AppBar + 侧边栏 + 内容区。
// v0.1 菜单收窄为后台核心三块：概览 / 学员 / 立项。

import 'package:flutter/material.dart';

import 'admin_api.dart';
import 'pages/dashboard_page.dart';
import 'pages/learners_page.dart';
import 'pages/proposals_page.dart';

class AdminShell extends StatefulWidget {
  const AdminShell({super.key, required this.api});

  final AdminApi api;

  @override
  State<AdminShell> createState() => _AdminShellState();
}

class _AdminShellState extends State<AdminShell> {
  int _index = 0;

  @override
  Widget build(BuildContext context) {
    final pages = [
      DashboardPage(api: widget.api),
      LearnersPage(api: widget.api),
      ProposalsPage(api: widget.api),
    ];
    return Scaffold(
      appBar: AppBar(
        title: const Text('量潮学习云 · LMS 后台'),
        actions: [
          Padding(
            padding: const EdgeInsets.only(right: 16),
            child: Center(child: Text('v0.1', style: Theme.of(context).textTheme.bodySmall)),
          ),
        ],
      ),
      body: Row(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          SizedBox(
            width: 200,
            child: NavigationRail(
              selectedIndex: _index,
              onDestinationSelected: (i) => setState(() => _index = i),
              labelType: NavigationRailLabelType.all,
              destinations: const [
                NavigationRailDestination(
                  icon: Icon(Icons.dashboard_outlined),
                  selectedIcon: Icon(Icons.dashboard),
                  label: Text('概览'),
                ),
                NavigationRailDestination(
                  icon: Icon(Icons.person_outline),
                  selectedIcon: Icon(Icons.person),
                  label: Text('学员'),
                ),
                NavigationRailDestination(
                  icon: Icon(Icons.rocket_launch_outlined),
                  selectedIcon: Icon(Icons.rocket_launch),
                  label: Text('立项'),
                ),
              ],
            ),
          ),
          const VerticalDivider(width: 1),
          Expanded(child: pages[_index]),
        ],
      ),
    );
  }
}
