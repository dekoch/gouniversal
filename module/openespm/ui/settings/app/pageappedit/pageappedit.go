package pageappedit

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dekoch/gouniversal/module/openespm/app"
	"github.com/dekoch/gouniversal/module/openespm/appconfig"
	"github.com/dekoch/gouniversal/module/openespm/global"
	"github.com/dekoch/gouniversal/module/openespm/lang"
	"github.com/dekoch/gouniversal/module/openespm/typeoespm"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

func RegisterPage(page *typeoespm.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "App:openESPM:Settings:App:Edit", page.Lang.Settings.App.Edit.Title)
}

func Render(page *typeoespm.Page, nav *navigation.Navigation, r *http.Request) {

	button := r.FormValue("edit")

	type content struct {
		Lang     lang.SettingsAppEdit
		App      appconfig.App
		CmbApps  template.HTML
		CmbState template.HTML
	}
	var c content

	c.Lang = page.Lang.Settings.App.Edit

	// Form input
	id := nav.Parameter("UUID")

	if button == "" {

		if id == "new" {

			id = newApp()
			alert.Message(alert.INFO, page.Lang.Alert.Info, page.Lang.Settings.App.Edit.InfoGroup, nav.CurrentPath, nav.User.UUID)

			nav.RedirectPath(strings.Replace(nav.Path, "UUID=new", "UUID="+id, 1), false)
		}
	} else if button == "apply" {

		err := editApp(r, id)
		if err == nil {
			nav.RedirectPath("App:openESPM:Settings:App:List", false)
		} else {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err.Error(), nav.CurrentPath, nav.User.UUID)
		}

	} else if button == "delete" {

		err := deleteApp(id)
		if err == nil {
			nav.RedirectPath("App:openESPM:Settings:App:List", false)
		}
	}

	// copy app from array
	var err error
	c.App, err = global.AppConfig.Get(id)

	// combobox App
	cmbApps := "<select name=\"app\">"
	apps := app.UiAppList

	for i := 0; i < len(apps); i++ {

		cmbApps += "<option value=\"" + apps[i] + "\""

		if c.App.App == apps[i] {
			cmbApps += " selected"
		}

		cmbApps += ">" + apps[i] + "</option>"
	}
	cmbApps += "</select>"
	c.CmbApps = template.HTML(cmbApps)

	// combobox State
	cmbState := "<select name=\"state\">"
	statetext := ""

	for i := 1; i <= 2; i++ {

		switch i {
		case 1:
			statetext = page.Lang.Settings.Device.Edit.States.Active
		case 2:
			statetext = page.Lang.Settings.Device.Edit.States.Inactive
		}

		cmbState += "<option value=\"" + strconv.Itoa(i) + "\""

		if c.App.State == i {
			cmbState += " selected"
		}

		cmbState += ">" + statetext + "</option>"
	}
	cmbState += "</select>"
	c.CmbState = template.HTML(cmbState)

	// display app
	p, err := functions.PageToString(global.UiConfig.AppFileRoot+"settings/appedit.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func newApp() string {

	u := uuid.Must(uuid.NewRandom())

	var newApp appconfig.App
	newApp.UUID = u.String()
	newApp.Name = u.String()
	newApp.State = 1 // active

	apps := app.UiAppList

	// select first app as default
	if len(apps) > 0 {
		newApp.App = apps[0]
	}

	global.AppConfig.Add(newApp)

	err := global.AppConfig.SaveConfig()
	if err == nil {
		appDataFolder := global.AppDataFolder + newApp.UUID + "/"
		os.MkdirAll(appDataFolder, os.ModePerm)
	}

	return u.String()
}

func editApp(r *http.Request, uid string) error {

	name, _ := functions.CheckFormInput("name", r)
	app, _ := functions.CheckFormInput("app", r)
	state, _ := functions.CheckFormInput("state", r)
	comment, errComment := functions.CheckFormInput("comment", r)

	// check input
	if functions.IsEmpty(name) ||
		functions.IsEmpty(app) ||
		functions.IsEmpty(state) ||
		govalidator.IsNumeric(state) == false ||
		// content not required
		errComment != nil {

		return errors.New("bad input")
	}

	intState, err := strconv.Atoi(state)
	if err != nil {
		return err
	}

	var a appconfig.App
	a.UUID = uid
	a.Name = name
	a.App = app
	a.State = intState
	a.Comment = comment

	err = global.AppConfig.Edit(uid, a)
	if err != nil {
		return err
	}

	return global.AppConfig.SaveConfig()
}

func deleteApp(uid string) error {

	global.AppConfig.Delete(uid)

	err := global.AppConfig.SaveConfig()
	if err != nil {
		return err
	}

	appDataFolder := global.AppDataFolder + uid + "/"

	if _, err := os.Stat(appDataFolder); os.IsNotExist(err) {
		return nil
	}

	return os.RemoveAll(appDataFolder)
}
