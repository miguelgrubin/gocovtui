## 1. Project Setup

- [x] 1.1 Create `pkg/coverage/` package directory structure
- [x] 1.2 Create `parser.go` with coverprofile parsing types and signatures
- [x] 1.3 Create `stats.go` with statistics calculation types and signatures

## 2. Coverprofile Parser Implementation

- [x] 2.1 Implement `ParseCoverprofile(data string)` function to parse coverprofile format
- [x] 2.2 Implement `FileCoverage` struct to store file-level coverage data
- [x] 2.3 Implement `StatementCoverage` struct to store individual statement coverage
- [x] 2.4 Add line parsing logic to extract filename, line/col ranges, statement count, execution count
- [x] 2.5 Add error handling for malformed coverprofile lines
- [x] 2.6 Implement `ParseFile(path string)` function to read from file path
- [x] 2.7 Implement `Parse(reader io.Reader)` function to read from io.Reader

## 3. Function-Level Coverage Extraction

- [x] 3.1 Implement function boundary detection logic (grouping statements by proximity)
- [x] 3.2 Implement `FunctionCoverage` struct to represent function-level coverage
- [x] 3.3 Create `ExtractFunctions()` method to identify functions from statement ranges
- [x] 3.4 Handle edge cases (nested functions, closures)

## 4. Coverage Statistics Implementation

- [x] 4.1 Implement `FileStats` struct with coverage percentage calculation
- [x] 4.2 Implement `SummaryStats` struct for aggregate project statistics
- [x] 4.3 Implement percentage calculation logic (covered/total × 100)
- [x] 4.4 Implement `CalculateStats()` function to aggregate all file statistics
- [x] 4.5 Implement `GetSummary()` function to return project-wide stats
- [x] 4.6 Add support for merging coverage from multiple files
- [x] 4.7 Implement range-based coverage queries

## 5. Unit Tests

- [x] 5.1 Create `parser_test.go` with coverprofile parsing tests
- [x] 5.2 Write tests for valid coverprofile format
- [x] 5.3 Write tests for malformed coverprofile handling
- [x] 5.4 Write tests for file not found errors
- [x] 5.5 Create `stats_test.go` with statistics calculation tests
- [x] 5.6 Write tests for percentage calculation (edge cases: 0 statements, partial coverage)
- [x] 5.7 Write tests for aggregation across multiple files
- [x] 5.8 Write tests for merging coverage data

## 6. Integration & Documentation

- [x] 6.1 Run existing test suite to ensure no regressions (`make test`)
- [x] 6.2 Run coverage report generation (`make coverage`) to validate with real data
- [x] 6.3 Add godoc comments to exported types and functions
- [x] 6.4 Create example usage in package documentation
- [x] 6.5 Verify code passes linting (`make lint`)
- [x] 6.6 Verify code passes security scan (`make sec`)

## 7. Integration with App

- [x] 7.1 Update `pkg/app.go` to optionally load coverage data
- [x] 7.2 Add example in TUI to display loaded coverage statistics
- [x] 7.3 Create end-to-end test: generate coverage, parse, display in TUI
