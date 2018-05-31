package alert

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// SSE writes Server-Sent Events to an HTTP client.
type SSE struct{}

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
	mut      sync.Mutex
	clients  int
)

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

	mut.Lock()
	clients += 1
	mut.Unlock()

	for {
		select {
		case <-cn.CloseNotify():
			mut.Lock()
			clients -= 1
			mut.Unlock()
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

	m := fmt.Sprintf("%v", message)

	if t == NONE {
		return
	}

	go func() {

		time.Sleep(500 * time.Millisecond)

		alert := ""

		if t == SUCCESS {
			alert += "<div class=\"alert alert-success alert-dismissible\">"
			alert += "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a>"
			alert += "<strong>" + title + "</strong> " + m
			alert += "</div>"
		} else if t == INFO {
			alert += "<div class=\"alert alert-info alert-dismissible\">"
			alert += "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a>"
			alert += "<strong>" + title + "</strong> " + m
			alert += "</div>"
		} else if t == WARNING {
			alert += "<div class=\"alert alert-warning alert-dismissible\">"
			alert += "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a>"
			alert += "<strong>" + title + "</strong> " + m
			alert += "</div>"
		} else {
			alert += "<div class=\"alert alert-danger alert-dismissible\">"
			alert += "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a>"
			alert += "<strong>" + title + "</strong> " + m
			alert += "</div>"
		}

		var am alertMessage
		am.ClientUUID = clientuuid
		am.Content = alert

		mut.Lock()
		c := clients
		mut.Unlock()

		// send alert to all waiting clients
		for i := 0; i < c; i++ {

			messages <- am
		}
	}()
}

func Start() {

	clients = 0

	http.Handle("/alert/", &SSE{})
}
