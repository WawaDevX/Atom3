package parser

import (
	"fmt"
	"os"
)

const (
	ErrImport              = "Import Error"
	ErrInternal            = "Internal Error"
	ErrRuntime             = "Runtime Error"
	ErrParse               = "Parse Error"
	ErrType                = "Type Error"
	ErrSyntax              = "Syntax Error"
	ErrFileNotFound        = "File Not Found Error"
	ErrIndexOutOfRange     = "Index Out of Range Error"
	ErrDivisionByZero      = "Division By Zero Error"
	ErrInvalidOperation    = "Invalid Operation Error"
	ErrInvalidFunctionCall = "Invalid Function Call Error"
	ErrNotFound            = "Not Found Error"
)

func RaiseException(message string, errType string) error {
	errMsg := fmt.Sprintf("[%s]: %s (Fl: %s, Ln: %d)", errType, message, currentFile, currentLine)
	fmt.Println(errMsg)
	os.Exit(1)
	return fmt.Errorf("%s", errMsg) // satisfies compiler
}
