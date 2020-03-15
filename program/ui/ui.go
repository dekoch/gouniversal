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
	"github.com/dekoch/gouniversal/program/ui/pageuseracct"
	"github.com/dekoch/gouniversal/program/ui/settings"
	"github.com/dekoch/gouniversal/program/userconfig"
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

func notFound(message interface{}, w http.ResponseWriter) {

	console.Log(message, "Error 404")

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
}

func internalError(message interface{}, w http.ResponseWriter) {

	console.Log(message, "Error 500")

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusInternalServerError)
}

func isHTTPS(r *http.Request) bool {
	return r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
}

func redirectToApp(w http.ResponseWriter, r *http.Request) bool {

	host := strings.Split(r.Host, ":")
	hostName := ""

	if len(host) > 0 {
		hostName = host[0]
	}

	// redirect to HTTPS if enabled
	if (isHTTPS(r) == false && global.UIConfig.HTTPS.Enabled) ||
		strings.HasPrefix(hostName, "www.") ||
		strings.HasPrefix(r.URL.Path, "/app/") == false {

		hostName = strings.Replace(hostName, "www.", "", 1)
		hostRedi := "http://" + hostName
		port := strconv.Itoa(global.UIConfig.HTTP.Port)

		if global.UIConfig.HTTPS.Enabled {

			port = strconv.Itoa(global.UIConfig.HTTPS.Port)
			hostRedi = "https://" + hostName
		}

		http.Redirect(w, r, hostRedi+":"+port+"/app/", http.StatusSeeOther)

		return true
	}

	return false
}

func setCookie(parameter string, value string, w http.ResponseWriter, r *http.Request) error {

	session, err := store.Get(r, cookieName)
	if err != nil {
		return err
	}

	session.Values[parameter] = value

	return session.Save(r, w)
}

func getCookie(parameter string, w http.ResponseWriter, r *http.Request) (string, error) {

	session, err := store.Get(r, cookieName)
	if err != nil {
		return "", err
	}

	return session.Values[parameter].(string), nil
}

func setSession(nav *navigation.Navigation, w http.ResponseWriter, r *http.Request) error {

	err := setCookie("navigation", nav.Path, w, r)
	if err != nil {
		return err
	}

	err = setCookie("authenticated", nav.User.UUID, w, r)
	if err != nil {
		return err
	}

	isGod := "false"

	if nav.GodMode {
		isGod = "true"
	}

	return setCookie("isgod", isGod, w, r)
}

func getSession(nav *navigation.Navigation, w http.ResponseWriter, r *http.Request) error {

	var (
		err     error
		session *sessions.Session
		uid     string
		isGod   string
	)

	func() {

		for i := 0; i <= 8; i++ {

			switch i {
			case 0:
				session, err = store.Get(r, cookieName)

			case 1:
				if session.IsNew {
					err = initCookies(nav, w, r)
				}

			case 2:
				// get stored path
				nav.Path, err = getCookie("navigation", w, r)

			case 3:
				// get stored user UUID
				uid, err = getCookie("authenticated", w, r)

			case 4:
				nav.User, _ = global.UserConfig.Get(uid)

			case 5:
				if nav.User.State < userconfig.StatePublic {
					// if no user found
					nav.User = global.Guests.SelectGuest(uid)
				}

			case 6:
				if nav.User.State < userconfig.StatePublic {
					// if no guest found, create new
					var guest userconfig.User
					guest, _ = global.UserConfig.GetWithState(userconfig.StatePublic)
					nav.User = global.Guests.NewGuest(guest, global.UIConfig.MaxGuests)
				}

			case 7:
				// for debugging
				nav.GodMode = false
				isGod, err = getCookie("isgod", w, r)

			case 8:
				if isGod == "true" {
					nav.GodMode = true
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

func initCookies(nav *navigation.Navigation, w http.ResponseWriter, r *http.Request) error {

	u, _ := global.UserConfig.GetWithState(userconfig.StatePublic)
	nav.User = global.Guests.NewGuest(u, global.UIConfig.MaxGuests)
	nav.Path = "init"
	nav.GodMode = false

	return setSession(nav, w, r)
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
		Link      template.HTML
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

	var mDropdown []menuDropdown

	// put each allowed page into menu slice
	for _, page := range nav.Sitemap.GetPages() {

		if page.Menu == "" {
			continue
		}

		if usermanagement.IsPageAllowed(page.Path, nav.User) == false &&
			nav.GodMode == false {

			continue
		}

		dropdownFound := false

		for d := 0; d < len(mDropdown); d++ {

			if page.Menu == mDropdown[d].Title {

				mDropdown[d].Items = append(mDropdown[d].Items, page)

				dropdownFound = true
			}
		}

		if dropdownFound {
			continue
		}

		// create new dropdown
		var newDropdown menuDropdown
		newDropdown.Title = page.Menu

		if page.MenuOrder > -1 {

			newDropdown.Order = page.MenuOrder
		} else {
			// set predefined order
			switch newDropdown.Title {
			case "Program":
				newDropdown.Order = 0

			case "App":
				newDropdown.Order = 1

			case "Account":
				newDropdown.Order = 999

			default:
				newDropdown.Order = 999 - len(mDropdown)
			}
		}

		newDropdown.Items = append(newDropdown.Items, page)

		mDropdown = append(mDropdown, newDropdown)
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

			dropDownTitle = page.Lang.Menu.Account.Title

			if functions.IsEmpty(nav.User.Name) == false {

				dropDownTitle = nav.User.Name

			} else if functions.IsEmpty(nav.User.LoginName) == false {

				dropDownTitle = nav.User.LoginName
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

	link := ""

	if nav.UIConfig.HTTPS.Enabled {
		link = "https://"
	} else {
		link = "http://"
	}

	link += nav.Server.Host + "?path=" + nav.Path
	c.Link = template.HTML("<!-- " + link + " -->")

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

		// set new path to cookie
		if newPath != "" {

			nav := new(navigation.Navigation)

			err := getSession(nav, w, r)
			if err != nil {
				internalError(err, w)
				return
			}

			nav.NavigatePath(newPath)

			err = setSession(nav, w, r)
			if err != nil {
				internalError(err, w)
				return
			}

			console.Output(newPath, "app")
		}

		redirectToApp(w, r)

		return
	}

	if r.URL.Path == "/favicon.ico" {

		icon := global.UIConfig.StaticFileRoot + "favicon.ico"

		if _, err := os.Stat(icon); os.IsNotExist(err) == false {

			http.ServeFile(w, r, icon)
			return
		}
	} else {

		if redirectToApp(w, r) {
			return
		}
	}

	notFound("\""+r.URL.Path+"\" ("+clientinfo.String(r)+")", w)
}

func handleApp(w http.ResponseWriter, r *http.Request) {

	if redirectToApp(w, r) {
		return
	}

	start := time.Now()

	page := new(types.Page)
	nav := new(navigation.Navigation)

	nav.UIConfig = global.UIConfig

	nav.Server.Host = r.Host
	nav.Server.URLPath = r.URL.Path

	nav.ResponseWriter = w

	err := getSession(nav, w, r)
	if err != nil {
		internalError(err, w)
		return
	}

	newPath := r.FormValue("navigation")

	if newPath != "" {
		nav.NavigatePath(newPath)

		console.Output(newPath, "App")
	}

	nav.Redirect = "init"
	loopDetection := 0

	for functions.IsEmpty(nav.Redirect) == false {

		loopDetection++

		nav.CurrentPath = ""
		nav.Redirect = ""
		page.Content = ""

		global.Lang.SelectLang(nav.User.Lang, &page.Lang)

		nav.Sitemap.Clear()
		pagehome.RegisterPage(page, nav)
		module.RegisterPage(page, nav)
		settings.RegisterPage(page, nav)
		nav.Sitemap.Register("Program", "Program:Exit", page.Lang.Exit.Title)
		pagelogin.RegisterPage(page, nav)
		pageuseracct.RegisterPage(page, nav)

		// select first allowed page
		if nav.Path == "init" {

			for _, p := range nav.Sitemap.GetPages() {

				if p.Menu != "" && nav.Path == "init" {

					if usermanagement.IsPageAllowed(p.Path, nav.User) ||
						nav.GodMode {

						nav.Path = p.Path
					}
				}
			}
		}

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
				if nav.User.LoginName != "" {
					console.Log("\""+nav.User.LoginName+"\" logged out", "Logout")
				}

				u, _ := global.UserConfig.GetWithState(userconfig.StatePublic)
				nav.User = global.Guests.NewGuest(u, global.UIConfig.MaxGuests)

				nav.RedirectPath("Account:Login", false)

			case "UserAccount":
				pageuseracct.Render(page, nav, r)

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
				// prevent redirect loops through form input
				var f http.Request
				r = &f
				nav.NavigatePath(nav.Redirect)
			}
		}
	}

	if nav.Path == "Program:Exit" {

		page.Content = "goodbye"

		console.Input("exit")
	} else {

		err = setSession(nav, w, r)
		if err != nil {
			internalError(err, w)
			return
		}
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

	if global.UIConfig.Recovery == false {
		notFound("Recovery Mode is DISABLED!", w)
		return
	}

	switch r.FormValue("recovery") {
	case "goback":
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case "disablerecovery":
		global.UIConfig.Recovery = false
		global.UIConfig.SaveConfig()

	case "cookies":
		nav := new(navigation.Navigation)
		initCookies(nav, w, r)

	case "god":
		setCookie("isgod", "true", w, r)

	default:
		console.Log(r.FormValue("recovery"), "ui.go")
	}

	type recovery struct {
		Temp string
	}
	var re recovery

	re.Temp = ""

	templ, err := template.ParseFiles(global.UIConfig.ProgramFileRoot + "recovery.html")
	if err != nil {
		internalError(err, w)
	}
	templ.Execute(w, re)
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
	if err != nil {
		console.Log("no interface found", " ")
		return
	}

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
}

func Exit() {

}

func GetUserUUIDList() []string {

	return global.UserConfig.GetUUIDList()
}
