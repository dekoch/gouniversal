package pageserver

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/dekoch/gouniversal/module/mesh/global"
	"github.com/dekoch/gouniversal/module/mesh/lang"
	"github.com/dekoch/gouniversal/module/mesh/server"
	"github.com/dekoch/gouniversal/module/mesh/typemesh"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemesh.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.Title, "App:Mesh:Server", page.Lang.Server.Title)
}

func Render(page *typemesh.Page, nav *navigation.Navigation, r *http.Request) {

	var err error

	type Content struct {
		Lang             lang.Server
		ID               template.HTML
		Port             template.HTML
		ExposePort       template.HTML
		SetManualAddress template.HTML
		PubAddrUpdInterv template.HTML
		Addresses        template.HTML
	}
	var c Content

	c.Lang = page.Lang.Server

	switch r.FormValue("edit") {
	case "apply":
		err = edit(r)
	}

	if err != nil {
		alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
	}

	global.Config.Server.Update()
	this := global.Config.Server.Get()

	c.ID = template.HTML(this.ID)
	c.Port = template.HTML(strconv.Itoa(this.Port))

	if this.ExposePort > 0 {
		c.ExposePort = template.HTML(strconv.Itoa(this.ExposePort))
	} else {
		c.ExposePort = template.HTML("")
	}

	c.SetManualAddress = template.HTML(global.Config.GetManualAddress())
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
		sExposePort       string
		iExposePort       int
		address           string
		restartServer     bool
	)

	func() {

		for i := 0; i <= 13; i++ {

			switch i {
			case 0:
				sPubAddrUpdInterv, err = functions.CheckFormInput("PubAddrUpdInterv", r)

			case 1:
				sPort, err = functions.CheckFormInput("Port", r)

			case 2:
				sExposePort, err = functions.CheckFormInput("ExposePort", r)

			case 3:
				address, err = functions.CheckFormInput("SetManualAddress", r)

			case 4:
				// check input
				if functions.IsEmpty(sPort) ||
					functions.IsEmpty(sPubAddrUpdInterv) {

					err = errors.New("bad input")
				}

			case 5:
				iPubAddrUpdInterv, err = strconv.Atoi(sPubAddrUpdInterv)

			case 6:
				iPort, err = strconv.Atoi(sPort)

			case 7:
				if functions.IsEmpty(sExposePort) {
					iExposePort = -1
				} else {
					iExposePort, err = strconv.Atoi(sExposePort)
				}

			case 8:
				// check converted input
				if iPubAddrUpdInterv < 0 ||
					iPubAddrUpdInterv > 1440 ||
					iPort < 1 ||
					iPort > 65535 ||
					iExposePort < -1 ||
					iExposePort > 65535 {

					err = errors.New("bad input")
				}

			case 9:
				if global.Config.Server.GetPort() != iPort {

					global.Config.Server.SetPort(iPort)
					restartServer = true
				}

			case 10:
				if iExposePort <= 0 {
					iExposePort = -1
				}

				global.Config.Server.SetExposePort(iExposePort)

			case 11:
				global.Config.PubAddrUpdInterv = iPubAddrUpdInterv
				global.Config.Server.SetPubAddrUpdInterv(iPubAddrUpdInterv)

			case 12:
				address = strings.Trim(address, " ")
				global.Config.SetManualAddress(address)
				global.Config.Server.SetManualAddress(address)

				err = global.Config.SaveConfig()

			case 13:
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
