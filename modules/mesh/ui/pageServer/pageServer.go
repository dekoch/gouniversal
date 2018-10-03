package pageServer

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/dekoch/gouniversal/modules/mesh/global"
	"github.com/dekoch/gouniversal/modules/mesh/lang"
	"github.com/dekoch/gouniversal/modules/mesh/server"
	"github.com/dekoch/gouniversal/modules/mesh/typesMesh"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typesMesh.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.Title, "App:Mesh:Server", page.Lang.Server.Title)
}

func Render(page *typesMesh.Page, nav *navigation.Navigation, r *http.Request) {

	button := r.FormValue("edit")

	type Content struct {
		Lang             lang.Server
		ID               template.HTML
		Port             template.HTML
		PubAddrUpdInterv template.HTML
		Addresses        template.HTML
	}
	var c Content

	c.Lang = page.Lang.Server

	if button == "apply" {

		err := edit(r)
		if err != nil {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
		}
	}

	this := global.Config.Server.Get()

	c.ID = template.HTML(this.ID)
	c.Port = template.HTML(strconv.Itoa(this.Port))
	c.PubAddrUpdInterv = template.HTML(strconv.Itoa(global.Config.PubAddrUpdInterv))

	addr := ""

	for _, a := range this.Address {
		addr += a + "<br>"
	}

	c.Addresses = template.HTML(addr)

	p, err := functions.PageToString(global.Config.UIFileRoot+"server.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func edit(r *http.Request) error {

	var (
		err               error
		sPubAddrUpdInterv string
		iPubAddrUpdInterv int
		sPort             string
		iPort             int
		restartServer     bool
	)

	func() {

		for i := 0; i <= 7; i++ {

			switch i {
			case 0:
				sPubAddrUpdInterv, err = functions.CheckFormInput("PubAddrUpdInterv", r)

			case 1:
				sPort, err = functions.CheckFormInput("Port", r)

			case 2:
				// check input
				if functions.IsEmpty(sPort) ||
					functions.IsEmpty(sPubAddrUpdInterv) {

					err = errors.New("bad input")
				}

			case 3:
				iPubAddrUpdInterv, err = strconv.Atoi(sPubAddrUpdInterv)

			case 4:
				iPort, err = strconv.Atoi(sPort)

			case 5:
				// check converted input
				if iPubAddrUpdInterv < 0 ||
					iPubAddrUpdInterv > 1440 ||
					iPort < 1 ||
					iPort > 65535 {

					err = errors.New("bad input")
				}

			case 6:
				if global.Config.Server.GetPort() != iPort {

					global.Config.Server.SetPort(iPort)
					restartServer = true
				}

				global.Config.PubAddrUpdInterv = iPubAddrUpdInterv
				global.Config.Server.SetPubAddrUpdInterv(iPubAddrUpdInterv)

				err = global.Config.SaveConfig()

			case 7:
				if restartServer {
					server.Restart()
				}
			}

			if err != nil {
				console.Log(err, "")
				return
			}
		}
	}()

	return err
}
