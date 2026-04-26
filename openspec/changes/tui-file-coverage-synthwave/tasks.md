## 1. Dependencies

- [x] 1.1 Add `github.com/charmbracelet/bubbletea` to `go.mod` via `go get`
- [x] 1.2 Add `github.com/charmbracelet/bubbles` to `go.mod` via `go get`
- [x] 1.3 Run `go mod tidy` to clean up `go.sum`

## 2. Synthwave Theme & Styles

- [x] 2.1 Create `pkg/tui/styles.go` with Synthwave color constants (dark navy, neon pink, neon cyan, purple, bright yellow)
- [x] 2.2 Define lipgloss styles for: header, list item (normal/selected), border, title, coverage values (high/mid/low thresholds)

## 3. TUI Model

- [x] 3.1 Create `pkg/tui/model.go` with `Model` struct holding `list.Model`, coverage summary, and terminal dimensions
- [x] 3.2 Implement `Init() tea.Cmd` returning `nil`
- [x] 3.3 Implement `Update(msg tea.Msg) (tea.Model, tea.Cmd)` handling `tea.KeyMsg` (q/ctrl+c to quit), `tea.WindowSizeMsg` (resize), and delegating to the list
- [x] 3.4 Implement `View() string` composing header + list rendered with Synthwave styles

## 4. List Item

- [x] 4.1 Create `pkg/tui/item.go` defining a `fileItem` type implementing `list.Item` and `list.DefaultDelegate` (or custom delegate)
- [x] 4.2 Implement item rendering: filename (left-aligned), statement count, and coverage % colored by threshold

## 5. Summary Header

- [x] 5.1 Implement `renderHeader(summary coverage.SummaryStats) string` in `pkg/tui/model.go` using Synthwave styles
- [x] 5.2 Header SHALL show: total coverage %, covered/total statements, file count

## 6. App Wiring

- [x] 6.1 Refactor `pkg/app.go`: remove static `lipgloss.Println` calls, expose a `NewTUIModel(stats *coverage.Stats) tui.Model` constructor
- [x] 6.2 Update `main.go` to accept a coverage file path argument (os.Args), call `pkg.NewApp(path)`, build TUI model, and run `tea.NewProgram(model)`
- [x] 6.3 Remove the `RenderTable` demo function from `pkg/app.go`

## 7. Verification

- [x] 7.1 Run `go build ./...` and confirm no compilation errors
- [x] 7.2 Run `go test ./...` and confirm existing tests pass
- [ ] 7.3 Manually test with `go run . coverage.out` and verify interactive TUI renders with Synthwave theme, list is scrollable, and `q` quits
