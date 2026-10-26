package main

import (
	"testing"
)

func TestPackageToSelection(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"quick", "1"},
		{"1", "1"},
		{"deep", "2"},
		{"2", "2"},
		{"look", "3"},
		{"3", "3"},
		{"full", "4"},
		{"4", "4"},
		{"QUICK", "1"},
		{"Deep", "2"},
		{"invalid", ""},
		{"", ""},
		{"5", ""},
	}

	for _, tt := range tests {
		got := packageToSelection(tt.input)
		if got != tt.want {
			t.Errorf("packageToSelection(%q) = %q; want %q", tt.input, got, tt.want)
		}
	}
}
