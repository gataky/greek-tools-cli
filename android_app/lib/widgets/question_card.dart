import 'package:flutter/material.dart';
import '../models/sentence.dart';

/// Widget for displaying a practice question
class QuestionCard extends StatelessWidget {
  final Sentence sentence;
  final bool showAnswer;

  const QuestionCard({
    super.key,
    required this.sentence,
    this.showAnswer = false,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 4,
      margin: const EdgeInsets.all(16),
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            // English prompt
            Text(
              showAnswer
                  ? _getFullEnglishPrompt()
                  : sentence.englishPrompt,
              style: const TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.w500,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),

            // Greek sentence
            Text(
              showAnswer ? _getCompleteSentence() : _getSentenceWithBlanks(),
              style: const TextStyle(
                fontSize: 24,
                fontWeight: FontWeight.bold,
                letterSpacing: 0.5,
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }

  /// Get sentence with blanks (for question state)
  String _getSentenceWithBlanks() {
    // Replace the correct answer with blanks
    return sentence.greekSentence
        .replaceAll(sentence.correctAnswer, '___ ___');
  }

  /// Get complete sentence with answer (for answer revealed state)
  String _getCompleteSentence() {
    return sentence.greekSentence;
  }

  /// Get full English prompt without blanks (for answer revealed state)
  String _getFullEnglishPrompt() {
    // Remove the blank marker and just show clean English
    return sentence.englishPrompt
        .replaceAll('___', '')
        .replaceAll(RegExp(r'\s*\([^)]*\)'), '')
        .trim();
  }
}
