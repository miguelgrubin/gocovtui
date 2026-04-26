## Why

When `gocovtui` is launched without a file argument, the coverage profile path is empty, leading to a broken or blank UI experience. An interactive file picker lets users select a `.out` coverage file directly from the terminal, making the tool usable without memorising file names.

## What Changes

- When no CLI argument is provided, display an interactive file picker (using `charmbracelet/bubbles`) instead of launching the main TUI immediately.
- The file picker lists only `*.out` files found in the current working directory.
- After the user selects a file, the main coverage TUI launches with the chosen path.
- If the user cancels or quits the picker without selecting a file, the program exits cleanly.

## Capabilities

### New Capabilities

- `file-picker`: Interactive TUI file picker that lists `*.out` files in the CWD, allowing the user to choose one before the main coverage viewer starts.

### Modified Capabilities

(none)

## Impact

- `main.go`: Add a conditional branch — if no argument is provided, run the file picker first; use its result as the coverage profile path.
- `pkg/filepicker/`: New package encapsulating the bubbles-based file picker model.
- `go.mod` / `go.sum`: No new dependencies; `charmbracelet/bubbles` is already present.
