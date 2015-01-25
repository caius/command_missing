package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
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

type ExtendedFileMode struct {
	os.FileMode
}

// Checks filename exists, is a regular file and has the executable bit set
func (mode ExtendedFileMode) IsExecutable() bool {
	return (mode.FileMode & 0111) != 0
}

func NibbleStart(source string) (string, string) {
	result := ""
	nibble := ""
	if source != "" {
		result = source[1:]
		nibble = source[0:1]
	}
	return result, nibble
}

func NibbleEnd(source string) (string, string) {
	result, nibble := "", ""
	if source != "" {
		result = source[0 : len(source)-1]
		nibble = source[len(source)-1:]
	}
	return result, nibble
}

func Dup(source []string) []string {
	destination := []string{}
	for _, name := range source {
		destination = append(destination, name)
	}
	return destination
}

func main() {

	/* Setup logging */
	debugHandle := ioutil.Discard
	if os.Getenv("DEBUG") != "" {
		debugHandle = os.Stderr
	}
	Debug = log.New(debugHandle, "DEBUG: ", 0)
	Fatal = log.New(os.Stderr, "ERROR: ", 0)
	Info = log.New(os.Stdout, "", 0)

	/* We need a PATH to search */

	if Path = os.Getenv("PATH"); Path != "" {
		PathComponents = strings.Split(Path, string(os.PathListSeparator))
	} else {
		Fatal.Println("Couldn't find PATH directories from $PATH")
		os.Exit(2)
	}

	/* Grab ARGV for us to mutilate */

	original_arguments := []string{"", ""}
	if len(os.Args) >= 3 {
		original_arguments[0] = os.Args[1]
		original_arguments[1] = os.Args[2]
	}

	if original_arguments[0] == "" && original_arguments[1] == "" {
		// No Match, due to nothing to match
		// TODO: handle "gitst" => "git st"
		Debug.Println("no args to deal with")
		os.Exit(1)
	}

	/* Generate matches for us to check laster */

	var(
		start, end []string
		possibles [][]string
	)

	// Try moving up to three characters each way to solve issue.
	for i := 0; i < 3; i++ {
		if len(possibles) >= 2 {
			twofer := possibles[len(possibles)-2:]
			start, end = Dup(twofer[0]), Dup(twofer[1])
		} else {
			start, end = Dup(original_arguments), Dup(original_arguments)
		}

		remainder, nibble := NibbleStart(start[1])
		start[0] = start[0] + nibble
		start[1] = remainder
		if start[0] != "" {
			possibles = append(possibles, start)
		}

		remainder, nibble = NibbleEnd(end[0])
		end[0] = remainder
		end[1] = nibble + end[1]
		if end[0] != "" {
			possibles = append(possibles, end)
		}
	}

	Debug.Printf("possibles: %q", possibles)

	foundMatch := false

	// And now check each possibility
	for _, possible := range possibles {
		// And check the possible name in each path dir
		for _, dir := range PathComponents {
			// Debug.Println("Testing ", dir)

			possibleFilename := filepath.Join(dir, possible[0])
			// Debug.Println("possibleFilename", possibleFilename)

			if stat, err := os.Stat(possibleFilename); err == nil {
				mode := ExtendedFileMode{stat.Mode()}
				if mode.IsRegular() && mode.IsExecutable() {
					Debug.Println("Found", possibleFilename)

					foundMatch = true
					possible[0] = possibleFilename
					break
				}
			}
		}

		// Searched all paths for this name, output a result & exit if we have one
		if foundMatch {
			cmd := possible[0]
			possible[0] = "" // Exec#argv appears to need an empty [0] argument. WTF Go?!
			args := possible
			env := os.Environ()
			syscall.Exec(cmd, args, env)
		}
	}

	// If we get here, we didn't find a binary anywhere :'(
	Debug.Printf("Not found: %q", os.Args)
	os.Exit(1)
}
