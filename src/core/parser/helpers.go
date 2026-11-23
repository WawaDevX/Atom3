package parser

import (
	"fmt"
	"strconv"
	"strings"

	atomlexer "Atom3/src/core/lexer"
)

// Router function to determine a values true value.
func ResolveValue(arg string, doexpr bool) string {

	if strings.HasPrefix(arg, "[") && strings.HasSuffix(arg, "]") && doexpr {
		SetVariable("__temp__", arg, "expr")
		return Variables["__temp__"].Value //TODO: make a better way of doing this.
	} else if strings.Contains(arg, "[") {

		parts := strings.SplitN(arg, "[", 2)

		arrayname := parts[0]
		existing, ok := Variables[arrayname]

		if !ok {
			RaiseException(
				fmt.Sprintf("%s is not defined", arrayname),
				ErrRuntime,
			)
			return ""
		} else if existing.ttype != "array" {
			RaiseException(
				fmt.Sprintf("%s is not an array", arrayname),
				ErrRuntime,
			)
			return ""
		}

		index := ResolveValue(strings.TrimSuffix(parts[1], "]"), false)
		indexresult, err := strconv.Atoi(index)
		if err != nil {
			RaiseException(
				fmt.Sprintf("Index %s is invalid or empty", parts[1]),
				ErrRuntime,
			)
			return ""
		} else if indexresult < 0 || indexresult >= len(Arrays[arrayname].Elements) {
			RaiseException(
				fmt.Sprintf("Index %d is out of bounds for array %s", indexresult, arrayname),
				ErrIndexOutOfRange,
			)
			return ""
		} else {
			return ResolveValue(Arrays[arrayname].Elements[indexresult].Value, false)
		}

	}
	// Variables.
	if v, exists := Variables[arg]; exists {
		return v.Value
	}

	trimmed := strings.TrimSpace(arg)
	if strings.HasPrefix(trimmed, `"`) && strings.HasSuffix(trimmed, `"`) {
		return strings.Trim(trimmed, `"`)
	}
	if trimmed == "true" || trimmed == "false" || isNumeric(trimmed) {
		return trimmed
	}
	if len(trimmed) > 0 && atomlexer.IsIdentifierStart(trimmed[0]) {
		valid := true
		for i := 1; i < len(trimmed); i++ {
			if !atomlexer.IsIdentifierPart(trimmed[i]) {
				valid = false
				break
			}
		}
		if valid {
			RaiseException(fmt.Sprintf("Variable '%s' not found or has not been defined yet!", trimmed), ErrNotFound)
			return ""
		}
	}

	return trimmed
}

func Getenvready() {
	fmt.Println("Atom3")
}

func Print(text string, ttype string) {
	val := ResolveValue(text, true)
	fmt.Println(val)
}

func DetermineType(val string) string {
	t := "var"
	if strings.HasPrefix(val, "{") && strings.HasSuffix(val, "}") {
		t = "array"
	} else if isNumeric(val) {
		t = "int"
	} else if val == "true" || val == "false" {
		t = "bool"
	} else {
		t = "string"
	}
	return t
}
