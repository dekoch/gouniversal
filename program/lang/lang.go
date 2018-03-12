package lang

import (
	"encoding/json"
	"gouniversal/shared/config"
	"gouniversal/shared/io/file"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

const LangDir = "data/lang/"

// menu
type MenuProgram struct {
	Title string
}

type MenuAccount struct {
	Title  string
	Login  string
	Logout string
}

type Menu struct {
	Program MenuProgram
	Account MenuAccount
}

// home
type Home struct {
	Title string
}

// settings
type SettingsGeneral struct {
	Title string
	Apply string
}

type SettingsUserList struct {
	Title     string
	AddUser   string
	LoginName string
	Name      string
	Options   string
	Edit      string
}

type SettingsUserEdit struct {
	Title     string
	LoginName string
	Name      string
	State     string
	Language  string
	Password  string
	Groups    string
	Apply     string
	Delete    string
}

type SettingsUser struct {
	Title    string
	UserList SettingsUserList
	UserEdit SettingsUserEdit
}

type SettingsGroupList struct {
	Title    string
	AddGroup string
	Name     string
	Options  string
	Edit     string
}

type SettingsGroupEdit struct {
	Title   string
	Name    string
	State   string
	Comment string
	Pages   string
	Apply   string
	Delete  string
}

type SettingsGroup struct {
	Title     string
	GroupList SettingsGroupList
	GroupEdit SettingsGroupEdit
}

type SettingsAbout struct {
	Title string
}

type Settings struct {
	Title       string
	GeneralEdit SettingsGeneral
	User        SettingsUser
	Group       SettingsGroup
	About       SettingsAbout
}

type Login struct {
	Title    string
	User     string
	Password string
	Login    string
}

type Logout struct {
	Title string
}

type Exit struct {
	Title string
}

type File struct {
	Header   config.FileHeader
	Menu     Menu
	Home     Home
	Settings Settings
	Login    Login
	Logout   Logout
	Exit     Exit
}

type Global struct {
	Mut  sync.Mutex
	File []File
}

func SaveLang(lang File, n string) error {

	lang.Header = config.BuildHeader(n, "LangGlobal", 1.0, "Language File")

	if _, err := os.Stat(LangDir + n); os.IsNotExist(err) {
		// if not found, create default file
		lang.Menu.Program.Title = "Program"

		lang.Menu.Account.Title = "Account"
		lang.Menu.Account.Login = "Login"
		lang.Menu.Account.Logout = "Logout"

		lang.Home.Title = "Home"

		lang.Settings.Title = "Settings"
		lang.Settings.GeneralEdit.Title = "General"
		lang.Settings.GeneralEdit.Apply = "Apply"

		lang.Settings.User.Title = "User"
		lang.Settings.User.UserList.Title = "User List"
		lang.Settings.User.UserList.AddUser = "Add User"
		lang.Settings.User.UserList.LoginName = "Login Name"
		lang.Settings.User.UserList.Name = "Name"
		lang.Settings.User.UserList.Options = "Options"
		lang.Settings.User.UserList.Edit = "Edit"

		lang.Settings.User.UserEdit.Title = "User Edit"
		lang.Settings.User.UserEdit.LoginName = "Login Name"
		lang.Settings.User.UserEdit.Name = "Name"
		lang.Settings.User.UserEdit.State = "State"
		lang.Settings.User.UserEdit.Language = "Language"
		lang.Settings.User.UserEdit.Password = "Password"
		lang.Settings.User.UserEdit.Groups = "Groups"
		lang.Settings.User.UserEdit.Apply = "Apply"
		lang.Settings.User.UserEdit.Delete = "Delete"

		lang.Settings.Group.Title = "Group"
		lang.Settings.Group.GroupList.Title = "Group List"
		lang.Settings.Group.GroupList.AddGroup = "Add Group"
		lang.Settings.Group.GroupList.Name = "Name"
		lang.Settings.Group.GroupList.Options = "Options"
		lang.Settings.Group.GroupList.Edit = "Edit"

		lang.Settings.Group.GroupEdit.Title = "Group Edit"
		lang.Settings.Group.GroupEdit.Name = "Name"
		lang.Settings.Group.GroupEdit.State = "State"
		lang.Settings.Group.GroupEdit.Comment = "Comment"
		lang.Settings.Group.GroupEdit.Pages = "Pages"
		lang.Settings.Group.GroupEdit.Apply = "Apply"
		lang.Settings.Group.GroupEdit.Delete = "Delete"

		lang.Settings.About.Title = "About"

		lang.Login.Title = "Login"
		lang.Login.User = "User"
		lang.Login.Password = "Password"
		lang.Login.Login = "Login"

		lang.Logout.Title = "Logout"

		lang.Exit.Title = "Exit"
	}

	b, err := json.Marshal(lang)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(LangDir+n, b)

	return err
}

func LoadLang(n string) File {

	var lg File

	if _, err := os.Stat(LangDir + n); os.IsNotExist(err) {
		// if not found, create default file
		SaveLang(lg, n)
	}

	f := new(file.File)
	b, err := f.ReadFile(LangDir + n)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &lg)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(lg.Header, "LangGlobal") == false {
		log.Fatal("wrong config")
	}

	return lg
}

func LoadLangFiles() []File {

	lg := make([]File, 0)

	if _, err := os.Stat(LangDir + "en"); os.IsNotExist(err) {
		// if not found, create default file
		var newlg File
		SaveLang(newlg, "en")
	}

	files, err := ioutil.ReadDir(LangDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, fl := range files {

		var langfile File

		f := new(file.File)
		b, err := f.ReadFile(LangDir + fl.Name())
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(b, &langfile)
		if err != nil {
			log.Fatal(err)
		}

		if config.CheckHeader(langfile.Header, "LangGlobal") == false {
			log.Fatal("wrong config")
		} else {

			lg = append(lg, langfile)
		}

	}

	return lg
}
