package utils

import "testing"

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		length   int
		expected string
	}{
		{"cat", 10, "cat"},
		{"poppy", 5, "poppy"},
		{"discord bot", 7, "discord"},
		{"", 3, ""},
	}

	for _, test := range tests {
		result := Truncate(test.input, test.length)
		if result != test.expected {
			t.Errorf("Truncate(%q, %d) = %q; want %q", test.input, test.length, result, test.expected)
		}
	}
}
