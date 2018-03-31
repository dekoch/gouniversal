package deviceManagement

import (
	"encoding/json"
	"gouniversal/modules/openespm/oespmGlobal"
	"gouniversal/modules/openespm/oespmTypes"
	"gouniversal/shared/config"
	"gouniversal/shared/io/file"
	"log"
	"os"
)

const DeviceFile = "data/config/openespm/devices"

func SaveDevices(dc oespmTypes.DeviceConfigFile) error {

	dc.Header = config.BuildHeader("devices", "devices", 1.0, "device config file")

	if _, err := os.Stat(DeviceFile); os.IsNotExist(err) {
		// if not found, create default file

		newDevice := make([]oespmTypes.Device, 1)

		newDevice[0].UUID = "test"
		newDevice[0].Key = "1234"
		newDevice[0].Name = "Test"
		newDevice[0].State = 1 // active
		newDevice[0].App = "SimpleSwitch_v1_0"

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

func LoadDevices() oespmTypes.DeviceConfigFile {

	var dc oespmTypes.DeviceConfigFile

	if _, err := os.Stat(DeviceFile); os.IsNotExist(err) {
		// if not found, create default file
		SaveDevices(dc)
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

func SelectDevice(uid string) oespmTypes.Device {

	oespmGlobal.DeviceConfig.Mut.Lock()
	defer oespmGlobal.DeviceConfig.Mut.Unlock()

	for u := 0; u < len(oespmGlobal.DeviceConfig.File.Devices); u++ {

		// search user with UUID
		if uid == oespmGlobal.DeviceConfig.File.Devices[u].UUID {

			return oespmGlobal.DeviceConfig.File.Devices[u]
		}
	}

	var device oespmTypes.Device
	device.State = -1
	return device
}
