package iptracker

import (
	"github.com/dekoch/gouniversal/module/iptracker/global"
	"github.com/dekoch/gouniversal/module/iptracker/tracker"
)

func LoadConfig() {

	global.Config.LoadConfig()

	tracker.LoadConfig()
}

func Exit() {

	global.Config.SaveConfig()
}
