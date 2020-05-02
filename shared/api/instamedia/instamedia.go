package instamedia

import (
	"encoding/json"
	"errors"

	"github.com/dekoch/gouniversal/shared/api/instaclient"
	"github.com/dekoch/gouniversal/shared/api/instamedia/query"
)

type Response struct {
	Data   Data   `json:"data"`
	Status string `json:"status"`
}

type Data struct {
	User User `json:"user"`
}

type User struct {
	EottMedia EottMedia `json:"edge_owner_to_timeline_media"`
}

// EottMedia edge_owner_to_timeline_media
type EottMedia struct {
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
	DisplayURL  string      `json:"display_url"`
	IsVideo     bool        `json:"is_video"`
	VideoURL    string      `json:"video_url"`
	Owner       Owner       `json:"owner"`
	EstChildren EstChildren `json:"edge_sidecar_to_children"`
}

type Owner struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
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

func GetResponse(id, queryhash string, first int, after string, ic *instaclient.InstaClient) (Response, error) {

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
					b, err = query.Send(id, queryhash, irFirst, irAfter, ic)

				case 1:
					err = json.Unmarshal(b, &ir)
					if err != nil {
						err = errors.New(err.Error() + " \"" + string(b) + "\"")
					}

				case 2:
					if ir.Status != "ok" {
						err = errors.New(ir.Status)
					}

				case 3:
					ret.Status = ir.Status
					ret.Data.User.EottMedia.PageInfo = ir.Data.User.EottMedia.PageInfo

					if len(ir.Data.User.EottMedia.Edges) > irFirst {
						// got more than requested
						for i := 0; i < irFirst; i++ {

							ret.Data.User.EottMedia.Edges = append(ret.Data.User.EottMedia.Edges, ir.Data.User.EottMedia.Edges[i])
						}
					} else {
						ret.Data.User.EottMedia.Edges = append(ret.Data.User.EottMedia.Edges, ir.Data.User.EottMedia.Edges...)
					}

				case 4:
					if len(ret.Data.User.EottMedia.Edges) >= first {
						return
					}

					if ret.Data.User.EottMedia.PageInfo.HasNextPage == false {
						return
					}

					irFirst = first - len(ret.Data.User.EottMedia.Edges)
					irAfter = ret.Data.User.EottMedia.PageInfo.EndCursor
				}

				if err != nil {
					return
				}
			}
		}
	}()

	return ret, nil
}
