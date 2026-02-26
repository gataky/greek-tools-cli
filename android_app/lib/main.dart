import 'package:flutter/material.dart';
import 'screens/setup_screen.dart';

void main() {
  runApp(const GreekCasePracticeApp());
}

/// Main app widget
class GreekCasePracticeApp extends StatelessWidget {
  const GreekCasePracticeApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Greek Case Practice',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(
          seedColor: Colors.blue,
          brightness: Brightness.light,
        ),
        useMaterial3: true,
        cardTheme: CardThemeData(
          elevation: 2,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
        ),
      ),
      home: const SetupScreen(),
    );
  }
}
