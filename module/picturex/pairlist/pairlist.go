package pairlist

import (
	"errors"
	"sort"
	"sync"

	"github.com/dekoch/gouniversal/module/picturex/pair"
)

type PairList struct {
	Pairs    []pair.Pair
	maxPairs int
}

var mut sync.RWMutex

func SetSourcePath(p string) {
	pair.SetSourcePath(p)
}

func SetDestinationPath(p string) {
	pair.SetDestinationPath(p)
}

func SetStaticPath(p string) {
	pair.SetStaticPath(p)
}

func (pl *PairList) SetMaxPairs(cnt int) {

	mut.Lock()
	defer mut.Unlock()

	pl.maxPairs = cnt
}

func (pl *PairList) NewPair(user string) (string, error) {

	mut.Lock()
	defer mut.Unlock()

	pairs := len(pl.Pairs)
	// if number of pairs exceeds maximum, remove the oldest
	if pairs > pl.maxPairs {

		sort.Slice(pl.Pairs, func(i, j int) bool { return pl.Pairs[i].GetTimeStamp().Unix() < pl.Pairs[j].GetTimeStamp().Unix() })

		err := pl.deletePair(pl.Pairs[pairs-1].GetUUID())
		if err != nil {
			return "", err
		}
	}

	p := pair.New(user)
	pl.Pairs = append(pl.Pairs, p)

	return p.GetUUID(), nil
}

func (pl *PairList) DeletePair(p string) error {

	mut.Lock()
	defer mut.Unlock()

	return pl.deletePair(p)
}

func (pl *PairList) deletePair(p string) error {

	var l []pair.Pair

	for i := 0; i < len(pl.Pairs); i++ {

		if p == pl.Pairs[i].GetUUID() {

			err := pl.Pairs[i].Delete()
			if err != nil {
				return err
			}
		} else {
			l = append(l, pl.Pairs[i])
		}
	}

	pl.Pairs = l

	return nil
}

func (pl *PairList) GetFirstPairFromUser(user string) (string, error) {

	mut.RLock()
	defer mut.RUnlock()

	sort.Slice(pl.Pairs, func(i, j int) bool { return pl.Pairs[i].GetTimeStamp().Unix() > pl.Pairs[j].GetTimeStamp().Unix() })

	for i := 0; i < len(pl.Pairs); i++ {

		if pl.Pairs[i].IsFirstUser(user) ||
			pl.Pairs[i].IsSecondUser(user) {

			return pl.Pairs[i].GetUUID(), nil
		}
	}

	return "", errors.New("pair not found")
}

func (pl *PairList) GetPairsFromUser(user string) ([]string, error) {

	mut.RLock()
	defer mut.RUnlock()

	var l []string

	for i := 0; i < len(pl.Pairs); i++ {

		if pl.Pairs[i].IsFirstUser(user) ||
			pl.Pairs[i].IsSecondUser(user) {

			l = append(l, pl.Pairs[i].GetUUID())
		}
	}

	if len(l) > 0 {
		return l, nil
	}

	return l, errors.New("no pair found")
}

func (pl *PairList) IsFirstUser(p, user string) (bool, error) {

	mut.RLock()
	defer mut.RUnlock()

	for i := 0; i < len(pl.Pairs); i++ {

		if p == pl.Pairs[i].GetUUID() {
			return pl.Pairs[i].IsFirstUser(user), nil
		}
	}

	return false, errors.New("pair not found")
}

func (pl *PairList) GetFirstUser(p string) (string, error) {

	mut.RLock()
	defer mut.RUnlock()

	for i := 0; i < len(pl.Pairs); i++ {

		if p == pl.Pairs[i].GetUUID() {
			return pl.Pairs[i].GetFirstUser(), nil
		}
	}

	return "", errors.New("pair not found")
}

func (pl *PairList) SetSecondUser(p, user string) error {

	mut.Lock()
	defer mut.Unlock()

	for i := 0; i < len(pl.Pairs); i++ {

		if p == pl.Pairs[i].GetUUID() {
			return pl.Pairs[i].SetSecondUser(user)
		}
	}

	return errors.New("pair not found")
}

func (pl *PairList) GetSecondUser(p string) (string, error) {

	mut.RLock()
	defer mut.RUnlock()

	for i := 0; i < len(pl.Pairs); i++ {

		if p == pl.Pairs[i].GetUUID() {
			return pl.Pairs[i].GetSecondUser(), nil
		}
	}

	return "", errors.New("pair not found")
}

func (pl *PairList) SetPicture(p, user, name string) error {

	mut.Lock()
	defer mut.Unlock()

	for i := 0; i < len(pl.Pairs); i++ {

		if p == pl.Pairs[i].GetUUID() {
			return pl.Pairs[i].SetPicture(user, name)
		}
	}

	return errors.New("pair not found")
}

func (pl *PairList) GetFirstPicture(p, user string) (string, error) {

	mut.Lock()
	defer mut.Unlock()

	for i := 0; i < len(pl.Pairs); i++ {

		if p == pl.Pairs[i].GetUUID() {
			return pl.Pairs[i].GetFirstPicture(user)
		}
	}

	return "", errors.New("pair not found")
}

func (pl *PairList) GetSecondPicture(p, user string) (string, error) {

	mut.Lock()
	defer mut.Unlock()

	for i := 0; i < len(pl.Pairs); i++ {

		if p == pl.Pairs[i].GetUUID() {
			return pl.Pairs[i].GetSecondPicture(user)
		}
	}

	return "", errors.New("pair not found")
}

func (pl *PairList) UnlockPicture(p, user string) error {

	mut.Lock()
	defer mut.Unlock()

	for i := 0; i < len(pl.Pairs); i++ {

		if p == pl.Pairs[i].GetUUID() {
			return pl.Pairs[i].UnlockPicture(user)
		}
	}

	return errors.New("pair not found")
}
