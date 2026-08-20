package parser

import (
	"fmt"
	"os"
	"strconv"
)

type UserFunc struct {
	Params []string
	Body   []string
}

var UserFunctions = make(map[string]UserFunc)

// Router
func ExecuteFunction(name string, args []string) {
	switch name {
	case "print":
		for _, arg := range args {
			val := ResolveValue(arg)
			Print(val, "string")
		}

	case "exit":
		os.Exit(0)

	case "wait":
		for _, arg := range args {
			valStr := ResolveValue(arg)
			if num, err := strconv.Atoi(valStr); err == nil {
				Wait(num)
			}
		}

	default:
		// User-defined functions lookup
		fn, exists := UserFunctions[name]
		if !exists {
			fmt.Printf("[Runtime Error]: Unknown function '%s'\n", name)
			return
		}

		if len(args) != len(fn.Params) {
			fmt.Printf("[Runtime Error]: Function '%s' expects %d args, got %d\n", name, len(fn.Params), len(args))
			return
		}

		// Param
		for i, param := range fn.Params {
			rawArg := args[i]
			// If it's already resolved.
			val := ResolveValue(rawArg)

			// Check type
			varType := "var"
			if isNumber(val) {
				varType = "int"
			} else if val == "true" || val == "false" {
				varType = "bool"
			} else {
				varType = "string"
			}

			SetVariable(param, val, varType)
		}

		for _, line := range fn.Body {
			RunLine(line)
		}
	}
}
