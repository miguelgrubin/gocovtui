## Context

`pkg/coverage/stats.go` tracks per-file coverage via `FileStats` and a flat `map[string]*FileStats`. The TUI list in `pkg/tui/model.go` iterates `FilesSortedByCoverage` and produces one `fileItem` per file. There is no concept of folder/package grouping in either the data layer or the presentation layer.

Go coverage profile filenames follow the module path convention, e.g. `github.com/miguelgrubin/gocovtui/pkg/coverage/parser.go`. The "folder" is everything up to the last `/`.

## Goals / Non-Goals

**Goals:**
- Add `FolderStats` type to `pkg/coverage/stats.go` with `Dir`, `TotalStatements`, `CoveredStatements`, `CoveragePercent`, and `FileCount` fields
- Add `FolderStats() []*FolderStats` and `FoldersSortedByCoverage(ascending bool) []*FolderStats` methods on `Stats`
- Add `folderItem` list item type in `pkg/tui/item.go` with a distinct visual style (bold, folder-color, no statement count detail)
- Update `NewModel` in `pkg/tui/model.go` to build a grouped list: for each folder (sorted by coverage ascending), emit one `folderItem` then its `fileItem`s sorted by coverage ascending
- Add a `folderRowStyle` in `pkg/tui/styles.go`

**Non-Goals:**
- Collapsible/expandable folder groups
- Keyboard navigation that skips folder rows (folder items are navigable like any list item; selecting one does nothing special)
- Recursive sub-folder grouping (one level: immediate parent directory only)
- Filtering by folder

## Decisions

### Folder key = `path.Dir(filename)`
**Rationale**: The Go standard `path` package (not `filepath`) handles the forward-slash module paths in coverage filenames correctly on all platforms. Alternative of using `strings.LastIndex` is equivalent but less idiomatic.

### Folder aggregation lives in `pkg/coverage/stats.go`
**Rationale**: Keeps presentation-agnostic data logic in the data layer. The TUI should not need to group files itself — it just iterates the pre-grouped data.

### Grouped list order: folders by coverage ascending, files within folder by coverage ascending
**Rationale**: Worst folders/files bubble to the top so the engineer immediately sees where coverage is lowest.

### `folderItem` implements `list.Item` but uses `fileDelegate.Render` via a type switch
**Rationale**: Single delegate with a type switch keeps rendering logic co-located. Alternative of separate delegate per type would require a custom list adapter.

## Risks / Trade-offs

- [Module path as folder key] → Files at the module root (no subdirectory) will share folder key equal to the module path itself. This is acceptable — they render as one folder group. Mitigation: none needed.
- [list.Model does not support non-selectable rows] → Folder rows are selectable (cursor can land on them) but pressing enter/space will have no effect since no action is wired. This is a minor UX imperfection. Mitigation: document in help text or style folder rows clearly as headers.
