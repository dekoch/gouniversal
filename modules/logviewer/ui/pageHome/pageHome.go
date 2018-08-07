package pageHome

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/dekoch/gouniversal/modules/logviewer/global"
	"github.com/dekoch/gouniversal/modules/logviewer/lang"
	"github.com/dekoch/gouniversal/modules/logviewer/typesLogViewer"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/io/file"
	"github.com/dekoch/gouniversal/shared/io/fileInfo"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typesLogViewer.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program", "App:LogViewer:Home", page.Lang.Home.Title)
}

func parentDir(path string) string {

	newPath := path

	// remove last /
	if strings.HasSuffix(newPath, "/") {
		index := strings.LastIndex(newPath, "/")
		newPath = newPath[:index]
	}

	// remove last directory from path
	index := strings.LastIndex(newPath, "/")
	if index < 0 {
		return ""
	}

	cnt := len(newPath)
	if cnt > 0 {
		newPath = newPath[:index]
	}

	return newPath + "/"
}

func Render(page *typesLogViewer.Page, nav *navigation.Navigation, r *http.Request) {

	type Content struct {
		Lang        lang.Home
		List        template.HTML
		FileHidden  template.HTML
		FileName    template.HTML
		FileContent template.HTML
	}
	var c Content

	c.Lang = page.Lang.Home

	fileRoot := global.Config.File.LogFileRoot

	selFolder := nav.Parameter("Folder")
	path := ""

	if functions.IsEmpty(selFolder) == false {
		path = selFolder

		if path == "/" {
			path = ""
		}
	}

	// scan directory
	folders, files := fileInfo.Get(fileRoot + path)

	htmlFolders := ""

	if path != "" {
		// fileshare root
		htmlFolders += "<tr>"
		htmlFolders += "<td></td>"
		htmlFolders += "<td><button class=\"btn btn-link\" type=\"submit\" name=\"navigation\" value=\"App:LogViewer:Home$Folder=/\">..</button></td>"
		htmlFolders += "<td></td>"
		htmlFolders += "</tr>"

		// parent directory
		htmlFolders += "<tr>"
		htmlFolders += "<td></td>"
		htmlFolders += "<td><button class=\"btn btn-link\" type=\"submit\" name=\"navigation\" value=\"App:LogViewer:Home$Folder=" + parentDir(path) + "\">.</button></td>"
		htmlFolders += "<td></td>"
		htmlFolders += "</tr>"
	}

	for _, f := range folders {

		htmlFolders += "<tr>"
		htmlFolders += "<td><i class=\"fa fa-folder\" aria-hidden=\"true\"></td>"
		htmlFolders += "<td><button class=\"btn btn-link\" type=\"submit\" name=\"navigation\" value=\"App:LogViewer:Home$Folder=" + path + f.Name + "/" + "\">" + f.Name + "</button></td>"
		htmlFolders += "<td></td>"
		htmlFolders += "</tr>"
	}

	htmlFiles := ""

	for _, f := range files {

		htmlFiles += "<tr>"
		htmlFiles += "<td><i class=\"fa fa-file\" aria-hidden=\"true\"></td>"
		htmlFiles += "<td><button class=\"btn btn-link\" type=\"submit\" name=\"view\" value=\"" + f.Name + "\">" + f.Name + "</button></td>"
		htmlFiles += "<td>" + f.Size + "</td>"
		htmlFiles += "</tr>"
	}

	c.List = template.HTML(htmlFolders + htmlFiles)

	c.FileHidden = template.HTML("hidden")

	fileName := ""

	view := r.FormValue("view")
	if functions.IsEmpty(view) == false {
		fileName = view
	} else {
		// if no file is selected, load last file from list
		fileCnt := len(files)
		if fileCnt > 0 {
			fileName = files[fileCnt-1].Name
		}
	}

	if fileName != "" {
		// open and convert log file to HTML
		c.FileHidden = template.HTML("")
		c.FileName = template.HTML(fileName)

		fileRaw, err := file.ReadFile(fileRoot + path + fileName)
		if err == nil {

			fileString := string(fileRaw[:])
			fileString = strings.Replace(fileString, "\r\n", "<br>", -1)
			fileString = strings.Replace(fileString, "\n", "<br>", -1)
			fileString = strings.Replace(fileString, "\r", "<br>", -1)
			c.FileContent = template.HTML(fileString)
		}
	}

	p, err := functions.PageToString(global.Config.File.UIFileRoot+"home.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}