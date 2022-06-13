package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall/js"
)

var (
	packageName string
	funcName    string

	inFilenames []string
	inFiles     []*os.File

	outFilename string

	verbose bool

	backup, fixSource bool
)

func init() {
	flag.StringVar(&packageName, "package", "", "the Go package name to produce code for; defaults to the GOPACKAGE environment variable, or 'main' if not set")
	flag.BoolVar(&verbose, "v", false, "more verbose output (to stderr)")
	flag.StringVar(&outFilename, "out", "", "filename to write to instead of stdout")
	flag.BoolVar(&fixSource, "fix", true, "cleans up common JSON syntax errors comming from MongoDB Compass; see README.md for details")
	flag.BoolVar(&backup, "backup", false, "when used with -fix, writes a copy of the original file before attempting to fix and format")
}

func main() {
	js.Global().Set("fixSourceErrors", wrap_fixSourceErrors())
	js.Global().Set("buildFunction", wrap_buildFunction())

	<-make(chan bool) // wait on a channel to block exiting
}

func processFileNames() {
	if len(flag.Args()) == 0 {
		// if verbose {
		// 	// fmt.Fprintln(os.Stderr, "No input files provided, reading from stdin")
		// }
		inFilenames = []string{"-"}
		inFiles = []*os.File{os.Stdin}
		return
	} else {
		for _, path := range flag.Args() {
			fi, err := os.Stat(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
				os.Exit(1)
			}
			if fi.IsDir() {
				dir, err := os.Open(fi.Name())
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not open directory %s: %s\n", path, err.Error())
					os.Exit(1)
				}
				entries, err := dir.Readdir(0)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not read directory entries %s: %s\n", path, err.Error())
					os.Exit(1)
				}

				for _, entry := range entries {
					if strings.HasSuffix(entry.Name(), ".json") && !entry.IsDir() {
						fullName := fmt.Sprint(path, string(os.PathSeparator), entry.Name())
						inFilenames = append(inFilenames, fullName)
					}
				}
				dir.Close()
			} else {
				inFilenames = append(inFilenames, path)
			}
		}
	}
}
