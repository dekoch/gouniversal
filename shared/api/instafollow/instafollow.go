package instafollow

import (
	"encoding/json"
	"errors"

	"github.com/dekoch/gouniversal/shared/api/instaclient"
	"github.com/dekoch/gouniversal/shared/api/instafollow/query"
)

type response struct {
	Result string `json:"result"`
	Status string `json:"status"`
}

func Follow(id, username string, ic *instaclient.InstaClient) error {

	return send(id, username, true, ic)
}

func Unfollow(id, username string, ic *instaclient.InstaClient) error {

	return send(id, username, false, ic)
}

func send(id, username string, follow bool, ic *instaclient.InstaClient) error {

	var (
		err  error
		b    []byte
		resp response
	)

	for i := 0; i <= 2; i++ {

		switch i {
		case 0:
			b, err = query.Send(id, username, follow, ic)

		case 1:
			err = json.Unmarshal(b, &resp)
			if err != nil {
				err = errors.New(err.Error() + " \"" + string(b) + "\"")
			}

		case 2:
			if resp.Status != "ok" {
				err = errors.New(resp.Status)
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}
