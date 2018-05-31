package groupConfig

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const configFilePath = "data/config/group"

type Group struct {
	UUID         string
	Name         string
	State        int
	Comment      string
	CanModify    bool
	AllowedPages []string
}

type GroupConfigFile struct {
	Header config.FileHeader
	Group  []Group
}

type GroupConfig struct {
	Mut  sync.Mutex
	File GroupConfigFile
}

func (c *GroupConfig) SaveConfig() error {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	c.File.Header = config.BuildHeader("group", "groups", 1.0, "group config file")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file

		newgroup := make([]Group, 1)

		newgroup[0].UUID = "admin"
		newgroup[0].Name = "admin"
		newgroup[0].State = 1 // active

		pages := []string{"Program:Settings:User", "Program:Settings:User:List", "Program:Settings:User:Edit", "Program:Settings:Group", "Program:Settings:Group:List", "Program:Settings:Group:Edit"}
		newgroup[0].AllowedPages = pages

		c.File.Group = newgroup
	}

	b, err := json.Marshal(c.File)
	if err != nil {
		console.Log(err, "groupConfig.SaveConfig()")
	}

	f := new(file.File)
	err = f.WriteFile(configFilePath, b)

	return err
}

func (c *GroupConfig) LoadConfig() error {

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file
		c.SaveConfig()
	}

	c.Mut.Lock()
	defer c.Mut.Unlock()

	f := new(file.File)
	b, err := f.ReadFile(configFilePath)
	if err != nil {
		console.Log(err, "groupConfig.LoadConfig()")
	}

	err = json.Unmarshal(b, &c.File)
	if err != nil {
		console.Log(err, "groupConfig.LoadConfig()")
	}

	if config.CheckHeader(c.File.Header, "groups") == false {
		console.Log("wrong config \""+configFilePath+"\"", "groupConfig.LoadConfig()")
	}

	return err
}

func (c *GroupConfig) Add(g Group) {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	newGroup := make([]Group, 1)

	newGroup[0] = g

	c.File.Group = append(c.File.Group, newGroup...)
}

func (c *GroupConfig) Edit(g Group) error {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	for i := 0; i < len(c.File.Group); i++ {

		if g.UUID == c.File.Group[i].UUID {

			c.File.Group[i] = g
			return nil
		}
	}

	return errors.New("Edit() group \"" + g.UUID + "\" not found")
}

func (c *GroupConfig) Get(uid string) (Group, error) {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	for i := 0; i < len(c.File.Group); i++ {

		if uid == c.File.Group[i].UUID {

			return c.File.Group[i], nil
		}
	}

	var g Group
	g.State = -1
	return g, errors.New("Get() group \"" + uid + "\" not found")
}

func (c *GroupConfig) List() []Group {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	return c.File.Group
}

func (c *GroupConfig) Delete(uid string) {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	var l []Group
	n := make([]Group, 1)

	for i := 0; i < len(c.File.Group); i++ {

		if uid != c.File.Group[i].UUID {

			n[0] = c.File.Group[i]

			l = append(l, n...)
		}
	}

	c.File.Group = l
}
