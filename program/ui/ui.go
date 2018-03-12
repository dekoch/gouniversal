package ui

import (
	"encoding/json"
	"fmt"
	"gouniversal/program/global"
	"gouniversal/program/lang"
	"gouniversal/program/programTypes"
	"gouniversal/program/ui/pageHome"
	"gouniversal/program/ui/pageLogin"
	"gouniversal/program/ui/pageSettings"
	"gouniversal/program/ui/uifunc"
	"gouniversal/program/userManagement"
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
)

func SaveUiConfig(uc programTypes.UiConfig) error {

	uc.Header = config.BuildHeader("ui", "ui", 1.0, "UI config file")

	if _, err := os.Stat(UiConfigFile); os.IsNotExist(err) {
		// if not found, create default file

		uc.FileRoot = "data/ui/1.0/"
		uc.Port = 8080
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

	path, err := getCookie("navigation", w, r)

	if err == nil {
		nav.Path = path
	}

	u, err := getCookie("authenticated", w, r)

	if err == nil {

		nav.User = userManagement.SelectUser(u)
	}

	g, err := getCookie("isgod", w, r)

	nav.GodMode = false

	if err == nil {
		if g == "true" {
			nav.GodMode = true
		}
	}
}

func initCookies(w http.ResponseWriter, r *http.Request) {
	nav := new(navigation.Navigation)
	nav.Path = "Program:Home"
	nav.User.UUID = ""
	nav.GodMode = false

	setSession(nav, w, r)
}

func renderProgram(page *types.Page, nav *navigation.Navigation) []byte {

	type program struct {
		Lang        lang.Menu
		Title       string
		MenuProgram template.HTML
		MenuAccount template.HTML
		Content     template.HTML
	}
	var p program

	p.Title = nav.Sitemap.PageTitle(nav.Path)
	p.Lang = page.Lang.Menu

	menuprogram := ""
	var depth int
	lastDepth := 1
	for i := len(nav.Sitemap.Pages) - 1; i >= 0; i-- {

		depth = nav.Sitemap.Pages[i].Depth

		if strings.HasPrefix(nav.Sitemap.Pages[i].Path, "Program:") &&
			depth <= 2 {

			if userManagement.IsPageAllowed(nav.Sitemap.Pages[i].Path, nav.User) ||
				nav.GodMode {

				if depth != lastDepth {
					lastDepth = depth

					menuprogram += "<div class=\"dropdown-divider\"></div>"
				}

				menuprogram += "<button class=\"dropdown-item\" type=\"submit\" name=\"navigation\" value=\"" + nav.Sitemap.Pages[i].Path + "\">" + nav.Sitemap.Pages[i].Title + "</button>"
			}
		}
	}

	p.MenuProgram = template.HTML(menuprogram)
	p.Content = template.HTML(page.Content)

	templ, err := template.ParseFiles(global.UiConfig.FileRoot + "program/program.html")
	if err != nil {
		fmt.Println(err)
	}

	return []byte(functions.TemplToString(templ, p))
}

func selectLang(l string) lang.File {

	global.Lang.Mut.Lock()
	defer global.Lang.Mut.Unlock()

	for i := 0; i < len(global.Lang.File); i++ {

		if l != "" {

			if l == global.Lang.File[i].Header.FileName {

				return global.Lang.File[i]
			}
		} else {

			if "en" == global.Lang.File[i].Header.FileName {

				return global.Lang.File[i]
			}
		}
	}

	return lang.LoadLang("en")
}

func handleroot(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/" {
		session, _ := store.Get(r, cookieName)
		if session.IsNew {
			initCookies(w, r)
		}

		http.Redirect(w, r, "/app/", http.StatusSeeOther)
	}

	fmt.Println(r.URL.Path)
}

func handleapp(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	page := new(types.Page)
	nav := new(navigation.Navigation)

	getSession(nav, w, r)

	page.Lang = selectLang(nav.User.Lang)

	pageHome.RegisterPage(page, nav)
	pageSettings.RegisterPage(page, nav)
	pageLogin.RegisterPage(page, nav)
	nav.Sitemap.Register("Program:Logout", page.Lang.Logout.Title)
	nav.Sitemap.Register("Program:Exit", page.Lang.Exit.Title)

	var strNavigation string
	strNavigation = r.FormValue("navigation")

	if strNavigation != "" {
		nav.NavigatePath(strNavigation)

		fmt.Println(strNavigation)
	}

	nav.Redirect = "init"

	for nav.Redirect != "" {

		nav.CurrentPath = ""
		nav.Redirect = ""
		page.Content = ""

		if nav.IsNext("Program") {

			if nav.IsNext("Home") {

				pageHome.Render(page, nav, r)

			} else if nav.IsNext("Settings") {

				pageSettings.Render(page, nav, r)

			} else if nav.IsNext("Login") {

				user := r.FormValue("user")
				pwd := r.FormValue("pwd")

				if uifunc.CheckLogin(user, pwd) {

					nav.User = userManagement.SelectUser(uifunc.LoginNameToUUID(user))
					nav.RedirectPath("Program:Home", false)
				} else {
					pageLogin.Render(page, nav, r)
				}

			} else if nav.IsNext("Logout") {

				nav.User = userManagement.SelectUser("")
				nav.RedirectPath("Program:Home", false)

			} else if nav.IsNext("Exit") {
			} else {

				nav.RedirectPath("Program:Login", true)
			}
		} else {
			nav.RedirectPath("Program:Home", true)
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
	fmt.Print(nav.Path)
	fmt.Println(" " + strconv.FormatFloat(f, 'f', 1, 64) + "ms")
}

func handleRecovery(w http.ResponseWriter, r *http.Request) {

	if global.UiConfig.Recovery {

		var strButton string
		strButton = r.FormValue("recovery")

		if strButton == "goback" {

			http.Redirect(w, r, "/", http.StatusSeeOther)

		} else if strButton == "disablerecovery" {

			global.UiConfig.Recovery = false
			SaveUiConfig(global.UiConfig)

		} else if strButton == "cookies" {

			initCookies(w, r)

		} else if strButton == "god" {

			setCookie("isgod", "true", w, r)

		} else {
			fmt.Println(strButton)
		}

		type recovery struct {
			Temp string
		}
		var re recovery

		re.Temp = ""

		templ, err := template.ParseFiles(global.UiConfig.FileRoot + "program/recovery.html")
		if err != nil {
			fmt.Println(err)
		}
		templ.Execute(w, re)
	} else {
		w.Write([]byte("recovery mode is DISABLED!"))
	}
}

func (ui UI) StartServer() {
	fmt.Println("starting webserver...")

	fmt.Print("FileRoot: ")
	fmt.Println(global.UiConfig.FileRoot)

	if _, err := os.Stat(global.UiConfig.FileRoot); os.IsNotExist(err) {
		// if not found, exit program
		fmt.Println("error: FileRoot not found")
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

		// configure server
		fs := http.FileServer(http.Dir(global.UiConfig.FileRoot + "static/"))
		http.Handle("/static/", http.StripPrefix("/static/", fs))
		http.HandleFunc("/", handleroot)
		http.HandleFunc("/app/", handleapp)
		http.HandleFunc("/recovery/", handleRecovery)

		http.ListenAndServe(":"+strconv.Itoa(global.UiConfig.Port), nil)
	} else {
		fmt.Println("no interface found")
	}
}
