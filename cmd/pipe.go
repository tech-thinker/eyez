package cmd

import (
	"image"
	"os"

	"github.com/tech-thinker/eyez/consts"
	"github.com/tech-thinker/eyez/renderer"
	"github.com/tech-thinker/eyez/utils"
)

func Pipe(in *os.File, width int64) error {
	if width <= 0 {
		width = consts.DEFAULT_WIDTH
	}

	img, _, err := image.Decode(in)
	if err != nil {
		return err
	}

	img = utils.Resize(img, int(width))
	renderer.Render(img)

	return nil
}
