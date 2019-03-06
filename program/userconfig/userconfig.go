package userconfig

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"sync"

	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"

	"github.com/google/uuid"
)

const configFilePath = "data/config/"

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

type UserConfig struct {
	Header config.FileHeader
	User   []User
}

var (
	header config.FileHeader
	mut    sync.RWMutex
)

func init() {
	header = config.FileHeader{HeaderVersion: 0.0, FileName: "user", ContentName: "users", ContentVersion: 1.0, Comment: "user config file"}
}

func (c *UserConfig) loadDefaults() {

	console.Log("loading defaults \""+configFilePath+header.FileName+"\"", " ")

	newuser := make([]User, 2)

	// admin
	u := uuid.Must(uuid.NewRandom())
	newuser[0].UUID = u.String()
	newuser[0].Lang = "en"
	newuser[0].State = 1 // active
	// admin/admin
	newuser[0].LoginName = "admin"
	newuser[0].PWDHash = "$2a$14$ueP7ISwguEjrGcHI0SKjO2Jn/A2CjFsWA7LEWgV0FcPNwI7tetde"

	groups := []string{"admin"}
	newuser[0].Groups = groups

	// guest
	u = uuid.Must(uuid.NewRandom())
	newuser[1].UUID = u.String()
	newuser[1].Lang = "en"
	newuser[1].State = 0 // public
	// guest
	newuser[1].LoginName = "guest"
	newuser[1].PWDHash = ""

	groups = []string{"admin"}
	newuser[1].Groups = groups

	c.User = newuser
}

func (c UserConfig) SaveConfig() error {

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

func (c *UserConfig) LoadConfig() error {

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

func (c *UserConfig) Add(u User) {

	mut.Lock()
	defer mut.Unlock()

	newUser := make([]User, 1)

	newUser[0] = u

	c.User = append(c.User, newUser...)
}

func (c *UserConfig) Edit(u User) error {

	mut.Lock()
	defer mut.Unlock()

	for i := 0; i < len(c.User); i++ {

		if u.UUID == c.User[i].UUID {

			c.User[i] = u
			return nil
		}
	}

	return errors.New("Edit() user \"" + u.UUID + "\" not found")
}

func (c *UserConfig) Get(uid string) (User, error) {

	mut.RLock()
	defer mut.RUnlock()

	for i := 0; i < len(c.User); i++ {

		if uid == c.User[i].UUID {

			return c.User[i], nil
		}
	}

	var u User
	u.State = -1
	return u, errors.New("Get() user \"" + uid + "\" not found")
}

func (c *UserConfig) GetWithName(name string) (User, error) {

	mut.RLock()
	defer mut.RUnlock()

	for i := 0; i < len(c.User); i++ {

		if name == c.User[i].LoginName {

			return c.User[i], nil
		}
	}

	var u User
	u.State = -1
	return u, errors.New("GetWithName() user \"" + name + "\" not found")
}

func (c *UserConfig) GetWithState(state int) (User, error) {

	mut.RLock()
	defer mut.RUnlock()

	for i := 0; i < len(c.User); i++ {

		if state == c.User[i].State {

			return c.User[i], nil
		}
	}

	var u User
	u.State = -1
	sState := strconv.Itoa(state)
	return u, errors.New("GetWithState() no user with state \"" + sState + "\" found")
}

func (c *UserConfig) GetUserCnt() int {

	mut.RLock()
	defer mut.RUnlock()

	return len(c.User)
}

func (c *UserConfig) List() []User {

	mut.RLock()
	defer mut.RUnlock()

	return c.User
}

func (c *UserConfig) Delete(uid string) {

	mut.Lock()
	defer mut.Unlock()

	var l []User
	n := make([]User, 1)

	for i := 0; i < len(c.User); i++ {

		if uid != c.User[i].UUID {

			n[0] = c.User[i]

			l = append(l, n...)
		}
	}

	c.User = l
}
