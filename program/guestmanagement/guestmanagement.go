package guestmanagement

import (
	"sort"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/program/userconfig"

	"github.com/google/uuid"
)

type guestUser struct {
	user          userconfig.User
	loginAttempts int
	timeStamp     time.Time
}

type GuestManagement struct {
	user []guestUser
}

var mut sync.Mutex

func (c *GuestManagement) NewGuest(user userconfig.User, maxuser int) userconfig.User {

	var newGuest guestUser
	newGuest.user = user
	newGuest.loginAttempts = 0

	// set new uuid for guest
	u := uuid.Must(uuid.NewRandom())
	newGuest.user.UUID = u.String()

	newGuest.timeStamp = time.Now()

	// add new guest to list
	mut.Lock()
	defer mut.Unlock()

	guests := len(c.user)
	// if number of guests exceeds maximum, remove the oldest
	if guests > maxuser {

		sort.Slice(c.user, func(i, j int) bool { return c.user[i].timeStamp.Unix() < c.user[j].timeStamp.Unix() })

		c.user = c.user[1:guests]
	}

	c.user = append(c.user, newGuest)

	return newGuest.user
}

func (c *GuestManagement) SelectGuest(uid string) userconfig.User {

	mut.Lock()
	defer mut.Unlock()

	for u := 0; u < len(c.user); u++ {

		// search guest with UUID
		if uid == c.user[u].user.UUID {

			c.user[u].timeStamp = time.Now()

			return c.user[u].user
		}
	}

	var user userconfig.User
	user.State = -1
	return user
}

func (c *GuestManagement) MaxLoginAttempts(uid string, maxattempts int) bool {

	mut.Lock()
	defer mut.Unlock()

	for u := 0; u < len(c.user); u++ {

		// search guest with UUID
		if uid == c.user[u].user.UUID {

			c.user[u].loginAttempts++

			if c.user[u].loginAttempts > maxattempts {
				return true
			}

			return false
		}
	}

	// if another user is already logged in (no guest)
	return false
}
