## 1. Create file picker package

- [x] 1.1 Create `pkg/filepicker/` directory and `filepicker.go` file
- [x] 1.2 Define `Model` struct wrapping `bubbles/filepicker.Model` with `AllowedTypes: []string{".out"}`
- [x] 1.3 Implement `Init()`, `Update()`, `View()` to satisfy `tea.Model`
- [x] 1.4 Define `SelectedMsg` (carries chosen file path) and `CancelMsg` types
- [x] 1.5 Emit `SelectedMsg` on file selection and `CancelMsg` on quit/cancel

## 2. Handle empty file list

- [x] 2.1 On `Init`, glob CWD for `*.out` files
- [x] 2.2 If no files found, return an error from a helper so `main.go` can print to stderr and exit non-zero

## 3. Wire into main.go

- [x] 3.1 When `len(os.Args) == 1`, run the file picker as a separate `tea.NewProgram` (no alt screen)
- [x] 3.2 Collect the result: if `SelectedMsg`, use path; if `CancelMsg` or error, exit cleanly
- [x] 3.3 Pass chosen path into existing `pkg.NewApp(coverprofilePath)` flow

## 4. Verify behaviour

- [x] 4.1 Manual smoke test: run without args, pick a file, confirm main TUI opens
- [x] 4.2 Manual smoke test: run without args, press `q`, confirm clean exit
- [x] 4.3 Manual smoke test: run with arg, confirm picker is skipped
- [x] 4.4 Run `go build ./...` to confirm no compile errors
