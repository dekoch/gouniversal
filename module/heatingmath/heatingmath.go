package heatingmath

import (
	"fmt"

	"github.com/dekoch/gouniversal/shared/aes"
)

func LoadConfig() {

	for i := 0; i <= 50; i++ {
		go calc()
	}
}

func calc() {

	for {
		key, err := aes.NewKey(32)
		if err != nil {
			fmt.Println(err)
			return
		}

		enc, err := aes.Encrypt(key, key)
		if err != nil {
			fmt.Println(err)
			return
		}

		dec, err := aes.Decrypt(key, enc)
		if err != nil {
			fmt.Println(err)
			return
		}

		dec = dec
	}
}
