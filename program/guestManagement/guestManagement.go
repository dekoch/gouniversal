package guestManagement

import (
	"gouniversal/program/global"
	"gouniversal/program/userConfig"
	"sync"

	"github.com/google/uuid"
)

type guestUser struct {
	User          userConfig.User
	LoginAttempts int
}

type guest struct {
	Mut  sync.Mutex
	User []guestUser
}

var Guest guest

func NewGuest() userConfig.User {

	newGuest := make([]guestUser, 1)

	// search for public user
	global.UserConfig.Mut.Lock()
	for u := 0; u < len(global.UserConfig.File.User); u++ {

		if global.UserConfig.File.User[u].State == 0 {
			newGuest[0].User = global.UserConfig.File.User[u]
			newGuest[0].LoginAttempts = 0
		}
	}
	global.UserConfig.Mut.Unlock()

	// set new uuid for guest
	u := uuid.Must(uuid.NewRandom())
	newGuest[0].User.UUID = u.String()

	// add new guest to list
	Guest.Mut.Lock()
	defer Guest.Mut.Unlock()

	guests := len(Guest.User)
	// if number of guests exceeds maximum, remove the oldest
	if guests > global.UiConfig.File.MaxGuests {
		Guest.User = Guest.User[0 : guests-1]
	}

	Guest.User = append(newGuest, Guest.User...)

	return newGuest[0].User
}

func SelectGuest(uid string) userConfig.User {

	Guest.Mut.Lock()
	defer Guest.Mut.Unlock()

	for u := 0; u < len(Guest.User); u++ {

		// search guest with UUID
		if uid == Guest.User[u].User.UUID {

			return Guest.User[u].User
		}
	}

	var user userConfig.User
	user.State = -1
	return user
}

func MaxLoginAttempts(uid string) bool {

	Guest.Mut.Lock()
	defer Guest.Mut.Unlock()

	for u := 0; u < len(Guest.User); u++ {

		// search guest with UUID
		if uid == Guest.User[u].User.UUID {

			Guest.User[u].LoginAttempts++

			if Guest.User[u].LoginAttempts > global.UiConfig.File.MaxLoginAttempts {
				return true
			}

			return false
		}
	}

	return true
}
