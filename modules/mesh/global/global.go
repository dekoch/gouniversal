package global

import (
	"github.com/dekoch/gouniversal/modules/mesh/keyfile"
	"github.com/dekoch/gouniversal/modules/mesh/moduleConfig"
	"github.com/dekoch/gouniversal/modules/mesh/networkConfig"
	"github.com/dekoch/gouniversal/modules/mesh/typesMesh"
)

var (
	Config        moduleConfig.ModuleConfig
	NetworkConfig networkConfig.NetworkConfig
	Keyfile       keyfile.Keyfile

	ChanMessenger = make(chan typesMesh.ServerMessage)
)
