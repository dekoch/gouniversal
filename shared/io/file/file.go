package file

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dekoch/gouniversal/shared/console"
)

func ReadFile(path string) ([]byte, error) {

	file, err := os.Open(path)
	if err != nil {
		console.Log(err, "ReadFile()")
	} else {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			console.Log(err, "ReadFile()")
		}

		return content, nil
	}
	defer file.Close()

	b := make([]byte, 0)

	return b, err
}

func WriteFile(path string, content []byte) error {

	// directory from path
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// if not found, create dir
		os.MkdirAll(dir, os.ModePerm)
	}

	file, err := os.Create(path)
	if err != nil {
		console.Log(err, "WriteFile()")
	}
	defer file.Close()

	file.Write(content)
	file.Close()

	return err
}
