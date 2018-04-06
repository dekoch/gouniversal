package request

import (
	"errors"
	"fmt"
	"gouniversal/modules/openespm/app/appManagement"
	"gouniversal/modules/openespm/deviceManagement"
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

	req.UUID = req.Values.Get("id")
	req.Key = req.Values.Get("key")

	resp := new(typesOESPM.Response)
	resp.Type = typesOESPM.PLAIN
	resp.Content = ""
	resp.Status = http.StatusOK
	resp.Err = nil

	for i := 0; i <= 5; i++ {
		if resp.Err == nil {
			switch i {
			case 0:
				if functions.IsEmpty(req.UUID) {
					resp.Err = errors.New("UUID not set")
					resp.Status = http.StatusForbidden
				}

			case 1:
				if functions.IsEmpty(req.Key) {
					resp.Err = errors.New("key not set")
					resp.Status = http.StatusForbidden
				}

			case 2:
				req.Device, resp.Err = deviceManagement.LoadDevice(req.UUID)
				if resp.Err != nil {
					resp.Status = http.StatusForbidden
				}

			case 3:
				if req.Device.Key != req.Key {
					resp.Err = errors.New("key mismatch")
					resp.Status = http.StatusForbidden
				}

			case 4:
				appManagement.Request(resp, req)

			case 5:
				resp.Err = deviceManagement.SaveDevice(req.UUID, req.Device)
			}
		}
	}

	if resp.Err != nil {
		if resp.Status == http.StatusOK {
			resp.Status = http.StatusInternalServerError
		}

		http.Error(w, resp.Err.Error(), resp.Status)
		fmt.Println(req.UUID + "\t" + resp.Err.Error())
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

	http.HandleFunc("/request/", handleRequest)
}

func Exit() {
	// save device config on exit
	globalOESPM.DeviceConfig.Mut.Lock()
	deviceManagement.SaveConfig(globalOESPM.DeviceConfig.File)
	globalOESPM.DeviceConfig.Mut.Unlock()
}
