package cmd

import (
	"image"
	"os"
	"strconv"

	"github.com/tech-thinker/eyez/consts"
	"github.com/tech-thinker/eyez/renderer"
	"github.com/tech-thinker/eyez/utils"
)

func Pipe(in *os.File) {
	var width int64 = consts.DEFAULT_WIDTH

	img, _, err := image.Decode(os.Stdin)
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 1 {
		width, _ = strconv.ParseInt(os.Args[1], 10, 64)
	}

	img = utils.Resize(img, int(width))
	renderer.Render(img)
}
