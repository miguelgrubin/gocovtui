package tui

// rowKind distinguishes folder summary rows from file rows.
type rowKind int

const (
	kindFolder rowKind = iota
	kindFile
)

// rowData holds the display data for a single table row.
type rowData struct {
	kind     rowKind
	name     string
	total    int
	covered  int
	coverPct float64
}
