package parser

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Wait pauses execution for a given number of seconds
func wait(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

// ExecuteFunction acts as the central router for built-in functions
func ExecuteFunction(name string, args []string) {
	switch name {
	case "print":
		for _, arg := range args {
			val := ResolveValue(arg)
			Print(val, "string")
		}

	case "exit":
		os.Exit(0)

	case "wait":
		for _, arg := range args {
			valStr := ResolveValue(arg)
			if num, err := strconv.Atoi(valStr); err == nil {
				Wait(num)
			}
		}

	default:
		fmt.Printf("[Runtime Error]: Unknown function '%s'\n", name)
	}
}
