package ui

import (
	"html/template"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dekoch/gouniversal/module"
	"github.com/dekoch/gouniversal/program/global"
	"github.com/dekoch/gouniversal/program/ui/pagehome"
	"github.com/dekoch/gouniversal/program/ui/pagelogin"
	"github.com/dekoch/gouniversal/program/ui/settings"
	"github.com/dekoch/gouniversal/program/usermanagement"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/clientinfo"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/sitemap"
	"github.com/dekoch/gouniversal/shared/types"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type menuDropdown struct {
	Order int
	Title string
	Items []sitemap.Page
}

var (
	store      = new(sessions.CookieStore)
	cookieName string
)

func setCookie(parameter string, value string, w http.ResponseWriter, r *http.Request) {
	//console.Log("set cookie " + parameter + " to " + value)

	session, _ := store.Get(r, cookieName)
	session.Values[parameter] = value
	session.Save(r, w)
}

func getCookie(parameter string, w http.ResponseWriter, r *http.Request) (string, error) {
	//console.Log("get cookie " + parameter)

	session, err := store.Get(r, cookieName)

	if err == nil {
		return session.Values[parameter].(string), err
	}

	return "", err
}

func setSession(nav *navigation.Navigation, w http.ResponseWriter, r *http.Request) {
	setCookie("navigation", nav.Path, w, r)

	setCookie("authenticated", nav.User.UUID, w, r)

	if nav.GodMode {
		setCookie("isgod", "true", w, r)
	} else {
		setCookie("isgod", "", w, r)
	}
}

func getSession(nav *navigation.Navigation, w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, cookieName)
	if session.IsNew {
		initCookies(w, r)
	}

	// get stored path
	path, err := getCookie("navigation", w, r)
	if err == nil {
		nav.Path = path
	}

	// get stored user UUID
	uid, err := getCookie("authenticated", w, r)
	if err == nil {

		nav.User, err = global.UserConfig.Get(uid)

		if nav.User.State < 0 {
			// if no user found
			nav.User = global.Guests.SelectGuest(uid)
		}

		if nav.User.State < 0 {
			// if no guest found, create new
			u, err := global.UserConfig.GetWithState(0)
			if err != nil {
				console.Log(err, "")
			}
			nav.User = global.Guests.NewGuest(u, global.UIConfig.MaxGuests)
		}
	}

	// for debugging
	nav.GodMode = false
	g, err := getCookie("isgod", w, r)
	if err == nil {
		if g == "true" {
			nav.GodMode = true
		}
	}
}

func initCookies(w http.ResponseWriter, r *http.Request) {
	nav := new(navigation.Navigation)
	nav.Path = "init"

	u, err := global.UserConfig.GetWithState(0)
	if err != nil {
		console.Log(err, "")
	}
	nav.User = global.Guests.NewGuest(u, global.UIConfig.MaxGuests)
	nav.GodMode = false

	setSession(nav, w, r)
}

func renderProgram(page *types.Page, nav *navigation.Navigation) []byte {

	type content struct {
		Title     string
		MenuBrand string
		MenuLeft  template.HTML
		MenuRight template.HTML
		UUID      template.HTML
		Token     template.HTML
		Content   template.HTML
	}
	var c content

	c.MenuBrand = nav.Sitemap.PageTitle(nav.Path)

	if c.MenuBrand == "" {
		c.MenuBrand = nav.Path
	}

	// title
	if global.UIConfig.Title != "" {
		c.Title = global.UIConfig.Title + " - " + c.MenuBrand
	} else {
		c.Title = c.MenuBrand
	}

	if c.Title == "" {
		c.Title = nav.Path
	}

	pages := nav.Sitemap.Pages
	mDropdown := make([]menuDropdown, 0)

	// put each allowed page into menu slice
	for i := 0; i < len(pages); i++ {

		if pages[i].Menu != "" {

			if usermanagement.IsPageAllowed(pages[i].Path, nav.User) ||
				nav.GodMode {

				dropdownFound := false
				for d := 0; d < len(mDropdown); d++ {

					if pages[i].Menu == mDropdown[d].Title {
						dropdownFound = true
					}
				}

				if dropdownFound == false {
					// create new dropdown
					newDropdown := make([]menuDropdown, 1)
					newDropdown[0].Title = pages[i].Menu

					if pages[i].MenuOrder > -1 {

						newDropdown[0].Order = pages[i].MenuOrder
					} else {
						// set predefined order
						if newDropdown[0].Title == "Program" {

							newDropdown[0].Order = 0

						} else if newDropdown[0].Title == "App" {

							newDropdown[0].Order = 1

						} else if newDropdown[0].Title == "Account" {

							newDropdown[0].Order = 999

						} else {
							newDropdown[0].Order = 999 - len(mDropdown)
						}
					}

					mDropdown = append(mDropdown, newDropdown...)
				}

				// add items to dropdown
				for d := 0; d < len(mDropdown); d++ {

					if pages[i].Menu == mDropdown[d].Title {

						newItem := make([]sitemap.Page, 1)
						newItem[0] = pages[i]

						mDropdown[d].Items = append(mDropdown[d].Items, newItem...)
					}
				}
			}
		}
	}

	htmlMenuLeft := ""
	htmlMenuRight := ""

	sort.Slice(mDropdown, func(i, j int) bool { return mDropdown[i].Order < mDropdown[j].Order })

	// DropDown to HTML
	for d := 0; d < len(mDropdown); d++ {

		alignRight := false
		dropDownTitle := mDropdown[d].Title

		if dropDownTitle == "Program" {

			dropDownTitle = page.Lang.Menu.Program.Title

		} else if dropDownTitle == "App" {

			dropDownTitle = page.Lang.Menu.App.Title

		} else if dropDownTitle == "Account" {

			alignRight = true

			if nav.User.UUID != "" {
				dropDownTitle = nav.User.LoginName
			} else {
				dropDownTitle = page.Lang.Menu.Account.Title
			}
		}

		dropDown := "<li class=\"nav-item dropdown\">\n"
		dropDown += "<a class=\"nav-link dropdown-toggle\" href=\"\" id=\"navbar" + mDropdown[d].Title + "\" role=\"button\" data-toggle=\"dropdown\" aria-haspopup=\"true\" aria-expanded=\"false\">\n"
		dropDown += dropDownTitle + "\n"
		dropDown += "</a>\n"
		dropDown += "<div class=\"dropdown-menu"

		if alignRight {
			dropDown += " dropdown-menu-right"
		}

		dropDown += "\" aria-labelledby=\"navbar" + mDropdown[d].Title + "\">\n"

		for i := 0; i < len(mDropdown[d].Items); i++ {

			dropDown += "<button class=\"dropdown-item\" type=\"submit\" name=\"navigation\" value=\"" + mDropdown[d].Items[i].Path + "\">" + mDropdown[d].Items[i].Title + "</button>\n"
		}

		dropDown += "</div>\n</li>\n"

		if alignRight {

			htmlMenuRight += dropDown
		} else {

			htmlMenuLeft += dropDown
		}
	}

	c.MenuLeft = template.HTML(htmlMenuLeft)
	c.MenuRight = template.HTML(htmlMenuRight)
	c.UUID = template.HTML(nav.User.UUID)
	c.Content = template.HTML(page.Content)

	alert.Tokens.SetMaxTokens(global.UserConfig.GetUserCnt() + global.UIConfig.MaxGuests)
	c.Token = template.HTML(alert.Tokens.New(nav.User.UUID))

	p, err := functions.PageToString(global.UIConfig.ProgramFileRoot+"program.html", c)
	if err != nil {
		console.Log(err, "ui.go")
		p = err.Error()
	}

	return []byte(p)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {

	newPath := r.URL.Query().Get("path")

	if r.URL.Path == "/" ||
		newPath != "" {

		session, _ := store.Get(r, cookieName)
		if session.IsNew {
			initCookies(w, r)
		}

		// set new path to cookie
		if newPath != "" {
			nav := new(navigation.Navigation)
			getSession(nav, w, r)
			nav.NavigatePath(newPath)
			setSession(nav, w, r)
			console.Output(newPath, "app")
		}

		// redirect to HTTPS if enabled
		if global.UIConfig.HTTPS.Enabled {
			host := strings.Split(r.Host, ":")
			port := strconv.Itoa(global.UIConfig.HTTPS.Port)

			http.Redirect(w, r, "https://"+host[0]+":"+port+"/app/", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/app/", http.StatusSeeOther)
		}

	} else if r.URL.Path == "/favicon.ico" {
		icon := global.UIConfig.StaticFileRoot + "favicon.ico"

		if _, err := os.Stat(icon); os.IsNotExist(err) == false {
			http.ServeFile(w, r, icon)
		}
	} else {
		console.Log("Error 404 \""+r.URL.Path+"\" ("+clientinfo.String(r)+")", "handleRoot()")

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleApp(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	page := new(types.Page)
	nav := new(navigation.Navigation)

	getSession(nav, w, r)

	global.Lang.SelectLang(nav.User.Lang, &page.Lang)

	pagehome.RegisterPage(page, nav)
	module.RegisterPage(page, nav)
	settings.RegisterPage(page, nav)
	nav.Sitemap.Register("Program", "Program:Exit", page.Lang.Exit.Title)
	pagelogin.RegisterPage(page, nav)

	newPath := r.FormValue("navigation")

	if newPath != "" {
		nav.NavigatePath(newPath)

		console.Output(newPath, "app")
	}

	// select first allowed page
	if nav.Path == "init" {

		for _, p := range nav.Sitemap.Pages {

			if p.Menu != "" && nav.Path == "init" {

				if usermanagement.IsPageAllowed(p.Path, nav.User) ||
					nav.GodMode {

					nav.Path = p.Path
				}
			}
		}
	}

	nav.Redirect = "init"
	loopDetection := 0

	for functions.IsEmpty(nav.Redirect) == false {

		loopDetection++

		nav.CurrentPath = ""
		nav.Redirect = ""
		page.Content = ""

		// select next page
		// e.g. nav.NavigatePath("Program:Settings:User:List")
		//
		// Program
		// Program:Settings
		// Program:Settings:User
		// Program:Settings:User:List
		switch nav.GetNextPage() {
		case "Program":
			switch nav.GetNextPage() {
			case "Home":
				pagehome.Render(page, nav, r)

			case "Settings":
				settings.Render(page, nav, r)

			case "Exit":
				// do nothing

			default:
				nav.RedirectPath("404", true)
			}

		case "App":
			module.Render(page, nav, r)

		case "Account":
			switch nav.GetNextPage() {
			case "Login":
				pagelogin.Render(page, nav, r)

			case "Logout":
				console.Log("\""+nav.User.LoginName+"\" logged out", "Logout")

				u, err := global.UserConfig.GetWithState(0)
				if err != nil {
					console.Log(err, "")
				}
				nav.User = global.Guests.NewGuest(u, global.UIConfig.MaxGuests)

				nav.RedirectPath("Account:Login", false)

			default:
				nav.RedirectPath("404", true)
			}

		default:
			if nav.Path == "400" {
				// Bad Request
				console.Log("400 \""+nav.LastPath+"\"", "")
				page.Content = "<h1>400</h1><br>"
				page.Content += page.Lang.Error.CE400BadRequest
			} else if nav.Path == "404" {
				// NotFound404
				console.Log("404 \""+nav.LastPath+"\"", "")
				page.Content = "<h1>404</h1><br>"
				page.Content += page.Lang.Error.CE404NotFound
			} else if nav.Path == "508" {
				// LoopDetected508
				console.Log("508 \""+nav.Path+"\" ("+nav.LastPath+")", "")
				page.Content = "<h1>508</h1><br>"
				page.Content += page.Lang.Error.SE508LoopDetected
			} else {
				nav.RedirectPath("404", true)
			}
		}

		if loopDetection >= 5 {
			loopDetection = 0

			nav.RedirectPath("508", true)
		}

		if nav.Redirect != "" {

			if nav.Redirect == "400" ||
				nav.Redirect == "404" ||
				nav.Redirect == "508" {

				nav.Path = nav.Redirect
			} else {
				nav.NavigatePath(nav.Redirect)
			}
		}
	}

	if nav.Path == "Program:Exit" {

		page.Content = "goodbye"

		console.Input("exit")
	} else {

		setSession(nav, w, r)
	}

	w.Write(renderProgram(page, nav))

	// show all pages
	//nav.Sitemap.ShowMap()

	// show allowed pages
	/*var sm []string
	sm = nav.Sitemap.PageList()
	console.Output("####")
	for i := 0; i < len(sm); i++ {
		if usermanagement.IsPageAllowed(sm[i], nav.User) {
			console.Output(sm[i])
		}
	}
	console.Output("####")*/

	t := time.Now()
	elapsed := t.Sub(start)
	f := elapsed.Seconds() * 1000.0
	console.Output(nav.Path+" "+strconv.FormatFloat(f, 'f', 1, 64)+"ms", "app")
}

func handleRecovery(w http.ResponseWriter, r *http.Request) {

	if global.UIConfig.Recovery {

		button := r.FormValue("recovery")

		if button == "goback" {

			http.Redirect(w, r, "/", http.StatusSeeOther)

		} else if button == "disablerecovery" {

			global.UIConfig.Recovery = false
			global.UIConfig.SaveConfig()

		} else if button == "cookies" {

			initCookies(w, r)

		} else if button == "god" {

			setCookie("isgod", "true", w, r)

		} else {
			console.Log(button, "ui.go")
		}

		type recovery struct {
			Temp string
		}
		var re recovery

		re.Temp = ""

		templ, err := template.ParseFiles(global.UIConfig.ProgramFileRoot + "recovery.html")
		if err != nil {
			console.Log(err, "")
		}
		templ.Execute(w, re)
	} else {
		w.Write([]byte("recovery mode is DISABLED!"))
	}
}

func StartServer() {
	console.Log("starting webserver...", " ")

	if _, err := os.Stat(global.UIConfig.ProgramFileRoot); os.IsNotExist(err) {
		// if not found, exit program
		console.Log("error: ProgramFileRoot not found", "ui.go")
	}

	if _, err := os.Stat(global.UIConfig.StaticFileRoot); os.IsNotExist(err) {
		// if not found, exit program
		console.Log("error: StaticFileRoot not found", "ui.go")
	}

	addrs, err := net.InterfaceAddrs()

	if err == nil {

		// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
		key := securecookie.GenerateRandomKey(32)
		store = sessions.NewCookieStore(key)

		u := uuid.Must(uuid.NewRandom())
		cookieName = u.String()

		global.UserConfig.LoadConfig()
		global.GroupConfig.LoadConfig()

		alert.Start()
		module.LoadConfig()

		// if HTTPS is enabled, check cert and key file
		if global.UIConfig.HTTPS.Enabled {
			if _, err = os.Stat(global.UIConfig.HTTPS.CertFile); os.IsNotExist(err) {
				global.UIConfig.HTTPS.Enabled = false
				console.Log("missing CertFile \""+global.UIConfig.HTTPS.CertFile+"\"", " ")
			}

			if _, err = os.Stat(global.UIConfig.HTTPS.KeyFile); os.IsNotExist(err) {
				global.UIConfig.HTTPS.Enabled = false
				console.Log("missing KeyFile \""+global.UIConfig.HTTPS.KeyFile+"\"", " ")
			}
		}

		// configure server
		fs := http.FileServer(http.Dir(global.UIConfig.StaticFileRoot))
		http.Handle("/static/", http.StripPrefix("/static/", fs))
		http.HandleFunc("/", handleRoot)
		http.HandleFunc("/app/", handleApp)
		http.HandleFunc("/recovery/", handleRecovery)

		// start HTTP server
		if global.UIConfig.HTTP.Enabled {

			go http.ListenAndServe(":"+strconv.Itoa(global.UIConfig.HTTP.Port), nil)
		}

		// start HTTPS server
		if global.UIConfig.HTTPS.Enabled {

			go http.ListenAndServeTLS(":"+strconv.Itoa(global.UIConfig.HTTPS.Port), global.UIConfig.HTTPS.CertFile, global.UIConfig.HTTPS.KeyFile, nil)
		}

		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ipnet.IP.To4() != nil {

					if global.UIConfig.HTTP.Enabled {

						console.Log("http://"+ipnet.IP.String()+":"+strconv.Itoa(global.UIConfig.HTTP.Port), " ")
					}

					if global.UIConfig.HTTPS.Enabled {

						console.Log("https://"+ipnet.IP.String()+":"+strconv.Itoa(global.UIConfig.HTTPS.Port), " ")
					}
				}
			}
		}

		if global.UIConfig.Recovery {
			console.Log("WARNING! Recovery Mode ist enabled", " ")
		}
	} else {
		console.Log("no interface found", " ")
	}
}

func Exit() {

}
