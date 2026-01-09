package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/yudgnahk/euro21/utils/stringutil"
)

// Alignment defines text alignment
type Alignment int

const (
	// AlignLeft aligns text to the left
	AlignLeft Alignment = iota
	// AlignCenter centers text
	AlignCenter
	// AlignRight aligns text to the right
	AlignRight
)

// Table renders tabular data
type Table struct {
	headers []string
	rows    [][]string
	colors  []Color // Per-row colors
	theme   *Theme

	// Options
	headerVisible bool
	borders       bool
	alternateRows bool
	minColWidths  []int
	maxColWidths  []int
	align         []Alignment
}

// NewTable creates a new table widget
func NewTable() *Table {
	return &Table{
		headers:       []string{},
		rows:          [][]string{},
		colors:        []Color{},
		theme:         DefaultTheme,
		headerVisible: true,
		borders:       true,
		alternateRows: false,
	}
}

// SetHeaders sets the table headers
func (t *Table) SetHeaders(headers ...string) *Table {
	t.headers = headers
	return t
}

// AddRow adds a row to the table
func (t *Table) AddRow(cells ...string) *Table {
	t.rows = append(t.rows, cells)
	return t
}

// AddRowWithColor adds a row with a specific color
func (t *Table) AddRowWithColor(color Color, cells ...string) *Table {
	t.rows = append(t.rows, cells)
	t.colors = append(t.colors, color)
	return t
}

// SetTheme sets the table theme
func (t *Table) SetTheme(theme *Theme) *Table {
	t.theme = theme
	return t
}

// SetAlignment sets the alignment for a column
func (t *Table) SetAlignment(col int, align Alignment) *Table {
	if len(t.align) <= col {
		newAlign := make([]Alignment, col+1)
		copy(newAlign, t.align)
		t.align = newAlign
	}
	t.align[col] = align
	return t
}

// SetBorders enables or disables borders
func (t *Table) SetBorders(enabled bool) *Table {
	t.borders = enabled
	return t
}

// SetHeaderVisible shows or hides the header
func (t *Table) SetHeaderVisible(visible bool) *Table {
	t.headerVisible = visible
	return t
}

// SetMinColWidths sets the minimum width for each column
func (t *Table) SetMinColWidths(widths ...int) *Table {
	t.minColWidths = widths
	return t
}

// Render implements the Widget interface
func (t *Table) Render(ctx context.Context, buf *Buffer) error {
	// Check context
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Calculate column widths
	colWidths := t.calculateColumnWidths()

	// Render top border
	t.renderBorder(buf, BorderTop, colWidths)

	// Render header
	if t.headerVisible && len(t.headers) > 0 {
		t.renderHeader(buf, colWidths)
		t.renderBorder(buf, BorderMiddle, colWidths)
	}

	// Render rows
	for i, row := range t.rows {
		t.renderRow(buf, row, colWidths, i)
	}

	// Render bottom border
	t.renderBorder(buf, BorderBottom, colWidths)

	return nil
}

// MinSize implements the Widget interface
func (t *Table) MinSize() Size {
	colWidths := t.calculateColumnWidths()

	// Calculate width
	width := 0
	for _, w := range colWidths {
		width += w
	}
	width += len(colWidths) + 1 // Vertical borders

	// Calculate height
	height := 0
	if t.headerVisible && len(t.headers) > 0 {
		height += 2 // Header + separator line
	}
	height += len(t.rows) // Data rows
	height += len(t.rows) // Separator lines between rows
	height += 2           // Top + bottom borders

	return Size{Width: width, Height: height}
}

// calculateColumnWidths calculates the width needed for each column
func (t *Table) calculateColumnWidths() []int {
	if len(t.headers) == 0 {
		return []int{}
	}

	widths := make([]int, len(t.headers))

	// Consider headers
	for i, h := range t.headers {
		widths[i] = stringutil.GetPrintableLength(h)
	}

	// Consider all rows
	for _, row := range t.rows {
		for i := 0; i < len(row) && i < len(widths); i++ {
			width := stringutil.GetPrintableLength(row[i])
			if width > widths[i] {
				widths[i] = width
			}
		}
	}

	// Apply min/max constraints
	for i := range widths {
		if i < len(t.minColWidths) && widths[i] < t.minColWidths[i] {
			widths[i] = t.minColWidths[i]
		}
		if i < len(t.maxColWidths) && widths[i] > t.maxColWidths[i] {
			widths[i] = t.maxColWidths[i]
		}
	}

	// Add padding (2 spaces per cell)
	for i := range widths {
		widths[i] += 2
	}

	return widths
}

// renderBorder renders a border line
func (t *Table) renderBorder(buf *Buffer, borderType BorderType, colWidths []int) {
	border := t.theme.Border

	var left, right, join string
	switch borderType {
	case BorderTop:
		left = border.TopLeft
		right = border.TopRight
		join = border.TopJoin
	case BorderMiddle:
		left = border.MiddleLeft
		right = border.MiddleRight
		join = border.Cross
	case BorderBottom:
		left = border.BottomLeft
		right = border.BottomRight
		join = border.BottomJoin
	}

	buf.WriteString(left)
	for i, width := range colWidths {
		buf.WriteString(strings.Repeat(border.Horizontal, width))
		if i < len(colWidths)-1 {
			buf.WriteString(join)
		}
	}
	buf.WriteString(right)
	buf.NewLine()
}

// renderHeader renders the table header
func (t *Table) renderHeader(buf *Buffer, colWidths []int) {
	border := t.theme.Border

	buf.WriteString(border.Vertical)
	for i, header := range t.headers {
		// Add padding
		padded := " " + header + " "
		aligned := t.alignText(padded, colWidths[i], i)

		buf.WriteStringWithStyle(aligned, t.theme.TableHeader)
		buf.WriteString(border.Vertical)
	}
	buf.NewLine()
}

// renderRow renders a table row
func (t *Table) renderRow(buf *Buffer, row []string, colWidths []int, rowIndex int) {
	border := t.theme.Border

	// Get row color
	var style Style
	if rowIndex < len(t.colors) && t.colors[rowIndex] != nil {
		style = Style{Foreground: t.colors[rowIndex]}
	} else if t.alternateRows && rowIndex%2 == 1 {
		style = t.theme.TableRowAlt
	} else {
		style = t.theme.TableRow
	}

	// Apply color to entire row (including borders)
	if style.Foreground != nil {
		buf.WriteString(style.Foreground.ANSI())
	}

	buf.WriteString(border.Vertical)
	for i, cell := range row {
		if i >= len(colWidths) {
			break
		}

		// Add padding
		padded := " " + cell + " "
		aligned := t.alignText(padded, colWidths[i], i)

		buf.WriteString(aligned)
		buf.WriteString(border.Vertical)
	}

	// Reset color at end of row
	if style.Foreground != nil {
		buf.WriteString(ColorReset.ANSI())
	}

	buf.NewLine()
}

// alignText aligns text within the given width
func (t *Table) alignText(text string, width int, colIndex int) string {
	// Get alignment for this column
	align := AlignLeft
	if colIndex < len(t.align) {
		align = t.align[colIndex]
	}

	switch align {
	case AlignLeft:
		return Pad(text, width)
	case AlignCenter:
		return PadCenter(text, width)
	case AlignRight:
		return PadLeft(text, width)
	default:
		return Pad(text, width)
	}
}

// getColor returns the color for the given row index
func (t *Table) getColor(index int) Color {
	if index < len(t.colors) {
		return t.colors[index]
	}
	return ColorWhite
}

// String returns a string representation (for debugging)
func (t *Table) String() string {
	return fmt.Sprintf("Table{headers=%d, rows=%d}", len(t.headers), len(t.rows))
}
