package query

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dekoch/gouniversal/shared/api/instaclient"
)

func Send(username string, ic *instaclient.InstaClient) ([]byte, error) {

	var (
		err  error
		ur   string
		resp *http.Response
		ret  []byte
	)

	for i := 0; i <= 4; i++ {

		switch i {
		case 0:
			ur = "https://www.instagram.com/" + username + "/?utm_source=ig_seo&utm_campaign=profiles&utm_medium="
			//ur = "https://www.instagram.com/" + username + "/"

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
