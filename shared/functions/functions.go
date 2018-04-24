package functions

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
)

// PageToString converts a page and struct to string
func PageToString(path string, data interface{}) (string, error) {

	// check file exist
	if _, err := os.Stat(path); os.IsNotExist(err) {

		fmt.Println(err)
		return "", err
	}

	// read template
	templ, err := template.ParseFiles(path)
	if err != nil {

		fmt.Println(err)
		return "", err
	}

	// template to buffer
	var tpl bytes.Buffer
	if err := templ.Execute(&tpl, data); err != nil {

		fmt.Println(err)
		return "", err
	}

	// buffer to string
	return tpl.String(), nil
}

// TemplToString converts a Template and struct to string
func TemplToString(templ *template.Template, data interface{}) string {

	var tpl bytes.Buffer
	if err := templ.Execute(&tpl, data); err != nil {
		fmt.Println(err)
	}

	return tpl.String()
}

// IsEmpty returns true, if a string is empty or whitespace
func IsEmpty(s string) bool {

	str := strings.Replace(s, " ", "", -1)

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

	if govalidator.IsUTFLetterNumeric(inputCheck) {
		return input, nil
	}

	return "", errors.New("bad input")
}

// readDir is a helper for ReadDir()
func readDir(dir string, maxdepth int, currdepth int) ([]os.FileInfo, error) {

	files, err := ioutil.ReadDir(dir)

	for _, fl := range files {

		if err == nil {

			if fl.IsDir() {

				if currdepth < maxdepth {

					var sub []os.FileInfo
					sub, err = readDir(dir+fl.Name()+"/", maxdepth, currdepth+1)

					if err == nil {
						files = append(sub, files...)
					}
				}
			}
		}
	}

	return files, err
}

// ReadDir is the same as ioutil.ReadDir() but recursive
// with a max depth option
func ReadDir(dir string, maxdepth int) ([]os.FileInfo, error) {

	return readDir(dir, maxdepth, 0)
}
