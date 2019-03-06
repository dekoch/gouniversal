package global

import (
	"github.com/dekoch/gouniversal/program/groupconfig"
	"github.com/dekoch/gouniversal/program/guestmanagement"
	"github.com/dekoch/gouniversal/program/uiconfig"
	"github.com/dekoch/gouniversal/program/userconfig"
	"github.com/dekoch/gouniversal/shared/language"
)

type Global struct{}

var (
	UIConfig    uiconfig.UIConfig
	UserConfig  userconfig.UserConfig
	GroupConfig groupconfig.GroupConfig
	Lang        language.Language
	Guests      guestmanagement.GuestManagement
)
