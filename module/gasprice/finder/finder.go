package finder

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/dekoch/gouniversal/module/gasprice/finder/adac"
	"github.com/dekoch/gouniversal/module/gasprice/price"
	"github.com/dekoch/gouniversal/module/gasprice/station"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
)

func GetPrices(st station.Station) ([]price.Price, error) {

	var (
		err  error
		resp *http.Response
		page string
		ret  []price.Price
	)

	func() {

		for i := 0; i <= 5; i++ {

			switch i {
			case 0:
				// check input
				if functions.IsEmpty(st.URL) {
					err = errors.New("Please Enter Url")
				}

			case 1:
				_, err = url.Parse(st.URL)

			case 2:
				resp, err = http.Get(st.URL)
				if err == nil {
					defer resp.Body.Close()
				}

			case 3:
				ct := resp.Header.Get("Content-Type")

				if ct != "" {
					if strings.Contains(ct, "text/html") == false {
						err = errors.New("NotSupportedContentType: " + ct)
					}
				}

			case 4:
				var b []byte
				b, err = ioutil.ReadAll(resp.Body)
				if err == nil {
					page = string(b)
				}

			case 5:
				err = findOnPage(st, page, &ret)
			}

			if err != nil {
				console.Log(err, "")
				return
			}
		}
	}()

	return ret, err
}

func findOnPage(st station.Station, raw string, prices *[]price.Price) error {

	if strings.HasPrefix(st.URL, "https://www.adac.de/") {

		return adac.FindOnPage(st, raw, prices)
	}

	return errors.New("page not supported")
}
