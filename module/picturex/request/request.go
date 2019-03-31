package request

import (
	"errors"
	"net/http"
	"net/url"
	"os"

	"github.com/dekoch/gouniversal/module/picturex/global"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
)

type request struct {
	Values url.Values
	UUID   string
	Token  string
	Name   string
}

type response struct {
	Content string
	Path    string
	Status  int
	Err     error
}

func LoadConfig() {

	http.HandleFunc("/picturex/req/", handleRequest)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	req := new(request)
	req.Values = r.URL.Query()
	//console.Output("GET params:", "")
	//console.Output(req.Values, "")

	req.UUID = req.Values.Get("uuid")
	req.Token = req.Values.Get("token")
	req.Name = req.Values.Get("name")

	resp := new(response)
	resp.Status = http.StatusOK
	resp.Err = nil
	resp.Content = ""
	resp.Path = global.Config.TempFileRoot + req.Name

	for i := 0; i <= 3; i++ {
		if resp.Err == nil {
			switch i {
			case 0:
				if functions.IsEmpty(req.UUID) {
					resp.Err = errors.New("UUID not set")
					resp.Status = http.StatusForbidden
				}

			case 1:
				if functions.IsEmpty(req.Token) {
					resp.Err = errors.New("token not set")
					resp.Status = http.StatusForbidden
				}

			case 2:
				if global.Tokens.Check(req.UUID, req.Token) == false {
					resp.Err = errors.New("invalid token")
					resp.Status = http.StatusForbidden
				}

			case 3:
				if _, err := os.Stat(resp.Path); os.IsNotExist(err) {
					resp.Err = errors.New("file not found")
					resp.Status = http.StatusNotFound
				}
			}
		}
	}

	if resp.Err != nil {
		if resp.Status == http.StatusOK {
			resp.Status = http.StatusInternalServerError
		}

		http.Error(w, resp.Err.Error(), resp.Status)
		console.Output(req.UUID+"\t"+resp.Err.Error(), "")
	} else {
		http.ServeFile(w, r, resp.Path)
	}
}
