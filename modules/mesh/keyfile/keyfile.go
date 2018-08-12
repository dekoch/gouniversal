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

func (hc Keyfile) GetKey() []byte {

	mut.Lock()
	defer mut.Unlock()

	return stringToKey(hc.Key)
}

func (hc *Keyfile) LoadDefaults() {

	mut.Lock()
	defer mut.Unlock()

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
		hc.LoadDefaults()
		hc.SaveConfig()
	}

	mut.Lock()
	defer mut.Unlock()

	b, err := file.ReadFile(configFilePath + header.FileName)
	if err != nil {
		console.Log(err, "")
	}

	err = json.Unmarshal(b, &hc)
	if err != nil {
		console.Log(err, "")
	}

	if config.CheckHeader(hc.Header, header.ContentName) == false {
		err = errors.New("wrong config \"" + configFilePath + header.FileName + "\"")
		console.Log(err, "")
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
func stringToKey(key string) []byte {

	strArr := strings.Split(key, " ")

	byArr := make([]byte, len(strArr))

	for i, str := range strArr {

		b, err := strconv.Atoi(str)
		if err == nil {

			byArr[i] = byte(b)
		}
	}

	return byArr
}
