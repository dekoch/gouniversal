package deviceManagement

import (
	"encoding/json"
	"errors"
	"fmt"
	"gouniversal/modules/openespm/globalOESPM"
	"gouniversal/modules/openespm/typesOESPM"
	"gouniversal/shared/config"
	"gouniversal/shared/functions"
	"gouniversal/shared/io/file"
	"html/template"
	"log"
	"os"
)

const DeviceFile = "data/config/openespm/devices"

func SaveConfig(dc typesOESPM.DeviceConfigFile) error {

	dc.Header = config.BuildHeader("devices", "devices", 1.0, "device config file")

	if _, err := os.Stat(DeviceFile); os.IsNotExist(err) {
		// if not found, create default file

		newDevice := make([]typesOESPM.Device, 1)

		newDevice[0].UUID = "test"
		newDevice[0].Key = "1234"
		newDevice[0].Name = "Test"
		newDevice[0].State = 1 // active
		newDevice[0].App = "SimpleSwitchV1x0"

		dc.Devices = newDevice
	}

	b, err := json.Marshal(dc)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(DeviceFile, b)

	return err
}

func LoadConfig() typesOESPM.DeviceConfigFile {

	var dc typesOESPM.DeviceConfigFile

	if _, err := os.Stat(DeviceFile); os.IsNotExist(err) {
		// if not found, create default file
		SaveConfig(dc)
	}

	f := new(file.File)
	b, err := f.ReadFile(DeviceFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &dc)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(dc.Header, "devices") == false {
		log.Fatal("wrong config")
	}

	return dc
}

func LoadDevice(uid string) (typesOESPM.Device, error) {

	globalOESPM.DeviceConfig.Mut.Lock()
	defer globalOESPM.DeviceConfig.Mut.Unlock()

	for u := 0; u < len(globalOESPM.DeviceConfig.File.Devices); u++ {

		// search device with UUID
		if uid == globalOESPM.DeviceConfig.File.Devices[u].UUID {

			return globalOESPM.DeviceConfig.File.Devices[u], nil
		}
	}

	var device typesOESPM.Device
	device.State = -1
	return device, errors.New("LoadDevice() device not found")
}

func SaveDevice(dev typesOESPM.Device) error {

	globalOESPM.DeviceConfig.Mut.Lock()
	defer globalOESPM.DeviceConfig.Mut.Unlock()

	for u := 0; u < len(globalOESPM.DeviceConfig.File.Devices); u++ {

		// search device with UUID
		if dev.UUID == globalOESPM.DeviceConfig.File.Devices[u].UUID {

			globalOESPM.DeviceConfig.File.Devices[u] = dev
			return nil
		}
	}

	return errors.New("SaveDevice() device not found")
}

func HTMLSelectDevice(name string, appname string, uid string) template.HTML {

	type content struct {
		Title  template.HTML
		Select template.HTML
	}
	var c content

	title := "..."

	sel := "<select name=\"" + name + "\">"

	globalOESPM.DeviceConfig.Mut.Lock()
	defer globalOESPM.DeviceConfig.Mut.Unlock()

	for u := 0; u < len(globalOESPM.DeviceConfig.File.Devices); u++ {

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

	templ, err := template.ParseFiles(globalOESPM.UiConfig.AppFileRoot + "selectdevice.html")
	if err != nil {
		fmt.Println(err)
	}

	p := functions.TemplToString(templ, c)

	return template.HTML(p)
}
