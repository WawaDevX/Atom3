package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var currentFile string
var currentModule string
var importedFiles = make(map[string]bool)

// Wait pauses execution for a given number of seconds
func Wait(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

func RunFile(filename string) {
	path, err := filepath.Abs(filename)
	if err != nil {
		fmt.Println("Error resolving file:", err)
		return
	}
	runFile(path, "")
}

func ImportFile(name string) {
	if filepath.Ext(name) != ".atom" {
		name += ".atom"
	}
	base := "."
	if currentFile != "" {
		base = filepath.Dir(currentFile)
	}
	path, err := filepath.Abs(filepath.Join(base, name))
	if err != nil {
		fmt.Println("Error resolving import:", err)
		return
	}
	if importedFiles[path] {
		return
	}
	module := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
	runFile(path, module)
}

func runFile(filename, module string) {
	if importedFiles[filename] {
		return
	}
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	importedFiles[filename] = true
	previousFile := currentFile
	previousModule := currentModule
	currentFile = filename
	currentModule = module
	defer func() {
		currentFile = previousFile
		currentModule = previousModule
	}()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		RunLine(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
