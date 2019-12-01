package monmotion

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/monmotion/global"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	global.Config.LoadConfig()

	for i := range global.Config.Cam {

		global.Config.Cam[i].LoadConfig()
	}
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

}

func Exit(em *types.ExitMessage) {

	err := global.Config.SaveConfig()
	if err != nil {
		console.Log(err, "")
		return
	}
}
