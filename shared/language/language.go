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

type lang struct {
	name string
	file []byte
}

type Language struct {
	mut     sync.RWMutex
	langDir string
	defLang string
	files   []lang
}

func New(dir string, def interface{}, defname string) Language {

	var lf Language

	lf.langDir = dir
	lf.defLang = defname

	err := functions.CreateDir(lf.langDir)
	if err != nil {
		console.Log(err, "")
	} else {

		if _, err = os.Stat(lf.langDir + defname); os.IsNotExist(err) {

			// if nothing found, create default file
			lf.SaveLang(def, defname)
		}
	}

	lf.LoadLangFiles()

	return lf
}

func (l *Language) SaveLang(lf interface{}, n string) error {

	l.mut.Lock()
	defer l.mut.Unlock()

	b, err := json.Marshal(lf)
	if err != nil {
		console.Log(err, "")
	}

	err = file.WriteFile(l.langDir+n, b)

	return err
}

func (l *Language) loadLang(n string) (lang, error) {

	var lf lang
	var err error

	lf.name = n

	lf.file, err = file.ReadFile(l.langDir + lf.name)
	if err != nil {
		console.Log(err, "")
	}

	return lf, err
}

func (l *Language) loadLangFiles() {

	files, err := ioutil.ReadDir(l.langDir)
	if err != nil {
		console.Log(err, "")
		return
	}

	l.files = nil

	for _, fl := range files {

		lf, err := l.loadLang(fl.Name())
		if err == nil {

			la := make([]lang, 1)
			la[0] = lf

			l.files = append(l.files, la...)
		}
	}
}

func (l *Language) LoadLangFiles() {
	l.mut.Lock()
	defer l.mut.Unlock()

	l.loadLangFiles()
}

func fileToStruct(lf lang, s interface{}) {

	err := json.Unmarshal(lf.file, &s)
	if err != nil {
		console.Log(err, "")
	}
}

func (l *Language) SelectLang(n string, s interface{}) {

	l.mut.RLock()

	// search lang
	for i := 0; i < len(l.files); i++ {

		if n == l.files[i].name {

			fileToStruct(l.files[i], s)

			l.mut.RUnlock()
			return
		}
	}

	l.mut.RUnlock()

	l.mut.Lock()
	defer l.mut.Unlock()

	// if nothing found
	// refresh list
	l.loadLangFiles()
	// search lang
	for i := 0; i < len(l.files); i++ {

		if n == l.files[i].name {

			fileToStruct(l.files[i], s)
			return
		}
	}

	// if nothing found
	// search default
	for i := 0; i < len(l.files); i++ {

		if l.defLang == l.files[i].name {

			fileToStruct(l.files[i], s)
			return
		}
	}
}

func (l *Language) ListNames() []string {

	l.mut.RLock()
	defer l.mut.RUnlock()

	fl := make([]string, 0)
	newName := make([]string, 1)

	for i := 0; i < len(l.files); i++ {

		newName[0] = l.files[i].name
		fl = append(fl, newName...)
	}

	return fl
}
