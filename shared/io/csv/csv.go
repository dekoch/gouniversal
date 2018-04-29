package csv

import (
	"encoding/csv"
	"gouniversal/shared/functions"
	"os"
)

func AddRow(fname string, row []string) error {

	err := functions.CreateDir(fname)
	if err != nil {
		return err
	}

	var rows [][]string

	if _, err = os.Stat(fname); os.IsNotExist(err) == false {
		// read the file
		f, err := os.Open(fname)
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
	f, err := os.Create(fname)
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
