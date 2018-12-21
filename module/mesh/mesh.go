package mesh

import (
	"net/http"

	"github.com/dekoch/gouniversal/module/mesh/client"
	"github.com/dekoch/gouniversal/module/mesh/global"
	"github.com/dekoch/gouniversal/module/mesh/lang"
	"github.com/dekoch/gouniversal/module/mesh/server"
	"github.com/dekoch/gouniversal/module/mesh/serverinfo"
	"github.com/dekoch/gouniversal/module/mesh/typemesh"
	"github.com/dekoch/gouniversal/module/mesh/ui"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/language"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"
)

func LoadConfig() {

	global.Config.LoadConfig()

	en := lang.DefaultEn()
	global.Lang = language.New(global.Config.LangFileRoot, en, "en")

	global.Keyfile.LoadConfig()
	global.NetworkConfig.LoadConfig()

	if global.NetworkConfig.Network.CheckKey(global.Keyfile.GetKey()) == false {
		console.Log("key<->hash mismatch", " ")
	}

	global.NetworkConfig.Network.SetKey(global.Keyfile.GetKey())
	global.NetworkConfig.Add(global.Config.Server)
	global.NetworkConfig.ServerList.SetMaxAge(global.NetworkConfig.Network.GetMaxClientAge())

	server.LoadConfig()
	client.LoadConfig()
}

func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	ui.RegisterPage(page, nav)
}

func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	ui.Render(page, nav, r)
}

func GetServerInfo() serverinfo.ServerInfo {
	global.Config.Server.Update()
	return global.Config.Server.Get()
}

func GetServerList() []serverinfo.ServerInfo {
	return global.NetworkConfig.ServerList.Get()
}

func GetServerWithID(id string) (serverinfo.ServerInfo, error) {
	return global.NetworkConfig.ServerList.GetWithID(id)
}

func SendMessage(output typemesh.ServerMessage) error {
	return client.SendMessage(output)
}

func NewMessage(receiver serverinfo.ServerInfo, t typemesh.MessageType, ver float32, c []byte) error {

	var output typemesh.ServerMessage

	output.Receiver = receiver
	output.Message.Type = t
	output.Message.Version = ver
	output.Message.Content = c

	return client.SendMessage(output)
}

func IsLoop(in serverinfo.ServerInfo) bool {
	return client.IsLoop(in)
}

func Exit() {

	global.Config.SaveConfig()

	global.NetworkConfig.ServerList.Clean()
	global.NetworkConfig.SaveConfig()
}
