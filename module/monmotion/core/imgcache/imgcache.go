package imgcache

import (
	"errors"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/dbcache"
	"github.com/dekoch/gouniversal/module/monmotion/dbstorage"
	"github.com/dekoch/gouniversal/module/monmotion/mdimg"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/sbool"
)

type ImgCache struct {
	device string

	ramMaxAge time.Duration

	dbMaxAge    time.Duration
	dbBlockSize int

	saving     sbool.Sbool
	dropFrames sbool.Sbool
	mut        sync.RWMutex
	mutSave    sync.RWMutex
}

func (ic *ImgCache) LoadConfig(device string) error {

	ic.mut.Lock()
	defer ic.mut.Unlock()

	ic.device = device

	return dbstorage.Stor.DeleteImages(ic.device, mdimg.CACHE, time.Now().AddDate(-999, 0, 0), time.Now())
}

func (ic *ImgCache) Exit() error {

	ic.mut.Lock()
	defer ic.mut.Unlock()

	return dbstorage.Stor.DeleteImages(ic.device, mdimg.CACHE, time.Now().AddDate(-999, 0, 0), time.Now())
}

func (ic *ImgCache) SetRAMSettings(maxage time.Duration) {

	ic.mut.Lock()
	defer ic.mut.Unlock()

	ic.ramMaxAge = maxage
}

func (ic *ImgCache) SetDBSettings(maxage time.Duration, blocksize int) {

	ic.mut.Lock()
	defer ic.mut.Unlock()

	ic.dbMaxAge = maxage
	ic.dbBlockSize = blocksize
}

func (ic *ImgCache) AddImage(img *mdimg.MDImage, todb bool) error {

	ic.mut.RLock()
	defer ic.mut.RUnlock()

	img.Device = ic.device
	img.State = mdimg.CACHE

	if ic.dropFrames.IsSet() && ic.saving.IsSet() && img.Trigger == false {
		return nil
	}

	err := dbcache.Cache.SaveImage(img)
	if err != nil {
		return err
	}

	if todb == false {

		toDate := img.Captured.Add(-ic.ramMaxAge * 2)
		return dbcache.Cache.DeleteImages(ic.device, mdimg.CACHE, time.Now().AddDate(-999, 0, 0), toDate)
	}

	if ic.saving.IsSet() {
		return nil
	}

	go func() {

		ic.mut.RLock()
		defer ic.mut.RUnlock()

		ids, err := dbcache.Cache.GetImageIDsWithState(ic.device, mdimg.CACHE)
		if err != nil {
			console.Log(err, "")
			return
		}

		if len(ids) >= ic.dbBlockSize {

			ic.saving.Set()
			defer ic.saving.UnSet()

			err := ic.saveBlock()
			if err != nil {
				console.Log(err, "")
				return
			}
		}
	}()

	return nil
}

func (ic *ImgCache) GetImageCnt() int {

	ic.mut.RLock()
	defer ic.mut.RUnlock()

	ids, err := dbcache.Cache.GetImageIDs(ic.device)
	if err != nil {
		return 0
	}

	return len(ids)
}

func (ic *ImgCache) GetLatestImage(img *mdimg.MDImage) error {

	ic.mut.RLock()
	defer ic.mut.RUnlock()

	var err error

	func() {

		var ids []string

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				ids, err = dbcache.Cache.GetImageIDs(ic.device)

			case 1:
				if len(ids) == 0 {
					err = errors.New("no image available")
					return
				}

			case 2:
				err = dbcache.Cache.LoadImage(ids[len(ids)-1], img)
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func (ic *ImgCache) GetOldImage(d time.Duration, img *mdimg.MDImage) error {

	ic.mut.RLock()
	defer ic.mut.RUnlock()

	var err error

	func() {

		var ids []string

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				ids, err = dbcache.Cache.GetImageIDsBetween(ic.device, time.Now().AddDate(-999, 0, 0), time.Now().Add(-d))

			case 1:
				if len(ids) == 0 {
					err = errors.New("no image available")
					return
				}

			case 2:
				err = dbcache.Cache.LoadImage(ids[len(ids)-1], img)
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func (ic *ImgCache) SaveImages(trigger time.Time, prerecoding, overrun time.Duration, keepallseq bool) error {

	console.Output("saving images", "MonMotion")

	if keepallseq == false {

		err := dbstorage.Stor.DeleteOldSequences(0)
		if err != nil {
			return err
		}
	}

	ic.saving.Set()
	defer ic.saving.UnSet()

	err := ic.saveBlock()
	if err != nil {
		return err
	}

	ic.mut.RLock()
	defer ic.mut.RUnlock()

	triggerID, err := dbstorage.Stor.GetIDByTime(ic.device, trigger)
	if err != nil {
		return err
	}

	return dbstorage.Stor.SetStateToSequence(triggerID, mdimg.SAVED)
}

func (ic *ImgCache) saveBlock() error {

	ic.mutSave.Lock()
	defer ic.mutSave.Unlock()

	ic.mut.RLock()
	defer ic.mut.RUnlock()

	ids, err := dbcache.Cache.GetImageIDsWithState(ic.device, mdimg.CACHE)
	if err != nil {
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	err = dbstorage.Stor.SaveBlock(ic.device, ids, ic.ramMaxAge, ic.dbMaxAge)
	if err != nil {
		return err
	}

	ids, err = dbcache.Cache.GetImageIDsWithState(ic.device, mdimg.CACHE)
	if err != nil {
		return err
	}

	if len(ids) > ic.dbBlockSize {
		ic.dropFrames.Set()
		console.Log("slow disk writing speed", "MonMotion")
	} else {
		ic.dropFrames.UnSet()
	}

	return nil
}
