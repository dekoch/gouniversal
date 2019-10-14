package instaresp

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/dekoch/gouniversal/shared/io/sqlite3"
)

const TableName = "instaresp"

type InstaResp struct {
	UserID   string
	UserName string
	Checked  time.Time
	Response Response
}

type RespFile struct {
	FileID     string
	DisplayURL string
	IsVideo    bool
	VideoURL   string
}

type Response struct {
	Data   Data   `json:"data"`
	Status string `json:"status"`
}

type Data struct {
	User User `json:"user"`
}

type User struct {
	Eottm Eottm `json:"edge_owner_to_timeline_media"`
}

type Eottm struct {
	Count    uint     `json:"count"`
	PageInfo PageInfo `json:"page_info"`
	Edges    []Edge   `json:"edges"`
}

type PageInfo struct {
	HasNextPage bool   `json:"has_next_page"`
	EndCursor   string `json:"end_cursor"`
}

type Edge struct {
	Node Node `json:"node"`
}

type Node struct {
	ID                    string `json:"id"`
	DisplayURL            string `json:"display_url"`
	IsVideo               bool   `json:"is_video"`
	VideoURL              string `json:"video_url"`
	Owner                 Owner  `json:"owner"`
	EdgeSidecarToChildren Estc   `json:"edge_sidecar_to_children"`
}

type Owner struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
}

type Estc struct {
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

func (ir *Response) Unmarshal(b []byte) error {

	return json.Unmarshal(b, &ir)
}

func (ir *Response) Marshal() ([]byte, error) {

	return json.Marshal(ir)
}

func (ir *InstaResp) GetFiles() []RespFile {

	var (
		ret []RespFile
		n   RespFile
	)

	for _, edge := range ir.Response.Data.User.Eottm.Edges {

		n.FileID = edge.Node.ID
		n.DisplayURL = edge.Node.DisplayURL
		n.IsVideo = edge.Node.IsVideo
		n.VideoURL = edge.Node.VideoURL

		ret = append(ret, n)

		for _, child := range edge.Node.EdgeSidecarToChildren.Edges {

			n.FileID = child.Node.ID
			n.DisplayURL = child.Node.DisplayURL
			n.IsVideo = child.Node.IsVideo
			n.VideoURL = child.Node.VideoURL

			ret = append(ret, n)
		}
	}

	return ret
}

func LoadConfig(dbconn *sqlite3.SQLite) error {

	var lyt sqlite3.Layout
	lyt.SetTableName(TableName)
	lyt.AddField("userid", sqlite3.TypeTEXT, true, false)
	lyt.AddField("username", sqlite3.TypeTEXT, false, false)
	lyt.AddField("checked", sqlite3.TypeDATE, false, false)
	lyt.AddField("response", sqlite3.TypeTEXT, false, false)

	return dbconn.CreateTableFromLayout(lyt)
}

func (ir *InstaResp) Save(tx *sql.Tx) error {

	resp, err := ir.Response.Marshal()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT OR REPLACE INTO `"+TableName+"` (userid, username, checked, response) values(?,?,?,?)", ir.UserID, ir.UserName, ir.Checked, string(resp))
	return err
}

func (ir *InstaResp) Load(dbconn *sqlite3.SQLite) (bool, error) {

	var (
		err   error
		found bool
		resp  string
	)

	func() {

		var rows *sql.Rows

		for i := 0; i <= 3; i++ {

			switch i {
			case 0:
				rows, err = dbconn.DB.Query("SELECT username, checked, response FROM `"+TableName+"` WHERE userid=?", ir.UserID)

			case 1:
				defer rows.Close()

			case 2:
				for rows.Next() {

					err = rows.Scan(&ir.UserName, &ir.Checked, &resp)

					found = true
				}

			case 3:
				if found {
					err = ir.Response.Unmarshal([]byte(resp))
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return found, err
}
