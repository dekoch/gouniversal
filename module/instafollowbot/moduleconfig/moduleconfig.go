package moduleconfig

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/shared/api/instaclient"
	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/hashstor"
	"github.com/dekoch/gouniversal/shared/io/file"
	"github.com/dekoch/gouniversal/shared/types"
)

const configFilePath = "data/config/instafollowbot/"

type ModuleConfig struct {
	Header         config.FileHeader
	UIFileRoot     string
	StaticFileRoot string
	LangFileRoot   string
	FileRoot       string
	DBFile         string
	CheckInterv    int // minutes (0=disabled)
	FollowCnt      int
	UnfollowCnt    int
	UnfollowAfter  int // hours
	Tags           []string
	TagQueryHash   hashstor.HashStor
	MediaQueryHash hashstor.HashStor
	Cookies        instaclient.Cookies
	FollowTime     time.Time
	UnfollowTime   time.Time
}

var (
	header config.FileHeader
	mut    sync.RWMutex
)

func init() {
	header = config.FileHeader{HeaderVersion: 0.0, FileName: "instafollowbot.json", ContentName: "instafollowbot", ContentVersion: 1.0, Comment: ""}
}

func (hc *ModuleConfig) loadDefaults() {

	console.Log("loading defaults \""+configFilePath+header.FileName+"\"", " ")

	hc.UIFileRoot = "data/ui/instafollowbot/1.0/"
	hc.StaticFileRoot = "data/ui/instafollowbot/1.0/static/"
	hc.LangFileRoot = "data/lang/instafollowbot/"

	hc.FileRoot = "data/instafollowbot/"
	hc.DBFile = "data/instafollowbot/instafollowbot.db"

	hc.CheckInterv = -1
	hc.FollowCnt = 15
	hc.UnfollowCnt = 15
	hc.UnfollowAfter = 24

	hc.Tags = append(hc.Tags, "")

	hc.TagQueryHash.Add("")
	hc.MediaQueryHash.Add("")
}

func (hc ModuleConfig) SaveConfig() error {

	mut.RLock()
	defer mut.RUnlock()

	return hc.saveConfig()
}

func (hc ModuleConfig) saveConfig() error {

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

func (hc *ModuleConfig) LoadConfig() error {

	mut.Lock()
	defer mut.Unlock()

	if _, err := os.Stat(configFilePath + header.FileName); os.IsNotExist(err) {
		// if not found, create default file
		hc.loadDefaults()
		hc.saveConfig()
	}

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

	hc.TagQueryHash.Init()
	hc.MediaQueryHash.Init()

	return err
}

func (hc *ModuleConfig) Exit(em *types.ExitMessage) error {

	mut.Lock()
	defer mut.Unlock()

	return nil
}

func (hc *ModuleConfig) GetDBFile() string {

	mut.RLock()
	defer mut.RUnlock()

	return hc.DBFile
}

func (hc *ModuleConfig) SetCheckInterval(minutes int) {

	mut.Lock()
	defer mut.Unlock()

	hc.CheckInterv = minutes
}

func (hc *ModuleConfig) GetCheckInterval() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(hc.CheckInterv) * time.Minute
}

func (hc *ModuleConfig) SetFollowCount(cnt int) {

	mut.Lock()
	defer mut.Unlock()

	hc.FollowCnt = cnt
}

func (hc *ModuleConfig) GetFollowCount() int {

	mut.RLock()
	defer mut.RUnlock()

	return hc.FollowCnt
}

func (hc *ModuleConfig) SetUnfollowCount(cnt int) {

	mut.Lock()
	defer mut.Unlock()

	hc.UnfollowCnt = cnt
}

func (hc *ModuleConfig) GetUnfollowCount() int {

	mut.RLock()
	defer mut.RUnlock()

	return hc.UnfollowCnt
}

func (hc *ModuleConfig) SetUnfollowAfter(hours int) {

	mut.Lock()
	defer mut.Unlock()

	hc.UnfollowAfter = hours
}

func (hc *ModuleConfig) GetUnfollowAfter() time.Duration {

	mut.RLock()
	defer mut.RUnlock()

	return time.Duration(hc.UnfollowAfter) * time.Hour
}

func (hc *ModuleConfig) SetTags(tags []string) {

	mut.Lock()
	defer mut.Unlock()

	hc.Tags = tags
}

func (hc *ModuleConfig) GetTags() []string {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Tags
}

func (hc *ModuleConfig) SetCookies(co instaclient.Cookies) {

	mut.Lock()
	defer mut.Unlock()

	hc.Cookies = co
}

func (hc *ModuleConfig) GetCookies() instaclient.Cookies {

	mut.RLock()
	defer mut.RUnlock()

	return hc.Cookies
}

func (hc *ModuleConfig) SetFollowTime(t time.Time) {

	mut.Lock()
	defer mut.Unlock()

	hc.FollowTime = t
}

func (hc *ModuleConfig) GetFollowTime() time.Time {

	mut.RLock()
	defer mut.RUnlock()

	return hc.FollowTime
}

func (hc *ModuleConfig) SetUnfollowTime(t time.Time) {

	mut.Lock()
	defer mut.Unlock()

	hc.UnfollowTime = t
}

func (hc *ModuleConfig) GetUnfollowTime() time.Time {

	mut.RLock()
	defer mut.RUnlock()

	return hc.UnfollowTime
}
