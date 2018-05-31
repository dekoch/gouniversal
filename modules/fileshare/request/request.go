package request

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dekoch/gouniversal/modules/fileshare/global"
	"github.com/dekoch/gouniversal/modules/fileshare/typesFileshare"
	"github.com/dekoch/gouniversal/shared/console"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()

	req := new(typesFileshare.Request)

	req.Values = r.URL.Query()
	console.Output("GET params:", "")
	console.Output(req.Values, "")

	req.ID = req.Values.Get("id")
	req.Key = req.Values.Get("key")
	req.FilePath = req.Values.Get("file")

	resp := new(typesFileshare.Response)
	resp.Status = http.StatusOK
	resp.Err = nil
	resp.Content = ""
	resp.FilePath = ""

	resp.FilePath = global.Config.File.FileRoot + req.FilePath

	/*for i := 0; i <= 6; i++ {
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

			}
		}
	}*/

	if resp.Err != nil {
		if resp.Status == http.StatusOK {
			resp.Status = http.StatusInternalServerError
		}

		http.Error(w, resp.Err.Error(), resp.Status)
		console.Output(req.ID+"\t"+resp.Err.Error(), "")
	} else {

		if _, err := os.Stat(resp.FilePath); os.IsNotExist(err) == false {
			http.ServeFile(w, r, resp.FilePath)
		}
	}

	t := time.Now()
	elapsed := t.Sub(startTime)
	f := elapsed.Seconds() * 1000.0
	console.Output(strconv.FormatFloat(f, 'f', 1, 64)+"ms", "")
}

func LoadConfig() {

	http.HandleFunc("/fileshare/req/", handleRequest)
}
