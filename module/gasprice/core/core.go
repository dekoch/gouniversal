package core

import (
	"time"

	"github.com/dekoch/gouniversal/module/gasprice/csv"
	"github.com/dekoch/gouniversal/module/gasprice/finder"
	"github.com/dekoch/gouniversal/module/gasprice/global"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
)

var (
	chanCheckStart = make(chan bool)
)

func LoadConfig() {

	go checkPrice()
	go job()

	//chanCheckStart <- true
}

func Exit() {

}

func job() {

	intvl := global.Config.GetUpdInterval()
	timer := time.NewTimer(intvl)

	for {

		if intvl > 0 {

			select {
			case <-timer.C:
				chanCheckStart <- true
				timer.Reset(intvl)
			}
		} else {
			// wait until enabled
			time.Sleep(1 * time.Minute)
			intvl = global.Config.GetUpdInterval()
		}
	}
}

func checkPrice() {

	for {
		<-chanCheckStart

		fileName := time.Now().Format("2006")
		fileName += ".csv"

		for _, st := range global.Config.Stations.GetList() {

			if functions.IsEmpty(st.Name) ||
				functions.IsEmpty(st.URL) {

				continue
			}

			prices, err := finder.GetPrices(st)
			if err != nil {
				console.Log(err, "checkPrice()")
				continue
			}

			for _, price := range prices {

				csv.Export(global.Config.FileRoot+fileName, price)
			}
		}
	}
}
