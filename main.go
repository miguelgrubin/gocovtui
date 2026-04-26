package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/miguelgrubin/gocovtui/pkg"
)

func main() {
	coverprofilePath := ""
	if len(os.Args) > 1 {
		coverprofilePath = os.Args[1]
	}

	app := pkg.NewApp(coverprofilePath)
	model := app.NewTUIModel()

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
