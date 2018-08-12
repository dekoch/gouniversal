package client

import (
	"encoding/json"
	"fmt"
	"net"
	"net/rpc"
	"strconv"
	"time"

	"github.com/dekoch/gouniversal/modules/mesh/global"
	"github.com/dekoch/gouniversal/modules/mesh/serverInfo"
	"github.com/dekoch/gouniversal/modules/mesh/typesMesh"
	"github.com/dekoch/gouniversal/shared/aes"
)

const localConnection = true

var (
	chanAnnounceStart  = make(chan bool)
	chanAnnounceFinish = make(chan bool)
	chanHelloStart     = make(chan bool)
	chanHelloFinish    = make(chan bool)
)

func LoadConfig() {

	go announce()
	go hello()
	go job()
}

func job() {

	timerAnnounce := time.NewTimer(global.NetworkConfig.Network.GetAnnounceInterval())
	timerHello := time.NewTimer(global.NetworkConfig.Network.GetHelloInterval())

	for {
		select {
		case <-timerAnnounce.C:
			timerAnnounce.Stop()
			chanAnnounceStart <- true

		case <-chanAnnounceFinish:
			timerAnnounce.Reset(global.NetworkConfig.Network.GetAnnounceInterval())

		case <-timerHello.C:
			timerHello.Stop()
			chanHelloStart <- true

		case <-chanHelloFinish:
			timerHello.Reset(global.NetworkConfig.Network.GetHelloInterval())
		}
	}
}

func announce() {

	for {
		<-chanAnnounceStart

		global.Config.Server.Update()
		global.NetworkConfig.ServerList.Add(global.Config.Server.Get())
		global.NetworkConfig.ServerList.Clean(global.NetworkConfig.Network.GetMaxClientAge())

		serverList := global.NetworkConfig.ServerList.Get()
		client := global.Config.Server.Get()
		var message typesMesh.ServerMessage
		message.Message.Type = typesMesh.MessAnnounce
		message.Message.Version = 1.0

		// announce to every server in list
		for _, server := range serverList {

			if localConnection == false {
				// only to other systems
				if server.ID == client.ID {
					continue
				}
			}

			message.Receiver = server

			b, err := json.Marshal(serverList)
			if err != nil {
				fmt.Println(err)
			} else {

				message.Message.Content = b

				input := SendMessage(message)
				if input.Error == typesMesh.ErrNil {

					global.NetworkConfig.Network.Update(input.Network)

					if input.Message.Type == typesMesh.MessAnnounce {

						var newList []serverInfo.ServerInfo

						err := json.Unmarshal(input.Message.Content, &newList)
						if err != nil {
							fmt.Println(err)
						} else {
							global.NetworkConfig.ServerList.AddList(newList)
						}
					}
				}
			}
		}

		chanAnnounceFinish <- true
	}
}

func hello() {

	for {
		<-chanHelloStart

		var message typesMesh.ServerMessage
		message.Message.Type = typesMesh.MessHello

		serverList := global.NetworkConfig.ServerList.Get()
		client := global.Config.Server.Get()

		// hello to every server in list
		for _, server := range serverList {

			if localConnection == false {
				// only to other systems
				if server.ID == client.ID {
					continue
				}
			}

			message.Receiver = server

			SendMessage(message)
		}

		chanHelloFinish <- true
	}
}

func SendMessage(output typesMesh.ServerMessage) typesMesh.ServerMessage {

	output.Sender = global.Config.Server.Get()
	output.Network = global.NetworkConfig.Network.Get()
	output.Error = typesMesh.ErrNil

	var input typesMesh.ServerMessage
	input.Error = typesMesh.ErrNoConnection

	// encrypt message content
	b, err := aes.Encrypt(global.Keyfile.GetKey(), string(output.Message.Content))
	if err != nil {
		fmt.Println(err)
		input.Error = typesMesh.ErrClientEncryption
		return input
	}

	output.Message.Content = []byte(b)

	serverOK := true

	// try all addresses from server
	for _, addr := range output.Receiver.Address {

		if serverOK {

			addressOK := true
			address := ""

			if localConnection == false {
				// send only to addresses from other systems
				for _, senderAddr := range output.Sender.Address {
					if senderAddr == addr {
						addressOK = false
					}
				}
			}

			// check for v4 or v6 addresses
			ip := net.ParseIP(addr)

			if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
				addressOK = false
			}

			if ip.To4() != nil {
				address = addr
			} else if ip.To16() != nil {
				// we have to add []
				address = "[" + addr + "]"
			} else {
				addressOK = false
			}

			if addressOK {
				// add port
				address += ":" + strconv.Itoa(global.NetworkConfig.Network.GetPort())

				switch output.Message.Type {
				case typesMesh.MessAnnounce:
					fmt.Print("announce")

				case typesMesh.MessHello:
					fmt.Print("hello")

				case typesMesh.MessRAW:
					fmt.Print("raw")
				}

				fmt.Print(" to \"" + output.Receiver.ID + "\" @" + address + "...")

				c, err := rpc.Dial("tcp", address)
				if err != nil {
					fmt.Println(err)
				} else {

					err = c.Call("Server.Message", output, &input)
					if err != nil {
						fmt.Println(err)
					} else {
						// error response from server
						switch input.Error {
						case typesMesh.ErrNil:

							// validate server response
							if global.NetworkConfig.Network.CheckID(input.Network.ID) == false {
								input.Error = typesMesh.ErrClientDifferentMeshID
							} else if global.NetworkConfig.Network.CheckHashWithLocalKey(input.Network.Hash) == false {
								input.Error = typesMesh.ErrClientWrongMeshKey
							} else if output.Receiver.ID != input.Sender.ID {
								input.Error = typesMesh.ErrClientWrongSender
							}

							switch input.Error {
							case typesMesh.ErrNil:
								// decrypt message content
								b, err := aes.Decrypt(global.Keyfile.GetKey(), string(input.Message.Content))
								if err != nil {
									fmt.Println(err)
									input.Error = typesMesh.ErrClientDecryption
								} else {
									input.Message.Content = []byte(b)

									fmt.Println("OK")
								}

							default:
								writeError(input.Error)

								global.NetworkConfig.ServerList.Delete(output.Receiver.ID)
							}

						default:
							writeError(input.Error)

							global.NetworkConfig.ServerList.Delete(output.Receiver.ID)
						}

						c.Close()

						// message sent
						serverOK = false
					}
				}
			}
		}
	}

	return input
}

func writeError(serr typesMesh.MessageType) {

	switch serr {
	case typesMesh.ErrServerDifferentMeshID:
		fmt.Println("server: different mesh ID")

	case typesMesh.ErrServerWrongMeshKey:
		fmt.Println("server: wrong mesh key")

	case typesMesh.ErrServerWrongReceiver:
		fmt.Println("server: wrong receiver")

	case typesMesh.ErrClientDifferentMeshID:
		fmt.Println("different mesh ID from server")

	case typesMesh.ErrClientWrongMeshKey:
		fmt.Println("wrong mesh key from server")

	case typesMesh.ErrClientWrongSender:
		fmt.Println("different server ID from server")
	}
}
