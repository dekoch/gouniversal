package modules

import (
	"net/http"

	"github.com/dekoch/gouniversal/build"
	"github.com/dekoch/gouniversal/modules/console"
	"github.com/dekoch/gouniversal/modules/fileshare"
	"github.com/dekoch/gouniversal/modules/homepage"
	"github.com/dekoch/gouniversal/modules/logviewer"
	"github.com/dekoch/gouniversal/modules/mesh"
	"github.com/dekoch/gouniversal/modules/meshFileSync"
	"github.com/dekoch/gouniversal/modules/messenger"
	"github.com/dekoch/gouniversal/modules/modbustest"
	"github.com/dekoch/gouniversal/modules/openespm"
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
		meshFileSync.LoadConfig()
	}

	if build.ModuleModbusTest {
		sharedConsole.Log("ModbusTest enabled", "Module")
		modbustest.LoadConfig()
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
}

// Exit is called before program exit
func Exit() {

	if build.ModuleLogViewer {
		logviewer.Exit()
	}

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
		meshFileSync.Exit()
	}

	if build.ModuleMesh || build.ModuleMessenger || build.ModuleMeshFS {
		mesh.Exit()
	}

	if build.ModuleConsole {
		console.Exit()
	}
}
