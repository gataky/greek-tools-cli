/// State for Setup screen
/// Tracks user's configuration choices
class SetupState {
  final String difficultyLevel; // "beginner", "intermediate", "advanced"
  final bool includePlural;
  final int questionCount; // 10, 20, 50, or 0 for endless

  SetupState({
    required this.difficultyLevel,
    required this.includePlural,
    required this.questionCount,
  });

  /// Create initial state with default values
  factory SetupState.initial() {
    return SetupState(
      difficultyLevel: 'beginner',
      includePlural: false,
      questionCount: 10,
    );
  }

  /// Create a copy of the state with updated fields
  SetupState copyWith({
    String? difficultyLevel,
    bool? includePlural,
    int? questionCount,
  }) {
    return SetupState(
      difficultyLevel: difficultyLevel ?? this.difficultyLevel,
      includePlural: includePlural ?? this.includePlural,
      questionCount: questionCount ?? this.questionCount,
    );
  }
}
