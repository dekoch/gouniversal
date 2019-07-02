package sitemap

import (
	"reflect"
	"testing"
)

func TestRegisterGetList(t *testing.T) {

	tests := []struct {
		inMenu   string
		inPath   string
		inTitle  string
		wantList []string
	}{
		// missing values
		{"", "", "", []string{}},
		{"fooMenu0", "", "", []string{}},
		{"", "fooPath1", "", []string{}},
		{"", "", "fooTitle2", []string{}},
		// required values
		{"", "fooPath3", "fooTitle3", []string{"fooPath3"}},
		{"fooMenu4", "fooPath4", "fooTitle4", []string{"fooPath4", "fooPath3"}},
	}

	var s Sitemap

	for _, tt := range tests {
		s.Register(tt.inMenu, tt.inPath, tt.inTitle)

		gotList := s.PageList()
		if len(gotList) > 0 ||
			len(tt.wantList) > 0 {

			if reflect.DeepEqual(gotList, tt.wantList) == false {
				t.Errorf("PageList(): got %v, want %v", gotList, tt.wantList)
			}
		}
	}
}

func TestRegisterGetTitle(t *testing.T) {

	tests := []struct {
		inMenu    string
		inPath    string
		inTitle   string
		testPath  string
		wantTitle string
	}{
		{"", "fooPath0", "fooTitle0", "fooPath0", "fooTitle0"},
		{"", "fooPath1", "fooTitle1", "fooPath1:Par1", "fooTitle1"},
		{"", "fooPath2", "fooTitle2", "fooPathX", ""},
	}

	var s Sitemap

	for _, tt := range tests {
		s.Register(tt.inMenu, tt.inPath, tt.inTitle)

		gotTitle := s.PageTitle(tt.testPath)
		if gotTitle != tt.wantTitle {
			t.Errorf("PageTitle(): got %s, want %s", gotTitle, tt.wantTitle)
		}
	}
}

func TestShowMap(t *testing.T) {

	tests := []struct {
		inMenu  string
		inPath  string
		inTitle string
	}{
		{"", "fooPath0", "fooTitle0"},
		{"", "fooPath1", "fooTitle1"},
	}

	var s Sitemap

	for _, tt := range tests {
		s.Register(tt.inMenu, tt.inPath, tt.inTitle)
		s.ShowMap()
	}
}

func TestClear(t *testing.T) {

	var s Sitemap
	s.Register("fooMenu", "fooPath0", "fooTitle0")
	s.Register("fooMenu", "fooPath1", "fooTitle1")

	gotList := s.PageList()
	if len(gotList) != 2 {
		t.Errorf("TestClear(): got %d, want 2", len(gotList))
	}

	s.Clear()

	gotList = s.PageList()
	if len(gotList) != 0 {
		t.Errorf("TestClear(): got %d, want 0", len(gotList))
	}
}
