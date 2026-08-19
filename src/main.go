package main

import (
	"Atom3/src/misc"
	"Atom3/src/parser"
	"fmt"
	"os"
)

func main() {
	misc.ShowBanner()

	if len(os.Args) < 2 {
		fmt.Println("Usage: atom3 <file.atom>")
		return
	}

	parser.RunFile(os.Args[1])
}
