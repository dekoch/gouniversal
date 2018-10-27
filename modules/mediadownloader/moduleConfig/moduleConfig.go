package moduleConfig

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/dekoch/gouniversal/shared/config"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/file"
)

const configFilePath = "data/config/mediadownloader/"

type ModuleConfig struct {
	Header          config.FileHeader
	UIFileRoot      string
	LangFileRoot    string
	DownloadEnabled bool
	FileRoot        string
	Extension       []string
}

var (
	header config.FileHeader
)

func init() {
	header = config.FileHeader{HeaderVersion: 0.0, FileName: "mediadownloader", ContentName: "mediadownloader", ContentVersion: 1.0, Comment: ""}
}

func (hc *ModuleConfig) loadDefaults() {

	console.Log("loading defaults \""+configFilePath+header.FileName+"\"", " ")

	hc.UIFileRoot = "data/ui/mediadownloader/1.0/"
	hc.LangFileRoot = "data/lang/mediadownloader/"

	hc.DownloadEnabled = true
	hc.FileRoot = "data/mediadownloader/"

	extension := make([]string, 24)
	extension[0] = ".png"
	extension[1] = ".jpg"
	extension[2] = ".jpeg"
	extension[3] = ".bmp"
	extension[4] = ".gif"
	extension[5] = ".mp3"
	extension[6] = ".ogg"
	extension[7] = ".mp4"
	extension[8] = ".avi"
	extension[9] = ".opus"
	extension[10] = ".flv"
	extension[11] = ".mkv"
	extension[12] = ".webm"
	extension[13] = ".wmv"
	extension[14] = ".mpg"
	extension[15] = ".mov"
	extension[16] = ".zip"
	extension[17] = ".7z"
	extension[18] = ".tar"
	extension[19] = ".tar.gz"
	extension[20] = ".rar"
	extension[21] = ".iso"
	extension[22] = ".m3u"
	extension[23] = ".m3u8"
	hc.Extension = extension
}

func (hc ModuleConfig) SaveConfig() error {

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

	if _, err := os.Stat(configFilePath + header.FileName); os.IsNotExist(err) {
		// if not found, create default file
		hc.loadDefaults()
		hc.SaveConfig()
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

	return err
}
