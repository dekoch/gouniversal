package pagehome

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/dekoch/gouniversal/module/picturex/global"
	"github.com/dekoch/gouniversal/module/picturex/lang"
	"github.com/dekoch/gouniversal/module/picturex/typemo"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/sbool"
	"github.com/dekoch/gouniversal/shared/stringarray"
)

// SSE writes Server-Sent Events to an HTTP client.
type sse struct{}

type sseMessage struct {
	ClientUUID string
	PairUUID   string
	Content    string
}

type sseJSON struct {
	First  string
	Second string
}

var (
	messages      = make(chan sseMessage)
	clients       stringarray.StringArray
	streamEnabled sbool.Sbool
)

func LoadConfig() {

	clients.RemoveAll()

	http.Handle("/picturex/sse/", &sse{})
}

func RegisterPage(page *typemo.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.Home.Menu, "App:PictureX:Home", page.Lang.Home.Title)

	pairs, err := global.PairList.GetPairsFromUser(nav.User.UUID)
	if err == nil {
		for _, pair := range pairs {
			nav.Sitemap.Register(page.Lang.Home.Menu, "App:PictureX:Home$Pair="+pair, page.Lang.Home.Title+" ("+pair+")")
		}
	}
}

func Render(page *typemo.Page, nav *navigation.Navigation, r *http.Request) {

	type Content struct {
		Lang             lang.Home
		UUID             template.HTML
		Token            template.HTML
		Pair             template.HTML
		Link             template.HTML
		ShowShareLink    template.HTML
		ShowLinkReceived template.HTML
	}
	var c Content

	c.Lang = page.Lang.Home

	var (
		err         error
		redirect    bool
		isFirstUser bool
		ssej        sseJSON
		ssem        sseMessage
		firstPic    string
		secondPic   string
		b           []byte
	)

	pair := nav.Parameter("Pair")
	token := global.Tokens.New(nav.User.UUID)

	func() {
		for i := 0; i <= 8; i++ {

			switch i {
			case 0:
				if pair == "" {
					pair, err = global.PairList.NewPair(nav.User.UUID)
					redirect = true
				}

			case 1:
				switch r.FormValue("edit") {
				case "newpair":
					pair, err = global.PairList.NewPair(nav.User.UUID)
					redirect = true

				case "unlock":
					global.PairList.UnlockPicture(pair, nav.User.UUID)

				case "deletepair":
					global.PairList.DeletePair(pair)
					pair, err = global.PairList.GetFirstPairFromUser(nav.User.UUID)
					if err != nil {
						pair, err = global.PairList.NewPair(nav.User.UUID)
					}

					redirect = true
				}

			case 2:
				isFirstUser, err = global.PairList.IsFirstUser(pair, nav.User.UUID)

			case 3:
				if isFirstUser {
					c.ShowShareLink = template.HTML("show")
				} else {
					err = global.PairList.SetSecondUser(pair, nav.User.UUID)
					c.ShowLinkReceived = template.HTML("show")
				}

			case 4:
				firstPic, err = global.PairList.GetFirstPicture(pair, nav.User.UUID)

			case 5:
				secondPic, err = global.PairList.GetSecondPicture(pair, nav.User.UUID)

			case 6:
				if firstPic != "" {
					ssej.First = "/picturex/req/?uuid=" + nav.User.UUID + "&token=" + token + "&name=" + firstPic
				}

				if secondPic != "" {
					ssej.Second = "/picturex/req/?uuid=" + nav.User.UUID + "&token=" + token + "&name=" + secondPic
				}

			case 7:
				b, err = json.Marshal(ssej)

			case 8:
				ssem.PairUUID = pair
				ssem.Content = string(b)

				go func(message sseMessage) {

					time.Sleep(500 * time.Millisecond)

					cl := clients.List()
					// send message to all waiting clients
					for i := 0; i < len(cl); i++ {

						ssem.ClientUUID = cl[i]

						messages <- ssem
					}
				}(ssem)
			}

			if err != nil {
				pair, _ = global.PairList.NewPair(nav.User.UUID)
				redirect = true

				alert.Message(alert.ERROR, page.Lang.Alert.Error, err, "", nav.User.UUID)
				return
			}
		}
	}()

	if redirect {
		nav.RedirectPath("App:PictureX:Home$Pair="+pair, false)
		return
	}

	c.UUID = template.HTML(nav.User.UUID)
	c.Token = template.HTML(token)
	c.Pair = template.HTML(pair)

	link := "http://"

	if nav.UIConfig.HTTPS.Enabled {
		link = "https://"
	}

	link += nav.Server.Host + "?path=App:PictureX:Home$Pair=" + pair
	c.Link = template.HTML(link)

	p, err := functions.PageToString(global.Config.UIFileRoot+"home.html", c)
	if err == nil {
		page.Content += p
	} else {
		nav.RedirectPath("404", true)
	}
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

	pair := r.FormValue("pair")

	clients.Add(id)

	for {
		select {
		case <-cn.CloseNotify():
			clients.Remove(id)
			//fmt.Println("done: closed connection")
			return
		case ssem := <-messages:
			// check client uuid
			if ssem.ClientUUID == id &&
				ssem.PairUUID == pair {

				fmt.Fprintf(w, "data: %s\n\n", ssem.Content)
				f.Flush()
			}
		}
	}
}
