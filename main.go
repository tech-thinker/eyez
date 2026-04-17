package main

import (
	"os"

	"github.com/tech-thinker/eyez/cmd"
)

func main() {

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	isPiped := (fi.Mode() & os.ModeCharDevice) == 0
	if isPiped {
		cmd.Pipe(os.Stdin)
		return
	}

	cmd.CommandArgs()

}
