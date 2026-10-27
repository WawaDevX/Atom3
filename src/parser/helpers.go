package parser

import (
	"fmt"
	"strings"
)

// Router function to determine a values true value.
func ResolveValue(arg string, doexpr bool) string {

	if strings.HasPrefix(arg, "[") && strings.HasSuffix(arg, "]") && doexpr {
		SetVariable("__temp__", arg, "expr")
		return Variables["__temp__"].Value //TODO: make a better way of doing this.
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
	if isNumeric(val) {
		t = "int"
	} else if val == "true" || val == "false" {
		t = "bool"
	} else {
		t = "string"
	}
	return t
}
