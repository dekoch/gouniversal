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
	"github.com/dekoch/gouniversal/modules/mesh/global"
	"github.com/dekoch/gouniversal/modules/mesh/serverInfo"
	"github.com/dekoch/gouniversal/modules/mesh/settings"
	"github.com/dekoch/gouniversal/modules/mesh/typesMesh"
	"github.com/dekoch/gouniversal/shared/aes"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/google/uuid"
	"github.com/tevino/abool"

	meshFSServer "github.com/dekoch/gouniversal/modules/meshfilesync/server"
)

const debug = false

type Server struct{}

var ln net.Listener
var started *abool.AtomicBool
var restart *abool.AtomicBool

func LoadConfig() {

	restart = abool.New()
	restart.UnSet()

	started = abool.New()
	started.UnSet()

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

func (this *Server) Message(input typesMesh.ServerMessage, output *string) error {

	var err error

	server := global.Config.Server.Get()

	if global.NetworkConfig.Network.CheckID(input.Network.ID) == false {
		err = errors.New("ServerDifferentMeshID")
	} else if global.NetworkConfig.Network.CheckHashWithLocalKey(input.Network.Hash) == false {
		err = errors.New("ServerWrongMeshKey")
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

			if input.Message.Type == typesMesh.MessAnnounce {
				if input.Message.Version == 1.0 {
					announce(input)
				}
			}

			if input.Receiver.ID != server.ID {
				err = errors.New("ServerWrongReceiver")
			} else {

				switch input.Message.Type {
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

	var s string

	switch t {
	case typesMesh.MessAnnounce:
		s = "announce"

	case typesMesh.MessHello:
		s = "hello"

	case typesMesh.MessRAW:
		s = "raw"

	case typesMesh.MessMessenger:
		s = "messenger"

	case typesMesh.MessFileSync:
		s = "meshFileSync"

	default:
		s = "unknown"
	}

	s += " from " + id

	console.Output(s, "mesh")
}
