import '../../models/session_config.dart';

/// Events for Practice screen
/// Handles the practice session flow
abstract class PracticeEvent {}

/// Load a new practice session with given configuration
class LoadSession extends PracticeEvent {
  final SessionConfig config;

  LoadSession(this.config);
}

/// User tapped Show Answer button
class ShowAnswer extends PracticeEvent {}

/// User marked answer as correct
class MarkCorrect extends PracticeEvent {}

/// User marked answer as incorrect
class MarkIncorrect extends PracticeEvent {}

/// User tapped Next Question button
class NextQuestion extends PracticeEvent {}
