# gocovtui

A terminal UI for visualizing Go coverage profiles, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss).

## Screenshot

<!-- screenshot -->

## Features

- Parses standard Go coverage profiles (`coverage.out`)
- Groups files by package/folder with summary rows
- Sorts folders and files by coverage percentage (ascending)
- Color-coded coverage: high (≥80%), mid (50–79%), low (<50%)
- Keyboard navigation: `↑`/`↓`, `g`/`G` (top/bottom), `q` to quit
- Interactive file picker when no argument is provided

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
