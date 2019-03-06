package groupconfig

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const configFilePath = "data/config/"

type Group struct {
	UUID         string
	Name         string
	State        int
	Comment      string
	CanModify    bool
	AllowedPages []string
}

type GroupConfig struct {
	Header config.FileHeader
	Group  []Group
}

var (
	header config.FileHeader
	mut    sync.RWMutex
)

func init() {
	header = config.FileHeader{HeaderVersion: 0.0, FileName: "group", ContentName: "groups", ContentVersion: 1.0, Comment: "group config file"}
}

func (c *GroupConfig) loadDefaults() {

	console.Log("loading defaults \""+configFilePath+header.FileName+"\"", " ")

	newgroup := make([]Group, 1)

	newgroup[0].UUID = "admin"
	newgroup[0].Name = "admin"
	newgroup[0].State = 1 // active

	pages := []string{"Program:Exit",
		"Program:Settings:User",
		"Program:Settings:User:List",
		"Program:Settings:User:Edit",
		"Program:Settings:Group",
		"Program:Settings:Group:List",
		"Program:Settings:Group:Edit"}
	newgroup[0].AllowedPages = pages

	c.Group = newgroup
}

func (c GroupConfig) SaveConfig() error {

	mut.RLock()
	defer mut.RUnlock()

	c.Header = config.BuildHeaderWithStruct(header)

	b, err := json.Marshal(c)
	if err != nil {
		console.Log(err, "")
		return err
	}

	err = file.WriteFile(configFilePath+header.FileName, b)
	if err != nil {
		console.Log(err, "")
	}

	return err
}

func (c *GroupConfig) LoadConfig() error {

	if _, err := os.Stat(configFilePath + header.FileName); os.IsNotExist(err) {
		// if not found, create default file
		c.loadDefaults()
		c.SaveConfig()
	}

	mut.Lock()
	defer mut.Unlock()

	b, err := file.ReadFile(configFilePath + header.FileName)
	if err != nil {
		console.Log(err, "")
		c.loadDefaults()
	} else {
		err = json.Unmarshal(b, &c)
		if err != nil {
			console.Log(err, "")
			c.loadDefaults()
		}
	}

	if config.CheckHeader(c.Header, header.ContentName) == false {
		err = errors.New("wrong config \"" + configFilePath + header.FileName + "\"")
		console.Log(err, "")
		c.loadDefaults()
	}

	return err
}

func (c *GroupConfig) Add(g Group) {

	mut.Lock()
	defer mut.Unlock()

	c.Group = append(c.Group, g)
}

func (c *GroupConfig) Edit(g Group) error {

	mut.Lock()
	defer mut.Unlock()

	for i := 0; i < len(c.Group); i++ {

		if g.UUID == c.Group[i].UUID {

			c.Group[i] = g
			return nil
		}
	}

	return errors.New("Edit() group \"" + g.UUID + "\" not found")
}

func (c *GroupConfig) Get(uid string) (Group, error) {

	mut.RLock()
	defer mut.RUnlock()

	for i := 0; i < len(c.Group); i++ {

		if uid == c.Group[i].UUID {

			return c.Group[i], nil
		}
	}

	var g Group
	g.State = -1
	return g, errors.New("Get() group \"" + uid + "\" not found")
}

func (c *GroupConfig) List() []Group {

	mut.RLock()
	defer mut.RUnlock()

	return c.Group
}

func (c *GroupConfig) Delete(uid string) {

	mut.Lock()
	defer mut.Unlock()

	var l []Group

	for i := 0; i < len(c.Group); i++ {

		if uid != c.Group[i].UUID {

			l = append(l, c.Group[i])
		}
	}

	c.Group = l
}
