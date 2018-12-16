package token

import (
	"sync"

	"github.com/dekoch/gouniversal/shared/functions"

	"github.com/google/uuid"
)

type tokenContent struct {
	uid   string
	token string
}

type Token struct {
	mut    sync.RWMutex
	tokens []tokenContent
}

// New returns and adds a unique token
func (t *Token) New(uid string) string {

	if functions.IsEmpty(uid) {
		return ""
	}

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

	var newToken tokenContent
	newToken.uid = uid
	newToken.token = ut

	t.tokens = append(t.tokens, newToken)

	return ut
}

// Check returns true, if the UUID and token match with items inside list
func (t *Token) Check(uid string, ut string) bool {

	if functions.IsEmpty(uid) ||
		functions.IsEmpty(ut) {
		return false
	}

	t.mut.RLock()
	defer t.mut.RUnlock()

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
func (t *Token) Remove(uid string) {

	if functions.IsEmpty(uid) {
		return
	}

	t.mut.Lock()
	defer t.mut.Unlock()

	var l []tokenContent

	for i := 0; i < len(t.tokens); i++ {

		if uid != t.tokens[i].uid {

			l = append(l, t.tokens[i])
		}
	}

	t.tokens = l
}
