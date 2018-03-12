package functions

import (
	"bytes"
	"html/template"
)

func TemplToString(templ *template.Template, data interface{}) string {

	var tpl bytes.Buffer
	if err := templ.Execute(&tpl, data); err != nil {

	}

	return tpl.String()
}
