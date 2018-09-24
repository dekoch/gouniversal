package messenger

import (
	"fmt"
	"time"

	"github.com/dekoch/gouniversal/modules/mesh"
	"github.com/dekoch/gouniversal/modules/mesh/typesMesh"
	"github.com/dekoch/gouniversal/modules/messenger/global"
	"github.com/dekoch/gouniversal/shared/io/file"
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
			//timerHello.Reset(1 * time.Second)
		}
	}
}

func hello() {

	for {
		<-cHelloStart

		func() {

			var message typesMesh.ServerMessage
			var err error

			message.Receiver, err = mesh.GetServerWithID("9c3fc567-c4c7-43b9-a239-87de716d2d86")
			if err != nil {
				fmt.Println(err)
				return
			}

			message.Message.Type = typesMesh.MessMessenger
			message.Message.Version = 1.0

			b, err := file.ReadFile("test.tar.gz")
			if err != nil {
				fmt.Println(err)
				return
			}

			message.Message.Content = b

			t := time.Now()

			mesh.SendMessage(message)

			fmt.Print("sent ")
			fmt.Println(time.Since(t).Seconds())

		}()

		cHelloFinish <- true
	}
}

func Exit() {

	global.Config.SaveConfig()
}
