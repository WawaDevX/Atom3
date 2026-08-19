package parser

import (
	"bufio"
	"fmt"
	"os"
)

func RunFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		RunLine(scanner.Text())
	}
}
