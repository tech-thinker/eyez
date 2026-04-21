package main

import (
	"fmt"
	"os"

	"github.com/tech-thinker/eyez/cmd"
	"github.com/urfave/cli/v2"
)

var (
	AppVersion = "v0.0.0"
	CommitHash = "unknown"
	BuildDate  = "unknown"
)

func main() {
	var width int64

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("Version: %s\n", AppVersion)
		fmt.Printf("Commit: %s\n", CommitHash)
		fmt.Printf("Build Date: %s\n", BuildDate)
	}

	app := &cli.App{
		Name:        "eyez",
		Version:     AppVersion,
		Description: "eyez is the powerful CLI tool for viewing image in the terminal (Mostly any terminal).",
		Usage:       "eyez -w <width> <file-name>",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:        "width",
				Aliases:     []string{"w"},
				Destination: &width,
			},
		},
		Action: func(c *cli.Context) error {
			stat, _ := os.Stdin.Stat()
			isPiped := (stat.Mode() & os.ModeCharDevice) == 0
			if isPiped {
				cmd.Pipe(os.Stdin, width)
			} else {
				if c.Args().Len() > 0 {
					cmd.CommandArgs(c.Args().First(), width)
				} else {
					fmt.Println("Usage: eyez -w <width> <file-name>")
					return nil
				}
			}
			return nil
		},
	}

	app.Run(os.Args)
}
