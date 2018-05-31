package language

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/io/file"
)

type Lang struct {
	Name string
	File []byte
}

type Language struct {
	Mut     sync.Mutex
	LangDir string
	Def     string
	Files   []Lang
}

func New(dir string, def interface{}, defname string) Language {

	var lf Language

	lf.LangDir = dir
	lf.Def = defname

	err := functions.CreateDir(lf.LangDir)
	if err != nil {
		console.Log(err, "")
	} else {

		if _, err = os.Stat(lf.LangDir + defname); os.IsNotExist(err) {

			// if nothing found, create default file
			lf.SaveLang(def, defname)
		}
	}

	lf.LoadLangFiles()

	return lf
}

func (l *Language) SaveLang(lf interface{}, n string) error {

	l.Mut.Lock()
	defer l.Mut.Unlock()

	b, err := json.Marshal(lf)
	if err != nil {
		console.Log(err, "")
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
		console.Log(err, "")
	}

	return lf, err
}

func (l *Language) loadLangFiles() {

	files, err := ioutil.ReadDir(l.LangDir)
	if err != nil {
		console.Log(err, "")
		return
	}

	l.Files = nil

	for _, fl := range files {

		lf, err := l.loadLang(fl.Name())
		if err == nil {

			la := make([]Lang, 1)
			la[0] = lf

			l.Files = append(l.Files, la...)
		}
	}
}

func (l *Language) LoadLangFiles() {
	l.Mut.Lock()
	defer l.Mut.Unlock()

	l.loadLangFiles()
}

func fileToStruct(lf Lang, s interface{}) {

	err := json.Unmarshal(lf.File, &s)
	if err != nil {
		console.Log(err, "")
	}
}

func (l *Language) SelectLang(n string, s interface{}) {

	l.Mut.Lock()
	defer l.Mut.Unlock()

	// search lang
	for i := 0; i < len(l.Files); i++ {

		if n == l.Files[i].Name {

			fileToStruct(l.Files[i], s)
			return
		}
	}

	// if nothing found
	// refresh list
	l.loadLangFiles()
	// search lang
	for i := 0; i < len(l.Files); i++ {

		if n == l.Files[i].Name {

			fileToStruct(l.Files[i], s)
			return
		}
	}

	// if nothing found
	// search default
	for i := 0; i < len(l.Files); i++ {

		if l.Def == l.Files[i].Name {

			fileToStruct(l.Files[i], s)
			return
		}
	}
}

func (l *Language) ListNames() []string {

	l.Mut.Lock()
	defer l.Mut.Unlock()

	fl := make([]string, 0)
	newName := make([]string, 1)

	for i := 0; i < len(l.Files); i++ {

		newName[0] = l.Files[i].Name
		fl = append(fl, newName...)
	}

	return fl
}
