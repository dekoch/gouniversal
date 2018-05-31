package home

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/dekoch/gouniversal/modules/fileshare/global"
	"github.com/dekoch/gouniversal/modules/fileshare/lang"
	"github.com/dekoch/gouniversal/modules/fileshare/typesFileshare"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/datasize"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"

	"github.com/google/uuid"
)

type fileInfo struct {
	Name string
	Size string
}

func RegisterPage(page *typesFileshare.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Fileshare", "App:Fileshare:Home", page.Lang.Home.Title)
}

func searchContent(path string) ([]fileInfo, []fileInfo) {

	list, _ := functions.ReadDir(path, 0)

	folders := []fileInfo{}
	files := []fileInfo{}

	var fi fileInfo

	for _, l := range list {

		fi.Name = l.Name()

		if l.IsDir() {

			fi.Size = ""
			folders = append(folders, fi)
		} else {

			s := datasize.ByteSize(l.Size()).HumanReadable()
			fi.Size = s
			files = append(files, fi)
		}
	}

	return folders, files
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

func LoadConfig() {

}

func Render(page *typesFileshare.Page, nav *navigation.Navigation, r *http.Request) {

	type Content struct {
		Lang  lang.Home
		UUID  template.HTML
		Token template.HTML
		Path  template.HTML
		List  template.HTML
	}
	var c Content

	c.Lang = page.Lang.Home

	fileRoot := global.Config.File.FileRoot + nav.User.UUID + "/"

	selFolder := nav.Parameter("Folder")
	path := ""

	if functions.IsEmpty(selFolder) == false {
		path = selFolder

		if path == "/" {
			path = ""
		}
	}

	edit := r.FormValue("edit")

	if edit == "newfolder" {

		u := uuid.Must(uuid.NewRandom())

		err := functions.CreateDir(fileRoot + path + u.String() + "/")
		if err != nil {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "home.go", nav.User.UUID)
		}

	} else if strings.HasPrefix(edit, "deletefolder") {

		folder := strings.Replace(edit, "deletefolder", "", 1)
		err := os.RemoveAll(fileRoot + path + folder)
		if err != nil {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "home.go", nav.User.UUID)
		}

	} else if strings.HasPrefix(edit, "deletefile") {

		file := strings.Replace(edit, "deletefile", "", 1)
		err := os.Remove(fileRoot + path + file)
		if err != nil {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "home.go", nav.User.UUID)
		}
	}

	// scan directory
	folders, files := searchContent(fileRoot + path)

	htmlFolders := ""

	if path != "" {
		// fileshare root
		htmlFolders += "<tr>"
		htmlFolders += "<td>"
		htmlFolders += "<td><button class=\"btn btn-link\" type=\"submit\" name=\"navigation\" value=\"App:Fileshare:Home$Folder=/\">..</button></td>"
		htmlFolders += "<td></td>"
		htmlFolders += "<td></td>"
		htmlFolders += "</tr>"

		// parent directory
		htmlFolders += "<tr>"
		htmlFolders += "<td>"
		htmlFolders += "<td><button class=\"btn btn-link\" type=\"submit\" name=\"navigation\" value=\"App:Fileshare:Home$Folder=" + parentDir(path) + "\">.</button></td>"
		htmlFolders += "<td></td>"
		htmlFolders += "<td></td>"
		htmlFolders += "</tr>"
	}

	for _, f := range folders {

		htmlFolders += "<tr>"
		htmlFolders += "<td><i class=\"fa fa-folder\" aria-hidden=\"true\">"
		htmlFolders += "<td><button class=\"btn btn-link\" type=\"submit\" name=\"navigation\" value=\"App:Fileshare:Home$Folder=" + path + f.Name + "/" + "\">" + f.Name + "</button></td>"
		htmlFolders += "<td>" + f.Size + "</td>"
		htmlFolders += "<td><button class=\"btn btn-danger fa fa-trash\" type=\"submit\" name=\"edit\" value=\"deletefolder" + f.Name + "\" title=\"" + c.Lang.Delete + "\"></button></td>"
		htmlFolders += "</tr>"
	}

	htmlFiles := ""

	for _, f := range files {

		htmlFiles += "<tr>"
		htmlFiles += "<td><i class=\"fa fa-file\" aria-hidden=\"true\">"
		htmlFiles += "<td><a href=\"/fileshare/req/?file=" + nav.User.UUID + "/" + path + f.Name + "\" download=\"" + f.Name + "\">" + f.Name + "</a></td>"
		htmlFiles += "<td>" + f.Size + "</td>"
		htmlFiles += "<td><button class=\"btn btn-danger fa fa-trash\" type=\"submit\" name=\"edit\" value=\"deletefile" + f.Name + "\" title=\"" + c.Lang.Delete + "\"></button></td>"
		htmlFiles += "</tr>"
	}

	c.List = template.HTML(htmlFolders + htmlFiles)
	c.Path = template.HTML(nav.User.UUID + "/" + path)

	c.UUID = template.HTML(nav.User.UUID)
	c.Token = template.HTML(global.Tokens.New(nav.User.UUID))

	content, err := functions.PageToString(global.Config.File.UIFileRoot+"home.html", c)
	if err == nil {
		page.Content += content
	} else {
		nav.RedirectPath("404", true)
	}
}
