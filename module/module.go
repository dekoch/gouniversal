package module

import (
	"net/http"

	"github.com/dekoch/gouniversal/build"
	"github.com/dekoch/gouniversal/module/backup"
	"github.com/dekoch/gouniversal/module/console"
	"github.com/dekoch/gouniversal/module/fileshare"
	"github.com/dekoch/gouniversal/module/gasprice"
	"github.com/dekoch/gouniversal/module/gpsnav"
	"github.com/dekoch/gouniversal/module/heatingmath"
	"github.com/dekoch/gouniversal/module/homepage"
	"github.com/dekoch/gouniversal/module/instabackup"
	"github.com/dekoch/gouniversal/module/iptracker"
	"github.com/dekoch/gouniversal/module/logviewer"
	"github.com/dekoch/gouniversal/module/mark"
	"github.com/dekoch/gouniversal/module/mediadownloader"
	"github.com/dekoch/gouniversal/module/mesh"
	"github.com/dekoch/gouniversal/module/meshfilesync"
	"github.com/dekoch/gouniversal/module/messenger"
	"github.com/dekoch/gouniversal/module/modbustest"
	"github.com/dekoch/gouniversal/module/monmotion"
	"github.com/dekoch/gouniversal/module/openespm"
	"github.com/dekoch/gouniversal/module/paratest"
	"github.com/dekoch/gouniversal/module/picturex"
	"github.com/dekoch/gouniversal/shared/navigation"
	"github.com/dekoch/gouniversal/shared/types"

	sharedConsole "github.com/dekoch/gouniversal/shared/console"
)

// Modules provide a interface to nest apps and modules

// LoadConfig is called before UI starts
func LoadConfig() {

	if build.ModuleConsole {
		sharedConsole.Log("Console enabled", "Module")
		console.LoadConfig()
	}

	if build.ModuleLogViewer {
		sharedConsole.Log("LogViewer enabled", "Module")
		logviewer.LoadConfig()
	}

	if build.ModuleBackup {
		sharedConsole.Log("Backup enabled", "Module")
		backup.LoadConfig()
	}

	if build.ModuleOpenESPM {
		sharedConsole.Log("openESPM enabled", "Module")
		openespm.LoadConfig()
	}

	if build.ModuleFileshare {
		sharedConsole.Log("Fileshare enabled", "Module")
		fileshare.LoadConfig()
	}

	if build.ModuleHomepage {
		sharedConsole.Log("Homepage enabled", "Module")
		homepage.LoadConfig()
	}

	if build.ModuleMesh || build.ModuleMessenger || build.ModuleMeshFS {
		sharedConsole.Log("Mesh enabled", "Module")
		mesh.LoadConfig()
	}

	if build.ModuleMessenger {
		sharedConsole.Log("Messenger enabled", "Module")
		messenger.LoadConfig()
	}

	if build.ModuleMeshFS {
		sharedConsole.Log("MeshFileSync enabled", "Module")
		meshfilesync.LoadConfig()
	}

	if build.ModuleMediaDownloader {
		sharedConsole.Log("MediaDownloader enabled", "Module")
		mediadownloader.LoadConfig()
	}

	if build.ModuleIPTracker {
		sharedConsole.Log("IPTracker enabled", "Module")
		iptracker.LoadConfig()
	}

	if build.ModuleGasPrice {
		sharedConsole.Log("GasPrice enabled", "Module")
		gasprice.LoadConfig()
	}

	if build.ModulePictureX {
		sharedConsole.Log("PictureX enabled", "Module")
		picturex.LoadConfig()
	}

	if build.ModuleInstaBackup {
		sharedConsole.Log("InstaBackup enabled", "Module")
		instabackup.LoadConfig()
	}

	if build.ModuleMonMotion {
		sharedConsole.Log("MonMotion enabled", "Module")
		monmotion.LoadConfig()
	}

	if build.ModuleGPSNav {
		sharedConsole.Log("GPSNav enabled", "Module")
		gpsnav.LoadConfig()
	}

	if build.ModuleModbusTest {
		sharedConsole.Log("ModbusTest enabled", "Module")
		modbustest.LoadConfig()
	}

	if build.ModuleHeatingMath {
		sharedConsole.Log("HeatingMath enabled", "Module")
		heatingmath.LoadConfig()
	}

	if build.ModuleMark {
		sharedConsole.Log("Mark enabled", "Module")
		mark.LoadConfig()
	}

	if build.ModuleParaTest {
		sharedConsole.Log("ParaTest enabled", "Module")
		paratest.LoadConfig()
	}
}

// RegisterPage adds pages to a sitemap
// use
// nav.Sitemap.Register("App:MyApp", "MyApp")
// nav.Sitemap.Register("App:Program:MyApp", "MyApp")
// nav.Sitemap.Register("App:Account:MyApp", "MyApp")
//
// to nest your app or module inside menus
func RegisterPage(page *types.Page, nav *navigation.Navigation) {

	if build.ModuleConsole {
		console.RegisterPage(page, nav)
	}

	if build.ModuleLogViewer {
		logviewer.RegisterPage(page, nav)
	}

	if build.ModuleBackup {
		backup.RegisterPage(page, nav)
	}

	if build.ModuleOpenESPM {
		openespm.RegisterPage(page, nav)
	}

	if build.ModuleFileshare {
		fileshare.RegisterPage(page, nav)
	}

	if build.ModuleHomepage {
		homepage.RegisterPage(page, nav)
	}

	if build.ModuleMeshFS {
		meshfilesync.RegisterPage(page, nav)
	}

	if build.ModuleMesh || build.ModuleMessenger || build.ModuleMeshFS {
		mesh.RegisterPage(page, nav)
	}

	if build.ModuleMediaDownloader {
		mediadownloader.RegisterPage(page, nav)
	}

	if build.ModuleIPTracker {
		iptracker.RegisterPage(page, nav)
	}

	if build.ModuleGasPrice {
		gasprice.RegisterPage(page, nav)
	}

	if build.ModulePictureX {
		picturex.RegisterPage(page, nav)
	}

	if build.ModuleInstaBackup {
		instabackup.RegisterPage(page, nav)
	}

	if build.ModuleGPSNav {
		gpsnav.RegisterPage(page, nav)
	}
}

// Render UI page
func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	switch nav.GetNextPage() {
	case "Console":
		if build.ModuleConsole {
			console.Render(page, nav, r)
			return
		}

	case "LogViewer":
		if build.ModuleLogViewer {
			logviewer.Render(page, nav, r)
			return
		}

	case "Backup":
		if build.ModuleBackup {
			backup.Render(page, nav, r)
			return
		}

	case "openESPM":
		if build.ModuleOpenESPM {
			openespm.Render(page, nav, r)
			return
		}

	case "Fileshare":
		if build.ModuleFileshare {
			fileshare.Render(page, nav, r)
			return
		}

	case "Homepage":
		if build.ModuleHomepage {
			homepage.Render(page, nav, r)
			return
		}

	case "MeshFS":
		if build.ModuleMeshFS {
			meshfilesync.Render(page, nav, r)
			return
		}

	case "Mesh":
		if build.ModuleMesh ||
			build.ModuleMessenger ||
			build.ModuleMeshFS {

			mesh.Render(page, nav, r)
			return
		}

	case "MediaDownloader":
		if build.ModuleMediaDownloader {
			mediadownloader.Render(page, nav, r)
			return
		}

	case "IPTracker":
		if build.ModuleIPTracker {
			iptracker.Render(page, nav, r)
			return
		}

	case "GasPrice":
		if build.ModuleGasPrice {
			gasprice.Render(page, nav, r)
			return
		}

	case "PictureX":
		if build.ModulePictureX {
			picturex.Render(page, nav, r)
			return
		}

	case "InstaBackup":
		if build.ModuleInstaBackup {
			instabackup.Render(page, nav, r)
			return
		}

	case "GPSNav":
		if build.ModuleGPSNav {
			gpsnav.Render(page, nav, r)
			return
		}
	}

	nav.RedirectPath("404", true)
}

// Exit is called before program exit
func Exit(em *types.ExitMessage) {

	if build.ModuleOpenESPM {
		openespm.Exit()
	}

	if build.ModuleFileshare {
		fileshare.Exit()
	}

	if build.ModuleHomepage {
		homepage.Exit()
	}

	if build.ModuleMessenger {
		messenger.Exit()
	}

	if build.ModuleMeshFS {
		meshfilesync.Exit()
	}

	if build.ModuleMesh || build.ModuleMessenger || build.ModuleMeshFS {
		mesh.Exit()
	}

	if build.ModuleMediaDownloader {
		mediadownloader.Exit()
	}

	if build.ModuleIPTracker {
		iptracker.Exit()
	}

	if build.ModuleGasPrice {
		gasprice.Exit()
	}

	if build.ModulePictureX {
		picturex.Exit()
	}

	if build.ModuleInstaBackup {
		instabackup.Exit(em)
	}

	if build.ModuleMonMotion {
		monmotion.Exit(em)
	}

	if build.ModuleGPSNav {
		gpsnav.Exit()
	}

	if build.ModuleParaTest {
		paratest.Exit()
	}

	if build.ModuleBackup {
		backup.Exit(em)
	}

	if build.ModuleLogViewer {
		logviewer.Exit()
	}

	if build.ModuleConsole {
		console.Exit()
	}
}
