package graphics

import (
	"fmt"
	"image"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

type Unicode struct {
}

func (*Unicode) clamp(v float64) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return int(v)
}

func (*Unicode) brightness(r, g, b uint32) float64 {
	rf := float64(r >> 8)
	gf := float64(g >> 8)
	bf := float64(b >> 8)

	return 0.299*rf + 0.587*gf + 0.114*bf
}

func (*Unicode) Draw(img image.Image) error {
	rgba := img.(*image.RGBA)
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
	return nil
}
