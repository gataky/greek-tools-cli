# Greek Case Master

Greek Case Master is a command-line educational tool designed to help English speakers learning Modern Greek transition from rote memorization of noun declension tables to instinctive application of cases (Nominative, Genitive, Accusative) in real sentence contexts.

The application presents English sentence prompts with missing Greek nouns. Users must provide the correctly declined Greek article + noun combination. The system validates answers with exact matching and provides detailed grammar explanations (translation, syntactic role, morphology) for every answer.

## Features

- **AI-Powered Generation**: Uses Claude API to generate accurate declensions and natural practice sentences for any Greek noun.
- **Interactive TUI**: A modern terminal interface for practice sessions with instant feedback.
- **Offline Practice**: Once imported, all nouns and sentences are stored in a local SQLite database for offline study.
- **Grammar Explanations**: Every answer is followed by a detailed breakdown of the grammar, helping you understand *why* a specific case is used.
- **Progressive Difficulty**: Choose between Beginner (Accusative focus), Intermediate (Genitive focus), and Advanced (Mixed cases with prepositions) levels.
- **Resume Capability**: Large imports can be resumed if interrupted, ensuring no progress (or API credits) are lost.

## Prerequisites

- **Go 1.21+**
- **Anthropic API Key**: Required for importing new nouns or adding nouns manually.
- **SQLite 3**: Used for local data storage.

## Installation

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/gataky/greekmaster.git
   cd greekmaster
   ```

2. Build the application using the Makefile:
   ```bash
   make build
   ```

3. (Optional) Install the binary to your `$GOPATH/bin`:
   ```bash
   make install
   ```

## Getting Started

### 1. Set up your API Key

Export your Anthropic API key as an environment variable:

```bash
export ANTHROPIC_API_KEY='your-api-key-here'
```

### 2. Import Nouns

Create a CSV file (e.g., `nouns.csv`) with the nouns you want to practice. The format should be: `english,greek,attribute` (where attribute is the gender: masculine, feminine, neuter, or invariable).

```csv
english,greek,attribute
teacher,δάσκαλος,masculine
book,βιβλίο,neuter
woman,γυναίκα,feminine
```

Run the import command:

```bash
./greekmaster import nouns.csv
```

*Note: This process uses the Claude API to generate practice data and may take a few minutes depending on the number of nouns.*

### 3. Start Practicing

Once you have imported some nouns, start an interactive practice session:

```bash
./greekmaster practice
```

Follow the on-screen prompts to select your difficulty, session length, and whether to include plural forms.

## Usage

### Commands

- `import <csv-file>`: Import nouns from a CSV and generate practice data.
- `practice`: Start an interactive TUI practice session.
- `add`: Interactively add a single noun with AI-generated data.
- `list`: List all nouns currently in the database.
- `--help`: Show help for any command.

### Global Flags

- `--db-path <path>`: Specify a custom path for the SQLite database (default: `~/.greekmaster/greekmaster.db`).

## Development

### Makefile Commands

- `make build`: Compiles the binary to the root directory.
- `make install`: Installs the binary to `$GOPATH/bin`.
- `make test`: Runs all unit tests.
- `make clean`: Removes the compiled binary and build artifacts.
- `make run`: Builds and runs the application.

## License

[MIT License](LICENSE) (Replace with your actual license if different)
