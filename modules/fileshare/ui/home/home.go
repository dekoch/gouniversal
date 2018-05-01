package home

import (
	"fmt"
	"gouniversal/modules/fileshare/global"
	"gouniversal/modules/fileshare/lang"
	"gouniversal/modules/fileshare/typesFileshare"
	"gouniversal/shared/functions"
	"gouniversal/shared/navigation"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
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

			fi.Size = strconv.Itoa(int(l.Size()))
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

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	path := r.FormValue("path")
	path = global.Config.File.FileRoot + path

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(w, "path does not exist")
		return
	}

	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("uploadfile")

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	defer file.Close()

	out, err := os.Create(path + header.Filename)
	if err != nil {
		fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
		return
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	//fmt.Fprintf(w, "File uploaded successfully: ")
	//fmt.Fprintf(w, header.Filename)

	fmt.Fprintf(w, "<html><head><meta http-equiv=\"refresh\" content=\"0; url=/app\" /></head></html>")
}

func LoadConfig() {

	//disabled ToDo: upload authentication
	//http.HandleFunc("/fileserver/upload", uploadHandler) // Handle the incoming file
}

func Render(page *typesFileshare.Page, nav *navigation.Navigation, r *http.Request) {

	type Content struct {
		Lang lang.Home
		Path template.HTML
		List template.HTML
	}
	var c Content

	c.Lang = page.Lang.Home

	selFolder := nav.Parameter("Folder")
	path := ""

	if functions.IsEmpty(selFolder) == false {
		path = selFolder

		if path == "/" {
			path = ""
		}
	}

	// scan directory
	folders, files := searchContent(global.Config.File.FileRoot + path)

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
		htmlFolders += "<td></td>"
		htmlFolders += "</tr>"
	}

	htmlFiles := ""

	for _, f := range files {

		htmlFiles += "<tr>"
		htmlFiles += "<td><i class=\"fa fa-file\" aria-hidden=\"true\">"
		htmlFiles += "<td><a href=\"/fileshare/req/?file=" + path + f.Name + "\" download=\"" + f.Name + "\">" + f.Name + "</a></td>"
		htmlFiles += "<td>" + f.Size + "</td>"
		htmlFiles += "<td></td>"
		htmlFiles += "</tr>"
	}

	c.List = template.HTML(htmlFolders + htmlFiles)
	c.Path = template.HTML(path)

	content, err := functions.PageToString(global.Config.File.UIFileRoot+"home.html", c)
	if err == nil {
		page.Content += content
	} else {
		nav.RedirectPath("404", true)
	}
}
