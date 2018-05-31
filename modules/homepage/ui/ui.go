package ui

import (
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"

	"github.com/dekoch/gouniversal/modules/homepage/global"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	path := global.Config.File.UIFileRoot + "static/"
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

func registerMenuItems(menu string, filepath string, navpath string, nav *navigation.Navigation) {

	menuItems, _ := functions.ReadDir(filepath, 0)

	for _, i := range menuItems {

		if i.IsDir() {

			registerMenuItems(menu, filepath+i.Name()+"/", navpath+":"+i.Name(), nav)
		} else {

			if strings.HasSuffix(i.Name(), ".html") {

				menuName, menuOrder := getNameAndOrder(menu)

				file := strings.Replace(i.Name(), ".html", "", -1)
				fileName, fileOrder := getNameAndOrder(file)

				nav.Sitemap.RegisterWithOrder(menuName, menuOrder, navpath+":"+i.Name(), fileName, fileOrder)
			}
		}
	}
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	// autogenerate menu entries from folders and files in module UIFileRoot
	menuFolders, _ := functions.ReadDir(global.Config.File.UIFileRoot, 0)

	for _, f := range menuFolders {

		if f.IsDir() {

			registerMenuItems(f.Name(), global.Config.File.UIFileRoot+f.Name()+"/", "App:Homepage:"+f.Name(), nav)
		} else {

			if strings.HasSuffix(f.Name(), ".html") {

				file := strings.Replace(f.Name(), ".html", "", -1)
				fileName, fileOrder := getNameAndOrder(file)

				nav.Sitemap.RegisterWithOrder(fileName, fileOrder, "App:Homepage:"+f.Name(), fileName, -1)
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
	path = global.Config.File.UIFileRoot + path

	p, err := functions.PageToString(path, c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
