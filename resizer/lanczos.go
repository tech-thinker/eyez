package resizer

import (
	"image"
	"image/draw"

	"github.com/disintegration/imaging"
)

type Lanczos struct {
}

func (r *Lanczos) Resize(src image.Image, newW int) image.Image {
	resized := imaging.Resize(src, newW, 0, imaging.Lanczos)
	// Step 2: convert to RGBA
	b := resized.Bounds()
	dst := image.NewRGBA(b)

	draw.Draw(dst, b, resized, b.Min, draw.Src)

	return dst
}
