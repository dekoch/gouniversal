package main

import (
	"bufio"
	"fmt"
	"gouniversal/modules"
	"gouniversal/program/global"
	"gouniversal/program/lang"
	"gouniversal/program/ui"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("App starting...")

	web := new(ui.UI)

	global.Console.Input = ""

	global.UiConfig.LoadConfig()

	global.Lang.File = lang.LoadLangFiles()

	if global.UiConfig.File.UIEnabled {
		go web.StartServer()
	} else {
		fmt.Println("UI is disabled")

		modules.LoadConfig()
	}

	exitApp := false
	reader := bufio.NewReader(os.Stdin)

	for exitApp == false {
		input, _ := reader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)

		global.Console.Mut.Lock()
		if input != "" {
			global.Console.Input = input
		}

		if global.Console.Input != "" {

			if global.Console.Input == "help" {
				printHelp()
			} else if global.Console.Input == "exit" {
				exitApp = true
			} else {
				fmt.Println("")
				fmt.Println("unrecognized command \"" + global.Console.Input + "\"")
				fmt.Println("Type 'help' for a list of available commands.")
			}

			global.Console.Input = ""
		}
		global.Console.Mut.Unlock()

		time.Sleep(100 * time.Millisecond)
	}

	if global.UiConfig.File.UIEnabled {
		web.Exit()
	}

	modules.Exit()

	fmt.Println("App ended")
	os.Exit(1)
}

func printHelp() {
	fmt.Println("")
	fmt.Println("Command\t\tMeaning")
	fmt.Println("")
	fmt.Println("help\t\tShow this help text")
	fmt.Println("exit\t\tExit this program")
}
