package ui

import (
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"

	"github.com/dekoch/gouniversal/module/homepage/global"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/io/fileinfo"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	path := global.Config.UIFileRoot + "static/"
	// if exist handle static folder
	if _, err := os.Stat(path); os.IsNotExist(err) == false {
		fs := http.FileServer(http.Dir(path))
		http.Handle("/homepage/static/", http.StripPrefix("/homepage/static/", fs))
	}
}

func getNameAndOrder(s string) (string, int) {

	name := s
	no := ""
	order := -1

	if strings.Contains(s, ".") {
		s := strings.SplitAfterN(s, ".", -1)
		no = s[0]
		name = s[1]

		no = strings.Replace(no, ".", "", -1)
	}

	if no != "" && govalidator.IsNumeric(no) {
		newOrder, err := strconv.Atoi(no)
		if err != nil {
			console.Log(err, "homepage/ui.go")
		} else {
			order = newOrder
		}
	}

	return name, order
}

func registerMenuItems(menu string, filepath string, navpath string, nav *navigation.Navigation) error {

	menuItems, err := fileinfo.Get(filepath, 0, true)
	if err != nil {
		return err
	}

	for _, i := range menuItems {

		if i.IsDir {

			err = registerMenuItems(menu, filepath+i.Name+"/", navpath+":"+i.Name, nav)
			if err != nil {
				return err
			}
		} else {

			if strings.HasSuffix(i.Name, ".html") {

				menuName, menuOrder := getNameAndOrder(menu)

				file := strings.Replace(i.Name, ".html", "", -1)
				fileName, fileOrder := getNameAndOrder(file)

				nav.Sitemap.RegisterWithOrder(menuName, menuOrder, navpath+":"+i.Name, fileName, fileOrder)
			}
		}
	}

	return nil
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	// autogenerate menu entries from folders and files in module UIFileRoot
	menuFolders, err := fileinfo.Get(global.Config.UIFileRoot, 0, true)
	if err != nil {
		return
	}

	for _, f := range menuFolders {

		if f.IsDir {

			if f.Name == "static" {
				continue
			}

			err = registerMenuItems(f.Name, global.Config.UIFileRoot+f.Name+"/", "App:Homepage:"+f.Name, nav)
			if err != nil {
				return
			}
		} else {

			if strings.HasSuffix(f.Name, ".html") {

				file := strings.Replace(f.Name, ".html", "", -1)
				fileName, fileOrder := getNameAndOrder(file)

				nav.Sitemap.RegisterWithOrder(fileName, fileOrder, "App:Homepage:"+f.Name, fileName, -1)
			}
		}
	}
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	type content struct {
		Index template.HTML
	}
	var c content

	// use navigation path to get path to html
	path := strings.Replace(nav.Path, "App:Homepage:", "", -1)
	path = strings.Replace(path, ":", "/", -1)
	path = global.Config.UIFileRoot + path

	p, err := functions.PageToString(path, c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
