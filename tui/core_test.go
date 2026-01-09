package tui

import (
	"testing"
)

func TestDetect(t *testing.T) {
	term := Detect()
	if term == nil {
		t.Fatal("Detect() returned nil")
	}

	// Width should be positive (or default 80)
	if term.Width <= 0 {
		t.Errorf("Expected positive width, got %d", term.Width)
	}

	// Height should be positive (or default 24)
	if term.Height <= 0 {
		t.Errorf("Expected positive height, got %d", term.Height)
	}

	// ColorSupport should be a valid level
	if term.ColorSupport < ColorNone || term.ColorSupport > ColorTrueColor {
		t.Errorf("Invalid ColorSupport level: %d", term.ColorSupport)
	}
}

func TestBasicColorANSI(t *testing.T) {
	tests := []struct {
		name  string
		color BasicColor
		want  string
	}{
		{"white", ColorWhite, "\033[37m"},
		{"green", ColorGreen, "\033[32m"},
		{"yellow", ColorYellow, "\033[33m"},
		{"red", ColorRed, "\033[31m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.ANSI(); got != tt.want {
				t.Errorf("ANSI() = %v, want %v", got, tt.want)
			}
		})
	}
}
