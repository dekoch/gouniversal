package pageviewer

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/dekoch/gouniversal/module/monmotion/dbstorage"
	"github.com/dekoch/gouniversal/module/monmotion/global"
	"github.com/dekoch/gouniversal/module/monmotion/lang"
	"github.com/dekoch/gouniversal/module/monmotion/typemd"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func RegisterPage(page *typemd.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.Viewer.Menu, "App:MonMotion:Viewer", page.Lang.Viewer.Title)
}

func Render(page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	type Content struct {
		Lang       lang.Viewer
		UUID       template.HTML
		Token      template.HTML
		CmbTrigger template.HTML
		Pictures   template.HTML
		Interval   template.JS
	}
	var c Content

	c.Lang = page.Lang.Viewer

	c.UUID = template.HTML(nav.User.UUID)

	var (
		err error
	)

	func() {

		var (
			selID    string
			seqInfos []dbstorage.SequenceImage
		)

		for i := 0; i <= 5; i++ {

			switch i {
			case 0:
				selID, err = functions.CheckFormInput("trigger", r)

			case 1:
				c.CmbTrigger, err = cmbTrigger(selID)

			case 2:
				if functions.IsEmpty(selID) {
					return
				}

			case 3:
				seqInfos, err = dbstorage.GetSequenceInfos(selID)

			case 4:
				c.Pictures, err = pictures(seqInfos, nav.User.UUID, global.UIRequest.GetNewToken(nav.User.UUID))

			case 5:
				c.Interval, err = getInterval(seqInfos)
			}

			if err != nil {
				alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
				return
			}
		}
	}()

	p, err := functions.PageToString(global.Config.UIFileRoot+"viewer.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func cmbTrigger(selid string) (template.HTML, error) {

	ids, err := dbstorage.GetTriggerIDs()
	if err != nil {
		return template.HTML(""), err
	}

	tag := "<select name=\"trigger\">"
	// empty
	tag += "<option value=\"\""
	if functions.IsEmpty(selid) {
		tag += " selected"
	}
	tag += "></option>"

	for _, id := range ids {

		img, err := dbstorage.LoadImage(id)
		if err != nil {
			return template.HTML(""), err
		}

		tag += "<option value=\"" + id + "\""
		if selid == id {
			tag += " selected"
		}
		tag += ">" + img.Captured.Format("2006.01.02 15:04:05.0000") + "</option>"
	}

	tag += "</select>"

	return template.HTML(tag), nil
}

func pictures(seqinfos []dbstorage.SequenceImage, uuid, token string) (template.HTML, error) {

	var tag string

	for _, seqInfo := range seqinfos {

		path := "/monmotion/viewer/?uuid=" + uuid + "&token=" + token + "&imageid=" + seqInfo.ID
		name := seqInfo.Captured.Format("2006.01.02 15:04:05.0000")

		tag += "<li><img data-original=\"" + path + "\" src=\"" + path + "\" alt=\"" + name + "\"></li>"
	}

	return template.HTML(tag), nil
}

func getInterval(seqinfos []dbstorage.SequenceImage) (template.JS, error) {

	l := len(seqinfos)

	if l <= 1 {
		return template.JS("33"), nil
	}

	t := seqinfos[l-1].Captured.Sub(seqinfos[0].Captured).Milliseconds()

	interval := float64(t) / float64(l)

	if interval <= 0.0 {
		interval = 33.0
	}

	return template.JS(strconv.FormatFloat(interval, 'f', 0, 64)), nil
}

func edit(r *http.Request) error {

	var err error

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:

			case 1:

			}

			if err != nil {
				return
			}
		}
	}()

	return err
}
