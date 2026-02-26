/// Sentence model representing a practice sentence
/// Mirrors the Go Sentence struct from internal/models/sentence.go
class Sentence {
  final int? id;
  final int nounId;
  final String englishPrompt;
  final String greekSentence;
  final String correctAnswer;
  final String caseType;
  final String number;
  final int difficultyPhase;
  final String contextType;
  final String? preposition;

  Sentence({
    this.id,
    required this.nounId,
    required this.englishPrompt,
    required this.greekSentence,
    required this.correctAnswer,
    required this.caseType,
    required this.number,
    required this.difficultyPhase,
    required this.contextType,
    this.preposition,
  });

  /// Factory constructor to create Sentence from map (for deserialization)
  factory Sentence.fromMap(Map<String, dynamic> map) {
    return Sentence(
      id: map['id'] as int?,
      nounId: map['noun_id'] as int,
      englishPrompt: map['english_prompt'] as String,
      greekSentence: map['greek_sentence'] as String,
      correctAnswer: map['correct_answer'] as String,
      caseType: map['case_type'] as String,
      number: map['number'] as String,
      difficultyPhase: map['difficulty_phase'] as int,
      contextType: map['context_type'] as String,
      preposition: map['preposition'] as String?,
    );
  }

  /// Convert Sentence to map (for serialization)
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'noun_id': nounId,
      'english_prompt': englishPrompt,
      'greek_sentence': greekSentence,
      'correct_answer': correctAnswer,
      'case_type': caseType,
      'number': number,
      'difficulty_phase': difficultyPhase,
      'context_type': contextType,
      'preposition': preposition,
    };
  }
}
