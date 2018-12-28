package mark

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/shared/aes"
	"github.com/dekoch/gouniversal/shared/console"
)

const sender = "mark"
const testtime = 30.0 // seconds

var (
	tStart time.Time
	tEnd   time.Time
)

func LoadConfig() {

	console.Log("start", sender)

	singleEnc()
	multiEnc()
}

func singleEnc() {

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

		enc(key)

		mut.Lock()
		c++
		mut.Unlock()
	}

	tEnd = time.Now()

	console.Log("single enc\tc:"+strconv.Itoa(c)+" t:"+getTime(), sender)
}

func multiEnc() {

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

	for i := 0; i < 20; i++ {

		wg.Add(1)

		go func() {

			for elapsed() == false {

				enc(key)

				mut.Lock()
				c++
				mut.Unlock()
			}

			wg.Done()
		}()
	}

	wg.Wait()

	tEnd = time.Now()

	console.Log("multi enc\tc:"+strconv.Itoa(c)+" t:"+getTime(), sender)
}

func enc(key []byte) {

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

func getTime() string {

	elapsed := tEnd.Sub(tStart)
	f := elapsed.Seconds() * 1000.0
	return strconv.FormatFloat(f, 'f', 1, 64) + "ms"
}

/*
RPi B OpenWRT
arm5
2018/12/28 21:02:03 mark: single enc	c:140089 t:30000.1ms
2018/12/28 21:02:33 mark: multi enc		c:136903 t:30002.7ms
arm6
2018/12/28 21:00:30 mark: single enc	c:166846 t:30000.2ms
2018/12/28 21:01:00 mark: multi enc		c:163833 t:30002.1ms

RPi 3 B OpenWRT
aarch64
2018/12/28 20:57:08 mark: single enc	c:1083464 t:30000.0ms
2018/12/28 20:57:38 mark: multi enc		c:4691490 t:30000.4ms

i7-7500U debian
amd64
2018/12/28 22:17:20 mark: single enc	c:17078423 t:30000.0ms
2018/12/28 22:17:50 mark: multi enc		c:32936586 t:30000.0ms
*/
