package alert

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dekoch/gouniversal/shared/counter"
	"github.com/dekoch/gouniversal/shared/functions"
	token "github.com/dekoch/gouniversal/shared/token"
)

// SSE writes Server-Sent Events to an HTTP client.
type sse struct{}

type alertType int

const (
	NONE alertType = 1 + iota
	SUCCESS
	INFO
	WARNING
	ERROR
)

type alertMessage struct {
	ClientUUID string
	Content    string
}

var (
	messages = make(chan alertMessage)
	clients  counter.Counter
	Tokens   token.Token
)

func Start() {

	clients.Reset()

	http.Handle("/alert/", &sse{})
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

	if Tokens.Check(id, tok) == false {
		fmt.Fprintf(w, "error: invalid token")
		f.Flush()
		return
	}

	clients.Add()

	for {
		select {
		case <-cn.CloseNotify():
			clients.Remove()
			//fmt.Println("done: closed connection")
			return
		case am := <-messages:
			// check client uuid
			if am.ClientUUID == id ||
				am.ClientUUID == "all" {

				fmt.Fprintf(w, "data: %s\n\n", am.Content)
				f.Flush()
			}
		}
	}
}

func Message(t alertType, title string, message interface{}, sender string, clientuuid string) {

	if t == NONE {
		return
	}

	go func() {

		m := fmt.Sprintf("%v", message)
		ti := time.Now().Format("15:04:05")
		alert := ""

		if t == SUCCESS {
			alert += "<div class=\"alert alert-success alert-dismissible\">"
			alert += "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a>"
			alert += ti + " <strong>" + title + ":</strong> " + m
			alert += "</div>"
		} else if t == INFO {
			alert += "<div class=\"alert alert-info alert-dismissible\">"
			alert += "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a>"
			alert += ti + " <strong>" + title + ":</strong> " + m
			alert += "</div>"
		} else if t == WARNING {
			alert += "<div class=\"alert alert-warning alert-dismissible\">"
			alert += "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a>"
			alert += ti + " <strong>" + title + ":</strong> " + m
			alert += "</div>"
		} else {
			alert += "<div class=\"alert alert-danger alert-dismissible\">"
			alert += "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a>"
			alert += ti + " <strong>" + title + ":</strong> " + m
			alert += "</div>"
		}

		time.Sleep(500 * time.Millisecond)

		var am alertMessage
		am.ClientUUID = clientuuid
		am.Content = alert

		// send message to all waiting clients
		for i := 0; i < clients.GetCount(); i++ {

			messages <- am
		}
	}()
}
