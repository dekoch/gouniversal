package token

import (
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {

	tests := []struct {
		test int
		in   string
	}{
		{0, ""},
		{1, "1234"},
		{2, "1234"},
	}

	var (
		ut    Token
		token string
	)

	for _, tt := range tests {

		got := ut.New(tt.in)

		switch tt.test {
		case 0:
			if got != "" {
				t.Errorf("New(): got %v, want \"\"", got)
			}

		case 1:
			token = got

			if got == "" {
				t.Errorf("New(): got %v, want \"...\"", got)
			}

		case 2:
			if got == token {
				t.Errorf("New(): got %v, want %v", got, token)
			}
		}
	}
}

func TestCheck(t *testing.T) {

	tests := []struct {
		inNew   string
		inCheck string
		want    bool
	}{
		{"", "", false},
		{"", "1234", false},
		{"1234", "", false},
		{"1234", "1111", false},
		{"1234", "1234", true},
	}

	var ut Token

	for _, tt := range tests {

		token := ut.New(tt.inNew)

		got := ut.Check(tt.inCheck, token)

		if got != tt.want {
			t.Errorf("Check(): got %t, want %t", got, tt.want)
		}

		got = ut.Check(tt.inCheck, "2222")

		if got != false {
			t.Errorf("Check(): got %t, want false", got)
		}
	}
}

func TestRemove(t *testing.T) {

	tests := []struct {
		inNew    string
		inRemove string
		want     bool
	}{
		{"1", "", true},
		{"2", "2", false},
		{"3", "3333", true},
	}

	var ut Token
	var token [3]string

	for i, tt := range tests {

		token[i] = ut.New(tt.inNew)
	}

	for i, tt := range tests {

		got := ut.Check(tt.inNew, token[i])

		if got != true {
			t.Errorf("Remove(): got %t, want true", got)
		}

		ut.Remove(tt.inRemove)

		got = ut.Check(tt.inNew, token[i])

		if got != tt.want {
			t.Errorf("Remove(): got %t, want %t", got, tt.want)
		}
	}
}

func TestSetMaxTokens(t *testing.T) {

	tests := []struct {
		in   int
		want bool
	}{
		{0, true},
		{4, true},
		{3, false},
	}

	for i, tt := range tests {

		var ut Token
		ut.SetMaxTokens(tt.in)

		l := make([]string, 0)

		for i := 0; i < 5; i++ {

			l = append(l, ut.New(strconv.Itoa(i)))
		}

		got := ut.Check("0", l[0])

		if got != tt.want {
			t.Errorf("SetMaxTokens(): test %d: got %t, want %t", i, got, tt.want)
		}
	}
}
