import 'package:flutter/material.dart';

/// Results screen showing final score after practice session
class ResultsScreen extends StatelessWidget {
  final int totalQuestions;
  final int correctCount;
  final int incorrectCount;

  const ResultsScreen({
    super.key,
    required this.totalQuestions,
    required this.correctCount,
    required this.incorrectCount,
  });

  @override
  Widget build(BuildContext context) {
    final percentage = totalQuestions > 0
        ? (correctCount / totalQuestions * 100).toStringAsFixed(0)
        : '0';

    return Scaffold(
      appBar: AppBar(
        title: const Text('Practice Complete!'),
        centerTitle: true,
        automaticallyImplyLeading: false, // Remove back button
      ),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(24.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              // Celebration icon
              Icon(
                Icons.emoji_events,
                size: 80,
                color: Colors.amber.shade600,
              ),
              const SizedBox(height: 32),

              // Total questions
              Text(
                'Total Questions: $totalQuestions',
                style: const TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.w500,
                ),
              ),
              const SizedBox(height: 16),

              // Correct count
              Text(
                'Correct: $correctCount ($percentage%)',
                style: TextStyle(
                  fontSize: 22,
                  fontWeight: FontWeight.bold,
                  color: Colors.green.shade700,
                ),
              ),
              const SizedBox(height: 8),

              // Incorrect count
              Text(
                'Incorrect: $incorrectCount',
                style: TextStyle(
                  fontSize: 18,
                  color: Colors.red.shade700,
                ),
              ),
              const SizedBox(height: 48),

              // Back to Setup button
              ElevatedButton(
                onPressed: () {
                  // Pop all routes and return to setup
                  Navigator.of(context).popUntil((route) => route.isFirst);
                },
                style: ElevatedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 48,
                    vertical: 16,
                  ),
                  textStyle: const TextStyle(fontSize: 18),
                ),
                child: const Text('Back to Setup'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
