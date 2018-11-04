package meshFileSync

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/dekoch/gouniversal/modules/mesh"
	"github.com/dekoch/gouniversal/modules/meshFileSync/client"
	"github.com/dekoch/gouniversal/modules/meshFileSync/global"
	"github.com/dekoch/gouniversal/modules/meshFileSync/lang"
	"github.com/dekoch/gouniversal/modules/meshFileSync/server"
	"github.com/dekoch/gouniversal/modules/meshFileSync/ui"
	"github.com/dekoch/gouniversal/shared/language"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	rand.Seed(time.Now().UnixNano())

	global.Config.LoadConfig()

	en := lang.DefaultEn()
	global.Lang = language.New(global.Config.LangFileRoot, en, "en")

	global.LocalFiles.SetPath(global.Config.FileRoot)
	global.LocalFiles.SetServerID(mesh.GetServerInfo().ID)
	global.LocalFiles.AddList(global.Config.LocalFiles)
	global.LocalFiles.Scan()

	server.LoadConfig()
	client.LoadConfig()
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	ui.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	ui.Render(page, nav, r)
}

func Exit() {

	global.Config.LocalFiles = global.LocalFiles.Get()
	global.Config.SaveConfig()
}
