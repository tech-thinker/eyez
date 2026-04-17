package main

import (
	"fmt"
	"image"
	"os"
	"strconv"

	"github.com/tech-thinker/eyez/renderer"
	"github.com/tech-thinker/eyez/utils"
	"github.com/tech-thinker/eyez/validator"
)

func main() {

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	isPiped := (fi.Mode() & os.ModeCharDevice) == 0

	var img image.Image

	var width int64 = 80

	if isPiped {
		img, _, err = image.Decode(os.Stdin)
		if len(os.Args) > 1 {
			width, _ = strconv.ParseInt(os.Args[1], 10, 64)
		}
	} else {
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
		img, _, err = image.Decode(f)
		if err != nil {
			panic(err)
		}

		if len(os.Args) > 2 {
			width, _ = strconv.ParseInt(os.Args[2], 10, 64)
		}
	}

	img = utils.Resize(img, int(width))
	renderer.Render(img)

}
