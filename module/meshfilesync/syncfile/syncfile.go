package syncfile

import (
	"errors"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/mesh/serverinfo"
	"github.com/dekoch/gouniversal/shared/io/file"
	"github.com/google/uuid"
)

type SyncFile struct {
	ID          string
	Path        string
	Checksum    []byte
	Size        int64
	AddTime     time.Time
	ModTime     time.Time
	DelTime     time.Time
	Deleted     bool
	Sources     []string
	destination string
}

var (
	mut sync.Mutex
)

func (f *SyncFile) New(root string, path string) error {

	mut.Lock()
	defer mut.Unlock()

	u := uuid.Must(uuid.NewRandom())
	f.ID = u.String()

	f.Path = path

	return f.update(root)
}

func (f *SyncFile) Update(root string) error {

	mut.Lock()
	defer mut.Unlock()

	return f.update(root)
}

func (f *SyncFile) update(root string) error {

	sum, err := file.Checksum(root + f.Path)
	if err != nil {
		return err
	}

	f.Checksum = sum

	fileInfo, err := os.Stat(root + f.Path)
	if err != nil {
		return err
	}

	f.Size = fileInfo.Size()
	f.ModTime = fileInfo.ModTime()

	return nil
}

func (f *SyncFile) AddSource(id string) {

	mut.Lock()
	defer mut.Unlock()

	f.addSource(id)
}

func (f *SyncFile) AddSourceList(ids []string) {

	mut.Lock()
	defer mut.Unlock()

	for _, id := range ids {
		f.addSource(id)
	}
}

func (f *SyncFile) addSource(id string) {

	for _, s := range f.Sources {

		if id == s {
			return
		}
	}

	f.Sources = append(f.Sources, id)
}

func (f *SyncFile) DeleteSource(id string) {

	mut.Lock()
	defer mut.Unlock()

	var newList []string

	for _, s := range f.Sources {

		if s != id {
			newList = append(newList, s)
		}
	}

	f.Sources = newList
}

func (f *SyncFile) CleanSources(servers []serverinfo.ServerInfo) {

	mut.Lock()
	defer mut.Unlock()

	var newList []string

	for _, src := range f.Sources {
		for _, server := range servers {

			if src == server.ID {
				newList = append(newList, src)
			}
		}
	}

	f.Sources = newList
}

func (f *SyncFile) GetSources() []string {

	mut.Lock()
	defer mut.Unlock()

	return f.Sources
}

func (f *SyncFile) SelectSource(thisID string) (string, error) {

	mut.Lock()
	defer mut.Unlock()

	l := len(f.Sources)

	if l == 1 {
		// if only one source in list
		if f.Sources[0] != thisID {
			return f.Sources[0], nil
		}
	} else if l > 1 {

		// select randomly
		for range f.Sources {

			index := rand.Intn(l)

			if f.Sources[index] != thisID {
				return f.Sources[index], nil
			}
		}

		// if random fails
		for _, src := range f.Sources {

			if src != thisID {
				return src, nil
			}
		}
	}

	return "", errors.New("no source found")
}

func (f *SyncFile) SetDestination(id string) {

	mut.Lock()
	defer mut.Unlock()

	f.destination = id
}

func (f *SyncFile) GetDestination() string {

	mut.Lock()
	defer mut.Unlock()

	return f.destination
}