package config

import (
	"encoding/json"

	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
)

// FileHeader default config file header
type FileHeader struct {
	HeaderVersion  float32
	FileName       string
	ContentName    string
	ContentVersion float32
	Comment        string
}

type File struct {
	Header  FileHeader
	Content interface{}
}

// BuildHeader builds a default config file header
func BuildHeader(filename string, conname string, conver float32, comment string) FileHeader {
	var h FileHeader

	h.HeaderVersion = 1.0

	h.FileName = filename
	h.ContentName = conname
	h.ContentVersion = conver
	h.Comment = comment

	return h
}

// BuildHeader builds a default config file header with a struct
func BuildHeaderWithStruct(h FileHeader) FileHeader {

	h.HeaderVersion = 1.0

	return h
}

// CheckHeader returns true, if conname and Header ContentName is identical
func CheckHeader(fh FileHeader, conname string) bool {
	if fh.ContentName != conname {
		return false
	}

	return true
}

func Save(path string, h FileHeader, content interface{}) error {

	var f File

	f.Header = BuildHeaderWithStruct(h)
	f.Content = content

	b, err := json.Marshal(f)
	if err != nil {
		console.Log(err, "config.Save()")
	}

	err = file.WriteFile(path+h.FileName, b)

	return err
}

func Load(path string, h FileHeader) (interface{}, error) {

	b, err := file.ReadFile(path + h.FileName)
	if err != nil {
		console.Log(err, "config.Load()")
	}

	var f File

	err = json.Unmarshal(b, &f)
	if err != nil {
		console.Log(err, "config.Load()")
	}

	if CheckHeader(f.Header, h.ContentName) == false {
		console.Log("wrong config \""+path+h.FileName+"\"", "config.Load()")
	}

	return f, err
}
