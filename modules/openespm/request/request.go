package request

import (
	"errors"
	"fmt"
	"gouniversal/modules/openespm/app"
	"gouniversal/modules/openespm/globalOESPM"
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/shared/functions"
	"net/http"
	"strconv"
	"time"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()

	req := new(typesOESPM.Request)

	req.Values = r.URL.Query()
	fmt.Println("GET params:", req.Values)

	req.ID = req.Values.Get("id")
	req.Key = req.Values.Get("key")

	resp := new(typesOESPM.Response)
	resp.Type = typesOESPM.PLAIN
	resp.Content = ""
	resp.Status = http.StatusOK
	resp.Err = nil

	for i := 0; i <= 6; i++ {
		if resp.Err == nil {
			switch i {
			case 0:
				if functions.IsEmpty(req.ID) {
					resp.Err = errors.New("ID not set")
					resp.Status = http.StatusForbidden
				}

			case 1:
				if functions.IsEmpty(req.Key) {
					resp.Err = errors.New("key not set")
					resp.Status = http.StatusForbidden
				}

			case 2:
				req.Device, resp.Err = globalOESPM.DeviceConfig.GetWithReqID(req.ID)
				if resp.Err != nil {
					resp.Status = http.StatusForbidden
				}

			case 3:
				if req.Device.RequestKey != req.Key {
					resp.Err = errors.New("key mismatch")
					resp.Status = http.StatusForbidden
				}

			case 4:
				if req.Device.State == 2 {
					resp.Err = errors.New("device is inactive")
					resp.Status = http.StatusForbidden
				}

			case 5:
				req.DeviceDataFolder = globalOESPM.DeviceDataFolder + req.Device.UUID + "/"
				app.Request(resp, req)

			case 6:
				resp.Err = globalOESPM.DeviceConfig.Edit(req.Device.UUID, req.Device)
			}
		}
	}

	if resp.Err != nil {
		if resp.Status == http.StatusOK {
			resp.Status = http.StatusInternalServerError
		}

		http.Error(w, resp.Err.Error(), resp.Status)
		fmt.Println(req.ID + "\t" + resp.Err.Error())
	} else {

		switch resp.Type {
		case typesOESPM.JSON:
			w.Header().Set("Content-Type", "application/json")

		case typesOESPM.XML:
			w.Header().Set("Content-Type", "application/xml")
		}

		w.Write([]byte(resp.Content))
	}

	t := time.Now()
	elapsed := t.Sub(startTime)
	f := elapsed.Seconds() * 1000.0
	fmt.Println(strconv.FormatFloat(f, 'f', 1, 64) + "ms")
}

func LoadConfig() {

	http.HandleFunc("/oespmreq/", handleRequest)
}

func Exit() {
	// save device config on exit
	globalOESPM.DeviceConfig.SaveConfig()
}
