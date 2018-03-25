package test

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// HTML5 SSE

/*
<p id="p1"></p>

<script type="text/javascript">

    // Create a new HTML5 EventSource
    var source = new EventSource('/events/');

    // Create a callback for when a new message is received.
    source.onmessage = function(e) {

        // Append the `data` attribute of the message to the DOM.
        document.getElementById("p1").innerHTML = e.data;
    };
</script>
*/

// SSE writes Server-Sent Events to an HTTP client.
type SSE struct{}

var messages = make(chan string)

func (s *SSE) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	f, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "cannot stream", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")

	cn, ok := rw.(http.CloseNotifier)
	if !ok {
		http.Error(rw, "cannot stream", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-cn.CloseNotify():
			log.Println("done: closed connection")
			return
		case msg := <-messages:
			fmt.Fprintf(rw, "data: %s\n\n", msg)
			f.Flush()
		}
	}
}

func Start() {
	http.Handle("/events/", &SSE{})

	go func() {
		for i := 0; ; i++ {

			// Create a little message to send to clients,
			// including the current time.
			messages <- fmt.Sprintf("%d - the time is %v", i, time.Now())

			// Print a nice log message and sleep for 5s.
			log.Printf("Sent message %d ", i)
			time.Sleep(5 * 1e9)

		}
	}()
}
