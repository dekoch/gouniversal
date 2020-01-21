package imgcache

import (
	"errors"
	"image"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/dbstorage"
	"github.com/dekoch/gouniversal/module/monmotion/mdimg"
	"github.com/dekoch/gouniversal/shared/console"
)

type ImgCache struct {
	device string

	ramImages []mdimg.MDImage
	ramMaxAge time.Duration

	dbImages    []mdimg.MDImage
	dbMaxAge    time.Duration
	dbBlockSize int

	mut     sync.RWMutex
	mutSave sync.Mutex
}

func (ic *ImgCache) LoadConfig(device string) error {

	ic.mut.Lock()
	defer ic.mut.Unlock()

	ic.device = device

	err := dbstorage.LoadConfig()
	if err != nil {
		return err
	}

	return dbstorage.DeleteImages(ic.device, mdimg.CACHE, time.Now().AddDate(-999, 0, 0), time.Now())
}

func (ic *ImgCache) Exit() error {

	ic.mut.Lock()
	defer ic.mut.Unlock()

	ic.ramImages = []mdimg.MDImage{}
	ic.dbImages = []mdimg.MDImage{}

	return dbstorage.DeleteImages(ic.device, mdimg.CACHE, time.Now().AddDate(-999, 0, 0), time.Now())
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

func (ic *ImgCache) AddImage(img mdimg.MDImage, todb bool) {

	ic.mut.Lock()
	defer ic.mut.Unlock()

	ic.addToCache(img)

	if todb {
		ic.addToDB(img)
	}
}

func (ic *ImgCache) addToCache(img mdimg.MDImage) {

	var n []mdimg.MDImage
	n = append(n, img)

	for i := range ic.ramImages {

		if time.Since(ic.ramImages[i].Captured) > ic.ramMaxAge {

			if i >= 1 {
				ic.ramImages = append(n, ic.ramImages[:i-1]...)
				return
			}
		}
	}

	ic.ramImages = append(n, ic.ramImages...)
}

func (ic *ImgCache) addToDB(img mdimg.MDImage) {

	if len(ic.dbImages) >= ic.dbBlockSize {

		go func(images []mdimg.MDImage) {

			ic.saveImagesToDB(images, mdimg.CACHE)
		}(ic.dbImages)

		var n []mdimg.MDImage
		ic.dbImages = append(n, img)
		return
	}

	ic.dbImages = append(ic.dbImages, img)
}

func (ic *ImgCache) GetImageCnt() int {

	ic.mut.RLock()
	defer ic.mut.RUnlock()

	return len(ic.ramImages)
}

func (ic *ImgCache) GetLatestImage() (mdimg.MDImage, error) {

	ic.mut.RLock()
	defer ic.mut.RUnlock()

	if len(ic.ramImages) == 0 {
		var nw mdimg.MDImage
		nw.Captured = time.Now()

		upLeft := image.Point{0, 0}
		lowRight := image.Point{100, 100}

		nw.EncodeImage(image.NewRGBA(image.Rectangle{upLeft, lowRight}))

		return nw, nil
	}

	return ic.ramImages[0], nil
}

func (ic *ImgCache) GetOldImage(d time.Duration) (mdimg.MDImage, error) {

	ic.mut.RLock()
	defer ic.mut.RUnlock()

	if len(ic.ramImages) == 0 {
		var nw mdimg.MDImage
		return nw, errors.New("no image available")
	}

	for i := range ic.ramImages {

		if time.Since(ic.ramImages[i].Captured) > d {
			return ic.ramImages[i], nil
		}
	}

	return ic.ramImages[0], nil
}

func (ic *ImgCache) GetFPS() float32 {

	ic.mut.RLock()
	defer ic.mut.RUnlock()

	return ic.getFPS()
}

func (ic *ImgCache) getFPS() float32 {

	l := len(ic.ramImages)

	if l <= 1 {
		return 0.0
	}

	t := ic.ramImages[0].Captured.Sub(ic.ramImages[l-1].Captured).Milliseconds()

	fps := float32(l) / float32(t) * 1000.0

	return fps
}

func (ic *ImgCache) SaveImages(trigger time.Time, prerecoding, overrun time.Duration) error {

	ic.mut.Lock()
	defer ic.mut.Unlock()

	console.Output("saving images", "MonMotion")

	go func(images []mdimg.MDImage, trigger time.Time, prerecoding, overrun time.Duration) {

		ic.saveImagesToDB(images, mdimg.CACHE)
		dbstorage.SetStateToImages(ic.device, mdimg.SAVED, trigger.Add(-prerecoding), trigger.Add(overrun))
	}(ic.dbImages, trigger, prerecoding, overrun)

	ic.dbImages = []mdimg.MDImage{}

	return nil
}

func (ic *ImgCache) saveImagesToDB(images []mdimg.MDImage, state mdimg.ImageState) error {

	if len(images) == 0 {
		return errors.New("no image available")
	}

	ic.mutSave.Lock()
	defer ic.mutSave.Unlock()

	t := images[len(images)-1].Captured.Add(-ic.dbMaxAge)
	dbstorage.DeleteImages(ic.device, mdimg.CACHE, time.Now().AddDate(-999, 0, 0), t)

	for i := range images {

		images[i].Device = ic.device
		images[i].State = state
	}

	return dbstorage.SaveImages(images)
}
