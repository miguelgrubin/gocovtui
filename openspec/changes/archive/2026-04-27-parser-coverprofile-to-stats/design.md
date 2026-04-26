## Context

Go's coverprofile format is standardized text output from `go test -coverprofile=coverage.out`. Each line contains:
```
<filename>:<startLine>.<startCol>,<endLine>.<endCol> <numStmt> <count>
```

The gocovtui project currently has a proof-of-concept TUI structure but no data ingestion capability. Building the parser creates the foundation for visualization features. The project uses Go stdlib only, has a structured package layout (`pkg/`), and includes a Makefile with test/coverage automation.

## Goals / Non-Goals

**Goals:**
- Parse Go coverprofile files (format defined by `go test -cover` tooling)
- Extract file-level coverage (which files are tested, coverage %)
- Extract function-level coverage (which functions are tested, coverage %)
- Calculate aggregate statistics (total coverage %, total statements, covered statements)
- Provide clean Go API for consuming parsed data
- Support multiple coverprofile files (merging coverage data)
- Be fast and memory-efficient for large test suites

**Non-Goals:**
- Visualization of coverage data (that's the TUI layer, separate concern)
- Custom/proprietary coverage formats (only Go's standard format)
- Integration with CI/CD systems
- HTML/JSON export formats (can be added later)

## Decisions

### 1. Package Structure: `pkg/coverage` for all coverage-related code
**Rationale**: Keeps coverage logic separate from TUI concerns. Aligns with existing `pkg/app.go` structure. Easy to test independently.
**Alternatives**: Put in `internal/coverage` (more restrictive) - rejected because this may be exported for library use in future.

### 2. Separate `parser.go` and `stats.go` modules
**Rationale**: Parser handles file I/O and line parsing (unit testable). Stats handles aggregation logic (testable independently). Clear separation of concerns.
**Alternatives**: Single monolithic file - rejected for maintainability and testability.

### 3. Data Models: Use simple Go structs (no external dependencies)
**Rationale**: Coverprofile is simple format; stdlib `bufio.Scanner` is sufficient. Keeps binary size small. Project has no external deps currently.
**Alternatives**: Use external parsing library - rejected to match project philosophy of minimal deps.

### 4. No Caching Layer Initially
**Rationale**: Start simple. Profile later if needed. Re-parsing is O(n) with small constant; not a bottleneck.
**Alternatives**: Add memoization - premature optimization; defer to v2.

### 5. Handle Missing Function Names Gracefully
**Rationale**: Go coverprofile doesn't always include function names directly. Parse what we can; allow partial coverage data.
**Alternatives**: Fail on incomplete data - rejected; better to work with imperfect data.

## Risks / Trade-offs

**[Risk] Large coverprofile files**: Parsing very large files (millions of statements) could use significant memory.
→ Mitigation: Use streaming parser with `bufio.Scanner`, not loading entire file into memory. Profile later if needed.

**[Risk] Multiple file coverage merging**: Merging coverage from multiple test runs could have subtle bugs (e.g., overlapping line ranges).
→ Mitigation: Start with single-file parsing. Add merging logic only when needed; include comprehensive unit tests.

**[Risk] Coverprofile format changes**: Future Go versions might change format slightly.
→ Mitigation: Add version detection; document format expectations. Monitor Go releases.

**[Trade-off] Accuracy vs Speed**: Fully resolving function names requires AST parsing (slow). We'll extract what coverprofile gives us.
→ Acceptance: This is acceptable for MVP. Can improve later with AST analysis if needed.

## Migration Plan

1. Add `pkg/coverage/` directory with parser and stats modules
2. Add unit tests in `pkg/coverage/*_test.go`
3. Integration test: run `make coverage`, parse output, verify statistics
4. Update `pkg/app.go` to optionally load coverage data (no breaking changes)
5. Merge to main; document in README

Rollback: Simple - remove `pkg/coverage/` directory.

## Open Questions

1. Should we support merging multiple coverprofile files (e.g., from parallel test runs)?
   - **Decision pending**: Start with single file, add merging if TUI requires it.

2. How precise do we need function-level coverage?
   - **Decision pending**: Coverprofile doesn't give us full function names reliably. Acceptable to work with line ranges for MVP.
