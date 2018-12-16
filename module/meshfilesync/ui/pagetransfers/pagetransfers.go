package pagetransfers

import (
	"html/template"
	"net/http"

	"github.com/dekoch/gouniversal/module/meshfilesync/global"
	"github.com/dekoch/gouniversal/module/meshfilesync/lang"
	"github.com/dekoch/gouniversal/module/meshfilesync/typesmfs"
	"github.com/dekoch/gouniversal/shared/datasize"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typesmfs.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.Transfers.Menu, "App:MeshFS:Transfers", page.Lang.Transfers.Title)
}

func Render(page *typesmfs.Page, nav *navigation.Navigation, r *http.Request) {

	button := r.FormValue("edit")

	type content struct {
		Lang         lang.Transfers
		DownloadList template.HTML
		UploadList   template.HTML
	}
	var c content

	c.Lang = page.Lang.Transfers

	if button == "abort" {

		abort(r)
	}

	// download list
	files := global.DownloadFiles.Get()
	htmlFiles := ""

	for _, file := range files {

		size := datasize.ByteSize(file.Size).HumanReadable()

		htmlFiles += "<tr>"
		htmlFiles += "<td><input class=\"form-check-input\" type=\"checkbox\" name=\"download" + file.ID + "\" value=\"checked\"></td>"
		htmlFiles += "<td>" + file.Path + "</td>"
		htmlFiles += "<td>" + size + "</td>"
		htmlFiles += "<td>" + file.ModTime.Format("2006-01-02 15:04:05") + "</td>"
		htmlFiles += "</tr>"
	}

	c.DownloadList = template.HTML(htmlFiles)

	// upload list
	files = global.UploadFiles.Get()
	htmlFiles = ""

	for _, file := range files {

		size := datasize.ByteSize(file.Size).HumanReadable()

		htmlFiles += "<tr>"
		htmlFiles += "<td>" + file.Path + "</td>"
		htmlFiles += "<td>" + size + "</td>"
		htmlFiles += "<td>" + file.ModTime.Format("2006-01-02 15:04:05") + "</td>"
		htmlFiles += "</tr>"
	}

	c.UploadList = template.HTML(htmlFiles)

	cont, err := functions.PageToString(global.Config.UIFileRoot+"transfers.html", c)
	if err == nil {
		page.Content += cont
	} else {
		nav.RedirectPath("404", true)
	}
}

func abort(r *http.Request) error {

	if r.FormValue("ckbDownloadAll") == "checked" {

		global.DownloadFiles.Reset()
	} else {

		for _, file := range global.DownloadFiles.Get() {

			if r.FormValue("download"+file.ID) == "checked" {
				global.DownloadFiles.Delete(file.Path)
			}
		}
	}

	return nil
}
