package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/rpc"
	"strconv"

	"github.com/dekoch/gouniversal/modules/mesh/global"
	"github.com/dekoch/gouniversal/modules/mesh/serverInfo"
	"github.com/dekoch/gouniversal/modules/mesh/settings"
	"github.com/dekoch/gouniversal/modules/mesh/typesMesh"
	"github.com/dekoch/gouniversal/shared/aes"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/google/uuid"
)

type Server struct{}

func LoadConfig() {

	go start()
}

func start() {

	console.Log("mesh network listening on port "+strconv.Itoa(global.Config.Server.GetPort()), " ")

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

func (this *Server) Message(input typesMesh.ServerMessage, output *typesMesh.ServerMessage) error {

	// update server info
	global.Config.Server.Update()

	var out typesMesh.ServerMessage
	out.Sender = global.Config.Server.Get()
	out.Receiver = input.Sender
	out.Network = global.NetworkConfig.Network.Get()
	out.Error = typesMesh.ErrNil
	out.Message.Type = typesMesh.MessNil

	if settings.LocalConnection == false {
		// check IDs, if we have the same inside a network, change
		if input.Sender.ID == out.Sender.ID {
			fmt.Println("change ID")

			global.NetworkConfig.ServerList.Delete(out.Sender.ID)

			u := uuid.Must(uuid.NewRandom())
			global.Config.Server.SetID(u.String())
			out.Sender = global.Config.Server.Get()
		}
	}

	if global.NetworkConfig.Network.CheckID(input.Network.ID) == false {
		out.Error = typesMesh.ErrServerDifferentMeshID
	} else if global.NetworkConfig.Network.CheckHashWithLocalKey(input.Network.Hash) == false {
		out.Error = typesMesh.ErrServerWrongMeshKey
	} else if input.Receiver.ID != out.Sender.ID {
		out.Error = typesMesh.ErrServerWrongReceiver
	}

	switch out.Error {
	case typesMesh.ErrNil:
		// decrypt message content
		b, err := aes.Decrypt(global.Keyfile.GetKey(), string(input.Message.Content))
		if err != nil {
			fmt.Println(err)
			out.Error = typesMesh.ErrServerDecryption
		} else {
			input.Message.Content = []byte(b)

			switch input.Message.Type {
			case typesMesh.MessAnnounce:
				fmt.Print("announce")

			case typesMesh.MessHello:
				fmt.Print("hello")

			case typesMesh.MessRAW:
				fmt.Print("raw")

			case typesMesh.MessMessenger:
				fmt.Print("messenger")

			default:
				fmt.Print("unknown")
			}

			fmt.Println(" from \"" + input.Sender.ID + "\"")

			switch input.Message.Type {
			case typesMesh.MessAnnounce:

				if input.Message.Version == 1.0 {
					announce(input, &out)
				}

			case typesMesh.MessHello:

				global.NetworkConfig.ServerList.Add(input.Sender)

			case typesMesh.MessMessenger:

				if input.Message.Version == 1.0 {
					global.ChanMessenger <- input
				}
			}

			// encrypt message content
			b, err := aes.Encrypt(global.Keyfile.GetKey(), string(out.Message.Content))
			if err != nil {
				fmt.Println(err)
				out.Error = typesMesh.ErrServerEncryption
			} else {
				out.Message.Content = []byte(b)
			}
		}

	default:
		global.NetworkConfig.ServerList.Delete(input.Receiver.ID)
	}

	*output = out

	return nil
}

func announce(input typesMesh.ServerMessage, output *typesMesh.ServerMessage) {

	// update network
	global.NetworkConfig.Network.Update(input.Network)

	global.NetworkConfig.ServerList.Clean(global.NetworkConfig.Network.GetMaxClientAge())

	// input
	// update server list
	var newList []serverInfo.ServerInfo

	err := json.Unmarshal(input.Message.Content, &newList)
	if err != nil {
		fmt.Println(err)
	} else {
		global.NetworkConfig.ServerList.AddList(newList)
	}

	// output
	// server list
	global.NetworkConfig.ServerList.Add(global.Config.Server.Get())

	b, err := json.Marshal(global.NetworkConfig.ServerList.Get())
	if err != nil {
		fmt.Println(err)
	} else {
		output.Message.Type = typesMesh.MessAnnounce
		output.Message.Content = b
	}
}
