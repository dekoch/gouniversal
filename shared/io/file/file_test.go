package file

import (
	"crypto/rand"
	"os"
	"reflect"
	"testing"
)

func TestWriteFile(t *testing.T) {

	path := "test/file"

	tests := []struct {
		in   string
		want bool
	}{
		{path, false},
		{"", true},
	}

	b := make([]byte, 512)

	_, err := rand.Read(b)
	if err != nil {
		t.Errorf("WriteFile() %v", err.Error())
	}

	for i, tt := range tests {

		got := true

		err = WriteFile(tt.in, b)
		if err == nil {
			got = false
		}

		if got != tt.want {
			t.Errorf("WriteFile() test %d: got: %t, want: %t", i, got, tt.want)
		}
	}

	Remove("test/")
}

func TestReadFile(t *testing.T) {

	path := "test/file"

	b := make([]byte, 512)

	_, err := rand.Read(b)
	if err != nil {
		t.Errorf("ReadFile() %v", err.Error())
	}

	err = WriteFile(path, b)
	if err != nil {
		t.Errorf("ReadFile() %v", err.Error())
	}

	tests := []struct {
		in   string
		want bool
	}{
		{path, false},
		{"", true},
		{"test/f", true},
	}

	for i, tt := range tests {

		got := true

		_, err := ReadFile(tt.in)
		if err == nil {
			got = false
		}

		if got != tt.want {
			t.Errorf("ReadFile() test %d: got: %t, want: %t", i, got, tt.want)
		}
	}

	Remove("test/")
}

func TestWriteReadFile(t *testing.T) {

	path := "test/file"

	b := make([]byte, 512)

	_, err := rand.Read(b)
	if err != nil {
		t.Errorf("WriteReadFile() %v", err.Error())
	}

	err = WriteFile(path, b)
	if err != nil {
		t.Errorf("WriteReadFile() %v", err.Error())
	}

	got, err := ReadFile(path)
	if err != nil {
		t.Errorf("WriteReadFile() %v", err.Error())
	}

	if reflect.DeepEqual(got, b) == false {
		t.Errorf("WriteReadFile() got: %s, want: %s", got, b)
	}

	Remove("test/")
}

func TestChecksum(t *testing.T) {

	path := "test/file"

	tests := []struct {
		inPath    string
		inString  string
		wantError bool
		wantSum   []byte
	}{
		{path, "", false, []byte{227, 176, 196, 66, 152, 252, 28, 20, 154, 251, 244, 200, 153, 111, 185, 36, 39, 174, 65, 228, 100, 155, 147, 76, 164, 149, 153, 27, 120, 82, 184, 85}},
		{path, "test string", false, []byte{213, 87, 156, 70, 223, 204, 127, 24, 32, 112, 19, 230, 91, 68, 228, 203, 78, 44, 34, 152, 244, 172, 69, 123, 168, 248, 39, 67, 243, 30, 147, 11}},
		{"", "", true, []byte{}},
		{"test/f", "test string", true, []byte{}},
	}

	for i, tt := range tests {

		err := WriteFile(path, []byte(tt.inString))
		if err != nil {
			t.Errorf("Checksum() %v", err.Error())
			continue
		}

		gotErr := true

		gotSum, err := Checksum(tt.inPath)
		if err == nil {
			gotErr = false
		}

		if gotErr != tt.wantError {
			t.Errorf("Checksum() test %d: got: %t, want: %t", i, gotErr, tt.wantError)
			continue
		}

		if len(gotSum) == 0 &&
			len(tt.wantSum) == 0 {
			continue
		}

		if reflect.DeepEqual(gotSum, tt.wantSum) == false {
			t.Errorf("Checksum() test %d: got: %v, want: %v", i, gotSum, tt.wantSum)
		}
	}

	Remove("test/")
}

func TestRemove(t *testing.T) {

	path := "test/file"

	b := make([]byte, 512)

	_, err := rand.Read(b)
	if err != nil {
		t.Errorf("Remove() %v", err.Error())
	}

	err = WriteFile(path, b)
	if err != nil {
		t.Errorf("Remove() %v", err.Error())
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Remove() %v", err.Error())
	}

	tests := []struct {
		in   string
		want bool
	}{
		{path, false},
		{"", true},
		{"test/f", true},
	}

	for i, tt := range tests {

		Remove(tt.in)

		if _, err := os.Stat(tt.in); os.IsExist(err) {
			t.Errorf("Remove() test %d: %v", i, err.Error())
		}
	}

	Remove("test/")
}
