package pageDeviceEdit

import (
	"errors"
	"fmt"
	"gouniversal/modules/openespm/app"
	"gouniversal/modules/openespm/deviceManagement"
	"gouniversal/modules/openespm/globalOESPM"
	"gouniversal/modules/openespm/langOESPM"
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/shared/alert"
	"gouniversal/shared/functions"
	"gouniversal/shared/navigation"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

func RegisterPage(page *typesOESPM.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("", "App:openESPM:Settings:Device:Edit", page.Lang.Settings.Device.Edit.Title)
}

func Render(page *typesOESPM.Page, nav *navigation.Navigation, r *http.Request) {

	button := r.FormValue("edit")

	type deviceEdit struct {
		Lang     langOESPM.SettingsDeviceEdit
		Device   typesOESPM.Device
		CmbApps  template.HTML
		CmbState template.HTML
	}
	var de deviceEdit

	de.Lang = page.Lang.Settings.Device.Edit

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
	globalOESPM.DeviceConfig.Mut.Lock()
	for i := 0; i < len(globalOESPM.DeviceConfig.File.Devices); i++ {

		if id == globalOESPM.DeviceConfig.File.Devices[i].UUID {

			de.Device = globalOESPM.DeviceConfig.File.Devices[i]
		}
	}
	globalOESPM.DeviceConfig.Mut.Unlock()

	// combobox App
	cmbApps := "<select name=\"app\">"
	apps := app.List()

	for i := 0; i < len(apps); i++ {

		cmbApps += "<option value=\"" + apps[i] + "\""

		if de.Device.App == apps[i] {
			cmbApps += " selected"
		}

		cmbApps += ">" + apps[i] + "</option>"
	}
	cmbApps += "</select>"
	de.CmbApps = template.HTML(cmbApps)

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

		if de.Device.State == i {
			cmbState += " selected"
		}

		cmbState += ">" + statetext + "</option>"
	}
	cmbState += "</select>"
	de.CmbState = template.HTML(cmbState)

	// display device
	templ, err := template.ParseFiles(globalOESPM.UiConfig.AppFileRoot + "settings/deviceedit.html")
	if err != nil {
		fmt.Println(err)
	}

	page.Content += functions.TemplToString(templ, de)
}

func newDevice() string {

	globalOESPM.DeviceConfig.Mut.Lock()
	defer globalOESPM.DeviceConfig.Mut.Unlock()

	u := uuid.Must(uuid.NewRandom())
	key := uuid.Must(uuid.NewRandom())

	newDevice := make([]typesOESPM.Device, 1)
	newDevice[0].UUID = u.String()
	newDevice[0].Name = u.String()
	newDevice[0].State = 1 // active
	newDevice[0].Key = key.String()

	apps := app.List()

	// select first app as default
	if len(apps) > 0 {
		newDevice[0].App = apps[0]
	}

	globalOESPM.DeviceConfig.File.Devices = append(newDevice, globalOESPM.DeviceConfig.File.Devices...)

	deviceManagement.SaveConfig(globalOESPM.DeviceConfig.File)

	return u.String()
}

func editDevice(r *http.Request, u string) error {

	name, _ := functions.CheckFormInput("name", r)
	app, _ := functions.CheckFormInput("app", r)
	state, _ := functions.CheckFormInput("state", r)
	comment, errComment := functions.CheckFormInput("comment", r)
	id, _ := functions.CheckFormInput("uuid", r)
	key, _ := functions.CheckFormInput("key", r)

	// check input
	if functions.IsEmpty(name) ||
		functions.IsEmpty(app) ||
		functions.IsEmpty(state) ||
		govalidator.IsNumeric(state) == false ||
		functions.IsEmpty(id) ||
		functions.IsEmpty(key) ||
		// content not required
		errComment != nil {

		return errors.New("bad input")
	}

	globalOESPM.DeviceConfig.Mut.Lock()
	defer globalOESPM.DeviceConfig.Mut.Unlock()

	for i := 0; i < len(globalOESPM.DeviceConfig.File.Devices); i++ {

		if u == globalOESPM.DeviceConfig.File.Devices[i].UUID {

			intState, err := strconv.Atoi(state)
			if err != nil {
				return err
			}

			globalOESPM.DeviceConfig.File.Devices[i].Name = name
			globalOESPM.DeviceConfig.File.Devices[i].App = app
			globalOESPM.DeviceConfig.File.Devices[i].State = intState
			globalOESPM.DeviceConfig.File.Devices[i].Comment = comment
			globalOESPM.DeviceConfig.File.Devices[i].UUID = id
			globalOESPM.DeviceConfig.File.Devices[i].Key = key

			return deviceManagement.SaveConfig(globalOESPM.DeviceConfig.File)
		}
	}

	return errors.New("UUID not found")
}

func deleteDevice(u string) error {

	globalOESPM.DeviceConfig.Mut.Lock()
	defer globalOESPM.DeviceConfig.Mut.Unlock()

	var dl []typesOESPM.Device
	n := make([]typesOESPM.Device, 1)

	for i := 0; i < len(globalOESPM.DeviceConfig.File.Devices); i++ {

		if u != globalOESPM.DeviceConfig.File.Devices[i].UUID {

			n[0] = globalOESPM.DeviceConfig.File.Devices[i]

			dl = append(dl, n...)
		}
	}

	globalOESPM.DeviceConfig.File.Devices = dl

	return deviceManagement.SaveConfig(globalOESPM.DeviceConfig.File)
}
