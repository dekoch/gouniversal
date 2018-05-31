package respDevConfig

import (
	"time"

	"github.com/dekoch/gouniversal/shared/functions"
)

type RespDevConfig struct {
	Lastseen        time.Time
	SetIntvl        float64
	Intvl           float64
	IntvlAdjustStep float64
}

func (c *RespDevConfig) Init() {
	c.Lastseen = time.Now()
	c.SetIntvl = 3.0
	c.Intvl = 3.0
	c.IntvlAdjustStep = 0.1
}

func (c *RespDevConfig) SetInterval(sec float64, prc float64) {

	c.SetIntvl = sec
	c.Intvl = sec

	c.IntvlAdjustStep = prc
}

func (c *RespDevConfig) Interval() float64 {

	sec := time.Now().Sub(c.Lastseen)
	c.Lastseen = time.Now()
	dif := sec.Seconds()

	//console.Output(dif, "")

	if dif < c.SetIntvl {

		c.Intvl += c.IntvlAdjustStep

	} else {

		c.Intvl -= c.IntvlAdjustStep
	}

	c.Intvl = functions.Round(c.Intvl, 0.0, 3)

	return c.Intvl
}
