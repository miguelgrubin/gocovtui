package pkg

import (
	"os"
	"strings"

	"github.com/miguelgrubin/gocovtui/pkg/coverage"
	"github.com/miguelgrubin/gocovtui/pkg/tui"
)

// App holds loaded coverage data.
type App struct {
	CoverageResult *coverage.CoverageResult
	CoverageStats  *coverage.Stats
}

// NewApp creates a new App, optionally loading coverage data from coverprofilePath.
func NewApp(coverprofilePath string) *App {
	app := &App{}

	if coverprofilePath != "" {
		result, err := coverage.ParseFile(coverprofilePath)
		if err == nil {
			app.CoverageResult = result
			app.CoverageStats = coverage.CalculateStats(result)
		}
	}

	return app
}

// NewTUIModel builds a tui.Model from the app's coverage stats.
func (a *App) NewTUIModel() tui.Model {
	return tui.NewModel(a.CoverageStats, readModulePath())
}

// readModulePath reads the module path from go.mod in the current directory.
func readModulePath() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return ""
}
