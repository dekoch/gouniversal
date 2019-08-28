package functions

import (
	"net/http"
	"net/url"
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

func TestCheckFormInput(t *testing.T) {

	tests := []struct {
		inQuery    string
		inKey      string
		wantString string
		wantError  bool
	}{
		{"key0=value0", "key0", "value0", false},
		{"key1=value1 z", "key1", "value1 z", false},
		{"key2=value2 z<", "key2", "", true},
	}

	for i, tt := range tests {

		v, _ := url.ParseQuery(tt.inQuery)

		var r http.Request
		r.Form = v

		gotString, gotError := CheckFormInput(tt.inKey, &r)
		if gotString != tt.wantString {
			t.Errorf("CheckFormInput() string test %d: inQuery %s, inKey %s, got %s, want %s", i, tt.inQuery, tt.inKey, gotString, tt.wantString)
		}

		if (gotError == nil) == tt.wantError {
			t.Errorf("CheckFormInput() error test %d: inQuery %s, inKey %s, got %v, want %t", i, tt.inQuery, tt.inKey, gotError, tt.wantError)
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
