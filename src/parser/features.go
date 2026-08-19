package parser

import (
	"fmt"
	"strings"
)

func ResolveValue(arg string) string {

	if strings.HasPrefix(arg, "[") && strings.HasSuffix(arg, "]") {
		SetVariable("__temp__", arg, "expr")
		return Variables["__temp__"].Value //TODO: make a better way of doing this.
	}
	// 2. Variables: x
	if v, exists := Variables[arg]; exists {
		return v.Value
	}
	// 3. Raw Strings ("hello"), Raw Numbers (10), or Raw Booleans (true/false)
	return strings.Trim(arg, `"`)
}

func Getenvready() {
	fmt.Println("Atom3")
}

func Print(text string, ttype string) {
	val := ResolveValue(text)
	fmt.Println(val)
}
