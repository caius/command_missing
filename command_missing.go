package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	Debug *log.Logger
	Fatal *log.Logger
	Info  *log.Logger

	Path           string
	PathComponents []string
)

// http://stackoverflow.com/a/9561388
func appendIfMissing(slice []string, str string) []string {
	for _, ele := range slice {
		if ele == str {
			return slice
		}
	}
	return append(slice, str)
}

// Checks filename exists, is a regular file and has the executable bit set
func isExecutableFile(filename string) bool {
	if fileinfo, err := os.Stat(filename); err == nil {
		if mode := fileinfo.Mode(); mode.IsRegular() && (mode & 0111) != 0 {
			return true
		}
	}

	return false
}

func main() {
	debugHandle := ioutil.Discard
	if os.Getenv("DEBUG") != "" {
		debugHandle = os.Stderr
	}
	Debug = log.New(debugHandle, "DEBUG: ", 0)
	Fatal = log.New(os.Stderr, "ERROR: ", 0)
	Info = log.New(os.Stdout, "", 0)

	if Path = os.Getenv("PATH"); Path != "" {
		PathComponents = strings.Split(Path, string(os.PathListSeparator))
	} else {
		Fatal.Println("Couldn't find PATH directories from $PATH")
		os.Exit(2)
	}

	one, two := "", ""
	// 0 is program name, 1+ is ARGV
	switch {
	case len(os.Args) == 2:
		one = os.Args[1]
	case len(os.Args) >= 3:
		one = os.Args[1]
		two = os.Args[2]
	}

	if one == "" && two == "" {
		// No Match, due to nothing to match
		Debug.Println("no args to deal with")
		os.Exit(1)
	}

	one_possibles := []string{}
	two_possibles := []string{}

	for i := len(one) - 1; i > 0; i-- {
		one_possibles = append(one_possibles, one[0:i])
	}

	for i := 1; i < len(two)+1; i++ {
		two_possibles = append(two_possibles, one+two[0:i])
	}

	Debug.Printf("one_possibles: %q", one_possibles)
	Debug.Printf("two_possibles: %q", two_possibles)

	possibles := []string{}
	for i, _ := range one_possibles {
		possibles = appendIfMissing(possibles, one_possibles[i])
		if i < len(two_possibles) && two_possibles[i] != "" {
			possibles = appendIfMissing(possibles, two_possibles[i])
		}
	}

	Debug.Printf("%q", possibles)

	cmd := ""

	// And now check each possibility
	for _, possible := range possibles {
		// And check the possible name in each path dir
		for _, dir := range PathComponents {
			// Debug.Println("Testing ", dir)

			possibleFilename := filepath.Join(dir, possible)
			// Debug.Println("possibleFilename", possibleFilename)

			if isExecutableFile(possibleFilename) {
				Debug.Println("Found ", possibleFilename, " exists!")

				cmd = possibleFilename
			}
		}

		// Searched all paths for this name, output a result & exit if we have one
		if cmd != "" {
			fmt.Println(cmd)
			os.Exit(0)
		}
	}

	// If we get here, we didn't find a binary anywhere :'(
	Debug.Printf("Not found: %q", os.Args)
	os.Exit(1)
}
