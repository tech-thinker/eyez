package cmd

import (
	"fmt"
	"image"
	"os"
	"strconv"

	"github.com/tech-thinker/eyez/consts"
	"github.com/tech-thinker/eyez/renderer"
	"github.com/tech-thinker/eyez/utils"
	"github.com/tech-thinker/eyez/validator"
)

func CommandArgs() {
	var width int64 = consts.DEFAULT_WIDTH

	if len(os.Args) < 2 {
		fmt.Printf("Usages: %s <file-name> <width>\n", os.Args[0])
		return
	}
	filename := os.Args[1]

	err := validator.Validate(filename)
	if err != nil {
		fmt.Println(err)
	}

	f, _ := os.Open(filename)
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 2 {
		width, _ = strconv.ParseInt(os.Args[2], 10, 64)
	}

	img = utils.Resize(img, int(width))
	renderer.Render(img)
}
