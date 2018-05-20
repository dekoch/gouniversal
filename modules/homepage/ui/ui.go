package ui

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/dekoch/gouniversal/modules/homepage/global"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func registerMenuItems(menu string, filepath string, navpath string, nav *navigation.Navigation) {

	menuItems, _ := functions.ReadDir(filepath, 0)

	for _, i := range menuItems {

		if i.IsDir() {

			registerMenuItems(menu, filepath+i.Name()+"/", navpath+":"+i.Name(), nav)
		} else {

			if strings.HasSuffix(i.Name(), ".html") {
				title := strings.Replace(i.Name(), ".html", "", -1)
				nav.Sitemap.Register(menu, navpath+":"+i.Name(), title)
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
				title := strings.Replace(f.Name(), ".html", "", -1)
				nav.Sitemap.Register(title, "App:Homepage:"+f.Name(), title)
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
