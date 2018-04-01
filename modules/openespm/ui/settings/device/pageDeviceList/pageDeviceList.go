package pageDeviceList

import (
	"fmt"
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

	nav.Sitemap.Register("App:Program:openESPM:Settings:Device:List", page.Lang.Settings.Device.List.Title)
}

func Render(page *typesOESPM.Page, nav *navigation.Navigation, r *http.Request) {

	type devicelist struct {
		Lang       langOESPM.SettingsDeviceList
		DeviceList template.HTML
	}
	var dl devicelist

	dl.Lang = page.Lang.Settings.Device.List

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
		tbody += "<td><button class=\"btn btn-default fa fa-wrench\" type=\"submit\" name=\"navigation\" value=\"App:Program:openESPM:Settings:Device:Edit$UUID=" + dev.UUID + "\" title=\"" + dl.Lang.Edit + "\"></button></td>"
		tbody += "</tr>"
	}
	globalOESPM.DeviceConfig.Mut.Unlock()

	dl.DeviceList = template.HTML(tbody)

	templ, err := template.ParseFiles(globalOESPM.UiConfig.AppFileRoot + "settings/devicelist.html")
	if err != nil {
		fmt.Println(err)
	}

	page.Content += functions.TemplToString(templ, dl)
}
