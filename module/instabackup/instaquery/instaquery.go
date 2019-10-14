package instaquery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type InstaQuery struct {
	QueryHash string
	Variables Variables
}

type Variables struct {
	ID    string `json:"id"`
	First uint   `json:"first"`
	After string `json:"after"`
}

func (iv *Variables) marshal() (string, error) {

	b, err := json.Marshal(iv)

	return string(b), err
}

func (iq *InstaQuery) getQuery() (string, error) {

	v, err := iq.Variables.marshal()
	if err != nil {
		return "", err
	}

	ret := "https://www.instagram.com/graphql/query/?"
	ret += "query_hash=" + iq.QueryHash
	ret += "&variables=" + v

	return ret, nil
}

func (iq *InstaQuery) SendQuery() ([]byte, error) {

	var (
		err  error
		ur   string
		resp *http.Response
		ret  []byte
	)

	func() {

		for i := 0; i <= 3; i++ {

			switch i {
			case 0:
				ur, err = iq.getQuery()

				fmt.Println(ur)

			case 1:
				_, err = url.Parse(ur)

			case 2:
				resp, err = http.Get(ur)
				if err == nil {
					defer resp.Body.Close()
				}

			case 3:
				ret, err = ioutil.ReadAll(resp.Body)
			}

			if err != nil {
				return
			}
		}
	}()

	return ret, err
}
