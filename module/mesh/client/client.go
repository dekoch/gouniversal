package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"strconv"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/mesh/global"
	"github.com/dekoch/gouniversal/module/mesh/serverInfo"
	"github.com/dekoch/gouniversal/module/mesh/settings"
	"github.com/dekoch/gouniversal/module/mesh/typesMesh"
	"github.com/dekoch/gouniversal/shared/aes"
	"github.com/dekoch/gouniversal/shared/console"
)

const debug = false

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

		serverList := global.NetworkConfig.ServerList.Get()
		var message typesMesh.ServerMessage
		message.Message.Type = typesMesh.MessAnnounce
		message.Message.Version = 1.0

		b, err := json.Marshal(serverList)
		if err != nil {
			fmt.Println(err)
		} else {

			message.Message.Content = b

			var wg sync.WaitGroup

			// announce to every server in list
			for _, server := range serverList {

				// only to other systems
				if IsLoop(server) {
					continue
				}

				if debug {
					console.Output("announce to "+server.ID, "mesh")
				}

				wg.Add(1)

				message.Receiver = server

				go func(msg typesMesh.ServerMessage) {

					SendMessage(msg)

					wg.Done()
				}(message)
			}

			wg.Wait()
		}

		global.NetworkConfig.ServerList.Clean()

		chanAnnounceFinish <- true
	}
}

func hello() {

	for {
		<-chanHelloStart

		global.Config.Server.Update()

		var message typesMesh.ServerMessage
		message.Message.Type = typesMesh.MessHello

		serverList := global.NetworkConfig.ServerList.Get()

		var wg sync.WaitGroup

		// hello to every server in list
		for _, server := range serverList {

			// only to other systems
			if IsLoop(server) {
				continue
			}

			if debug {
				console.Output("hello to "+server.ID, "mesh")
			}

			wg.Add(1)

			message.Receiver = server

			go func(msg typesMesh.ServerMessage) {

				SendMessage(msg)

				wg.Done()
			}(message)
		}

		wg.Wait()

		chanHelloFinish <- true
	}
}

func SendMessage(output typesMesh.ServerMessage) error {

	// send only to addresses from other systems
	if IsLoop(output.Receiver) {
		return errors.New("IsLoopback")
	}

	output.Sender = global.Config.Server.Get()
	output.Network = global.NetworkConfig.Network.Get()

	// encrypt message content
	b, err := aes.Encrypt(global.Keyfile.GetKey(), string(output.Message.Content))
	if err != nil {
		return err
	}

	output.Message.Content = []byte(b)

	// try preferred address
	prefAddr := global.NetworkConfig.ServerList.GetPrefAddress(output.Receiver)

	if prefAddr != "" {

		err = send(output, prefAddr)
		if err == nil {
			return err
		}
	}

	// try all addresses from server
	for _, addr := range output.Receiver.Address {

		err = send(output, addr)
		if err == nil {
			global.NetworkConfig.ServerList.SetPrefAddress(output.Receiver, addr)
			return err
		}
	}

	return errors.New("not send")
}

func send(output typesMesh.ServerMessage, addr string) error {

	// check for v4 or v6 addresses
	ip := net.ParseIP(addr)
	if ip == nil {
		return errors.New("bad address")
	}

	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return errors.New("bad address")
	}

	address := ""

	if ip.To4() != nil {
		address = addr
	} else if ip.To16() != nil {
		// we have to add []
		address = "[" + addr + "]"
	} else {
		return errors.New("bad address")
	}

	receiverPort := strconv.Itoa(output.Receiver.Port)

	conn, err := net.DialTimeout("tcp", address+":"+receiverPort, 5*time.Second)
	if err != nil {
		return err
	}

	c := rpc.NewClient(conn)

	var inputErr string
	err = c.Call("Server.Message", output, &inputErr)
	if err != nil {
		return err
	}

	c.Close()

	// message sent
	return nil
}

func IsLoop(in serverInfo.ServerInfo) bool {

	if settings.LocalConnection {
		return false
	}

	this := global.Config.Server.Get()

	if this.ID == in.ID {
		return true
	}

	for _, thisAddr := range this.Address {

		for _, inAddr := range in.Address {

			if thisAddr+strconv.Itoa(this.Port) == inAddr+strconv.Itoa(in.Port) {

				return true
			}
		}
	}

	return false
}
