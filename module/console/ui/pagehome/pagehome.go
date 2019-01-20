package pagehome

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/console/global"
	"github.com/dekoch/gouniversal/module/console/typeconsole"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/stringarray"
	token "github.com/dekoch/gouniversal/shared/token"
)

// SSE writes Server-Sent Events to an HTTP client.
type sse struct{}

type consoleMessage struct {
	ClientUUID string
	Content    string
}

var (
	mut           sync.Mutex
	messages      = make(chan consoleMessage)
	clients       stringarray.StringArray
	mytoken       token.Token
	streamEnabled bool
)

func LoadConfig() {

	clients.RemoveAll()

	http.Handle("/console/", &sse{})
}

func RegisterPage(page *typeconsole.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program", "App:Console:Home", page.Lang.Home.Title)
}

func Render(page *typeconsole.Page, nav *navigation.Navigation, r *http.Request) {

	type content struct {
		UUID  template.HTML
		Token template.HTML
	}
	var c content

	c.UUID = template.HTML(nav.User.UUID)
	c.Token = template.HTML(mytoken.New(nav.User.UUID))

	p, err := functions.PageToString(global.Config.UIFileRoot+"console.html", c)
	if err == nil {
		page.Content += p
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

	if mytoken.Check(id, tok) == false {
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
		case cm := <-messages:
			// check client uuid
			if cm.ClientUUID == id {

				fmt.Fprintf(w, "data: %s\n\n", cm.Content)
				f.Flush()
			}
		}
	}
}

func stream() {

	mut.Lock()
	defer mut.Unlock()

	if streamEnabled {
		return
	}

	streamEnabled = true

	go func() {

		for {
			time.Sleep(global.Config.RefreshInterval * time.Millisecond)

			c := clients.List()

			if len(c) > 0 {

				consoleOuput := console.GetOutput()

				var o string

				for i := 0; i < len(consoleOuput); i++ {
					o += consoleOuput[i] + "<br>"
				}

				var cm consoleMessage
				cm.Content = o

				// send message to all waiting clients
				for i := 0; i < len(c); i++ {

					cm.ClientUUID = c[i]

					messages <- cm
				}
			}
		}
	}()
}
