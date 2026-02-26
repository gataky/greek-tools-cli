import 'package:flutter/material.dart';

/// Widget for selecting difficulty level with radio buttons
class DifficultySelector extends StatelessWidget {
  final String selectedDifficulty;
  final Function(String) onChanged;

  const DifficultySelector({
    super.key,
    required this.selectedDifficulty,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Difficulty Level',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 8),
        RadioListTile<String>(
          title: const Text('Beginner'),
          value: 'beginner',
          groupValue: selectedDifficulty,
          onChanged: (value) => onChanged(value!),
        ),
        RadioListTile<String>(
          title: const Text('Intermediate'),
          value: 'intermediate',
          groupValue: selectedDifficulty,
          onChanged: (value) => onChanged(value!),
        ),
        RadioListTile<String>(
          title: const Text('Advanced'),
          value: 'advanced',
          groupValue: selectedDifficulty,
          onChanged: (value) => onChanged(value!),
        ),
      ],
    );
  }
}
