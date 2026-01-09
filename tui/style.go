package tui

import (
	"fmt"

	"github.com/yudgnahk/euro21/constants"
)

// Style defines visual appearance
type Style struct {
	Foreground Color
	Background Color
	Bold       bool
	Italic     bool
	Underline  bool
}

// Color represents a terminal color
type Color interface {
	// ANSI returns the ANSI escape code for this color
	ANSI() string
	// String returns a string representation
	String() string
}

// BasicColor represents a basic ANSI color
type BasicColor struct {
	code int
}

// ANSI returns the ANSI escape code
func (b BasicColor) ANSI() string {
	return fmt.Sprintf("\033[%dm", b.code)
}

// String returns a string representation
func (b BasicColor) String() string {
	return fmt.Sprintf("Color(%d)", b.code)
}

// NewColor creates a Color from constants.Color
func NewColor(c constants.Color) Color {
	// Parse ANSI code from constants.Color (format: "\033[XXm")
	switch c {
	case constants.ColorWhite:
		return ColorWhite
	case constants.ColorGreen:
		return ColorGreen
	case constants.ColorYellow:
		return ColorYellow
	case constants.ColorRed:
		return ColorRed
	case constants.ColorBlue:
		return ColorBlue
	case constants.ColorCyan:
		return ColorCyan
	case constants.ColorPurple:
		return ColorPurple
	default:
		return ColorWhite
	}
}

// Predefined colors
var (
	ColorWhite  = BasicColor{code: 37}
	ColorGreen  = BasicColor{code: 32}
	ColorYellow = BasicColor{code: 33}
	ColorRed    = BasicColor{code: 31}
	ColorBlue   = BasicColor{code: 34}
	ColorCyan   = BasicColor{code: 36}
	ColorPurple = BasicColor{code: 35}
	ColorReset  = BasicColor{code: 0}
)

// BorderStyle defines the characters used for borders
type BorderStyle struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
	MiddleLeft  string
	MiddleRight string
	Cross       string
	TopJoin     string
	BottomJoin  string
}

// Theme is a collection of styles
type Theme struct {
	Name        string
	TableHeader Style
	TableRow    Style
	TableRowAlt Style // Alternating row
	Qualified   Style // Green for qualified teams
	Conditional Style // Yellow for conditional
	Eliminated  Style // White/default
	Border      BorderStyle
}

// DefaultTheme is the default theme with Unicode box-drawing characters
var DefaultTheme = &Theme{
	Name:        "Default",
	TableHeader: Style{Foreground: ColorWhite},
	TableRow:    Style{Foreground: ColorWhite},
	TableRowAlt: Style{Foreground: ColorWhite},
	Qualified:   Style{Foreground: ColorGreen},
	Conditional: Style{Foreground: ColorYellow},
	Eliminated:  Style{Foreground: ColorWhite},
	Border: BorderStyle{
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "└",
		BottomRight: "┘",
		Horizontal:  "─",
		Vertical:    "│",
		MiddleLeft:  "├",
		MiddleRight: "┤",
		Cross:       "┼",
		TopJoin:     "┬",
		BottomJoin:  "┴",
	},
}

// MinimalTheme is an ASCII-only theme for terminals without Unicode support
var MinimalTheme = &Theme{
	Name:        "Minimal",
	TableHeader: Style{Foreground: ColorWhite},
	TableRow:    Style{Foreground: ColorWhite},
	TableRowAlt: Style{Foreground: ColorWhite},
	Qualified:   Style{Foreground: ColorGreen},
	Conditional: Style{Foreground: ColorYellow},
	Eliminated:  Style{Foreground: ColorWhite},
	Border: BorderStyle{
		TopLeft:     "+",
		TopRight:    "+",
		BottomLeft:  "+",
		BottomRight: "+",
		Horizontal:  "-",
		Vertical:    "|",
		MiddleLeft:  "+",
		MiddleRight: "+",
		Cross:       "+",
		TopJoin:     "+",
		BottomJoin:  "+",
	},
}

// ModernTheme is a modern theme with enhanced styling
var ModernTheme = DefaultTheme
