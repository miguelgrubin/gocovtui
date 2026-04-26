## Why

Go's `go test -coverprofile` output format is a raw, unstructured text format that requires custom parsing to extract coverage metrics. Currently, gocovtui cannot ingest or analyze Go coverage data, limiting its ability to be a functional coverage visualization tool. By building a robust parser, we enable the core data ingestion layer needed for all downstream features—visualization, filtering, and reporting.

## What Changes

- Add a new `coverage` package with a parser that reads Go's coverprofile format and extracts structured coverage statistics
- Parse coverprofile files to extract file-level and function-level coverage data
- Calculate aggregate coverage metrics (total statements covered, total statements, coverage percentage)
- Provide a structured API for accessing parsed coverage data
- Enable the TUI to consume coverage data for visualization and analysis

## Capabilities

### New Capabilities
- `coverprofile-parser`: Parse Go test coverprofile files and extract coverage statistics (files, functions, statements, coverage percentages)
- `coverage-stats`: Aggregate and calculate coverage metrics from parsed data (total statements, covered statements, coverage percentage, per-file stats)

### Modified Capabilities

## Impact

- **New code**: `pkg/coverage/parser.go`, `pkg/coverage/stats.go`, data models for coverage results
- **Tests**: Unit tests for parser and stats calculation
- **API changes**: New public types and functions in coverage package
- **Existing code**: No breaking changes to current app structure; coverage package is additive
- **Dependencies**: No new external dependencies required (uses Go stdlib)
