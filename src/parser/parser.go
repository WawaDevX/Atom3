package parser

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// Wait pauses execution for a given number of seconds
func Wait(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

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
