package pagehome

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/dekoch/gouniversal/module/gpsnav/core"
	"github.com/dekoch/gouniversal/module/gpsnav/global"
	"github.com/dekoch/gouniversal/module/gpsnav/lang"
	"github.com/dekoch/gouniversal/module/gpsnav/typenav"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/sbool"
	"github.com/dekoch/gouniversal/shared/stringarray"
)

// SSE writes Server-Sent Events to an HTTP client.
type sse struct{}

type sseMessage struct {
	ClientUUID string
	Content    string
}

type sseJSON struct {
	Time         string
	State        string
	Step         string
	CurrentPos   string
	NextWaypoint string
	Wpt          string
	Bearing      string
	Distance     string
}

var (
	messages      = make(chan sseMessage)
	clients       stringarray.StringArray
	streamEnabled sbool.Sbool
)

func LoadConfig() {

	clients.RemoveAll()

	http.Handle("/gpsnav/sse/", &sse{})
}

func RegisterPage(page *typenav.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.Home.Menu, "App:GPSNav:Home", page.Lang.Home.Title)
}

func Render(page *typenav.Page, nav *navigation.Navigation, r *http.Request) {

	type content struct {
		Lang  lang.Home
		UUID  template.HTML
		Token template.HTML
	}
	var c content

	c.Lang = page.Lang.Home

	c.UUID = template.HTML(nav.User.UUID)
	c.Token = template.HTML(global.Tokens.New(nav.User.UUID))

	switch r.FormValue("edit") {
	case "start":
		core.Start(0)

	case "stop":
		core.Stop()
	}

	cont, err := functions.PageToString(global.Config.UIFileRoot+"home.html", c)
	if err == nil {
		page.Content += cont
	} else {
		nav.RedirectPath("404", true)
	}

	stream()
}

func (s *sse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "cannot stream", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	cn, ok := w.(http.CloseNotifier)
	if !ok {
		http.Error(w, "cannot stream", http.StatusInternalServerError)
		return
	}

	id := r.FormValue("uuid")

	if functions.IsEmpty(id) {
		fmt.Fprintf(w, "error: UUID not set")
		f.Flush()
		return
	}

	tok := r.FormValue("token")

	if functions.IsEmpty(tok) {
		fmt.Fprintf(w, "error: token not set")
		f.Flush()
		return
	}

	if global.Tokens.Check(id, tok) == false {
		fmt.Fprintf(w, "error: invalid token")
		f.Flush()
		return
	}

	clients.Add(id)

	for {
		select {
		case <-cn.CloseNotify():
			clients.Remove(id)
			//fmt.Println("done: closed connection")
			return
		case ssem := <-messages:
			// check client uuid
			if ssem.ClientUUID == id {

				fmt.Fprintf(w, "data: %s\n\n", ssem.Content)
				f.Flush()
			}
		}
	}
}

func stream() {

	if streamEnabled.IsSet() {
		return
	}

	streamEnabled.Set()

	go func() {

		for {
			time.Sleep(100 * time.Millisecond)

			c := clients.List()

			if len(c) > 0 {

				var ssej sseJSON

				ssej.Time = time.Now().Format("15:04:05.000")

				switch core.GetState() {
				case typenav.STOPPED:
					ssej.State = "Stopped"

				case typenav.RUNNING:
					ssej.State = "Running"
				}

				ssej.Step = string(core.GetStep())

				pos := core.GetCurrentPos()
				ssej.CurrentPos = posToString(pos)

				pos, _ = core.GetNextWaypoint()
				ssej.NextWaypoint = posToString(pos)

				no := core.GetWaypointNo()
				cnt := core.GetWaypointCnt()
				ssej.Wpt = strconv.Itoa(no) + "/" + strconv.Itoa(cnt)

				bearing, _ := core.GetBearing()
				ssej.Bearing = strconv.FormatFloat(bearing, 'f', 1, 64)

				dist, _ := core.GetDistance()
				ssej.Distance = strconv.FormatFloat(dist, 'f', 1, 64)

				var ssem sseMessage

				b, err := json.Marshal(ssej)
				if err != nil {
					continue
				}

				ssem.Content = string(b)

				// send message to all waiting clients
				for i := 0; i < len(c); i++ {

					ssem.ClientUUID = c[i]

					messages <- ssem
				}
			}
		}
	}()
}

func posToString(p typenav.Pos) string {

	ret := "lat:" + strconv.FormatFloat(p.Lat, 'f', 6, 64)
	ret += " lon:" + strconv.FormatFloat(p.Lon, 'f', 6, 64)
	ret += " ele:" + strconv.FormatFloat(p.Ele, 'f', 1, 64)

	if len(p.Name) > 0 {
		ret += " (" + p.Name + ")"
	}

	return ret
}
