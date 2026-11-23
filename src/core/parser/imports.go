package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func RunFile(filename string) {
	path, err := filepath.Abs(filename)
	if err != nil {
		RaiseException("Error resolving file path: "+err.Error(), ErrInternal)
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
		RaiseException("Error resolving import path: "+err.Error(), ErrInternal)
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
	file, err := os.Open(filename)
	if err != nil {
		RaiseException("Error opening file: "+err.Error(), ErrInternal)
		return
	}
	defer file.Close()

	importedFiles[filename] = true
	previousFile := currentFile
	previousModule := currentModule
	currentFile = filename
	currentModule = module
	currentLine = 0
	defer func() {
		currentFile = previousFile
		currentModule = previousModule
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine++
		RunLine(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		RaiseException("Error reading file: "+err.Error(), ErrInternal)
	}
}
