## Why

The current app renders coverage data as a static lipgloss table with no interactivity or visual polish. A proper interactive TUI using Bubble Tea with a Synthwave color theme will make exploring file-level coverage stats more engaging and navigable.

## What Changes

- Replace static `lipgloss.Println` rendering in `pkg/app.go` with an interactive Bubble Tea model
- Add a Synthwave-themed color palette (deep purples, hot pinks, cyan neons, dark backgrounds)
- Implement a scrollable, keyboard-navigable file coverage list with per-file stats
- Display a summary header (total coverage %, covered/total statements, file count)
- Add `github.com/charmbracelet/bubbletea` and `github.com/charmbracelet/bubbles` as dependencies
- Wire `main.go` to run the Bubble Tea program instead of the current static render

## Capabilities

### New Capabilities

- `tui-coverage-viewer`: Interactive TUI screen that displays file coverage statistics in a scrollable list with Synthwave theme, keyboard navigation, and a summary header

### Modified Capabilities

- (none)

## Impact

- `pkg/app.go`: Refactored to expose a Bubble Tea `Model` instead of direct rendering
- `main.go`: Updated to run `tea.NewProgram(model)`
- `go.mod` / `go.sum`: New dependencies `charmbracelet/bubbletea` and `charmbracelet/bubbles`
- `pkg/tui/` (new): TUI model, styles, and view logic
