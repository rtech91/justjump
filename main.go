package main

import (
	"fmt"
	"os"

	"github.com/rtech91/justjump/cmd/justjump/cmd"

	"github.com/ttacon/chalk"
)

func main() {
	if os.Geteuid() == 0 {
		red := chalk.Red.NewStyle()
		fmt.Print(red.Style("JustJump is not designed to be run as root. Please run it as a normal user.\n"))
		os.Exit(1)
	}

	cmd.Execute()
}
