package pagehome

import (
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/dekoch/gouniversal/module/mediadownloader/downloader"
	"github.com/dekoch/gouniversal/module/mediadownloader/finder"
	"github.com/dekoch/gouniversal/module/mediadownloader/global"
	"github.com/dekoch/gouniversal/module/mediadownloader/lang"
	"github.com/dekoch/gouniversal/module/mediadownloader/typemd"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.Home.Menu, "App:MediaDownloader:Home", page.Lang.Home.Title)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	button := r.FormValue("edit")
	ur := r.FormValue("url")
	ur = strings.Trim(ur, " ")

	type Content struct {
		Lang                    lang.Home
		Url                     template.HTML
		DownloadHidden          template.HTML
		Link                    template.HTML
		SupportedFileExtensions template.HTML
		Extensions              template.HTML
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
		files []typemd.DownloadFile
	)

	if button == "find" {

		files, err = find(ur, false, page, nav)
		if err != nil {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
		}
	} else if button == "download" && global.Config.DownloadEnabled {

		files, err = find(ur, true, page, nav)
		if err != nil {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
		}
	}

	// extensions
	c.SupportedFileExtensions = template.HTML(page.Lang.Home.SupportedFileExtensions)

	e := ""

	for _, ex := range global.Config.Extension {

		e += ex + " "
	}

	c.Extensions = template.HTML(e)

	// links
	f := ""

	for _, file := range files {

		f += "<a href=\"" + file.Url + "\" download=\"" + file.Filename + "\" target=\"_blank\">" + file.Filename + "</a><br>"
	}

	c.Link = template.HTML(f)

	p, err := functions.PageToString(global.Config.UIFileRoot+"home.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func find(ur string, download bool, page *typemd.Page, nav *navigation.Navigation) ([]typemd.DownloadFile, error) {

	var (
		err  error
		resp *http.Response
		b    []byte
		p    string
		ret  []typemd.DownloadFile
	)

	func() {

		for i := 0; i <= 6; i++ {

			switch i {
			case 0:
				// check input
				if functions.IsEmpty(ur) {
					err = errors.New(page.Lang.Home.PleaseEnterUrl)
				}

			case 1:
				_, err = url.Parse(ur)

			case 2:
				resp, err = http.Get(ur)
				if err == nil {
					defer resp.Body.Close()
				}

			case 3:
				ct := resp.Header.Get("Content-Type")

				if ct != "" {
					if strings.Contains(ct, "text/html") == false {
						err = errors.New(page.Lang.Home.NotSupportedContentType + ": " + ct)
					}
				}

			case 4:
				b, err = ioutil.ReadAll(resp.Body)
				p = string(b)

			case 5:
				ret, err = finder.Find(ur, p)
				if len(ret) == 0 {
					alert.Message(alert.INFO, page.Lang.Alert.Info, page.Lang.Home.NoFileFound, "", nav.User.UUID)
				}

			case 6:
				if download {
					go func() {
						for _, file := range ret {
							downloader.Download(file)
						}

						alert.Message(alert.INFO, page.Lang.Alert.Info, page.Lang.Home.DownloadFinished, "", nav.User.UUID)
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
