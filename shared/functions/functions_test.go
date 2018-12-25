package functions

import (
	"testing"
)

func TestIsEmtpy(t *testing.T) {

	tests := []struct {
		in   string
		want bool
	}{
		{"", true},
		{" ", true},
		{"  ", true},
		{"\r", true},
		{"\n", true},
		{"foo", false},
	}

	for _, tt := range tests {
		got := IsEmpty(tt.in)
		if got != tt.want {
			t.Errorf("IsEmtpy(): in: %s, got %t, want %t", tt.in, got, tt.want)
		}
	}
}

func TestRound(t *testing.T) {

	tests := []struct {
		inVal     float64
		inRoundOn float64
		inPlaces  int
		want      float64
	}{
		{0, .5, 0, 0},
		{.4, .5, 0, 0},
		{.5, .5, 0, 1},
		{.44, .5, 1, .4},
		{.45, .5, 1, .5},
		{1, .5, 0, 1},
	}

	for _, tt := range tests {
		got := Round(tt.inVal, tt.inRoundOn, tt.inPlaces)
		if got != tt.want {
			t.Errorf("IsEmtpy(): in %g %g %d, got %g, want %g", tt.inVal, tt.inRoundOn, tt.inPlaces, got, tt.want)
		}
	}
}
