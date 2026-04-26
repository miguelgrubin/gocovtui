## ADDED Requirements

### Requirement: FolderStats type
The system SHALL provide a `FolderStats` struct in `pkg/coverage` that aggregates coverage across all files sharing the same parent directory.

#### Scenario: FolderStats fields populated
- **WHEN** `Stats.FolderStats()` is called after files have been added
- **THEN** each returned `FolderStats` SHALL have `Dir`, `TotalStatements`, `CoveredStatements`, `CoveragePercent`, and `FileCount` correctly summed from all files whose `path.Dir(Filename)` matches `Dir`

#### Scenario: Coverage percentage calculation
- **WHEN** a folder contains files with varying coverage
- **THEN** `CoveragePercent` SHALL equal `CoveredStatements / TotalStatements * 100` across all files in that folder

#### Scenario: Empty stats
- **WHEN** no files have been added to `Stats`
- **THEN** `FolderStats()` SHALL return an empty slice

### Requirement: FoldersSortedByCoverage method
The `Stats` type SHALL provide a `FoldersSortedByCoverage(ascending bool) []*FolderStats` method.

#### Scenario: Ascending sort
- **WHEN** `FoldersSortedByCoverage(true)` is called
- **THEN** folders with the lowest `CoveragePercent` SHALL appear first

#### Scenario: Descending sort
- **WHEN** `FoldersSortedByCoverage(false)` is called
- **THEN** folders with the highest `CoveragePercent` SHALL appear first
