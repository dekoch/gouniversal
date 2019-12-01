package typemd

import (
	"image"
	"time"
)

/*type Page struct {
	Content string
	Lang    lang.LangFile
}*/

type MoImage struct {
	Img      image.Image
	Captured time.Time
}
