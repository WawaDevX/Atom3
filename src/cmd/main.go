package main

import (
	"Atom3/src/core/parser"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: atom3 <file.atom>")
		return
	}

	if os.Args[1] == "version" {
		ShowVersion()
	} else {
		parser.RunFile(os.Args[1])
	}
}
