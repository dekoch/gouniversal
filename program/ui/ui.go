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

	"github.com/dekoch/gouniversal/modules"
	"github.com/dekoch/gouniversal/program/global"
	"github.com/dekoch/gouniversal/program/guestManagement"
	"github.com/dekoch/gouniversal/program/ui/pageHome"
	"github.com/dekoch/gouniversal/program/ui/pageLogin"
	"github.com/dekoch/gouniversal/program/ui/settings"
	"github.com/dekoch/gouniversal/program/ui/uifunc"
	"github.com/dekoch/gouniversal/program/userManagement"
	"github.com/dekoch/gouniversal/shared/alert"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/sitemap"
	"github.com/dekoch/gouniversal/shared/types"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type UI struct{}

type menuDropdown struct {
	Order int
	Title string
	Items []sitemap.Page
}

var (
	store      = new(sessions.CookieStore)
	cookieName string
	guest      guestManagement.GuestManagement
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
			nav.User = guest.SelectGuest(uid)
		}

		if nav.User.State < 0 {
			// if no guest found, create new
			nav.User = guest.NewGuest()
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
	nav.Path = "Account:Login"
	nav.User = guest.NewGuest()
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
		Content   template.HTML
	}
	var c content

	c.MenuBrand = nav.Sitemap.PageTitle(nav.Path)

	// title
	if global.UiConfig.File.Title != "" {
		c.Title = global.UiConfig.File.Title + " - " + c.MenuBrand
	} else {
		c.Title = c.MenuBrand
	}

	pages := nav.Sitemap.Pages
	mDropdown := make([]menuDropdown, 0)

	// put each allowed page into menu slice
	for i := 0; i < len(pages); i++ {

		if pages[i].Menu != "" {

			if userManagement.IsPageAllowed(pages[i].Path, nav.User) ||
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

	p, err := functions.PageToString(global.UiConfig.File.ProgramFileRoot+"program.html", c)
	if err != nil {
		console.Log(err, "ui.go")
		p = err.Error()
	}

	return []byte(p)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/" {
		session, _ := store.Get(r, cookieName)
		if session.IsNew {
			initCookies(w, r)
		}

		// redirect to HTTPS if enabled
		if global.UiConfig.File.HTTPS.Enabled {
			host := strings.Split(r.Host, ":")
			port := strconv.Itoa(global.UiConfig.File.HTTPS.Port)

			http.Redirect(w, r, "https://"+host[0]+":"+port+"/app/", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/app/", http.StatusSeeOther)
		}

	} else if r.URL.Path == "/favicon.ico" {
		f := global.UiConfig.File.StaticFileRoot + "favicon.ico"

		if _, err := os.Stat(f); os.IsNotExist(err) == false {
			http.ServeFile(w, r, f)
		}
	}

	console.Output(r.URL.Path, "")
}

func handleApp(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	page := new(types.Page)
	nav := new(navigation.Navigation)

	getSession(nav, w, r)

	global.Lang.SelectLang(nav.User.Lang, &page.Lang)

	pageHome.RegisterPage(page, nav)
	modules.RegisterPage(page, nav)
	settings.RegisterPage(page, nav)
	nav.Sitemap.Register("Program", "Program:Exit", page.Lang.Exit.Title)
	pageLogin.RegisterPage(page, nav)

	newPath := r.FormValue("navigation")

	if newPath != "" {
		nav.NavigatePath(newPath)

		console.Output(newPath, "")
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
		if nav.IsNext("Program") {

			if nav.IsNext("Home") {

				pageHome.Render(page, nav, r)

			} else if nav.IsNext("Settings") {

				settings.Render(page, nav, r)

			} else if nav.IsNext("Exit") {
			} else {
				nav.RedirectPath("404", true)
			}
		} else if nav.IsNext("App") {

			modules.Render(page, nav, r)

		} else if nav.IsNext("Account") {

			if nav.IsNext("Login") {

				name := r.FormValue("name")
				pwd := r.FormValue("pwd")
				maxAttempts := guest.MaxLoginAttempts(nav.User.UUID)

				if uifunc.CheckLogin(name, pwd) && maxAttempts == false {

					var err error
					nav.User, err = global.UserConfig.Get(uifunc.LoginNameToUUID(name))
					if err != nil {
						pageLogin.Render(page, nav, r)
					} else {
						nav.RedirectPath("Program:Home", false)
					}
				} else {
					pageLogin.Render(page, nav, r)
				}

			} else if nav.IsNext("Logout") {

				nav.User = guest.NewGuest()
				nav.RedirectPath("Account:Login", false)
			} else {
				nav.RedirectPath("404", true)
			}

		} else if nav.Path == "404" {
			page.Content = "<h1>404</h1><br>"
			page.Content += page.Lang.Error.NotFound404
		} else if nav.Path == "508" {
			page.Content = "<h1>508</h1><br>"
			page.Content += page.Lang.Error.LoopDetected508
		} else {
			nav.RedirectPath("404", true)
		}

		if loopDetection >= 5 {
			loopDetection = 0

			nav.RedirectPath("508", true)
		}

		if nav.Redirect != "" {

			if nav.Redirect == "404" ||
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
		if userManagement.IsPageAllowed(sm[i], nav.User) {
			console.Output(sm[i])
		}
	}
	console.Output("####")*/

	t := time.Now()
	elapsed := t.Sub(start)
	f := elapsed.Seconds() * 1000.0
	console.Output(nav.Path+" "+strconv.FormatFloat(f, 'f', 1, 64)+"ms", "")
}

func handleRecovery(w http.ResponseWriter, r *http.Request) {

	if global.UiConfig.File.Recovery {

		button := r.FormValue("recovery")

		if button == "goback" {

			http.Redirect(w, r, "/", http.StatusSeeOther)

		} else if button == "disablerecovery" {

			global.UiConfig.File.Recovery = false
			global.UiConfig.SaveUiConfig()

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

		templ, err := template.ParseFiles(global.UiConfig.File.ProgramFileRoot + "recovery.html")
		if err != nil {
			console.Log(err, "")
		}
		templ.Execute(w, re)
	} else {
		w.Write([]byte("recovery mode is DISABLED!"))
	}
}

func (ui *UI) StartServer() {
	console.Log("starting webserver...", " ")

	console.Log("ProgramFileRoot: "+global.UiConfig.File.ProgramFileRoot, " ")
	console.Log("StaticFileRoot: "+global.UiConfig.File.StaticFileRoot, " ")

	if _, err := os.Stat(global.UiConfig.File.ProgramFileRoot); os.IsNotExist(err) {
		// if not found, exit program
		console.Log("error: ProgramFileRoot not found", "ui.go")
	}

	if _, err := os.Stat(global.UiConfig.File.StaticFileRoot); os.IsNotExist(err) {
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
		modules.LoadConfig()

		// if HTTPS is enabled, check cert and key file
		if global.UiConfig.File.HTTPS.Enabled {
			if _, err = os.Stat(global.UiConfig.File.HTTPS.CertFile); os.IsNotExist(err) {
				global.UiConfig.File.HTTPS.Enabled = false
				console.Log("missing CertFile \""+global.UiConfig.File.HTTPS.CertFile+"\"", " ")
			}

			if _, err = os.Stat(global.UiConfig.File.HTTPS.KeyFile); os.IsNotExist(err) {
				global.UiConfig.File.HTTPS.Enabled = false
				console.Log("missing KeyFile \""+global.UiConfig.File.HTTPS.KeyFile+"\"", " ")
			}
		}

		// configure server
		fs := http.FileServer(http.Dir(global.UiConfig.File.StaticFileRoot))
		http.Handle("/static/", http.StripPrefix("/static/", fs))
		http.HandleFunc("/", handleRoot)
		http.HandleFunc("/app/", handleApp)
		http.HandleFunc("/recovery/", handleRecovery)

		// start HTTP server
		if global.UiConfig.File.HTTP.Enabled {

			go http.ListenAndServe(":"+strconv.Itoa(global.UiConfig.File.HTTP.Port), nil)
		}

		// start HTTPS server
		if global.UiConfig.File.HTTPS.Enabled {

			go http.ListenAndServeTLS(":"+strconv.Itoa(global.UiConfig.File.HTTPS.Port), global.UiConfig.File.HTTPS.CertFile, global.UiConfig.File.HTTPS.KeyFile, nil)
		}

		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ipnet.IP.To4() != nil {

					if global.UiConfig.File.HTTP.Enabled {

						console.Log("http://"+ipnet.IP.String()+":"+strconv.Itoa(global.UiConfig.File.HTTP.Port), " ")
					}

					if global.UiConfig.File.HTTPS.Enabled {

						console.Log("https://"+ipnet.IP.String()+":"+strconv.Itoa(global.UiConfig.File.HTTPS.Port), " ")
					}
				}
			}
		}
	} else {
		console.Log("no interface found", " ")
	}
}

func (ui *UI) Exit() {

}
