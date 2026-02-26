import '../../models/sentence.dart';
import '../../models/session_config.dart';

/// Abstract base class for all Practice screen states
abstract class PracticeState {}

/// Loading state - generating sentences
class PracticeLoading extends PracticeState {}

/// Showing a question (answer hidden)
class QuestionState extends PracticeState {
  final Sentence currentSentence;
  final int currentIndex;
  final int totalQuestions;
  final int correctCount;
  final int incorrectCount;

  QuestionState({
    required this.currentSentence,
    required this.currentIndex,
    required this.totalQuestions,
    required this.correctCount,
    required this.incorrectCount,
  });
}

/// Answer has been revealed (waiting for correct/incorrect input)
class AnswerRevealedState extends PracticeState {
  final Sentence currentSentence;
  final int currentIndex;
  final int totalQuestions;
  final int correctCount;
  final int incorrectCount;

  AnswerRevealedState({
    required this.currentSentence,
    required this.currentIndex,
    required this.totalQuestions,
    required this.correctCount,
    required this.incorrectCount,
  });
}

/// Showing explanation (waiting for next question)
class ExplanationState extends PracticeState {
  final Sentence currentSentence;
  final int currentIndex;
  final int totalQuestions;
  final int correctCount;
  final int incorrectCount;
  final bool wasCorrect; // Whether user marked this as correct

  ExplanationState({
    required this.currentSentence,
    required this.currentIndex,
    required this.totalQuestions,
    required this.correctCount,
    required this.incorrectCount,
    required this.wasCorrect,
  });
}

/// Practice session complete
class PracticeComplete extends PracticeState {
  final int totalQuestions;
  final int correctCount;
  final int incorrectCount;

  PracticeComplete({
    required this.totalQuestions,
    required this.correctCount,
    required this.incorrectCount,
  });

  /// Calculate percentage correct
  double get percentage =>
      totalQuestions > 0 ? (correctCount / totalQuestions) * 100 : 0;
}

/// Error occurred during practice
class PracticeError extends PracticeState {
  final String message;

  PracticeError(this.message);
}
