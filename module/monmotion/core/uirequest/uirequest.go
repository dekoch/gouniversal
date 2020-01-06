package uirequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/dekoch/gouniversal/module/monmotion/typemd"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/token"
)

type UIRequest struct {
	tokens      token.Token
	latestImage typemd.MoImage
	active      bool
}

var mut sync.RWMutex

func (rt *UIRequest) LoadConfig(device string) error {

	mut.Lock()
	defer mut.Unlock()

	if rt.active {
		return nil
	}

	http.HandleFunc("/monmotion/"+device+"/", rt.serve)

	rt.active = true

	return nil
}

func (rt *UIRequest) GetNewToken(uid string) string {

	return rt.tokens.New(uid)
}

func (rt *UIRequest) SetLatestImage(img typemd.MoImage) {

	mut.Lock()
	defer mut.Unlock()

	rt.latestImage = img
}

func (rt *UIRequest) getLatestImage() typemd.MoImage {

	mut.RLock()
	defer mut.RUnlock()

	return rt.latestImage
}

func (rt *UIRequest) serve(w http.ResponseWriter, r *http.Request) {

	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "cannot stream", http.StatusInternalServerError)
		return
	}

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

	if rt.tokens.Check(id, tok) == false {
		fmt.Fprintf(w, "error: invalid token")
		f.Flush()
		return
	}

	content := r.FormValue("content")

	if functions.IsEmpty(content) {
		fmt.Fprintf(w, "error: content not set")
		f.Flush()
		return
	}

	switch content {
	case "image":
		w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")

	case "info":
		w.Header().Set("Content-Type", "text/event-stream")
	}

	timerStream := time.NewTimer(100 * time.Millisecond)

	for {
		select {
		case <-cn.CloseNotify():
			return

		case <-timerStream.C:
			switch content {
			case "image":
				s, err := rt.getImage()
				if err == nil {
					fmt.Fprintf(w, "%s", []byte("--frame\r\n  Content-Type: image/jpeg\r\n\r\n"+s+"\r\n\r\n"))
					f.Flush()
				}

			case "info":
				s, err := rt.getInfo()
				if err == nil {
					fmt.Fprintf(w, "data: %s\n\n", s)
					f.Flush()
				}
			}

			timerStream.Reset(100 * time.Millisecond)
		}
	}
}

func (rt *UIRequest) getImage() (string, error) {

	img := rt.getLatestImage()

	if img.Img == nil {
		return "", nil
	}

	buf := &bytes.Buffer{}
	err := jpeg.Encode(buf, img.Img, nil)
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}

func (rt *UIRequest) getInfo() (string, error) {

	img := rt.getLatestImage()

	if img.Img == nil {
		return "", nil
	}

	type jsonInfo struct {
		Captured string
		Size     string
	}

	var jn jsonInfo
	jn.Captured = img.Captured.Format("15:04:05.0000")
	jn.Size = strconv.Itoa(img.Img.Bounds().Size().X) + "x" + strconv.Itoa(img.Img.Bounds().Size().Y)

	b, err := json.Marshal(jn)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
