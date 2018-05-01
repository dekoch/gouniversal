package deviceManagement

import (
	"errors"
	"fmt"
	"gouniversal/modules/openespm/deviceConfig"
	"gouniversal/modules/openespm/globalOESPM"
	"gouniversal/shared/functions"
	"html/template"
)

func LoadDevice(uid string) (deviceConfig.Device, error) {

	globalOESPM.DeviceConfig.Mut.Lock()
	defer globalOESPM.DeviceConfig.Mut.Unlock()

	for u := 0; u < len(globalOESPM.DeviceConfig.File.Devices); u++ {

		// search device with UUID
		if uid == globalOESPM.DeviceConfig.File.Devices[u].UUID {

			return globalOESPM.DeviceConfig.File.Devices[u], nil
		}
	}

	var device deviceConfig.Device
	device.State = -1
	return device, errors.New("LoadDevice() device \"" + uid + "\" not found")
}

func LoadDeviceWithReqID(id string) (deviceConfig.Device, error) {

	globalOESPM.DeviceConfig.Mut.Lock()
	defer globalOESPM.DeviceConfig.Mut.Unlock()

	for u := 0; u < len(globalOESPM.DeviceConfig.File.Devices); u++ {

		// search device with RequestID
		if id == globalOESPM.DeviceConfig.File.Devices[u].RequestID {

			return globalOESPM.DeviceConfig.File.Devices[u], nil
		}
	}

	var device deviceConfig.Device
	device.State = -1
	return device, errors.New("LoadDeviceWithReqID() device \"" + id + "\" not found")
}

func SaveDevice(dev deviceConfig.Device) error {

	globalOESPM.DeviceConfig.Mut.Lock()
	defer globalOESPM.DeviceConfig.Mut.Unlock()

	for u := 0; u < len(globalOESPM.DeviceConfig.File.Devices); u++ {

		// search device with UUID
		if dev.UUID == globalOESPM.DeviceConfig.File.Devices[u].UUID {

			globalOESPM.DeviceConfig.File.Devices[u] = dev
			return nil
		}
	}

	return errors.New("SaveDevice() device \"" + dev.UUID + "\" not found")
}

func HTMLSelectDevice(name string, appname string, uid string) template.HTML {

	type content struct {
		Title  template.HTML
		Select template.HTML
	}
	var c content

	title := "..."

	sel := "<select name=\"" + name + "\">"

	if uid == "" {
		sel += "<option value=\"\"></option>"
	}

	globalOESPM.DeviceConfig.Mut.Lock()
	defer globalOESPM.DeviceConfig.Mut.Unlock()

	for u := 0; u < len(globalOESPM.DeviceConfig.File.Devices); u++ {

		// list only devices with the same app
		if appname == globalOESPM.DeviceConfig.File.Devices[u].App {
			sel += "<option value=\"" + globalOESPM.DeviceConfig.File.Devices[u].UUID + "\""

			if uid == globalOESPM.DeviceConfig.File.Devices[u].UUID {
				sel += " selected"

				title = globalOESPM.DeviceConfig.File.Devices[u].Name
			}

			sel += ">" + globalOESPM.DeviceConfig.File.Devices[u].Name + "</option>"
		}
	}

	sel += "</select>"

	c.Title = template.HTML(title)
	c.Select = template.HTML(sel)

	p, err := functions.PageToString(globalOESPM.UiConfig.AppFileRoot+"selectdevice.html", c)
	if err != nil {
		fmt.Println(err)
		p = err.Error()
	}

	return template.HTML(p)
}
