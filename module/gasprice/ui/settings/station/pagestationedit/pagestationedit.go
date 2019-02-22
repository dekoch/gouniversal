package pagestationedit

import (
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/dekoch/gouniversal/module/gasprice/global"
	"github.com/dekoch/gouniversal/module/gasprice/lang"
	"github.com/dekoch/gouniversal/module/gasprice/station"
	"github.com/dekoch/gouniversal/module/gasprice/typemd"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"

	"github.com/google/uuid"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "App:GasPrice:Settings:Station:Edit", page.Lang.StationEdit.Title)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	type Content struct {
		Lang    lang.StationEdit
		UUID    template.HTML
		Name    template.HTML
		Company template.HTML
		Street  template.HTML
		City    template.HTML
		URL     template.HTML
	}
	var c Content

	c.Lang = page.Lang.StationEdit

	// Form input
	id := nav.Parameter("UUID")

	switch r.FormValue("edit") {
	case "apply":
		err := editStation(r, id)
		if err == nil {
			nav.RedirectPath("App:GasPrice:Settings:Station:List", false)
		} else {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
		}

	case "delete":
		err := deleteStation(id)
		if err == nil {
			nav.RedirectPath("App:GasPrice:Settings:Station:List", false)
		} else {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
		}

	default:
		if id == "new" {

			id = newStation()
			nav.RedirectPath(strings.Replace(nav.Path, "UUID=new", "UUID="+id, 1), false)
		}
	}

	station, err := global.Config.Stations.GetStation(id)
	c.UUID = template.HTML(station.UUID)
	c.Name = template.HTML(station.Name)
	c.Company = template.HTML(station.Company)
	c.Street = template.HTML(station.Street)
	c.City = template.HTML(station.City)
	c.URL = template.HTML(station.URL)

	p, err := functions.PageToString(global.Config.UIFileRoot+"settings/stationedit.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func newStation() string {

	id := uuid.Must(uuid.NewRandom()).String()

	var station station.Station
	station.UUID = id
	station.Name = id

	global.Config.Stations.Add(station)
	global.Config.SaveConfig()

	return id
}

func editStation(r *http.Request, uid string) error {

	var (
		err     error
		name    string
		company string
		street  string
		city    string
		ur      string
		station station.Station
	)

	func() {

		for i := 0; i <= 9; i++ {

			switch i {
			case 0:
				name, err = functions.CheckFormInput("Name", r)

			case 1:
				company, err = functions.CheckFormInput("Company", r)

			case 2:
				street, err = functions.CheckFormInput("Street", r)

			case 3:
				city, err = functions.CheckFormInput("City", r)

			case 4:
				ur = r.FormValue("Url")

			case 5:
				_, err = url.Parse(ur)

			case 6:
				// check input
				if functions.IsEmpty(name) ||
					functions.IsEmpty(ur) {

					err = errors.New("bad input")
				}

			case 7:
				station, err = global.Config.Stations.GetStation(uid)

			case 8:
				station.Name = name
				station.Company = company
				station.Street = street
				station.City = city
				station.URL = ur

			case 9:
				global.Config.Stations.Add(station)
				err = global.Config.SaveConfig()
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func deleteStation(uid string) error {

	global.Config.Stations.Remove(uid)

	return global.Config.SaveConfig()
}
