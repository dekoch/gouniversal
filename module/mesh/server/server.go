package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"strconv"
	"time"

	"github.com/dekoch/gouniversal/build"
	"github.com/dekoch/gouniversal/module/mesh/global"
	"github.com/dekoch/gouniversal/module/mesh/serverinfo"
	"github.com/dekoch/gouniversal/module/mesh/settings"
	"github.com/dekoch/gouniversal/module/mesh/typemesh"
	"github.com/dekoch/gouniversal/shared/aes"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/sbool"

	meshFSServer "github.com/dekoch/gouniversal/module/meshfilesync/server"
)

const debug = false

type Server struct{}

var ln net.Listener
var started sbool.Sbool
var restart sbool.Sbool

func LoadConfig() {

	started.UnSet()
	restart.UnSet()

	go start()
}

func start() {

	defer started.UnSet()

	console.Log("mesh network listening on port: "+strconv.Itoa(global.Config.Server.GetPort()), " ")
	console.Log("mesh ID: "+global.Config.Server.ID, " ")

	rpc.Register(new(Server))

	var err error

	ln, err = net.Listen("tcp", ":"+strconv.Itoa(global.Config.Server.GetPort()))
	if err != nil {
		console.Log(err, " ")
		return
	}

	func() {

		started.Set()

		for {
			if restart.IsSet() {
				restart.UnSet()
				return
			}

			c, err := ln.Accept()
			if err != nil {
				continue
			}

			go rpc.ServeConn(c)
		}
	}()

	console.Log("mesh closed", " ")
}

func Restart() {

	console.Log("mesh restart", " ")

	if started.IsSet() {
		ln.Close()
	}

	restart.Set()

	for started.IsSet() {
		time.Sleep(100 * time.Millisecond)
	}

	restart.UnSet()

	go start()
}

func (this *Server) Message(input typemesh.ServerMessage, output *string) error {

	var (
		err error
		key []byte
	)

	server := global.Config.Server.Get()

	func() {

		for i := 0; i <= 8; i++ {

			switch i {
			case 0:
				if input.Receiver.ID != server.ID {
					err = errors.New("ServerWrongReceiver")
				}

			case 1:
				// check IDs
				if settings.LocalConnection == false {

					if input.Sender.ID == server.ID {
						err = errors.New("ServerSameID")
					}
				}

			case 2:
				if global.NetworkConfig.Network.CheckID(input.Network.ID) == false {
					err = errors.New("ServerDifferentMeshID")
				}

			case 3:
				if global.NetworkConfig.Network.CheckHashWithLocalKey(input.Network.Hash) == false {
					err = errors.New("ServerWrongMeshKey")
				}

			case 4:
				if debug {
					writeDebug(input.Message.Type, input.Sender.ID)
				}

			case 5:
				switch input.Message.Type {
				case typemesh.MessHello:

					global.NetworkConfig.ServerList.Add(input.Sender)
					return
				}

			case 6:
				key, err = global.Keyfile.GetKey()

			case 7:
				// decrypt message content
				input.Message.Content, err = aes.Decrypt(key, input.Message.Content)

			case 8:
				switch input.Message.Type {
				case typemesh.MessAnnounce:

					if input.Message.Version == 1.0 {
						announce(input)
					}

				case typemesh.MessFileSync:

					if build.ModuleMeshFS {
						if input.Message.Version == 1.0 {
							err = meshFSServer.Server(input)
						}
					} else {
						err = errors.New("ServerModuleDisabled")
					}
				}
			}

			if err != nil {
				if debug {
					console.Output(err, "mesh")
				}

				return
			}
		}
	}()

	if err == nil {
		err = errors.New("nil")
	}

	*output = err.Error()

	return nil
}

func announce(input typemesh.ServerMessage) {

	// update network
	global.NetworkConfig.Network.Update(input.Network)

	// update server list
	var newList []serverinfo.ServerInfo

	err := json.Unmarshal(input.Message.Content, &newList)
	if err != nil {
		fmt.Println(err)
	} else {
		global.NetworkConfig.ServerList.SetMaxAge(global.NetworkConfig.Network.GetMaxClientAge())
		global.NetworkConfig.ServerList.AddList(newList)
	}
}

func writeDebug(t typemesh.MessageType, id string) {

	var s string

	switch t {
	case typemesh.MessAnnounce:
		s = "announce"

	case typemesh.MessHello:
		s = "hello"

	case typemesh.MessRAW:
		s = "raw"

	case typemesh.MessMessenger:
		s = "messenger"

	case typemesh.MessFileSync:
		s = "meshFileSync"

	default:
		s = "unknown"
	}

	s += " from " + id

	console.Output(s, "mesh")
}
