package main

import (
	"bufio"
	"fmt"
	"gouniversal/program/global"
	"gouniversal/program/groupManagement"
	"gouniversal/program/programConfig"
	"gouniversal/program/ui"
	"gouniversal/program/userManagement"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("App starting...")

	global.Console.Input = ""

	//time.Sleep(15 * time.Second)
	//fmt.Println("exit")
	//os.Exit(1)

	global.ProgramConfig = programConfig.LoadConfig()
	global.UiConfig = ui.LoadUiConfig()
	global.UserConfig.File = userManagement.LoadUser()
	global.GroupConfig.File = groupManagement.LoadGroup()

	web := new(ui.UI)
	go web.StartServer()

	var boExit bool
	reader := bufio.NewReader(os.Stdin)

	for boExit == false {
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
				boExit = true
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

	web.Exit()

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
