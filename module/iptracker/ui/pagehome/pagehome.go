package pagehome

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/dekoch/gouniversal/module/iptracker/global"
	"github.com/dekoch/gouniversal/module/iptracker/lang"
	"github.com/dekoch/gouniversal/module/iptracker/typeiptracker"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func LoadConfig() {

}

func RegisterPage(page *typeiptracker.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.Home.Menu, "App:IPTracker:Home", page.Lang.Home.Title)
}

func Render(page *typeiptracker.Page, nav *navigation.Navigation, r *http.Request) {

	button := r.FormValue("edit")

	type content struct {
		Lang           lang.Home
		UpdateInterval template.HTML
	}
	var c content

	c.Lang = page.Lang.Home

	if button == "apply" {

		err := edit(r)
		if err != nil {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
		}
	}

	c.UpdateInterval = template.HTML(strconv.Itoa(global.Config.GetUpdIntervalInt()))

	cont, err := functions.PageToString(global.Config.UIFileRoot+"home.html", c)
	if err == nil {
		page.Content += cont
	} else {
		nav.RedirectPath("404", true)
	}
}

func edit(r *http.Request) error {

	var (
		err             error
		sUpdateInterval string
		iUpdateInterval int
	)

	func() {

		for i := 0; i <= 4; i++ {

			switch i {
			case 0:
				sUpdateInterval, err = functions.CheckFormInput("UpdateInterval", r)

			case 1:
				// check input
				if functions.IsEmpty(sUpdateInterval) {

					err = errors.New("bad input")
				}

			case 2:
				iUpdateInterval, err = strconv.Atoi(sUpdateInterval)

			case 3:
				// check converted input
				if iUpdateInterval < 0 {

					err = errors.New("bad input")
				}

			case 4:
				global.Config.SetUpdInterval(iUpdateInterval)

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
