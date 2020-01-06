package menu

import (
	"html/template"
	"net/http"

	"github.com/dekoch/gouniversal/module/monmotion/global"
	"github.com/dekoch/gouniversal/module/monmotion/lang"
	"github.com/dekoch/gouniversal/module/monmotion/typemd"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func Render(src string, page *typemd.Page, nav *navigation.Navigation, r *http.Request) {

	id := nav.Parameter("UUID")

	switch r.FormValue("menu") {
	case "acquire":
		nav.RedirectPath("App:MonMotion:Acquire$UUID="+id, false)
		return

	case "trigger":
		nav.RedirectPath("App:MonMotion:Trigger$UUID="+id, false)
		return
	}

	type Content struct {
		Lang             lang.DeviceMenu
		BtnAcquireActive template.HTMLAttr
		BtnTriggerActive template.HTMLAttr
	}
	var c Content

	c.Lang = page.Lang.Device.DeviceMenu

	c.BtnAcquireActive = template.HTMLAttr("")
	c.BtnTriggerActive = template.HTMLAttr("")

	switch src {
	case "acquire":
		c.BtnAcquireActive = template.HTMLAttr(" btn-lg active")

	case "trigger":
		c.BtnTriggerActive = template.HTMLAttr(" btn-lg active")
	}

	p, err := functions.PageToString(global.Config.UIFileRoot+"menu.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}
