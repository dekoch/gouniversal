package main

import (
	"bufio"
	"fmt"
	"gouniversal/modules"
	"gouniversal/program/global"
	"gouniversal/program/lang"
	"gouniversal/program/ui"
	"gouniversal/shared/language"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("App starting...")

	web := new(ui.UI)

	global.Console.Input("")

	global.UiConfig.LoadConfig()

	en := lang.DefaultEn()
	global.Lang = language.New("data/lang/program/", en, "en")

	if global.UiConfig.File.UIEnabled {
		go web.StartServer()
	} else {
		fmt.Println("UI is disabled")

		modules.LoadConfig()
	}

	go consoleInput()

	exitApp := false

	for exitApp == false {

		if global.Console.Get() != "" {

			s := global.Console.Get()

			if s == "help" {
				printHelp()
			} else if s == "exit" {
				exitApp = true
			} else {
				fmt.Println("")
				fmt.Println("unrecognized command \"" + s + "\"")
				fmt.Println("Type 'help' for a list of available commands.")
			}

			global.Console.Input("")
		}

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

func consoleInput() {
	var input string
	scanner := bufio.NewScanner(os.Stdin)

	for {
		if scanner.Scan() {
			input = scanner.Text()

			if input != "" {

				input = strings.Replace(input, "\n", "", -1)
				input = strings.Replace(input, "\r", "", -1)

				global.Console.Input(input)
			}
		}
	}
}
