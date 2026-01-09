package tui

import (
	"io"
	"strings"

	"github.com/yudgnahk/euro21/utils/stringutil"
)

// Buffer is an in-memory render target
type Buffer struct {
	lines   []string
	cursorX int
	cursorY int
	width   int
	style   Style
}

// NewBuffer creates a new buffer with the specified dimensions
func NewBuffer(width, height int) *Buffer {
	lines := make([]string, 0, height)
	return &Buffer{
		lines: lines,
		width: width,
	}
}

// WriteString writes a string at the current cursor position
func (b *Buffer) WriteString(s string) {
	// Ensure we have enough lines
	for len(b.lines) <= b.cursorY {
		b.lines = append(b.lines, "")
	}

	// Append to current line
	b.lines[b.cursorY] += s
}

// WriteStringWithStyle writes a string with the given style
func (b *Buffer) WriteStringWithStyle(s string, style Style) {
	if style.Foreground != nil {
		b.WriteString(style.Foreground.ANSI())
	}
	b.WriteString(s)
	if style.Foreground != nil {
		b.WriteString(ColorReset.ANSI())
	}
}

// NewLine moves to the next line
func (b *Buffer) NewLine() {
	b.cursorY++
	b.cursorX = 0
}

// MoveCursor moves the cursor to the specified position
func (b *Buffer) MoveCursor(x, y int) {
	b.cursorX = x
	b.cursorY = y
}

// Flush writes the buffer contents to the writer
func (b *Buffer) Flush(w io.Writer) error {
	var sb strings.Builder
	for i, line := range b.lines {
		sb.WriteString(line)
		if i < len(b.lines)-1 {
			sb.WriteString("\n")
		}
	}
	// Add final newline
	sb.WriteString("\n")

	_, err := w.Write([]byte(sb.String()))
	return err
}

// Clear resets the buffer
func (b *Buffer) Clear() {
	b.lines = b.lines[:0]
	b.cursorX = 0
	b.cursorY = 0
}

// Width returns the buffer width
func (b *Buffer) Width() int {
	return b.width
}

// Height returns the current buffer height (number of lines)
func (b *Buffer) Height() int {
	return len(b.lines)
}

// Pad adds padding to a string to reach the desired width
func Pad(s string, width int) string {
	actualLen := stringutil.GetPrintableLength(s)
	if actualLen >= width {
		return s
	}
	return s + strings.Repeat(" ", width-actualLen)
}

// PadLeft adds left padding to a string
func PadLeft(s string, width int) string {
	actualLen := stringutil.GetPrintableLength(s)
	if actualLen >= width {
		return s
	}
	return strings.Repeat(" ", width-actualLen) + s
}

// PadCenter centers a string within the given width
func PadCenter(s string, width int) string {
	actualLen := stringutil.GetPrintableLength(s)
	if actualLen >= width {
		return s
	}
	leftPad := (width - actualLen) / 2
	rightPad := width - actualLen - leftPad
	return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
}
