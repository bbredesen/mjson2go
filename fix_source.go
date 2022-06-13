package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

var (
	fixUnquotedKeys       = regexp.MustCompile(`([\s\{])(\$?\w+)(\s?:)`)
	fixSingleQuoteStrings = regexp.MustCompile(`'(\$?\w+)'`)
	fixTrailingCommas     = regexp.MustCompile(`,(\n?\s*[\]|\}])`)
)

func fixSourceFile(filename string) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: (%s) Could not read source file for fixing: %s\n", filename, err.Error())
		os.Exit(2)
	}

	if backup {
		err = ioutil.WriteFile(filename+".backup", b, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: (%s) Could not write backup file before fixing source: %s\n", filename, err.Error())
			os.Exit(2)
		}
	}

	b = fixSourceErrors(b)

	// buf := bytes.NewBuffer([]byte{})
	// err = json.Indent(buf, b, "", "  ")
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "ERROR: (%s) Could not format JSON after fixing source: %s\n", filename, err.Error())
	// 	os.Exit(2)
	// }

	ioutil.WriteFile(filename, b, 0666)
}

func fixSourceErrors(b []byte) []byte {

	b = fixUnquotedKeys.ReplaceAll(b, []byte(`$1"$2"$3`))
	b = fixSingleQuoteStrings.ReplaceAll(b, []byte(`"$1"`))
	b = fixTrailingCommas.ReplaceAll(b, []byte(`$1`))

	buf := bytes.NewBuffer([]byte{})

	err := json.Indent(buf, b, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not format JSON after fixing source: %s\n", err.Error())
		os.Exit(2)
	}

	return buf.Bytes()
}
