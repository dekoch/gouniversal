package main

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/dekoch/gouniversal/build"
	"github.com/dekoch/gouniversal/module"
	"github.com/dekoch/gouniversal/program/global"
	"github.com/dekoch/gouniversal/program/lang"
	"github.com/dekoch/gouniversal/program/ui"
	"github.com/dekoch/gouniversal/shared/console"
	"github.com/dekoch/gouniversal/shared/language"
)

func main() {
	console.LoadConfig()
	console.Log("App starting...", " ")
	console.Log("Build: "+build.BuildTime, " ")
	console.Log("Commit: "+build.Commit, " ")
	console.Input("")

	en := lang.DefaultEn()
	global.Lang = language.New("data/lang/program/", en, "en")

	if build.UIEnabled {
		global.UIConfig.LoadConfig()
	}

	// UI or console mode
	if build.UIEnabled && global.UIConfig.UIEnabled {
		// start UI
		go ui.StartServer()
	} else {
		console.Log("UI is disabled", " ")

		module.LoadConfig()
	}

	go consoleInput()

	exitApp := false

	for exitApp == false {

		s := console.GetInput()

		if s != "" {

			switch s {
			case "help":
				printHelp()

			case "build":
				console.Output(build.BuildTime, " ")

			case "gover":
				console.Output(runtime.Version(), " ")

			case "gocpu":
				console.Output(strconv.Itoa(runtime.NumCPU()), " ")

			case "gonum":
				console.Output(strconv.Itoa(runtime.NumGoroutine()), " ")

			case "exit":
				exitApp = true

			default:
				console.Output("unrecognized command \""+s+"\"", " ")
				console.Output("Type 'help' for a list of available commands.", " ")
			}

			console.Input("")
		}

		time.Sleep(100 * time.Millisecond)
	}

	if build.UIEnabled && global.UIConfig.UIEnabled {
		ui.Exit()
	}

	module.Exit()

	console.Log("App ended", " ")
	os.Exit(1)
}

func printHelp() {
	console.Output("", " ")
	console.Output("Command\t\tMeaning", " ")
	console.Output("", " ")
	console.Output("help\t\tShow this help text", " ")
	console.Output("build\t\tReturns the build timestamp.", " ")
	console.Output("gover\t\tReturns the Go tree's version string.", " ")
	console.Output("gocpu\t\tReturns the number of logical CPUs usable by the current process.", " ")
	console.Output("gonum\t\tReturns the number of goroutines that currently exist.", " ")
	console.Output("exit\t\tExit this program", " ")
}

func consoleInput() {
	var input string
	scanner := bufio.NewScanner(os.Stdin)

	for {
		if scanner.Scan() {
			input = scanner.Text()

			if input != "" {

				input = strings.Replace(input, "\n", "", -1)
				input = strings.Replace(input, "\r", "", -1)

				console.Input(input)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}
