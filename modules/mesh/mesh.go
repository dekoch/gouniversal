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

func SendMessage(output typesMesh.ServerMessage) typesMesh.ServerMessage {
	return client.SendMessage(output)
}

func Exit() {

	global.Config.SaveConfig()
	global.NetworkConfig.SaveConfig()
}
