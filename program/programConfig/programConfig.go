package programConfig

import (
	"encoding/json"
	"gouniversal/config"
	"gouniversal/io/file"
	"gouniversal/program/types"
	"log"
	"os"
)

const ConfigFilePath = "data/config/program"

func SaveConfig(mc types.ProgramConfig) error {

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

func LoadConfig() types.ProgramConfig {

	var mc types.ProgramConfig

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
