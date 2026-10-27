package main

import (
	"Atom3/src/misc"
	"Atom3/src/parser"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: atom3 <file.atom>")
		return
	}

	if os.Args[1] == "version" {
		misc.ShowVersion()
	} else {
		parser.RunFile(os.Args[1])
	}
}
