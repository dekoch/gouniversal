package config

import (
	"testing"
)

func TestBuildHeader(t *testing.T) {

	tests := []struct {
		inFileName       string
		inContentName    string
		inContentVersion float32
		inComment        string
	}{
		{"", "", 0.0, ""},
		{"file1234", "content1234", 1.234, "comment1234"},
	}

	for i, tt := range tests {

		got := BuildHeader(tt.inFileName, tt.inContentName, tt.inContentVersion, tt.inComment)

		if got.HeaderVersion != 1.0 {
			t.Errorf("BuildHeader() HeaderVersion test %d: got %f, want %f", i, got.HeaderVersion, 1.0)
		}

		if got.FileName != tt.inFileName {
			t.Errorf("BuildHeader() FileName test %d: got %s, want %s", i, got.FileName, tt.inFileName)
		}

		if got.ContentName != tt.inContentName {
			t.Errorf("BuildHeader() ContentName test %d: got %s, want %s", i, got.ContentName, tt.inContentName)
		}

		if got.ContentVersion != tt.inContentVersion {
			t.Errorf("BuildHeader() ContentVersion test %d: got %f, want %f", i, got.ContentVersion, tt.inContentVersion)
		}

		if got.Comment != tt.inComment {
			t.Errorf("BuildHeader() Comment test %d: got %s, want %s", i, got.Comment, tt.inComment)
		}
	}
}

func TestBuildHeaderWithStruct(t *testing.T) {

	tests := []struct {
		inFileName       string
		inContentName    string
		inContentVersion float32
		inComment        string
	}{
		{"", "", 0.0, ""},
		{"file1234", "content1234", 1.234, "comment1234"},
	}

	for i, tt := range tests {

		var h FileHeader
		h.FileName = tt.inFileName
		h.ContentName = tt.inContentName
		h.ContentVersion = tt.inContentVersion
		h.Comment = tt.inComment

		got := BuildHeaderWithStruct(h)

		if got.HeaderVersion != 1.0 {
			t.Errorf("BuildHeaderWithStruct() HeaderVersion test %d: got %f, want %f", i, got.HeaderVersion, 1.0)
		}

		if got.FileName != tt.inFileName {
			t.Errorf("BuildHeaderWithStruct() FileName test %d: got %s, want %s", i, got.FileName, tt.inFileName)
		}

		if got.ContentName != tt.inContentName {
			t.Errorf("BuildHeaderWithStruct() ContentName test %d: got %s, want %s", i, got.ContentName, tt.inContentName)
		}

		if got.ContentVersion != tt.inContentVersion {
			t.Errorf("BuildHeaderWithStruct() ContentVersion test %d: got %f, want %f", i, got.ContentVersion, tt.inContentVersion)
		}

		if got.Comment != tt.inComment {
			t.Errorf("BuildHeaderWithStruct() Comment test %d: got %s, want %s", i, got.Comment, tt.inComment)
		}
	}
}

func TestCheckHeader(t *testing.T) {

	tests := []struct {
		inHeader string
		inName   string
		want     bool
	}{
		{"", "", true},
		{"1234", "", false},
		{"", "1234", false},
		{"1234", "1234", true},
	}

	for i, tt := range tests {

		var h FileHeader
		h.ContentName = tt.inHeader

		got := CheckHeader(h, tt.inName)

		if got != tt.want {
			t.Errorf("CheckHeader() test %d: inHeader %s, inName %s, got %t, want %t", i, tt.inHeader, tt.inName, got, tt.want)
		}
	}
}
