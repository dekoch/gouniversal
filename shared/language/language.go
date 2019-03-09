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
	langDir string
	defLang string
	files   []lang
}

var mut sync.RWMutex

// New creates a new language instance and adds all existing language files from a given path
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

// SaveLang saves a language struct as file
func (l *Language) SaveLang(lf interface{}, name string) error {

	mut.Lock()
	defer mut.Unlock()

	b, err := json.Marshal(lf)
	if err != nil {
		console.Log(err, "")
	}

	err = file.WriteFile(l.langDir+name, b)

	return err
}

func (l *Language) loadLang(name string) (lang, error) {

	var lf lang
	var err error

	lf.name = name

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
			l.files = append(l.files, lf)
		}
	}
}

// LoadLangFiles adds all existing language files from a given path
func (l *Language) LoadLangFiles() {
	mut.Lock()
	defer mut.Unlock()

	l.loadLangFiles()
}

func fileToStruct(lf lang, s interface{}) {

	err := json.Unmarshal(lf.file, &s)
	if err != nil {
		console.Log(err, "")
	}
}

// SelectLang parses a language file to struct
func (l *Language) SelectLang(name string, s interface{}) {

	mut.RLock()

	// search lang
	for i := 0; i < len(l.files); i++ {

		if name == l.files[i].name {

			fileToStruct(l.files[i], s)

			mut.RUnlock()
			return
		}
	}

	mut.RUnlock()

	mut.Lock()
	defer mut.Unlock()

	// if nothing found
	// refresh list
	l.loadLangFiles()
	// search lang
	for i := 0; i < len(l.files); i++ {

		if name == l.files[i].name {

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

// ListNames returns a list of all language files
func (l *Language) ListNames() []string {

	mut.RLock()
	defer mut.RUnlock()

	fl := make([]string, 0)

	for i := 0; i < len(l.files); i++ {

		fl = append(fl, l.files[i].name)
	}

	return fl
}
