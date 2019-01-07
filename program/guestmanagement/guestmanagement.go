package guestmanagement

import (
	"sync"

	"github.com/dekoch/gouniversal/program/global"
	"github.com/dekoch/gouniversal/program/userconfig"
	"github.com/dekoch/gouniversal/shared/console"

	"github.com/google/uuid"
)

type guestUser struct {
	User          userconfig.User
	LoginAttempts int
}

type GuestManagement struct {
	User []guestUser
}

var mut sync.RWMutex

func (c *GuestManagement) NewGuest() userconfig.User {

	newGuest := make([]guestUser, 1)

	// search for public user
	var err error
	newGuest[0].User, err = global.UserConfig.GetWithState(0)
	if err != nil {
		console.Log(err, "")
	}

	newGuest[0].LoginAttempts = 0

	// set new uuid for guest
	u := uuid.Must(uuid.NewRandom())
	newGuest[0].User.UUID = u.String()

	// add new guest to list
	mut.Lock()
	defer mut.Unlock()

	guests := len(c.User)
	// if number of guests exceeds maximum, remove the oldest
	if guests > global.UIConfig.MaxGuests {
		c.User = c.User[1:guests]
	}

	c.User = append(c.User, newGuest...)

	return newGuest[0].User
}

func (c *GuestManagement) SelectGuest(uid string) userconfig.User {

	mut.RLock()
	defer mut.RUnlock()

	for u := 0; u < len(c.User); u++ {

		// search guest with UUID
		if uid == c.User[u].User.UUID {

			return c.User[u].User
		}
	}

	var user userconfig.User
	user.State = -1
	return user
}

func (c *GuestManagement) MaxLoginAttempts(uid string) bool {

	mut.Lock()
	defer mut.Unlock()

	for u := 0; u < len(c.User); u++ {

		// search guest with UUID
		if uid == c.User[u].User.UUID {

			c.User[u].LoginAttempts++

			if c.User[u].LoginAttempts > global.UIConfig.MaxLoginAttempts {
				return true
			}

			return false
		}
	}

	// if another user is already logged in (no guest)
	return false
}