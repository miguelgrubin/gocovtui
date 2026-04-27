# gocovtui

A terminal UI for visualizing Go coverage profiles, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss).

## Screenshot

<!-- screenshot -->

## Features

- Parses standard Go coverage profiles (`coverage.out`)
- Groups files by package/folder with summary rows
- Sorts folders and files by coverage percentage (ascending)
- Color-coded coverage: high (≥80%), mid (50–79%), low (<50%)
- Keyboard navigation with Vim-style keybindings
- Search/filter mode for quick file lookup
- Open files in your editor directly from the TUI
- Interactive file picker when no argument is provided

## Keyboard Shortcuts

### Navigation

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Home` / `g` | Go to top |
| `End` / `G` | Go to bottom |
| `f` | Enter find mode |
| `e` | Open file in editor |
| `q` / `Ctrl+C` | Quit |

### Find Mode

| Key | Action |
|-----|--------|
| Type | Filter by name |
| `↑` / `Shift+Tab` | Previous match |
| `↓` / `Tab` | Next match |
| `Backspace` | Delete character |
| `Enter` | Confirm selection |
| `Esc` / `Ctrl+C` | Cancel |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `EDITOR` | Editor used when pressing `e` to open a file |

## Installation

```sh
go install github.com/miguelgrubin/gocovtui@latest
```

Or build from source:

```sh
git clone https://github.com/miguelgrubin/gocovtui.git
cd gocovtui
make build
```

## Usage

Generate a coverage profile and pass it to `gocovtui`:

```sh
go test ./... -coverprofile=coverage.out
gocovtui coverage.out
```

Run without arguments to open an interactive file picker:

```sh
gocovtui
```
