package fileinfo

import (
	"strconv"
	"testing"

	"github.com/dekoch/gouniversal/shared/io/file"
)

func TestGet(t *testing.T) {

	tests := []struct {
		fileNo       int
		wantName     string
		wantByteSize int64
		wantIsDir    bool
		wantPath     string
	}{
		{0, "file0", 0, false, "test/"},
		{1, "file1", 1, false, "test/"},
		{4, "file4", 4, false, "test/"},
		{5, "foo", -1, true, "test/"},
		{8, "file2", 2, false, "test/foo/"},
	}

	err := writeFiles("test/", 5)
	if err != nil {
		t.Errorf("Get() error %v", err)
		return
	}

	err = writeFiles("test/foo/", 3)
	if err != nil {
		t.Errorf("Get() error %v", err)
		return
	}

	got, err := Get("test/", 1, true)
	if err != nil {
		t.Errorf("Get() error %v", err)
		return
	}

	if len(got) == 0 {
		t.Errorf("Get() got %v", got)
		return
	}

	for i, tt := range tests {

		if got[tt.fileNo].Name != tt.wantName {
			t.Errorf("Get() Name test %d: got %s, want %s", i, got[tt.fileNo].Name, tt.wantName)
			return
		}

		if tt.wantByteSize > 0 {
			if got[tt.fileNo].ByteSize != tt.wantByteSize {
				t.Errorf("Get() ByteSize test %d: got %d, want %d", i, got[tt.fileNo].ByteSize, tt.wantByteSize)
				return
			}
		}

		if got[tt.fileNo].IsDir != tt.wantIsDir {
			t.Errorf("Get() IsDir test %b: got %t, want %t", i, got[tt.fileNo].IsDir, tt.wantIsDir)
			return
		}

		if got[tt.fileNo].Path != tt.wantPath {
			t.Errorf("Get() Path test %d: got %s, want %s", i, got[tt.fileNo].Path, tt.wantPath)
			return
		}
	}

	// test without dirs
	got, err = Get("test", 1, false)
	if err != nil {
		t.Errorf("Get() error %v", err)
		return
	}

	for _, fi := range got {
		if fi.IsDir {
			t.Errorf("Get() IsDir got true, want false")
			return
		}
	}

	// remove files and folders
	err = file.Remove("test/")
	if err != nil {
		t.Errorf("Get() error %v", err)
		return
	}

	// test error
	_, err = Get("test/", 1, true)
	if err == nil {
		t.Errorf("Get() got nil, want error")
		return
	}
}

func writeFiles(path string, nfiles int) error {

	var err error

	for i := 0; i < nfiles; i++ {

		var content string

		for n := 0; n < i; n++ {
			content += "0"
		}

		err = file.WriteFile(path+"file"+strconv.Itoa(i), []byte(content))
		if err != nil {
			return err
		}
	}

	return nil
}
