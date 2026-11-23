package parser

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"Atom3/src/core/evaluator"
)

type UserFunc struct {
	Params   []string
	Defaults map[string]string
	Body     []string
	Module   string
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

func expressionFunctions() map[string]evaluator.Function {
	functions := make(map[string]evaluator.Function)
	for registeredName := range UserFunctions {
		name := registeredName
		if strings.Contains(name, ".") {
			name = name[strings.LastIndex(name, ".")+1:]
		}
		functionName := name
		functions[name] = func(params ...any) (any, error) {
			args := make([]string, len(params))
			for i, param := range params {
				args[i] = fmt.Sprint(param)
			}
			result := ExecuteFunction(functionName, args)
			if !result.returned {
				return nil, nil
			}
			switch result.ttype {
			case "int":
				if number, err := strconv.ParseFloat(result.value, 64); err == nil {
					return number, nil
				}
			case "bool":
				return result.value == "true", nil
			case "string":
				return strings.Trim(result.value, `"`), nil
			}
			return result.value, nil
		}
	}
	return functions
}

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

	case "arrlen":
		if Arr, exists := Arrays[args[0]]; exists {
			return Result{value: strconv.Itoa(len(Arr.Elements)), ttype: "int", returned: true}
		} else {
			RaiseException(fmt.Sprintf("Array '%s' not found or has not been defined yet!", args[0]), ErrNotFound)
			return Result{returned: false}
		}

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
			RaiseException(fmt.Sprintf("Unknown function '%s'", name), ErrRuntime)
			return Result{returned: false}
		}

		if len(args) > len(fn.Params) {
			RaiseException(fmt.Sprintf("Function '%s' expects at most %d args, got %d", name, len(fn.Params), len(args)), ErrNotFound)
			return Result{returned: false}
		}
		for i, param := range fn.Params {
			if i >= len(args) {
				if _, optional := fn.Defaults[param]; !optional {
					RaiseException(fmt.Sprintf("Function '%s' is missing required parameter '%s'", name, param), ErrNotFound)
					return Result{returned: false}
				}
			}
		}

		for i, param := range fn.Params {
			rawArg := ""
			if i < len(args) {
				rawArg = args[i]
			} else {
				rawArg = fn.Defaults[param]
			}
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
		return Result{value: frame.value, ttype: frame.ttype, returned: true}
	}
}
