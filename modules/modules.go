package modules

import (
	"gouniversal/modules/openespm"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"net/http"
)

// Modules provide a interface to nest apps and modules
type Modules struct{}

const modOpenESPM = true

// LoadConfig is called before UI starts
func (m *Modules) LoadConfig() {

	if modOpenESPM {
		openespm.LoadConfig()
	}
}

// RegisterPage adds pages to a sitemap
// use
// nav.Sitemap.Register("App:MyApp", "MyApp")
// nav.Sitemap.Register("App:Program:MyApp", "MyApp")
// nav.Sitemap.Register("App:Account:MyApp", "MyApp")
//
// to nest your app or module inside menus
func (m *Modules) RegisterPage(page *types.Page, nav *navigation.Navigation) {

	if modOpenESPM {
		openespm.RegisterPage(page, nav)
	}
}

// Render UI page
func (m *Modules) Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	if modOpenESPM {
		if nav.IsNext("openESPM") {

			openespm.Render(page, nav, r)
		}
	}
}

// Exit is called before program exit
func (m *Modules) Exit() {

	if modOpenESPM {
		openespm.Exit()
	}
}
