## ADDED Requirements

### Requirement: User can trigger find action with 'f' key
The system SHALL activate an interactive find/search mode when the user presses 'f' while browsing files in the file picker.

#### Scenario: Find mode activated
- **WHEN** user presses 'f' key in file picker
- **THEN** find mode activates showing a search input field (e.g., "Find: ")

#### Scenario: Find mode cancelled
- **WHEN** user presses ESC or Ctrl+C in find mode
- **THEN** find mode closes and returns to normal file picker view

### Requirement: User can search files by name
The system SHALL filter the file and folder list based on substring matching against the user's search term.

#### Scenario: Filter by file name
- **WHEN** user types a search term in find mode (e.g., "test")
- **THEN** the file list updates to show only files/folders containing "test" in their name

#### Scenario: Live filtering as user types
- **WHEN** user types or modifies search term character by character
- **THEN** file list updates in real-time with matching files

#### Scenario: Case handling
- **WHEN** user searches regardless of case (e.g., "Test" vs "test")
- **THEN** search is case-insensitive and matches files like "test.go", "Test.md", "TEST_FILE"

### Requirement: User can confirm and select from filtered results
The system SHALL allow the user to navigate and select files from the filtered list.

#### Scenario: Navigate filtered results
- **WHEN** user presses arrow keys in find mode after filtering
- **THEN** selection highlights move through the filtered file list

#### Scenario: Select file from filtered results
- **WHEN** user presses Enter on a highlighted file in filtered results
- **THEN** file is selected and find mode closes, returning to normal view with the file selected

### Requirement: User can clear search to return to full file list
The system SHALL provide a way to reset the filter and show all files again.

#### Scenario: Clear search and show all
- **WHEN** user clears the search input (e.g., deletes all typed characters)
- **THEN** file list returns to showing all files in the directory

#### Scenario: Search history or editing
- **WHEN** user edits the search term (adds/removes characters)
- **THEN** filtered list updates immediately to match the new search term
