package tracker

import (
	"time"

	"github.com/dekoch/gouniversal/module/iptracker/global"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/getpublicip"
)

var (
	chanStart  = make(chan bool)
	chanFinish = make(chan bool)
	lastIP     string
)

func LoadConfig() {

	if global.Config.GetUpdInterval() != 0 {

		go run()
		go job()
	}
}

func job() {

	timer := time.NewTimer(global.Config.GetUpdInterval())

	for {
		select {
		case <-timer.C:
			timer.Stop()
			chanStart <- true

		case <-chanFinish:
			timer.Reset(global.Config.GetUpdInterval())
		}
	}
}

func run() {

	for {
		<-chanStart

		ip, err := getpublicip.Get()
		if err != nil {
			console.Log(err, "")
		} else {

			if ip != lastIP {

				lastIP = ip

				console.Log(ip, "IPTracker")
			}
		}

		chanFinish <- true
	}
}
