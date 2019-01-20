package sbool

import (
	"testing"
)

func TestSetUnSetAndIsSet(t *testing.T) {

	tests := []struct {
		in   bool
		want bool
	}{
		{true, true},
		{false, false},
		{true, true},
		{false, false},
	}

	var b Sbool

	for _, tt := range tests {

		if tt.in {
			b.Set()
		} else {
			b.UnSet()
		}

		got := b.IsSet()

		if got != tt.want {
			t.Errorf("SetUnSetAndIsSet(): got %t, want %t", got, tt.want)
		}
	}
}

func TestSetStateAndIsSet(t *testing.T) {

	tests := []struct {
		in   bool
		want bool
	}{
		{true, true},
		{false, false},
		{true, true},
		{false, false},
	}

	var b Sbool

	for _, tt := range tests {

		b.SetState(tt.in)

		got := b.IsSet()

		if got != tt.want {
			t.Errorf("SetStateAndIsSet(): got %t, want %t", got, tt.want)
		}
	}
}
