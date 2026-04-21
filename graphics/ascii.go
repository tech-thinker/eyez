package graphics

import (
	"fmt"
	"image"
	"math"
)

type ASCII struct{}

var asciiRamp = "@%#*+=-:. " // dark → light

func luminance(r, g, b uint8) float64 {
	return 0.2126*float64(r) +
		0.7152*float64(g) +
		0.0722*float64(b)
}

func gammaCorrect(v uint8) uint8 {
	f := float64(v) / 255.0
	f = math.Pow(f, 1.0/2.2)
	return uint8(f * 255.0)
}

func (*ASCII) Draw(img image.Image) error {
	nrgba, ok := img.(*image.RGBA)
	if !ok {
		return fmt.Errorf("expected *image.NRGBA")
	}

	b := nrgba.Bounds()
	pix := nrgba.Pix
	stride := nrgba.Stride

	for y := 0; y < b.Dy(); y += 2 { // step 2 → fix terminal aspect ratio
		for x := 0; x < b.Dx(); x++ {

			i := y*stride + x*4

			r := gammaCorrect(pix[i])
			g := gammaCorrect(pix[i+1])
			bb := gammaCorrect(pix[i+2])

			// brightness
			lum := luminance(r, g, bb)

			// map to ascii
			idx := int((lum / 255.0) * float64(len(asciiRamp)-1))

			fmt.Print(string(asciiRamp[idx]))
		}
		fmt.Println()
	}

	return nil
}
