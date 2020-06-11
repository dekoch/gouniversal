package instafollowbot

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/instafollowbot/core"
	"github.com/dekoch/gouniversal/module/instafollowbot/global"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	global.Config.LoadConfig()

	coreIDs := global.Config.GetCoreList()

	for i := range coreIDs {

		go func(uid string) {

			conf, err := global.Config.GetCoreConfig(uid)
			if err != nil {
				return
			}

			var co core.Core
			co.LoadConfig(conf)
		}(coreIDs[i])
	}
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	//ui.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	//ui.Render(page, nav, r)
}

func Exit(em *types.ExitMessage) {

	err := global.Config.Exit(em)
	if err != nil {
		console.Log(err, "")
		return
	}

	err = global.Config.SaveConfig()
	if err != nil {
		console.Log(err, "")
	}
}
