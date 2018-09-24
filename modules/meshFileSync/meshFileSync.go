package meshFileSync

import (
	"math/rand"
	"time"

	"github.com/dekoch/gouniversal/modules/mesh"
	"github.com/dekoch/gouniversal/modules/meshFileSync/client"
	"github.com/dekoch/gouniversal/modules/meshFileSync/global"
	"github.com/dekoch/gouniversal/modules/meshFileSync/server"
)

func LoadConfig() {

	rand.Seed(time.Now().UnixNano())

	global.Config.LoadConfig()

	global.LocalFiles.SetPath(global.Config.FileRoot)
	global.LocalFiles.SetServerID(mesh.GetServerInfo().ID)
	global.LocalFiles.AddList(global.Config.LocalFiles)
	global.LocalFiles.Scan()

	server.LoadConfig()
	client.LoadConfig()
}

func Exit() {

	global.Config.LocalFiles = global.LocalFiles.Get()
	global.Config.SaveConfig()
}
