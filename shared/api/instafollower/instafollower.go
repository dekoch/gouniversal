package instafollower

import (
	"encoding/json"
	"errors"

	"github.com/dekoch/gouniversal/shared/api/instaclient"
	"github.com/dekoch/gouniversal/shared/api/instafollower/query"
)

type Response struct {
	Data   Data   `json:"data"`
	Status string `json:"status"`
}

type Data struct {
	User User `json:"user"`
}

type User struct {
	EdgeFollow EdgeFollow `json:"edge_follow"`
}

type EdgeFollow struct {
	PageInfo PageInfo    `json:"page_info"`
	Edges    []EdgesNode `json:"edges"`
}

type PageInfo struct {
	HasNextPage bool   `json:"has_next_page"`
	EndCursor   string `json:"end_cursor"`
}

type EdgesNode struct {
	Node EdgesNodeContent `json:"node"`
}

type EdgesNodeContent struct {
	ID                string `json:"id"`
	UserName          string `json:"username"`
	FullName          string `json:"full_name"`
	ProfilePicURL     string `json:"profile_pic_url"`
	IsPrivate         bool   `json:"is_private"`
	IsVerified        bool   `json:"is_verified"`
	FollowedByViewer  bool   `json:"followed_by_viewer"`
	RequestedByViewer bool   `json:"requested_by_viewer"`
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
					ret.Data.User.EdgeFollow.PageInfo = ir.Data.User.EdgeFollow.PageInfo

					if len(ir.Data.User.EdgeFollow.Edges) > irFirst {
						// got more than requested
						for i := 0; i < irFirst; i++ {

							ret.Data.User.EdgeFollow.Edges = append(ret.Data.User.EdgeFollow.Edges, ir.Data.User.EdgeFollow.Edges[i])
						}
					} else {
						ret.Data.User.EdgeFollow.Edges = append(ret.Data.User.EdgeFollow.Edges, ir.Data.User.EdgeFollow.Edges...)
					}

				case 4:
					if len(ret.Data.User.EdgeFollow.Edges) >= first {
						return
					}

					if ret.Data.User.EdgeFollow.PageInfo.HasNextPage == false {
						return
					}

					irFirst = first - len(ret.Data.User.EdgeFollow.Edges)
					irAfter = ret.Data.User.EdgeFollow.PageInfo.EndCursor
				}

				if err != nil {
					return
				}
			}
		}
	}()

	return ret, nil
}
