import 'package:flutter_bloc/flutter_bloc.dart';
import 'setup_event.dart';
import 'setup_state.dart';

/// SetupBloc manages the state of the Setup screen
/// Handles user configuration choices before starting practice
class SetupBloc extends Bloc<SetupEvent, SetupState> {
  SetupBloc() : super(SetupState.initial()) {
    // Register event handlers
    on<DifficultyChanged>(_onDifficultyChanged);
    on<PluralToggled>(_onPluralToggled);
    on<QuestionCountChanged>(_onQuestionCountChanged);
    on<StartPractice>(_onStartPractice);
  }

  /// Handle difficulty level change
  void _onDifficultyChanged(
    DifficultyChanged event,
    Emitter<SetupState> emit,
  ) {
    emit(state.copyWith(difficultyLevel: event.difficulty));
  }

  /// Handle plural toggle
  void _onPluralToggled(
    PluralToggled event,
    Emitter<SetupState> emit,
  ) {
    emit(state.copyWith(includePlural: !state.includePlural));
  }

  /// Handle question count change
  void _onQuestionCountChanged(
    QuestionCountChanged event,
    Emitter<SetupState> emit,
  ) {
    emit(state.copyWith(questionCount: event.count));
  }

  /// Handle start practice button
  /// This just acknowledges the event - navigation is handled in UI layer
  void _onStartPractice(
    StartPractice event,
    Emitter<SetupState> emit,
  ) {
    // State doesn't change, but event is processed
    // UI will listen for this event and navigate to Practice screen
  }
}
