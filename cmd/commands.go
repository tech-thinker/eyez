package cmd

import (
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/tech-thinker/eyez/consts"
	"github.com/tech-thinker/eyez/graphics"
	"github.com/tech-thinker/eyez/resizer"
	"github.com/tech-thinker/eyez/validator"
)

type Commands interface {
	ByArgs(filename string, width int64) error
	ByStdin(in *os.File, width int64) error
}

type commands struct {
	g    graphics.Graphics
	algo resizer.Resizer
}

func (c *commands) ByArgs(filename string, width int64) error {
	err := validator.Validate(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}

	f, _ := os.Open(filename)
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}

	img = c.algo.Resize(img, int(width))
	c.g.Draw(img)
	return nil
}

func (c *commands) ByStdin(in *os.File, width int64) error {
	img, _, err := image.Decode(in)
	if err != nil {
		return err
	}

	img = c.algo.Resize(img, int(width))
	c.g.Draw(img)
	return nil
}

func (c *commands) applySettings(g string, a string) {
	if strings.EqualFold(g, consts.GRAPHICS_UNICODE) {
		c.g = new(graphics.Unicode)
	} else if strings.EqualFold(g, consts.GRAPHICS_ASCII) {
		c.g = new(graphics.ASCII)
	} else {
		fmt.Printf("Graphics must be in [%s|%s].\n", consts.GRAPHICS_UNICODE, consts.GRAPHICS_ASCII)
		os.Exit(-1)
	}

	if strings.EqualFold(a, consts.ALGO_CATMULL_ROM) {
		c.algo = new(resizer.CatmullRom)
	} else if strings.EqualFold(a, consts.ALGO_LANCZOS) {
		c.algo = new(resizer.Lanczos)
	} else {
		fmt.Printf("Algorithm must be in [%s|%s].\n", consts.ALGO_CATMULL_ROM, consts.ALGO_LANCZOS)
		os.Exit(-1)
	}
}

func NewCommands(g string, a string) Commands {
	c := &commands{}
	c.applySettings(g, a)
	return c
}
