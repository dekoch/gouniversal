// package to load/save the mesh key from/to file

package keyfile

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/dekoch/gouniversal/shared/aes"
	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const configFilePath = "data/config/mesh/"

type Keyfile struct {
	Header config.FileHeader
	Key    string
}

var (
	header config.FileHeader
	mut    sync.Mutex
)

func init() {
	header = config.FileHeader{HeaderVersion: 0.0, FileName: "network.key", ContentName: "networkkey", ContentVersion: 1.0, Comment: ""}
}

func (hc *Keyfile) SetKey(key []byte) {

	mut.Lock()
	defer mut.Unlock()

	hc.Key = keyToString(key)
}

func (hc Keyfile) GetKey() ([]byte, error) {

	mut.Lock()
	defer mut.Unlock()

	return stringToKey(hc.Key)
}

func (hc *Keyfile) loadDefaults() {

	mut.Lock()
	defer mut.Unlock()

	console.Log("loading defaults \""+configFilePath+header.FileName+"\"", " ")

	hc.Header = config.BuildHeaderWithStruct(header)

	key, err := aes.NewKey(32)
	if err != nil {
		fmt.Println(err)
	} else {

		hc.Key = keyToString(key)
	}
}

func (hc Keyfile) SaveConfig() error {

	mut.Lock()
	defer mut.Unlock()

	hc.Header = config.BuildHeaderWithStruct(header)

	b, err := json.Marshal(hc)
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

func (hc *Keyfile) LoadConfig() error {

	if _, err := os.Stat(configFilePath + header.FileName); os.IsNotExist(err) {
		// if not found, create default file
		hc.loadDefaults()
		hc.SaveConfig()
	}

	mut.Lock()
	defer mut.Unlock()

	b, err := file.ReadFile(configFilePath + header.FileName)
	if err != nil {
		console.Log(err, "")
		hc.loadDefaults()
	} else {
		err = json.Unmarshal(b, &hc)
		if err != nil {
			console.Log(err, "")
			hc.loadDefaults()
		}
	}

	if config.CheckHeader(hc.Header, header.ContentName) == false {
		err = errors.New("wrong config \"" + configFilePath + header.FileName + "\"")
		console.Log(err, "")
		hc.loadDefaults()
	}

	return err
}

// convert byte array to space separated string
func keyToString(key []byte) string {

	str := ""

	for _, b := range key {

		str += strconv.Itoa(int(b)) + " "
	}

	str = str[:len(str)-1]

	return str
}

// convert space separated string to byte array
func stringToKey(key string) ([]byte, error) {

	var (
		err error
		ret []byte
		b   int
	)

	for _, s := range strings.Split(key, " ") {

		b, err = strconv.Atoi(s)
		if err != nil {
			return ret, err
		}

		ret = append(ret, byte(b))
	}

	return ret, nil
}
