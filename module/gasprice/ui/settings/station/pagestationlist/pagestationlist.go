package pagestationlist

import (
	"html/template"
	"net/http"

	"github.com/dekoch/gouniversal/module/gasprice/global"
	"github.com/dekoch/gouniversal/module/gasprice/lang"
	"github.com/dekoch/gouniversal/module/gasprice/typemd"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.StationList.Menu, "App:GasPrice:Settings:Station:List", page.Lang.StationList.Title)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	type Content struct {
		Lang        lang.StationList
		StationList template.HTML
	}
	var c Content

	c.Lang = page.Lang.StationList

	tbody := ""

	for _, station := range global.Config.Stations.GetList() {

		tbody += "<tr>"
		tbody += "<td></td>"
		tbody += "<td>" + station.Name + "</td>"
		tbody += "<td>" + station.Company + "</td>"
		tbody += "<td>" + station.Street + "</td>"
		tbody += "<td>" + station.City + "</td>"
		tbody += "<td><button class=\"btn btn-default fa fa-wrench\" type=\"submit\" name=\"navigation\" value=\"App:GasPrice:Settings:Station:Edit$UUID=" + station.UUID + "\" title=\"" + c.Lang.Edit + "\"></button></td>"
		tbody += "</tr>"
	}

	c.StationList = template.HTML(tbody)

	p, err := functions.PageToString(global.Config.UIFileRoot+"settings/stationlist.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
