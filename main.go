package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/miguelgrubin/gocovtui/pkg"
	"github.com/miguelgrubin/gocovtui/pkg/filepicker"
)

func main() {
	coverprofilePath := ""

	if len(os.Args) > 1 {
		coverprofilePath = os.Args[1]
	} else {
		// No argument provided — show interactive file picker for *.out files.
		pickerModel, err := filepicker.NewModel()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		result, err := tea.NewProgram(pickerModel).Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error running file picker: %v\n", err)
			os.Exit(1)
		}

		m, ok := result.(filepicker.Model)
		if !ok {
			fmt.Fprintf(os.Stderr, "error: unexpected model type from file picker\n")
			os.Exit(1)
		}

		path, selected := m.Result()
		if !selected {
			// User cancelled — exit cleanly.
			os.Exit(0)
		}
		coverprofilePath = path
	}

	app := pkg.NewApp(coverprofilePath)
	model := app.NewTUIModel()

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
