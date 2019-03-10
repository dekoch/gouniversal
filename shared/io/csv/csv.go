package csv

import (
	"encoding/csv"
	"errors"
	"os"

	"github.com/dekoch/gouniversal/shared/functions"
)

// AddRow adds a new row of data to a csv file
func AddRow(path string, row []string) error {

	err := functions.CreateDir(path)
	if err != nil {
		return err
	}

	var rows [][]string

	if _, err = os.Stat(path); os.IsNotExist(err) == false {
		// read the file
		file, err := os.Open(path)
		if err != nil {
			return err
		}

		r := csv.NewReader(file)
		rows, err = r.ReadAll()
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}
	}

	rows = append(rows, row)

	// write the file
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	w := csv.NewWriter(file)
	err = w.WriteAll(rows)
	if err != nil {
		file.Close()
		return err
	}

	return file.Close()
}

// ReadAll returns the content of a csv file
func ReadAll(path string) ([][]string, error) {

	var ret [][]string

	if functions.IsEmpty(path) {
		return ret, errors.New("invalid path")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ret, err
	}

	// read the file
	file, err := os.Open(path)
	if err != nil {
		return ret, err
	}

	r := csv.NewReader(file)
	ret, err = r.ReadAll()
	if err != nil {
		return ret, err
	}

	return ret, file.Close()
}
