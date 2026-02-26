/// SentenceTemplate model representing a reusable sentence pattern
/// Mirrors the Go SentenceTemplate struct from internal/models/template.go
class SentenceTemplate {
  final int id;
  final String englishTemplate;
  final String greekTemplate;
  final String articleField;
  final String nounFormField;
  final String caseType;
  final String number;
  final int difficultyPhase;
  final String contextType;
  final String? preposition;

  SentenceTemplate({
    required this.id,
    required this.englishTemplate,
    required this.greekTemplate,
    required this.articleField,
    required this.nounFormField,
    required this.caseType,
    required this.number,
    required this.difficultyPhase,
    required this.contextType,
    this.preposition,
  });

  /// Factory constructor to create SentenceTemplate from SQLite query result
  factory SentenceTemplate.fromMap(Map<String, dynamic> map) {
    return SentenceTemplate(
      id: map['id'] as int,
      englishTemplate: map['english_template'] as String,
      greekTemplate: map['greek_template'] as String,
      articleField: map['article_field'] as String,
      nounFormField: map['noun_form_field'] as String,
      caseType: map['case_type'] as String,
      number: map['number'] as String,
      difficultyPhase: map['difficulty_phase'] as int,
      contextType: map['context_type'] as String,
      preposition: map['preposition'] as String?,
    );
  }
}
