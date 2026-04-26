## Context

The project already parses Go coverage profiles (`pkg/coverage/`) and has `charm.land/lipgloss/v2` for styling. Currently `pkg/app.go` renders a static table to stdout using `lipgloss.Println`. The goal is to replace this with a proper interactive TUI driven by Bubble Tea.

The app receives a coverage profile path as input, computes per-file stats, and needs to present them in a navigable list with a Synthwave visual identity.

## Goals / Non-Goals

**Goals:**
- Add Bubble Tea (`charmbracelet/bubbletea`) as the program runtime
- Add Bubbles (`charmbracelet/bubbles`) for the list component with scroll support
- Define a Synthwave color palette using lipgloss: dark background (`#1a1a2e`), neon pink (`#ff2d78`), neon cyan (`#00f5d4`), purple (`#7b2fff`), bright yellow (`#ffe600`)
- Implement a `pkg/tui` package with a `Model` that satisfies `tea.Model` (`Init`, `Update`, `View`)
- Show a header with aggregate stats (total coverage %, covered/total stmts, file count)
- Show a scrollable list of files with filename, statement count, and coverage % colored by threshold (green-ish if ≥80%, yellow if 50–79%, red/pink if <50%)
- Wire `main.go` to run `tea.NewProgram`
- Support `q` / `ctrl+c` to quit

**Non-Goals:**
- Drill-down into individual file line coverage
- Mouse support
- Config file or persistent settings
- Windows-specific terminal handling

## Decisions

### Use `charmbracelet/bubbletea` as the TUI runtime
**Rationale**: The project is already in the Charm ecosystem (lipgloss). Bubble Tea provides the standard Elm-architecture event loop for interactive terminal UIs. Alternative (tcell/tview) would diverge from existing dependency set.

### Use `bubbles/list` for the file list
**Rationale**: Provides keyboard navigation, scrolling, and item rendering hooks out of the box. Writing a custom scroller would duplicate effort. The list's `Styles` can be overridden with Synthwave colors.

### New `pkg/tui` package
**Rationale**: Separates TUI concerns from `pkg/app.go` which currently mixes data loading with rendering. `app.go` will remain as the entry point that wires coverage data → TUI model, but rendering logic moves to `pkg/tui`.

### Synthwave palette as constants in `pkg/tui/styles.go`
**Rationale**: Centralizing colors makes it easy to adjust the theme without hunting through view code.

### Coverage threshold coloring
- ≥ 80%: neon cyan (`#00f5d4`)
- 50–79%: bright yellow (`#ffe600`)
- < 50%: neon pink (`#ff2d78`)

## Risks / Trade-offs

- [lipgloss v2 API drift] → `charm.land/lipgloss/v2` may differ from `github.com/charmbracelet/lipgloss` used by bubbles. Mitigation: pin compatible versions; use only stable APIs (NewStyle, Foreground, Bold, etc.).
- [Terminal width] → List item widths need to adapt to terminal size. Mitigation: handle `tea.WindowSizeMsg` in `Update` and set list width/height accordingly.
- [bubbles/list default styles] → Default styles may clash with Synthwave theme. Mitigation: fully override `list.DefaultStyles()` in styles setup.
