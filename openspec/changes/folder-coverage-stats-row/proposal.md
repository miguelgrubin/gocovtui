## Why

The TUI currently shows a flat list of individual file coverage stats. When a project has many files spread across many packages/folders, navigating the list is tedious and gives no quick sense of which packages are most problematic. Adding folder-level aggregate rows makes it easy to spot low-coverage packages at a glance.

## What Changes

- Compute per-folder aggregate coverage stats by grouping `FileStats` by their directory path
- Introduce a `FolderStats` type in `pkg/coverage/stats.go` that aggregates statements and coverage across all files under a folder
- Add a `FoldersSortedByCoverage` method on `Stats` mirroring the existing files API
- Display folder rows in the TUI list as visual section separators or summary rows, styled differently from file rows (bold, indented differently, distinct Synthwave color)
- Folder rows are non-interactive (no selection action), purely informational

## Capabilities

### New Capabilities

- `folder-coverage-stats`: Aggregation of coverage statistics by folder/package, with a new `FolderStats` type and query methods on `Stats`

### Modified Capabilities

- `tui-coverage-viewer`: The file list now interleaves folder summary rows above each group of files that belong to that folder, with a distinct visual style

## Impact

- `pkg/coverage/stats.go`: New `FolderStats` type, new `FolderStats()` and `FoldersSortedByCoverage()` methods on `Stats`
- `pkg/tui/item.go`: New `folderItem` list item type rendered with a distinct style (bold folder name + aggregate coverage)
- `pkg/tui/model.go`: List population logic updated to interleave folder rows before each file group
- `pkg/tui/styles.go`: New style for folder rows
