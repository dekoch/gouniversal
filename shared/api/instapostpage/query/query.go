package query

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dekoch/gouniversal/shared/api/instaclient"
)

func Send(id string, ic *instaclient.InstaClient) ([]byte, error) {

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
				ur = "https://www.instagram.com/p/" + id

				fmt.Println(ur)

				_, err = url.Parse(ur)

			case 1:
				var g instaclient.Get
				g.URL = ur
				g.SetCookies = false

				resp, err = ic.SendGet(g)

			case 2:
				if resp != nil {
					defer resp.Body.Close()
				} else {
					err = errors.New("resp nil")
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
