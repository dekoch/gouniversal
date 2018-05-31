package userToken

import (
	"sync"

	"github.com/google/uuid"
)

type uToken struct {
	uid   string
	token string
}

type UserToken struct {
	mut    sync.Mutex
	tokens []uToken
}

// New returns and adds a unique token
func (t *UserToken) New(uid string) string {

	t.mut.Lock()
	defer t.mut.Unlock()

	u := uuid.Must(uuid.NewRandom())
	ut := u.String()

	for i := 0; i < len(t.tokens); i++ {

		if uid == t.tokens[i].uid {
			t.tokens[i].token = ut
			return ut
		}
	}

	newToken := make([]uToken, 1)
	newToken[0].uid = uid
	newToken[0].token = ut

	t.tokens = append(t.tokens, newToken...)

	return ut
}

// Check returns true, if the UUID and token match with items inside list
func (t *UserToken) Check(uid string, ut string) bool {

	t.mut.Lock()
	defer t.mut.Unlock()

	for i := 0; i < len(t.tokens); i++ {

		if uid == t.tokens[i].uid {

			if ut == t.tokens[i].token {
				return true
			}

			return false
		}
	}

	return false
}

// Remove removes UUID and token from list
func (t *UserToken) Remove(uid string) {

	t.mut.Lock()
	defer t.mut.Unlock()

	var l []uToken
	n := make([]uToken, 1)

	for i := 0; i < len(t.tokens); i++ {

		if uid != t.tokens[i].uid {

			n[0] = t.tokens[i]

			l = append(l, n...)
		}
	}

	t.tokens = l
}
