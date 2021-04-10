package main

import (
	"flag"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"strings"
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
	flag.BoolVar(&fixSource, "fix", false, "cleans up common JSON syntax errors comming from MongoDB Compass; see README.md for details")
	flag.BoolVar(&backup, "backup", false, "when used with -fix, writes a copy of the original file before attempting to fix and format")
}

func main() {
	flag.Parse()

	if packageName == "" {
		packageName = os.Getenv("GOPACKAGE")
	}
	if packageName == "" {
		packageName = "main"
	}

	processFileNames()

	if verbose {
		if len(inFilenames) == 1 {
			if inFilenames[0] == "-" {
				fmt.Fprint(os.Stderr, "Reading from stdin\n")
			} else {
				fmt.Fprintf(os.Stderr, "Reading from %s\n", inFilenames[0])
			}
		} else {
			fmt.Fprintf(os.Stderr, "Reading from %d files\n", len(inFilenames))
		}
	}

	commandLine := ""
	for _, arg := range os.Args {
		commandLine = commandLine + " " + arg
	}
	outString := fmt.Sprintf(fileTemplate, strings.TrimLeft(commandLine, " "), packageName)

	for _, filename := range inFilenames {
		outString += buildFunction(filename)
	}

	output, err := format.Source([]byte(outString))

	if err != nil {
		fmt.Fprintln(os.Stderr, "WARNING: Could not gofmt output: ", err.Error())
		output = []byte(outString)
	}

	var outFile *os.File
	if outFilename == "" {
		outFilename = "<stdout>"
		outFile = os.Stdout
	} else {
		var err error
		outFile, err = os.OpenFile(outFilename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
		defer outFile.Close()
		if err != nil {
			panic(err)
		}
	}

	if verbose {
		fmt.Fprintln(os.Stderr, "Writing to", outFilename)
	}
	_, err = outFile.Write(output)
	if err != nil {
		panic(err)
	}
	if outFilename != "<stdout>" {
		outFile.Close()
		// if flag "pls"...
		cmd := exec.Command("goimports", "-w", outFilename)
		cmd.Stderr = os.Stderr

		goimpErr := cmd.Run()

		if goimpErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to process imports for %s: %s\n", outFilename, goimpErr.Error())
		}
	}
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
