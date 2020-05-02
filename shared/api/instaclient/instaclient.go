package instaclient

import (
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"
)

type InstaClient struct {
	cookies     []*http.Cookie
	lastRequest time.Time
}

type Cookies struct {
	DsUserID  string
	CsrfToken string
	SessionID string
}

type Post struct {
	URL     string
	Referer string
}

type Get struct {
	URL string
}

var mut sync.Mutex

func (ic *InstaClient) SetCookies(src Cookies) error {

	mut.Lock()
	defer mut.Unlock()

	if len(ic.cookies) > 0 {

		for i := range ic.cookies {

			switch ic.cookies[i].Name {
			case "ds_user_id":
				ic.cookies[i].Value = src.DsUserID

			case "csrftoken":
				ic.cookies[i].Value = src.CsrfToken

			case "sessionid":
				ic.cookies[i].Value = src.SessionID
			}

			ic.cookies[i].Path = "/"
			ic.cookies[i].Domain = ".instagram.com"
		}

		return nil
	}

	cookie := &http.Cookie{
		Name:   "ds_user_id",
		Value:  src.DsUserID,
		Path:   "/",
		Domain: ".instagram.com",
	}
	ic.cookies = append(ic.cookies, cookie)

	cookie = &http.Cookie{
		Name:   "csrftoken",
		Value:  src.CsrfToken,
		Path:   "/",
		Domain: ".instagram.com",
	}
	ic.cookies = append(ic.cookies, cookie)

	cookie = &http.Cookie{
		Name:   "sessionid",
		Value:  src.SessionID,
		Path:   "/",
		Domain: ".instagram.com",
	}
	ic.cookies = append(ic.cookies, cookie)

	return nil
}

func (ic *InstaClient) GetCookies() (Cookies, error) {

	mut.Lock()
	defer mut.Unlock()

	return ic.getCookies()
}

func (ic *InstaClient) getCookies() (Cookies, error) {

	var (
		err error
		ret Cookies
	)

	for i := range ic.cookies {

		switch ic.cookies[i].Name {
		case "ds_user_id":
			ret.DsUserID = ic.cookies[i].Value

		case "csrftoken":
			ret.CsrfToken = ic.cookies[i].Value

		case "sessionid":
			ret.SessionID = ic.cookies[i].Value
		}
	}

	return ret, err
}

func (ic *InstaClient) SendPost(p Post) (*http.Response, error) {

	mut.Lock()
	defer mut.Unlock()

	ic.delayClient()

	var ret *http.Response

	client, err := ic.newClient()
	if err != nil {
		return ret, err
	}

	req, err := http.NewRequest("POST", p.URL, nil)
	if err != nil {
		return ret, err
	}

	co, err := ic.getCookies()
	if err != nil {
		return ret, err
	}

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Length", "0")
	req.Header.Add("Host", "www.instagram.com")
	req.Header.Add("Referer", p.Referer)
	req.Header.Add("X-Instagram-AJAX", "f875f4c886d7")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("X-CSRFToken", co.CsrfToken)

	ret, err = client.Do(req)
	if err != nil {
		return ret, err
	}

	ic.cookies = req.Cookies()

	return ret, nil
}

func (ic *InstaClient) SendGet(g Get) (*http.Response, error) {

	mut.Lock()
	defer mut.Unlock()

	ic.delayClient()

	var ret *http.Response

	client, err := ic.newClient()
	if err != nil {
		return ret, err
	}

	req, err := http.NewRequest("GET", g.URL, nil)
	if err != nil {
		return ret, err
	}

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Host", "www.instagram.com")

	ret, err = client.Do(req)
	if err != nil {
		return ret, err
	}

	ic.cookies = req.Cookies()

	return ret, nil
}

func (ic *InstaClient) newClient() (*http.Client, error) {

	var client *http.Client

	u, err := url.Parse("https://www.instagram.com")
	if err != nil {
		return client, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return client, err
	}

	client = &http.Client{
		Jar: jar,
	}

	client.Jar.SetCookies(u, ic.cookies)

	return client, nil
}

func (ic *InstaClient) delayClient() {

	if time.Since(ic.lastRequest) < time.Duration(300*time.Millisecond) {

		delay := rand.Intn(500) + 300
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	ic.lastRequest = time.Now()
}
