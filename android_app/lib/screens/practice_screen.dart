import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../blocs/practice/practice_bloc.dart';
import '../blocs/practice/practice_event.dart';
import '../blocs/practice/practice_state.dart';
import '../data/database_helper.dart';
import '../data/repository.dart';
import '../models/session_config.dart';
import '../widgets/question_card.dart';
import '../widgets/explanation_card.dart';
import 'results_screen.dart';

/// Practice screen for answering questions
class PracticeScreen extends StatelessWidget {
  final SessionConfig config;

  const PracticeScreen({
    super.key,
    required this.config,
  });

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (context) {
        final repository = DatabaseRepository(DatabaseHelper());
        final bloc = PracticeBloc(repository);
        bloc.add(LoadSession(config));
        return bloc;
      },
      child: BlocConsumer<PracticeBloc, PracticeState>(
        listener: (context, state) {
          // Navigate to results when practice complete
          if (state is PracticeComplete) {
            Navigator.pushReplacement(
              context,
              MaterialPageRoute(
                builder: (context) => ResultsScreen(
                  totalQuestions: state.totalQuestions,
                  correctCount: state.correctCount,
                  incorrectCount: state.incorrectCount,
                ),
              ),
            );
          }
        },
        builder: (context, state) {
          return Scaffold(
            appBar: AppBar(
              title: _buildTitle(state),
              centerTitle: true,
            ),
            body: _buildBody(context, state),
          );
        },
      ),
    );
  }

  /// Build app bar title based on state
  Widget _buildTitle(PracticeState state) {
    if (state is QuestionState ||
        state is AnswerRevealedState ||
        state is ExplanationState) {
      final currentIndex = _getCurrentIndex(state);
      final totalQuestions = _getTotalQuestions(state);

      if (totalQuestions == 0) {
        // Endless mode
        return Text('Question $currentIndex');
      } else {
        return Text('Question $currentIndex of $totalQuestions');
      }
    }

    return const Text('Practice');
  }

  /// Build body based on state
  Widget _buildBody(BuildContext context, PracticeState state) {
    if (state is PracticeLoading) {
      return const Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            CircularProgressIndicator(),
            SizedBox(height: 16),
            Text('Generating questions...'),
          ],
        ),
      );
    }

    if (state is PracticeError) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(24.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(Icons.error_outline, size: 64, color: Colors.red),
              const SizedBox(height: 16),
              Text(
                'Error: ${state.message}',
                textAlign: TextAlign.center,
                style: const TextStyle(fontSize: 16),
              ),
              const SizedBox(height: 24),
              ElevatedButton(
                onPressed: () => Navigator.pop(context),
                child: const Text('Back to Setup'),
              ),
            ],
          ),
        ),
      );
    }

    if (state is QuestionState) {
      return _buildQuestionView(context, state);
    }

    if (state is AnswerRevealedState) {
      return _buildAnswerRevealedView(context, state);
    }

    if (state is ExplanationState) {
      return _buildExplanationView(context, state);
    }

    return const Center(child: Text('Unknown state'));
  }

  /// Build view for question state (answer hidden)
  Widget _buildQuestionView(BuildContext context, QuestionState state) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        QuestionCard(
          sentence: state.currentSentence,
          showAnswer: false,
        ),
        const SizedBox(height: 24),
        ElevatedButton(
          onPressed: () {
            context.read<PracticeBloc>().add(ShowAnswer());
          },
          style: ElevatedButton.styleFrom(
            padding: const EdgeInsets.symmetric(horizontal: 48, vertical: 16),
            textStyle: const TextStyle(fontSize: 18),
          ),
          child: const Text('Show Answer'),
        ),
      ],
    );
  }

  /// Build view for answer revealed state (waiting for correct/incorrect)
  Widget _buildAnswerRevealedView(
      BuildContext context, AnswerRevealedState state) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        QuestionCard(
          sentence: state.currentSentence,
          showAnswer: true,
        ),
        const SizedBox(height: 24),
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            ElevatedButton.icon(
              onPressed: () {
                context.read<PracticeBloc>().add(MarkCorrect());
              },
              icon: const Icon(Icons.check),
              label: const Text('I Got It Right'),
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.green,
                foregroundColor: Colors.white,
                padding:
                    const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
                textStyle: const TextStyle(fontSize: 16),
              ),
            ),
            const SizedBox(width: 16),
            ElevatedButton.icon(
              onPressed: () {
                context.read<PracticeBloc>().add(MarkIncorrect());
              },
              icon: const Icon(Icons.close),
              label: const Text('I Got It Wrong'),
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.red,
                foregroundColor: Colors.white,
                padding:
                    const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
                textStyle: const TextStyle(fontSize: 16),
              ),
            ),
          ],
        ),
      ],
    );
  }

  /// Build view for explanation state (showing explanation)
  Widget _buildExplanationView(BuildContext context, ExplanationState state) {
    return SingleChildScrollView(
      child: Column(
        children: [
          const SizedBox(height: 16),
          QuestionCard(
            sentence: state.currentSentence,
            showAnswer: true,
          ),
          const SizedBox(height: 16),
          ExplanationCard(sentence: state.currentSentence),
          const SizedBox(height: 24),
          ElevatedButton(
            onPressed: () {
              context.read<PracticeBloc>().add(NextQuestion());
            },
            style: ElevatedButton.styleFrom(
              padding: const EdgeInsets.symmetric(horizontal: 48, vertical: 16),
              textStyle: const TextStyle(fontSize: 18),
            ),
            child: const Text('Next Question'),
          ),
          const SizedBox(height: 24),
        ],
      ),
    );
  }

  /// Helper to get current index from state
  int _getCurrentIndex(PracticeState state) {
    if (state is QuestionState) return state.currentIndex;
    if (state is AnswerRevealedState) return state.currentIndex;
    if (state is ExplanationState) return state.currentIndex;
    return 0;
  }

  /// Helper to get total questions from state
  int _getTotalQuestions(PracticeState state) {
    if (state is QuestionState) return state.totalQuestions;
    if (state is AnswerRevealedState) return state.totalQuestions;
    if (state is ExplanationState) return state.totalQuestions;
    return 0;
  }
}
