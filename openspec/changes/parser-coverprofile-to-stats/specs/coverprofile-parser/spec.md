## ADDED Requirements

### Requirement: Parse coverprofile file format
The system SHALL parse Go test coverprofile files in the standard format produced by `go test -coverprofile=<file>`. Each line follows the pattern: `<filename>:<startLine>.<startCol>,<endLine>.<endCol> <numStmt> <count>`.

#### Scenario: Parse valid coverprofile file
- **WHEN** a valid coverprofile file is provided to the parser
- **THEN** the parser successfully reads all coverage records without error

#### Scenario: Handle file not found error
- **WHEN** the coverprofile file path does not exist
- **THEN** the parser returns an error indicating file not found

#### Scenario: Handle malformed line format
- **WHEN** a line in the coverprofile does not match the expected format
- **THEN** the parser returns an error with the problematic line number and content

### Requirement: Extract file-level coverage data
The system SHALL extract and organize coverage information by source file, tracking which statements were covered and which were not.

#### Scenario: Group statements by file
- **WHEN** parsing a coverprofile with statements from multiple files
- **THEN** the parser organizes statements into file-level records with unique filename keys

#### Scenario: Track coverage status per statement
- **WHEN** parsing coverage data
- **THEN** for each statement range (startLine:startCol to endLine:endCol), the parser records whether it was covered (count > 0) or uncovered (count = 0)

### Requirement: Extract function-level coverage data
The system SHALL extract function-level coverage information from the coverprofile data. Function ranges are identified by analyzing statement ranges within files.

#### Scenario: Identify function boundaries
- **WHEN** parsing coverprofile statements with overlapping ranges
- **THEN** the parser groups statements by logical function boundaries (statements within close line ranges belong to same function)

#### Scenario: Handle functions with no coverage
- **WHEN** a function exists but has uncovered statements
- **THEN** the parser records it with coverage count of 0

### Requirement: Return structured data
The system SHALL return parsed coverage data as Go structs with clear semantics for files, functions, and statements.

#### Scenario: Return file coverage records
- **WHEN** parsing completes successfully
- **THEN** the parser returns a slice of structs with fields: filename, total statements, covered statements, coverage percentage

#### Scenario: Return statement-level details
- **WHEN** file coverage is requested with details
- **THEN** for each file, the parser provides statement records with: startLine, startCol, endLine, endCol, covered (bool)

### Requirement: Support reading from file path
The system SHALL provide a function to read and parse a coverprofile from a file path.

#### Scenario: Parse from file path
- **WHEN** ParseFile(filename string) is called with a valid file path
- **THEN** the function opens, reads, and parses the file, returning coverage data or error

### Requirement: Support reading from io.Reader
The system SHALL provide a function to parse coverprofile data from an io.Reader for flexibility (stdin, bytes.Buffer, etc.).

#### Scenario: Parse from reader
- **WHEN** Parse(reader io.Reader) is called with a valid reader
- **THEN** the function reads and parses the coverprofile, returning coverage data or error
