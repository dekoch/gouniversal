package modules

import (
	"gouniversal/modules/openespm"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"net/http"
)

type Modules struct{}

const modOpenESPM = true

func (m *Modules) LoadConfig() {

}

func (m *Modules) RegisterPage(page *types.Page, nav *navigation.Navigation) {

	if modOpenESPM {
		openespm.RegisterPage(page, nav)
	}
}

func (m *Modules) Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	if modOpenESPM {
		openespm.Render(page, nav, r)
	}
}

func (m *Modules) Exit() {

}
