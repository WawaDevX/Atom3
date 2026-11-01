package parser

import (
	"fmt"
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
)

func RaiseException(message string, errType string) error {
	fmt.Printf("[%s]: %s (Fl: %s, Ln: %d)\n", errType, message, currentFile, currentLine)
	return fmt.Errorf("[%s]: %s (Fl: %s, Ln: %d)", errType, message, currentFile, currentLine)
}
