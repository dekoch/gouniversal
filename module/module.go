package module

import (
	"net/http"

	"github.com/dekoch/gouniversal/build"
	"github.com/dekoch/gouniversal/module/console"
	"github.com/dekoch/gouniversal/module/fileshare"
	"github.com/dekoch/gouniversal/module/heatingmath"
	"github.com/dekoch/gouniversal/module/homepage"
	"github.com/dekoch/gouniversal/module/iptracker"
	"github.com/dekoch/gouniversal/module/logviewer"
	"github.com/dekoch/gouniversal/module/mark"
	"github.com/dekoch/gouniversal/module/mediadownloader"
	"github.com/dekoch/gouniversal/module/mesh"
	"github.com/dekoch/gouniversal/module/meshfilesync"
	"github.com/dekoch/gouniversal/module/messenger"
	"github.com/dekoch/gouniversal/module/modbustest"
	"github.com/dekoch/gouniversal/module/openespm"
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
}

// Render UI page
func Render(page *types.Page, nav *navigation.Navigation, r *http.Request) {

	if build.ModuleConsole {
		if nav.IsNext("Console") {

			console.Render(page, nav, r)
		}
	}

	if build.ModuleLogViewer {
		if nav.IsNext("LogViewer") {

			logviewer.Render(page, nav, r)
		}
	}

	if build.ModuleOpenESPM {
		if nav.IsNext("openESPM") {

			openespm.Render(page, nav, r)
		}
	}

	if build.ModuleFileshare {
		if nav.IsNext("Fileshare") {

			fileshare.Render(page, nav, r)
		}
	}

	if build.ModuleHomepage {
		if nav.IsNext("Homepage") {

			homepage.Render(page, nav, r)
		}
	}

	if build.ModuleMeshFS {

		if nav.IsNext("MeshFS") {

			meshfilesync.Render(page, nav, r)
		}
	}

	if build.ModuleMesh || build.ModuleMessenger || build.ModuleMeshFS {
		if nav.IsNext("Mesh") {

			mesh.Render(page, nav, r)
		}
	}

	if build.ModuleMediaDownloader {
		if nav.IsNext("MediaDownloader") {

			mediadownloader.Render(page, nav, r)
		}
	}

	if build.ModuleIPTracker {
		if nav.IsNext("IPTracker") {

			iptracker.Render(page, nav, r)
		}
	}
}

// Exit is called before program exit
func Exit() {

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

	if build.ModuleLogViewer {
		logviewer.Exit()
	}

	if build.ModuleConsole {
		console.Exit()
	}
}
