package tui

import (
	"testing"
)

func TestNewBuffer(t *testing.T) {
	buf := NewBuffer(80, 24)
	if buf == nil {
		t.Fatal("NewBuffer() returned nil")
	}
	if buf.Width() != 80 {
		t.Errorf("Expected width 80, got %d", buf.Width())
	}
}

func TestBufferWriteString(t *testing.T) {
	buf := NewBuffer(80, 24)
	buf.WriteString("Hello")
	buf.NewLine()
	buf.WriteString("World")

	if buf.Height() != 2 {
		t.Errorf("Expected 2 lines, got %d", buf.Height())
	}
}

func TestBufferClear(t *testing.T) {
	buf := NewBuffer(80, 24)
	buf.WriteString("Hello")
	buf.NewLine()
	buf.WriteString("World")

	buf.Clear()

	if buf.Height() != 0 {
		t.Errorf("Expected 0 lines after Clear(), got %d", buf.Height())
	}
}

func TestPad(t *testing.T) {
	tests := []struct {
		name  string
		input string
		width int
		want  int // expected length
	}{
		{"short string", "Hi", 10, 10},
		{"exact length", "Hello", 5, 5},
		{"long string", "HelloWorld", 5, 10}, // Should not truncate
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Pad(tt.input, tt.width)
			if len(result) != tt.want {
				t.Errorf("Pad() length = %d, want %d", len(result), tt.want)
			}
		})
	}
}

func TestPadCenter(t *testing.T) {
	result := PadCenter("Hi", 10)
	if len(result) != 10 {
		t.Errorf("PadCenter() length = %d, want 10", len(result))
	}
	// Check roughly centered (should have spaces on both sides)
	if result[0] != ' ' || result[len(result)-1] != ' ' {
		t.Error("PadCenter() did not center the text properly")
	}
}
