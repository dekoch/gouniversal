package ui

import (
	"fmt"
	"gouniversal/modules"
	"gouniversal/program/global"
	"gouniversal/program/guestManagement"
	"gouniversal/program/lang"
	"gouniversal/program/ui/pageHome"
	"gouniversal/program/ui/pageLogin"
	"gouniversal/program/ui/settings"
	"gouniversal/program/ui/uifunc"
	"gouniversal/program/userManagement"
	"gouniversal/shared/alert"
	"gouniversal/shared/functions"
	"gouniversal/shared/navigation"
	"gouniversal/shared/sitemap"
	"gouniversal/shared/types"
	"html/template"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

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
)

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
		Title     string
		MenuLeft  template.HTML
		MenuRight template.HTML
		UUID      template.HTML
		Content   template.HTML
	}
	var p program

	p.Title = nav.Sitemap.PageTitle(nav.Path)

	// remove leading numbers
	if strings.Contains(p.Title, ".") {
		s := strings.SplitAfterN(p.Title, ".", -1)
		p.Title = s[1]
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

					mDropdown = append(newDropdown, mDropdown...)
				}

				// add items to dropdown
				for d := 0; d < len(mDropdown); d++ {

					if pages[i].Menu == mDropdown[d].Title {

						newItem := make([]sitemap.Page, 1)
						newItem[0] = pages[i]

						mDropdown[d].Items = append(newItem, mDropdown[d].Items...)
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

		// remove leading numbers
		if strings.Contains(dropDownTitle, ".") {
			s := strings.SplitAfterN(dropDownTitle, ".", -1)
			dropDownTitle = s[1]
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

			itemTitle := mDropdown[d].Items[i].Title
			// remove leading numbers
			if strings.Contains(itemTitle, ".") {
				s := strings.SplitAfterN(itemTitle, ".", -1)
				itemTitle = s[1]
			}

			dropDown += "<button class=\"dropdown-item\" type=\"submit\" name=\"navigation\" value=\"" + mDropdown[d].Items[i].Path + "\">" + itemTitle + "</button>\n"
		}

		dropDown += "</div>\n</li>\n"

		if alignRight {

			htmlMenuRight += dropDown
		} else {

			htmlMenuLeft += dropDown
		}
	}

	p.MenuLeft = template.HTML(htmlMenuLeft)
	p.MenuRight = template.HTML(htmlMenuRight)
	p.UUID = template.HTML(nav.User.UUID)
	p.Content = template.HTML(page.Content)

	templ, err := template.ParseFiles(global.UiConfig.File.ProgramFileRoot + "program.html")
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
	modules.RegisterPage(page, nav)
	settings.RegisterPage(page, nav)
	nav.Sitemap.Register("Program", "Program:Exit", page.Lang.Exit.Title)
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

			modules.Render(page, nav, r)

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
			fmt.Println(button)
		}

		type recovery struct {
			Temp string
		}
		var re recovery

		re.Temp = ""

		templ, err := template.ParseFiles(global.UiConfig.File.ProgramFileRoot + "recovery.html")
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

	fmt.Println("ProgramFileRoot: " + global.UiConfig.File.ProgramFileRoot)
	fmt.Println("StaticFileRoot: " + global.UiConfig.File.StaticFileRoot)

	if _, err := os.Stat(global.UiConfig.File.ProgramFileRoot); os.IsNotExist(err) {
		// if not found, exit program
		fmt.Println("error: ProgramFileRoot not found")
		os.Exit(1)
	}

	if _, err := os.Stat(global.UiConfig.File.StaticFileRoot); os.IsNotExist(err) {
		// if not found, exit program
		fmt.Println("error: StaticFileRoot not found")
		os.Exit(1)
	}

	addrs, err := net.InterfaceAddrs()

	if err == nil {

		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ipnet.IP.To4() != nil {

					fmt.Println("http://" + ipnet.IP.String() + ":" + strconv.Itoa(global.UiConfig.File.Port))
				}
			}
		}

		// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
		key := securecookie.GenerateRandomKey(32)
		store = sessions.NewCookieStore(key)

		u := uuid.Must(uuid.NewRandom())
		cookieName = u.String()

		global.UserConfig.LoadConfig()
		global.GroupConfig.LoadConfig()

		alert.Start()
		modules.LoadConfig()

		// configure server
		fs := http.FileServer(http.Dir(global.UiConfig.File.StaticFileRoot))
		http.Handle("/static/", http.StripPrefix("/static/", fs))
		http.HandleFunc("/", handleRoot)
		http.HandleFunc("/app/", handleApp)
		http.HandleFunc("/recovery/", handleRecovery)

		http.ListenAndServe(":"+strconv.Itoa(global.UiConfig.File.Port), nil)
	} else {
		fmt.Println("no interface found")
	}
}

func (ui *UI) Exit() {

}
