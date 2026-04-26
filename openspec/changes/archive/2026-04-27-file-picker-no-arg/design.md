## Context

`gocovtui` is a TUI Go coverage viewer built with `charmbracelet/bubbletea`. Currently `main.go` reads `os.Args[1]` as the coverage file path. When the program is run without arguments, the path is empty and the UI shows nothing useful. The `charmbracelet/bubbles` library is already a dependency and includes a `filepicker` component.

## Goals / Non-Goals

**Goals:**
- Show an interactive file picker when no CLI argument is provided.
- Filter the picker to display only `*.out` files in the current working directory.
- Hand off the selected path to the existing app flow seamlessly.
- Exit cleanly when the user cancels the picker.

**Non-Goals:**
- Recursive directory traversal or multi-directory navigation.
- Remembering the last-used file across sessions.
- Changing the existing behaviour when an argument IS supplied.

## Decisions

### 1. Use `bubbles/filepicker` component

`charmbracelet/bubbles` already ships a `filepicker.Model`. We use it with `AllowedTypes: []string{".out"}` to restrict visible entries. Alternative: build a custom list — rejected because the built-in component already handles keyboard navigation, filtering, and styling.

### 2. New `pkg/filepicker` package

Encapsulate the picker in `pkg/filepicker/filepicker.go` as a standalone `tea.Model`. It emits a `SelectedMsg` (carrying the chosen path) or a `CancelMsg` when the user quits. `main.go` runs a short bubbletea program just for the picker, then launches the main app with the returned path.

**Alternative considered**: embed the picker inside the existing TUI model as a startup phase. Rejected — it couples unrelated concerns and complicates the existing model.

### 3. Two-program approach in `main.go`

Run a small picker program first (`tea.NewProgram(filepicker.NewModel())`), collect the result, then run the main TUI program. This keeps each program single-purpose and easy to test.

## Risks / Trade-offs

- **No `.out` files in CWD** → The picker displays an empty list. Risk: user confusion. Mitigation: show a status message "No *.out files found" and exit with a clear error.
- **bubbles filepicker API changes** → Pinned to `v1.0.0` already in `go.mod`. Mitigation: no action needed.
- **Two sequential `tea.NewProgram` calls** → Works fine in practice with `tea.WithAltScreen()` on the second call; the terminal is restored between runs.
