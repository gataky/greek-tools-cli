/// Events for Setup screen
/// Used to update configuration before starting practice session
abstract class SetupEvent {}

/// User changed difficulty level
class DifficultyChanged extends SetupEvent {
  final String difficulty; // "beginner", "intermediate", "advanced"

  DifficultyChanged(this.difficulty);
}

/// User toggled the plural checkbox
class PluralToggled extends SetupEvent {}

/// User changed question count
class QuestionCountChanged extends SetupEvent {
  final int count; // 10, 20, 50, or 0 for endless

  QuestionCountChanged(this.count);
}

/// User tapped Start Practice button
class StartPractice extends SetupEvent {}
