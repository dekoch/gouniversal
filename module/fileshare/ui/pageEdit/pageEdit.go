package pageEdit

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dekoch/gouniversal/module/fileshare/global"
	"github.com/dekoch/gouniversal/module/fileshare/lang"
	"github.com/dekoch/gouniversal/module/fileshare/typesFileshare"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func LoadConfig() {

}

func RegisterPage(page *typesFileshare.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "App:Fileshare:Edit", page.Lang.Edit.Title)
}

func Render(page *typesFileshare.Page, nav *navigation.Navigation, r *http.Request) {

	type content struct {
		Lang lang.Edit
		Name template.HTML
	}
	var c content

	c.Lang = page.Lang.Edit

	fileRoot := global.Config.FileRoot + nav.User.UUID + "/"

	var (
		err      error
		fullPath string
		dir      string
	)

	selFolder := nav.Parameter("Folder")
	selFile := nav.Parameter("File")

	if selFolder != "" {

		fullPath = selFolder
		// get containing directory
		selDir := filepath.Base(fullPath) + "/"
		dir = strings.TrimSuffix(fullPath, selDir)

	} else if selFile != "" {

		fullPath = selFile
		dir = filepath.Dir(fullPath)
	}

	// check if exist
	if _, oserr := os.Stat(fileRoot + fullPath); os.IsNotExist(oserr) {
		err = oserr
	}

	if err == nil {

		if strings.HasSuffix(dir, "/") == false {
			dir += "/"
		}

		name := filepath.Base(fullPath)

		edit := r.FormValue("edit")

		if edit == "apply" {

			// rename
			newName := r.FormValue("name")

			if newName != name && functions.IsEmpty(newName) == false {

				err = os.Rename(fileRoot+fullPath, fileRoot+dir+newName)
			}

			// redirect to containing folder
			nav.RedirectPath("App:Fileshare:Home$Folder="+dir, false)
		}

		c.Name = template.HTML(name)

		cont, err := functions.PageToString(global.Config.UIFileRoot+"edit.html", c)
		if err == nil {
			page.Content += cont
		} else {
			nav.RedirectPath("404", true)
		}

	} else {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "edit.go", nav.User.UUID)
		nav.RedirectPath("400", true)
	}
}
