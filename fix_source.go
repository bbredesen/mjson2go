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
		fmt.Fprintf(os.Stderr, "ERROR: (%s) Could not read source file for fixing\n", filename)
		os.Exit(2)
	}

	if backup {
		err = ioutil.WriteFile(filename+".backup", b, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: (%s) Could not write backup file before fixing source\n", filename)
			os.Exit(2)
		}
	}

	b = fixSourceErrors(b)

	buf := bytes.NewBuffer([]byte{})
	err = json.Indent(buf, b, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: (%s) Could not format JSON after fixing source\n", filename)
		os.Exit(2)
	}

	ioutil.WriteFile(filename, buf.Bytes(), 0666)
}

func fixSourceErrors(bytes []byte) []byte {

	bytes = fixUnquotedKeys.ReplaceAll(bytes, []byte(`$1"$2"$3`))
	bytes = fixSingleQuoteStrings.ReplaceAll(bytes, []byte(`"$1"`))
	bytes = fixTrailingCommas.ReplaceAll(bytes, []byte(`$1`))

	return bytes
}
