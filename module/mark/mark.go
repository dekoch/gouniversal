package mark

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/shared/aes"
	"github.com/dekoch/gouniversal/shared/console"
)

const sender = "mark v1.1"
const testtime = 30.0 // seconds

var (
	tStart time.Time
	tEnd   time.Time
)

func LoadConfig() {

	console.Log("start", sender)

	singleAesEnc()
	multiAesEnc()
}

func singleAesEnc() {

	key, err := aes.NewKey(32)
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		c   int
		mut sync.Mutex
	)

	tStart = time.Now()

	for elapsed() == false {

		aesEnc(key)

		mut.Lock()
		c++
		mut.Unlock()
	}

	tEnd = time.Now()

	console.Log("AES single enc"+log(c), sender)
}

func multiAesEnc() {

	key, err := aes.NewKey(32)
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		wg  sync.WaitGroup
		c   int
		mut sync.Mutex
	)

	tStart = time.Now()

	for i := 0; i < runtime.NumCPU(); i++ {

		wg.Add(1)

		go func() {

			for elapsed() == false {

				aesEnc(key)

				mut.Lock()
				c++
				mut.Unlock()
			}

			wg.Done()
		}()
	}

	wg.Wait()

	tEnd = time.Now()

	console.Log("AES multi enc"+log(c), sender)
}

func aesEnc(key []byte) {

	_, err := aes.Encrypt(key, string(key))
	if err != nil {
		fmt.Println(err)
	}
}

func elapsed() bool {

	elapsed := time.Now().Sub(tStart)
	f := elapsed.Seconds() * 1000.0

	if f < testtime*1000.0 {
		return false
	}

	return true
}

func log(c int) string {

	elapsed := tEnd.Sub(tStart)
	t := elapsed.Seconds()

	count := float64(c) / 1000.0
	count = count / t

	return "\t" + strconv.FormatFloat(count, 'f', 3, 64) + " KC/s"
}

/*
Raspberry Pi B OpenWRT
arm5
2018/12/29 13:09:38 mark: single enc	c:166992 t:30000.2ms
2018/12/29 13:10:08 mark: multi enc		c:164590 t:30003.0ms
arm6
2018/12/29 13:08:00 mark: single enc	c:193839 t:30000.1ms
2018/12/29 13:08:30 mark: multi enc		c:192503 t:30002.4ms

Raspberry Pi B buildroot
arm5
1970/01/01 00:02:05 mark: single enc	c:159194 t:30000.1ms
1970/01/01 00:02:35 mark: multi enc		c:157103 t:30002.7ms
arm6
1970/01/01 00:05:35 mark: single enc	c:182115 t:30000.1ms
1970/01/01 00:06:05 mark: multi enc		c:182365 t:30002.2ms

Raspberry Pi 3 B OpenWRT
aarch64
2018/12/28 20:57:08 mark: single enc	c:1083464 t:30000.0ms
2018/12/28 20:57:38 mark: multi enc		c:4691490 t:30000.4ms

NanoPi NEO OpenWRT
arm5
2019/01/01 15:37:21 mark: single enc	c:598780 t:30000.0ms
2019/01/01 15:37:51 mark: multi enc		c:2174463 t:30000.4ms
arm6
2019/01/01 15:39:24 mark: single enc	c:729607 t:30000.0ms
2019/01/01 15:39:54 mark: multi enc		c:2572255 t:30000.3ms

i7-7500U debian
amd64
2018/12/28 22:17:20 mark: single enc	c:17078423 t:30000.0ms
2018/12/28 22:17:50 mark: multi enc		c:32936586 t:30000.0ms
i386
2019/01/01 16:53:42 mark: single enc	c:7954930 t:30000.0ms
2019/01/01 16:54:12 mark: multi enc		c:18893605 t:30000.0ms
*/
