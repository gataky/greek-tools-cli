/// SessionConfig model for practice session configuration
/// Based on PRD requirements and mirrors Go session configuration logic
class SessionConfig {
  final String difficultyLevel; // "beginner", "intermediate", "advanced"
  final bool includePlural;
  final int questionCount; // 0 for endless

  SessionConfig({
    required this.difficultyLevel,
    required this.includePlural,
    required this.questionCount,
  });

  /// Map difficulty level to phase number
  /// Beginner = 1, Intermediate = 2, Advanced = 3
  int get phase {
    switch (difficultyLevel) {
      case 'beginner':
        return 1;
      case 'intermediate':
        return 2;
      case 'advanced':
        return 3;
      default:
        return 1;
    }
  }

  /// Get number filter for database queries
  /// Returns 'singular' if plural not included, empty string otherwise
  String get numberFilter => includePlural ? '' : 'singular';

  /// Convert SessionConfig to map for serialization
  Map<String, dynamic> toJson() {
    return {
      'difficulty_level': difficultyLevel,
      'include_plural': includePlural,
      'question_count': questionCount,
    };
  }

  /// Factory constructor to create SessionConfig from map
  factory SessionConfig.fromJson(Map<String, dynamic> json) {
    return SessionConfig(
      difficultyLevel: json['difficulty_level'] as String,
      includePlural: json['include_plural'] as bool,
      questionCount: json['question_count'] as int,
    );
  }
}
