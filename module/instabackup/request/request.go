package request

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/dekoch/gouniversal/module/instabackup/global"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
)

type request struct {
	Values   url.Values
	UUID     string
	Token    string
	FilePath string
}

type response struct {
	Content  string
	FilePath string
	Status   int
	Err      error
}

func LoadConfig() {

	http.HandleFunc("/instabackup/req/", handleRequest)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()

	var req request
	req.Values = r.URL.Query()
	req.UUID = req.Values.Get("uuid")
	req.Token = req.Values.Get("token")
	req.FilePath = req.Values.Get("file")

	var resp response
	resp.Status = http.StatusOK
	resp.Err = nil
	resp.Content = ""
	resp.FilePath = global.Config.FileRoot + req.FilePath

	var fi os.FileInfo

	for i := 0; i <= 4; i++ {
		if resp.Err == nil {
			switch i {
			case 0:
				fi, resp.Err = os.Stat(resp.FilePath)

			case 1:
				if fi.IsDir() {
					resp.Err = errors.New("path forbidden")
					resp.Status = http.StatusForbidden
				}

			case 2:
				if functions.IsEmpty(req.UUID) {
					resp.Err = errors.New("UUID not set")
					resp.Status = http.StatusForbidden
				}

			case 3:
				if functions.IsEmpty(req.Token) {
					resp.Err = errors.New("token not set")
					resp.Status = http.StatusForbidden
				}

			case 4:
				if global.Tokens.Check(req.UUID, req.Token) == false {
					resp.Err = errors.New("invalid token")
					resp.Status = http.StatusForbidden
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

		http.ServeFile(w, r, resp.FilePath)
	}

	elapsed := time.Now().Sub(startTime)
	f := elapsed.Seconds() * 1000.0
	console.Output(strconv.FormatFloat(f, 'f', 1, 64)+"ms", "")
}
