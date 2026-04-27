package tui

import (
	"testing"
)

func TestCoverageStyle_High(t *testing.T) {
	tests := []struct {
		name string
		pct  float64
	}{
		{"exactly 80", 80.0},
		{"above 80", 95.5},
		{"100 percent", 100.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := coverageStyle(tt.pct)
			if s.GetForeground() != coverageHighStyle.GetForeground() {
				t.Errorf("coverageStyle(%.1f) should return high style", tt.pct)
			}
		})
	}
}

func TestCoverageStyle_Mid(t *testing.T) {
	tests := []struct {
		name string
		pct  float64
	}{
		{"exactly 50", 50.0},
		{"mid range", 65.0},
		{"just below 80", 79.9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := coverageStyle(tt.pct)
			if s.GetForeground() != coverageMidStyle.GetForeground() {
				t.Errorf("coverageStyle(%.1f) should return mid style", tt.pct)
			}
		})
	}
}

func TestCoverageStyle_Low(t *testing.T) {
	tests := []struct {
		name string
		pct  float64
	}{
		{"zero", 0.0},
		{"below 50", 30.0},
		{"just below 50", 49.9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := coverageStyle(tt.pct)
			if s.GetForeground() != coverageLowStyle.GetForeground() {
				t.Errorf("coverageStyle(%.1f) should return low style", tt.pct)
			}
		})
	}
}
