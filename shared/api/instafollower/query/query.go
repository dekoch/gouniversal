package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dekoch/gouniversal/shared/api/instaclient"
)

type variables struct {
	ID          string `json:"id"`
	IncludeReel bool   `json:"include_reel"`
	FetchMutual bool   `json:"fetch_mutual"`
	First       int    `json:"first"`
	After       string `json:"after"`
}

func (iv *variables) marshal() (string, error) {

	b, err := json.Marshal(iv)

	return string(b), err
}

func getURL(queryhash string, v variables) (string, error) {

	va, err := v.marshal()
	if err != nil {
		return "", err
	}

	ret := "https://www.instagram.com/graphql/query/?"
	ret += "query_hash=" + queryhash
	ret += "&variables=" + va

	return ret, nil
}

func Send(id, queryhash string, first int, after string, ic *instaclient.InstaClient) ([]byte, error) {

	var (
		err  error
		ur   string
		resp *http.Response
		ret  []byte
	)

	for i := 0; i <= 4; i++ {

		switch i {
		case 0:
			var va variables
			va.ID = id
			va.IncludeReel = false
			va.FetchMutual = false
			va.First = first
			va.After = after

			ur, err = getURL(queryhash, va)

			fmt.Println(ur)

		case 1:
			_, err = url.Parse(ur)

		case 2:
			var g instaclient.Get
			g.URL = ur

			resp, err = ic.SendGet(g)

		case 3:
			if resp != nil {
				defer resp.Body.Close()
			} else {
				err = errors.New("resp nil")
			}

		case 4:
			ret, err = ioutil.ReadAll(resp.Body)
		}

		if err != nil {
			return ret, err
		}
	}

	return ret, err
}
