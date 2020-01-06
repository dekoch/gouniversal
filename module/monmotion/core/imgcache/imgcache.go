package imgcache

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"strconv"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/typemd"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/io/file"
)

type ImgCache struct {
	images []typemd.MoImage
	maxAge time.Duration
	mu     sync.RWMutex
}

func (ic *ImgCache) SetMaxAge(d time.Duration) {

	ic.mu.Lock()
	defer ic.mu.Unlock()

	ic.maxAge = d
}

func (ic *ImgCache) AddImage(img typemd.MoImage) {

	ic.mu.Lock()
	defer ic.mu.Unlock()

	var n []typemd.MoImage
	n = append(n, img)

	for i := range ic.images {

		if time.Since(ic.images[i].Captured) > ic.maxAge {

			if i >= 1 {

				ic.images = append(n, ic.images[:i-1]...)

				return
			}
		}
	}

	ic.images = append(n, ic.images...)
}

func (ic *ImgCache) Clear() {

	ic.mu.Lock()
	defer ic.mu.Unlock()

	var n []typemd.MoImage
	ic.images = n
}

func (ic *ImgCache) GetImageCnt() int {

	ic.mu.RLock()
	defer ic.mu.RUnlock()

	return len(ic.images)
}

func (ic *ImgCache) GetLatestImage() (typemd.MoImage, error) {

	ic.mu.RLock()
	defer ic.mu.RUnlock()

	if len(ic.images) == 0 {
		var nw typemd.MoImage
		nw.Captured = time.Now()

		upLeft := image.Point{0, 0}
		lowRight := image.Point{100, 100}

		nw.Img = image.NewRGBA(image.Rectangle{upLeft, lowRight})

		return nw, nil
	}

	return ic.images[0], nil
}

func (ic *ImgCache) GetOldImage(d time.Duration) (typemd.MoImage, error) {

	ic.mu.RLock()
	defer ic.mu.RUnlock()

	if len(ic.images) == 0 {
		var nw typemd.MoImage
		return nw, errors.New("no image available")
	}

	for i := range ic.images {

		if time.Since(ic.images[i].Captured) > d {
			return ic.images[i], nil
		}
	}

	return ic.images[0], nil
}

func (ic *ImgCache) GetFPS() float32 {

	ic.mu.RLock()
	defer ic.mu.RUnlock()

	return ic.getFPS()
}

func (ic *ImgCache) getFPS() float32 {

	l := len(ic.images)

	if l == 0 {
		return 0
	}

	t := time.Since(ic.images[l-1].Captured).Milliseconds()

	fps := float32(l) / float32(t) * 1000.0

	return fps
}

func (ic *ImgCache) SaveImages(path, name string) error {

	ic.mu.RLock()
	defer ic.mu.RUnlock()

	if len(ic.images) == 0 {
		return errors.New("no image available")
	}

	err := functions.CreateDir(path)
	if err != nil {
		return err
	}

	go saveImageCopy(path, name, ic.getFPS(), ic.images)

	return nil
}

func saveImageCopy(path, name string, fps float32, images []typemd.MoImage) error {

	if len(images) == 0 {
		return errors.New("no image available")
	}

	console.Output("saving images "+strconv.Itoa(len(images))+" (FPS:"+strconv.FormatFloat(float64(fps), 'f', 1, 32)+")", "MonMotion")

	for i := len(images) - 1; i >= 0; i-- {

		buf := &bytes.Buffer{}
		err := jpeg.Encode(buf, images[i].Img, nil)
		if err != nil {
			return err
		}

		err = file.WriteFile(path+images[i].Captured.Format("20060102_150405.0000")+"_"+name+".jpg", buf.Bytes())
		if err != nil {
			return err
		}
	}

	console.Output("images saved to "+path, "MonMotion")

	return nil
}
