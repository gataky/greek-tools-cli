import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../blocs/setup/setup_bloc.dart';
import '../blocs/setup/setup_event.dart';
import '../blocs/setup/setup_state.dart';
import '../models/session_config.dart';
import '../widgets/difficulty_selector.dart';
import 'practice_screen.dart';

/// Setup screen for configuring practice session
class SetupScreen extends StatelessWidget {
  const SetupScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (context) => SetupBloc(),
      child: BlocListener<SetupBloc, SetupState>(
        listenWhen: (previous, current) => false, // We'll handle navigation differently
        listener: (context, state) {},
        child: const _SetupScreenContent(),
      ),
    );
  }
}

class _SetupScreenContent extends StatelessWidget {
  const _SetupScreenContent();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Greek Case Practice'),
        centerTitle: true,
      ),
      body: BlocBuilder<SetupBloc, SetupState>(
        builder: (context, state) {
          return SingleChildScrollView(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                // Difficulty selector
                DifficultySelector(
                  selectedDifficulty: state.difficultyLevel,
                  onChanged: (difficulty) {
                    context.read<SetupBloc>().add(DifficultyChanged(difficulty));
                  },
                ),
                const SizedBox(height: 24),

                // Plural checkbox
                CheckboxListTile(
                  title: const Text('Include plural forms'),
                  value: state.includePlural,
                  onChanged: (value) {
                    context.read<SetupBloc>().add(PluralToggled());
                  },
                ),
                const SizedBox(height: 24),

                // Question count selector
                const Text(
                  'Number of Questions',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 8),
                DropdownButtonFormField<int>(
                  value: state.questionCount,
                  decoration: const InputDecoration(
                    border: OutlineInputBorder(),
                    contentPadding: EdgeInsets.symmetric(
                      horizontal: 16,
                      vertical: 12,
                    ),
                  ),
                  items: const [
                    DropdownMenuItem(value: 10, child: Text('10 questions')),
                    DropdownMenuItem(value: 20, child: Text('20 questions')),
                    DropdownMenuItem(value: 50, child: Text('50 questions')),
                    DropdownMenuItem(value: 0, child: Text('Endless mode')),
                  ],
                  onChanged: (value) {
                    if (value != null) {
                      context
                          .read<SetupBloc>()
                          .add(QuestionCountChanged(value));
                    }
                  },
                ),
                const SizedBox(height: 32),

                // Start button
                ElevatedButton(
                  onPressed: () {
                    // Create session config from current state
                    final config = SessionConfig(
                      difficultyLevel: state.difficultyLevel,
                      includePlural: state.includePlural,
                      questionCount: state.questionCount,
                    );

                    // Navigate to practice screen
                    Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (context) => PracticeScreen(config: config),
                      ),
                    );
                  },
                  style: ElevatedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(vertical: 16),
                    textStyle: const TextStyle(fontSize: 18),
                  ),
                  child: const Text('Start Practice'),
                ),
              ],
            ),
          );
        },
      ),
    );
  }
}
