package pageSearch

import (
	"html/template"
	"net/http"

	"github.com/dekoch/gouniversal/modules/meshFileSync/global"
	"github.com/dekoch/gouniversal/modules/meshFileSync/lang"
	"github.com/dekoch/gouniversal/modules/meshFileSync/typesMFS"
	"github.com/dekoch/gouniversal/shared/datasize"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typesMFS.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.Search.Menu, "App:MeshFS:Search", page.Lang.Search.Title)
}

func Render(page *typesMFS.Page, nav *navigation.Navigation, r *http.Request) {

	button := r.FormValue("edit")

	type content struct {
		Lang         lang.Search
		OutdatedList template.HTML
		MeshList     template.HTML
	}
	var c content

	c.Lang = page.Lang.Search

	if button == "download" {

		download(r)

	} else if button == "clear" {

		global.OutdatedFiles.Reset()
		global.MeshFiles.Reset()
	}

	// remove pending files from lists
	for _, file := range global.DownloadFiles.Get() {

		global.OutdatedFiles.Delete(file.Path)
		global.MeshFiles.Delete(file.Path)
	}

	// remove existing files from list
	for _, file := range global.LocalFiles.Get() {

		global.MeshFiles.Delete(file.Path)
	}

	// outdated list
	files := global.OutdatedFiles.Get()
	htmlFiles := ""

	for _, file := range files {

		size := datasize.ByteSize(file.Size).HumanReadable()

		htmlFiles += "<tr>"
		htmlFiles += "<td><input class=\"form-check-input\" type=\"checkbox\" name=\"outdated" + file.ID + "\" value=\"checked\"></td>"
		htmlFiles += "<td>" + file.Path + "</td>"
		htmlFiles += "<td>" + size + "</td>"
		htmlFiles += "<td>" + file.ModTime.Format("2006-01-02 15:04:05") + "</td>"
		htmlFiles += "</tr>"
	}

	c.OutdatedList = template.HTML(htmlFiles)

	// mesh list
	files = global.MeshFiles.Get()
	htmlFiles = ""

	for _, file := range files {

		size := datasize.ByteSize(file.Size).HumanReadable()

		htmlFiles += "<tr>"
		htmlFiles += "<td><input class=\"form-check-input\" type=\"checkbox\" name=\"mesh" + file.ID + "\" value=\"checked\"></td>"
		htmlFiles += "<td>" + file.Path + "</td>"
		htmlFiles += "<td>" + size + "</td>"
		htmlFiles += "<td>" + file.ModTime.Format("2006-01-02 15:04:05") + "</td>"
		htmlFiles += "</tr>"
	}

	c.MeshList = template.HTML(htmlFiles)

	cont, err := functions.PageToString(global.Config.UIFileRoot+"search.html", c)
	if err == nil {
		page.Content += cont
	} else {
		nav.RedirectPath("404", true)
	}
}

func download(r *http.Request) error {

	// outdated
	files := global.OutdatedFiles.Get()

	if r.FormValue("ckbOutdatedAll") == "checked" {

		global.DownloadFiles.AddList(files)
	} else {

		for _, file := range files {

			if r.FormValue("outdated"+file.ID) == "checked" {
				global.DownloadFiles.Add(file)
			}
		}
	}

	// mesh
	files = global.MeshFiles.Get()

	if r.FormValue("ckbMeshAll") == "checked" {

		global.DownloadFiles.AddList(files)
	} else {

		for _, file := range files {

			if r.FormValue("mesh"+file.ID) == "checked" {
				global.DownloadFiles.Add(file)
			}
		}
	}

	return nil
}
