package csv

import (
	"reflect"
	"testing"

	"github.com/dekoch/gouniversal/shared/io/file"
)

func TestAddRow(t *testing.T) {

	path := "test/file"

	tests := []struct {
		inPath string
		inData []string
		want   bool
	}{
		{path, []string{"11", "12", "13"}, false},
		{"", []string{"11", "12", "13"}, true},
	}

	for i, tt := range tests {

		got := true

		err := AddRow(tt.inPath, tt.inData)
		if err == nil {
			got = false
		}

		if got != tt.want {
			t.Errorf("AddRow() test %d: got: %t, want: %t", i, got, tt.want)
		}
	}

	file.Remove("test/")
}

func TestReadAll(t *testing.T) {

	path := "test/file"

	tests := []struct {
		inPath    string
		inData    []string
		wantError bool
		wantData  [][]string
	}{
		{path, []string{"11", "12", "13"}, false, [][]string{[]string{"11", "12", "13"}}},
		{path, []string{"21", "22", "23"}, false, [][]string{[]string{"11", "12", "13"}, []string{"21", "22", "23"}}},
		{"", []string{"11", "12", "13"}, true, [][]string{}},
		{"test/f", []string{"11", "12", "13"}, true, [][]string{}},
	}

	for i, tt := range tests {

		err := AddRow(path, tt.inData)
		if err != nil {
			t.Errorf("ReadAll() test %d: %v", i, err.Error())
			continue
		}

		gotError := true

		gotData, err := ReadAll(tt.inPath)
		if err == nil {
			gotError = false
		}

		if gotError != tt.wantError {
			t.Errorf("ReadAll() test %d: got: %t, want: %t", i, gotError, tt.wantError)
		}

		if len(gotData) == 0 &&
			len(tt.wantData) == 0 {
			continue
		}

		if reflect.DeepEqual(gotData, tt.wantData) == false {
			t.Errorf("ReadAll() got: %s, want: %s", gotData, tt.wantData)
		}
	}

	file.Remove("test/")
}

func TestAddRowReadAll(t *testing.T) {

	path := "test/file"

	tests := []struct {
		in   []string
		want [][]string
	}{
		{[]string{"11", "12", "13"}, [][]string{[]string{"11", "12", "13"}}},
		{[]string{"21", "22", "23"}, [][]string{[]string{"11", "12", "13"}, []string{"21", "22", "23"}}},
	}

	for i, tt := range tests {

		err := AddRow(path, tt.in)
		if err != nil {
			t.Errorf("AddRowReadAll() test %d: %v", i, err.Error())
			continue
		}

		got, err := ReadAll(path)
		if err != nil {
			t.Errorf("AddRowReadAll() test %d: %v", i, err.Error())
			continue
		}

		if reflect.DeepEqual(got, tt.want) == false {
			t.Errorf("AddRowReadAll() got: %s, want: %s", got, tt.want)
		}
	}

	file.Remove("test/")
}
