package pagesettings

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/dekoch/gouniversal/modules/meshfilesync/global"
	"github.com/dekoch/gouniversal/modules/meshfilesync/lang"
	"github.com/dekoch/gouniversal/modules/meshfilesync/typesmfs"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typesmfs.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.Settings.Menu, "App:MeshFS:Settings", page.Lang.Settings.Title)
}

func Render(page *typesmfs.Page, nav *navigation.Navigation, r *http.Request) {

	button := r.FormValue("edit")

	type content struct {
		Lang              lang.Settings
		List              template.HTML
		MaxFileSize       template.HTML
		AutoAddChecked    template.HTML
		AutoUpdateChecked template.HTML
		AutoDeleteChecked template.HTML
	}
	var c content

	c.Lang = page.Lang.Settings

	if button == "apply" {

		err := edit(r)
		if err != nil {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
		}
	}

	c.MaxFileSize = template.HTML(strconv.FormatFloat(global.Config.GetMaxFileSize(), 'f', 0, 64))

	if global.Config.GetAutoAdd() {
		c.AutoAddChecked = template.HTML("checked")
	}

	if global.Config.GetAutoUpdate() {
		c.AutoUpdateChecked = template.HTML("checked")
	}

	if global.Config.GetAutoDelete() {
		c.AutoDeleteChecked = template.HTML("checked")
	}

	cont, err := functions.PageToString(global.Config.UIFileRoot+"settings.html", c)
	if err == nil {
		page.Content += cont
	} else {
		nav.RedirectPath("404", true)
	}
}

func edit(r *http.Request) error {

	var (
		err          error
		sMaxFileSize string
		fMaxFileSize float64
	)

	func() {

		for i := 0; i <= 4; i++ {

			switch i {
			case 0:
				sMaxFileSize, err = functions.CheckFormInput("MaxFileSize", r)

			case 1:
				// check input
				if functions.IsEmpty(sMaxFileSize) {

					err = errors.New("bad input")
				}

			case 2:
				fMaxFileSize, err = strconv.ParseFloat(sMaxFileSize, 64)

			case 3:
				// check converted input
				if fMaxFileSize < 0.1 {

					err = errors.New("bad input")
				}

			case 4:
				global.Config.SetMaxFileSize(fMaxFileSize)

				if r.FormValue("chbAutoAdd") == "checked" {
					global.Config.SetAutoAdd(true)
				} else {
					global.Config.SetAutoAdd(false)
				}

				if r.FormValue("chbAutoUpdate") == "checked" {
					global.Config.SetAutoUpdate(true)
				} else {
					global.Config.SetAutoUpdate(false)
				}

				if r.FormValue("chbAutoDelete") == "checked" {
					global.Config.SetAutoDelete(true)
				} else {
					global.Config.SetAutoDelete(false)
				}

				err = global.Config.SaveConfig()
			}

			if err != nil {
				console.Log(err, "")
				return
			}
		}
	}()

	return err
}
