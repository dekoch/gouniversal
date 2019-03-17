package navigation

import "testing"

func TestGetNextPage(t *testing.T) {

	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"App:Module:Settings$Par1=1", "App"},
		{"App:Module:Settings$Par1=1", "Module"},
		{"App:Module:Settings$Par1=1", "Settings"},
	}

	var nav Navigation

	for i, tt := range tests {

		nav.Path = tt.in

		got := nav.GetNextPage()

		if got != tt.want {
			t.Errorf("GetNextPage() test %d: got %v, want %v", i, got, tt.want)
		}
	}
}

func TestIsNext(t *testing.T) {

	tests := []struct {
		inPath string
		inPage string
		want   bool
	}{
		{"", "App", false},
		{"App:Module:Settings$Par1=1", "", false},
		{"App:Module:Settings$Par1=1", "App", true},
		{"App:Module:Settings$Par1=1", "App", false},
		{"App:Module:Settings$Par1=1", "Module", true},
		{"App:Module:Settings$Par1=1", "Module", false},
		{"App:Module:Settings$Par1=1", "Settings", true},
		{"App:Module:Settings$Par1=1", "Settings", false},
	}

	var nav Navigation

	for i, tt := range tests {

		nav.Path = tt.inPath

		got := nav.IsNext(tt.inPage)

		if got != tt.want {
			t.Errorf("IsNext() test %d: got %v, want %v", i, got, tt.want)
		}
	}
}

func TestNavigatePath(t *testing.T) {

	tests := []struct {
		in   string
		want string
	}{
		{"", "Program:Home"},
		{"Account:Login", "Account:Login"},
		{"", "Account:Login"},
	}

	var nav Navigation
	nav.Path = "Program:Home"

	for i, tt := range tests {

		nav.NavigatePath(tt.in)

		got := nav.Path

		if got != tt.want {
			t.Errorf("NavigatePath() test %d: got %v, want %v", i, got, tt.want)
		}
	}
}

func TestRedirectPath(t *testing.T) {

	tests := []struct {
		inRedirect  string
		inPath      string
		inOverwrite bool
		want        string
	}{
		{"App:Module", "App:Module:Settings", true, "App:Module:Settings"},
		{"App:Module", "", true, "App:Module"},
		{"", "App:Module:Settings", false, "App:Module:Settings"},
	}

	var nav Navigation

	for i, tt := range tests {

		nav.Redirect = tt.inRedirect

		nav.RedirectPath(tt.inPath, tt.inOverwrite)

		got := nav.Redirect

		if got != tt.want {
			t.Errorf("RedirectPath() test %d: got %v, want %v", i, got, tt.want)
		}
	}
}

func TestParameter(t *testing.T) {

	tests := []struct {
		inPath      string
		inParameter string
		want        string
	}{
		{"App:Module:Settings$Par1=1", "Par", ""},
		{"App:Module:Settings$Par1=1", "Par1", "1"},
		{"App:Module:Settings$Par1=1$Par2=2", "Par1", "1"},
		{"App:Module:Settings$Par1=1$Par2=2", "Par2", "2"},
		{"App:Module:Settings$Par1=$Par2=2", "Par1", ""},
	}

	var nav Navigation

	for i, tt := range tests {

		nav.Path = tt.inPath

		got := nav.Parameter(tt.inParameter)

		if got != tt.want {
			t.Errorf("Parameter() test %d: got %v, want %v", i, got, tt.want)
		}
	}
}
