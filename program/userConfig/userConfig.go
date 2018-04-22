package userConfig

import (
	"encoding/json"
	"gouniversal/shared/config"
	"gouniversal/shared/io/file"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
)

const configFilePath = "data/config/user"

// User stores all information about a single user
type User struct {
	UUID      string
	LoginName string
	Name      string
	PWDHash   string
	Groups    []string
	State     int
	Lang      string
	Comment   string
}

type UserConfigFile struct {
	Header config.FileHeader
	User   []User
}

type UserConfig struct {
	Mut  sync.Mutex
	File UserConfigFile
}

func (uc UserConfig) SaveConfig() error {

	uc.File.Header = config.BuildHeader("user", "users", 1.0, "user config file")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file

		newuser := make([]User, 1)

		u := uuid.Must(uuid.NewRandom())

		newuser[0].UUID = u.String()
		newuser[0].Lang = "en"
		newuser[0].State = 1 // active
		// admin/admin
		newuser[0].LoginName = "admin"
		newuser[0].PWDHash = "$2a$14$ueP7ISwguEjrGHcHI0SKjO2Jn/A2CjFsWA7LEWgV0FcPNwI7tetde"

		groups := []string{"admin"}
		newuser[0].Groups = groups

		uc.File.User = newuser
	}

	b, err := json.Marshal(uc.File)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(configFilePath, b)

	return err
}

func (uc *UserConfig) LoadConfig() error {

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file
		uc.SaveConfig()
	}

	f := new(file.File)
	b, err := f.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &uc.File)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(uc.File.Header, "users") == false {
		log.Fatal("wrong config \"" + configFilePath + "\"")
	}

	return err
}
