package homepage

import (
	"net/http"

	"github.com/dekoch/gouniversal/modules/homepage/global"
	"github.com/dekoch/gouniversal/modules/homepage/ui"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	global.Config.LoadConfig()

	ui.LoadConfig()
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	ui.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	ui.Render(page, nav, r)
}

func Exit() {

}
