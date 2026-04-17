package main

import (
	"fmt"
	"image"
	"os"
	"strconv"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"golang.org/x/image/draw"
)

func clamp(v float64) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return int(v)
}

func brightness(r, g, b uint32) float64 {
	rf := float64(r >> 8)
	gf := float64(g >> 8)
	bf := float64(b >> 8)

	return 0.299*rf + 0.587*gf + 0.114*bf
}

func resize(src image.Image, newW int) image.Image {
	ration := float64(src.Bounds().Max.X) / float64(src.Bounds().Max.Y)
	newH := int(float64(newW)/ration + 0.5)

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))

	// high-quality scaling
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

	return dst
}

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

		if !(strings.HasSuffix(filename, ".png") || strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") || strings.HasSuffix(filename, ".gif") || strings.HasSuffix(filename, ".bmp") || strings.HasSuffix(filename, ".webp")) {
			fmt.Printf("Unsupported file type: %s\n", filename)
			return
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

	rgba := resize(img, int(width)).(*image.RGBA)
	b := rgba.Bounds()
	pix := rgba.Pix
	stride := rgba.Stride

	for y := 0; y < b.Dy(); y += 2 {
		for x := 0; x < b.Dx(); x++ {

			i1 := y*stride + x*4
			i2 := i1 + stride

			// top pixel
			tr, tg, tb := pix[i1], pix[i1+1], pix[i1+2]

			// bottom pixel
			br, bg, bb := uint8(0), uint8(0), uint8(0)
			if y+1 < b.Dy() {
				br, bg, bb = pix[i2], pix[i2+1], pix[i2+2]
			}

			fmt.Printf(
				"\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm▀",
				tr, tg, tb,
				br, bg, bb,
			)
		}
		fmt.Print("\033[0m\n") // reset once per line
	}
}
