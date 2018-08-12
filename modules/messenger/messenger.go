package messenger

import (
	"fmt"
	"time"

	"github.com/dekoch/gouniversal/modules/mesh"
	"github.com/dekoch/gouniversal/modules/mesh/typesMesh"
	"github.com/dekoch/gouniversal/modules/messenger/global"

	meshGlobal "github.com/dekoch/gouniversal/modules/mesh/global"
)

var (
	cHelloStart  = make(chan bool)
	cHelloFinish = make(chan bool)
)

func LoadConfig() {

	global.Config.LoadConfig()

	//go hello()
	//go job()
}

func job() {

	timerHello := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timerHello.C:
			timerHello.Stop()
			cHelloStart <- true

		case <-cHelloFinish:
			timerHello.Reset(1 * time.Second)

		case input := <-meshGlobal.ChanMessenger:
			fmt.Println("from C")
			fmt.Println(string(input.Message.Content))
		}
	}
}

func hello() {

	for {
		<-cHelloStart

		var message typesMesh.ServerMessage
		message.Receiver = mesh.GetServerInfo()
		message.Message.Type = typesMesh.MessMessenger
		message.Message.Version = 1.0
		message.Message.Content = []byte("huhu vom messenger")

		mesh.SendMessage(message)

		cHelloFinish <- true
	}
}

func Exit() {

	global.Config.SaveConfig()
}
