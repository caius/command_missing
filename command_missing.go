package main

import (
	"fmt"
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

func candidates(one, two string) []string {
	matches := []string{}

	for i := 0; i < len(two) + 1; i++ {
		matches = append(matches, fmt.Sprintf("%s%s", one, two[0:i]))
	}

	Debug.Printf("%q", matches)

	return matches
}

func main() {
	debugHandle := ioutil.Discard
	// TODO: check for -d from $@ instead
	if os.Getenv("DEBUG") != "" {
		debugHandle = os.Stderr
	}
	Init(debugHandle, os.Stderr, os.Stdout)

	possibles := candidates(os.Args[1], os.Args[2])

	for _, possible := range possibles {
		cmd := which(possible)
		if cmd != "" {
			fmt.Println(cmd)
			os.Exit(0)
		}
	}

	Debug.Println("Not found: ", os.Args[1], os.Args[2])
	os.Exit(1)
}
