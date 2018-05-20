package userConfig

import (
	"encoding/json"
	"errors"
	"gouniversal/shared/config"
	"gouniversal/shared/io/file"
	"log"
	"os"
	"strconv"
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

func (c *UserConfig) SaveConfig() error {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	c.File.Header = config.BuildHeader("user", "users", 1.0, "user config file")

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

		c.File.User = newuser
	}

	b, err := json.Marshal(c.File)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(configFilePath, b)

	return err
}

func (c *UserConfig) LoadConfig() error {

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// if not found, create default file
		c.SaveConfig()
	}

	c.Mut.Lock()
	defer c.Mut.Unlock()

	f := new(file.File)
	b, err := f.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &c.File)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(c.File.Header, "users") == false {
		log.Fatal("wrong config \"" + configFilePath + "\"")
	}

	return err
}

func (c *UserConfig) Add(u User) {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	newUser := make([]User, 1)

	newUser[0] = u

	c.File.User = append(c.File.User, newUser...)
}

func (c *UserConfig) Edit(u User) error {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	for i := 0; i < len(c.File.User); i++ {

		if u.UUID == c.File.User[i].UUID {

			c.File.User[i] = u
			return nil
		}
	}

	return errors.New("Edit() user \"" + u.UUID + "\" not found")
}

func (c *UserConfig) Get(uid string) (User, error) {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	for i := 0; i < len(c.File.User); i++ {

		if uid == c.File.User[i].UUID {

			return c.File.User[i], nil
		}
	}

	var u User
	u.State = -1
	return u, errors.New("Get() user \"" + uid + "\" not found")
}

func (c *UserConfig) GetWithName(name string) (User, error) {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	for i := 0; i < len(c.File.User); i++ {

		if name == c.File.User[i].LoginName {

			return c.File.User[i], nil
		}
	}

	var u User
	u.State = -1
	return u, errors.New("GetWithName() user \"" + name + "\" not found")
}

func (c *UserConfig) GetWithState(state int) (User, error) {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	for i := 0; i < len(c.File.User); i++ {

		if state == c.File.User[i].State {

			return c.File.User[i], nil
		}
	}

	var u User
	u.State = -1
	sState := strconv.Itoa(state)
	return u, errors.New("GetWithState() user \"" + sState + "\" not found")
}

func (c *UserConfig) List() []User {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	return c.File.User
}

func (c *UserConfig) Delete(uid string) {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	var l []User
	n := make([]User, 1)

	for i := 0; i < len(c.File.User); i++ {

		if uid != c.File.User[i].UUID {

			n[0] = c.File.User[i]

			l = append(l, n...)
		}
	}

	c.File.User = l
}
