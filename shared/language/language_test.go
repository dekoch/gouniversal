package language

import (
	"reflect"
	"testing"

	"github.com/dekoch/gouniversal/shared/io/file"
)

type langFile struct {
	Text string
}

func TestNewListNames(t *testing.T) {

	tests := []struct {
		in   string
		want []string
	}{
		{"en", []string{"en"}},
		{"de", []string{"de", "en"}},
	}

	var def langFile

	for i, tt := range tests {

		lang := New("test/lang/", def, tt.in)

		got := lang.ListNames()

		if reflect.DeepEqual(got, tt.want) == false {
			t.Errorf("NewListNames() test %d: got %s, want %s", i, got, tt.want)
		}
	}

	file.Remove("test/")
}

func TestSelectLang(t *testing.T) {

	tests := []struct {
		in   string
		want string
	}{
		{"en", "hello"},
		{"de", "hallo"},
		{"fr", "allô"},
		{"cn", "hello"},
	}

	var (
		def langFile
		lf  langFile
	)

	def.Text = "hello"
	langDef := New("test/lang/", def, "en")

	lf.Text = "hallo"
	langRead := New("test/lang/", lf, "de")

	lf.Text = "allô"
	langRead = New("test/lang/", lf, "fr")

	langRead.SelectLang("en", &lf)

	for i, tt := range tests {

		langDef.SelectLang(tt.in, &lf)

		got := lf.Text

		if reflect.DeepEqual(got, tt.want) == false {
			t.Errorf("SelectLang() test %d: got %s, want %s", i, got, tt.want)
		}
	}

	file.Remove("test/")
}
