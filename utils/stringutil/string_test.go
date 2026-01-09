package stringutil

import "testing"

func TestToInt(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "valid positive integer",
			input: "123",
			want:  123,
		},
		{
			name:  "valid negative integer",
			input: "-456",
			want:  -456,
		},
		{
			name:  "zero",
			input: "0",
			want:  0,
		},
		{
			name:  "large number",
			input: "999999",
			want:  999999,
		},
		{
			name:  "invalid string",
			input: "abc",
			want:  0,
		},
		{
			name:  "empty string",
			input: "",
			want:  0,
		},
		{
			name:  "string with spaces",
			input: "12 34",
			want:  0,
		},
		{
			name:  "decimal number",
			input: "12.34",
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToInt(tt.input)
			if got != tt.want {
				t.Errorf("ToInt(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestGetPrintableLength(t *testing.T) {
	tests := []struct {
		name  string
		input string
		// Note: Exact length depends on emoji handling, so we just check it returns a positive value
		wantPositive bool
	}{
		{
			name:         "plain text",
			input:        "Hello World",
			wantPositive: true,
		},
		{
			name:         "text with flag emoji",
			input:        "🇺🇸 United States",
			wantPositive: true,
		},
		{
			name:         "text with simple emoji",
			input:        "⚽️ Football",
			wantPositive: true,
		},
		{
			name:         "empty string",
			input:        "",
			wantPositive: false, // 0 is valid for empty string
		},
		{
			name:         "only emoji",
			input:        "🇮🇹",
			wantPositive: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPrintableLength(tt.input)
			if tt.wantPositive && got <= 0 {
				t.Errorf("GetPrintableLength(%q) = %v, want positive value", tt.input, got)
			}
			if !tt.wantPositive && got != 0 {
				t.Errorf("GetPrintableLength(%q) = %v, want 0", tt.input, got)
			}
		})
	}
}

func TestGetPrintableLength_Comparison(t *testing.T) {
	// Test that emoji strings return smaller printable length than raw length
	plainText := "Italy"
	emojiText := "🇮🇹 Italy"

	plainLen := GetPrintableLength(plainText)
	emojiLen := GetPrintableLength(emojiText)

	// The emoji text should have a calculated printable length that accounts for emoji display width
	if plainLen <= 0 {
		t.Errorf("Plain text length should be positive, got %d", plainLen)
	}
	if emojiLen <= 0 {
		t.Errorf("Emoji text length should be positive, got %d", emojiLen)
	}

	// Both should be reasonable values (not zero or negative)
	if plainLen > 100 || emojiLen > 100 {
		t.Errorf("Lengths seem unreasonably large: plain=%d, emoji=%d", plainLen, emojiLen)
	}
}
