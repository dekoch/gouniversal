package query

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dekoch/gouniversal/shared/api/instaclient"
)

func Send(id, username string, follow bool, ic *instaclient.InstaClient) ([]byte, error) {

	var (
		err       error
		followURL string
		userURL   string
		resp      *http.Response
		ret       []byte
	)

	for i := 0; i <= 4; i++ {

		switch i {
		case 0:
			if follow {
				followURL = "https://www.instagram.com/web/friendships/" + id + "/follow/"
			} else {
				followURL = "https://www.instagram.com/web/friendships/" + id + "/unfollow/"
			}

			_, err = url.Parse(followURL)

			fmt.Println(followURL)

		case 1:
			userURL = "https://www.instagram.com/" + username + "/"
			_, err = url.Parse(userURL)

		case 2:
			var p instaclient.Post
			p.URL = followURL
			p.Referer = userURL

			resp, err = ic.SendPost(p)

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
