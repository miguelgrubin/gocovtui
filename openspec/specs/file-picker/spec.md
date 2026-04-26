## ADDED Requirements

### Requirement: Show file picker when no argument is given
When `gocovtui` is launched without any CLI arguments, the program SHALL display an interactive TUI file picker listing all `*.out` files in the current working directory before starting the main coverage viewer.

#### Scenario: Launch without argument
- **WHEN** the user runs `gocovtui` with no arguments
- **THEN** an interactive file picker is shown listing all `*.out` files in the CWD

#### Scenario: User selects a file
- **WHEN** the user navigates the picker and confirms a selection
- **THEN** the file picker closes and the main coverage TUI opens with the selected file path

#### Scenario: User cancels the picker
- **WHEN** the user presses `q` or `ctrl+c` in the file picker without selecting a file
- **THEN** the program exits cleanly with no error output

### Requirement: Filter picker to only show `.out` files
The file picker SHALL display only files with the `.out` extension. Other files and directories SHALL NOT appear in the list.

#### Scenario: CWD contains mixed files
- **WHEN** the CWD contains `coverage.out`, `main.go`, and `go.mod`
- **THEN** only `coverage.out` is shown in the picker

### Requirement: Handle empty file list
When no `*.out` files exist in the CWD, the program SHALL display an informative message and exit with a non-zero status code.

#### Scenario: No `.out` files present
- **WHEN** the user runs `gocovtui` with no arguments and there are no `*.out` files in the CWD
- **THEN** the program prints an error message to stderr and exits with a non-zero status code

### Requirement: Argument bypasses file picker
When `gocovtui` is launched with a file path argument, the file picker SHALL NOT be shown and existing behaviour SHALL be preserved.

#### Scenario: Launch with argument
- **WHEN** the user runs `gocovtui coverage.out`
- **THEN** the file picker is skipped and the main TUI opens directly with `coverage.out`
