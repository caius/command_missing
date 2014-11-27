package main

import (
	"log"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// POSIX shells should use : for $PATH separator
const PathSeparator string = ":"

var (
	Debug *log.Logger
	Fatal *log.Logger
	Info *log.Logger

	PathComponents []string
)

func Init(debugHandle io.Writer, fatalHandle io.Writer, infoHandle io.Writer) {
	Debug = log.New(debugHandle, "DEBUG: ", 0)
	Fatal = log.New(fatalHandle, "ERROR: ", 0)
	Info = log.New(infoHandle, "", 0)

	PathComponents = strings.Split(os.Getenv("PATH"), PathSeparator)
}

func which(possibleCmd string) (actualCmd string) {
	for _, dir := range PathComponents {
		Debug.Println("Testing ", dir)

		possibleFilename := filepath.Join(dir, possibleCmd)
		Debug.Println("possibleFilename", possibleFilename)

		if _, err := os.Stat(possibleFilename); err == nil {
			Debug.Println("Found ", possibleFilename, " exists!")

			actualCmd = possibleFilename
			return
		}
	}

	return actualCmd
}

func main() {
	debugHandle := ioutil.Discard
	// TODO: check for -d from $@ instead
	if "d" == "d" {
		debugHandle = os.Stderr
	}
	Init(debugHandle, os.Stderr, os.Stdout)

	cmd := which("bash2")
	if cmd == "" {
		Fatal.Println("Not found")
		os.Exit(1)
	} else {
		Info.Println(cmd)
	}
}
