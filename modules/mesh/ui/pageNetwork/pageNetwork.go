package pageNetwork

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/dekoch/gouniversal/modules/mesh/global"
	"github.com/dekoch/gouniversal/modules/mesh/lang"
	"github.com/dekoch/gouniversal/modules/mesh/typesMesh"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typesMesh.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.Title, "App:Mesh:Network", page.Lang.Network.Title)
}

func Render(page *typesMesh.Page, nav *navigation.Navigation, r *http.Request) {

	button := r.FormValue("edit")

	type Content struct {
		Lang             lang.Network
		ID               template.HTML
		AnnounceInterval template.HTML
		HelloInterval    template.HTML
		MaxClientAge     template.HTML
		ServerList       template.HTML
	}
	var c Content

	c.Lang = page.Lang.Network

	if button == "apply" {

		err := edit(r)
		if err != nil {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
		}
	}

	c.ID = template.HTML(global.NetworkConfig.Network.GetID())
	c.AnnounceInterval = template.HTML(strconv.Itoa(global.NetworkConfig.Network.GetAnnounceIntervalInt()))
	c.HelloInterval = template.HTML(strconv.Itoa(global.NetworkConfig.Network.GetHelloIntervalInt()))
	c.MaxClientAge = template.HTML(strconv.FormatFloat(global.NetworkConfig.Network.GetMaxClientAge(), 'f', 0, 64))

	list := ""

	for _, server := range global.NetworkConfig.Get() {
		list += "<tr>"
		list += "<td>" + server.ID + "</td>"
		list += "<td>" + server.TimeStamp.Format("2006-01-02 15:04:05") + "</td>"
		list += "</tr>"
	}

	c.ServerList = template.HTML(list)

	p, err := functions.PageToString(global.Config.UIFileRoot+"network.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func edit(r *http.Request) error {

	var (
		err               error
		sAnnounceInterval string
		iAnnounceInterval int
		sHelloInterval    string
		iHelloInterval    int
		sMaxClientAge     string
		fMaxClientAge     float64
	)

	func() {

		for i := 0; i <= 8; i++ {

			switch i {
			case 0:
				sAnnounceInterval, err = functions.CheckFormInput("AnnounceInterval", r)

			case 1:
				sHelloInterval, err = functions.CheckFormInput("HelloInterval", r)

			case 2:
				sMaxClientAge, err = functions.CheckFormInput("MaxClientAge", r)

			case 3:
				// check input
				if functions.IsEmpty(sAnnounceInterval) ||
					functions.IsEmpty(sHelloInterval) ||
					functions.IsEmpty(sMaxClientAge) {

					err = errors.New("bad input")
				}

			case 4:
				iAnnounceInterval, err = strconv.Atoi(sAnnounceInterval)

			case 5:
				iHelloInterval, err = strconv.Atoi(sHelloInterval)

			case 6:
				fMaxClientAge, err = strconv.ParseFloat(sMaxClientAge, 64)

			case 7:
				// check converted input
				if iAnnounceInterval < 1 ||
					iAnnounceInterval > 900 ||
					iHelloInterval < 0 ||
					iHelloInterval > 900 ||
					fMaxClientAge < 1.0 ||
					fMaxClientAge > 365.0 {

					err = errors.New("bad input")
				}

			case 8:
				global.NetworkConfig.Network.SetAnnounceInterval(iAnnounceInterval)
				global.NetworkConfig.Network.SetHelloInterval(iHelloInterval)
				global.NetworkConfig.Network.SetMaxClientAge(fMaxClientAge)
				global.NetworkConfig.Network.SetTimeStamp(time.Now())

				global.NetworkConfig.ServerList.SetMaxAge(fMaxClientAge)

				err = global.NetworkConfig.SaveConfig()
			}

			if err != nil {
				console.Log(err, "")
				return
			}
		}
	}()

	return err
}
