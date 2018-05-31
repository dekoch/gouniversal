package console

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const logFilePath = "data/log/"
const logOutputBuffer = 50

var (
	logFile   *os.File
	logLogger *log.Logger
	mut       sync.Mutex
	input     string
	output    []string
)

func init() {

	t := time.Now()

	var fileName string
	fileName += t.Format("20060102")
	fileName += "_"
	fileName += t.Format("150405")
	fileName += ".log"

	var err error
	if _, err = os.Stat(logFilePath); os.IsNotExist(err) {
		// if not found, create dir
		err = os.MkdirAll(logFilePath, os.ModePerm)
	}
	if err != nil {
		log.Fatal(err)
	}

	logFile, _ = os.Create(logFilePath + fileName)
	logLogger = log.New(logFile, "", 0)
}

func Input(s string) {

	mut.Lock()
	defer mut.Unlock()

	input = s
}

func GetInput() string {

	mut.Lock()
	defer mut.Unlock()

	return input
}

func appendOutput(s string) {

	// remove older entries
	cnt := len(output)
	if cnt > logOutputBuffer {
		output = output[1:cnt]
	}

	newOutput := make([]string, 1)
	newOutput[0] = s
	output = append(output, newOutput...)
}

func caller() string {

	_, file, line, ok := runtime.Caller(3)
	if ok == false {
		file = "???"
		line = 0
	}

	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	return file + ":" + strconv.Itoa(line)
}

func consoleOutput(message interface{}, sender string) string {

	m := fmt.Sprintf("%v", message)
	s := m

	if sender != "" {

		str := strings.Replace(sender, " ", "", -1)

		if len(str) > 0 {
			s = sender + ": " + m
		}
	} else {
		s = caller() + ": " + m
	}

	t := time.Now()
	return t.Format("2006/01/02") + " " + t.Format("15:04:05") + " " + s
}

func Output(message interface{}, sender string) {

	mut.Lock()
	defer mut.Unlock()

	s := consoleOutput(message, sender)

	appendOutput(s)

	fmt.Println(s)
}

func Log(message interface{}, sender string) {

	mut.Lock()
	defer mut.Unlock()

	s := consoleOutput(message, sender)

	appendOutput(s)

	fmt.Println(s)
	logLogger.Println(s)
}

func GetOutput() []string {

	mut.Lock()
	defer mut.Unlock()

	return output
}
