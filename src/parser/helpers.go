package parser

import (
	"fmt"
	"strconv"
	"strings"
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
			fmt.Printf("Error: %s is not defined\n", arrayname)
			return ""
		} else if existing.ttype != "array" {
			fmt.Printf("Error: %s is not an array\n", arrayname)
			return ""
		}

		index := ResolveValue(strings.TrimSuffix(parts[1], "]"), false)
		indexresult, err := strconv.Atoi(index)
		if err != nil {
			fmt.Printf("Error: index %s is invalid or empty\n", parts[1])
			return ""
		} else if indexresult < 0 || indexresult >= len(Arrays[arrayname].Elements) {
			fmt.Printf("Error: index %d is out of bounds for array %s\n", indexresult, arrayname)
			return ""
		} else {
			return ResolveValue(Arrays[arrayname].Elements[indexresult].Value, false)
		}

	}
	// Variables.
	if v, exists := Variables[arg]; exists {
		return v.Value
	}
	// Booleans and strings
	return strings.Trim(arg, `"`)
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
