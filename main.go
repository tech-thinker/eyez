package main

import (
	"fmt"
	"os"

	"github.com/tech-thinker/eyez/cmd"
	"github.com/tech-thinker/eyez/consts"
	"github.com/urfave/cli/v2"
)

var (
	AppVersion = "v0.0.0"
	CommitHash = "unknown"
	BuildDate  = "unknown"
)

func main() {
	var width int64
	var graphics string
	var algo string

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("Version: %s\n", AppVersion)
		fmt.Printf("Commit: %s\n", CommitHash)
		fmt.Printf("Build Date: %s\n", BuildDate)
	}

	app := &cli.App{
		Name:        "eyez",
		Version:     AppVersion,
		Description: "eyez is the powerful CLI tool for viewing image in the terminal (Mostly any terminal).",
		Usage:       "The image viewer in the terminal",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:        "width",
				Aliases:     []string{"w"},
				Value:       consts.DEFAULT_WIDTH,
				Destination: &width,
			},
			&cli.StringFlag{
				Name:        "graphics",
				Usage:       fmt.Sprintf("graphics type [%s, %s, %s]", consts.GRAPHICS_UNICODE, consts.GRAPHICS_ASCII, consts.GRAPHICS_KITTY),
				Aliases:     []string{"g"},
				Value:       consts.DEFAULT_GRAPHICS,
				Destination: &graphics,
			},
			&cli.StringFlag{
				Name:        "algo",
				Usage:       fmt.Sprintf("algo type [%s, %s]", consts.ALGO_CATMULL_ROM, consts.ALGO_LANCZOS),
				Aliases:     []string{"a"},
				Value:       consts.DEFAULT_ALGORITHM,
				Destination: &algo,
			},
		},
		Action: func(c *cli.Context) error {
			stat, _ := os.Stdin.Stat()
			isPiped := (stat.Mode() & os.ModeCharDevice) == 0
			cmd := cmd.NewCommands(graphics, algo)
			if isPiped {
				cmd.ByStdin(os.Stdin, width)
			} else {
				if c.Args().Len() > 0 {
					cmd.ByArgs(c.Args().First(), width)
				} else {
					fmt.Println("Error: missing required arguments")
					cli.ShowAppHelp(c)
					return nil
				}
			}
			return nil
		},
	}

	app.Run(os.Args)
}
