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

	go run()
	go job()
}

func job() {

	intvl := global.Config.GetUpdInterval()
	timer := time.NewTimer(intvl)

	for {

		if intvl != 0 {

			select {
			case <-timer.C:
				timer.Stop()
				chanStart <- true

			case <-chanFinish:
				intvl = global.Config.GetUpdInterval()
				timer.Reset(intvl)
			}
		} else {
			// wait until enabled
			time.Sleep(1 * time.Minute)
			intvl = global.Config.GetUpdInterval()
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
