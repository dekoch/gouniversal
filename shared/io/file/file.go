package file

import (
	"crypto/sha256"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dekoch/gouniversal/shared/functions"
)

// ReadFile reads the file and returns the contents
func ReadFile(path string) ([]byte, error) {

	var ret []byte

	if functions.IsEmpty(path) {
		return ret, errors.New("invalid path")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ret, err
	}

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return ret, err
	}

	ret, err = ioutil.ReadFile(path)

	return ret, err
}

// WriteFile writes the contents to file
func WriteFile(path string, content []byte) error {

	if functions.IsEmpty(path) {
		return errors.New("invalid path")
	}

	// directory from path
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// if not found, create dir
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return err
	}

	_, err = file.Write(content)

	return err
}

// Checksum returns the checksum of a file
func Checksum(path string) ([]byte, error) {

	var ret []byte

	if functions.IsEmpty(path) {
		return ret, errors.New("invalid path")
	}

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return ret, err
	}

	h := sha256.New()

	_, err = io.Copy(h, file)
	if err != nil {
		return ret, err
	}

	return h.Sum(nil), nil
}

// Remove removes path and any children it contains
func Remove(path string) error {

	if functions.IsEmpty(path) {
		return errors.New("invalid path")
	}

	// directory from path
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}

	return os.RemoveAll(path)
}
