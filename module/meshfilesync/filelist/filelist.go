package filelist

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/mesh/serverinfo"
	"github.com/dekoch/gouniversal/module/meshfilesync/syncfile"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/io/fileinfo"
)

const debug = false

type fileState int

const (
	stateError fileState = 1 + iota
	stateIdentical
	stateNewer
	stateOlder
)

type FileList struct {
	Files    []syncfile.SyncFile
	path     string
	serverID string
}

var (
	mut sync.Mutex
)

func (fl *FileList) Lock() {

	mut.Lock()
}

func (fl *FileList) Unlock() {

	mut.Unlock()
}

func (fl *FileList) SetPath(p string) {

	mut.Lock()
	defer mut.Unlock()

	fl.path = p
}

func (fl *FileList) SetServerID(id string) {

	mut.Lock()
	defer mut.Unlock()

	fl.serverID = id
}

func (fl *FileList) Scan() error {

	mut.Lock()
	defer mut.Unlock()

	// directory from path
	dir := filepath.Dir(fl.path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// if not found, create dir
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	localFiles, err := fileinfo.Get(fl.path, -1, false)
	if err != nil {
		return err
	}

	if debug {
		fmt.Println("---")
	}

	now := time.Now()

	// update files
	for i := range fl.Files {

		for _, lf := range localFiles {

			if lf.IsDir {
				continue
			}

			lf.Path = strings.TrimPrefix(lf.Path, fl.path)

			if fl.Files[i].Path == lf.Path+lf.Name {

				// file was deleted
				if fl.Files[i].Deleted {

					console.Output("DA: "+fl.Files[i].Path, "meshFS")

					fl.Files[i].Deleted = false
					fl.Files[i].AddTime = now
				} else if debug {

					fmt.Println("U: " + fl.Files[i].Path)
				}

				updated, err := fl.Files[i].Update(fl.path)
				if err != nil {
					return err
				}

				if updated {
					console.Output("U: "+fl.Files[i].Path, "meshFS")
				}

				fl.Files[i].AddSource(fl.serverID)
			}
		}
	}

	var found bool
	var n syncfile.SyncFile

	// new files
	for _, lf := range localFiles {

		if lf.IsDir {
			continue
		}

		lf.Path = strings.TrimPrefix(lf.Path, fl.path)

		found = false

		for i := range fl.Files {

			if lf.Path+lf.Name == fl.Files[i].Path {
				found = true
			}
		}

		if found == false {

			console.Output("A: "+lf.Path+lf.Name, "meshFS")

			// new file to list
			n.New(fl.path, lf.Path+lf.Name)
			n.AddSource(fl.serverID)
			n.AddTime = now
			n.Deleted = false
			fl.add(n)
		}
	}

	// deleted files
	for i := range fl.Files {

		found = false

		for _, lf := range localFiles {

			if lf.IsDir {
				continue
			}

			lf.Path = strings.TrimPrefix(lf.Path, fl.path)

			if fl.Files[i].Path == lf.Path+lf.Name {
				found = true
			}
		}

		if found == false {

			if fl.Files[i].Deleted == false {

				console.Output("D: "+fl.Files[i].Path, "meshFS")

				fl.Files[i].Deleted = true
				fl.Files[i].DelTime = now
			} else if debug {

				fmt.Println("D: " + fl.Files[i].Path)
			}

			fl.Files[i].Size = 0
			fl.Files[i].ClearChecksum()
			fl.Files[i].ClearSources()
		}
	}

	return nil
}

func (fl *FileList) Add(n syncfile.SyncFile) {

	mut.Lock()
	defer mut.Unlock()

	fl.add(n)
}

func (fl *FileList) AddList(files []syncfile.SyncFile) {

	mut.Lock()
	defer mut.Unlock()

	for _, n := range files {

		fl.add(n)
	}
}

func (fl *FileList) add(n syncfile.SyncFile) {

	found := false

	for i, o := range fl.Files {

		if n.Path == o.Path {

			found = true

			fl.Files[i] = n
			// keep old ID
			fl.Files[i].ID = o.ID

			if bytes.Compare(fl.Files[i].Checksum, o.Checksum) == 0 {
				fl.Files[i].AddSourceList(o.Sources)
			}

			return
		}
	}

	if found == false {
		fl.Files = append(fl.Files, n)
	}
}

func (fl *FileList) SourceUpdateList(files []syncfile.SyncFile) {

	mut.Lock()
	defer mut.Unlock()

	for _, n := range files {

		fl.sourceUpdate(n)
	}
}

func (fl *FileList) SourceUpdate(n syncfile.SyncFile) {

	mut.Lock()
	defer mut.Unlock()

	fl.sourceUpdate(n)
}

func (fl *FileList) sourceUpdate(n syncfile.SyncFile) {

	for i, o := range fl.Files {

		if n.Path == o.Path {

			if compare(o, n) == stateIdentical {

				for _, src := range o.Sources {

					fl.Files[i].AddSource(src)
				}
			}

			return
		}
	}
}

func (fl *FileList) SourceClean(servers []serverinfo.ServerInfo) {

	mut.Lock()
	defer mut.Unlock()

	for i, _ := range fl.Files {

		fl.Files[i].CleanSources(servers)
	}
}

func (fl *FileList) Get() []syncfile.SyncFile {

	mut.Lock()
	defer mut.Unlock()

	return fl.Files
}

func (fl *FileList) GetFileInfo(path string) (syncfile.SyncFile, error) {

	mut.Lock()
	defer mut.Unlock()

	for i := range fl.Files {

		if fl.Files[i].Path == path {
			return fl.Files[i], nil
		}
	}

	var ret syncfile.SyncFile
	return ret, errors.New("file not found")
}

func (fl *FileList) Reset() {

	mut.Lock()
	defer mut.Unlock()

	var empty []syncfile.SyncFile
	fl.Files = empty
}

func (fl *FileList) GetLocalMissing(files []syncfile.SyncFile) []syncfile.SyncFile {

	mut.Lock()
	defer mut.Unlock()

	var ret []syncfile.SyncFile

	var add bool
	var found bool

	for _, remote := range files {

		if remote.Deleted {
			continue
		}

		add = false
		found = false

		for _, local := range fl.Files {

			if local.Path == remote.Path {

				found = true

				if local.Deleted &&
					remote.AddTime.After(local.DelTime) {

					add = true
				}
			}
		}

		if add || found == false {
			ret = append(ret, remote)
		}
	}

	return ret
}

func (fl *FileList) GetLocalOutdated(files []syncfile.SyncFile) []syncfile.SyncFile {

	mut.Lock()
	defer mut.Unlock()

	var ret []syncfile.SyncFile

	var add bool

	for _, remote := range files {

		if remote.Deleted {
			continue
		}

		add = false

		for _, local := range fl.Files {

			if local.Deleted {
				continue
			}

			if local.Path == remote.Path {

				if compare(local, remote) == stateNewer {
					// remote is newer
					add = true
				}
			}
		}

		if add {
			ret = append(ret, remote)
		}
	}

	return ret
}

func (fl *FileList) GetRemoteDeleted(files []syncfile.SyncFile) []syncfile.SyncFile {

	mut.Lock()
	defer mut.Unlock()

	var ret []syncfile.SyncFile

	for _, local := range fl.Files {

		for _, remote := range files {

			if local.Path == remote.Path &&
				local.Deleted == false && remote.Deleted &&
				remote.DelTime.After(local.AddTime) {

				ret = append(ret, remote)
			}
		}
	}

	return ret
}

func (fl *FileList) Delete(path string) {

	mut.Lock()
	defer mut.Unlock()

	fl.delete(path)
}

func (fl *FileList) delete(path string) {

	var l []syncfile.SyncFile

	for _, file := range fl.Files {

		if path != file.Path {

			l = append(l, file)
		}
	}

	fl.Files = l
}

func (fl *FileList) MarkAsDeleted(path string, deltime time.Time) {

	mut.Lock()
	defer mut.Unlock()

	for i := range fl.Files {

		if fl.Files[i].Path == path {

			fl.Files[i].DelTime = deltime
			fl.Files[i].Deleted = true

			return
		}
	}
}

func (fl *FileList) Exists(path string) bool {

	mut.Lock()
	defer mut.Unlock()

	for i := range fl.Files {

		if fl.Files[i].Path == path {
			return true
		}
	}

	return false
}

func (fl *FileList) SetIncomingTime(path string, t time.Time) error {

	mut.Lock()
	defer mut.Unlock()

	for i := range fl.Files {

		if fl.Files[i].Path == path {

			fl.Files[i].SetIncomingTime(t)

			return nil
		}
	}

	return errors.New("file not found")
}

func (fl *FileList) GetIncomingTime(path string) (time.Time, error) {

	mut.Lock()
	defer mut.Unlock()

	for i := range fl.Files {

		if fl.Files[i].Path == path {

			return fl.Files[i].GetIncomingTime(), nil
		}
	}

	t := time.Now()
	return t, errors.New("file not found")
}

func (fl *FileList) ClearIncomingFiles(sec float64) {

	mut.Lock()
	defer mut.Unlock()

	for i := range fl.Files {

		if time.Since(fl.Files[i].GetIncomingTime()).Seconds() > sec {

			fl.delete(fl.Files[i].Path)
		}
	}
}

func compare(in1 syncfile.SyncFile, in2 syncfile.SyncFile) fileState {

	ret := stateError

	if in1.Path == in2.Path {

		if bytes.Compare(in1.Checksum, in2.Checksum) == 0 {

			ret = stateIdentical
		} else {

			if in2.ModTime.After(in1.ModTime) {
				ret = stateNewer
			} else {
				ret = stateOlder
			}
		}
	}

	if debug {
		switch ret {
		case stateError:
			fmt.Println("error " + in2.Path)

		/*case stateIdentical:
		fmt.Println("identical " + in2.Path)*/

		case stateNewer:
			fmt.Println("newer " + in2.Path)

		case stateOlder:
			fmt.Println("older " + in2.Path)
		}
	}

	return ret
}
