package parser

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type UserFunc struct {
	Params []string
	Body   []string
	Module string
}

type Result struct {
	value    string
	ttype    string
	returned bool
}

var UserFunctions = make(map[string]UserFunc)

type functionReturn struct {
	value    string
	ttype    string
	returned bool
}

var functionReturnStack []*functionReturn

func ExecuteFunction(name string, args []string) Result {
	switch name {
	case "print":
		for _, arg := range args {
			Print(arg, "string")
		}
		return Result{returned: false}

	case "exit":
		os.Exit(0)
		return Result{returned: false}

	case "input":
		if len(args) > 0 {
			fmt.Print(ResolveValue(args[0], true))
		}
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			return Result{returned: false}
		}
		return Result{value: input, ttype: DetermineType(input), returned: true}

	case "wait":
		for _, arg := range args {
			valStr := ResolveValue(arg, true)
			if num, err := strconv.Atoi(valStr); err == nil {
				Wait(num)
			}
		}
		return Result{returned: false}

	default:
		functionName := name
		if !strings.Contains(name, ".") && currentModule != "" {
			moduleFunction := currentModule + "." + name
			if _, exists := UserFunctions[moduleFunction]; exists {
				functionName = moduleFunction
			}
		}
		fn, exists := UserFunctions[functionName]
		if !exists {
			fmt.Printf("[Runtime Error]: Unknown function '%s'\n", name)
			return Result{returned: false}
		}

		if len(args) != len(fn.Params) {
			fmt.Printf("[Runtime Error]: Function '%s' expects %d args, got %d\n", name, len(fn.Params), len(args))
			return Result{returned: false}
		}

		for i, param := range fn.Params {
			rawArg := args[i]
			val := ResolveValue(rawArg, true)
			SetVariable(param, val, DetermineType(val))
		}

		frame := &functionReturn{}
		functionReturnStack = append(functionReturnStack, frame)
		callerBlocks := blockStack
		callerLastIfClosed := lastIfClosed
		callerModule := currentModule
		currentModule = fn.Module
		defer func() {
			functionReturnStack = functionReturnStack[:len(functionReturnStack)-1]
			blockStack = callerBlocks
			lastIfClosed = callerLastIfClosed
			currentModule = callerModule
		}()

		for _, line := range fn.Body {
			RunLine(line)
			if frame.returned {
				break
			}
		}
		if !frame.returned {
			return Result{returned: false}
		}
		return Result{value: frame.value, ttype: DetermineType(frame.value), returned: true}
	}
}
