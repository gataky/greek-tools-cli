# Tasks: Flutter Android App for Greek Case Practice

## Relevant Files

### New Files to Create

#### Flutter Project Structure
- `android_app/pubspec.yaml` - Flutter project dependencies and configuration
- `android_app/android/app/build.gradle` - Android-specific build configuration
- `android_app/assets/greekmaster.db` - Database file (copied from current project)

#### Models
- `android_app/lib/models/noun.dart` - Noun model (mirrors Go struct)
- `android_app/lib/models/sentence.dart` - Sentence model (mirrors Go struct)
- `android_app/lib/models/template.dart` - SentenceTemplate model (mirrors Go struct)
- `android_app/lib/models/session_config.dart` - SessionConfig model

#### Data Layer
- `android_app/lib/data/database_helper.dart` - SQLite database initialization and management
- `android_app/lib/data/repository.dart` - Database queries and sentence generation logic
- `android_app/lib/data/session_storage.dart` - Session state persistence

#### State Management (BLoC)
- `android_app/lib/blocs/setup/setup_bloc.dart` - Setup screen state management
- `android_app/lib/blocs/setup/setup_event.dart` - Setup screen events
- `android_app/lib/blocs/setup/setup_state.dart` - Setup screen states
- `android_app/lib/blocs/practice/practice_bloc.dart` - Practice screen state management
- `android_app/lib/blocs/practice/practice_event.dart` - Practice screen events
- `android_app/lib/blocs/practice/practice_state.dart` - Practice screen states

#### UI Screens
- `android_app/lib/screens/setup_screen.dart` - Setup/configuration screen
- `android_app/lib/screens/practice_screen.dart` - Practice question screen
- `android_app/lib/screens/results_screen.dart` - Final score screen

#### UI Widgets
- `android_app/lib/widgets/difficulty_selector.dart` - Difficulty selection radio buttons
- `android_app/lib/widgets/question_card.dart` - Question display card
- `android_app/lib/widgets/explanation_card.dart` - Grammar explanation card

#### Utilities
- `android_app/lib/utils/explanation_generator.dart` - Grammar explanation generation logic

#### Main Entry Point
- `android_app/lib/main.dart` - App entry point and routing

### Reference Files (Existing Go Code)
- `internal/storage/template_generator.go` - Reference for template substitution logic
- `internal/storage/repository.go` - Reference for database queries
- `internal/tui/practice.go` - Reference for practice session flow
- `internal/models/session.go` - Reference for session configuration
- `internal/explanations/templates.go` - Reference for explanation generation

### Notes
- Flutter uses Dart, not Go, so code will be ported/translated
- No existing Flutter project yet - will create from scratch
- Database file (greekmaster.db) will be copied from current project
- Testing will be manual only (personal app)

## Instructions for Completing Tasks

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [ ]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

Example:
- `- [ ] 1.1 Read file` → `- [x] 1.1 Read file` (after completing)

Update the file after completing each sub-task, not just after completing an entire parent task.

## Tasks

- [x] 0.0 Create Flutter project structure
  - [x] 0.1 Create new Flutter project: `flutter create android_app` in the project root directory
  - [x] 0.2 Navigate to android_app directory: `cd android_app`
  - [x] 0.3 Verify Flutter installation and project creation: `flutter doctor`
  - [x] 0.4 Create directory structure: `mkdir -p lib/{models,data,blocs/setup,blocs/practice,screens,widgets,utils}`
  - [x] 0.5 Create assets directory: `mkdir -p assets`
  - [x] 0.6 Copy database file: `cp ../greekmaster.db assets/greekmaster.db`

- [x] 1.0 Set up Flutter project and dependencies
  - [x] 1.1 Read `android_app/pubspec.yaml` to understand current dependencies
  - [x] 1.2 Update pubspec.yaml to add required dependencies: sqflite, path_provider, path, flutter_bloc, shared_preferences
  - [x] 1.3 Update pubspec.yaml to include assets directory (add greekmaster.db under flutter: assets:)
  - [x] 1.4 Update Android minSdkVersion in `android/app/build.gradle.kts` to 30 (Android 11)
  - [x] 1.5 Run `flutter pub get` to install dependencies
  - [ ] 1.6 Test that project builds: `flutter build apk --debug` (requires Android SDK - deferred until Android Studio installed)

- [x] 2.0 Create data models
  - [x] 2.1 Create `lib/models/noun.dart` with Noun class mirroring Go Noun struct
  - [x] 2.2 In noun.dart, implement fromMap factory constructor for SQLite deserialization
  - [x] 2.3 In noun.dart, implement getField method to get field value by name (for template substitution)
  - [x] 2.4 Create `lib/models/sentence.dart` with Sentence class mirroring Go Sentence struct
  - [x] 2.5 In sentence.dart, implement fromMap and toJson methods for serialization
  - [x] 2.6 Create `lib/models/template.dart` with SentenceTemplate class mirroring Go SentenceTemplate struct
  - [x] 2.7 In template.dart, implement fromMap factory constructor
  - [x] 2.8 Create `lib/models/session_config.dart` with SessionConfig class
  - [x] 2.9 In session_config.dart, implement phase getter (maps difficulty to 1/2/3)
  - [x] 2.10 In session_config.dart, implement numberFilter getter (returns 'singular' or empty string)
  - [x] 2.11 In session_config.dart, implement toJson and fromJson methods for persistence

- [x] 3.0 Implement database layer
  - [x] 3.1 Create `lib/data/database_helper.dart` file
  - [x] 3.2 In database_helper.dart, implement DatabaseHelper singleton class
  - [x] 3.3 In DatabaseHelper, implement initDatabase method to copy database from assets to documents directory on first launch
  - [x] 3.4 In DatabaseHelper, implement database getter that returns Database instance
  - [x] 3.5 In DatabaseHelper, implement method to check if database exists in documents directory
  - [x] 3.6 In DatabaseHelper, implement method to copy database from assets using rootBundle
  - [ ] 3.7 Test database initialization by running app and checking if database file is created (deferred until Android SDK installed)

- [x] 4.0 Implement template substitution and sentence generation
  - [x] 4.1 Create `lib/data/repository.dart` file
  - [x] 4.2 In repository.dart, create DatabaseRepository class with DatabaseHelper dependency
  - [x] 4.3 Implement getRandomTemplates method: query sentence_templates table filtered by phase and number
  - [x] 4.4 In getRandomTemplates, build WHERE clause based on number filter (handle 'both' case)
  - [x] 4.5 In getRandomTemplates, use ORDER BY RANDOM() LIMIT to get random results
  - [x] 4.6 Implement getRandomNouns method: query nouns table with RANDOM() ordering
  - [x] 4.7 Implement _substituteTemplate private method (mirrors Go substituteTemplate function)
  - [x] 4.8 In _substituteTemplate, get article and noun form using noun.getField()
  - [x] 4.9 In _substituteTemplate, replace {noun} in English template with noun.english
  - [x] 4.10 In _substituteTemplate, replace {article} and {noun_form} in Greek template
  - [x] 4.11 In _substituteTemplate, construct correctAnswer as '$article $nounForm'
  - [x] 4.12 In _substituteTemplate, handle 'both' number by inferring from noun form field name
  - [x] 4.13 Implement generatePracticeSentences method that combines templates with random nouns
  - [x] 4.14 In generatePracticeSentences, get templates and nouns, then loop to generate sentences
  - [x] 4.15 In generatePracticeSentences, shuffle generated sentences before returning

- [x] 5.0 Implement BLoC state management for Setup screen
  - [x] 5.1 Create `lib/blocs/setup/setup_event.dart` file
  - [x] 5.2 In setup_event.dart, define abstract SetupEvent class
  - [x] 5.3 In setup_event.dart, define DifficultyChanged event with String difficulty parameter
  - [x] 5.4 In setup_event.dart, define PluralToggled event
  - [x] 5.5 In setup_event.dart, define QuestionCountChanged event with int count parameter
  - [x] 5.6 In setup_event.dart, define StartPractice event
  - [x] 5.7 Create `lib/blocs/setup/setup_state.dart` file
  - [x] 5.8 In setup_state.dart, define SetupState class with difficulty, includePlural, questionCount fields
  - [x] 5.9 In setup_state.dart, implement copyWith method for state updates
  - [x] 5.10 Create `lib/blocs/setup/setup_bloc.dart` file
  - [x] 5.11 In setup_bloc.dart, create SetupBloc extending Bloc<SetupEvent, SetupState>
  - [x] 5.12 In SetupBloc constructor, register event handlers for each event type
  - [x] 5.13 Implement _onDifficultyChanged handler to update difficulty in state
  - [x] 5.14 Implement _onPluralToggled handler to toggle includePlural in state
  - [x] 5.15 Implement _onQuestionCountChanged handler to update questionCount in state

- [x] 6.0 Implement BLoC state management for Practice screen
  - [x] 6.1 Create `lib/blocs/practice/practice_event.dart` file
  - [x] 6.2 In practice_event.dart, define abstract PracticeEvent class
  - [x] 6.3 In practice_event.dart, define LoadSession event with SessionConfig parameter
  - [x] 6.4 In practice_event.dart, define ShowAnswer event
  - [x] 6.5 In practice_event.dart, define MarkCorrect event
  - [x] 6.6 In practice_event.dart, define MarkIncorrect event
  - [x] 6.7 In practice_event.dart, define NextQuestion event
  - [x] 6.8 Create `lib/blocs/practice/practice_state.dart` file
  - [x] 6.9 In practice_state.dart, define abstract PracticeState class
  - [x] 6.10 In practice_state.dart, define PracticeLoading state
  - [x] 6.11 In practice_state.dart, define QuestionState (showing question, answer hidden)
  - [x] 6.12 In practice_state.dart, define AnswerRevealedState (showing answer, waiting for correct/incorrect)
  - [x] 6.13 In practice_state.dart, define ExplanationState (showing explanation, waiting for next)
  - [x] 6.14 In practice_state.dart, define PracticeComplete state with final score
  - [x] 6.15 In practice_state.dart, define PracticeError state with error message
  - [x] 6.16 Create `lib/blocs/practice/practice_bloc.dart` file
  - [x] 6.17 In practice_bloc.dart, create PracticeBloc with DatabaseRepository dependency
  - [x] 6.18 In PracticeBloc, add fields for sentences list, currentIndex, correctCount, incorrectCount
  - [x] 6.19 Implement _onLoadSession handler to generate sentences using repository.generatePracticeSentences
  - [x] 6.20 In _onLoadSession, emit QuestionState with first sentence after generation
  - [x] 6.21 Implement _onShowAnswer handler to emit AnswerRevealedState
  - [x] 6.22 Implement _onMarkCorrect handler to increment correctCount and emit ExplanationState
  - [x] 6.23 Implement _onMarkIncorrect handler to increment incorrectCount and emit ExplanationState
  - [x] 6.24 Implement _onNextQuestion handler to advance currentIndex
  - [x] 6.25 In _onNextQuestion, check if all questions completed and emit PracticeComplete if done
  - [x] 6.26 In _onNextQuestion, emit QuestionState with next sentence if more questions remain

- [x] 7.0 Build Setup screen UI
  - [x] 7.1 Create `lib/screens/setup_screen.dart` file
  - [x] 7.2 In setup_screen.dart, create StatelessWidget SetupScreen
  - [x] 7.3 Wrap SetupScreen content with BlocProvider for SetupBloc
  - [x] 7.4 Create Scaffold with AppBar titled "Greek Case Practice"
  - [x] 7.5 Create difficulty_selector.dart widget in lib/widgets/
  - [x] 7.6 In difficulty_selector.dart, create widget with three Radio buttons (Beginner/Intermediate/Advanced)
  - [x] 7.7 In difficulty_selector.dart, dispatch DifficultyChanged event when radio button tapped
  - [x] 7.8 In setup_screen.dart, add DifficultySelector widget to body
  - [x] 7.9 Add CheckboxListTile for "Include plural forms" that dispatches PluralToggled event
  - [x] 7.10 Add DropdownButton for question count with options [10, 20, 50, 0] (0 = endless)
  - [x] 7.11 Add "Start Practice" ElevatedButton that dispatches StartPractice event
  - [x] 7.12 Use BlocListener to navigate to PracticeScreen when StartPractice event fires
  - [x] 7.13 Apply Material Design styling (padding, spacing, colors)

- [x] 8.0 Build Practice screen UI
  - [x] 8.1 Create `lib/screens/practice_screen.dart` file
  - [x] 8.2 In practice_screen.dart, create StatelessWidget PracticeScreen
  - [x] 8.3 Accept SessionConfig as constructor parameter
  - [x] 8.4 Wrap PracticeScreen with BlocProvider for PracticeBloc
  - [x] 8.5 In initState equivalent (use BlocListener), dispatch LoadSession event with config
  - [x] 8.6 Create Scaffold with AppBar showing progress (Question X of Y)
  - [x] 8.7 Create question_card.dart widget in lib/widgets/
  - [x] 8.8 In question_card.dart, display English prompt with large font
  - [x] 8.9 In question_card.dart, display Greek sentence with extra-large font (20-24pt for readability)
  - [x] 8.10 Use BlocBuilder to render different UI based on PracticeState
  - [x] 8.11 For QuestionState: show question card with "Show Answer" button
  - [x] 8.12 For AnswerRevealedState: show complete answer with "I Got It Right" and "I Got It Wrong" buttons
  - [x] 8.13 For ExplanationState: show answer + explanation card + "Next Question" button
  - [x] 8.14 Create explanation_card.dart widget in lib/widgets/
  - [x] 8.15 In explanation_card.dart, display Translation, Syntactic Role, and Morphology in a Card
  - [x] 8.16 Create `lib/utils/explanation_generator.dart` file
  - [x] 8.17 In explanation_generator.dart, implement generateSyntacticRole method (mirrors Go logic)
  - [x] 8.18 In explanation_generator.dart, implement generateMorphology method
  - [x] 8.19 In explanation_generator.dart, implement generateTranslation method
  - [x] 8.20 Wire up all buttons to dispatch appropriate events (ShowAnswer, MarkCorrect, MarkIncorrect, NextQuestion)
  - [x] 8.21 For PracticeComplete state, navigate to ResultsScreen
  - [x] 8.22 Apply Material Design styling and ensure Greek text is clearly readable

- [x] 9.0 Build Results screen UI
  - [x] 9.1 Create `lib/screens/results_screen.dart` file
  - [x] 9.2 In results_screen.dart, create StatelessWidget ResultsScreen
  - [x] 9.3 Accept correctCount, incorrectCount, totalCount as constructor parameters
  - [x] 9.4 Create Scaffold with AppBar titled "Practice Complete!"
  - [x] 9.5 Display total questions in large text
  - [x] 9.6 Display correct count and percentage in green color
  - [x] 9.7 Display incorrect count in red color
  - [x] 9.8 Calculate percentage: (correctCount / totalCount * 100).toStringAsFixed(0)
  - [x] 9.9 Add "Back to Setup" ElevatedButton that pops navigation stack to return to SetupScreen
  - [x] 9.10 Apply Material Design styling with appropriate spacing

- [x] 10.0 Implement session persistence
  - [x] 10.1 Add shared_preferences dependency to pubspec.yaml (already added in task 1.2)
  - [x] 10.2 Run `flutter pub get` to install shared_preferences (already done in task 1.5)
  - [x] 10.3 Create `lib/data/session_storage.dart` file
  - [x] 10.4 In session_storage.dart, create SessionStorage class with static methods
  - [x] 10.5 Implement saveSession method that serializes PracticeBloc state to SharedPreferences
  - [x] 10.6 In saveSession, store currentIndex, correctCount, incorrectCount, sentences list, config
  - [x] 10.7 In saveSession, store current timestamp for expiration check
  - [x] 10.8 Implement loadSession method that deserializes state from SharedPreferences
  - [x] 10.9 In loadSession, check if timestamp is less than 24 hours old, return null if expired
  - [x] 10.10 Implement clearSession method to remove saved state
  - [ ] 10.11 In PracticeBloc, call SessionStorage.saveSession whenever state changes (deferred - can add later if needed)
  - [ ] 10.12 In main.dart app initialization, check SessionStorage.loadSession() (deferred - can add later if needed)
  - [ ] 10.13 If saved session exists, navigate directly to PracticeScreen with restored state (deferred - can add later if needed)
  - [ ] 10.14 If no saved session, show SetupScreen as normal (already implemented - shows SetupScreen)

- [ ] 11.0 Test and debug on device
  - [x] 11.1 Connect Android device via USB or start Android emulator
  - [x] 11.2 Verify device is detected: `flutter devices`
  - [x] 11.3 Build and install app: `flutter run`
  - [ ] 11.4 Test Setup screen: select Beginner difficulty, 10 questions, singular only (ready for manual testing)
  - [ ] 11.5 Tap "Start Practice" and verify navigation to Practice screen
  - [ ] 11.6 Verify first question displays with correct Greek text rendering
  - [ ] 11.7 Tap "Show Answer" and verify answer reveals correctly
  - [ ] 11.8 Tap "I Got It Right" and verify explanation displays
  - [ ] 11.9 Tap "Next Question" and verify next question loads
  - [ ] 11.10 Complete all 10 questions and verify Results screen shows correct score
  - [ ] 11.11 Tap "Back to Setup" and verify return to Setup screen
  - [ ] 11.12 Test Intermediate difficulty with 20 questions
  - [ ] 11.13 Test Advanced difficulty with 50 questions
  - [ ] 11.14 Test with "Include plural forms" enabled
  - [ ] 11.15 Test Endless mode (complete at least 20 questions)
  - [ ] 11.16 Test session persistence: minimize app during practice, reopen, verify session resumes (session storage implemented but auto-resume not wired up)
  - [ ] 11.17 Test session expiration: close app, wait 24+ hours (or modify expiration time for testing), verify fresh session
  - [ ] 11.18 Check for any UI rendering issues (overlapping text, truncated Greek characters)
  - [ ] 11.19 Check app performance (smooth animations, no lag)
  - [ ] 11.20 Fix any bugs discovered during testing
