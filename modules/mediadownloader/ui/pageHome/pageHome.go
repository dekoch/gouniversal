package pageHome

import (
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dekoch/gouniversal/modules/mediadownloader/downloader"
	"github.com/dekoch/gouniversal/modules/mediadownloader/finder"
	"github.com/dekoch/gouniversal/modules/mediadownloader/global"
	"github.com/dekoch/gouniversal/modules/mediadownloader/lang"
	"github.com/dekoch/gouniversal/modules/mediadownloader/typesMD"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typesMD.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.Home.Menu, "App:MediaDownloader:Home", page.Lang.Home.Title)
}

func Render(page *typesMD.Page, nav *navigation.Navigation, r *http.Request) {

	button := r.FormValue("edit")
	ur := r.FormValue("url")

	type Content struct {
		Lang           lang.Home
		Url            template.HTML
		DownloadHidden template.HTML
		Link           template.HTML
	}
	var c Content

	c.Lang = page.Lang.Home
	c.Url = template.HTML(ur)

	if global.Config.DownloadEnabled {
		c.DownloadHidden = template.HTML("")
	} else {
		c.DownloadHidden = template.HTML("hidden")
	}

	var (
		err   error
		files []typesMD.DownloadFile
	)

	if button == "find" {

		files, err = find(ur, false)
		if err != nil {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
		}
	} else if button == "download" && global.Config.DownloadEnabled {

		files, err = find(ur, true)
		if err != nil {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
		}
	}

	f := ""

	for _, file := range files {

		f += "<a href=\"" + file.Url + "\" download=\"" + file.Filename + "\" target=\"_blank\">" + file.Url + "</a><br>"
	}

	c.Link = template.HTML(f)

	p, err := functions.PageToString(global.Config.UIFileRoot+"home.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func find(ur string, download bool) ([]typesMD.DownloadFile, error) {

	var (
		err  error
		resp *http.Response
		b    []byte
		page string
		ret  []typesMD.DownloadFile
	)

	func() {

		for i := 0; i <= 5; i++ {

			switch i {
			case 0:
				// check input
				if functions.IsEmpty(ur) {

					err = errors.New("bad input")
				}

			case 1:
				_, err = url.Parse(ur)

			case 2:
				resp, err = http.Get(ur)
				if err == nil {
					defer resp.Body.Close()
				}

			case 3:
				b, err = ioutil.ReadAll(resp.Body)
				page = string(b)

			case 4:
				ret, err = finder.Find(ur, page)

			case 5:
				if download {
					go func() {
						for _, file := range ret {
							downloader.Download(file)
						}
					}()
				}
			}

			if err != nil {
				console.Log(err, "")
				return
			}
		}
	}()

	return ret, err
}