package devicemanagement

import (
	"html/template"

	"github.com/dekoch/gouniversal/module/openespm/global"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
)

func HTMLSelectDevice(name string, appname string, uid string) template.HTML {

	type content struct {
		Title  template.HTML
		Select template.HTML
	}
	var c content

	title := "..."

	sel := "<select name=\"" + name + "\">"

	if uid == "" {
		sel += "<option value=\"\"></option>"
	}

	devices := global.DeviceConfig.List()

	for u := 0; u < len(devices); u++ {

		// list only devices with the same app
		if appname == devices[u].App {
			sel += "<option value=\"" + devices[u].UUID + "\""

			if uid == devices[u].UUID {
				sel += " selected"

				title = devices[u].Name
			}

			sel += ">" + devices[u].Name + "</option>"
		}
	}

	sel += "</select>"

	c.Title = template.HTML(title)
	c.Select = template.HTML(sel)

	p, err := functions.PageToString(global.UiConfig.AppFileRoot+"selectdevice.html", c)
	if err != nil {
		console.Log(err, "")
		p = err.Error()
	}

	return template.HTML(p)
}
