package pageDeviceList

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/dekoch/gouniversal/module/openespm/global"
	"github.com/dekoch/gouniversal/module/openespm/lang"
	"github.com/dekoch/gouniversal/module/openespm/typesOESPM"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typesOESPM.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "App:openESPM:Settings:Device:List", page.Lang.Settings.Device.List.Title)
}

func Render(page *typesOESPM.Page, nav *navigation.Navigation, r *http.Request) {

	type content struct {
		Lang       lang.SettingsDeviceList
		DeviceList template.HTML
	}
	var c content

	c.Lang = page.Lang.Settings.Device.List

	var tbody string
	tbody = ""
	var intIndex int
	intIndex = 0

	devices := global.DeviceConfig.List()

	for u := 0; u < len(devices); u++ {

		dev := devices[u]

		intIndex += 1

		tbody += "<tr>"
		tbody += "<th scope='row'>" + strconv.Itoa(intIndex) + "</th>"
		tbody += "<td>" + dev.Name + "</td>"
		tbody += "<td>" + dev.App + "</td>"
		tbody += "<td><button class=\"btn btn-default fa fa-wrench\" type=\"submit\" name=\"navigation\" value=\"App:openESPM:Settings:Device:Edit$UUID=" + dev.UUID + "\" title=\"" + c.Lang.Edit + "\"></button></td>"
		tbody += "</tr>"
	}

	c.DeviceList = template.HTML(tbody)

	p, err := functions.PageToString(global.UiConfig.AppFileRoot+"settings/devicelist.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
