// Package tui provides terminal user interface widgets for displaying
// structured data like tables, trees, and brackets.
//
// # Basic Usage
//
// Create a table:
//
//	table := tui.NewTable().
//	    SetHeaders("Name", "Score", "Status").
//	    AddRow("Italy", "9", "Qualified").
//	    AddRow("Spain", "8", "Qualified")
//
// Render to stdout:
//
//	renderer := tui.NewRenderer(os.Stdout)
//	renderer.Render(context.Background(), table)
//
// # Theming
//
//	renderer.SetTheme(tui.ModernTheme)
//
// # Knockout Brackets
//
//	bracket := tui.NewBracketTree()
//	bracket.BuildFromMatches(matches)
//	renderer.Render(context.Background(), bracket)
package tui

import (
	"context"
)

// Widget is the base interface for all TUI components
type Widget interface {
	// Render writes the widget to the buffer
	Render(ctx context.Context, buf *Buffer) error

	// MinSize returns the minimum size required
	MinSize() Size
}

// Size represents widget dimensions
type Size struct {
	Width  int
	Height int
}

// Position represents widget coordinates
type Position struct {
	X int
	Y int
}

// Cell represents a single terminal cell
type Cell struct {
	Content string // Can be multi-rune (emoji, flags)
	Width   int    // Display width (1 for ASCII, 2 for emoji)
}

// BorderType indicates the type of border to render
type BorderType int

const (
	// BorderTop is the top border
	BorderTop BorderType = iota
	// BorderMiddle is a middle/separator border
	BorderMiddle
	// BorderBottom is the bottom border
	BorderBottom
)
