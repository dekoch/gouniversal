package groupManagement

import (
	"encoding/json"
	"gouniversal/config"
	"gouniversal/io/file"
	"gouniversal/program/global"
	"gouniversal/program/types"
	"log"
	"os"
)

const GroupFile = "data/config/group"

func SaveGroup(gc types.GroupConfigFile) error {

	gc.Header = config.BuildHeader("group", "groups", 1.0, "group config file")

	if _, err := os.Stat(GroupFile); os.IsNotExist(err) {
		// if not found, create default file

		newgroup := make([]types.Group, 1)

		newgroup[0].UUID = "admin"
		newgroup[0].Name = "admin"

		pages := []string{"Program:Settings:User", "Program:Settings:User:List", "Program:Settings:User:Edit", "Program:Settings:Group", "Program:Settings:Group:List", "Program:Settings:Group:Edit"}
		newgroup[0].AllowedPages = pages

		gc.Group = newgroup
	}

	b, err := json.Marshal(gc)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(GroupFile, b)

	return err
}

func LoadGroup() types.GroupConfigFile {

	var gc types.GroupConfigFile

	if _, err := os.Stat(GroupFile); os.IsNotExist(err) {
		// if not found, create default file
		SaveGroup(gc)
	}

	f := new(file.File)
	b, err := f.ReadFile(GroupFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &gc)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(gc.Header, "groups") == false {
		log.Fatal("wrong config")
	}

	return gc
}

func IsPageAllowed(pname string, gid string) bool {

	global.GroupConfig.Mut.Lock()
	defer global.GroupConfig.Mut.Unlock()

	for g := 0; g < len(global.GroupConfig.File.Group); g++ {

		if gid == global.GroupConfig.File.Group[g].UUID {

			for p := 0; p < len(global.GroupConfig.File.Group[g].AllowedPages); p++ {

				if pname == global.GroupConfig.File.Group[g].AllowedPages[p] {

					return true
				}
			}

			return false
		}
	}

	return false
}
