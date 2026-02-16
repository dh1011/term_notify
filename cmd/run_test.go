package cmd

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Duration
		expected string
	}{
		{
			name:     "zero duration",
			input:    0,
			expected: "0.0s",
		},
		{
			name:     "sub-second",
			input:    500 * time.Millisecond,
			expected: "0.5s",
		},
		{
			name:     "one second",
			input:    1 * time.Second,
			expected: "1.0s",
		},
		{
			name:     "just under a minute",
			input:    59*time.Second + 900*time.Millisecond,
			expected: "59.9s",
		},
		{
			name:     "exactly one minute",
			input:    1 * time.Minute,
			expected: "1m 0s",
		},
		{
			name:     "minutes and seconds",
			input:    2*time.Minute + 30*time.Second,
			expected: "2m 30s",
		},
		{
			name:     "just under an hour",
			input:    59*time.Minute + 59*time.Second,
			expected: "59m 59s",
		},
		{
			name:     "exactly one hour",
			input:    1 * time.Hour,
			expected: "1h 0m 0s",
		},
		{
			name:     "hours minutes seconds",
			input:    1*time.Hour + 15*time.Minute + 30*time.Second,
			expected: "1h 15m 30s",
		},
		{
			name:     "multiple hours",
			input:    3*time.Hour + 5*time.Minute + 12*time.Second,
			expected: "3h 5m 12s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.input)
			if got != tt.expected {
				t.Errorf("formatDuration(%v) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
