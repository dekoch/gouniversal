package language

import (
	"encoding/json"
	"fmt"
	"gouniversal/shared/functions"
	"gouniversal/shared/io/file"
	"io/ioutil"
	"os"
	"sync"
)

type Lang struct {
	Name string
	File []byte
}

type Storage struct {
	Mut   sync.Mutex
	Files []Lang
}

type Language struct {
	LangDir string
	Def     string
	Storage Storage
}

func New(dir string, def interface{}, defname string) Language {

	var lf Language

	lf.LangDir = dir
	lf.Def = defname

	err := functions.CreateDir(lf.LangDir)
	if err != nil {
		fmt.Println(err)
	} else {

		if _, err = os.Stat(lf.LangDir + defname); os.IsNotExist(err) {

			// if nothing found, create default file
			lf.SaveLang(def, defname)
		}
	}

	lf.LoadLangFiles()

	return lf
}

func (l Language) SaveLang(lf interface{}, n string) error {

	b, err := json.Marshal(lf)
	if err != nil {
		fmt.Println(err)
	}

	f := new(file.File)
	err = f.WriteFile(l.LangDir+n, b)

	return err
}

func (l *Language) loadLang(n string) (Lang, error) {

	var lf Lang
	var err error

	lf.Name = n

	f := new(file.File)
	lf.File, err = f.ReadFile(l.LangDir + lf.Name)
	if err != nil {
		fmt.Println(err)
	}

	return lf, err
}

func (l *Language) LoadLangFiles() {

	files, err := ioutil.ReadDir(l.LangDir)
	if err != nil {
		fmt.Println(err)
	}

	l.Storage.Mut.Lock()
	defer l.Storage.Mut.Unlock()

	l.Storage.Files = nil

	for _, fl := range files {

		lf, err := l.loadLang(fl.Name())
		if err == nil {

			la := make([]Lang, 1)
			la[0] = lf

			l.Storage.Files = append(la, l.Storage.Files...)
		}
	}
}

func fileToStruct(lf Lang, s interface{}) {

	err := json.Unmarshal(lf.File, &s)
	if err != nil {
		fmt.Println(err)
	}
}

func (l *Language) SelectLang(n string, s interface{}) {

	// search lang
	for i := 0; i < len(l.Storage.Files); i++ {

		if n == l.Storage.Files[i].Name {

			fileToStruct(l.Storage.Files[i], s)
			return
		}
	}

	// if nothing found
	// refresh list
	l.LoadLangFiles()
	// search lang
	for i := 0; i < len(l.Storage.Files); i++ {

		if n == l.Storage.Files[i].Name {

			fileToStruct(l.Storage.Files[i], s)
			return
		}
	}

	// if nothing found
	// search default
	for i := 0; i < len(l.Storage.Files); i++ {

		if l.Def == l.Storage.Files[i].Name {

			fileToStruct(l.Storage.Files[i], s)
			return
		}
	}
}
