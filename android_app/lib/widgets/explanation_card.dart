import 'package:flutter/material.dart';
import '../models/sentence.dart';
import '../utils/explanation_generator.dart';

/// Widget for displaying grammar explanation after answer
class ExplanationCard extends StatelessWidget {
  final Sentence sentence;

  const ExplanationCard({
    super.key,
    required this.sentence,
  });

  @override
  Widget build(BuildContext context) {
    final translation =
        ExplanationGenerator.generateTranslation(sentence.englishPrompt);
    final syntacticRole = ExplanationGenerator.generateSyntacticRole(
      sentence.contextType,
      sentence.caseType,
      sentence.preposition,
    );
    final morphology = ExplanationGenerator.generateMorphology(sentence);

    return Card(
      elevation: 2,
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      color: Colors.blue.shade50,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Header
            Row(
              children: [
                Icon(Icons.lightbulb_outline, color: Colors.blue.shade700),
                const SizedBox(width: 8),
                Text(
                  'Explanation',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: Colors.blue.shade900,
                  ),
                ),
              ],
            ),
            const Divider(height: 24),

            // Translation
            _buildSection('Translation', translation),
            const SizedBox(height: 12),

            // Syntactic Role
            _buildSection('Syntactic Role', syntacticRole),
            const SizedBox(height: 12),

            // Morphology
            _buildSection('Morphology', morphology),
          ],
        ),
      ),
    );
  }

  Widget _buildSection(String title, String content) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          title,
          style: const TextStyle(
            fontSize: 14,
            fontWeight: FontWeight.w600,
            color: Colors.black87,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          content,
          style: const TextStyle(
            fontSize: 14,
            color: Colors.black87,
          ),
        ),
      ],
    );
  }
}
