package modbusConfig

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const configFilePath = "data/config/modbustest/modbus"

type Station struct {
	Active      bool
	ReadOffset  uint16
	WriteOffset uint16
}

type ModbusConfigFile struct {
	Header   config.FileHeader
	IP       string
	Port     string
	Station1 Station
	Station2 Station
}

type ModbusConfig struct {
	Mut  sync.Mutex
	File ModbusConfigFile
}

func (mc *ModbusConfig) SaveConfig() error {

	mc.Mut.Lock()
	defer mc.Mut.Unlock()

	mc.File.Header = config.BuildHeader("modbus", "ModbusConfig", 1.0, "Modbus Settings")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file

		mc.File.IP = "127.0.0.1"
		mc.File.Port = "502"

		mc.File.Station1.Active = true
		mc.File.Station1.ReadOffset = 0
		mc.File.Station1.WriteOffset = 0

		mc.File.Station2.Active = true
		mc.File.Station2.ReadOffset = 64
		mc.File.Station2.WriteOffset = 64
	}

	b, err := json.Marshal(mc.File)
	if err != nil {
		console.Log(err, "modbus SaveConfig()")
	}

	f := new(file.File)
	err = f.WriteFile(configFilePath, b)

	return err
}

func (mc *ModbusConfig) LoadConfig() {

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file
		mc.SaveConfig()
	}

	mc.Mut.Lock()
	defer mc.Mut.Unlock()

	f := new(file.File)
	b, err := f.ReadFile(configFilePath)
	if err != nil {
		console.Log(err, "modbus LoadConfig()")
	}

	err = json.Unmarshal(b, &mc.File)
	if err != nil {
		console.Log(err, "modbus LoadConfig()")
	}

	if config.CheckHeader(mc.File.Header, "ModbusConfig") == false {
		console.Log("wrong config \""+configFilePath+"\"", "modbus LoadConfig()")
	}
}
