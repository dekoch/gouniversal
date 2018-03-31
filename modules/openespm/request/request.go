package request

import (
	"errors"
	"fmt"
	"gouniversal/modules/openespm/app/appManagement"
	"gouniversal/modules/openespm/deviceManagement"
	"gouniversal/modules/openespm/oespmTypes"
	"gouniversal/shared/functions"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {

	req := new(oespmTypes.Request)

	req.Values = r.URL.Query()
	fmt.Println("GET params:", req.Values)

	req.UUID = req.Values.Get("id")
	req.Key = req.Values.Get("key")

	resp := new(oespmTypes.Response)
	resp.Type = oespmTypes.PLAIN
	resp.Content = ""
	resp.Status = http.StatusOK
	resp.Err = nil

	for i := 0; i <= 4; i++ {
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
				req.Device = deviceManagement.SelectDevice(req.UUID)
				if req.Device.State < 0 {
					resp.Err = errors.New("device not found")
					resp.Status = http.StatusForbidden
				}

			case 3:
				if req.Device.Key != req.Key {
					resp.Err = errors.New("key mismatch")
					resp.Status = http.StatusForbidden
				}

			case 4:
				appManagement.Request(resp, req)
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
		case oespmTypes.JSON:
			w.Header().Set("Content-Type", "application/json")

		case oespmTypes.XML:
			w.Header().Set("Content-Type", "application/xml")
		}

		w.Write([]byte(resp.Content))
	}
}

func LoadConfig() {

	http.HandleFunc("/request/", handleRequest)
}

func Exit() {

}
