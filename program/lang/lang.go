package lang

import "github.com/dekoch/gouniversal/shared/config"

// menu
type MenuProgram struct {
	Title string
}

type MenuApp struct {
	Title string
}

type MenuAccount struct {
	Title  string
	Login  string
	Logout string
}

type Menu struct {
	Program MenuProgram
	App     MenuApp
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

type SettingsUserState struct {
	Public   string
	Active   string
	Inactive string
}

type SettingsUserEdit struct {
	Title     string
	LoginName string
	Name      string
	State     string
	States    SettingsUserState
	Language  string
	Password  string
	Comment   string
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
	Comment  string
	Options  string
	Edit     string
}

type SettingsGroupState struct {
	Active   string
	Inactive string
}

type SettingsGroupEdit struct {
	Title   string
	Name    string
	State   string
	States  SettingsGroupState
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

type Alert struct {
	Success string
	Info    string
	Warning string
	Error   string
}

type Error struct {
	CE400BadRequest   string
	CE404NotFound     string
	SE508LoopDetected string
}

type LangFile struct {
	Header   config.FileHeader
	Menu     Menu
	Home     Home
	Settings Settings
	Login    Login
	Logout   Logout
	Exit     Exit
	Alert    Alert
	Error    Error
}

func DefaultEn() LangFile {

	var l LangFile

	l.Header = config.BuildHeader("en", "LangProgram", 1.0, "Language File")

	l.Menu.Program.Title = "Program"

	l.Menu.App.Title = "Application"

	l.Menu.Account.Title = "Account"
	l.Menu.Account.Login = "Login"
	l.Menu.Account.Logout = "Logout"

	l.Home.Title = "Home"

	l.Settings.Title = "Settings"
	l.Settings.GeneralEdit.Title = "General"
	l.Settings.GeneralEdit.Apply = "Apply"

	l.Settings.User.Title = "User"
	l.Settings.User.UserList.Title = "User List"
	l.Settings.User.UserList.AddUser = "Add User"
	l.Settings.User.UserList.LoginName = "Login Name"
	l.Settings.User.UserList.Name = "Name"
	l.Settings.User.UserList.Options = "Options"
	l.Settings.User.UserList.Edit = "Edit"

	l.Settings.User.UserEdit.Title = "User Edit"
	l.Settings.User.UserEdit.LoginName = "Login Name"
	l.Settings.User.UserEdit.Name = "Name"
	l.Settings.User.UserEdit.State = "State"
	l.Settings.User.UserEdit.States.Public = "Public"
	l.Settings.User.UserEdit.States.Active = "Active"
	l.Settings.User.UserEdit.States.Inactive = "Inactive"
	l.Settings.User.UserEdit.Language = "Language"
	l.Settings.User.UserEdit.Password = "Password"
	l.Settings.User.UserEdit.Comment = "Comment"
	l.Settings.User.UserEdit.Groups = "Groups"
	l.Settings.User.UserEdit.Apply = "Apply"
	l.Settings.User.UserEdit.Delete = "Delete"

	l.Settings.Group.Title = "Group"
	l.Settings.Group.GroupList.Title = "Group List"
	l.Settings.Group.GroupList.AddGroup = "Add Group"
	l.Settings.Group.GroupList.Name = "Name"
	l.Settings.Group.GroupList.Comment = "Comment"
	l.Settings.Group.GroupList.Options = "Options"
	l.Settings.Group.GroupList.Edit = "Edit"

	l.Settings.Group.GroupEdit.Title = "Group Edit"
	l.Settings.Group.GroupEdit.Name = "Name"
	l.Settings.Group.GroupEdit.State = "State"
	l.Settings.Group.GroupEdit.States.Active = "Active"
	l.Settings.Group.GroupEdit.States.Inactive = "Inactive"
	l.Settings.Group.GroupEdit.Comment = "Comment"
	l.Settings.Group.GroupEdit.Pages = "Pages"
	l.Settings.Group.GroupEdit.Apply = "Apply"
	l.Settings.Group.GroupEdit.Delete = "Delete"

	l.Settings.About.Title = "About"

	l.Login.Title = "Login"
	l.Login.User = "User"
	l.Login.Password = "Password"
	l.Login.Login = "Login"

	l.Logout.Title = "Logout"

	l.Exit.Title = "Exit"

	l.Alert.Success = "Success"
	l.Alert.Info = "Info"
	l.Alert.Warning = "Warning"
	l.Alert.Error = "Error"

	l.Error.CE400BadRequest = "Bad Request"
	l.Error.CE404NotFound = "Not Found"
	l.Error.SE508LoopDetected = "Loop Detected"

	return l
}
