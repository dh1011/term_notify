package cmd

import "testing"

func TestDisplayValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string returns not set",
			input:    "",
			expected: "(not set)",
		},
		{
			name:     "non-empty string returns as-is",
			input:    "my-topic",
			expected: "my-topic",
		},
		{
			name:     "whitespace is returned as-is",
			input:    "  ",
			expected: "  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := displayValue(tt.input)
			if got != tt.expected {
				t.Errorf("displayValue(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestMaskToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty token returns not set",
			input:    "",
			expected: "(not set)",
		},
		{
			name:     "short token (4 chars) is fully masked",
			input:    "abcd",
			expected: "****",
		},
		{
			name:     "token exactly 8 chars is fully masked",
			input:    "abcdefgh",
			expected: "****",
		},
		{
			name:     "token 9 chars shows first and last 4",
			input:    "abcdefghi",
			expected: "abcd****fghi",
		},
		{
			name:     "long token shows first and last 4",
			input:    "tk_abcdefghijklmnop",
			expected: "tk_a****mnop",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maskToken(tt.input)
			if got != tt.expected {
				t.Errorf("maskToken(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
