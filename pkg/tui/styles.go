package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Synthwave color palette
const (
	colorBackground = lipgloss.Color("#1a1a2e")
	colorPink       = lipgloss.Color("#ff2d78")
	colorCyan       = lipgloss.Color("#00f5d4")
	colorPurple     = lipgloss.Color("#7b2fff")
	colorYellow     = lipgloss.Color("#ffe600")
	colorGray       = lipgloss.Color("#a0a0c0")
	colorDimGray    = lipgloss.Color("#4a4a6a")
)

var (
	// titleStyle is used for the application title.
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPink).
			Background(colorBackground).
			Padding(0, 1)

	// headerStyle wraps the summary bar at the top.
	headerStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Background(colorBackground).
			Bold(true).
			Padding(0, 1)

	// headerLabelStyle styles key labels inside the header.
	headerLabelStyle = lipgloss.NewStyle().
				Foreground(colorGray).
				Background(colorBackground)

	// headerValueStyle styles numeric values inside the header.
	headerValueStyle = lipgloss.NewStyle().
				Foreground(colorYellow).
				Background(colorBackground).
				Bold(true)

	// coverageHighStyle colors coverage values ≥ 80%.
	coverageHighStyle = lipgloss.NewStyle().Foreground(colorCyan).Bold(true).Padding(0, 1)

	// coverageMidStyle colors coverage values 50–79%.
	coverageMidStyle = lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Padding(0, 1)

	// coverageLowStyle colors coverage values < 50%.
	coverageLowStyle = lipgloss.NewStyle().Foreground(colorPink).Bold(true).Padding(0, 1)

	// Table styles
	tableHeaderStyle = lipgloss.NewStyle().
				Foreground(colorPurple).
				Bold(true).
				Align(lipgloss.Center).
				Padding(0, 1)

	tableFolderStyle = lipgloss.NewStyle().
				Foreground(colorPink).
				Bold(true).
				Padding(0, 1)

	tableFileStyle = lipgloss.NewStyle().
			Foreground(colorGray).
			Padding(0, 1)

	tableSelectedFolderStyle = lipgloss.NewStyle().
					Foreground(colorBackground).
					Background(colorPink).
					Bold(true).
					Padding(0, 1)

	tableSelectedFileStyle = lipgloss.NewStyle().
				Foreground(colorBackground).
				Background(colorPurple).
				Bold(true).
				Padding(0, 1)

	// Find mode styles
	findInputStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Background(colorBackground).
			Padding(0, 1)

	findLabelStyle = lipgloss.NewStyle().
			Foreground(colorPurple).
			Bold(true).
			Background(colorBackground).
			Padding(0, 1)

	findNoResultsStyle = lipgloss.NewStyle().
				Foreground(colorPink).
				Background(colorBackground).
				Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorPink).
			Bold(true).
			Background(colorBackground).
			Padding(0, 1)
)

// coverageStyle returns the appropriate lipgloss style based on the coverage percentage.
func coverageStyle(pct float64) lipgloss.Style {
	switch {
	case pct >= 80:
		return coverageHighStyle
	case pct >= 50:
		return coverageMidStyle
	default:
		return coverageLowStyle
	}
}
