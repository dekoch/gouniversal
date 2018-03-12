package programConfig

import (
	"encoding/json"
	"gouniversal/program/programTypes"
	"gouniversal/shared/config"
	"gouniversal/shared/io/file"
	"log"
	"os"
)

const ConfigFilePath = "data/config/program"

func SaveConfig(mc programTypes.ProgramConfig) error {

	mc.Header = config.BuildHeader("program", "ProgramConfig", 1.0, "Program Settings")

	if _, err := os.Stat(ConfigFilePath); os.IsNotExist(err) {
		// if not found, create default file
	}

	b, err := json.Marshal(mc)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(ConfigFilePath, b)

	return err
}

func LoadConfig() programTypes.ProgramConfig {

	var mc programTypes.ProgramConfig

	if _, err := os.Stat(ConfigFilePath); os.IsNotExist(err) {
		// if not found, create default file
		SaveConfig(mc)
	}

	f := new(file.File)
	b, err := f.ReadFile(ConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &mc)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(mc.Header, "ProgramConfig") == false {
		log.Fatal("wrong config")
	}

	return mc
}
