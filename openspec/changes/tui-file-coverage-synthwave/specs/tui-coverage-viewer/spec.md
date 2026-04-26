## ADDED Requirements

### Requirement: Interactive TUI program entry
The system SHALL launch a Bubble Tea interactive TUI program when invoked, replacing static stdout rendering.

#### Scenario: Program starts with coverage file
- **WHEN** the app is run with a valid coverage profile path
- **THEN** a full-screen interactive TUI SHALL be displayed with file coverage data loaded

#### Scenario: Program starts without coverage file
- **WHEN** the app is run with an empty or missing coverage profile path
- **THEN** the TUI SHALL display an empty state message indicating no coverage data

#### Scenario: User quits the program
- **WHEN** the user presses `q` or `ctrl+c`
- **THEN** the TUI SHALL exit and return control to the shell

### Requirement: Synthwave color theme
The TUI SHALL use a Synthwave-inspired color palette throughout all UI elements.

#### Scenario: Color palette applied
- **WHEN** the TUI is rendered
- **THEN** the background SHALL use dark navy (`#1a1a2e`), headings neon pink (`#ff2d78`), accents neon cyan (`#00f5d4`), borders purple (`#7b2fff`), and highlights bright yellow (`#ffe600`)

### Requirement: Summary header display
The TUI SHALL display an aggregate summary header above the file list.

#### Scenario: Summary shows aggregate stats
- **WHEN** coverage data is loaded
- **THEN** the header SHALL display total coverage percentage, covered/total statement count, and file count

### Requirement: Scrollable file coverage list
The TUI SHALL display a scrollable list of files with per-file coverage statistics.

#### Scenario: List shows all files
- **WHEN** coverage data is loaded
- **THEN** each list item SHALL display the filename, total statements, and coverage percentage

#### Scenario: Keyboard navigation
- **WHEN** the user presses arrow keys or `j`/`k`
- **THEN** the list selection SHALL move up or down respectively, scrolling as needed

#### Scenario: Coverage percentage color coding
- **WHEN** a file's coverage is ≥ 80%
- **THEN** the coverage value SHALL render in neon cyan (`#00f5d4`)
- **WHEN** a file's coverage is between 50% and 79%
- **THEN** the coverage value SHALL render in bright yellow (`#ffe600`)
- **WHEN** a file's coverage is below 50%
- **THEN** the coverage value SHALL render in neon pink (`#ff2d78`)

### Requirement: Terminal resize handling
The TUI SHALL adapt its layout when the terminal window is resized.

#### Scenario: Window resize event
- **WHEN** the terminal is resized
- **THEN** the list and header SHALL adjust their width and height to fill the new terminal dimensions
