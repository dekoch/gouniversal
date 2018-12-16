package stringarray

import (
	"reflect"
	"testing"
)

func TestAddAndList(t *testing.T) {

	tests := []struct {
		in   string
		want []string
	}{
		{"", []string{""}},
		{"111", []string{"", "111"}},
		{"222", []string{"", "111", "222"}},
	}

	var a StringArray

	for _, tt := range tests {

		a.Add(tt.in)

		got := a.List()

		if reflect.DeepEqual(got, tt.want) == false {
			t.Errorf("AddAndList(): got %v, want %v", got, tt.want)
		}
	}
}

func TestCount(t *testing.T) {

	tests := []struct {
		in   string
		want int
	}{
		{"", 1},
		{"111", 2},
		{"222", 3},
	}

	var a StringArray

	for _, tt := range tests {

		a.Add(tt.in)

		got := a.Count()

		if got != tt.want {
			t.Errorf("Count(): got %d, want %d", got, tt.want)
		}
	}
}

func TestRemove(t *testing.T) {

	list := make([]string, 3)
	list[0] = ""
	list[1] = "111"
	list[2] = "222"

	tests := []struct {
		in   string
		want []string
	}{
		{"", []string{"", "111", "222"}},
		{"111", []string{"", "222"}},
		{"", []string{"", "222"}},
	}

	var a StringArray
	a.AddList(list)

	for _, tt := range tests {

		if len(tt.in) > 0 {
			a.Remove(tt.in)
		}

		got := a.List()

		if reflect.DeepEqual(got, tt.want) == false {
			t.Errorf("Remove(): got %v, want %v", got, tt.want)
		}
	}
}

func TestRemoveAll(t *testing.T) {

	list := make([]string, 3)
	list[0] = ""
	list[1] = "111"
	list[2] = "222"

	var a StringArray
	a.AddList(list)

	got := a.Count()

	if got != 3 {
		t.Errorf("RemoveAll(): got %d, want 3", got)
	}

	a.RemoveAll()

	got = a.Count()

	if got != 0 {
		t.Errorf("RemoveAll(): got %d, want 0", got)
	}
}
