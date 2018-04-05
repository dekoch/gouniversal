package ui

import (
	"encoding/json"
	"fmt"
	"gouniversal/modules"
	"gouniversal/program/global"
	"gouniversal/program/guestManagement"
	"gouniversal/program/lang"
	"gouniversal/program/programTypes"
	"gouniversal/program/ui/pageHome"
	"gouniversal/program/ui/pageLogin"
	"gouniversal/program/ui/settings"
	"gouniversal/program/ui/uifunc"
	"gouniversal/program/userManagement"
	"gouniversal/shared/alert"
	"gouniversal/shared/config"
	"gouniversal/shared/functions"
	"gouniversal/shared/io/file"
	"gouniversal/shared/navigation"
	"gouniversal/shared/types"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type UI struct{}

const UiConfigFile = "data/config/ui"

var (
	store      = new(sessions.CookieStore)
	cookieName string
	mod        = new(modules.Modules)
)

func SaveUiConfig(uc programTypes.UiConfig) error {

	uc.Header = config.BuildHeader("ui", "ui", 1.0, "UI config file")

	if _, err := os.Stat(UiConfigFile); os.IsNotExist(err) {
		// if not found, create default file
		uc.ProgramFileRoot = "data/ui/program/1.0/"
		uc.StaticFileRoot = "data/ui/static/1.0/"
		uc.Port = 8080
		uc.MaxGuests = 20
		uc.Recovery = false
	}

	b, err := json.Marshal(uc)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(UiConfigFile, b)

	return err
}

func LoadUiConfig() programTypes.UiConfig {

	var uc programTypes.UiConfig

	if _, err := os.Stat(UiConfigFile); os.IsNotExist(err) {
		// if not found, create default file
		SaveUiConfig(uc)
	}

	f := new(file.File)
	b, err := f.ReadFile(UiConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &uc)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(uc.Header, "ui") == false {
		log.Fatal("wrong config")
	}

	return uc
}

func setCookie(parameter string, value string, w http.ResponseWriter, r *http.Request) {
	//fmt.Println("set cookie " + parameter + " to " + value)

	session, _ := store.Get(r, cookieName)
	session.Values[parameter] = value
	session.Save(r, w)
}

func getCookie(parameter string, w http.ResponseWriter, r *http.Request) (string, error) {
	//fmt.Println("get cookie " + parameter)

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
	u, err := getCookie("authenticated", w, r)
	if err == nil {

		nav.User = userManagement.SelectUser(u)

		if nav.User.State < 0 {
			// if no user found
			nav.User = guestManagement.SelectGuest(u)
		}

		if nav.User.State < 0 {
			// if no guest found, create new
			nav.User = guestManagement.NewGuest()
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
	nav.User = guestManagement.NewGuest()
	nav.GodMode = false

	setSession(nav, w, r)
}

func renderProgram(page *types.Page, nav *navigation.Navigation) []byte {

	type program struct {
		Lang              lang.Menu
		Title             string
		MenuProgram       template.HTML
		MenuProgramHidden template.HTML
		MenuApp           template.HTML
		MenuAppHidden     template.HTML
		MenuAccountTitle  template.HTML
		MenuAccount       template.HTML
		UUID              template.HTML
		Content           template.HTML
	}
	var p program

	p.Title = nav.Sitemap.PageTitle(nav.Path)
	p.Lang = page.Lang.Menu

	p.MenuProgramHidden = "hidden"
	menuProgram := ""
	lastProgramDepth := -1

	p.MenuAppHidden = "hidden"
	menuApp := ""
	lastAppDepth := -1

	menuAccount := ""
	lastAccountDepth := -1

	path := ""
	title := ""
	depth := -1

	for i := len(nav.Sitemap.Pages) - 1; i >= 0; i-- {

		path = nav.Sitemap.Pages[i].Path
		title = nav.Sitemap.Pages[i].Title
		depth = nav.Sitemap.Pages[i].Depth

		// menuProgram
		if (strings.HasPrefix(path, "Program:") &&
			depth <= 2) ||
			(strings.HasPrefix(path, "App:Program:") &&
				depth <= 4) {

			if userManagement.IsPageAllowed(path, nav.User) ||
				nav.GodMode {

				if lastProgramDepth == -1 {
					lastProgramDepth = depth
				}

				// if menu depth changed, add divider
				if depth != lastProgramDepth {
					lastProgramDepth = depth

					menuProgram += "<div class=\"dropdown-divider\"></div>"
				}

				menuProgram += "<button class=\"dropdown-item\" type=\"submit\" name=\"navigation\" value=\"" + path + "\">" + title + "</button>"
				p.MenuProgramHidden = ""
			}
		}

		// menuApp
		if (strings.HasPrefix(path, "App:") &&
			depth <= 2) &&
			strings.HasPrefix(path, "App:Program:") == false &&
			strings.HasPrefix(path, "App:Account:") == false {

			if userManagement.IsPageAllowed(path, nav.User) ||
				nav.GodMode {

				if lastAppDepth == -1 {
					lastAppDepth = depth
				}

				// if menu depth changed, add divider
				if depth != lastAppDepth {
					lastAppDepth = depth

					menuApp += "<div class=\"dropdown-divider\"></div>"
				}

				menuApp += "<button class=\"dropdown-item\" type=\"submit\" name=\"navigation\" value=\"" + path + "\">" + title + "</button>"
				p.MenuAppHidden = ""
			}
		}

		// menuAccount
		if (strings.HasPrefix(path, "Account:") &&
			depth <= 2) ||
			(strings.HasPrefix(path, "App:Account:") &&
				depth <= 3) {

			if userManagement.IsPageAllowed(path, nav.User) ||
				nav.GodMode {

				if lastAccountDepth == -1 {
					lastAccountDepth = depth
				}

				// if menu depth changed, add divider
				if depth != lastAccountDepth {
					lastAccountDepth = depth

					menuAccount += "<div class=\"dropdown-divider\"></div>"
				}

				menuAccount += "<button class=\"dropdown-item\" type=\"submit\" name=\"navigation\" value=\"" + path + "\">" + title + "</button>"
			}
		}
	}

	// menuAccount
	if nav.User.UUID != "" {
		p.MenuAccountTitle = template.HTML(nav.User.LoginName)
	} else {
		p.MenuAccountTitle = template.HTML(page.Lang.Menu.Account.Title)
	}

	p.MenuProgram = template.HTML(menuProgram)
	p.MenuApp = template.HTML(menuApp)
	p.MenuAccount = template.HTML(menuAccount)
	p.UUID = template.HTML(nav.User.UUID)
	p.Content = template.HTML(page.Content)

	templ, err := template.ParseFiles(global.UiConfig.ProgramFileRoot + "program.html")
	if err != nil {
		fmt.Println(err)
	}

	return []byte(functions.TemplToString(templ, p))
}

func selectLang(l string) lang.File {

	global.Lang.Mut.Lock()
	defer global.Lang.Mut.Unlock()

	// search lang
	for i := 0; i < len(global.Lang.File); i++ {

		if l == global.Lang.File[i].Header.FileName {

			return global.Lang.File[i]
		}
	}

	// if nothing found
	// search "en"
	for i := 0; i < len(global.Lang.File); i++ {

		if "en" == global.Lang.File[i].Header.FileName {

			return global.Lang.File[i]
		}
	}

	// if nothing found
	// load or create "en"
	return lang.LoadLang("en")
}

func handleRoot(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/" {
		session, _ := store.Get(r, cookieName)
		if session.IsNew {
			initCookies(w, r)
		}

		http.Redirect(w, r, "/app/", http.StatusSeeOther)
	}

	fmt.Println(r.URL.Path)
}

func handleApp(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	page := new(types.Page)
	nav := new(navigation.Navigation)

	getSession(nav, w, r)

	page.Lang = selectLang(nav.User.Lang)

	pageHome.RegisterPage(page, nav)
	mod.RegisterPage(page, nav)
	settings.RegisterPage(page, nav)
	nav.Sitemap.Register("Program:Exit", page.Lang.Exit.Title)
	pageLogin.RegisterPage(page, nav)

	newPath := r.FormValue("navigation")

	if newPath != "" {
		nav.NavigatePath(newPath)

		fmt.Println(newPath)
	}

	nav.Redirect = "init"

	for functions.IsEmpty(nav.Redirect) == false {

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

				nav.RedirectPath("Account:Login", true)
			}
		} else if nav.IsNext("App") {

			mod.Render(page, nav, r)

		} else if nav.IsNext("Account") {

			if nav.IsNext("Login") {

				user := r.FormValue("user")
				pwd := r.FormValue("pwd")

				if uifunc.CheckLogin(user, pwd) {

					nav.User = userManagement.SelectUser(uifunc.LoginNameToUUID(user))
					nav.RedirectPath("Program:Home", false)
				} else {
					pageLogin.Render(page, nav, r)
				}

			} else if nav.IsNext("Logout") {

				nav.User = guestManagement.NewGuest()
				nav.RedirectPath("Account:Login", false)
			}

		} else {
			nav.RedirectPath("Account:Login", true)
		}

		if nav.Redirect != "" {

			nav.NavigatePath(nav.Redirect)
		}
	}

	if nav.Path == "Program:Exit" {

		page.Content = "goodbye"

		global.Console.Mut.Lock()
		global.Console.Input = "exit"
		global.Console.Mut.Unlock()
	} else {

		setSession(nav, w, r)
	}

	w.Write(renderProgram(page, nav))

	// show all pages
	//nav.Sitemap.ShowMap()

	// show allowed pages
	/*var sm []string
	sm = nav.Sitemap.PageList()
	fmt.Println("####")
	for i := 0; i < len(sm); i++ {
		if userManagement.IsPageAllowed(sm[i], nav.User) {
			fmt.Println(sm[i])
		}
	}
	fmt.Println("####")*/

	t := time.Now()
	elapsed := t.Sub(start)
	f := elapsed.Seconds() * 1000.0
	fmt.Println(nav.Path + " " + strconv.FormatFloat(f, 'f', 1, 64) + "ms")
}

func handleRecovery(w http.ResponseWriter, r *http.Request) {

	if global.UiConfig.Recovery {

		button := r.FormValue("recovery")

		if button == "goback" {

			http.Redirect(w, r, "/", http.StatusSeeOther)

		} else if button == "disablerecovery" {

			global.UiConfig.Recovery = false
			SaveUiConfig(global.UiConfig)

		} else if button == "cookies" {

			initCookies(w, r)

		} else if button == "god" {

			setCookie("isgod", "true", w, r)

		} else {
			fmt.Println(button)
		}

		type recovery struct {
			Temp string
		}
		var re recovery

		re.Temp = ""

		templ, err := template.ParseFiles(global.UiConfig.ProgramFileRoot + "recovery.html")
		if err != nil {
			fmt.Println(err)
		}
		templ.Execute(w, re)
	} else {
		w.Write([]byte("recovery mode is DISABLED!"))
	}
}

func (ui *UI) StartServer() {
	fmt.Println("starting webserver...")

	fmt.Println("ProgramFileRoot: " + global.UiConfig.ProgramFileRoot)
	fmt.Println("StaticFileRoot: " + global.UiConfig.StaticFileRoot)

	if _, err := os.Stat(global.UiConfig.ProgramFileRoot); os.IsNotExist(err) {
		// if not found, exit program
		fmt.Println("error: ProgramFileRoot not found")
		os.Exit(1)
	}

	if _, err := os.Stat(global.UiConfig.StaticFileRoot); os.IsNotExist(err) {
		// if not found, exit program
		fmt.Println("error: StaticFileRoot not found")
		os.Exit(1)
	}

	addrs, err := net.InterfaceAddrs()

	if err == nil {

		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ipnet.IP.To4() != nil {

					fmt.Println("http://" + ipnet.IP.String() + ":" + strconv.Itoa(global.UiConfig.Port))
				}
			}
		}

		global.Lang.File = lang.LoadLangFiles()

		// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
		key := securecookie.GenerateRandomKey(32)
		store = sessions.NewCookieStore(key)

		u := uuid.Must(uuid.NewRandom())
		cookieName = u.String()

		alert.Start()
		mod.LoadConfig()

		// configure server
		fs := http.FileServer(http.Dir(global.UiConfig.StaticFileRoot))
		http.Handle("/static/", http.StripPrefix("/static/", fs))
		http.HandleFunc("/", handleRoot)
		http.HandleFunc("/app/", handleApp)
		http.HandleFunc("/recovery/", handleRecovery)

		http.ListenAndServe(":"+strconv.Itoa(global.UiConfig.Port), nil)
	} else {
		fmt.Println("no interface found")
	}
}

func (ui *UI) Exit() {
	mod.Exit()
}
