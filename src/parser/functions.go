package parser

import (
	"fmt"
	"strings"
)

func Getenvready() {
	fmt.Println("Atom3")
}

func Print(text string, ttype string) {
	if ttype == "string" {
		text = strings.Replace(text, `"`, "", -1)
		fmt.Println(text)
	} else if ttype == "int" {
		fmt.Println(text)
	} else if ttype == "var" {
		variable, exists := Variables[text]
		if !exists {
			fmt.Printf("[Runtime Error]: Variable '%s' is not defined\n", text)
			return
		}
		if variable.ttype == "string" {
			val := strings.Replace(variable.Value, `"`, "", -1)
			fmt.Println(val)
		} else {
			fmt.Println(variable.Value)
		}
	}
}
