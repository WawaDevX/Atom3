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
func ShowBanner() {
	println("Atom3 Compiler.")
	println("Release: 2026.0")
	println("Made by WawaDev (formerly spacecat) and the Atom3 Team.")
	println("Three years of nothing, Third attempt at a new programming language and three mugs of tea.")

	ClearTerminal()

}
