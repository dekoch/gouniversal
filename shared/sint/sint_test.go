package sint

import (
	"testing"
)

func TestSetAndGet(t *testing.T) {

	tests := []struct {
		in   int
		want int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{-1, -1},
	}

	var si Sint

	for _, tt := range tests {

		si.Set(tt.in)

		got := si.Get()

		if got != tt.want {
			t.Errorf("SetAndGet(): got %d, want %d", got, tt.want)
		}
	}
}
