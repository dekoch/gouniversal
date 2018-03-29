package functions

import (
	"bytes"
	"fmt"
	"html/template"
)

// TemplToString converts a Template and struct to string
func TemplToString(templ *template.Template, data interface{}) string {

	var tpl bytes.Buffer
	if err := templ.Execute(&tpl, data); err != nil {
		fmt.Println(err)
	}

	return tpl.String()
}
