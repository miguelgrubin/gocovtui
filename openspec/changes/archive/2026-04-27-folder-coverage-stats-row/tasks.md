## 1. Coverage Data Layer

- [x] 1.1 Add `FolderStats` struct to `pkg/coverage/stats.go` with fields: `Dir string`, `TotalStatements int`, `CoveredStatements int`, `CoveragePercent float64`, `FileCount int`
- [x] 1.2 Add `FolderStats() []*FolderStats` method on `Stats` that groups files by `path.Dir(filename)` and aggregates statement counts
- [x] 1.3 Add `FoldersSortedByCoverage(ascending bool) []*FolderStats` method on `Stats`
- [x] 1.4 Write unit tests for `FolderStats()` and `FoldersSortedByCoverage()` in `pkg/coverage/stats_test.go`

## 2. TUI Styles

- [x] 2.1 Add `folderRowStyle` to `pkg/tui/styles.go` (bold, Synthwave purple/pink accent, distinct from file item styles)
- [x] 2.2 Add `folderSelectedStyle` for when the cursor is on a folder row

## 3. TUI List Item

- [x] 3.1 Add `folderItem` struct to `pkg/tui/item.go` implementing `list.Item` with `dir string`, `coverPct float64`, `fileCount int`, `totalStmts int`, `coveredStmts int`
- [x] 3.2 Implement `FilterValue()`, `Title()`, `Description()` on `folderItem`
- [x] 3.3 Update `fileDelegate.Render` to type-switch on `folderItem` and render it with `folderRowStyle` / `folderSelectedStyle` showing folder path and aggregate coverage %

## 4. TUI Model

- [x] 4.1 Update `NewModel` in `pkg/tui/model.go` to build a grouped list: iterate `FoldersSortedByCoverage(true)`, and for each folder emit one `folderItem` followed by its `fileItem`s sorted by coverage ascending

## 5. Verification

- [x] 5.1 Run `go build ./...` and confirm no compilation errors
- [x] 5.2 Run `go test ./...` and confirm all tests pass including new `FolderStats` tests
- [ ] 5.3 Manually run `go run . coverage.out` and verify folder rows appear above each file group with distinct styling
