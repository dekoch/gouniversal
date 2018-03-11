package file

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type File struct {
	String string
	Bytes  []byte
}

func (f File) ReadFile(path string) ([]byte, error) {

	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	} else {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		return content, nil
	}
	defer file.Close()

	b := make([]byte, 0)

	return b, err
}

func (f File) WriteFile(path string, content []byte) error {

	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// if not found, create dir
		os.MkdirAll(dir, os.ModePerm)
	}

	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Cannot create file " + path)
	}
	defer file.Close()

	file.Write(content)

	file.Close()

	return err
}
