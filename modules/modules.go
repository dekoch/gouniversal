package modules

import (
	"gouniversal/modules/fileshare"
	"gouniversal/modules/homepage"
	"gouniversal/modules/modbustest"
	"gouniversal/modules/openespm"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"net/http"
)

// Modules provide a interface to nest apps and modules

const modOpenESPM = true
const modFileshare = true
const modHomepage = false
const modModbusTest = false

// LoadConfig is called before UI starts
func LoadConfig() {

	if modOpenESPM {
		openespm.LoadConfig()
	}

	if modFileshare {
		fileshare.LoadConfig()
	}

	if modHomepage {
		homepage.LoadConfig()
	}

	if modModbusTest {
		modbustest.LoadConfig()
	}
}

// RegisterPage adds pages to a sitemap
// use
// nav.Sitemap.Register("App:MyApp", "MyApp")
// nav.Sitemap.Register("App:Program:MyApp", "MyApp")
// nav.Sitemap.Register("App:Account:MyApp", "MyApp")
//
// to nest your app or module inside menus
func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	if modOpenESPM {
		openespm.RegisterPage(page, nav)
	}

	if modFileshare {
		fileshare.RegisterPage(page, nav)
	}

	if modHomepage {
		homepage.RegisterPage(page, nav)
	}
}

// Render UI page
func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	if modOpenESPM {
		if nav.IsNext("openESPM") {

			openespm.Render(page, nav, r)
		}
	}

	if modFileshare {
		if nav.IsNext("Fileshare") {

			fileshare.Render(page, nav, r)
		}
	}

	if modHomepage {
		if nav.IsNext("Homepage") {

			homepage.Render(page, nav, r)
		}
	}
}

// Exit is called before program exit
func Exit() {

	if modOpenESPM {
		openespm.Exit()
	}

	if modFileshare {
		fileshare.Exit()
	}

	if modHomepage {
		homepage.Exit()
	}
}
