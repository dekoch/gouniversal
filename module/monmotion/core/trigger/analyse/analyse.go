package analyse

import (
	"runtime"
	"strconv"
	"sync"

	"github.com/dekoch/gouniversal/module/monmotion/typemd"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/timeout"
)

type Analyse struct {
	wg            sync.WaitGroup
	chanResult    chan Result
	tuneEnabled   bool
	tuneThreshold uint32
	tuneTimeOut   timeout.TimeOut
	tuneStep      uint32
}

type work struct {
	startX, startY, endX, endY int
	threshold                  uint32
}

type Result struct {
	Px, Rpx, Gpx, Bpx uint
	Threshold         uint32
}

var mut sync.RWMutex

func (an *Analyse) LoadConfig() {

}

func (an *Analyse) EnableAutoTune(timeout int, tunestep uint32) {

	mut.Lock()
	defer mut.Unlock()

	an.tuneEnabled = true
	an.tuneThreshold = 0
	an.tuneStep = tunestep
	an.tuneTimeOut.Start(timeout)
}

func (an *Analyse) AnalyseImage(oldImg, newImg *typemd.MoImage, threshold uint32) (Result, error) {

	mut.Lock()
	defer mut.Unlock()

	var (
		err error
		ret Result
	)

	workerCnt := runtime.NumCPU()
	//workerCnt := 1

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				if newImg.Img.Bounds().Max.X != oldImg.Img.Bounds().Max.X ||
					newImg.Img.Bounds().Max.Y != oldImg.Img.Bounds().Max.Y {

					return
				}

			case 1:
				an.chanResult = make(chan Result, workerCnt)

				bounds := newImg.Img.Bounds().Max
				partY := bounds.Y / workerCnt

				for n := 0; n < workerCnt; n++ {

					var w work
					w.startX = 0
					w.startY = n * partY
					w.endX = bounds.X
					w.endY = (n + 1) * partY

					if an.tuneEnabled {
						threshold = an.tuneThreshold
					}

					w.threshold = threshold

					an.wg.Add(1)

					go an.worker(oldImg, newImg, w)
				}

			case 2:
				an.wg.Wait()

				close(an.chanResult)

				for r := range an.chanResult {

					ret.Rpx += r.Rpx
					ret.Gpx += r.Gpx
					ret.Bpx += r.Bpx

					ret.Px += r.Px
				}

				if an.tuneEnabled {

					if ret.Px > 0 {

						an.tuneThreshold += an.tuneStep

						console.Output("autotuning threshold "+strconv.FormatUint(uint64(an.tuneThreshold), 10), "MonMotion")

						an.tuneTimeOut.Reset()

					} else if an.tuneTimeOut.Elapsed() {

						console.Output("autotuning finished", "MonMotion")

						an.tuneEnabled = false

						ret.Threshold = threshold
					}
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return ret, err
}

func (an *Analyse) worker(oldImg, newImg *typemd.MoImage, w work) {

	var r Result

	for y := w.startY; y < w.endY; y++ {

		for x := w.startX; x < w.endX; x++ {

			rNew, gNew, bNew, _ := newImg.Img.At(x, y).RGBA()
			rOld, gOld, bOld, _ := oldImg.Img.At(x, y).RGBA()

			if an.threshold(rNew, rOld, w.threshold) {

				r.Rpx++
				r.Px++
			}

			if an.threshold(gNew, gOld, w.threshold) {

				r.Gpx++
				r.Px++
			}

			if an.threshold(bNew, bOld, w.threshold) {

				r.Bpx++
				r.Px++
			}
		}
	}

	an.chanResult <- r

	an.wg.Done()
}

func (an *Analyse) threshold(new, old, thld uint32) bool {

	if new > old+thld {
		return true
	}

	if new+thld < old {
		return true
	}

	return false
}
