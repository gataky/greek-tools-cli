# PRD: Flutter Android App for Greek Case Practice

## Introduction/Overview

Create a personal Flutter Android app that provides Greek noun case practice using the existing template-based database system. The app will be a simplified, mobile version of the CLI practice functionality, focusing on self-checking (user reports if they got the answer right or wrong) rather than input validation. The app is for personal use only, not public release.

The core workflow: user selects difficulty/settings → practices questions by viewing Greek sentences with blanks → reveals answer → self-reports correct/incorrect → views grammar explanation → moves to next question.

## Goals

### Functional Goals
1. Enable Greek case practice on Android devices without internet connectivity
2. Support all three difficulty levels (Beginner, Intermediate, Advanced) from the existing system
3. Allow users to practice with singular-only or include plural forms
4. Provide grammar explanations after each answer (matching CLI behavior)
5. Support session lengths of 10, 20, 50 questions, or Endless mode

### Technical Goals
1. App runs fully offline with 244KB database bundled in assets
2. App size under 30MB total
3. Smooth UI performance with no lag when rendering Greek text
4. Session state persists through app minimization (resume where you left off)
5. Support Android 11+ (API 30+) devices

## Tech Stack & Architecture

### Languages & Frameworks
- **Flutter 3.24+**: Cross-platform framework (though only targeting Android initially)
- **Dart 3.5+**: Programming language for Flutter
- **sqflite 2.3+**: SQLite plugin for Flutter
- **path_provider**: For accessing app directories

### Architectural Pattern
- **BLoC Pattern (Business Logic Component)**: For state management
  - Setup screen has SetupBloc
  - Practice screen has PracticeBloc
  - Simple, testable separation of UI and business logic
- **Repository Pattern**: Database access abstracted through repository layer
  - Mirrors existing Go code structure: `DatabaseRepository` with methods like `generatePracticeSentences()`
- **SQLite Database**: Embedded in app assets, copied to device on first launch

### Data Storage
- **Embedded Database**: `greekmaster.db` (244KB) bundled with app
- **App-level storage**: Session state saved to local storage for resume functionality
- **No remote sync**: Fully offline, no cloud services

### External Dependencies
```yaml
# pubspec.yaml
dependencies:
  flutter:
    sdk: flutter
  sqflite: ^2.3.0
  path_provider: ^2.1.0
  path: ^1.8.3
  flutter_bloc: ^8.1.3
```

### Existing Patterns to Follow
The Go codebase provides reference implementation:
- **Practice logic**: `internal/tui/practice.go` - maps difficulty to phase, handles question flow
- **Template substitution**: `internal/storage/template_generator.go` - generates sentences from templates
- **Session config**: `internal/models/session.go` - configuration structure
- **Database queries**: `internal/storage/repository.go` - GeneratePracticeSentences method

## Functional Requirements

### FR1: Setup Screen
The app must provide a setup screen where users configure their practice session:
- **Difficulty selection**: Radio buttons for Beginner / Intermediate / Advanced
- **Plural inclusion**: Checkbox for "Include plural forms"
- **Question count**: Dropdown or segmented control with options: 10, 20, 50, Endless
- **Start button**: Navigates to Practice screen when tapped

Map difficulty to phase: Beginner=1, Intermediate=2, Advanced=3 (matching Go code).

### FR2: Practice Question Display
For each question, the app must display:
- **Progress indicator**: "Question X of Y" at top (or "Question X" for endless)
- **English prompt**: e.g., "I see {noun}" with blanks showing the missing noun
- **Greek sentence**: e.g., "Βλέπω ___ ___" with blanks where article + noun should go
- **Show Answer button**: Centered below the Greek sentence

Greek text must render correctly using appropriate fonts (Noto Sans Greek or system default).

### FR3: Answer Reveal & Self-Checking
When user taps "Show Answer":
1. Button disappears
2. Greek sentence updates to show complete answer: "Βλέπω τον δάσκαλο"
3. English prompt updates to show full translation: "I see teacher"
4. Two new buttons appear:
   - "✓ I Got It Right" (green)
   - "✗ I Got It Wrong" (red)

User taps one button to self-report correctness.

### FR4: Grammar Explanation Display
After user marks answer as correct/incorrect:
1. Display grammar explanation in a card/panel below the answer
2. Explanation includes:
   - **Translation**: Full English sentence
   - **Syntactic Role**: Rule explanation (e.g., "Direct objects use accusative case")
   - **Morphology**: Case analysis (e.g., "Accusative masculine singular with article τον")
3. "Next Question" button appears to advance

Generate explanation using same logic as Go code:
- Syntactic role from `internal/explanations/templates.go:SyntacticRoleTemplate()`
- Morphology describes the case, gender, number, article used

### FR5: Session Completion
When all questions are answered (or user completes endless session):
1. Navigate to Results screen
2. Display:
   - Total questions: X
   - Correct: Y (Z%)
   - Incorrect: W
   - "Back to Setup" button
3. Tapping button returns to Setup screen

### FR6: Sentence Generation
The app must generate practice sentences at runtime using templates:
1. Query `sentence_templates` table filtered by phase and number
2. Query `nouns` table for random nouns
3. For each template+noun pair, substitute placeholders:
   - Replace `{noun}` in English template with `noun.english`
   - Replace `{article}` in Greek template with appropriate article field (e.g., `noun.acc_sg_article`)
   - Replace `{noun_form}` in Greek template with appropriate noun form (e.g., `noun.accusative_sg`)
4. Shuffle generated sentences
5. Limit to question count (or generate 1000 for endless mode)

This mirrors the Go `substituteTemplate()` function in `internal/storage/template_generator.go`.

### FR7: Session Persistence
The app must save session state when minimized or closed:
- Current question index
- List of generated sentences
- Correct/incorrect counts
- User's answers so far

When app resumes, check if saved state exists (<24 hours old). If yes, restore session. If no, show Setup screen.

### FR8: Database Initialization
On first app launch:
1. Check if database exists in app's documents directory
2. If not, copy `greekmaster.db` from assets to documents directory
3. Open database connection
4. Verify tables exist (sentence_templates, nouns)

Database remains in documents directory for all future sessions.

## Technical Specifications

### Data Models

#### Dart Models (Mirror Go Structs)

```dart
// lib/models/sentence.dart
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
}

// lib/models/template.dart
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

  factory SentenceTemplate.fromMap(Map<String, dynamic> map) {
    return SentenceTemplate(
      id: map['id'],
      englishTemplate: map['english_template'],
      greekTemplate: map['greek_template'],
      articleField: map['article_field'],
      nounFormField: map['noun_form_field'],
      caseType: map['case_type'],
      number: map['number'],
      difficultyPhase: map['difficulty_phase'],
      contextType: map['context_type'],
      preposition: map['preposition'],
    );
  }
}

// lib/models/noun.dart
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

  factory Noun.fromMap(Map<String, dynamic> map) {
    return Noun(
      id: map['id'],
      english: map['english'],
      gender: map['gender'],
      nominativeSg: map['nominative_sg'],
      genitiveSg: map['genitive_sg'],
      accusativeSg: map['accusative_sg'],
      nominativePl: map['nominative_pl'],
      genitivePl: map['genitive_pl'],
      accusativePl: map['accusative_pl'],
      nomSgArticle: map['nom_sg_article'],
      genSgArticle: map['gen_sg_article'],
      accSgArticle: map['acc_sg_article'],
      nomPlArticle: map['nom_pl_article'],
      genPlArticle: map['gen_pl_article'],
      accPlArticle: map['acc_pl_article'],
    );
  }

  // Get field value by name (for template substitution)
  String getField(String fieldName) {
    switch (fieldName) {
      case 'NomSgArticle': return nomSgArticle;
      case 'GenSgArticle': return genSgArticle;
      case 'AccSgArticle': return accSgArticle;
      case 'NomPlArticle': return nomPlArticle;
      case 'GenPlArticle': return genPlArticle;
      case 'AccPlArticle': return accPlArticle;
      case 'NominativeSg': return nominativeSg;
      case 'GenitiveSg': return genitiveSg;
      case 'AccusativeSg': return accusativeSg;
      case 'NominativePl': return nominativePl;
      case 'GenitivePl': return genitivePl;
      case 'AccusativePl': return accusativePl;
      default: throw Exception('Unknown field: $fieldName');
    }
  }
}

// lib/models/session_config.dart
class SessionConfig {
  final String difficultyLevel; // "beginner", "intermediate", "advanced"
  final bool includePlural;
  final int questionCount; // 0 for endless

  SessionConfig({
    required this.difficultyLevel,
    required this.includePlural,
    required this.questionCount,
  });

  int get phase {
    switch (difficultyLevel) {
      case 'beginner': return 1;
      case 'intermediate': return 2;
      case 'advanced': return 3;
      default: return 1;
    }
  }

  String get numberFilter => includePlural ? '' : 'singular';
}
```

### Component/Module Structure

```
lib/
├── main.dart                      # App entry point
├── models/
│   ├── noun.dart                  # Noun model
│   ├── sentence.dart              # Sentence model
│   ├── template.dart              # SentenceTemplate model
│   └── session_config.dart        # SessionConfig model
├── data/
│   ├── database_helper.dart       # SQLite database management
│   └── repository.dart            # Database repository (queries + generation)
├── blocs/
│   ├── setup/
│   │   ├── setup_bloc.dart        # Setup screen state management
│   │   ├── setup_event.dart       # Setup events (DifficultyChanged, etc.)
│   │   └── setup_state.dart       # Setup states
│   └── practice/
│       ├── practice_bloc.dart     # Practice screen state management
│       ├── practice_event.dart    # Practice events (ShowAnswer, MarkCorrect, etc.)
│       └── practice_state.dart    # Practice states (QuestionState, AnswerRevealedState, etc.)
├── screens/
│   ├── setup_screen.dart          # Setup/configuration screen
│   ├── practice_screen.dart       # Practice question screen
│   └── results_screen.dart        # Final score screen
└── widgets/
    ├── difficulty_selector.dart   # Difficulty radio buttons
    ├── question_card.dart         # Question display widget
    └── explanation_card.dart      # Grammar explanation widget
```

### Database Repository Interface

```dart
// lib/data/repository.dart
class DatabaseRepository {
  final DatabaseHelper _db;

  DatabaseRepository(this._db);

  // Get random templates filtered by phase and number
  Future<List<SentenceTemplate>> getRandomTemplates(
    int phase,
    String number,
    int limit,
  ) async {
    final db = await _db.database;

    String whereClause;
    List<dynamic> whereArgs;

    if (number.isEmpty) {
      whereClause = 'difficulty_phase = ? AND (number = "singular" OR number = "plural" OR number = "both")';
      whereArgs = [phase];
    } else {
      whereClause = 'difficulty_phase = ? AND (number = ? OR number = "both")';
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

  // Get random nouns
  Future<List<Noun>> getRandomNouns(int limit) async {
    final db = await _db.database;
    final results = await db.query(
      'nouns',
      orderBy: 'RANDOM()',
      limit: limit,
    );
    return results.map((map) => Noun.fromMap(map)).toList();
  }

  // Generate practice sentences from templates
  Future<List<Sentence>> generatePracticeSentences(
    int phase,
    String number,
    int limit,
  ) async {
    // Get templates
    final templates = await getRandomTemplates(phase, number, limit * 2);

    // Get nouns
    final nouns = await getRandomNouns(50); // Get variety

    if (templates.isEmpty || nouns.isEmpty) {
      throw Exception('No templates or nouns found');
    }

    final sentences = <Sentence>[];
    final random = Random();

    // Generate up to 'limit' sentences
    for (int i = 0; i < limit && sentences.length < limit; i++) {
      final template = templates[random.nextInt(templates.length)];
      final noun = nouns[random.nextInt(nouns.length)];

      try {
        final sentence = _substituteTemplate(template, noun);
        sentences.add(sentence);
      } catch (e) {
        continue; // Skip invalid combinations
      }
    }

    // Shuffle sentences
    sentences.shuffle();

    return sentences;
  }

  // Template substitution (mirrors Go substituteTemplate function)
  Sentence _substituteTemplate(SentenceTemplate template, Noun noun) {
    // Get article and noun form
    final article = noun.getField(template.articleField);
    final nounForm = noun.getField(template.nounFormField);

    // Substitute placeholders
    final englishPrompt = template.englishTemplate.replaceAll('{noun}', noun.english);
    final greekSentence = template.greekTemplate
        .replaceAll('{article}', article)
        .replaceAll('{noun_form}', nounForm);
    final correctAnswer = '$article $nounForm';

    // Determine actual number (handle "both")
    String number = template.number;
    if (number == 'both') {
      number = template.nounFormField.contains('Sg') ? 'singular' : 'plural';
    }

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
}
```

### UI Layout Specifications

#### Setup Screen
```
┌─────────────────────────────────────┐
│  Greek Case Practice                │
├─────────────────────────────────────┤
│                                     │
│  Difficulty Level                   │
│  ○ Beginner                         │
│  ○ Intermediate                     │
│  ○ Advanced                         │
│                                     │
│  ☐ Include plural forms             │
│                                     │
│  Number of Questions                │
│  [10 ▼]                             │
│                                     │
│  [ Start Practice ]                 │
│                                     │
└─────────────────────────────────────┘
```

#### Practice Screen (Question State)
```
┌─────────────────────────────────────┐
│  Question 3 of 20                   │
├─────────────────────────────────────┤
│                                     │
│  I see ___ (the teacher)            │
│                                     │
│  Βλέπω ___ ___                      │
│                                     │
│  [ Show Answer ]                    │
│                                     │
└─────────────────────────────────────┘
```

#### Practice Screen (Answer Revealed State)
```
┌─────────────────────────────────────┐
│  Question 3 of 20                   │
├─────────────────────────────────────┤
│                                     │
│  I see teacher                      │
│                                     │
│  Βλέπω τον δάσκαλο                 │
│                                     │
│  [ ✓ I Got It Right ]               │
│  [ ✗ I Got It Wrong ]               │
│                                     │
└─────────────────────────────────────┘
```

#### Practice Screen (Explanation State)
```
┌─────────────────────────────────────┐
│  Question 3 of 20                   │
├─────────────────────────────────────┤
│                                     │
│  I see teacher                      │
│  Βλέπω τον δάσκαλο                 │
│                                     │
│  ┌───────────────────────────────┐ │
│  │ Explanation                   │ │
│  │                               │ │
│  │ Translation: I see the teacher│ │
│  │                               │ │
│  │ Syntactic Role:               │ │
│  │ Direct objects use accusative │ │
│  │                               │ │
│  │ Morphology:                   │ │
│  │ Accusative masculine singular │ │
│  │ Article: τον                  │ │
│  └───────────────────────────────┘ │
│                                     │
│  [ Next Question ]                  │
│                                     │
└─────────────────────────────────────┘
```

#### Results Screen
```
┌─────────────────────────────────────┐
│  Practice Complete!                 │
├─────────────────────────────────────┤
│                                     │
│  Total Questions: 20                │
│  Correct: 15 (75%)                  │
│  Incorrect: 5                       │
│                                     │
│  [ Back to Setup ]                  │
│                                     │
└─────────────────────────────────────┘
```

### Integration Points

#### Database Assets
- Place `greekmaster.db` in `assets/` directory in Flutter project
- Update `pubspec.yaml` to include asset:
  ```yaml
  flutter:
    assets:
      - assets/greekmaster.db
  ```

#### Greek Text Rendering
- Flutter's default fonts support Greek characters
- Explicitly set TextStyle font family to 'Noto Sans' or system default
- Test on device to ensure proper rendering of accented characters

### Explanation Generation Logic

Port logic from `internal/explanations/templates.go`:

```dart
// lib/utils/explanation_generator.dart
class ExplanationGenerator {
  static String generateSyntacticRole(String contextType, String caseType, String? preposition) {
    switch (contextType) {
      case 'direct_object':
        return 'Direct objects use accusative case';
      case 'possession':
        return 'Possession requires genitive case';
      case 'preposition':
        if (preposition != null && preposition.isNotEmpty) {
          return 'The preposition \'$preposition\' requires $caseType case';
        }
        return 'This preposition requires $caseType case';
      default:
        return 'This context uses $caseType case';
    }
  }

  static String generateMorphology(Sentence sentence, String article) {
    final caseLabel = sentence.caseType.capitalize();
    // Infer gender from article pattern (simplified)
    String gender = 'masculine'; // Could be enhanced

    return '$caseLabel $gender ${sentence.number}\nArticle: $article';
  }

  static String generateTranslation(String englishPrompt) {
    // Simple: just return the English prompt
    // Could be enhanced to generate full sentence
    return englishPrompt;
  }
}
```

### Session State Persistence

```dart
// lib/data/session_storage.dart
import 'package:shared_preferences/shared_preferences.dart';
import 'dart:convert';

class SessionStorage {
  static const String _keySessionState = 'session_state';
  static const String _keyTimestamp = 'session_timestamp';

  // Save session state
  static Future<void> saveSession(PracticeState state) async {
    final prefs = await SharedPreferences.getInstance();

    final sessionData = {
      'currentIndex': state.currentIndex,
      'correctCount': state.correctCount,
      'incorrectCount': state.incorrectCount,
      'sentences': state.sentences.map((s) => s.toJson()).toList(),
      'config': state.config.toJson(),
    };

    await prefs.setString(_keySessionState, jsonEncode(sessionData));
    await prefs.setInt(_keyTimestamp, DateTime.now().millisecondsSinceEpoch);
  }

  // Load session state (returns null if expired or not found)
  static Future<Map<String, dynamic>?> loadSession() async {
    final prefs = await SharedPreferences.getInstance();

    final timestamp = prefs.getInt(_keyTimestamp);
    if (timestamp == null) return null;

    // Check if session is less than 24 hours old
    final age = DateTime.now().millisecondsSinceEpoch - timestamp;
    if (age > 24 * 60 * 60 * 1000) {
      await clearSession();
      return null;
    }

    final json = prefs.getString(_keySessionState);
    if (json == null) return null;

    return jsonDecode(json);
  }

  // Clear session state
  static Future<void> clearSession() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_keySessionState);
    await prefs.remove(_keyTimestamp);
  }
}
```

### Security Requirements
- **No sensitive data**: App stores no personal information
- **Local storage only**: Database remains on device, no network transmission
- **Standard Android permissions**: No special permissions required

### Performance Requirements
- **App launch**: <1 second on modern devices (Android 11+)
- **Database copy**: <2 seconds on first launch
- **Question generation**: <100ms for generating 50 questions
- **UI rendering**: 60fps during navigation and animations
- **Memory usage**: <100MB RAM during practice session

## Non-Goals (Out of Scope)

1. **iOS version**: Android-only for initial version
2. **Cloud sync**: No backend, no user accounts, no syncing
3. **Statistics/analytics**: No tracking of performance over time
4. **Custom database**: Users cannot add their own nouns or templates
5. **Social features**: No sharing, leaderboards, or multiplayer
6. **Offline download**: Database is always bundled, no downloading
7. **Text-to-speech**: No audio features
8. **Customization**: No themes, color schemes, or font size adjustments
9. **Advanced spaced repetition**: Simple practice only, no SRS algorithm
10. **Export/import**: No data export or import functionality
11. **Widget support**: No Android home screen widgets
12. **Tablet optimization**: Phone-sized UI only
13. **Accessibility features**: Basic support only (no screen reader optimization)
14. **Offline explanations AI**: Explanations generated from templates, not AI

## Testing Requirements

### Manual Testing Checklist
Since this is a personal app, manual testing only:

#### Setup Screen
- [ ] Can select each difficulty level
- [ ] Can toggle plural checkbox
- [ ] Can select each question count option
- [ ] Start button navigates to practice

#### Practice Flow
- [ ] Questions display correctly (English + Greek)
- [ ] Greek text renders properly (accents, special characters)
- [ ] Show Answer button reveals answer
- [ ] I Got It Right/Wrong buttons appear after reveal
- [ ] Explanation displays after marking answer
- [ ] Next Question button advances to next question
- [ ] Progress counter updates correctly

#### Session Types
- [ ] 10 question session completes correctly
- [ ] 20 question session completes correctly
- [ ] 50 question session completes correctly
- [ ] Endless mode continues generating questions

#### Results Screen
- [ ] Correct count matches actual correct answers
- [ ] Percentage calculated correctly
- [ ] Back to Setup button returns to setup

#### Session Persistence
- [ ] Minimizing app and returning resumes session
- [ ] Closing app and reopening within 24 hours resumes session
- [ ] Reopening after 24 hours starts fresh

#### Database
- [ ] First launch copies database successfully
- [ ] Subsequent launches use existing database
- [ ] All difficulty levels generate sentences
- [ ] Both singular and plural options work

### Test Devices
- Test on 2-3 physical devices with different screen sizes
- Minimum: One device with Android 11+
- Recommended: One phone, one tablet (if available)

## Success Metrics

### Functional Success Criteria
1. **✅ App builds and installs** on Android 11+ device
2. **✅ All three difficulty levels work** (Beginner, Intermediate, Advanced)
3. **✅ Questions generate correctly** with proper Greek text rendering
4. **✅ Explanations display** after each answer
5. **✅ Session persistence works** (resume after minimizing)
6. **✅ All session lengths work** (10, 20, 50, Endless)
7. **✅ Results screen shows accurate score**

### Technical Success Criteria
1. **✅ Database bundled** successfully (244KB in assets)
2. **✅ No crashes** during 50-question session
3. **✅ App size under 30MB**
4. **✅ UI renders smoothly** on target device (60fps)
5. **✅ Practice session can be completed** end-to-end without issues

### Personal Use Criteria
1. **✅ Convenient to use** for daily practice
2. **✅ Questions are varied** (not too repetitive)
3. **✅ Explanations are helpful** for learning
4. **✅ No friction** in workflow (quick to start practice)

## Open Questions

1. **Database versioning**: If the CLI database is updated (more nouns/templates added), how should the app be updated? *(Recommendation: Manually replace database and rebuild app)*

2. **Error handling**: What should happen if database is corrupted or missing tables? *(Recommendation: Show error message with "Reinstall app" instruction)*

3. **Screen orientation**: Should app support landscape mode or portrait-only? *(Recommendation: Portrait-only for simplicity)*

4. **Back button behavior**: Should back button during practice exit immediately or show confirmation? *(Recommendation: Exit immediately, session state is saved)*

5. **Endless mode exit**: How should user exit endless mode? *(Recommendation: Back button or explicit "End Session" button)*

6. **Font size**: Should Greek text size be fixed or adapt to system settings? *(Recommendation: Respect system font scaling)*

7. **Dark mode**: Should app support dark theme? *(Recommendation: No for initial version, follow system theme later if needed)*
