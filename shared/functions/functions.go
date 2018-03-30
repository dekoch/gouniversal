package functions

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
)

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
