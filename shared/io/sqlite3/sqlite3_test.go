package sqlite3

import (
	"testing"

	"github.com/dekoch/gouniversal/shared/io/file"
)

func TestOpen(t *testing.T) {

	tests := []struct {
		in    string
		close bool
		want  bool
	}{
		{"", true, true},
		{" ", true, true},
		{"test/foo.db", true, false},
		{"test/foo.db", false, false},
		{"test/foo.db", true, true},
	}

	var sq SQLite

	for i, tt := range tests {
		got := sq.Open(tt.in)

		if tt.close {
			sq.Close()
		}

		if (got == nil) == tt.want {
			t.Errorf("Open() test %d: in %s, got %v, want %t", i, tt.in, got, tt.want)
		}
	}

	file.Remove("test/")
}

func TestClose(t *testing.T) {

	tests := []struct {
		open bool
		want bool
	}{
		{true, false},
		{true, false},
		{false, true},
		{true, false},
	}

	var sq SQLite

	for i, tt := range tests {

		if tt.open {
			sq.Open("test/foo.db")
		}

		got := sq.Close()

		if (got == nil) == tt.want {
			t.Errorf("Close() test %d: open %t, got %v, want %t", i, tt.open, got, tt.want)
		}
	}

	file.Remove("test/")
}

func TestCreateTable(t *testing.T) {

	tests := []struct {
		name  string
		table string
		want  bool
	}{
		{"", "id TEXT PRIMARY KEY, value TEXT", true},
		{" ", "id TEXT PRIMARY KEY, value TEXT", true},
		{"table0", "", true},
		{"table1", " ", true},
		{"table2", "id TEXT PRIMARY KEY, value TEXT", false},
	}

	var sq SQLite
	sq.Open("test/foo.db")

	defer sq.Close()

	for i, tt := range tests {

		got := sq.CreateTable(tt.name, tt.table)

		if (got == nil) == tt.want {
			t.Errorf("CreateTable() test %d: name %s, table %s, got %v, want %t", i, tt.name, tt.table, got, tt.want)
		}

		if tt.want == false {
			exists, _ := sq.TableExists(tt.name)
			if exists == false {
				t.Errorf("CreateTable() test %d: table does not exist", i)
			}
		}
	}

	file.Remove("test/")
}

/*
func TestCreateTableFromLayout(t *testing.T) {

	tests := []struct {
		name   string
		fields []field
		want   bool
	}{
		{"", []field{{"", TypeTEXT, false, false}}, true},
	}

	var sq SQLite
	sq.Open("test/foo.db")

	defer sq.Close()

	var lyt Layout

	for i, tt := range tests {

		got := sq.CreateTableFromLayout(lyt)

		if (got == nil) == tt.want {
			t.Errorf("CreateTableFromLayout() test %d: name %s, table %s, got %v, want %t", i, tt.name, tt.table, got, tt.want)
		}

		if tt.want == false {
			exists, _ := sq.TableExists(tt.name)
			if exists == false {
				t.Errorf("CreateTableFromLayout() test %d: table does not exist", i)
			}
		}
	}

	file.Remove("test/")
}
*/

func TestTableExists(t *testing.T) {

	tests := []struct {
		name      string
		wantRet   bool
		wantError bool
	}{
		{"", false, true},
		{" ", false, true},
		{"table0", true, false},
	}

	var sq SQLite
	sq.Open("test/foo.db")

	defer sq.Close()

	for i, tt := range tests {

		sq.CreateTable(tt.name, "id TEXT PRIMARY KEY, value TEXT")

		gotRet, gotError := sq.TableExists(tt.name)

		if gotRet != tt.wantRet {
			t.Errorf("TableExists() return test %d: name %s, got %t, want %t", i, tt.name, gotRet, tt.wantRet)
		}

		if (gotError == nil) == tt.wantError {
			t.Errorf("TableExists() error test %d: name %s, got %v, want %t", i, tt.name, gotError, tt.wantError)
		}
	}

	file.Remove("test/")
}

func TestDropTable(t *testing.T) {

	tests := []struct {
		in     string
		create bool
		want   bool
	}{
		{"", false, true},
		{" ", false, true},
		{"table0", true, false},
	}

	var sq SQLite
	sq.Open("test/foo.db")

	defer sq.Close()

	for i, tt := range tests {

		if tt.create {
			sq.CreateTable(tt.in, "id TEXT PRIMARY KEY, value TEXT")

			exists, _ := sq.TableExists(tt.in)
			if exists == false {
				t.Errorf("DropTable() test %d: table does not exist", i)
			}
		}

		got := sq.DropTable(tt.in)

		if (got == nil) == tt.want {
			t.Errorf("DropTable() test %d: in %s, got %v, want %t", i, tt.in, got, tt.want)
		}

		if tt.create {
			exists, _ := sq.TableExists(tt.in)
			if exists {
				t.Errorf("DropTable() test %d: table exists", i)
			}
		}
	}

	file.Remove("test/")
}
