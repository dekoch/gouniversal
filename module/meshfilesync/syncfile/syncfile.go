package syncfile

import (
	"bytes"
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
	incoming    time.Time
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
	_, err := f.update(root)

	return err
}

func (f *SyncFile) Update(root string) (bool, error) {

	mut.Lock()
	defer mut.Unlock()

	return f.update(root)
}

func (f *SyncFile) update(root string) (bool, error) {

	var updated bool

	sum, err := file.Checksum(root + f.Path)
	if err != nil {
		return updated, err
	}

	if bytes.Equal(f.Checksum, sum) == false {
		updated = true
		f.Checksum = sum
	}

	fileInfo, err := os.Stat(root + f.Path)
	if err != nil {
		return updated, err
	}

	if f.Size != fileInfo.Size() {
		updated = true
		f.Size = fileInfo.Size()
	}

	if f.ModTime != fileInfo.ModTime() {
		updated = true
		f.ModTime = fileInfo.ModTime()
	}

	if updated {
		f.clearSources()
	}

	return updated, nil
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

func (f *SyncFile) ClearSources() {

	mut.Lock()
	defer mut.Unlock()

	f.clearSources()
}

func (f *SyncFile) clearSources() {

	var n []string
	f.Sources = n
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

func (f *SyncFile) ClearChecksum() {

	mut.Lock()
	defer mut.Unlock()

	var n []byte
	f.Checksum = n
}

func (f *SyncFile) SetIncomingTime(t time.Time) {

	mut.Lock()
	defer mut.Unlock()

	f.incoming = t
}

func (f *SyncFile) GetIncomingTime() time.Time {

	mut.Lock()
	defer mut.Unlock()

	return f.incoming
}
