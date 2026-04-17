package utils

import (
	"image"

	"golang.org/x/image/draw"
)

func Resize(src image.Image, newW int) image.Image {
	ration := float64(src.Bounds().Max.X) / float64(src.Bounds().Max.Y)
	newH := int(float64(newW)/ration + 0.5)

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))

	// high-quality scaling
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

	return dst
}
