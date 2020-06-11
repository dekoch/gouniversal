package instapostpage

import (
	"encoding/json"
	"errors"

	"github.com/dekoch/gouniversal/shared/api/instaclient"
	"github.com/dekoch/gouniversal/shared/api/instapostpage/query"
	"github.com/dekoch/gouniversal/shared/htmlextract"
)

type Response struct {
	EntryData   EntryData `json:"entry_data"`
	RolloutHash string    `json:"rollout_hash"`
}

type EntryData struct {
	PostPage []PostPage `json:"PostPage"`
}

type PostPage struct {
	GraphQL GraphQL `json:"graphql"`
}

type GraphQL struct {
	ShortcodeMedia ShortcodeMedia `json:"shortcode_media"`
}

type ShortcodeMedia struct {
	ID          string      `json:"id"`
	DisplayURL  string      `json:"display_url"`
	IsVideo     bool        `json:"is_video"`
	VideoURL    string      `json:"video_url"`
	EstChildren EstChildren `json:"edge_sidecar_to_children"`
}

// EstChildren edge_sidecar_to_children
type EstChildren struct {
	Edges []ChildrenEdge `json:"edges"`
}

type ChildrenEdge struct {
	Node ChildrenNode `json:"node"`
}

type ChildrenNode struct {
	ID         string `json:"id"`
	DisplayURL string `json:"display_url"`
	IsVideo    bool   `json:"is_video"`
	VideoURL   string `json:"video_url"`
}

func GetResponse(id string, ic *instaclient.InstaClient) (Response, error) {

	var (
		err error
		ret Response
		b   []byte
	)

	for i := 0; i <= 2; i++ {

		switch i {
		case 0:
			b, err = query.Send(id, ic)

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
