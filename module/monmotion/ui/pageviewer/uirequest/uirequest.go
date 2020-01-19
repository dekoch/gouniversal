package uirequest

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/dekoch/gouniversal/module/monmotion/dbstorage"
	"github.com/dekoch/gouniversal/module/monmotion/mdimg"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/token"
)

type UIRequest struct {
	tokens token.Token
	active bool
}

type request struct {
	Values  url.Values
	UUID    string
	Token   string
	ImageID string
}

type response struct {
	Content []byte
	Status  int
	Err     error
}

var mut sync.RWMutex

func (rt *UIRequest) LoadConfig() error {

	mut.Lock()
	defer mut.Unlock()

	if rt.active {
		return nil
	}

	http.HandleFunc("/monmotion/viewer/", rt.serve)

	rt.active = true

	return nil
}

func (rt *UIRequest) GetNewToken(uid string) string {

	return rt.tokens.New(uid)
}

func (rt *UIRequest) serve(w http.ResponseWriter, r *http.Request) {

	var req request
	req.Values = r.URL.Query()
	req.UUID = req.Values.Get("uuid")
	req.Token = req.Values.Get("token")
	req.ImageID = req.Values.Get("imageid")

	var resp response
	resp.Status = http.StatusOK

	for i := 0; i <= 4; i++ {
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
				if rt.tokens.Check(req.UUID, req.Token) == false {
					resp.Err = errors.New("invalid token")
					resp.Status = http.StatusForbidden
				}

			case 3:
				if functions.IsEmpty(req.ImageID) {
					resp.Err = errors.New("imageid not set")
					resp.Status = http.StatusForbidden
				}

			case 4:
				var img mdimg.MDImage
				img, resp.Err = dbstorage.LoadImage(req.ImageID)
				resp.Content = img.Jpeg
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

		fmt.Fprintf(w, "%s", resp.Content)
	}
}
