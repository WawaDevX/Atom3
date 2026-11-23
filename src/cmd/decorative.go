package main

import (
	"math/rand"
	"os"
	"os/exec"
	"runtime"
)

var atoms = []string{"▀▄▀▄▀▄ ATOM³ ▄▀▄▀▄▀", "▐░▒▓ 𝗔𝗧𝗢𝗠 𝟯 ▓▒░▌", "ATOM³", "( ˘▽˘)っ 𝐴𝑇𝑂𝑀3", "「 ATOM3 」", "𝗔𝗧𝗢𝗠 𝟯 ¯\\_(ツ)_/¯"}
var random = atoms[rand.Intn(len(atoms))]

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
	println(random)
	println("Release: 2026.0")
	println("by the Atom3 contributors.")
}
