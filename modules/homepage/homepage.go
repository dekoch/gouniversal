package homepage

import (
	"gouniversal/modules/homepage/ui"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"net/http"
)

func LoadConfig() {

}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	ui.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	ui.Render(page, nav, r)
}

func Exit() {

}
