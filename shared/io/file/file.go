package file

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dekoch/gouniversal/shared/console"
)

func ReadFile(path string) ([]byte, error) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		b := make([]byte, 0)
		return b, err
	}

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		console.Log(err, "ReadFile()")
	} else {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			console.Log(err, "ReadFile()")
		} else {
			return content, err
		}
	}

	b := make([]byte, 0)
	return b, err
}

func WriteFile(path string, content []byte) error {

	// directory from path
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// if not found, create dir
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			console.Log(err, "WriteFile()")
			return err
		}
	}

	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		console.Log(err, "WriteFile()")
		return err
	}

	_, err = file.Write(content)
	if err != nil {
		console.Log(err, "WriteFile()")
	}

	return err
}
