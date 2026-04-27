## ADDED Requirements

### Requirement: User can trigger edit action with 'e' key
The system SHALL open the selected file in the user's configured external editor when the user presses 'e' key on a file.

#### Scenario: Edit action triggered with file selected
- **WHEN** user selects a file and presses 'e' key
- **THEN** the selected file opens in the configured `$EDITOR` environment variable

#### Scenario: Edit mode unavailable without file selection
- **WHEN** user presses 'e' key without a file selected or while on a directory
- **THEN** system shows an appropriate message indicating a file must be selected

### Requirement: Editor process is spawned with correct file path
The system SHALL pass the absolute or relative file path to the editor command.

#### Scenario: Editor receives correct file path
- **WHEN** user presses 'e' on a file (e.g., "coverage.txt")
- **THEN** the `$EDITOR` command is executed with the file path as an argument

#### Scenario: Editor inherits TUI's environment
- **WHEN** editor is spawned
- **THEN** the editor process inherits the current shell environment (including `$EDITOR` value)

### Requirement: TUI pauses while editor is active
The system SHALL suspend the TUI interface and allow the user to interact with the external editor.

#### Scenario: TUI pauses during editing
- **WHEN** editor is spawned and running
- **THEN** TUI becomes inactive (hidden or suspended) while user edits

#### Scenario: TUI resumes after editor closes
- **WHEN** user closes the editor and returns to the TUI
- **THEN** TUI resumes normal operation showing the same file picker state as before

### Requirement: Fallback behavior if EDITOR is unset
The system SHALL provide a reasonable behavior if the `$EDITOR` environment variable is not configured.

#### Scenario: EDITOR environment variable not set
- **WHEN** user presses 'e' but `$EDITOR` is not set
- **THEN** system shows error message (e.g., "EDITOR environment variable not configured") or falls back to default editor (e.g., vim)

#### Scenario: Invalid EDITOR configuration
- **WHEN** `$EDITOR` points to a non-existent command or fails to execute
- **THEN** system shows clear error message and returns to file picker
