package functions

import (
	"bytes"
	"errors"
	"html/template"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/dekoch/gouniversal/shared/console"
)

// PageToString converts a page and struct to string
func PageToString(path string, data interface{}) (string, error) {

	// check file exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		console.Log(err, "PageToString()")
		return "", err
	}

	// read template
	templ, err := template.ParseFiles(path)
	if err != nil {
		console.Log(err, "PageToString()")
		return "", err
	}

	// template to buffer
	var tpl bytes.Buffer
	if err := templ.Execute(&tpl, data); err != nil {
		console.Log(err, "PageToString()")
		return "", err
	}

	// buffer to string
	return tpl.String(), nil
}

// IsEmpty returns true, if a string is empty or whitespace
func IsEmpty(s string) bool {

	str := strings.Replace(s, " ", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	str = strings.Replace(str, "\n", "", -1)

	return len(str) == 0
}

// CheckFormInput returns a value for the named component of the query
// and checks for allowed characters
func CheckFormInput(key string, r *http.Request) (string, error) {

	input := r.FormValue(key)
	inputCheck := strings.Replace(input, " ", "", -1)
	inputCheck = strings.Replace(inputCheck, "+", "", -1)
	inputCheck = strings.Replace(inputCheck, "-", "", -1)
	inputCheck = strings.Replace(inputCheck, "_", "", -1)
	inputCheck = strings.Replace(inputCheck, ".", "", -1)
	inputCheck = strings.Replace(inputCheck, ":", "", -1)
	inputCheck = strings.Replace(inputCheck, ",", "", -1)
	inputCheck = strings.Replace(inputCheck, ";", "", -1)
	inputCheck = strings.Replace(inputCheck, "/", "", -1)

	if govalidator.IsUTFLetterNumeric(inputCheck) {
		return input, nil
	}

	return "", errors.New("bad input")
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

// CreateDir creates a missing directory
func CreateDir(path string) error {

	if IsEmpty(path) {
		return errors.New("invalid path")
	}

	// dir from path
	dir := filepath.Dir(path)

	var err error

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		// if not found, create dir
		err = os.MkdirAll(dir, os.ModePerm)
	}

	return err
}
