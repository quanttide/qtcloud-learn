import 'package:flutter/material.dart';

void main() {
  runApp(const QtcloudLearnApp());
}

class QtcloudLearnApp extends StatelessWidget {
  const QtcloudLearnApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: '量潮学习云',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.teal),
        useMaterial3: true,
      ),
      home: const HomePage(),
    );
  }
}

class HomePage extends StatelessWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('量潮学习云')),
      body: const Center(
        child: Text('AI 原生学员学习中心'),
      ),
    );
  }
}
