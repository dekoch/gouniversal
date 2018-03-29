package request

import (
	"fmt"
	"gouniversal/modules/openespm/deviceManagement"
	"gouniversal/modules/openespm/oespmTypes"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {

	req := new(oespmTypes.Request)

	req.Values = r.URL.Query()
	fmt.Println("GET params:", req.Values)

	req.UUID = req.Values.Get("id")
	req.Key = req.Values.Get("key")
	fmt.Println(req.UUID)
	fmt.Println(req.Key)

	dev := deviceManagement.SelectDevice(req.UUID)

	fmt.Println(dev.State)
}

func LoadConfig() {

	http.HandleFunc("/request/", handleRequest)
}

func Exit() {

}
