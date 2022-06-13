//go:build js

package main

import (
	"fmt"
	"go/format"
	"syscall/js"
)

func wrap_fixSourceErrors() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Invalid number of arguments passed"
		}

		input := args[0].String()

		output := fixSourceErrors([]byte(input))
		output, err := format.Source(output)

		if err != nil {
			fmt.Println("WARNING: Could not format output: ", err.Error())
		}

		return string(output)
	})
}

func wrap_buildFunction() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 2 {
			return "Invalid number of arguments passed"
		}

		input, fnName := args[0].Bytes(), args[1].String()

		output, err := buildFunction(input, fnName)
		if err != nil {
			fmt.Println("WARNING: Could not build function: ", err.Error())
		}

		output, err = format.Source(output)
		if err != nil {
			fmt.Println("WARNING: Could not format output: ", err.Error())
		}

		return string(output)
	})
}
