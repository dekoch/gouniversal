package instatag

import (
	"encoding/json"
	"errors"

	"github.com/dekoch/gouniversal/shared/api/instaclient"
	"github.com/dekoch/gouniversal/shared/api/instatag/query"
)

type Response struct {
	Data   Data   `json:"data"`
	Status string `json:"status"`
}

type Data struct {
	Hashtag Hashtag `json:"hashtag"`
}

type Hashtag struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	EhtMedia EhtMedia `json:"edge_hashtag_to_media"`
}

// EhtMedia edge_hashtag_to_media
type EhtMedia struct {
	PageInfo PageInfo    `json:"page_info"`
	Edges    []MediaEdge `json:"edges"`
}

type PageInfo struct {
	HasNextPage bool   `json:"has_next_page"`
	EndCursor   string `json:"end_cursor"`
}

type MediaEdge struct {
	Node MediaNode `json:"node"`
}

type MediaNode struct {
	ID          string      `json:"id"`
	EmtCaption  EmtCaption  `json:"edge_media_to_caption"`
	DisplayURL  string      `json:"display_url"`
	IsVideo     bool        `json:"is_video"`
	VideoURL    string      `json:"video_url"`
	Owner       Owner       `json:"owner"`
	EstChildren EstChildren `json:"edge_sidecar_to_children"`
}

type Owner struct {
	ID string `json:"id"`
}

// EstChildren edge_sidecar_to_children
type EstChildren struct {
	Edges []ChildrenEdge `json:"edges"`
}

type ChildrenEdge struct {
	Node ChildrenNode `json:"node"`
}

type ChildrenNode struct {
	ID         string     `json:"id"`
	EmtCaption EmtCaption `json:"edge_media_to_caption"`
	DisplayURL string     `json:"display_url"`
	IsVideo    bool       `json:"is_video"`
	VideoURL   string     `json:"video_url"`
}

// EmtCaption edge_media_to_caption
type EmtCaption struct {
	Edges []CaptionEdges `json:"edges"`
}

type CaptionEdges struct {
	Node CaptionNode `json:"node"`
}

type CaptionNode struct {
	Text string `json:"text"`
}

func GetResponse(tagname, queryhash string, first int, after string, ic *instaclient.InstaClient) (Response, error) {

	var (
		err error
		ret Response
	)

	func() {

		var (
			b  []byte
			ir Response
		)

		irFirst := first
		irAfter := after

		for {

			for i := 0; i <= 4; i++ {

				switch i {
				case 0:
					b, err = query.Send(tagname, queryhash, irFirst, irAfter, ic)

				case 1:
					err = json.Unmarshal(b, &ir)

				case 2:
					if ir.Status != "ok" {
						err = errors.New(ir.Status)
					}

				case 3:
					ret.Status = ir.Status
					ret.Data.Hashtag.ID = ir.Data.Hashtag.ID
					ret.Data.Hashtag.Name = ir.Data.Hashtag.Name
					ret.Data.Hashtag.EhtMedia.PageInfo = ir.Data.Hashtag.EhtMedia.PageInfo

					if len(ir.Data.Hashtag.EhtMedia.Edges) > irFirst {
						// got more than requested
						for i := 0; i < irFirst; i++ {

							ret.Data.Hashtag.EhtMedia.Edges = append(ret.Data.Hashtag.EhtMedia.Edges, ir.Data.Hashtag.EhtMedia.Edges[i])
						}
					} else {
						ret.Data.Hashtag.EhtMedia.Edges = append(ret.Data.Hashtag.EhtMedia.Edges, ir.Data.Hashtag.EhtMedia.Edges...)
					}

				case 4:
					if len(ret.Data.Hashtag.EhtMedia.Edges) >= first {
						return
					}

					if ret.Data.Hashtag.EhtMedia.PageInfo.HasNextPage == false {
						return
					}

					irFirst = first - len(ret.Data.Hashtag.EhtMedia.Edges)
					irAfter = ret.Data.Hashtag.EhtMedia.PageInfo.EndCursor
				}

				if err != nil {
					err = errors.New(err.Error() + " \"" + string(b) + "\"")
					return
				}
			}
		}
	}()

	return ret, err
}
