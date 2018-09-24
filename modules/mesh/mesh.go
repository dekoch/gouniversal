package mesh

import (
	"github.com/dekoch/gouniversal/modules/mesh/client"
	"github.com/dekoch/gouniversal/modules/mesh/global"
	"github.com/dekoch/gouniversal/modules/mesh/server"
	"github.com/dekoch/gouniversal/modules/mesh/serverInfo"
	"github.com/dekoch/gouniversal/modules/mesh/typesMesh"
	"github.com/dekoch/gouniversal/shared/console"
)

func LoadConfig() {

	global.Config.LoadConfig()
	global.Keyfile.LoadConfig()
	global.NetworkConfig.LoadConfig()

	if global.NetworkConfig.Network.CheckKey(global.Keyfile.GetKey()) == false {
		console.Log("key<->hash mismatch", " ")
	}

	global.NetworkConfig.Network.SetKey(global.Keyfile.GetKey())
	global.NetworkConfig.Add(global.Config.Server)
	global.NetworkConfig.ServerList.SetMaxAge(global.NetworkConfig.Network.GetMaxClientAge())

	if global.Config.ServerEnabled {
		server.LoadConfig()
	}

	if global.Config.ClientEnabled {
		client.LoadConfig()
	}
}

func GetServerInfo() serverInfo.ServerInfo {
	global.Config.Server.Update()
	return global.Config.Server.Get()
}

func GetServerList() []serverInfo.ServerInfo {
	return global.NetworkConfig.ServerList.Get()
}

func GetServerWithID(id string) (serverInfo.ServerInfo, error) {
	return global.NetworkConfig.ServerList.GetWithID(id)
}

func SendMessage(output typesMesh.ServerMessage) error {
	return client.SendMessage(output)
}

func NewMessage(receiver serverInfo.ServerInfo, t typesMesh.MessageType, ver float32, c []byte) error {

	var output typesMesh.ServerMessage

	output.Receiver = receiver
	output.Message.Type = t
	output.Message.Version = ver
	output.Message.Content = c

	return client.SendMessage(output)
}

func IsLoop(in serverInfo.ServerInfo) bool {
	return client.IsLoop(in)
}

func Exit() {

	global.Config.SaveConfig()

	global.NetworkConfig.ServerList.Clean()
	global.NetworkConfig.SaveConfig()
}
