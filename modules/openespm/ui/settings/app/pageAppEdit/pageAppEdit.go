package pageAppEdit

import (
	"errors"
	"gouniversal/modules/openespm/app"
	"gouniversal/modules/openespm/appConfig"
	"gouniversal/modules/openespm/globalOESPM"
	"gouniversal/modules/openespm/langOESPM"
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/shared/alert"
	"gouniversal/shared/functions"
	"gouniversal/shared/navigation"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

func RegisterPage(page *typesOESPM.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "App:openESPM:Settings:App:Edit", page.Lang.Settings.App.Edit.Title)
}

func Render(page *typesOESPM.Page, nav *navigation.Navigation, r *http.Request) {

	button := r.FormValue("edit")

	type content struct {
		Lang     langOESPM.SettingsAppEdit
		App      appConfig.App
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
	globalOESPM.AppConfig.Mut.Lock()
	for i := 0; i < len(globalOESPM.AppConfig.File.Apps); i++ {

		if id == globalOESPM.AppConfig.File.Apps[i].UUID {

			c.App = globalOESPM.AppConfig.File.Apps[i]
		}
	}
	globalOESPM.AppConfig.Mut.Unlock()

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
	p, err := functions.PageToString(globalOESPM.UiConfig.AppFileRoot+"settings/appedit.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func newApp() string {

	globalOESPM.AppConfig.Mut.Lock()
	defer globalOESPM.AppConfig.Mut.Unlock()

	u := uuid.Must(uuid.NewRandom())

	newApp := make([]appConfig.App, 1)
	newApp[0].UUID = u.String()
	newApp[0].Name = u.String()
	newApp[0].State = 1 // active

	apps := app.UiAppList

	// select first app as default
	if len(apps) > 0 {
		newApp[0].App = apps[0]
	}

	globalOESPM.AppConfig.File.Apps = append(newApp, globalOESPM.AppConfig.File.Apps...)

	err := globalOESPM.AppConfig.SaveConfig()
	if err == nil {
		appDataFolder := globalOESPM.AppDataFolder + newApp[0].UUID + "/"
		os.MkdirAll(appDataFolder, os.ModePerm)
	}

	return u.String()
}

func editApp(r *http.Request, u string) error {

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

	globalOESPM.AppConfig.Mut.Lock()
	defer globalOESPM.AppConfig.Mut.Unlock()

	for i := 0; i < len(globalOESPM.AppConfig.File.Apps); i++ {

		if u == globalOESPM.AppConfig.File.Apps[i].UUID {

			intState, err := strconv.Atoi(state)
			if err != nil {
				return err
			}

			globalOESPM.AppConfig.File.Apps[i].Name = name
			globalOESPM.AppConfig.File.Apps[i].App = app
			globalOESPM.AppConfig.File.Apps[i].State = intState
			globalOESPM.AppConfig.File.Apps[i].Comment = comment

			return globalOESPM.AppConfig.SaveConfig()
		}
	}

	return errors.New("UUID not found")
}

func deleteApp(u string) error {

	globalOESPM.AppConfig.Mut.Lock()
	defer globalOESPM.AppConfig.Mut.Unlock()

	var al []appConfig.App
	n := make([]appConfig.App, 1)

	for i := 0; i < len(globalOESPM.AppConfig.File.Apps); i++ {

		if u != globalOESPM.AppConfig.File.Apps[i].UUID {

			n[0] = globalOESPM.AppConfig.File.Apps[i]

			al = append(al, n...)
		}
	}

	globalOESPM.AppConfig.File.Apps = al

	err := globalOESPM.AppConfig.SaveConfig()
	if err != nil {
		return err
	}

	appDataFolder := globalOESPM.AppDataFolder + u + "/"

	if _, err := os.Stat(appDataFolder); os.IsNotExist(err) {
		return nil
	}

	return os.RemoveAll(appDataFolder)
}
