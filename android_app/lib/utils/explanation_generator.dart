import '../models/sentence.dart';

/// ExplanationGenerator creates grammar explanations for practice sentences
/// Mirrors logic from Go's internal/explanations/templates.go
class ExplanationGenerator {
  /// Generate syntactic role explanation based on context type
  static String generateSyntacticRole(
    String contextType,
    String caseType,
    String? preposition,
  ) {
    switch (contextType) {
      case 'direct_object':
        return 'Direct objects use accusative case';
      case 'possession':
        return 'Possession requires genitive case';
      case 'preposition':
        if (preposition != null && preposition.isNotEmpty) {
          return 'The preposition "$preposition" requires $caseType case';
        }
        return 'This preposition requires $caseType case';
      case 'subject':
        return 'Subjects use nominative case';
      default:
        return 'This context uses $caseType case';
    }
  }

  /// Generate morphology explanation (case, gender, number, article)
  static String generateMorphology(Sentence sentence) {
    // Capitalize case type
    final caseLabel =
        sentence.caseType[0].toUpperCase() + sentence.caseType.substring(1);

    // Extract article from correct answer (first word)
    final parts = sentence.correctAnswer.split(' ');
    final article = parts.isNotEmpty ? parts[0] : '';

    // Infer gender from article (simplified)
    String gender = _inferGender(article, sentence.caseType, sentence.number);

    return '$caseLabel $gender ${sentence.number}\nArticle: $article';
  }

  /// Generate translation from English prompt
  static String generateTranslation(String englishPrompt) {
    // Remove the blank marker "___" and parentheses for cleaner translation
    String translation = englishPrompt
        .replaceAll('___', 'the')
        .replaceAll(RegExp(r'\s*\([^)]*\)'), '');

    return translation;
  }

  /// Infer gender from article (simplified heuristic)
  static String _inferGender(String article, String caseType, String number) {
    // Masculine articles
    if (_isMasculine(article)) {
      return 'masculine';
    }
    // Feminine articles
    else if (_isFeminine(article)) {
      return 'feminine';
    }
    // Neuter articles
    else if (_isNeuter(article)) {
      return 'neuter';
    }

    return 'unknown';
  }

  static bool _isMasculine(String article) {
    return ['ο', 'του', 'τον', 'οι', 'των', 'τους'].contains(article);
  }

  static bool _isFeminine(String article) {
    return ['η', 'της', 'την', 'τη', 'οι', 'των', 'τις'].contains(article);
  }

  static bool _isNeuter(String article) {
    return ['το', 'του', 'τα', 'των'].contains(article);
  }
}
