package typemd

import (
	"image"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/lang"
)

type Page struct {
	Content string
	Lang    lang.LangFile
}

type MoImage struct {
	Img      image.Image
	Captured time.Time
}
