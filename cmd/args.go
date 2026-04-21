package cmd

import (
	"fmt"
	"image"
	"os"

	"github.com/tech-thinker/eyez/consts"
	"github.com/tech-thinker/eyez/renderer"
	"github.com/tech-thinker/eyez/utils"
	"github.com/tech-thinker/eyez/validator"
)

func CommandArgs(filename string, width int64) error {

	if width <= 0 {
		width = consts.DEFAULT_WIDTH
	}

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

	img = utils.Resize(img, int(width))
	renderer.Render(img)
	return nil
}
