package userToken

import "github.com/google/uuid"

type uToken struct {
	Uid   string
	Token string
}

type UserToken struct {
	Tokens []uToken
}

func (t *UserToken) New(uid string) string {

	u := uuid.Must(uuid.NewRandom())
	ut := u.String()

	for i := 0; i < len(t.Tokens); i++ {

		if uid == t.Tokens[i].Uid {
			t.Tokens[i].Token = ut
			return ut
		}
	}

	newToken := make([]uToken, 1)
	newToken[0].Uid = uid
	newToken[0].Token = ut

	t.Tokens = append(newToken, t.Tokens...)

	return ut
}

func (t *UserToken) Check(uid string, ut string) bool {

	for i := 0; i < len(t.Tokens); i++ {

		if uid == t.Tokens[i].Uid {

			if ut == t.Tokens[i].Token {
				return true
			}

			return false
		}
	}

	return false
}
