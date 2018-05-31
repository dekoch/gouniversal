package modules

import (
	"net/http"

	"github.com/dekoch/gouniversal/modules/console"
	"github.com/dekoch/gouniversal/modules/fileshare"
	"github.com/dekoch/gouniversal/modules/homepage"
	"github.com/dekoch/gouniversal/modules/modbustest"
	"github.com/dekoch/gouniversal/modules/openespm"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"

	sharedConsole "github.com/dekoch/gouniversal/shared/console"
)

// Modules provide a interface to nest apps and modules

const modConsole = true
const modOpenESPM = true
const modFileshare = true
const modHomepage = false
const modModbusTest = false

// LoadConfig is called before UI starts
func LoadConfig() {

	if modConsole {
		sharedConsole.Log("Console enabled", "Module")
		console.LoadConfig()
	}

	if modOpenESPM {
		sharedConsole.Log("openESPM enabled", "Module")
		openespm.LoadConfig()
	}

	if modFileshare {
		sharedConsole.Log("Fileshare enabled", "Module")
		fileshare.LoadConfig()
	}

	if modHomepage {
		sharedConsole.Log("Homepage enabled", "Module")
		homepage.LoadConfig()
	}

	if modModbusTest {
		sharedConsole.Log("ModbusTest enabled", "Module")
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

	if modConsole {
		console.RegisterPage(page, nav)
	}

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

	if modConsole {
		if nav.IsNext("Console") {

			console.Render(page, nav, r)
		}
	}

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

	if modConsole {
		console.Exit()
	}

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
