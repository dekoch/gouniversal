package timeout

import (
	"testing"
	"time"
)

func TestElapsed(t *testing.T) {

	tests := []struct {
		in   int
		wait int
		want bool
	}{
		{-1, 0, true},
		{0, 0, true},
		{2, 1, false},
		{2, 3, true},
	}

	var to TimeOut
	to.Enable(true)

	for i, tt := range tests {

		to.SetTimeOut(tt.in)
		to.Reset()

		time.Sleep(time.Duration(tt.wait) * time.Millisecond)

		got := to.Elapsed()

		if got != tt.want {
			t.Errorf("Elapsed() test %d: got %t, want %t", i, got, tt.want)
		}
	}
}

func TestEnable(t *testing.T) {

	tests := []struct {
		in   bool
		want bool
	}{
		{false, false},
		{true, true},
		{false, false},
	}

	var to TimeOut

	for i, tt := range tests {

		to.Enable(tt.in)
		to.SetTimeOut(1)
		to.Reset()

		time.Sleep(time.Duration(1) * time.Millisecond)

		got := to.Elapsed()

		if got != tt.want {
			t.Errorf("Enable() test %d: got %t, want %t", i, got, tt.want)
		}
	}
}
