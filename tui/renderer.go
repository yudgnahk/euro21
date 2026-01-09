package tui

import (
	"context"
	"fmt"
	"io"
)

// Renderer orchestrates widget rendering
type Renderer struct {
	writer   io.Writer
	terminal *Terminal
	theme    *Theme
}

// NewRenderer creates a new renderer that writes to the given writer
func NewRenderer(w io.Writer) *Renderer {
	return &Renderer{
		writer:   w,
		terminal: Detect(),
		theme:    DefaultTheme,
	}
}

// SetTheme sets the theme for rendering
func (r *Renderer) SetTheme(theme *Theme) {
	r.theme = theme
}

// GetTheme returns the current theme
func (r *Renderer) GetTheme() *Theme {
	return r.theme
}

// Render renders a widget to the output
func (r *Renderer) Render(ctx context.Context, widget Widget) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Calculate required size
	minSize := widget.MinSize()

	// Adjust for terminal width if needed
	width := minSize.Width
	if width > r.terminal.Width {
		width = r.terminal.Width
	}

	// Create buffer
	buf := NewBuffer(width, minSize.Height)

	// Render widget to buffer
	if err := widget.Render(ctx, buf); err != nil {
		return fmt.Errorf("render failed: %w", err)
	}

	// Flush buffer to writer
	return buf.Flush(r.writer)
}

// RenderMany renders multiple widgets sequentially
func (r *Renderer) RenderMany(ctx context.Context, widgets ...Widget) error {
	for _, widget := range widgets {
		if err := r.Render(ctx, widget); err != nil {
			return err
		}
	}
	return nil
}
