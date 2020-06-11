package instaclient

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/shared/functions"
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
	URL            string
	Referer        string
	SetCookies     bool
	XInstagramAJAX string
}

type Get struct {
	URL        string
	SetCookies bool
}

var mut sync.Mutex

func (ic *InstaClient) SetCookies(src Cookies) error {

	if functions.IsEmpty(src.DsUserID) {
		return errors.New("DsUserID not set")
	}

	if functions.IsEmpty(src.CsrfToken) {
		return errors.New("CsrfToken not set")
	}

	if functions.IsEmpty(src.SessionID) {
		return errors.New("SessionID not set")
	}

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

	if functions.IsEmpty(ret.DsUserID) {
		return ret, errors.New("DsUserID not set")
	}

	if functions.IsEmpty(ret.CsrfToken) {
		return ret, errors.New("CsrfToken not set")
	}

	if functions.IsEmpty(ret.SessionID) {
		return ret, errors.New("SessionID not set")
	}

	return ret, err
}

func (ic *InstaClient) AddCookies(cookies []*http.Cookie) error {

	for i := range cookies {

		err := ic.AddCookie(cookies[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (ic *InstaClient) AddCookie(c *http.Cookie) error {

	for i := range ic.cookies {

		if ic.cookies[i].Name == c.Name {

			if ic.cookies[i].Value != c.Value {
				fmt.Println(c.Name+" = "+ic.cookies[i].Value+" --> "+c.Value, " ")
			}

			ic.cookies[i] = c
			return nil
		}
	}

	cookies := make([]*http.Cookie, 0, len(ic.cookies)+1)

	for i := range ic.cookies {
		cookies = append(cookies, ic.cookies[i])
	}

	cookies = append(cookies, c)
	ic.cookies = cookies

	return nil
}

func (ic *InstaClient) DeleteCookies() {

	mut.Lock()
	defer mut.Unlock()

	ic.cookies = []*http.Cookie{}
}

func (ic *InstaClient) SendPost(p Post) (*http.Response, error) {

	var ret *http.Response

	if functions.IsEmpty(p.URL) {
		return ret, errors.New("URL not set")
	}

	if functions.IsEmpty(p.Referer) {
		return ret, errors.New("Referer not set")
	}

	if functions.IsEmpty(p.XInstagramAJAX) {
		return ret, errors.New("XInstagramAJAX not set")
	}

	mut.Lock()
	defer mut.Unlock()

	ic.delayClient()

	client, err := ic.newClient(p.SetCookies)
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
	req.Header.Add("X-Instagram-AJAX", p.XInstagramAJAX)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("X-CSRFToken", co.CsrfToken)

	ret, err = client.Do(req)
	if err != nil {
		return ret, err
	}

	if p.SetCookies {

		err = ic.AddCookies(ret.Cookies())
		if err != nil {
			return ret, err
		}
	}

	return ret, nil
}

func (ic *InstaClient) SendGet(g Get) (*http.Response, error) {

	var ret *http.Response

	if functions.IsEmpty(g.URL) {
		return ret, errors.New("URL not set")
	}

	mut.Lock()
	defer mut.Unlock()

	ic.delayClient()

	client, err := ic.newClient(g.SetCookies)
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

	if g.SetCookies {

		err = ic.AddCookies(ret.Cookies())
		if err != nil {
			return ret, err
		}
	}

	return ret, nil
}

func (ic *InstaClient) newClient(setcookies bool) (*http.Client, error) {

	var client *http.Client

	if setcookies {
		_, err := ic.getCookies()
		if err != nil {
			return client, err
		}
	}

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

	if setcookies {
		client.Jar.SetCookies(u, ic.cookies)
	}

	return client, nil
}

func (ic *InstaClient) delayClient() {

	if time.Since(ic.lastRequest) < time.Duration(300*time.Millisecond) {

		delay := rand.Intn(500) + 300
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	ic.lastRequest = time.Now()
}
