package hashstor

import (
	"errors"
	"sort"
	"sync"
	"time"
)

type HashStor struct {
	mut    sync.RWMutex
	Hashes []Hash
}

type Hash struct {
	Hash      string
	Expired   time.Time
	requested time.Time
}

func (hs *HashStor) Init() {

	hs.mut.Lock()
	defer hs.mut.Unlock()

	t := time.Now()

	for i := range hs.Hashes {

		hs.Hashes[i].requested = t
	}
}

func (hs *HashStor) Add(str string) {

	hs.mut.Lock()
	defer hs.mut.Unlock()

	for i := range hs.Hashes {

		if hs.Hashes[i].Hash == str {
			return
		}
	}

	var n Hash
	n.Hash = str

	hs.Hashes = append(hs.Hashes, n)
}

func (hs *HashStor) Remove(str string) {

	hs.mut.Lock()
	defer hs.mut.Unlock()

	var l []Hash

	for i := range hs.Hashes {

		if str != hs.Hashes[i].Hash {
			l = append(l, hs.Hashes[i])
		}
	}

	hs.Hashes = l
}

func (hs *HashStor) RemoveAll() {

	hs.mut.Lock()
	defer hs.mut.Unlock()

	var l []Hash
	hs.Hashes = l
}

func (hs *HashStor) SetAsExpired(str string, to time.Duration) {

	hs.mut.Lock()
	defer hs.mut.Unlock()

	for i := range hs.Hashes {

		if hs.Hashes[i].Hash == str {

			hs.Hashes[i].Expired = time.Now().Add(to)
			return
		}
	}
}

func (hs *HashStor) GetHash() (string, error) {

	hs.mut.Lock()
	defer hs.mut.Unlock()

	sort.Slice(hs.Hashes, func(i, j int) bool { return hs.Hashes[i].requested.Unix() < hs.Hashes[j].requested.Unix() })

	for i := range hs.Hashes {

		if time.Now().After(hs.Hashes[i].Expired) {

			hs.Hashes[i].requested = time.Now()

			return hs.Hashes[i].Hash, nil
		}
	}

	return "", errors.New("no hash found")
}
