import 'package:flutter_bloc/flutter_bloc.dart';
import '../../data/repository.dart';
import '../../models/sentence.dart';
import 'practice_event.dart';
import 'practice_state.dart';

/// PracticeBloc manages the state of the Practice screen
/// Handles question flow, answer checking, and session completion
class PracticeBloc extends Bloc<PracticeEvent, PracticeState> {
  final DatabaseRepository repository;

  // Session data
  List<Sentence> _sentences = [];
  int _currentIndex = 0;
  int _correctCount = 0;
  int _incorrectCount = 0;

  PracticeBloc(this.repository) : super(PracticeLoading()) {
    // Register event handlers
    on<LoadSession>(_onLoadSession);
    on<ShowAnswer>(_onShowAnswer);
    on<MarkCorrect>(_onMarkCorrect);
    on<MarkIncorrect>(_onMarkIncorrect);
    on<NextQuestion>(_onNextQuestion);
  }

  /// Load a new practice session
  Future<void> _onLoadSession(
    LoadSession event,
    Emitter<PracticeState> emit,
  ) async {
    emit(PracticeLoading());

    try {
      // Generate sentences using repository
      final config = event.config;
      final limit =
          config.questionCount == 0 ? 1000 : config.questionCount; // Endless mode = 1000 questions

      _sentences = await repository.generatePracticeSentences(
        config.phase,
        config.numberFilter,
        limit,
      );

      if (_sentences.isEmpty) {
        emit(PracticeError('No sentences generated. Check database.'));
        return;
      }

      // Reset counters
      _currentIndex = 0;
      _correctCount = 0;
      _incorrectCount = 0;

      // Emit first question
      emit(QuestionState(
        currentSentence: _sentences[_currentIndex],
        currentIndex: _currentIndex + 1, // Display as 1-indexed
        totalQuestions: config.questionCount == 0 ? 0 : _sentences.length,
        correctCount: _correctCount,
        incorrectCount: _incorrectCount,
      ));
    } catch (e) {
      emit(PracticeError('Failed to load session: $e'));
    }
  }

  /// Show answer for current question
  void _onShowAnswer(
    ShowAnswer event,
    Emitter<PracticeState> emit,
  ) {
    if (state is QuestionState) {
      final currentState = state as QuestionState;
      emit(AnswerRevealedState(
        currentSentence: currentState.currentSentence,
        currentIndex: currentState.currentIndex,
        totalQuestions: currentState.totalQuestions,
        correctCount: currentState.correctCount,
        incorrectCount: currentState.incorrectCount,
      ));
    }
  }

  /// User marked answer as correct
  void _onMarkCorrect(
    MarkCorrect event,
    Emitter<PracticeState> emit,
  ) {
    if (state is AnswerRevealedState) {
      final currentState = state as AnswerRevealedState;
      _correctCount++;

      emit(ExplanationState(
        currentSentence: currentState.currentSentence,
        currentIndex: currentState.currentIndex,
        totalQuestions: currentState.totalQuestions,
        correctCount: _correctCount,
        incorrectCount: _incorrectCount,
        wasCorrect: true,
      ));
    }
  }

  /// User marked answer as incorrect
  void _onMarkIncorrect(
    MarkIncorrect event,
    Emitter<PracticeState> emit,
  ) {
    if (state is AnswerRevealedState) {
      final currentState = state as AnswerRevealedState;
      _incorrectCount++;

      emit(ExplanationState(
        currentSentence: currentState.currentSentence,
        currentIndex: currentState.currentIndex,
        totalQuestions: currentState.totalQuestions,
        correctCount: _correctCount,
        incorrectCount: _incorrectCount,
        wasCorrect: false,
      ));
    }
  }

  /// Advance to next question
  void _onNextQuestion(
    NextQuestion event,
    Emitter<PracticeState> emit,
  ) {
    if (state is ExplanationState) {
      final currentState = state as ExplanationState;
      _currentIndex++;

      // Check if we've completed all questions
      final isEndless = currentState.totalQuestions == 0;
      if (!isEndless && _currentIndex >= _sentences.length) {
        // Session complete
        emit(PracticeComplete(
          totalQuestions: _sentences.length,
          correctCount: _correctCount,
          incorrectCount: _incorrectCount,
        ));
      } else {
        // More questions remaining
        emit(QuestionState(
          currentSentence: _sentences[_currentIndex],
          currentIndex: _currentIndex + 1, // Display as 1-indexed
          totalQuestions: currentState.totalQuestions,
          correctCount: _correctCount,
          incorrectCount: _incorrectCount,
        ));
      }
    }
  }
}
