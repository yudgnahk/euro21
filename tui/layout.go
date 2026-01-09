package tui

import (
	"context"
)

// TableSection represents a section with a header and table
type TableSection struct {
	Header string
	Table  *Table
}

// MultiTable renders multiple tables with section headers
type MultiTable struct {
	sections []TableSection
	theme    *Theme
	spacing  int // Lines between sections
}

// NewMultiTable creates a new multi-table widget
func NewMultiTable() *MultiTable {
	return &MultiTable{
		sections: []TableSection{},
		theme:    DefaultTheme,
		spacing:  1,
	}
}

// AddSection adds a table section with a header
func (mt *MultiTable) AddSection(header string, table *Table) *MultiTable {
	mt.sections = append(mt.sections, TableSection{
		Header: header,
		Table:  table,
	})
	return mt
}

// SetTheme sets the theme for all sections
func (mt *MultiTable) SetTheme(theme *Theme) *MultiTable {
	mt.theme = theme
	for i := range mt.sections {
		mt.sections[i].Table.SetTheme(theme)
	}
	return mt
}

// SetSpacing sets the number of blank lines between sections
func (mt *MultiTable) SetSpacing(lines int) *MultiTable {
	mt.spacing = lines
	return mt
}

// NormalizeWidths calculates the maximum width needed for each column
// across all tables and applies it as minimum width to ensure uniform table sizes
func (mt *MultiTable) NormalizeWidths() *MultiTable {
	if len(mt.sections) == 0 {
		return mt
	}

	// Find the maximum number of columns
	maxCols := 0
	for _, section := range mt.sections {
		if len(section.Table.headers) > maxCols {
			maxCols = len(section.Table.headers)
		}
	}

	// Calculate maximum width for each column across all tables
	maxWidths := make([]int, maxCols)
	for _, section := range mt.sections {
		colWidths := section.Table.calculateColumnWidths()
		for i, width := range colWidths {
			if i < maxCols && width > maxWidths[i] {
				maxWidths[i] = width
			}
		}
	}

	// Apply the maximum widths to all tables
	for i := range mt.sections {
		mt.sections[i].Table.SetMinColWidths(maxWidths...)
	}

	return mt
}

// Render implements the Widget interface
func (mt *MultiTable) Render(ctx context.Context, buf *Buffer) error {
	for i, section := range mt.sections {
		// Check context
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Render section header
		if section.Header != "" {
			buf.WriteString(section.Header)
			buf.NewLine()
		}

		// Render table
		if err := section.Table.Render(ctx, buf); err != nil {
			return err
		}

		// Add spacing between sections
		if i < len(mt.sections)-1 {
			for j := 0; j < mt.spacing; j++ {
				buf.NewLine()
			}
		}
	}

	return nil
}

// MinSize implements the Widget interface
func (mt *MultiTable) MinSize() Size {
	var totalWidth, totalHeight int

	for i, section := range mt.sections {
		size := section.Table.MinSize()

		// Track maximum width
		if size.Width > totalWidth {
			totalWidth = size.Width
		}

		// Add height
		if section.Header != "" {
			totalHeight++ // Header line
		}
		totalHeight += size.Height

		// Add spacing
		if i < len(mt.sections)-1 {
			totalHeight += mt.spacing
		}
	}

	return Size{Width: totalWidth, Height: totalHeight}
}
