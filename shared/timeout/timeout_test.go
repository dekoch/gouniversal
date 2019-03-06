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

func TestElapsedChan(t *testing.T) {

	tests := []struct {
		inMillis int
		inEnable bool
		want     float64
	}{
		{0, true, 0.0},
		{86, true, 86.0},
		{86, false, 0.0},
		{132, true, 132.0},
	}

	var to TimeOut

	for i, tt := range tests {

		to.Start(tt.inMillis)
		to.Enable(tt.inEnable)

		<-to.ElapsedChan()

		got := to.ElapsedMillis()

		if got < tt.want-20.0 ||
			got > tt.want+20.0 {
			t.Errorf("ElapsedChan() test %d: got %f, want %f +-20.0", i, got, tt.want)
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

		time.Sleep(time.Duration(2) * time.Millisecond)

		got := to.Elapsed()

		if got != tt.want {
			t.Errorf("Enable() test %d: got %t, want %t", i, got, tt.want)
		}
	}
}

func TestStartElapsedMillis(t *testing.T) {

	var to TimeOut
	to.Start(1000)

	time.Sleep(10 * time.Millisecond)

	got := to.ElapsedMillis()

	if got < 9.5 ||
		got > 10.5 {
		t.Errorf("StartElapsedMillis() got %f, want 9.5..10.5", got)
	}

	to.Enable(false)

	got = to.ElapsedMillis()

	if got != 0.0 {
		t.Errorf("StartElapsedMillis() got %f, want 0.0", got)
	}
}
