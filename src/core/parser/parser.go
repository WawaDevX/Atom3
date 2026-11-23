package parser

import (
	"time"
)

var currentFile string
var currentModule string
var importedFiles = make(map[string]bool)
var currentLine int

// Wait pauses execution for a given number of seconds
func Wait(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
