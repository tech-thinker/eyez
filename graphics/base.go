package graphics

import "image"

type Graphics interface {
	Draw(img image.Image) error
}
