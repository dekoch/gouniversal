package lang

import (
	"github.com/dekoch/gouniversal/shared/config"
)

type HomeShareLink struct {
	Title      string
	Step1      string
	Step2      string
	Step3      string
	Step4      string
	Step5      string
	NewPair    string
	CopyLink   string
	Upload     string
	Unlock     string
	DeletePair string
}

type HomeLinkReceived struct {
	Title      string
	Step1      string
	Step2      string
	Step3      string
	Step4      string
	Upload     string
	Unlock     string
	DeletePair string
}

type Home struct {
	Menu          string
	Title         string
	ShareLink     HomeShareLink
	LinkReceived  HomeLinkReceived
	FirstPicture  string
	SecondPicture string
}
type Alert struct {
	Success string
	Info    string
	Warning string
	Error   string
}

type LangFile struct {
	Header config.FileHeader
	Home   Home
	Alert  Alert
}

func DefaultEn() LangFile {

	var l LangFile

	l.Header = config.BuildHeader("en", "Lang", 1.0, "Language File")

	l.Home.Menu = "Tools"
	l.Home.Title = "PictureX"

	l.Home.ShareLink.Title = "Share Link"
	l.Home.ShareLink.Step1 = "1. Share Link"
	l.Home.ShareLink.Step2 = "2. Upload First Picture"
	l.Home.ShareLink.Step3 = "3. Wait for Second Picture"
	l.Home.ShareLink.Step4 = "4. Unlock Second Picture"
	l.Home.ShareLink.Step5 = "5. Delete Pair"
	l.Home.ShareLink.NewPair = "New Pair"
	l.Home.ShareLink.CopyLink = "Copy Link"
	l.Home.ShareLink.Upload = "Upload"
	l.Home.ShareLink.Unlock = "Unlock"
	l.Home.ShareLink.DeletePair = "Delete"

	l.Home.LinkReceived.Title = "You received a Link"
	l.Home.LinkReceived.Step1 = "1. Upload Second Picture"
	l.Home.LinkReceived.Step2 = "2. Wait for First Picture"
	l.Home.LinkReceived.Step3 = "3. Unlock First Picture"
	l.Home.LinkReceived.Step4 = "4. Delete Pair"
	l.Home.LinkReceived.Upload = "Upload"
	l.Home.LinkReceived.Unlock = "Unlock"
	l.Home.LinkReceived.DeletePair = "Delete"

	l.Home.FirstPicture = "First Picture"
	l.Home.SecondPicture = "Second Picture"

	l.Alert.Success = "Success"
	l.Alert.Info = "Info"
	l.Alert.Warning = "Warning"
	l.Alert.Error = "Error"

	return l
}
