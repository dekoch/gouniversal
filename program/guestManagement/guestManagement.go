package guestManagement

import (
	"fmt"
	"sync"

	"github.com/dekoch/gouniversal/program/global"
	"github.com/dekoch/gouniversal/program/userConfig"

	"github.com/google/uuid"
)

type guestUser struct {
	User          userConfig.User
	LoginAttempts int
}

type GuestManagement struct {
	Mut  sync.Mutex
	User []guestUser
}

func (c *GuestManagement) NewGuest() userConfig.User {

	newGuest := make([]guestUser, 1)

	// search for public user
	var err error
	newGuest[0].User, err = global.UserConfig.GetWithState(0)
	if err != nil {
		fmt.Println(err)
	}

	newGuest[0].LoginAttempts = 0

	// set new uuid for guest
	u := uuid.Must(uuid.NewRandom())
	newGuest[0].User.UUID = u.String()

	// add new guest to list
	c.Mut.Lock()
	defer c.Mut.Unlock()

	guests := len(c.User)
	// if number of guests exceeds maximum, remove the oldest
	if guests > global.UiConfig.File.MaxGuests {
		c.User = c.User[0 : guests-1]
	}

	c.User = append(c.User, newGuest...)

	return newGuest[0].User
}

func (c *GuestManagement) SelectGuest(uid string) userConfig.User {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	for u := 0; u < len(c.User); u++ {

		// search guest with UUID
		if uid == c.User[u].User.UUID {

			return c.User[u].User
		}
	}

	var user userConfig.User
	user.State = -1
	return user
}

func (c *GuestManagement) MaxLoginAttempts(uid string) bool {

	c.Mut.Lock()
	defer c.Mut.Unlock()

	for u := 0; u < len(c.User); u++ {

		// search guest with UUID
		if uid == c.User[u].User.UUID {

			c.User[u].LoginAttempts++

			if c.User[u].LoginAttempts > global.UiConfig.File.MaxLoginAttempts {
				return true
			}

			return false
		}
	}

	return true
}
