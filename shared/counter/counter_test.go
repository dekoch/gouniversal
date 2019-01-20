package counter

import (
	"testing"
)

func TestSetCountAndGetCount(t *testing.T) {

	tests := []struct {
		in   int
		want int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{-1, 0},
	}

	var c Counter

	for _, tt := range tests {

		c.SetCount(tt.in)

		got := c.GetCount()

		if got != tt.want {
			t.Errorf("SetCountAndGetCount(): got %d, want %d", got, tt.want)
		}
	}
}

func TestAddAndGetCount(t *testing.T) {

	tests := []struct {
		in   int
		want int
	}{
		{0, 0},
		{1, 1},
		{2, 3},
		{0, 3},
	}

	var c Counter

	for _, tt := range tests {

		for i := 0; i < tt.in; i++ {
			c.Add()
		}

		got := c.GetCount()

		if got != tt.want {
			t.Errorf("AddAndGetCount(): got %d, want %d", got, tt.want)
		}
	}
}

func TestAddCountAndGetCount(t *testing.T) {

	tests := []struct {
		in   int
		want int
	}{
		{0, 0},
		{1, 1},
		{2, 3},
		{0, 3},
	}

	var c Counter

	for _, tt := range tests {

		c.AddCount(tt.in)

		got := c.GetCount()

		if got != tt.want {
			t.Errorf("AddCountAndGetCount(): got %d, want %d", got, tt.want)
		}
	}
}

func TestRemoveAndGetCount(t *testing.T) {

	// start with 10
	tests := []struct {
		in   int
		want int
	}{
		{0, 10},
		{1, 9},
		{2, 7},
		{0, 7},
		{7, 0},
		{1, 0},
	}

	var c Counter
	c.AddCount(10)

	for _, tt := range tests {

		for i := 0; i < tt.in; i++ {
			c.Remove()
		}

		got := c.GetCount()

		if got != tt.want {
			t.Errorf("RemoveAndGetCount(): got %d, want %d", got, tt.want)
		}
	}
}

func TestRemoveCountAndGetCount(t *testing.T) {

	// start with 10
	tests := []struct {
		in   int
		want int
	}{
		{0, 10},
		{1, 9},
		{2, 7},
		{0, 7},
		{7, 0},
		{1, 0},
	}

	var c Counter
	c.AddCount(10)

	for _, tt := range tests {

		c.RemoveCount(tt.in)

		got := c.GetCount()

		if got != tt.want {
			t.Errorf("RemoveCountAndGetCount(): got %d, want %d", got, tt.want)
		}
	}
}

func TestResetAndGetCount(t *testing.T) {

	var c Counter
	c.AddCount(10)

	got := c.GetCount()

	if got != 10 {
		t.Errorf("ResetAndGetCount(): got %d, want 10", got)
	}

	c.Reset()

	got = c.GetCount()

	if got != 0 {
		t.Errorf("ResetAndGetCount(): got %d, want 0", got)
	}
}
