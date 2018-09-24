package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"strconv"

	"github.com/dekoch/gouniversal/build"
	"github.com/dekoch/gouniversal/modules/mesh/global"
	"github.com/dekoch/gouniversal/modules/mesh/serverInfo"
	"github.com/dekoch/gouniversal/modules/mesh/settings"
	"github.com/dekoch/gouniversal/modules/mesh/typesMesh"
	"github.com/dekoch/gouniversal/shared/aes"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/google/uuid"

	meshFSServer "github.com/dekoch/gouniversal/modules/meshFileSync/server"
)

const debug = false

type Server struct{}

func LoadConfig() {

	go start()
}

func start() {

	console.Log("mesh network listening on port: "+strconv.Itoa(global.Config.Server.GetPort()), " ")
	console.Log("mesh ID: "+global.Config.Server.ID, " ")

	rpc.Register(new(Server))

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(global.Config.Server.GetPort()))
	if err != nil {
		console.Log(err, " ")
		return
	}

	for {
		c, err := ln.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(c)
	}
}

func (this *Server) Message(input typesMesh.ServerMessage, output *string) error {

	var err error

	server := global.Config.Server.Get()

	if global.NetworkConfig.Network.CheckID(input.Network.ID) == false {
		err = errors.New("ServerDifferentMeshID")
	} else if global.NetworkConfig.Network.CheckHashWithLocalKey(input.Network.Hash) == false {
		err = errors.New("ServerWrongMeshKey")
	} else if input.Receiver.ID != server.ID {
		err = errors.New("ServerWrongReceiver")
	}

	if err == nil {
		if settings.LocalConnection == false {
			// check IDs, if we have the same inside a network, change own
			if input.Sender.ID == server.ID {

				global.NetworkConfig.ServerList.Delete(server.ID)

				u := uuid.Must(uuid.NewRandom())
				fmt.Println("change ID to " + u.String())
				global.Config.Server.SetID(u.String())
			}
		}

		// decrypt message content
		b, err := aes.Decrypt(global.Keyfile.GetKey(), string(input.Message.Content))
		if err != nil {
			fmt.Println(err)
		} else {
			input.Message.Content = []byte(b)

			if debug {
				writeDebug(input.Message.Type, input.Sender.ID)
			}

			switch input.Message.Type {
			case typesMesh.MessAnnounce:

				if input.Message.Version == 1.0 {
					announce(input)
				}

			case typesMesh.MessHello:

				global.NetworkConfig.ServerList.Add(input.Sender)

			case typesMesh.MessFileSync:

				if build.ModuleMeshFS {
					if input.Message.Version == 1.0 {

						err = meshFSServer.Server(input)
					}
				} else {
					err = errors.New("ServerModuleDisabled")
				}
			}
		}
	}

	if err == nil {
		err = errors.New("nil")
	}

	*output = err.Error()

	return nil
}

func announce(input typesMesh.ServerMessage) {

	// update network
	global.NetworkConfig.Network.Update(input.Network)

	// update server list
	var newList []serverInfo.ServerInfo

	err := json.Unmarshal(input.Message.Content, &newList)
	if err != nil {
		fmt.Println(err)
	} else {
		global.NetworkConfig.ServerList.SetMaxAge(global.NetworkConfig.Network.GetMaxClientAge())
		global.NetworkConfig.ServerList.AddList(newList)
	}
}

func writeDebug(t typesMesh.MessageType, id string) {

	switch t {
	case typesMesh.MessAnnounce:
		fmt.Print("announce")

	case typesMesh.MessHello:
		fmt.Print("hello")

	case typesMesh.MessRAW:
		fmt.Print("raw")

	case typesMesh.MessMessenger:
		fmt.Print("messenger")

	case typesMesh.MessFileSync:
		fmt.Print("meshFileSync")

	default:
		fmt.Print("unknown")
	}

	fmt.Println(" from \"" + id + "\"")
}
