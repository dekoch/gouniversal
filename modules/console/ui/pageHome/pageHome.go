package pageHome

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/dekoch/gouniversal/modules/console/global"
	"github.com/dekoch/gouniversal/modules/console/typesConsole"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/stringArray"
)

// SSE writes Server-Sent Events to an HTTP client.
type SSE struct{}

type consoleMessage struct {
	ClientUUID string
	Content    string
}

var (
	messages      = make(chan consoleMessage)
	clients       stringArray.StringArray
	streamEnabled bool
)

func LoadConfig() {

	clients.RemoveAll()

	http.Handle("/console/", &SSE{})
}

func RegisterPage(page *typesConsole.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register("Program", "App:Console:Home", page.Lang.Home.Title)
}

func Render(page *typesConsole.Page, nav *navigation.Navigation, r *http.Request) {

	type content struct {
		UUID template.HTML
	}
	var c content

	c.UUID = template.HTML(nav.User.UUID)

	p, err := functions.PageToString(global.Config.File.UIFileRoot+"console.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}

	if streamEnabled == false {

		streamEnabled = true

		stream()
	}
}

func (s *SSE) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	//fmt.Println("new " + id)

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

	go func() {

		for {
			time.Sleep(global.Config.File.RefreshInterval * time.Millisecond)

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
