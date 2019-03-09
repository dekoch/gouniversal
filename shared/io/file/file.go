package file

import (
	"crypto/sha256"
	"io"
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

func Checksum(path string) ([]byte, error) {

	f, err := os.Open(path)
	if err != nil {
		var e []byte
		return e, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		var e []byte
		return e, err
	}

	return h.Sum(nil), nil
}

func Remove(path string) error {

	// directory from path
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}

	return os.RemoveAll(path)
}
