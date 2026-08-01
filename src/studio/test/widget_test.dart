import 'package:flutter_test/flutter_test.dart';

import 'package:qtcloud_learn_studio/main.dart';

void main() {
  testWidgets('renders home page', (WidgetTester tester) async {
    await tester.pumpWidget(const QtcloudLearnApp());

    expect(find.text('量潮学习云'), findsWidgets);
    expect(find.text('AI 原生学员学习中心'), findsOneWidget);
  });
}
