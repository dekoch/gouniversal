package csv

import (
	"encoding/csv"
	"os"

	"github.com/dekoch/gouniversal/shared/functions"
)

func AddRow(filepath string, row []string) error {

	err := functions.CreateDir(filepath)
	if err != nil {
		return err
	}

	var rows [][]string

	if _, err = os.Stat(filepath); os.IsNotExist(err) == false {
		// read the file
		f, err := os.Open(filepath)
		if err != nil {
			return err
		}
		r := csv.NewReader(f)
		rows, err = r.ReadAll()
		if err != nil {
			return err
		}
		if err = f.Close(); err != nil {
			return err
		}
	}

	rows = append(rows, row)

	// write the file
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	if err = w.WriteAll(rows); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func ReadAll(filepath string) ([][]string, error) {

	var ret [][]string

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return ret, err
	}

	// read the file
	f, err := os.Open(filepath)
	if err != nil {
		return ret, err
	}

	r := csv.NewReader(f)
	ret, err = r.ReadAll()
	if err != nil {
		return ret, err
	}

	return ret, f.Close()
}
