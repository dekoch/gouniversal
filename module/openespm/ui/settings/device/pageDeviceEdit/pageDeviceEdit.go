package pageDeviceEdit

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dekoch/gouniversal/module/openespm/app"
	"github.com/dekoch/gouniversal/module/openespm/deviceConfig"
	"github.com/dekoch/gouniversal/module/openespm/global"
	"github.com/dekoch/gouniversal/module/openespm/lang"
	"github.com/dekoch/gouniversal/module/openespm/typesOESPM"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

func RegisterPage(page *typesOESPM.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "App:openESPM:Settings:Device:Edit", page.Lang.Settings.Device.Edit.Title)
}

func Render(page *typesOESPM.Page, nav *navigation.Navigation, r *http.Request) {

	button := r.FormValue("edit")

	type content struct {
		Lang     lang.SettingsDeviceEdit
		Device   deviceConfig.Device
		CmbApps  template.HTML
		CmbState template.HTML
	}
	var c content

	c.Lang = page.Lang.Settings.Device.Edit

	// Form input
	id := nav.Parameter("UUID")

	if button == "" {

		if id == "new" {

			id = newDevice()
			nav.RedirectPath(strings.Replace(nav.Path, "UUID=new", "UUID="+id, 1), false)
		}
	} else if button == "apply" {

		err := editDevice(r, id)
		if err == nil {
			nav.RedirectPath("App:openESPM:Settings:Device:List", false)
		} else {
			alert.Message(alert.ERROR, page.Lang.Alert.Error, err.Error(), nav.CurrentPath, nav.User.UUID)
		}

	} else if button == "delete" {

		err := deleteDevice(id)
		if err == nil {
			nav.RedirectPath("App:openESPM:Settings:Device:List", false)
		}
	}

	// copy device from array
	var err error
	c.Device, err = global.DeviceConfig.Get(id)

	// combobox App
	cmbApps := "<select name=\"app\">"
	apps := app.DeviceAppList

	for i := 0; i < len(apps); i++ {

		cmbApps += "<option value=\"" + apps[i] + "\""

		if c.Device.App == apps[i] {
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

		if c.Device.State == i {
			cmbState += " selected"
		}

		cmbState += ">" + statetext + "</option>"
	}
	cmbState += "</select>"
	c.CmbState = template.HTML(cmbState)

	// display device
	p, err := functions.PageToString(global.UiConfig.AppFileRoot+"settings/deviceedit.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
}

func newDevice() string {

	u := uuid.Must(uuid.NewRandom())
	key := uuid.Must(uuid.NewRandom())

	var newDevice deviceConfig.Device
	newDevice.UUID = u.String()
	newDevice.Name = u.String()
	newDevice.State = 1 // active
	newDevice.RequestID = u.String()
	newDevice.RequestKey = key.String()

	apps := app.DeviceAppList

	// select first app as default
	if len(apps) > 0 {
		newDevice.App = apps[0]
	}

	global.DeviceConfig.Add(newDevice)

	err := global.DeviceConfig.SaveConfig()
	if err == nil {
		deviceDataFolder := global.DeviceDataFolder + newDevice.UUID + "/"
		os.MkdirAll(deviceDataFolder, os.ModePerm)
	}

	return u.String()
}

func editDevice(r *http.Request, uid string) error {

	name, _ := functions.CheckFormInput("name", r)
	app, _ := functions.CheckFormInput("app", r)
	state, _ := functions.CheckFormInput("state", r)
	comment, errComment := functions.CheckFormInput("comment", r)
	reqId, _ := functions.CheckFormInput("requestid", r)
	reqKey, _ := functions.CheckFormInput("requestkey", r)

	// check input
	if functions.IsEmpty(name) ||
		functions.IsEmpty(app) ||
		functions.IsEmpty(state) ||
		govalidator.IsNumeric(state) == false ||
		functions.IsEmpty(reqId) ||
		functions.IsEmpty(reqKey) ||
		// content not required
		errComment != nil {

		return errors.New("bad input")
	}

	intState, err := strconv.Atoi(state)
	if err != nil {
		return err
	}

	var d deviceConfig.Device
	d.UUID = uid
	d.Name = name
	d.App = app
	d.State = intState
	d.Comment = comment
	d.RequestID = reqId
	d.RequestKey = reqKey

	err = global.DeviceConfig.Edit(uid, d)
	if err != nil {
		return err
	}

	return global.DeviceConfig.SaveConfig()
}

func deleteDevice(uid string) error {

	global.DeviceConfig.Delete(uid)

	err := global.DeviceConfig.SaveConfig()
	if err != nil {
		return err
	}

	deviceDataFolder := global.DeviceDataFolder + uid + "/"

	if _, err := os.Stat(deviceDataFolder); os.IsNotExist(err) {
		return nil
	}

	return os.RemoveAll(deviceDataFolder)
}
