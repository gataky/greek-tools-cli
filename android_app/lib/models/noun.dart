/// Noun model representing a Greek noun with all declined forms
/// Mirrors the Go Noun struct from internal/models/noun.go
class Noun {
  final int id;
  final String english;
  final String gender;
  final String nominativeSg;
  final String genitiveSg;
  final String accusativeSg;
  final String nominativePl;
  final String genitivePl;
  final String accusativePl;
  final String nomSgArticle;
  final String genSgArticle;
  final String accSgArticle;
  final String nomPlArticle;
  final String genPlArticle;
  final String accPlArticle;

  Noun({
    required this.id,
    required this.english,
    required this.gender,
    required this.nominativeSg,
    required this.genitiveSg,
    required this.accusativeSg,
    required this.nominativePl,
    required this.genitivePl,
    required this.accusativePl,
    required this.nomSgArticle,
    required this.genSgArticle,
    required this.accSgArticle,
    required this.nomPlArticle,
    required this.genPlArticle,
    required this.accPlArticle,
  });

  /// Factory constructor to create Noun from SQLite query result
  factory Noun.fromMap(Map<String, dynamic> map) {
    return Noun(
      id: map['id'] as int,
      english: map['english'] as String,
      gender: map['gender'] as String,
      nominativeSg: map['nominative_sg'] as String,
      genitiveSg: map['genitive_sg'] as String,
      accusativeSg: map['accusative_sg'] as String,
      nominativePl: map['nominative_pl'] as String,
      genitivePl: map['genitive_pl'] as String,
      accusativePl: map['accusative_pl'] as String,
      nomSgArticle: map['nom_sg_article'] as String,
      genSgArticle: map['gen_sg_article'] as String,
      accSgArticle: map['acc_sg_article'] as String,
      nomPlArticle: map['nom_pl_article'] as String,
      genPlArticle: map['gen_pl_article'] as String,
      accPlArticle: map['acc_pl_article'] as String,
    );
  }

  /// Get field value by name (for template substitution)
  /// Mirrors the getFieldValue function from Go's template_generator.go
  String getField(String fieldName) {
    switch (fieldName) {
      case 'NomSgArticle':
        return nomSgArticle;
      case 'GenSgArticle':
        return genSgArticle;
      case 'AccSgArticle':
        return accSgArticle;
      case 'NomPlArticle':
        return nomPlArticle;
      case 'GenPlArticle':
        return genPlArticle;
      case 'AccPlArticle':
        return accPlArticle;
      case 'NominativeSg':
        return nominativeSg;
      case 'GenitiveSg':
        return genitiveSg;
      case 'AccusativeSg':
        return accusativeSg;
      case 'NominativePl':
        return nominativePl;
      case 'GenitivePl':
        return genitivePl;
      case 'AccusativePl':
        return accusativePl;
      default:
        throw Exception('Unknown field: $fieldName');
    }
  }
}
