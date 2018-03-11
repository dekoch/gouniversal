package userManagement

import (
	"encoding/json"
	"gouniversal/config"
	"gouniversal/io/file"
	"gouniversal/program/global"
	"gouniversal/program/groupManagement"
	"gouniversal/program/types"
	"log"
	"os"

	"github.com/google/uuid"
)

const UserFile = "data/config/user"

func SaveUser(uc types.UserConfigFile) error {

	uc.Header = config.BuildHeader("user", "users", 1.0, "user config file")

	if _, err := os.Stat(UserFile); os.IsNotExist(err) {
		// if not found, create default file

		newuser := make([]types.User, 1)

		u := uuid.Must(uuid.NewRandom())

		newuser[0].UUID = u.String()
		newuser[0].Lang = "en"
		// admin/admin
		newuser[0].LoginName = "admin"
		newuser[0].PWDHash = "$2a$14$ueP7ISwguEjrGHcHI0SKjO2Jn/A2CjFsWA7LEWgV0FcPNwI7tetde"

		groups := []string{"admin"}
		newuser[0].Groups = groups

		uc.User = newuser
	}

	b, err := json.Marshal(uc)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(UserFile, b)

	return err
}

func LoadUser() types.UserConfigFile {

	var uc types.UserConfigFile

	if _, err := os.Stat(UserFile); os.IsNotExist(err) {
		// if not found, create default file
		SaveUser(uc)
	}

	f := new(file.File)
	b, err := f.ReadFile(UserFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &uc)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(uc.Header, "users") == false {
		log.Fatal("wrong config")
	}

	return uc
}

func SelectUser(uid string) types.User {

	global.UserConfig.Mut.Lock()
	defer global.UserConfig.Mut.Unlock()

	for u := 0; u < len(global.UserConfig.File.User); u++ {

		// search user with UUID
		if uid == global.UserConfig.File.User[u].UUID {

			return global.UserConfig.File.User[u]
		}
	}

	var user types.User
	return user
}

func IsUserInGroup(gid string, user types.User) bool {

	for i := 0; i < len(user.Groups); i++ {

		// search group with UUID
		if gid == user.Groups[i] {

			return true
		}
	}

	return false
}

func IsPageAllowed(pname string, user types.User) bool {

	if pname == "Program:Login" ||
		pname == "Program:Logout" ||
		pname == "Program:Home" {

		return true
	}

	for i := 0; i < len(user.Groups); i++ {

		// test each group
		if groupManagement.IsPageAllowed(pname, user.Groups[i]) {

			return true
		}
	}

	return false
}
