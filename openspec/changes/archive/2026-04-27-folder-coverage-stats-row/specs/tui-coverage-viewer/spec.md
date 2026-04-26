## ADDED Requirements

### Requirement: Folder summary rows in file list
The TUI file list SHALL display a folder summary row above each group of files that belong to the same directory.

#### Scenario: Folder row appears before its files
- **WHEN** the TUI list is rendered with coverage data containing files from multiple directories
- **THEN** each directory SHALL have one folder row immediately preceding its file rows

#### Scenario: Folder row visual distinction
- **WHEN** a folder row is rendered
- **THEN** it SHALL use a distinct bold style with a folder-specific Synthwave color (not the same as normal file rows)

#### Scenario: Folder row content
- **WHEN** a folder row is rendered
- **THEN** it SHALL display the folder/package path and its aggregate coverage percentage

#### Scenario: Grouping order
- **WHEN** the list is built
- **THEN** folders SHALL be sorted by coverage ascending (lowest first), and files within each folder SHALL also be sorted by coverage ascending

### Requirement: Folder rows are navigable
The folder rows SHALL be reachable via keyboard navigation (arrow keys / j/k) but SHALL NOT trigger any action when selected.

#### Scenario: Cursor lands on folder row
- **WHEN** the user navigates the list and the cursor lands on a folder row
- **THEN** the folder row SHALL be highlighted like any other selected item but no action SHALL occur
