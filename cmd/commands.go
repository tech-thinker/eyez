package cmd

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"os"
	"strings"

	"github.com/tech-thinker/eyez/consts"
	"github.com/tech-thinker/eyez/graphics"
	"github.com/tech-thinker/eyez/resizer"
	"github.com/tech-thinker/eyez/svg"
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

func decodeImage(r io.Reader, filename string) (image.Image, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(strings.ToLower(filename), ".svg") {
		svgImg, err := svg.Decode(bytes.NewReader(data))
		if err == nil {
			return svgImg, nil
		}
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err == nil {
		return img, nil
	}

	svgImg, svgErr := svg.Decode(bytes.NewReader(data))
	if svgErr == nil {
		return svgImg, nil
	}

	return nil, err
}

func (c *commands) ByArgs(filename string, width int64) error {
	err := validator.Validate(filename)
	if err != nil {
		return err
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	img, err := decodeImage(f, filename)
	if err != nil {
		return err
	}

	img = c.algo.Resize(img, int(width))
	return c.g.Draw(img)
}

func (c *commands) ByStdin(in *os.File, width int64) error {
	img, err := decodeImage(in, "")
	if err != nil {
		return err
	}

	img = c.algo.Resize(img, int(width))
	return c.g.Draw(img)
}

func (c *commands) applySettings(g string, a string) {
	if strings.EqualFold(g, consts.GRAPHICS_UNICODE) {
		c.g = new(graphics.Unicode)
	} else if strings.EqualFold(g, consts.GRAPHICS_ASCII) {
		c.g = new(graphics.ASCII)
	} else if strings.EqualFold(g, consts.GRAPHICS_KITTY) {
		c.g = new(graphics.Kitty)
	} else {
		fmt.Printf("Graphics must be in [%s, %s, %s].\n", consts.GRAPHICS_UNICODE, consts.GRAPHICS_ASCII, consts.GRAPHICS_KITTY)
		os.Exit(-1)
	}

	if strings.EqualFold(a, consts.ALGO_CATMULL_ROM) {
		c.algo = new(resizer.CatmullRom)
	} else if strings.EqualFold(a, consts.ALGO_LANCZOS) {
		c.algo = new(resizer.Lanczos)
	} else {
		fmt.Printf("Algorithm must be in [%s, %s].\n", consts.ALGO_CATMULL_ROM, consts.ALGO_LANCZOS)
		os.Exit(-1)
	}
}

func NewCommands(g string, a string) Commands {
	c := &commands{}
	c.applySettings(g, a)
	return c
}
