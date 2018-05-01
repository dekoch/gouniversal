package pageDeviceList

import (
	"gouniversal/modules/openespm/globalOESPM"
	"gouniversal/modules/openespm/langOESPM"
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/shared/functions"
	"gouniversal/shared/navigation"
	"html/template"
	"net/http"
	"strconv"
)

func RegisterPage(page *typesOESPM.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "App:openESPM:Settings:Device:List", page.Lang.Settings.Device.List.Title)
}

func Render(page *typesOESPM.Page, nav *navigation.Navigation, r *http.Request) {

	type content struct {
		Lang       langOESPM.SettingsDeviceList
		DeviceList template.HTML
	}
	var c content

	c.Lang = page.Lang.Settings.Device.List

	var tbody string
	tbody = ""
	var intIndex int
	intIndex = 0

	globalOESPM.DeviceConfig.Mut.Lock()
	for i := 0; i < len(globalOESPM.DeviceConfig.File.Devices); i++ {

		dev := globalOESPM.DeviceConfig.File.Devices[i]

		intIndex += 1

		tbody += "<tr>"
		tbody += "<th scope='row'>" + strconv.Itoa(intIndex) + "</th>"
		tbody += "<td>" + dev.Name + "</td>"
		tbody += "<td>" + dev.App + "</td>"
		tbody += "<td><button class=\"btn btn-default fa fa-wrench\" type=\"submit\" name=\"navigation\" value=\"App:openESPM:Settings:Device:Edit$UUID=" + dev.UUID + "\" title=\"" + c.Lang.Edit + "\"></button></td>"
		tbody += "</tr>"
	}
	globalOESPM.DeviceConfig.Mut.Unlock()

	c.DeviceList = template.HTML(tbody)

	p, err := functions.PageToString(globalOESPM.UiConfig.AppFileRoot+"settings/devicelist.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
