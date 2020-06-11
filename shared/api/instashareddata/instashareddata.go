package instashareddata

import (
	"encoding/json"
	"errors"

	"github.com/dekoch/gouniversal/shared/api/instaclient"
	"github.com/dekoch/gouniversal/shared/api/instashareddata/query"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/htmlextract"
)

type Response struct {
	Config      Config    `json:"config"`
	EntryData   EntryData `json:"entry_data"`
	RolloutHash string    `json:"rollout_hash"`
}

type Config struct {
	CsrfToken string `json:"csrf_token"`
}

type EntryData struct {
	ProfilePage []ProfilePage `json:"ProfilePage"`
}

type ProfilePage struct {
	GraphQL GraphQL `json:"graphql"`
}

type GraphQL struct {
	User User `json:"user"`
}

type User struct {
	EdgeFollowedBy    EdgeFollowedBy `json:"edge_followed_by"`
	EdgeFollow        EdgeFollow     `json:"edge_follow"`
	ID                string         `json:"id"`
	IsBusinessAccount bool           `json:"is_business_account"`
	IsPrivate         bool           `json:"is_private"`
	IsVerified        bool           `json:"is_verified"`
	UserName          string         `json:"username"`
	EottMedia         EottMedia      `json:"edge_owner_to_timeline_media"`
}

type EdgeFollowedBy struct {
	Count int `json:"count"`
}

type EdgeFollow struct {
	Count int `json:"count"`
}

// EottMedia edge_owner_to_timeline_media
type EottMedia struct {
	Count int `json:"count"`
}

func GetResponse(username string, ic *instaclient.InstaClient) (Response, error) {

	var (
		err error
		b   []byte
		ret Response
	)

	if functions.IsEmpty(username) {
		return ret, errors.New("username is empty")
	}

	for i := 0; i <= 2; i++ {

		switch i {
		case 0:
			b, err = query.Send(username, ic)

		case 1:
			var s string
			s, err = htmlextract.Extract(string(b), "<script type=\"text/javascript\">window._sharedData = ", ";</script>")
			b = []byte(s)

		case 2:
			err = json.Unmarshal(b, &ret)
		}

		if err != nil {
			return ret, errors.New(err.Error() + " \"" + string(b) + "\"")
		}
	}

	return ret, nil
}
