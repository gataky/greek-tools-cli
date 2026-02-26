import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import '../models/sentence.dart';
import '../models/session_config.dart';

/// SessionStorage handles saving and loading practice session state
/// Allows users to resume sessions after minimizing/closing app
class SessionStorage {
  static const String _keySessionState = 'session_state';
  static const String _keyTimestamp = 'session_timestamp';
  static const int _expirationHours = 24;

  /// Save current practice session state
  static Future<void> saveSession({
    required int currentIndex,
    required int correctCount,
    required int incorrectCount,
    required List<Sentence> sentences,
    required SessionConfig config,
  }) async {
    final prefs = await SharedPreferences.getInstance();

    final sessionData = {
      'current_index': currentIndex,
      'correct_count': correctCount,
      'incorrect_count': incorrectCount,
      'sentences': sentences.map((s) => s.toJson()).toList(),
      'config': config.toJson(),
    };

    await prefs.setString(_keySessionState, jsonEncode(sessionData));
    await prefs.setInt(_keyTimestamp, DateTime.now().millisecondsSinceEpoch);
  }

  /// Load saved practice session (returns null if expired or not found)
  static Future<Map<String, dynamic>?> loadSession() async {
    final prefs = await SharedPreferences.getInstance();

    // Check timestamp
    final timestamp = prefs.getInt(_keyTimestamp);
    if (timestamp == null) {
      return null;
    }

    // Check if session is less than 24 hours old
    final age = DateTime.now().millisecondsSinceEpoch - timestamp;
    final maxAge = _expirationHours * 60 * 60 * 1000;

    if (age > maxAge) {
      // Session expired, clear it
      await clearSession();
      return null;
    }

    // Load session data
    final json = prefs.getString(_keySessionState);
    if (json == null) {
      return null;
    }

    try {
      return jsonDecode(json) as Map<String, dynamic>;
    } catch (e) {
      // Invalid JSON, clear it
      await clearSession();
      return null;
    }
  }

  /// Clear saved session state
  static Future<void> clearSession() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_keySessionState);
    await prefs.remove(_keyTimestamp);
  }

  /// Check if a valid session exists
  static Future<bool> hasValidSession() async {
    final session = await loadSession();
    return session != null;
  }
}
