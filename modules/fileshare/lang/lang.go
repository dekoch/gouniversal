package lang

import (
	"encoding/json"
	"gouniversal/shared/config"
	"gouniversal/shared/io/file"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

const LangDir = "data/lang/fileshare/"

type Home struct {
	Title   string
	Name    string
	Size    string
	Options string
}

type Alert struct {
	Success string
	Info    string
	Warning string
	Error   string
}

type File struct {
	Header config.FileHeader
	Home   Home
	Alert  Alert
}

type Lang struct {
	Mut   sync.Mutex
	Files []File
}

func SaveLang(lang File, n string) error {

	lang.Header = config.BuildHeader(n, "LangOpenESPM", 1.0, "Language File")

	if _, err := os.Stat(LangDir + n); os.IsNotExist(err) {
		// if not found, create default file
		lang.Home.Title = "Fileshare"
		lang.Home.Name = "Name"
		lang.Home.Size = "Size"
		lang.Home.Options = "Options"

		lang.Alert.Success = "Success"
		lang.Alert.Info = "Info"
		lang.Alert.Warning = "Warning"
		lang.Alert.Error = "Error"
	}

	b, err := json.Marshal(lang)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(LangDir+n, b)

	return err
}

func LoadLang(n string) File {

	var lg File

	if _, err := os.Stat(LangDir + n); os.IsNotExist(err) {
		// if not found, create default file
		SaveLang(lg, n)
	}

	f := new(file.File)
	b, err := f.ReadFile(LangDir + n)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &lg)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(lg.Header, "LangOpenESPM") == false {
		log.Fatal("wrong config")
	}

	return lg
}

func LoadLangFiles() []File {

	lg := make([]File, 0)

	if _, err := os.Stat(LangDir + "en"); os.IsNotExist(err) {
		// if not found, create default file
		var newlg File
		SaveLang(newlg, "en")
	}

	files, err := ioutil.ReadDir(LangDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, fl := range files {

		var langfile File

		f := new(file.File)
		b, err := f.ReadFile(LangDir + fl.Name())
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(b, &langfile)
		if err != nil {
			log.Fatal(err)
		}

		if config.CheckHeader(langfile.Header, "LangOpenESPM") {

			lg = append(lg, langfile)
		}

	}

	return lg
}
