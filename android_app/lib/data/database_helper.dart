import 'dart:io';
import 'package:flutter/services.dart';
import 'package:path/path.dart';
import 'package:path_provider/path_provider.dart';
import 'package:sqflite/sqflite.dart';

/// DatabaseHelper manages SQLite database initialization and access
/// Implements singleton pattern for single database instance
class DatabaseHelper {
  static final DatabaseHelper _instance = DatabaseHelper._internal();
  static Database? _database;

  factory DatabaseHelper() {
    return _instance;
  }

  DatabaseHelper._internal();

  /// Get database instance, initializing if necessary
  Future<Database> get database async {
    if (_database != null) {
      return _database!;
    }

    _database = await initDatabase();
    return _database!;
  }

  /// Initialize database by copying from assets on first launch
  Future<Database> initDatabase() async {
    // Get app's documents directory
    final documentsDirectory = await getApplicationDocumentsDirectory();
    final path = join(documentsDirectory.path, 'greekmaster.db');

    // Check if database exists
    final exists = await databaseExists(path);

    if (!exists) {
      // Database doesn't exist, copy from assets
      print('Copying database from assets...');
      await _copyDatabaseFromAssets(path);
      print('Database copied successfully to $path');
    } else {
      print('Using existing database at $path');
    }

    // Open and return database
    return await openDatabase(
      path,
      version: 1,
      readOnly: false,
    );
  }

  /// Check if database file exists at path
  Future<bool> databaseExists(String path) async {
    return await File(path).exists();
  }

  /// Copy database from assets to documents directory
  Future<void> _copyDatabaseFromAssets(String path) async {
    try {
      // Ensure directory exists
      final directory = Directory(dirname(path));
      if (!await directory.exists()) {
        await directory.create(recursive: true);
      }

      // Read database from assets
      final data = await rootBundle.load('assets/greekmaster.db');
      final bytes = data.buffer.asUint8List();

      // Write to documents directory
      await File(path).writeAsBytes(bytes, flush: true);
    } catch (e) {
      throw Exception('Failed to copy database from assets: $e');
    }
  }

  /// Close database connection
  Future<void> close() async {
    if (_database != null) {
      await _database!.close();
      _database = null;
    }
  }
}
