package misc

import (
	"os"
	"os/exec"
	"runtime"
)

func ClearTerminal() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func ShowVersion() {
	println("Atom3")
	println("Release: 2026.0")
	println("by the Atom3 contributors.")
}
