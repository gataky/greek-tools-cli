import 'dart:math';
import 'package:sqflite/sqflite.dart';
import '../models/noun.dart';
import '../models/sentence.dart';
import '../models/template.dart';
import 'database_helper.dart';

/// DatabaseRepository handles all database queries and sentence generation
/// Mirrors the Go repository pattern from internal/storage/
class DatabaseRepository {
  final DatabaseHelper _dbHelper;

  DatabaseRepository(this._dbHelper);

  /// Get random templates filtered by phase and number
  /// Mirrors GetRandomTemplates from Go's templates.go
  Future<List<SentenceTemplate>> getRandomTemplates(
    int phase,
    String number,
    int limit,
  ) async {
    final db = await _dbHelper.database;

    String whereClause;
    List<dynamic> whereArgs;

    // Build query based on number filter (handle 'both' case)
    if (number.isEmpty || number == 'both') {
      // Include templates for singular, plural, and both
      whereClause =
          'difficulty_phase = ? AND (number = "singular" OR number = "plural" OR number = "both")';
      whereArgs = [phase];
    } else {
      // Filter by specific number (singular or plural), but also include templates marked as 'both'
      whereClause =
          'difficulty_phase = ? AND (number = ? OR number = "both")';
      whereArgs = [phase, number];
    }

    final results = await db.query(
      'sentence_templates',
      where: whereClause,
      whereArgs: whereArgs,
      orderBy: 'RANDOM()',
      limit: limit,
    );

    return results.map((map) => SentenceTemplate.fromMap(map)).toList();
  }

  /// Get random nouns from the database
  Future<List<Noun>> getRandomNouns(int limit) async {
    final db = await _dbHelper.database;
    final results = await db.query(
      'nouns',
      orderBy: 'RANDOM()',
      limit: limit,
    );
    return results.map((map) => Noun.fromMap(map)).toList();
  }

  /// Get all nouns (used for sentence generation)
  Future<List<Noun>> getAllNouns() async {
    final db = await _dbHelper.database;
    final results = await db.query('nouns', orderBy: 'id');
    return results.map((map) => Noun.fromMap(map)).toList();
  }

  /// Template substitution - generates a Sentence from a template and noun
  /// Mirrors substituteTemplate from Go's template_generator.go
  Sentence _substituteTemplate(SentenceTemplate template, Noun noun) {
    // 1. Get article and noun form using noun.getField()
    final article = noun.getField(template.articleField);
    final nounForm = noun.getField(template.nounFormField);

    // 2. Substitute placeholders in English template
    final englishPrompt =
        template.englishTemplate.replaceAll('{noun}', noun.english);

    // 3. Substitute placeholders in Greek template
    String greekSentence = template.greekTemplate;
    greekSentence = greekSentence.replaceAll('{article}', article);
    greekSentence = greekSentence.replaceAll('{noun_form}', nounForm);

    // 4. Generate correct answer (article + space + noun form)
    final correctAnswer = '$article $nounForm';

    // 5. Determine actual number (handle 'both' by inferring from field name)
    String number = template.number;
    if (number == 'both') {
      if (template.nounFormField.contains('Sg')) {
        number = 'singular';
      } else if (template.nounFormField.contains('Pl')) {
        number = 'plural';
      }
    }

    // 6. Create Sentence object
    return Sentence(
      nounId: noun.id,
      englishPrompt: englishPrompt,
      greekSentence: greekSentence,
      correctAnswer: correctAnswer,
      caseType: template.caseType,
      number: number,
      difficultyPhase: template.difficultyPhase,
      contextType: template.contextType,
      preposition: template.preposition,
    );
  }

  /// Generate practice sentences from templates
  /// Mirrors GeneratePracticeSentences from Go's template_generator.go
  Future<List<Sentence>> generatePracticeSentences(
    int phase,
    String number,
    int limit,
  ) async {
    // 1. Get all nouns
    final nouns = await getAllNouns();
    if (nouns.isEmpty) {
      throw Exception('No nouns found in database');
    }

    // 2. Get templates matching filters (get more than needed for variety)
    int templateLimit = limit * 2;
    if (templateLimit < 100) {
      templateLimit = 100;
    }
    final templates = await getRandomTemplates(phase, number, templateLimit);

    if (templates.isEmpty) {
      throw Exception('No templates found for phase $phase and number $number');
    }

    // 3. Generate sentences by combining templates with random nouns
    final sentences = <Sentence>[];
    final used = <String>{}; // Track used combinations to avoid duplicates
    final random = Random();

    // Try to generate the requested number of sentences
    int attempts = 0;
    final maxAttempts = limit * 10; // Prevent infinite loop

    while (sentences.length < limit && attempts < maxAttempts) {
      attempts++;

      // Pick random template and noun
      final template = templates[random.nextInt(templates.length)];
      final noun = nouns[random.nextInt(nouns.length)];

      // Create unique key for this combination
      final key = '${template.id}-${noun.id}';
      if (used.contains(key)) {
        continue; // Skip if we've already used this combination
      }

      try {
        // Generate sentence
        final sentence = _substituteTemplate(template, noun);
        sentences.add(sentence);
        used.add(key);
      } catch (e) {
        // Skip invalid combinations (e.g., field mismatch)
        continue;
      }
    }

    // Shuffle sentences for variety
    sentences.shuffle();

    // If we couldn't generate enough unique combinations, that's okay
    // Just return what we have
    return sentences;
  }
}
