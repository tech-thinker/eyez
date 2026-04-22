package resizer

import (
	"image"

	"golang.org/x/image/draw"
)

type CatmullRom struct {
}

func (r *CatmullRom) Resize(src image.Image, newW int) image.Image {
	ration := float64(src.Bounds().Max.X) / float64(src.Bounds().Max.Y)
	newH := int(float64(newW)/ration + 0.5)

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))

	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}
