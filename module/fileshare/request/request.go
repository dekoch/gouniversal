package request

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dekoch/gouniversal/module/fileshare/global"
	"github.com/dekoch/gouniversal/module/fileshare/typefileshare"
	"github.com/dekoch/gouniversal/shared/console"
)

func LoadConfig() {

	http.HandleFunc("/fileshare/req/", handleRequest)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()

	req := new(typefileshare.Request)

	req.Values = r.URL.Query()
	console.Output("GET params:", "")
	console.Output(req.Values, "")

	req.ID = req.Values.Get("id")
	req.Key = req.Values.Get("key")
	req.FilePath = req.Values.Get("file")

	resp := new(typefileshare.Response)
	resp.Status = http.StatusOK
	resp.Err = nil
	resp.Content = ""
	resp.FilePath = ""

	resp.FilePath = global.Config.FileRoot + req.FilePath

	var fi os.FileInfo

	for i := 0; i <= 1; i++ {
		if resp.Err == nil {
			switch i {
			case 0:
				fi, resp.Err = os.Stat(resp.FilePath)

			case 1:
				if fi.IsDir() {
					resp.Err = errors.New("path forbidden")
					resp.Status = http.StatusForbidden
				}

				/*case 2:
					if functions.IsEmpty(req.ID) {
						resp.Err = errors.New("ID not set")
						resp.Status = http.StatusForbidden
					}

				case 3:
					if functions.IsEmpty(req.Key) {
						resp.Err = errors.New("key not set")
						resp.Status = http.StatusForbidden
					}*/
			}
		}
	}

	if resp.Err != nil {
		if resp.Status == http.StatusOK {
			resp.Status = http.StatusInternalServerError
		}

		http.Error(w, resp.Err.Error(), resp.Status)
		console.Output(req.ID+"\t"+resp.Err.Error(), "")
	} else {

		http.ServeFile(w, r, resp.FilePath)
	}

	t := time.Now()
	elapsed := t.Sub(startTime)
	f := elapsed.Seconds() * 1000.0
	console.Output(strconv.FormatFloat(f, 'f', 1, 64)+"ms", "")
}
