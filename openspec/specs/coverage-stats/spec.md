## ADDED Requirements

### Requirement: Calculate file-level coverage percentage
The system SHALL calculate coverage percentage for each file as: (covered statements / total statements) × 100.

#### Scenario: Calculate percentage for single file
- **WHEN** a file has 80 covered statements out of 100 total
- **THEN** the coverage percentage is calculated as 80.0%

#### Scenario: Handle zero-statement files
- **WHEN** a file has zero total statements
- **THEN** the coverage percentage is handled gracefully (e.g., 0% or undefined)

### Requirement: Calculate aggregate coverage statistics
The system SHALL aggregate coverage across all files to produce total project statistics: total statements, covered statements, and overall coverage percentage.

#### Scenario: Aggregate from multiple files
- **WHEN** coverage data from multiple files is provided
- **THEN** the system sums all covered statements and all total statements, then calculates overall percentage

#### Scenario: Multiple files with varying coverage
- **WHEN** files have different coverage percentages (50%, 80%, 100%)
- **THEN** the aggregated statistics reflect correct total coverage across all files

### Requirement: Track coverage per file
The system SHALL maintain per-file coverage statistics for detailed analysis.

#### Scenario: Provide file list with statistics
- **WHEN** stats are generated
- **THEN** each file record includes: filename, total statements, covered statements, coverage percentage

#### Scenario: Sort files by coverage
- **WHEN** requesting file list
- **THEN** the system provides a method to sort files by coverage percentage (ascending/descending)

### Requirement: Generate summary statistics
The system SHALL generate a summary object with overall project metrics.

#### Scenario: Return summary object
- **WHEN** stats.Summary() is called
- **THEN** the function returns an object with: total covered statements, total statements, coverage percentage, number of files analyzed

#### Scenario: Summary reflects current data
- **WHEN** the coverage data is updated or modified
- **THEN** the summary statistics reflect the latest data

### Requirement: Support incremental statistics updates
The system SHALL allow adding new coverage data to existing statistics and recalculating aggregates.

#### Scenario: Add file coverage to existing stats
- **WHEN** a new file's coverage is added to an existing stats object
- **THEN** aggregated statistics are recalculated to include the new file

#### Scenario: Merge multiple coverprofiles
- **WHEN** coverage from multiple test runs is provided
- **THEN** the system can merge them into unified statistics

### Requirement: Provide coverage breakdown by ranges
The system SHALL allow querying coverage statistics for ranges of code (e.g., "coverage of lines 10-50").

#### Scenario: Query coverage for line range
- **WHEN** a line range is specified (e.g., lines 20-40)
- **THEN** the system returns coverage statistics for statements within that range

### Requirement: Return statistics as structured Go types
The system SHALL provide clear Go types for accessing statistics data.

#### Scenario: Access file statistics
- **WHEN** requesting file statistics
- **THEN** the system returns slice of FileStats with: Filename, TotalStatements, CoveredStatements, CoveragePercent

#### Scenario: Access summary statistics
- **WHEN** requesting summary
- **THEN** the system returns SummaryStats with: TotalStatements, CoveredStatements, CoveragePercent, FileCount
